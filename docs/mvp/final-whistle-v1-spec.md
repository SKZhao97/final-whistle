# Final Whistle v1 Spec

## 1. 范围定义

Final Whistle v1 是一个面向足球观众的赛后记录产品，核心目标是让用户完成一条轻量但结构化的观赛记录，并在比赛详情页和个人主页中回看这条记录及其聚合结果。

v1 只覆盖以下主流程：

1. 用户登录
2. 浏览比赛列表
3. 进入比赛详情
4. 创建或编辑自己的比赛记录
5. 查看比赛聚合结果
6. 查看自己的个人档案

---

## 2. 非目标

v1 不包含：

- 实时比分
- 第三方足球数据同步
- 社交关系
- 点赞、回复、评论楼中楼
- 用户自定义标签
- 长文内容
- 复杂后台
- 赛季统计大屏

---

## 3. 关键产品决策

### 3.1 数据来源

v1 仅使用内部 seed data。

数据范围建议：

- 1 个赛季
- 2 到 4 个赛事
- 20 到 50 场比赛
- 与比赛对应的球队、球员

### 3.2 认证方案

v1 使用 `HTTP-only Cookie Session`。

登录方式：

- 默认实现 `dev login`
- GitHub OAuth 不进入 v1 必须范围

### 3.3 记录模型

用户对一场比赛只能有一条主记录，即一条 `CheckIn`。

唯一约束：

- `(user_id, match_id)`

### 3.4 聚合策略

v1 聚合全部实时查询，不做缓存，不做异步任务，不做预计算表。

---

## 4. 核心实体

### 4.1 User

字段：

- `id`
- `name`
- `email`
- `avatarUrl`
- `createdAt`
- `updatedAt`

### 4.2 Team

字段：

- `id`
- `name`
- `shortName`
- `slug`
- `logoUrl`
- `createdAt`
- `updatedAt`

### 4.3 Player

字段：

- `id`
- `teamId`
- `name`
- `slug`
- `position`
- `avatarUrl`
- `createdAt`
- `updatedAt`

### 4.4 Match

字段：

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
- `createdAt`
- `updatedAt`

### 4.5 MatchPlayer

字段：

- `id`
- `matchId`
- `playerId`
- `teamId`

用途：

- 约束用户只能给本场比赛出现的球员打分

### 4.6 Tag

字段：

- `id`
- `name`
- `slug`
- `sortOrder`
- `isActive`

### 4.7 CheckIn

字段：

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

- `shortReview` 最大 280 字

### 4.8 PlayerRating

字段：

- `id`
- `checkInId`
- `playerId`
- `rating`
- `note`

约束：

- 单条 `CheckIn` 最多 5 条 `PlayerRating`
- `note` 最大 80 字

### 4.9 CheckInTag

字段：

- `id`
- `checkInId`
- `tagId`

---

## 5. 枚举定义

### 5.1 watchedType

- `FULL`
- `PARTIAL`
- `HIGHLIGHTS`

### 5.2 supporterSide

- `HOME`
- `AWAY`
- `NEUTRAL`

### 5.3 match status

- `SCHEDULED`
- `FINISHED`

v1 仅允许对 `FINISHED` 比赛创建 check-in。

---

## 6. 业务规则

### 6.1 创建和编辑

- 用户必须登录后才能创建或编辑 check-in
- 每个用户每场比赛仅允许一条记录
- 已存在记录时只能编辑，不能再次创建

### 6.2 评分规则

- `matchRating`: `1-10` 整数
- `homeTeamRating`: `1-10` 整数
- `awayTeamRating`: `1-10` 整数
- `playerRatings[].rating`: `1-10` 整数

### 6.3 球员评分规则

- 可选，不是必填
- 最多选择 5 名球员
- 被评分球员必须属于当前比赛

### 6.4 标签规则

- 可多选
- 标签来自固定字典
- v1 不支持用户创建标签

### 6.5 文本规则

- `shortReview` 选填，最大 280 字
- `playerRatings[].note` 选填，最大 80 字

### 6.6 聚合规则

