# Final Whistle Spec Coding 执行规范

## 1. 文档目的

本文档用于以 spec-driven / spec coding 的方式推进 Final Whistle v1 开发。

它不是 PRD 的重复版本，而是面向实施的执行规范，重点定义：

- 明确范围
- 明确需求编号
- 明确接口与数据边界
- 明确验收标准
- 明确不做什么

后续开发、拆解任务、提交代码和验收时，均以本文档为主。

---

## 2. 项目摘要

### 2.1 产品一句话

Final Whistle 是一个面向足球观众的赛后记录产品，帮助用户记录自己看过的比赛，并沉淀成个人观赛档案。

### 2.2 v1 目标

v1 只完成一条闭环：

`登录 -> 找比赛 -> 记录比赛 -> 查看聚合 -> 回看个人档案`

### 2.3 v1 非目标

以下能力明确不在 v1 内：

- 实时比分
- 外部足球数据同步
- 点赞、回复、评论楼中楼
- 社交关系
- 用户自定义标签
- 长文社区
- 复杂后台
- 缓存、消息队列、异步聚合

---

## 3. 固定技术决策

### 3.1 前后端架构

- 前端：Next.js
- 后端：Go + Gin
- 数据库：PostgreSQL

### 3.2 认证

- 使用 `HTTP-only Cookie Session`
- v1 登录方式为 `dev login`
- 不实现 OAuth
- 不实现 JWT 双方案并存

### 3.3 数据来源

- v1 仅使用内部 seed data
- 不接入第三方足球 API

### 3.4 聚合策略

- 全部实时查询
- 不做缓存
- 不做预计算
- 不做异步任务

---

## 4. 范围定义

### 4.1 v1 页面范围

- `/matches`
- `/matches/[matchId]`
- `/teams/[teamId]`
- `/players/[playerId]`
- `/me`
- `/login`

### 4.2 v1 接口范围

- `POST /auth/login`
- `POST /auth/logout`
- `GET /auth/me`
- `GET /matches`
- `GET /matches/:id`
- `GET /matches/:id/my-checkin`
- `POST /matches/:id/checkin`
- `PUT /matches/:id/checkin`
- `GET /teams/:id`
- `GET /players/:id`
- `GET /me/profile`
- `GET /me/checkins`

### 4.3 分页范围

需要分页：

- `GET /matches`
- `GET /me/checkins`

不分页：

- `GET /matches/:id`
- `GET /matches/:id/my-checkin`
- `GET /teams/:id`
- `GET /players/:id`
- `GET /me/profile`

详情页/主页内嵌列表采用限量返回，不单独做分页接口。

---

## 5. 核心实体

### 5.1 User

- `id`
- `name`
- `email`
- `avatarUrl`
- `createdAt`
- `updatedAt`

### 5.2 Team

- `id`
- `name`
- `shortName`
- `slug`
- `logoUrl`
- `externalSource`
- `externalId`
- `createdAt`
- `updatedAt`

### 5.3 Player

- `id`
- `teamId`
- `name`
- `slug`
- `position`
- `avatarUrl`
- `externalSource`
- `externalId`
- `createdAt`
- `updatedAt`

### 5.4 Match

- `id`
- `competition`
- `season`
- `round`
- `status`
- `kickoffAt`
- `homeTeamId`
- `awayTeamId`
- `homeScore`
- `awayScore`
- `venue`
- `externalSource`
- `externalId`
- `createdAt`
- `updatedAt`

### 5.5 MatchPlayer

- `id`
- `matchId`
- `playerId`
- `teamId`

### 5.6 Tag

- `id`
- `name`
- `slug`
- `sortOrder`
- `isActive`

### 5.7 CheckIn

- `id`
- `userId`
- `matchId`
- `watchedType`
- `supporterSide`
- `matchRating`
- `homeTeamRating`
- `awayTeamRating`
- `shortReview`
- `watchedAt`
- `createdAt`
- `updatedAt`

约束：

- 唯一约束：`(user_id, match_id)`

### 5.8 PlayerRating

- `id`
- `checkInId`
- `playerId`
- `rating`
- `note`

### 5.9 CheckInTag

- `id`
- `checkInId`
- `tagId`

### 5.10 Session

- `id`
- `userId`
- `token`
- `expiredAt`
- `createdAt`

---

## 6. 枚举与校验

### 6.1 枚举

`watchedType`

- `FULL`
- `PARTIAL`
- `HIGHLIGHTS`

`supporterSide`

- `HOME`
- `AWAY`
- `NEUTRAL`

`match.status`

- `SCHEDULED`
- `FINISHED`

### 6.2 校验

- `matchRating` 范围 `1-10`
- `homeTeamRating` 范围 `1-10`
- `awayTeamRating` 范围 `1-10`
- `playerRatings[].rating` 范围 `1-10`
- `shortReview` 最大 `280`
- `playerRatings[].note` 最大 `80`
- 单条 check-in 最多 `5` 条球员评分
- 用户只能给当前比赛中的球员打分
- 仅允许对 `FINISHED` 比赛创建或更新 check-in

---

## 7. 功能需求

### 7.1 Auth

#### FW-AUTH-001

系统必须支持用户通过 `dev login` 登录。

验收标准：

- 输入 `email` 和 `name` 后可登录
- 已存在用户复用账号
- 不存在用户自动创建
- 服务端通过 Cookie 建立会话

#### FW-AUTH-002

系统必须支持用户登出。

验收标准：

- 调用登出接口后会话失效
- 后续访问受保护接口返回未登录

