# Final Whistle 自适应同步引擎落地设计

## 1. 文档目标

本文档不是概念方案，而是按当前仓库现状给出一版可以直接开始实现的技术设计。

本文档回答以下问题：

- 系统如何自动运行
- 系统如何初始化
- 如何判断不同时间段、不同同步频率
- 如何设计配置
- 如何提供手动触发能力
- 如何避免重复执行和超额请求

约束条件：

- 当前后端是单个 Go API 进程入口：[main.go](/Users/sz/Code/final-whistle/backend/cmd/api/main.go)
- 数据库是 PostgreSQL
- 迁移方式是顺序执行 SQL 文件：[main.go](/Users/sz/Code/final-whistle/backend/cmd/migrate/main.go)
- 目前没有现成 job 框架、cron 框架或 worker 进程

第一阶段目标：

- 只支持英超
- 只接 `football-data.org`
- 同步 `teams`、`players`、`matches`、`match_players`
- 支持自动调度和手动触发

---

## 2. 最终运行方式

第一阶段采用单体进程内后台调度。

具体来说：

1. `cmd/api` 启动 HTTP 服务
2. 如果配置 `SYNC_ENABLED=true` 且 `SYNC_AUTO_START=true`
3. API 进程启动一个后台 `SyncEngine`
4. `SyncEngine` 周期性扫描数据库和比赛窗口
5. 根据规则生成应执行的同步任务
6. 任务进入数据库 `sync_jobs`
7. 同进程内 `JobRunner` 轮询 `sync_jobs` 并执行

这意味着系统“自动运行”并不是靠外部 cron，而是：

- API 服务进程启动时自动带起 scheduler + runner

手动触发能力通过两种方式提供：

- 管理 API
- 独立 CLI

这两个入口都不直接调用 provider，而是统一往 `sync_jobs` 写任务，再由 `JobRunner` 执行。

---

## 3. 第一阶段新增组件

### 3.1 新增包结构

建议新增以下目录：

```text
backend/internal/sync/
  provider/
  normalize/
  repository/
  service/
  scheduler/
  runner/
  handler/
```

建议新增以下命令：

```text
backend/cmd/sync/main.go
```

### 3.2 各模块职责

#### `provider`

封装 `football-data.org` HTTP client。

只做：

- 发请求
- 解析响应
- 返回内部 provider DTO

不做数据库写入。

#### `normalize`

把 provider DTO 转成内部 upsert payload。

#### `repository`

负责：

- `sync_jobs` 读写
- `sync_cursors` 读写
- advisory lock
- 业务实体 upsert

#### `scheduler`

负责：

- 定期扫描应同步对象
- 生成 `sync_jobs`

#### `runner`

负责：

- 拉取待执行 job
- 执行 job
- 更新 job 状态

#### `service`

封装单个 job 的实际执行流程。

#### `handler`

提供手动触发和状态查询 API。

---

## 4. 必须新增的数据库表

### 4.1 `sync_jobs`

```sql
CREATE TABLE sync_jobs (
  id BIGSERIAL PRIMARY KEY,
  job_type VARCHAR(50) NOT NULL,
  scope_type VARCHAR(50) NOT NULL,
  scope_key VARCHAR(200) NOT NULL,
  dedupe_key VARCHAR(255) NOT NULL,
  trigger_mode VARCHAR(20) NOT NULL CHECK (trigger_mode IN ('automatic', 'manual')),
  priority INT NOT NULL DEFAULT 100,
  status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'running', 'succeeded', 'failed', 'canceled')),
  scheduled_at TIMESTAMPTZ NOT NULL,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  attempt INT NOT NULL DEFAULT 0,
  max_attempts INT NOT NULL DEFAULT 3,
  payload JSONB NOT NULL DEFAULT '{}'::jsonb,
  last_error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX sync_jobs_dedupe_active_idx
ON sync_jobs (dedupe_key)
WHERE status IN ('pending', 'running');

CREATE INDEX sync_jobs_status_scheduled_idx
ON sync_jobs (status, scheduled_at, priority, id);
```

说明：

- 任何待执行任务都持久化为一条 job
- `dedupe_key` 用于避免重复创建同一任务
- `pending/running` 状态下唯一，任务成功或失败后允许再次创建

### 4.2 `sync_cursors`

```sql
CREATE TABLE sync_cursors (
  id BIGSERIAL PRIMARY KEY,
  provider VARCHAR(50) NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  scope_key VARCHAR(200) NOT NULL,
  last_success_at TIMESTAMPTZ,
  last_attempt_at TIMESTAMPTZ,
  last_error_at TIMESTAMPTZ,
  last_error TEXT,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (provider, resource_type, scope_key)
);
```

说明：

- `sync_jobs` 记录“执行历史”
- `sync_cursors` 记录“当前游标和最近成功时间”

### 4.3 核心实体需补字段

`teams`、`players`、`matches` 都增加：

```sql
external_source VARCHAR(50),
external_id VARCHAR(100),
external_updated_at TIMESTAMPTZ
```

索引：

```sql
CREATE UNIQUE INDEX teams_external_unique_idx
ON teams (external_source, external_id)
WHERE external_source IS NOT NULL AND external_id IS NOT NULL;

CREATE UNIQUE INDEX players_external_unique_idx
ON players (external_source, external_id)
WHERE external_source IS NOT NULL AND external_id IS NOT NULL;

CREATE UNIQUE INDEX matches_external_unique_idx
ON matches (external_source, external_id)
WHERE external_source IS NOT NULL AND external_id IS NOT NULL;
```

---

## 5. 配置设计

当前配置入口在 [config.go](/Users/sz/Code/final-whistle/backend/internal/config/config.go)，第一阶段直接扩展它，不引入新配置系统。

