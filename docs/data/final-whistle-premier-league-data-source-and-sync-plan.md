# Final Whistle 英超数据源与同步方案

## 1. 文档目的

本文档收敛 Final Whistle 当前关于比赛和球员数据接入的调研结论，并给出一版可执行的方案设计。

目标是：

- 先以英超为单联赛深做
- 尽量不花钱，优先使用免费数据源
- 在免费额度内尽量做到更实时
- 保持现有本地数据库读模型不变
- 为后续切换供应商预留空间

本文档聚焦：

- 当前项目实际需要哪些比赛和球员数据
- 为什么选择 `football-data.org`
- 为什么采用“本地同步入库 + 自适应轮询”
- 如何设计同步任务、表结构扩展和后续实现顺序

---

## 2. 当前项目需要的数据

结合现有后端模型、DTO 和前端页面，当前系统对比赛和球员数据的需求分成两层。

### 2.1 足球事实数据

这部分是外部数据源需要提供并持续更新的内容：

- 联赛
- 赛季
- 球队
- 球员
- 比赛
- 某场比赛的参赛球员名单或 lineup / roster

当前项目里已经落地的核心字段包括：

#### 比赛

- `competition`
- `season`
- `round`
- `status`
- `kickoff_at`
- `home_team_id`
- `away_team_id`
- `home_score`
- `away_score`
- `venue`

#### 球队

- `name`
- `short_name`
- `slug`
- `logo_url`

#### 球员

- `team_id`
- `name`
- `slug`
- `position`
- `avatar_url`

#### 比赛球员关系

- `match_id`
- `player_id`
- `team_id`

这层数据主要服务：

- 比赛列表
- 比赛详情
- 球队详情
- 球员详情
- Check-in 时的可评分球员 roster

### 2.2 用户生成数据

这部分继续完全保存在本地数据库，不依赖第三方数据源：

- `check_ins`
- `player_ratings`
- `tags`
- 比赛评分、球队评分、球员评分、短评、观看方式、观看时间

原则是：

- 外部数据回答“世界上发生了什么比赛、有哪些球员”
- 本地用户数据回答“用户如何记录和评价这些比赛与球员”

---

## 3. 方案结论

当前阶段的推荐方案如下：

1. 单联赛先做深，先做英超
2. 主数据源选择 `football-data.org`
3. 采用“本地同步入库 + 自适应轮询”
4. 前端和业务接口只读取本地数据库，不直接读取第三方 API
5. 同步范围先覆盖：
   - teams
   - players
   - matches
   - match_players
6. 免费层不足时，再评估是否切换到 `API-Football`

一句话结论：

**先用 `football-data.org` 免费层做英超基础数据同步，本地库作为唯一读模型，在比赛窗口内采用自适应高频轮询。**

---

## 4. 为什么不采用运行时直连第三方 API

不推荐在页面请求或业务接口处理时直接实时打第三方 API，原因如下：

- 当前后端结构已经围绕本地关系库和本地聚合展开
- 比赛详情页和球员详情页都依赖本地数据 join 和聚合
- Check-in 逻辑要求某位球员必须属于当前比赛 roster
- 运行时直连会引入更高的延迟、限流风险和供应商耦合
- 供应商返回结构变化会直接影响业务接口

因此推荐将第三方 API 仅作为采集层使用：

- 同步任务负责从第三方 API 拉取数据
- 标准化后写入本地数据库
- 产品侧只查本地表

这和当前项目结构最匹配，也最容易控制成本和演进风险。

---

## 5. 数据源调研结论

本轮重点看了以下数据源：

- `football-data.org`
- `API-Football`
- `TheSportsDB`
- `Sportmonks`

### 5.1 `football-data.org`

优点：

- 英超在其覆盖范围内
- 免费层即可获取基础比赛、球队、赛季、积分榜等数据
- 文档和数据模型相对清晰
- 适合作为低成本起步的数据主源

缺点：

- 免费层限额较紧
- 比分与状态更新实时性有限
- 更深入的 live 和 lineup 能力需要付费层

判断：

- 适合做英超 MVP 的基础数据主源
- 适合“先低成本跑通，再看是否升级”

### 5.2 `API-Football`

优点：

- 比赛、球员、lineups、livescore 等能力更完整
- 如果愿意付小额月费，覆盖面和实时性更均衡

缺点：

- 免费层日额度很低
- 纯免费条件下很难支持分钟级轮询

判断：

- 是未来升级路径
- 不是当前“尽量不花钱”前提下的首选

### 5.3 `TheSportsDB`

优点：

- 成本更低
- 适合原型验证

缺点：

- 足球垂直能力和字段一致性需要额外验证
- 更适合作为备选而不是主源

