# Final Whistle 模块化开发计划

## 目标

在最短路径内交付一个可上线、可演示、可继续扩展的 v1，确保核心闭环为：

`登录 -> 浏览比赛 -> 提交打卡 -> 查看聚合 -> 查看个人档案`

---

## 模块拆解

### 1. Foundation 模块

职责：

- 建立前后端工程骨架
- 配置环境变量
- 配置数据库连接
- 建立 migration 和 seed 机制
- 统一 API 错误格式和基础中间件

输出：

- `frontend/` 基础工程
- `backend/` 基础工程
- PostgreSQL schema 初始化能力
- seed 数据导入能力

依赖：

- 无

---

### 2. Domain Data 模块

职责：

- 建模 `users`, `teams`, `players`, `matches`, `match_players`, `tags`
- 准备 v1 种子数据
- 明确实体关系与约束

输出：

- 数据表结构
- 初始标签字典
- 初始比赛数据
- 查询基础 repository

依赖：

- Foundation

---

### 3. Auth 模块

职责：

- 提供登录、登出、当前用户接口
- 维护基于 Cookie Session 的登录态

输出：

- `POST /auth/login`
- `POST /auth/logout`
- `GET /auth/me`
- 前端登录页与登录态管理

依赖：

- Foundation
- Domain Data

---

### 4. Match Read 模块

职责：

- 提供比赛列表、比赛详情、球队详情、球员详情
- 输出比赛基础信息和只读聚合信息

输出：

- `GET /matches`
- `GET /matches/:id`
- `GET /teams/:id`
- `GET /players/:id`
- 前端比赛列表页
- 前端比赛详情页基础信息区

依赖：

- Domain Data

---

### 5. CheckIn Write 模块

职责：

- 创建和更新用户对比赛的唯一记录
- 校验球员是否属于该场比赛
- 保存标签和球员评分

输出：

- `GET /matches/:id/my-checkin`
- `POST /matches/:id/checkin`
- `PUT /matches/:id/checkin`
- 打卡表单 UI

依赖：

- Auth
- Match Read

---

### 6. Match Aggregation 模块

职责：

- 在比赛详情页展示聚合评分、球员评分榜和最近短评
- 在 check-in 提交后反映最新聚合结果

输出：

- 比赛详情页聚合区
- 比赛详情页短评流
- 样本不足提示

依赖：

- CheckIn Write
- Match Read

---

### 7. Profile 模块

职责：

- 聚合当前用户档案摘要
- 展示历史记录和最近短评

输出：

- `GET /me/profile`
- `GET /me/checkins`
- 前端个人主页

依赖：

- Auth
- CheckIn Write

---

### 8. QA & Release 模块

职责：

- 完成关键路径测试
- 部署前环境检查
- 演示数据验收

输出：

- 后端 service/handler 测试
- 前端 E2E
- 部署说明
- 发布检查清单

依赖：

- 全部模块

---

## 推荐开发顺序

1. Foundation
2. Domain Data
3. Auth
4. Match Read
5. CheckIn Write
6. Match Aggregation
7. Profile
8. QA & Release

---

## 并行策略

### 可并行阶段 A

- 前端基础工程初始化
- 后端基础工程初始化
- 数据模型与 seed 设计

前提：

- 先确认 v1 spec 中的核心实体和认证方案

### 可并行阶段 B

- 前端比赛列表/详情静态页面
- 后端只读接口实现

前提：

- 只读 DTO 和响应格式固定

### 可并行阶段 C

- 前端打卡表单
- 后端 check-in 写接口

前提：

- check-in request/response contract 固定

---

## 风险控制

### 风险 1

球队页和球员页占用过多时间。

控制：

- 严格做成简版只读页，不增加复杂统计

### 风险 2

登录方案拖慢整体交付。

控制：

- v1 先使用 dev login + cookie session

### 风险 3

数据源范围膨胀。

控制：

- v1 只允许 seed data，不接第三方足球 API

### 风险 4

聚合查询和页面展示口径不一致。

控制：

- 先定义 spec，再写接口和页面

---

## 验收原则

每个模块验收必须满足：

- 有明确输入输出
- 有最小可演示页面或接口
- 有错误路径处理
- 不依赖未完成的扩展功能