### 5.1 新增配置结构

```go
type SyncConfig struct {
    Enabled               bool
    AutoStart             bool
    Role                  string
    Provider              string
    CompetitionCode       string
    ScanIntervalSeconds   int
    AcquireIntervalSeconds int
    MaxWorkers            int
    SafeRateLimitPerMinute int
    MatchLookbackHours    int
    MatchLookaheadDays    int
    WindowFarMatchDays    int
    WindowPreMatchMinutes int
    WindowLiveAfterKickoffMinutes int
    WindowPostMatchMinutes int
    ScheduleFarMatchEveryMinutes int
    SchedulePreMatchEveryMinutes int
    ScheduleLiveEveryMinutes int
    SchedulePostMatchEveryMinutes int
    RosterWindowBeforeKickoffMinutes int
    RosterWindowAfterKickoffMinutes int
    RosterScheduleEveryMinutes int
    TeamSyncHours         int
    PlayerSyncHours       int
    AdminToken            string
}
```

### 5.2 环境变量

```bash
SYNC_ENABLED=true
SYNC_AUTO_START=true
SYNC_ROLE=all
SYNC_PROVIDER=football-data
SYNC_COMPETITION_CODE=PL
SYNC_SCAN_INTERVAL_SECONDS=60
SYNC_ACQUIRE_INTERVAL_SECONDS=15
SYNC_MAX_WORKERS=1
SYNC_SAFE_RATE_LIMIT_PER_MINUTE=8
SYNC_MATCH_LOOKBACK_HOURS=24
SYNC_MATCH_LOOKAHEAD_DAYS=7
SYNC_WINDOW_FAR_MATCH_DAYS=7
SYNC_WINDOW_PRE_MATCH_MINUTES=90
SYNC_WINDOW_LIVE_AFTER_KICKOFF_MINUTES=180
SYNC_WINDOW_POST_MATCH_MINUTES=360
SYNC_SCHEDULE_FAR_MATCH_EVERY_MINUTES=2880
SYNC_SCHEDULE_PRE_MATCH_EVERY_MINUTES=30
SYNC_SCHEDULE_LIVE_EVERY_MINUTES=2
SYNC_SCHEDULE_POST_MATCH_EVERY_MINUTES=20
SYNC_ROSTER_WINDOW_BEFORE_KICKOFF_MINUTES=75
SYNC_ROSTER_WINDOW_AFTER_KICKOFF_MINUTES=120
SYNC_ROSTER_SCHEDULE_EVERY_MINUTES=10
SYNC_TEAM_SYNC_HOURS=24
SYNC_PLAYER_SYNC_HOURS=24
SYNC_ADMIN_TOKEN=<secret>
FOOTBALL_DATA_API_TOKEN=<secret>
```

### 5.3 配置的实际含义

- `SYNC_SCAN_INTERVAL_SECONDS`
  - scheduler 多久扫描一次是否需要创建新 job
- `SYNC_ACQUIRE_INTERVAL_SECONDS`
  - runner 多久尝试获取待执行 job
- `SYNC_MAX_WORKERS`
  - 同时执行多少个 job，第一阶段固定 `1`
- `SYNC_SAFE_RATE_LIMIT_PER_MINUTE`
  - 内部自限速，建议 `8`，不要吃满官方 `10`
- `SYNC_MATCH_LOOKBACK_HOURS`
  - 每轮扫描时回看最近多少小时的比赛
- `SYNC_MATCH_LOOKAHEAD_DAYS`
  - 每轮扫描时向前看未来多少天的比赛
- `SYNC_WINDOW_*`
  - 定义比赛生命周期窗口
- `SYNC_SCHEDULE_*`
  - 定义每个窗口的最小同步间隔
- `SYNC_ROSTER_*`
  - 定义 roster 的关注窗口和同步间隔
- `SYNC_ROLE`
  - 定义当前进程承担的角色：`all/api/scheduler/runner`

这些配置不是装饰性的，下面的调度逻辑会直接用到。

### 5.4 为什么窗口和频率都必须配置化

同步频率一定会随着阶段变化而变化，不能写死在代码里。

例如：

- `far_match` 在开发期可以 2 天同步一次
- 上线后可能改成 1 天一次
- 如果免费额度紧张，`post_match` 可以从 20 分钟放宽到 60 分钟

因此必须把下面这些都做成配置：

- 看多远的比赛
- 比赛窗口如何划分
- 每个窗口多久同步一次
- roster 在什么时间段尝试、多久重试一次

### 5.5 推荐默认值

| 配置项 | 默认值 | 说明 |
|---|---:|---|
| `SYNC_MATCH_LOOKBACK_HOURS` | `24` | 回看最近 24 小时 |
| `SYNC_MATCH_LOOKAHEAD_DAYS` | `7` | 向前看未来 7 天 |
| `SYNC_WINDOW_FAR_MATCH_DAYS` | `7` | far match 最大关注范围 |
| `SYNC_WINDOW_PRE_MATCH_MINUTES` | `90` | 开赛前 90 分钟进入 pre-match |
| `SYNC_WINDOW_LIVE_AFTER_KICKOFF_MINUTES` | `180` | kickoff 后 3 小时视为 live window |
| `SYNC_WINDOW_POST_MATCH_MINUTES` | `360` | 再补 6 小时 post-match |
| `SYNC_SCHEDULE_FAR_MATCH_EVERY_MINUTES` | `2880` | far match 2 天一次 |
| `SYNC_SCHEDULE_PRE_MATCH_EVERY_MINUTES` | `30` | pre-match 30 分钟一次 |
| `SYNC_SCHEDULE_LIVE_EVERY_MINUTES` | `2` | live 2 分钟一次 |
| `SYNC_SCHEDULE_POST_MATCH_EVERY_MINUTES` | `20` | post-match 20 分钟一次 |
| `SYNC_ROSTER_WINDOW_BEFORE_KICKOFF_MINUTES` | `75` | 开赛前 75 分钟关注 roster |
| `SYNC_ROSTER_WINDOW_AFTER_KICKOFF_MINUTES` | `120` | 开赛后 120 分钟仍允许补拉 |
| `SYNC_ROSTER_SCHEDULE_EVERY_MINUTES` | `10` | roster 10 分钟一次 |
| `SYNC_TEAM_SYNC_HOURS` | `24` | 球队每日同步 |
| `SYNC_PLAYER_SYNC_HOURS` | `24` | 球员每日同步 |