### 5.4 `Sportmonks`

优点：

- 数据深度和足球能力都比较完整

缺点：

- 成本偏高
- 不符合当前阶段成本约束

### 5.5 本阶段选择

综合成本、覆盖、迁移成本和当前需求，当前阶段选择：

**`football-data.org` 作为英超接入主源。**

保留原则：

- 设计 provider 抽象，不把业务层绑死在单一供应商响应结构上
- 如果未来免费层无法满足时效和 lineup 需求，再切到 `API-Football`

---

## 6. 推荐架构

### 6.1 分层

建议将数据接入分成四层：

#### A. Provider 层

负责对接第三方 API：

- 拉取 competitions
- 拉取 teams
- 拉取 players
- 拉取 matches
- 拉取 lineups / roster

#### B. Normalizer 层

负责把第三方字段映射成 Final Whistle 自己的标准模型：

- Team
- Player
- Match
- MatchPlayer

#### C. Sync 层

负责：

- 调度同步任务
- 判断轮询频率
- 幂等 upsert
- 记录同步时间和错误

#### D. Product 层

即现有前后端业务层，只读本地数据库：

- 比赛列表和详情
- 球员详情
- Check-in 创建与编辑
- 用户聚合

### 6.2 核心原则

- 第三方 API 不直接暴露给前端
- 本地数据库是唯一读模型
- 用户数据与外部事实数据严格分层
- 同步任务可替换供应商，但不影响产品层查询接口

---

## 7. 数据模型扩展建议

为后续接入真实数据，建议先扩展 `teams`、`players`、`matches` 三张表。

### 7.1 建议新增字段

#### `teams`

- `external_source`
- `external_id`
- `external_updated_at`

#### `players`

- `external_source`
- `external_id`
- `external_updated_at`

#### `matches`

- `external_source`
- `external_id`
- `external_updated_at`

### 7.2 可选新增字段

#### `matches`

- 更细粒度的 `status`
- `last_synced_at`
- `lineup_confirmed_at`

#### `match_players`

- `source_type`
  - 例如 `SQUAD`
  - `LINEUP`
  - `BENCH`
- `external_updated_at`

### 7.3 设计原则

- 本地主键继续使用内部 ID
- 外部 ID 只作为映射和幂等 upsert 的依据
- 业务查询一律使用内部 ID

---

## 8. 英超一期建议同步范围

第一阶段不追求“一切都同步”，而是同步最影响产品主路径的内容。

### 8.1 必须同步

- 英超赛季列表
- 英超球队资料
- 英超球员基础资料
- 英超比赛赛程
- 比赛状态和比分
- 每场比赛 roster 或 lineup

### 8.2 建议同步

- 球场名称
- 比赛轮次文本
- 球队 logo
- 球员头像

### 8.3 暂不强求

- 详细技术统计
- 事件流
- xG
- 热区图
- 新闻和赔率

---

## 9. 自适应轮询策略

考虑当前目标是“尽量免费，同时尽量实时”，推荐采用基于英超赛历的自适应轮询。

### 9.1 远期赛程阶段

适用范围：

- 开赛前 7 天以上

同步频率：

- 每天 1 次到 2 次

目标：

- 拉取新增赛程
- 更新 kick-off 时间变化
- 补齐球队、球员、场地等静态数据

### 9.2 比赛日前窗口

适用范围：

- 开赛前 24 小时到开赛前 90 分钟

同步频率：

- 每 30 到 60 分钟

目标：

- 捕捉时间调整
- 捕捉状态变化
- 预热该场比赛的同步任务

### 9.3 比赛窗口

适用范围：

- 开赛前 90 分钟到终场后 30 分钟

同步频率建议：

- 比赛状态和比分：每 1 到 2 分钟
- lineup / roster：开赛前 60 分钟附近重点轮询

目标：

- 尽快拿到 lineup
- 尽快反映开赛、终场、比分变化

### 9.4 终场后补同步

适用范围：

- 终场后 30 分钟到 6 小时

同步频率：

- 每 15 到 30 分钟

目标：

- 确认最终比分
- 确认最终状态
- 补齐或校正阵容数据

### 9.5 非比赛窗口降频

适用范围：

- 没有英超比赛进行中的日期

同步频率：

- 保持低频

目标：

- 避免浪费免费额度

### 9.6 原则

- 只对“近期比赛”和“比赛日中的活跃比赛”高频轮询
- 只在必要时拉取 lineup
- 免费额度优先保留给比赛窗口

---

## 10. 推荐的同步任务拆分

建议拆成以下几类任务，而不是做一个大而全的同步器。

### 10.1 `sync:competitions`

