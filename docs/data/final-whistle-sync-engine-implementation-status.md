# Final Whistle Sync Engine 实现状态

## 当前分支

- `feat-premier-league-sync-engine`

## 当前状态

本轮已经完成同步引擎第一阶段的基础落地，并达到“可编译、可运行、可测试”的状态。

已完成范围：

- 新增同步基础表：
  - `sync_jobs`
  - `sync_cursors`
- 为核心实体补充外部映射字段：
  - `teams.external_source / external_id / external_updated_at`
  - `players.external_source / external_id / external_updated_at`
  - `matches.external_source / external_id / external_updated_at`
- 扩展同步配置：
  - 窗口
  - 频率
  - 角色
  - provider token
- 新增同步模块基础包：
  - `policy`
  - `provider`
  - `repository`
  - `service`
  - `scheduler`
  - `runner`
  - `handler`
  - `app`
- 新增 `cmd/sync`
  - `daemon`
  - `bootstrap`
  - `run-once`
- 将同步模块接入 API 进程启动路径
- 新增 `/admin/sync` 管理接口
- 实现真实落库能力：
  - `sync_teams`
  - `sync_players`
  - `sync_matches_range`
- 新增隔离测试基建：
  - 每次测试创建独立 PostgreSQL 临时数据库
  - 跑完整 migration
  - 测试结束自动 `DROP DATABASE`
- 新增 fake provider
- 新增同步集成测试与 API 端到端测试

## 已验证内容

已实际验证：

- `go build ./...`
- `go test ./internal/sync -run TestSyncWritesTeamsPlayersAndMatches -v`
- `go test ./internal/router -run TestPublicAPICanReadSyncedData -v`
- `go test ./...`

验证结论：

- 同步任务可以把 fake provider 数据真实写入数据库
- 写入后的数据可以通过公开 API 查询到
- 测试是隔离的，不污染现有数据库

## 本轮仍未完成

以下内容尚未完成：

### 1. `sync_match_roster`

当前状态：

- job 类型和调度骨架已存在
- 实际 provider 落库逻辑未完成

影响：

- `match_players` 还不能通过同步自动填充
- 比赛详情页的 `matchPlayers` 目前不会因为同步而自动出现

### 2. 真实 `football-data.org` 生产接入验证

当前状态：

- provider client 已接入基础 endpoint
- 目前测试使用的是 fake provider

影响：

- 还没有完成一轮针对真实 token 和真实数据的联调验证
- 还没有验证免费层在真实英超赛程上的字段覆盖边界

### 3. 同步服务中的更完整异常处理

当前状态：

- 已有基础失败处理、重试和 cursor 更新
- 还没有针对 provider 429 / 5xx 做更细粒度区分

### 4. 部署级别的运行脚本与运维文档

当前状态：

- 代码已支持 `cmd/api` 和 `cmd/sync`
- 但没有新增部署脚本、systemd/容器示例或 CI 集成

## 下一步建议

建议优先顺序：

1. 完成 `sync_match_roster`
2. 用真实 `football-data.org` token 跑一轮手动联调
3. 补充 roster 相关集成测试和 API 验证
4. 根据真实 provider 表现调优窗口与频率配置
5. 补充部署文档与运行脚本

## 备注

本轮提交聚焦于：

- 建立可运行的同步引擎主骨架
- 打通真实写库
- 打通 API 读路径
- 建立隔离测试基础设施

因此这是一个可继续迭代的“第一阶段可提交版本”，但不是同步模块的最终完成态。