#### FW-AUTH-003

系统必须支持获取当前登录用户信息。

验收标准：

- 登录后调用 `/auth/me` 返回当前用户
- 未登录时返回 `401`

### 7.2 Match List

#### FW-MATCH-001

系统必须提供比赛列表页。

验收标准：

- 展示主客队、比分、比赛时间、赛事、社区平均评分、打卡人数
- 支持按 `competition` 筛选
- 支持按 `season` 筛选
- 默认按 `kickoffAt desc` 排序

#### FW-MATCH-002

比赛列表接口必须支持分页。

验收标准：

- 支持 `page`
- 支持 `pageSize`
- 默认 `page=1`
- 默认 `pageSize=20`
- 最大 `pageSize=50`

### 7.3 Match Detail

#### FW-MATCH-003

系统必须提供比赛详情页。

验收标准：

- 展示比赛基础信息
- 展示平均比赛评分
- 展示主队平均评分
- 展示客队平均评分
- 展示球员评分排行
- 展示最近短评
- 展示打卡人数

#### FW-MATCH-004

比赛详情页必须展示当前用户记录状态。

验收标准：

- 未登录显示登录引导
- 已登录未打卡显示创建入口
- 已登录已打卡显示我的记录与编辑入口

### 7.4 CheckIn

#### FW-CHECKIN-001

登录用户必须可以为一场比赛创建唯一一条 check-in。

验收标准：

- 同一用户同一比赛只能创建一次
- 重复创建返回 `CONFLICT`

#### FW-CHECKIN-002

登录用户必须可以更新自己的 check-in。

验收标准：

- 已存在记录可编辑
- 不存在记录时更新返回 `NOT_FOUND`

#### FW-CHECKIN-003

check-in 表单必须支持以下字段：

- `watchedType`
- `supporterSide`
- `matchRating`
- `homeTeamRating`
- `awayTeamRating`
- `playerRatings`
- `tags`
- `shortReview`
- `watchedAt`

#### FW-CHECKIN-004

系统必须校验被评分球员属于当前比赛。

验收标准：

- 非本场比赛球员不可提交成功

#### FW-CHECKIN-005

创建和更新 check-in 时必须使用事务。

验收标准：

- `check_ins`
- `player_ratings`
- `check_in_tags`

任一写入失败时全部回滚。

### 7.5 Team / Player

#### FW-TEAM-001

系统必须提供球队简版详情页。

验收标准：

- 展示球队基本信息
- 展示最近比赛
- 展示平均评分摘要

#### FW-PLAYER-001

系统必须提供球员简版详情页。

验收标准：

- 展示球员基本信息
- 展示最近被评分比赛
- 展示平均评分摘要

### 7.6 Profile

#### FW-PROFILE-001

系统必须提供个人主页。

验收标准：

- 展示用户基础信息
- 展示总打卡数
- 展示平均比赛评分
- 展示最近打卡
- 展示最近短评
- 展示常打卡球队前 3
- 展示常评分球员前 3

#### FW-PROFILE-002

系统必须提供个人历史记录分页列表。

验收标准：

- `GET /me/checkins` 支持分页
- 可从个人主页进入完整历史记录

---

## 8. 非功能要求

### 8.1 API 规范

- 响应格式统一
- 错误码统一
- API 字段使用 `camelCase`

### 8.2 数据一致性

- check-in 的创建和更新必须具备事务原子性

### 8.3 安全

- Cookie 必须启用 `HttpOnly`
- 生产环境开启 `Secure`
- Session 存数据库，不用内存

### 8.4 可观测性

- 提供基础健康检查接口
- 使用结构化日志

---

## 9. API 响应规范

### 9.1 成功响应

```json
{
  "success": true,
  "data": {}
}
```

### 9.2 错误响应

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "invalid request body",
    "details": {}
  }
}
```

### 9.3 分页响应

```json
{
  "success": true,
  "data": {
    "items": [],
    "page": 1,
    "pageSize": 20,
    "total": 100
  }
}
```

---

## 10. 页面与接口映射

### `/login`

- `POST /auth/login`
- `POST /auth/logout`
- `GET /auth/me`

### `/matches`

- `GET /matches`

### `/matches/[matchId]`

- `GET /matches/:id`
- `GET /matches/:id/my-checkin`
- `POST /matches/:id/checkin`
- `PUT /matches/:id/checkin`

### `/teams/[teamId]`

- `GET /teams/:id`

### `/players/[playerId]`

- `GET /players/:id`

### `/me`

- `GET /auth/me`
- `GET /me/profile`
- `GET /me/checkins`

---

## 11. 验收定义

以下条件全部满足，视为 v1 达成：

1. 前后端工程可启动
2. seed data 可导入
3. 用户可登录并维持会话
4. 用户可浏览比赛列表并进入详情
5. 用户可创建和编辑唯一一条 check-in
6. 比赛详情页可展示聚合信息和短评
7. 个人主页可展示档案摘要与历史记录
8. 核心路径测试通过

---

## 12. 开发顺序

1. Foundation
2. Domain Data
3. Auth
4. Match Read
5. CheckIn Write
6. Match Aggregation
7. Profile
8. QA / Release

---

## 13. 后续文档关系

后续实现时，建议使用以下文档关系：

- 本文档：主 spec
- `final-whistle-schema-draft.md`：数据库设计
- `final-whistle-api-dto-draft.md`：接口契约
- `mvp/final-whistle-technical-design-main.md`：系统设计与开发策略

如果三者冲突，以本文档和后续确认后的接口/schema 文档为准。