职责：

- 拉取英超所属 competition 基础信息
- 更新赛季标识

频率：

- 很低频

### 10.2 `sync:teams`

职责：

- 拉取英超球队资料
- upsert 本地 `teams`

频率：

- 每周或每日低频

### 10.3 `sync:players`

职责：

- 按球队或赛季拉取球员基础资料
- upsert 本地 `players`

频率：

- 低频
- 转会窗期间可适当升频

### 10.4 `sync:matches`

职责：

- 拉取赛程、状态、比分、场地
- upsert 本地 `matches`

频率：

- 使用自适应轮询

### 10.5 `sync:match-roster`

职责：

- 拉取某场比赛的 roster / lineup
- upsert 本地 `match_players`

频率：

- 开赛前 60 到 90 分钟升频
- 终场后补同步

### 10.6 `sync:backfill`

职责：

- 一次性补历史赛季或历史比赛

频率：

- 手动执行

---

## 11. Provider 抽象建议

为了避免未来被单一供应商绑定，建议尽早抽出 provider interface。

建议 provider 层至少定义这些能力：

- `ListTeams(season)`
- `ListPlayers(team or season)`
- `ListMatches(season, dateRange)`
- `GetMatch(matchExternalID)`
- `GetMatchRoster(matchExternalID)`

标准化输出应该映射到项目内部结构，而不是把第三方 JSON 直接透传到 repository。

---

## 12. 落库与幂等策略

### 12.1 基本原则

- 所有同步任务都必须是幂等的
- 以 `(external_source, external_id)` 作为 upsert 键
- 本地 `id` 不变

### 12.2 Team

upsert 依据：

- `(external_source, external_id)`

更新字段：

- `name`
- `short_name`
- `logo_url`
- `updated_at`
- `external_updated_at`

### 12.3 Player

upsert 依据：

- `(external_source, external_id)`

更新字段：

- `team_id`
- `name`
- `position`
- `avatar_url`
- `updated_at`
- `external_updated_at`

说明：

- 如果球员发生转会，允许更新当前 `team_id`
- 某场比赛的历史参赛关系由 `match_players` 保留快照

### 12.4 Match

upsert 依据：

- `(external_source, external_id)`

更新字段：

- `status`
- `kickoff_at`
- `home_score`
- `away_score`
- `venue`
- `round`
- `updated_at`
- `external_updated_at`

### 12.5 MatchPlayer

upsert 依据建议：

- `(match_id, player_id)`

说明：

- 如果是按 lineup 全量覆盖更新，需明确是否需要先删除过期 roster
- 如果供应商能区分首发和替补，建议保留 `source_type`

---

## 13. 成本与实时性的平衡

由于当前优先目标是“尽量免费”，我们需要接受几个现实约束：

- 免费层通常不适合全时段高频轮询
- 分钟级同步只能用于小范围活跃比赛窗口
- lineup 和 live 更新可能不如付费源稳定

因此本方案不是追求绝对实时，而是追求：

- 在关键比赛窗口尽量敏捷
- 在非关键时间尽量省额度
- 在现有额度内把“感知上的实时性”做出来

这也是选择自适应轮询而非固定全量高频轮询的原因。

---

## 14. 额度预算与请求策略

### 14.1 `10 req/min` 的计算口径

`football-data.org` 免费层的 `10 requests/minute` 应理解为：

- 一次 HTTP API 请求，通常计为一次 request
- 不是按响应里返回了多少场比赛、多少支球队、多少名球员来计费

因此：

- 请求一次比赛列表接口，即使返回当天全部英超比赛，通常也只算 `1 request`
- 请求一次单场比赛详情接口，算 `1 request`
- 如果逐场请求 10 场比赛详情，则算 `10 requests`

这意味着：

- 批量拉列表接口非常节省额度
- 逐场高频拉详情会迅速消耗额度

### 14.2 推荐的请求原则

为适配免费额度，推荐使用以下原则：

1. 比分和状态优先使用批量列表接口
2. 单场详情接口只在必要时使用
3. 球队和球员基础资料低频同步
4. 历史回填单独跑，不与比赛窗口高频轮询混用

### 14.3 推荐的请求形态

#### 比分和状态同步

优先使用：

- `/competitions/PL/matches?dateFrom=...&dateTo=...`

特点：

- 一次请求可以拿到一个日期范围内的多场英超比赛
- 适合轮询比分、状态、开球时间

#### 单场补充同步

仅在必要时使用：

- `/matches/{id}`

适用场景：

- 某场详情页需要补更细字段
- 某场比赛临近开赛，需要确认更详细状态
- 某场比赛数据出现异常，需要补拉纠正