### 5.6 配置校验规则

在 `config.validate()` 里增加这些校验：

- `SYNC_ROLE` 必须是 `all/api/scheduler/runner`
- `SYNC_SAFE_RATE_LIMIT_PER_MINUTE` 必须在 `1..10`
- 所有 `*_EVERY_MINUTES` 必须大于 0
- `SYNC_MATCH_LOOKAHEAD_DAYS` 必须大于 0
- `SYNC_WINDOW_FAR_MATCH_DAYS` 不能大于 `SYNC_MATCH_LOOKAHEAD_DAYS`
- `SYNC_WINDOW_PRE_MATCH_MINUTES` 不能小于 `SYNC_ROSTER_WINDOW_BEFORE_KICKOFF_MINUTES`

### 5.7 运行时策略对象

建议在配置加载后组装成一个明确的策略对象，scheduler 只依赖这个对象：

```go
type MatchSchedulePolicy struct {
    Lookback            time.Duration
    Lookahead           time.Duration
    FarMatchWindow      time.Duration
    PreMatchWindow      time.Duration
    LiveWindow          time.Duration
    PostMatchWindow     time.Duration
    FarMatchEvery       time.Duration
    PreMatchEvery       time.Duration
    LiveEvery           time.Duration
    PostMatchEvery      time.Duration
    RosterBeforeKickoff time.Duration
    RosterAfterKickoff  time.Duration
    RosterEvery         time.Duration
}
```

---

## 6. 系统初始化流程

### 6.1 API 进程启动顺序

对 [main.go](/Users/sz/Code/final-whistle/backend/cmd/api/main.go) 的改造顺序如下：

1. `config.Load()`
2. `db.NewConnection()`
3. 初始化 provider client
4. 初始化 sync repositories
5. 初始化 sync services
6. 初始化 scheduler
7. 初始化 job runner
8. 注册管理 API 路由
9. 如果 `SYNC_ENABLED && SYNC_AUTO_START`
10. `go scheduler.Run(ctx)`
11. `go runner.Run(ctx)`
12. 启动 HTTP server

### 6.2 初始化伪代码

```go
cfg, _ := config.Load()
database, _ := db.NewConnection(cfg.Database.URL)

footballProvider := provider.NewFootballDataClient(cfg.Sync, httpClient)
syncRepo := syncrepository.New(database.DB)
syncService := syncservice.New(syncRepo, footballProvider, database.DB, cfg.Sync)
scheduler := syncscheduler.New(syncRepo, database.DB, cfg.Sync)
runner := syncrunner.New(syncRepo, syncService, cfg.Sync)

if cfg.Sync.Enabled && cfg.Sync.AutoStart {
    go scheduler.Run(ctx)
    go runner.Run(ctx)
}
```

### 6.3 首次启动时系统会不会自动有任务

会，但只会生成“当前窗口需要的任务”。

首次启动时 scheduler 做两件事：

1. 检查英超未来 7 天和最近 1 天的比赛
2. 基于窗口规则创建缺失 job

注意：

- scheduler 不直接调用第三方 API
- scheduler 只依赖本地数据库中的比赛表来判断任务

因此第一次真正接外部源前，需要先提供一个“初始化同步入口”。

---

## 7. 系统初始化与首轮数据导入

这是前一版最缺的部分，这里明确下来。

### 7.1 为什么自动调度前必须先初始化

自适应调度依赖本地 `matches` 表的 kickoff 时间。

如果本地还没有英超比赛数据，scheduler 无法判断：

- 哪天有比赛
- 哪些比赛处于 pre-match/live/post-match

所以系统第一天上线时，必须先执行一次初始化导入。

### 7.2 初始化命令

新增 CLI：

```bash
go run ./cmd/sync/main.go bootstrap --competition PL
```

该命令职责：

1. 拉英超球队
2. 拉英超球员
3. 拉当前赛季英超比赛
4. 对未来 7 天和最近 1 天的比赛补建基础 `sync_cursors`
5. 结束

### 7.3 Bootstrap 与正常调度的边界

- `bootstrap` 是一次性初始化
- 自动调度是长期增量更新

自动模式不负责“猜测系统从零开始该怎么导入”，这个动作必须明确触发。

### 7.4 首次上线流程

第一阶段推荐上线流程：

1. 执行数据库迁移
2. 配置 provider token
3. 执行 `bootstrap`
4. 启动 API
5. API 启动后自动进入定时调度

---

## 8. Job 类型定义

第一阶段只定义 5 类 job。

### 8.1 `sync_matches_range`

用途：

- 同步一个日期范围内的英超比赛列表

payload 示例：

```json
{
  "competitionCode": "PL",
  "dateFrom": "2026-03-27",
  "dateTo": "2026-03-27"
}
```

为什么按范围而不是按单场：

- 最省 `football-data.org` 请求额度
- 一次请求返回多场比赛

