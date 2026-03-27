# Final Whistle 里程碑

## Milestone 0: Spec Freeze

目标：

- 锁定 v1 范围、认证方案、数据来源、接口契约、聚合规则

完成标准：

- PRD review 完成
- v1 spec 完成
- 开发计划完成
- 里程碑和验收口径完成

交付物：

- `final-whistle-review.md`
- `final-whistle-development-plan.md`
- `mvp/final-whistle-v1-spec.md`
- `final-whistle-milestones.md`

---

## Milestone 1: 项目地基可运行

目标：

- 工程骨架搭好，数据库可初始化，seed 可导入

完成标准：

- 前端工程可启动
- 后端工程可启动
- PostgreSQL migration 可执行
- seed data 可写入
- 至少有健康检查和统一错误返回

验收方式：

- 本地启动前后端成功
- 数据库中可看到基础实体数据

---

## Milestone 2: 只读浏览闭环

目标：

- 用户可以浏览比赛、查看比赛详情、查看球队和球员简版页

完成标准：

- `GET /matches` 可用
- `GET /matches/:id` 可用
- `GET /teams/:id` 可用
- `GET /players/:id` 可用
- 前端比赛列表页和比赛详情页可正常展示

验收方式：

- 能从列表进入详情
- 详情页能展示基础信息、聚合占位或空状态

---

## Milestone 3: 登录与打卡闭环

目标：

- 登录用户可以为比赛创建和编辑唯一一条 check-in

完成标准：

- `POST /auth/login` 可用
- `GET /auth/me` 可用
- `POST /matches/:id/checkin` 可用
- `PUT /matches/:id/checkin` 可用
- `GET /matches/:id/my-checkin` 可用
- 前端打卡表单可创建和编辑记录

验收方式：

- 未登录用户无法提交打卡
- 登录后可成功创建记录
- 二次进入详情页可看到“我的记录”
- 编辑后页面展示更新

---

## Milestone 4: 聚合与档案闭环

目标：

- 提交后的数据能体现在比赛聚合和用户档案中

完成标准：

- 比赛详情页显示平均比赛评分、球队评分、球员评分排行、最近短评
- `GET /me/profile` 可用
- `GET /me/checkins` 可用
- 个人主页显示最近记录和基础统计

验收方式：

- 新建或编辑 check-in 后聚合结果变化正确
- 个人主页能看到自己的历史记录

---

## Milestone 5: 测试与发布准备

目标：

- 核心流程稳定，可部署演示

完成标准：

- 后端关键 service/handler 测试通过
- 前端关键 E2E 通过
- loading、empty、error 状态完善
- 部署说明完整

验收方式：

- 完整跑通核心链路：
  - 登录
  - 浏览比赛
  - 提交打卡
  - 查看聚合
  - 查看个人主页

---

## 推荐节奏

- M0: 0.5 天
- M1: 0.5 到 1 天
- M2: 1 天
- M3: 1 天
- M4: 1 天
- M5: 0.5 到 1 天

总计：

- 约 4.5 到 5.5 天

---

## Kill Criteria

出现以下情况时，必须收缩范围，不继续扩展：

- 认证方案仍未定
- 开始接第三方足球数据 API
- 球队页和球员页开始做复杂统计
- 打卡表单新增长文本、互动、点赞、评论楼中楼