#### 球队和球员同步

低频使用相关 team / squad 接口。

适用场景：

- 每日或每周同步球队资料
- 每日或每周同步球员资料
- 转会窗期间适度升频

### 14.4 英超比赛日预算示例

以下示例只用于说明请求量级，不代表最终必须固定这样实现。

#### 示例 A：当天比分分钟级轮询

假设当天有英超比赛，采用：

- 每分钟请求一次当日英超比赛列表

则：

- 每分钟消耗约 `1 request`

这远低于 `10 req/min` 上限。

#### 示例 B：当天比分 + 少量补充详情

假设同一分钟内：

- 1 次当日英超比赛列表请求
- 2 次单场详情补拉

则：

- 总计约 `3 req/min`

仍然在免费层可接受范围内。

#### 示例 C：逐场详情硬刷

假设某一分钟里：

- 对 8 场比赛分别拉 8 次单场详情
- 再拉 1 次比赛列表

则：

- 总计约 `9 req/min`

此时已经非常接近上限，几乎没有余量给其他任务。

#### 示例 D：逐场详情 + 其他同步并发

如果在示例 C 的基础上，再叠加：

- 1 次球队同步请求
- 1 次球员同步请求

则：

- 可能达到或超过 `10 req/min`

这种方式不适合作为常态。

### 14.5 结论

对于“英超单联赛先做深”这个范围：

- 如果主要通过批量比赛列表接口同步比分和状态，免费层额度大概率够用
- 如果主要通过逐场详情接口轮询，免费层额度会明显吃紧

因此应优先采用：

- 批量拉比赛列表作为主轮询方式
- 单场详情作为补充方式

### 14.6 实施上的限流建议

为了避免在调度抖动时触发免费层限流，建议同步器内部增加保护：

- 单 provider 全局限流
- 分任务优先级
- 比赛窗口内预留额度

推荐原则：

- 将 `matches` 同步设为最高优先级
- 将 `match-roster` 设为第二优先级
- 将 `teams` 和 `players` 设为低优先级后台任务

### 14.7 一个务实的额度分配建议

对于比赛窗口内的单分钟额度，可先按下列思路控制：

- `matches` 批量同步：预留 `1` 到 `2 req/min`
- 单场补拉：预留 `2` 到 `4 req/min`
- 其余任务：尽量不在比赛窗口执行

这样做的好处是：

- 始终给核心比分同步留足安全余量
- 避免球员或球队同步挤占关键比赛窗口额度

---

## 15. 实施顺序建议

### Phase 1

先补 schema 能力：

- 为 `teams`、`players`、`matches` 增加外部映射字段

### Phase 2

实现 provider 抽象和 `football-data.org` provider：

- 先只覆盖 teams
- players
- matches

### Phase 3

实现一次性 backfill：

- 导入当前英超赛季
- 导入英超球队与球员

### Phase 4

实现定时同步：

- 远期低频
- 比赛日前中频
- 比赛窗口高频

### Phase 5

实现 roster / lineup 同步：

- 让 Check-in 表单只依赖真实比赛 roster

### Phase 6

评估免费层是否满足：

- 比赛状态时效
- roster 获取稳定性
- 请求额度

如果不满足，再评估切换到 `API-Football`。

---

## 16. 当前推荐决策

当前建议直接确认以下决策：

1. 英超作为第一阶段唯一目标联赛
2. `football-data.org` 作为第一阶段主数据源
3. 本地同步入库作为唯一产品读模型方案
4. 自适应轮询作为同步策略
5. schema 先补外部映射字段
6. provider 设计为可替换

---

## 17. 后续 TODO

- 明确 `football-data.org` provider 的字段映射表
- 明确英超赛季 backfill 的边界
- 设计同步任务入口和定时调度方式
- 设计同步失败重试和日志
- 设计 lineup / roster 与 `match_players` 的精确映射策略
- 补 migration 和 model 变更

---

## 18. 参考来源

- `football-data.org` pricing: https://www.football-data.org/pricing
- `football-data.org` coverage: https://www.football-data.org/coverage
- `football-data.org` policies: https://docs.football-data.org/general/v4/policies.html
- `football-data.org` docs: https://docs.football-data.org/general/v4/team.html
- `API-Football` pricing: https://www.api-football.com/pricing
- `TheSportsDB` documentation: https://www.thesportsdb.com/documentation
- `TheSportsDB` API page: https://www.thesportsdb.com/api.php
- `Sportmonks` pricing: https://www.sportmonks.com/football-api/world-plan/
- `Sportmonks` fixtures docs: https://docs.sportmonks.com/v3/tutorials-and-guides/tutorials/livescores-and-fixtures/fixtures