### 8.2 `sync_match_detail`

用途：

- 补拉单场比赛更细字段

第一阶段不是主路径，只作为手动修复或特殊补偿。

### 8.3 `sync_match_roster`

用途：

- 同步某场比赛 roster / lineup

payload 示例：

```json
{
  "matchId": 123,
  "externalMatchId": "200123",
  "competitionCode": "PL"
}
```

### 8.4 `sync_teams`

用途：

- 同步英超球队资料

### 8.5 `sync_players`

用途：

- 同步英超球队对应球员资料

---

## 9. 自动调度算法

### 9.1 自动调度器做什么

调度器只做三件事：

1. 扫描“应该关心的比赛”
2. 计算这些比赛是否到了该同步的时间
3. 为需要同步的 scope 创建 job

### 9.2 自动调度器不做什么

调度器不做：

- 调 provider
- 写业务数据
- 并发跑很多任务

这些都交给 runner 和 service。

### 9.3 每轮扫描的固定步骤

每 `SYNC_SCAN_INTERVAL_SECONDS` 执行一次：

1. `ensureStaticJobs()`
2. `loadRelevantMatches()`
3. `scheduleMatchRangeJobs()`
4. `scheduleRosterJobs()`

### 9.4 `ensureStaticJobs()`

职责：

- 判断是否需要低频同步 `teams`
- 判断是否需要低频同步 `players`

实现逻辑：

- 查 `sync_cursors`
- 如果 `teams` 距离上次成功超过 `SYNC_TEAM_SYNC_HOURS`
  - 创建 `sync_teams` job
- 如果 `players` 距离上次成功超过 `SYNC_PLAYER_SYNC_HOURS`
  - 创建 `sync_players` job

### 9.5 `loadRelevantMatches()`

只查本地数据库，不查 provider。

SQL 条件建议：

- `competition = 'Premier League'` 或基于外部 code 映射
- `kickoff_at BETWEEN now() - lookback AND now() + lookahead`

为什么是这个窗口：

- lookback 用于 post-match 补同步
- lookahead 用于 pre-match 和 far-match
- 两者都来自配置，不写死

### 9.6 `scheduleMatchRangeJobs()`

这里是关键。

原则不是“为每场比赛建一个同步 job”，而是：

- 把需要同步的比赛按日期归组
- 为每个日期创建一个 `sync_matches_range` job

这样一个比赛日只需要 1 个 job，就可以同步当天全部英超比赛。

#### 具体算法

对 `loadRelevantMatches()` 返回的所有比赛：

1. 计算每场比赛当前窗口类型
2. 计算这场比赛要求的最小同步频率
3. 归组到 `YYYY-MM-DD`
4. 对同一天所有比赛取“最激进的频率”
5. 根据对应日期的 `sync_cursor.last_success_at` 判断是否该创建 job

#### 窗口与频率映射

假设当前配置：

- `far_future`: `2880 分钟`
- `pre_match`: `30 分钟`
- `live_window`: `2 分钟`
- `post_match`: `20 分钟`

则：

- 如果某一天里有一场比赛处于 live_window
  - 当天日期 job 的目标频率就是 `2 分钟`
- 如果当天没有 live，但有比赛处于 pre_match
  - 频率是 `30 分钟`
- 如果都是 far_future
  - 频率是 `720 分钟`

#### 判定条件

假设 scope_key 为：

- `matches_range:PL:2026-03-27`

如果：

- `last_success_at IS NULL`
  - 立即创建 job
- `now - last_success_at >= target_interval`
  - 创建 job
- 否则
  - 不创建

### 9.7 `scheduleRosterJobs()`

roster 不按日期批量，而是按比赛创建 job。

对每场比赛判断：

- 仅在 `kickoff_at - SYNC_ROSTER_WINDOW_BEFORE_KICKOFF_MINUTES` 到 `kickoff_at + SYNC_ROSTER_WINDOW_AFTER_KICKOFF_MINUTES` 内尝试
- 或者比赛结束后 roster 仍为空，允许补一次

第一阶段规则：

- 如果本地 `match_players` 已存在且数量大于 0
  - 默认不再重复建 roster job
- 如果本地 `match_players` 为空
  - 且进入 roster 窗口
  - 按 `10 分钟` 一次频率调度

scope_key：

- `match_roster:<matchID>`

去重键：

- `sync_match_roster:<matchID>`

### 9.8 时间窗口计算函数

建议写成显式函数：

```go
func classifyMatchWindow(now, kickoff time.Time, cfg SyncConfig) MatchWindow
```

返回：

- `far_future`
- `pre_match`
- `live_window`
- `post_match`
- `inactive`

计算规则：

```text
if kickoff - now > WindowFarMatchDays * 24h:
    inactive
else if now < kickoff - WindowPreMatchMinutes:
    far_future
else if now >= kickoff - WindowPreMatchMinutes && now < kickoff:
    pre_match
else if now >= kickoff && now < kickoff + WindowLiveAfterKickoffMinutes:
    live_window
else if now >= kickoff + WindowLiveAfterKickoffMinutes &&
        now < kickoff + WindowLiveAfterKickoffMinutes + WindowPostMatchMinutes:
    post_match
else:
    inactive
```

注意：

- 这里的 `LiveWindowMinutesAfterKickoff` 不依赖外部 provider 是否告诉你 `LIVE`
- 第一阶段直接基于 kickoff 时间推断，足够落地

---

## 10. 如何确定“什么时间同步”

这是同步引擎里最核心的判定逻辑，需要单独明确。

系统不是靠固定 cron 表达式直接执行同步，而是靠“高频扫描 + 条件判定”决定是否创建 job。

### 10.1 判定分成两层

第一层：系统多久检查一次