- 聚合基于当前所有有效 check-in 实时计算
- 平均分保留 1 位小数
- 当样本数少于 3 时前端显示“样本较少”
- 无数据时返回 `null` 或空数组，不返回伪默认值

---

## 7. 页面定义

### 7.1 `/matches`

目标：

- 让用户快速找到比赛

展示内容：

- 主客队
- 比分
- 比赛时间
- 赛事
- 社区平均评分
- 打卡人数

筛选：

- `competition`
- `season`

排序：

- `kickoffAt desc`

### 7.2 `/matches/[matchId]`

展示内容：

- 比赛基础信息
- 社区聚合评分
- 球员评分排行
- 最近短评
- 我的记录模块

行为：

- 未登录：提示登录
- 已登录且未打卡：显示创建入口
- 已登录且已打卡：显示记录卡片和编辑入口

### 7.3 `/teams/[teamId]`

简版只读页，展示：

- 球队基本信息
- 最近相关比赛
- 平均评分摘要

### 7.4 `/players/[playerId]`

简版只读页，展示：

- 球员基本信息
- 最近被评分比赛
- 平均评分摘要

### 7.5 `/me`

展示内容：

- 用户基础信息
- 总打卡数
- 平均比赛评分
- 最近打卡
- 最近短评
- 常打卡球队前 3
- 常评分球员前 3

---

## 8. API Contract

### 8.1 通用成功返回

```json
{
  "success": true,
  "data": {}
}
```

### 8.2 通用错误返回

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

### 8.3 接口列表

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

### 8.4 分页规范

列表接口统一返回：

- `items`
- `page`
- `pageSize`
- `total`

默认值：

- `page=1`
- `pageSize=20`
- 最大 `pageSize=50`

---

## 9. 关键接口规则

### 9.1 `GET /matches`

查询参数：

- `competition`
- `season`
- `page`
- `pageSize`

返回：

- 比赛列表
- 每场比赛的聚合摘要

### 9.2 `GET /matches/:id`

返回：

- 比赛基础信息
- 聚合评分
- 球员评分排行
- 最近短评
- 打卡人数

### 9.3 `GET /matches/:id/my-checkin`

鉴权：

- 需要登录

返回：

- 当前用户对该比赛的 check-in
- 未创建时返回 `data: null`

### 9.4 `POST /matches/:id/checkin`

鉴权：

- 需要登录

失败条件：

- 比赛不存在
- 比赛未结束
- 该用户已存在记录
- 球员不属于比赛
- 参数不合法

### 9.5 `PUT /matches/:id/checkin`

鉴权：

- 需要登录

失败条件：

- 记录不存在
- 比赛不存在
- 球员不属于比赛
- 参数不合法

### 9.6 `GET /me/profile`

返回：

- 用户基础信息
- 用户统计摘要
- 最近记录
- 最近短评

### 9.7 `GET /me/checkins`

返回：

- 当前用户历史记录列表

---

## 10. 前端状态要求

所有主要页面至少覆盖：

- loading
- empty
- error
- unauthorized

关键交互要求：

- 提交打卡时防重复提交
- 提交成功后刷新详情页和我的记录区
- 编辑成功后刷新聚合与记录卡片

---

## 11. 测试最低要求

### 后端

- `CreateCheckIn` 成功
- `CreateCheckIn` 重复创建失败
- `CreateCheckIn` 球员校验失败
- `UpdateCheckIn` 成功
- `GetMatchDetail` 聚合正确
- `GetProfileSummary` 正确

### 前端

- 登录成功
- 比赛列表加载成功
- 比赛详情加载成功
- 创建 check-in 成功
- 编辑 check-in 成功
- 个人主页加载成功

---

## 12. v1 完成定义

满足以下条件即视为 v1 完成：

1. 用户可登录
2. 用户可浏览比赛并进入详情
3. 用户可创建和编辑自己的 check-in
4. 比赛详情页可展示聚合评分和最近短评
5. 个人主页可展示个人档案摘要
6. 核心链路通过基础测试
7. 项目可在一个公开可访问环境中部署演示