- 由 `SYNC_SCAN_INTERVAL_SECONDS` 控制
- 例如每 `60 秒` 扫描一次

这不是同步频率，只是“检查频率”。

第二层：当前这个 scope 是否到了该同步的时间

- 由“窗口类型 + 上次成功同步时间”决定

真正控制同步时机的是第二层。

### 10.2 判定所需输入

对于 `sync_matches_range`，每次判断需要以下输入：

- `now`
- 某一天英超比赛列表
- 这些比赛各自的 `kickoff_at`
- 本地 `sync_cursors.last_success_at`

对于 `sync_match_roster`，每次判断需要以下输入：

- `now`
- 单场比赛的 `kickoff_at`
- 本地该场 `match_players` 是否已有数据
- 本地 `sync_cursors.last_success_at`

### 10.3 比赛列表同步的判定单位

比赛列表同步不按“每场比赛一个 job”判定，而是按“某一天的比赛范围”判定。

scope 设计为：

- `matches_range:PL:2026-03-27`

原因：

- `football-data.org` 的比赛列表接口一次可以返回某个日期范围的多场比赛
- 按日期归组最省请求额度

### 10.4 单场比赛窗口分类

先对每场比赛做窗口分类。

输入：

- `now`
- `kickoff_at`
- sync 配置

输出：

- `far_future`
- `pre_match`
- `live_window`
- `post_match`
- `inactive`

建议规则：

#### `far_future`

条件：

- 比赛在未来 `FarFutureDays` 内
- 但还没进入赛前窗口

也就是：

- `kickoff_at - now > PreMatchWindowMinutes`
- 且 `kickoff_at - now <= FarFutureDays * 24h`

#### `pre_match`

条件：

- `kickoff_at - now <= PreMatchWindowMinutes`
- 且 `now < kickoff_at`

#### `live_window`

条件：

- `now >= kickoff_at`
- 且 `now < kickoff_at + LiveWindowMinutesAfterKickoff`

#### `post_match`

条件：

- `now >= kickoff_at + LiveWindowMinutesAfterKickoff`
- 且 `now < kickoff_at + LiveWindowMinutesAfterKickoff + PostMatchWindowMinutes`

#### `inactive`

条件：

- 超过 `FarFutureDays` 还很远
- 或已经超过 `post_match` 补同步窗口

### 10.5 每个窗口对应的目标频率

每个窗口不是直接触发同步，而是映射成“目标同步间隔”。

例如：

- `far_future` -> `SYNC_FAR_FUTURE_SYNC_MINUTES`
- `pre_match` -> `SYNC_PRE_MATCH_SYNC_MINUTES`
- `live_window` -> `SYNC_LIVE_SYNC_MINUTES`
- `post_match` -> `SYNC_POST_MATCH_SYNC_MINUTES`

假设配置为：

- `far_future = 2880 分钟`
- `pre_match = 30 分钟`
- `live_window = 2 分钟`
- `post_match = 20 分钟`

则说明：

- far future 比赛一天只需要看几次
- pre-match 比赛半小时看一次
- live_window 比赛两分钟看一次
- post-match 比赛 20 分钟补一次

### 10.6 一天内多场比赛如何决定频率

某个日期可能有多场比赛，每场比赛窗口可能不同。

例如 `2026-03-27`：

- 比赛 A 处于 `live_window`
- 比赛 B 处于 `pre_match`
- 比赛 C 处于 `far_future`

这一天的日期级同步频率应取“最激进”的那个窗口。

即：

- 取所有比赛目标频率中的最小值

在上面的例子里：

- `min(2m, 30m, 720m) = 2m`

所以 `matches_range:PL:2026-03-27` 这个 scope 的目标同步频率是 `2 分钟`

### 10.7 真正决定是否创建 job 的条件

对某个 scope，比如：

- `matches_range:PL:2026-03-27`

系统读取：

- `sync_cursors.last_success_at`

然后判断：

#### 规则 1：从未成功同步过

如果：

- `last_success_at IS NULL`

则：

- 立即创建 job

#### 规则 2：已达到最小同步间隔

如果：

- `now - last_success_at >= target_interval`

则：

- 创建 job

#### 规则 3：未达到间隔

如果：

- `now - last_success_at < target_interval`

则：

- 不创建 job

### 10.8 为什么用“上次成功时间”而不是“上次尝试时间”

因为第一阶段我们的目标是“按成功更新的节奏决定下一次同步”。

如果用 `last_attempt_at`，会出现问题：

- 一次失败后，系统会误以为刚同步过
- 导致该 scope 在关键窗口里长时间不再尝试

因此调度判定的主依据必须是：

- `last_success_at`

而失败重试交给 job runner 处理。

### 10.9 Roster 的判定逻辑

roster 不按日期范围判定，而是按比赛判定。

scope：

- `match_roster:<matchID>`

判定条件：

1. 当前时间处于 roster 关注窗口
2. 本地 `match_players` 为空，或者明确要求重拉
3. `last_success_at` 未达到 roster 最小间隔

roster 的关注窗口和频率也来自配置：

- `kickoff_at - SYNC_ROSTER_WINDOW_BEFORE_KICKOFF_MINUTES`
- 到 `kickoff_at + SYNC_ROSTER_WINDOW_AFTER_KICKOFF_MINUTES`
- 最小同步间隔是 `SYNC_ROSTER_SCHEDULE_EVERY_MINUTES`

因此：

- 如果比赛不在这个窗口内，不建 roster job
- 如果比赛在窗口内，但本地已有 roster，默认不建 job
- 如果比赛在窗口内且本地 roster 为空，再根据 `last_success_at` 判定是否该建 job

### 10.10 举例

#### 例子 A：普通比赛日中午

当前时间：

- `2026-03-27 12:00`

当天一场比赛 `20:00` 开赛。

则：

- 距离开赛还有 8 小时
- 若配置 `PreMatchWindowMinutes=90`
- 这场比赛属于 `far_future`

假设：

- `far_future` 频率是 `720 分钟`
- 当天该 scope 上次成功时间是今天 `08:00`

则：

- `12:00 - 08:00 = 4 小时 < 12 小时`
- 不创建 job

#### 例子 B：开赛前 45 分钟

当前时间：

- `2026-03-27 19:15`

比赛 `20:00` 开赛。

则：

- 属于 `pre_match`
- 目标频率为 `30 分钟`

假设：

- 上次成功时间是 `18:30`

则：

- `19:15 - 18:30 = 45 分钟`
- 已超过 `30 分钟`
- 创建 `sync_matches_range` job

#### 例子 C：比赛进行中

当前时间：

- `2026-03-27 21:08`

比赛 `20:00` 开赛。

则：

- 属于 `live_window`
- 目标频率为 `2 分钟`

假设：

- 上次成功时间是 `21:06`

则：

- 已达到 `2 分钟`
- 创建 job

#### 例子 D：roster

当前时间：

- `2026-03-27 19:10`

比赛 `20:00` 开赛。

则：

- 进入 `kickoff - 75m` 的 roster 窗口

假设：

- 本地 `match_players` 为空
- 上次成功同步 roster 是空值

则：

- 立即创建 `sync_match_roster` job

如果：

- 本地 `match_players` 已经有 22 条数据

则：

- 默认不再创建 roster job

### 10.11 调度器最终判断伪代码

#### 比赛范围同步

```go
for _, dateGroup := range groupedMatchesByDate {
    targetInterval := minInterval(dateGroup.matches, now, cfg)
    if targetInterval == 0 {
        continue
    }

    cursor := repo.GetCursor("football-data", "matches_range", dateGroup.scopeKey)
    if cursor.LastSuccessAt == nil || now.Sub(*cursor.LastSuccessAt) >= targetInterval {
        repo.EnqueueJob(...)
    }
}
```

#### Roster 同步

```go
for _, match := range relevantMatches {
    if !isRosterWindow(now, match.KickoffAt) {
        continue
    }
    if repo.MatchRosterExists(match.ID) {
        continue
    }

    cursor := repo.GetCursor("football-data", "match_roster", rosterScopeKey(match.ID))
    if cursor.LastSuccessAt == nil || now.Sub(*cursor.LastSuccessAt) >= 10*time.Minute {
        repo.EnqueueJob(...)
    }
}
```

### 10.12 这套机制的本质

总结一下，“什么时间同步”不是一个静态时间表，而是：

1. 系统按固定短周期检查
2. 对每个 scope 计算当前所处窗口
3. 将窗口映射为目标同步频率
4. 用 `last_success_at` 判断是否到点
5. 到点才建 job

因此它本质上是：

- 高频扫描
- 低频或高频执行
- 按比赛生命周期自适应变化

---

## 11. Runner 设计

### 10.1 Runner 的职责

每 `SYNC_ACQUIRE_INTERVAL_SECONDS` 执行一次：

1. 从 `sync_jobs` 取出一个待执行 job
2. 抢占执行权
3. 执行 job
4. 更新状态

### 10.2 为什么第一阶段只开 1 个 worker

因为：

- 免费额度只有 `10 req/min`
- 最核心目标是稳定而不是吞吐
- 先把正确性做对

所以第一阶段：

- `SYNC_MAX_WORKERS=1`

### 10.3 获取 job 的 SQL 策略

建议使用事务 + `FOR UPDATE SKIP LOCKED`。

伪 SQL：

```sql
SELECT id
FROM sync_jobs
WHERE status = 'pending'
  AND scheduled_at <= NOW()
ORDER BY priority ASC, scheduled_at ASC, id ASC
FOR UPDATE SKIP LOCKED
LIMIT 1;
```

拿到后立即更新为：

- `status = 'running'`
- `started_at = NOW()`
- `attempt = attempt + 1`

### 10.4 执行失败后的处理

如果失败：

- 如果 `attempt < max_attempts`
  - 设为 `pending`
  - `scheduled_at = NOW() + retryBackoff(attempt)`
- 否则
  - 设为 `failed`
  - 写 `last_error`

### 10.5 重试退避

第一阶段固定：

- 第 1 次失败：5 分钟后
- 第 2 次失败：15 分钟后
- 第 3 次失败：60 分钟后

---

## 12. Service 执行细节

### 11.1 `sync_matches_range`

执行步骤：

1. 对 `scope_key` 获取 advisory lock
2. 调 provider 的 range 接口
3. 解析响应中的 team 和 match
4. 事务内 upsert teams
5. 事务内 upsert matches
6. 更新 `sync_cursors(provider='football-data', resource_type='matches_range', scope_key=...)`
7. job 标记成功

### 11.2 `sync_teams`

执行步骤：

1. lock `teams:PL`
2. 拉取英超球队
3. upsert teams
4. 更新 cursor

### 11.3 `sync_players`

执行步骤：

1. lock `players:PL`
2. 先查本地英超 teams
3. 按 team 逐个拉 squad
4. 按 team 分批 upsert players
5. 更新 cursor

### 11.4 `sync_match_roster`

执行步骤：

1. lock `match_roster:<matchID>`
2. 查本地 match
3. 调 provider 拉 roster
4. 如果 roster 中出现本地不存在的球员
   - 先补 player upsert
5. 事务内：
   - 删除该 match 原有 `match_players`
   - 插入新 `match_players`
6. 更新 cursor

### 11.5 为什么 roster 采用全量替换

因为第一阶段最稳妥，逻辑简单：

- 同步来源是当前快照
- 本地 `match_players` 也应视为快照

如果未来要区分首发/替补/名单来源，再扩表。

---

## 13. Provider 限流落地方案

### 12.1 限流实现位置

限流必须放在 provider client 内，而不是放在 scheduler 外层。

因为：

- 手动触发和自动触发都走 provider client
- 这样任何入口都天然受控

### 12.2 限流算法

直接使用 token bucket。

参数：

- 每分钟补充 `SYNC_SAFE_RATE_LIMIT_PER_MINUTE`
- 建议值 `8`

行为：

- 每次 HTTP 请求消耗 1 token
- 没 token 时阻塞等待或返回“稍后重试”

第一阶段推荐阻塞等待。

### 12.3 为什么用 8 而不是 10

给这些场景留余量：

- 手动触发
- provider 重试
- 时钟抖动
- 未来可能增加额外查询

---

## 14. 发布后的部署设计

同步模块在正式环境不能只考虑“代码怎么跑”，还要考虑“发布后以什么进程形态运行”。

第一阶段建议从一开始就支持角色化部署。

### 14.1 `SYNC_ROLE`

建议增加配置：

- `SYNC_ROLE=all`
- `SYNC_ROLE=api`
- `SYNC_ROLE=scheduler`
- `SYNC_ROLE=runner`

各角色定义如下：

- `all`
  - 启动 HTTP API
  - 启动 scheduler
  - 启动 runner
- `api`
  - 只启动 HTTP API
  - 不启动 scheduler
  - 不启动 runner
- `scheduler`
  - 只启动 scheduler
  - 不启动 HTTP API
  - 不启动 runner
- `runner`
  - 只启动 runner
  - 不启动 HTTP API
  - 不启动 scheduler

### 14.2 开发环境部署模式

开发环境推荐一体化模式：

```bash
SYNC_ENABLED=true
SYNC_AUTO_START=true
SYNC_ROLE=all
```

这样开发时只启动一个进程即可，最省事。

### 14.3 正式环境推荐部署模式

正式环境推荐：

- 多个 API 实例
- 一个专用 Sync Worker 实例

API 实例配置：

```bash
SYNC_ENABLED=true
SYNC_AUTO_START=false
SYNC_ROLE=api
```

Sync Worker 配置：

```bash
SYNC_ENABLED=true
SYNC_AUTO_START=true
SYNC_ROLE=all
SYNC_MAX_WORKERS=1
```

如果后续要拆得更细，也可以：

- 一个 `scheduler` 进程
- 一个 `runner` 进程

但第一阶段没有必要。

### 14.4 为什么生产环境不建议所有 API 实例都带同步

原因：

- 每个 API 实例都跑 scheduler 会产生重复调度
- 后台同步会和在线请求争资源
- 排查问题更困难
- 调整同步参数时不如单独 worker 清晰

虽然系统有：

- `dedupe_key`
- `FOR UPDATE SKIP LOCKED`
- advisory lock

这些机制能降低重复执行风险，但不应该依赖它们来替代正确部署。

### 14.5 推荐的进程入口

第一阶段保留两种启动方式：

#### 方式 A：开发环境

直接用现有 API 入口：

```bash
go run ./cmd/api/main.go
```

当 `SYNC_ROLE=all` 且 `SYNC_AUTO_START=true` 时，API 进程内部启动 scheduler 和 runner。

#### 方式 B：正式环境

新增专用入口：

```bash
go run ./cmd/sync/main.go daemon
```

`daemon` 模式读取同一套配置，并按 `SYNC_ROLE` 决定启动哪些组件。

建议：

- API 服务使用 `cmd/api`
- Sync Worker 使用 `cmd/sync daemon`

### 14.6 `cmd/sync daemon` 的具体行为

当执行：

```bash
go run ./cmd/sync/main.go daemon
```

程序行为：

1. 读取配置
2. 建立数据库连接
3. 初始化 provider、repository、service
4. 根据 `SYNC_ROLE` 启动 scheduler 和/或 runner
5. 阻塞运行直到收到退出信号

因此：

- `cmd/api` 适合在线服务
- `cmd/sync daemon` 适合后台 worker

### 14.7 首次发布流程

正式环境第一次上线建议流程：

1. 执行 migration
2. 部署 API 服务
3. 部署 Sync Worker
4. 配置 `FOOTBALL_DATA_API_TOKEN`
5. 执行 `bootstrap`
6. 检查 `/admin/sync/status`
7. 确认 `sync_jobs` 和 `sync_cursors` 开始流转

### 14.8 配置变更如何生效

第一阶段不做动态热更新。

规则：

- 修改同步配置后，重启对应进程生效

这是当前最稳的方案。

### 14.9 生产环境最少需要哪些配置

正式环境必须显式配置：

- `SYNC_ENABLED`
- `SYNC_AUTO_START`
- `SYNC_ROLE`
- `SYNC_PROVIDER`
- `SYNC_COMPETITION_CODE`
- `SYNC_SAFE_RATE_LIMIT_PER_MINUTE`
- `SYNC_MATCH_LOOKBACK_HOURS`
- `SYNC_MATCH_LOOKAHEAD_DAYS`
- `SYNC_WINDOW_*`
- `SYNC_SCHEDULE_*`
- `SYNC_ROSTER_*`
- `FOOTBALL_DATA_API_TOKEN`
- `SYNC_ADMIN_TOKEN`

---

## 15. 手动触发能力设计

### 13.1 管理 API

新增管理路由：

```text
POST /admin/sync/bootstrap
POST /admin/sync/jobs
POST /admin/sync/jobs/retry-failed
GET  /admin/sync/jobs
GET  /admin/sync/jobs/:id
GET  /admin/sync/status
```

鉴权方式：

- 第一阶段直接使用 `Authorization: Bearer <SYNC_ADMIN_TOKEN>`

这是最简单可落地的做法，不依赖现有用户系统扩展管理员角色。

### 13.2 `POST /admin/sync/jobs`

请求体：

```json
{
  "jobType": "sync_matches_range",
  "scopeType": "date",
  "scopeKey": "matches_range:PL:2026-03-27",
  "payload": {
    "competitionCode": "PL",
    "dateFrom": "2026-03-27",
    "dateTo": "2026-03-27"
  },
  "priority": 10
}
```

后端行为：

- 校验参数
- 计算 dedupe_key
- 插入 `sync_jobs`
- 返回 job id

### 13.3 `POST /admin/sync/bootstrap`

后端行为：

- 直接创建 3 个高优先级 job：
  - `sync_teams`
  - `sync_players`
  - `sync_matches_range` for current season relevant range

注意：

- 这里不在 handler 里直接调 provider
- 仍然只是入队

### 13.4 CLI

新增命令：

```bash
go run ./cmd/sync/main.go bootstrap --competition PL
go run ./cmd/sync/main.go enqueue --job sync_matches_range --date-from 2026-03-27 --date-to 2026-03-27
go run ./cmd/sync/main.go enqueue --job sync_match_roster --match-id 123
go run ./cmd/sync/main.go run-once
```

说明：

- `bootstrap`
  - 创建初始化任务
- `enqueue`
  - 手动创建指定任务
- `run-once`
  - 适合本地调试，执行一个待处理 job 后退出

---

## 16. 自动模式与手动模式如何共存

这是实现上的关键。

### 14.1 自动模式

- scheduler 自动创建 job
- runner 自动执行 job

### 14.2 手动模式

- API 或 CLI 手动创建 job
- 同一个 runner 自动执行

### 14.3 冲突处理

靠三层机制避免冲突：

1. `dedupe_key` 避免重复建同类活动任务
2. `FOR UPDATE SKIP LOCKED` 避免多个 runner 抢同一 job
3. advisory lock 避免同 scope 重复执行

---

## 17. 系统状态接口设计

为了让这个系统可用，必须能回答“现在同步到哪一步了”。

### 15.1 `GET /admin/sync/status`

返回：

```json
{
  "enabled": true,
  "autoStart": true,
  "provider": "football-data",
  "competitionCode": "PL",
  "schedulerRunning": true,
  "runnerRunning": true,
  "pendingJobs": 3,
  "runningJobs": 1,
  "failedJobs": 0,
  "lastMatchesRangeSuccessAt": "2026-03-27T12:01:00Z",
  "lastTeamsSuccessAt": "2026-03-27T00:10:00Z",
  "lastPlayersSuccessAt": "2026-03-27T00:20:00Z"
}
```

### 15.2 `GET /admin/sync/jobs`

支持参数：

- `status`
- `jobType`
- `limit`

用于查看：

- 失败任务
- 正在执行的任务
- 最近任务历史

---

## 18. 当前仓库内的具体代码落点

### 16.1 `backend/internal/config/config.go`

改动：

- 增加 `SyncConfig`
- `setDefaults()` 补默认值
- `bindEnvVars()` 绑定 sync env vars
- `validate()` 校验 provider 和频率参数

### 16.2 `backend/cmd/api/main.go`

改动：

- 初始化 provider / sync repositories / scheduler / runner
- 启动后台 goroutine
- 注入管理路由

### 16.3 新增 `backend/cmd/sync/main.go`

用途：

- bootstrap
- enqueue
- run-once

### 16.4 新增 migration SQL

至少需要：

- 创建 `sync_jobs`
- 创建 `sync_cursors`
- 给 `teams`、`players`、`matches` 增字段与索引

---

## 19. 第一阶段不做的事情

为了保证能落地，以下内容第一阶段明确不做：

- 多实例分布式部署优化
- 独立 worker 服务拆分
- 图形化运维后台
- 高级事件流同步
- 复杂 DAG 依赖调度
- 自定义 cron 表达式配置

这些都不是当前需要的最小可用方案。

---

## 20. 可直接开工的实现顺序

按依赖关系，建议这么落地：

1. migration
   - `sync_jobs`
   - `sync_cursors`
   - 外部映射字段

2. config
   - `SyncConfig`
   - env vars

3. repository
   - enqueue job
   - acquire pending job
   - update status
   - update cursor

4. provider
   - football-data client
   - rate limiter

5. service
   - `sync_matches_range`
   - `sync_teams`
   - `sync_players`
   - `sync_match_roster`

6. runner
   - `Run(ctx)`
   - `runOnce()`

7. scheduler
   - `Run(ctx)`
   - `ensureStaticJobs()`
   - `scheduleMatchRangeJobs()`
   - `scheduleRosterJobs()`

8. admin API
   - enqueue
   - status
   - list jobs

9. CLI
   - bootstrap
   - enqueue
   - run-once

---

## 21. 最终结论

这套设计的落地点非常具体：

- 自动运行依赖 API 进程内的 `scheduler + runner`
- 初始化依赖显式 `bootstrap`
- 不同时间段频率通过“窗口分类 + sync_cursor 判定”实现
- 手动触发通过 API/CLI 入队实现
- 配置通过扩展现有 `config.go` 实现
- 限流放在 provider client 内
- 去重靠 `dedupe_key + advisory lock + SKIP LOCKED`

按这个文档实现，不需要额外基础设施，也不需要先重构整个后端架构。
