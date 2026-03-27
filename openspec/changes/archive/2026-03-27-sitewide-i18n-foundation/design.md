## Context

Final Whistle 已完成 v1 功能闭环，但当前语言层仍是临时状态：静态 UI 主要是英文，tag 字典和部分领域文案仍是中文。用户已经明确要求：主页必须能主动切换语言、切换后全站 UI 文案应尽可能实时切换、UGC 不翻译，并且 tag 不应由前端通过 slug 临时映射，而应由服务端作为领域数据负责输出本地化显示名。

当前系统特点：

- 前端是 Next.js App Router，已经有全局 auth/provider 结构
- 后端是 Go + Gin + PostgreSQL，tag 为数据库字典表
- API 目前通过 cookie session 维持登录态
- 还没有 locale provider、全局翻译字典，也没有服务端 locale-aware 领域字典输出

## Goals / Non-Goals

**Goals:**
- 为 Final Whistle 建立全站英语 / 汉语双语基础能力
- 提供首页可见、全局可用的语言切换入口
- 在不翻译 UGC 的前提下，使全站 UI 文案随当前 locale 一致切换
- 让服务端对 tag 这类领域字典返回 locale-aware 显示名
- 让当前页在切换语言后尽可能实时更新，而不依赖重启或重新 seed
- 为后续体验重构提供稳定的语言层

**Non-Goals:**
- 不引入第三种语言
- 不实现完整的 Accept-Language 协商系统
- 不实现 URL locale routing
- 不建设翻译 CMS 或后台管理系统
- 不翻译用户短评、用户名或任何 UGC
- 不在本 spec 中实现真实数据接入

## Decisions

### 1. Locale 状态由前端全局 provider 管理，语言切换入口放在首页可见位置并复用于全局 Header

前端需要一个全局 locale state，支持：

- 当前语言读取
- 切换语言
- 共享 `t()` 或等价翻译访问能力
- 页面级即时重渲染

用户要求主页必须有语言切换能力，因此入口至少要在首页明显可见。同时为了保证后续页面一致性，Header 中也应复用同一个切换组件或入口样式。

选择这个方案而不是单页局部状态，是因为 spec 的目标是“全站双语体验”，不是单页演示。

### 2. Locale 持久化使用 cookie，前端即时切换与服务端本地化输出共用同一来源

前端切换语言时：

- 立即更新 provider state，使当前页面 UI 文案直接切换
- 同时写入 locale cookie，保证刷新后仍保持选择

服务端可以在后续请求中读取同一个 cookie，决定 tag 等领域字典返回哪种显示名。这样可以让“前端 UI 文案”和“后端领域显示名”基于同一 locale 来源，而不是前后端各自维护一套状态。

之所以不只用 localStorage，是因为服务端也需要读取 locale 来返回正确的 tag 显示名。

### 3. UI 文案走前端字典，领域字典走服务端本地化输出

这次 change 明确两类内容边界：

- **UI copy**：按钮、标题、导航、空态、错误态、表单标签等，使用前端字典
- **Domain dictionary labels**：tag 显示名由服务端输出 locale-aware 的 `name`
- **UGC**：原样显示，不翻译

这种划分比“全部前端映射”更清楚，也避免 tag 这种领域对象在多端场景下重复维护翻译表。

### 4. Tag 采用双语字段存储，服务端继续对前端暴露稳定的 `name` 字段

为了满足“tag 是服务端数据源负责”的要求，同时控制 spec7 复杂度，tag 本地化采用：

- `name_en`
- `name_zh`

服务端根据当前 locale 选择其一映射到 DTO 的 `name` 字段。这样前端不需要知道 tag 的多语言存储细节，仍然拿到稳定 shape：

```json
{ "id": 1, "slug": "classic", "name": "Classic" }
```

这个方案比单独建翻译表简单，足够覆盖当前仅支持中英双语的需求。

### 5. 当前页切换语言的目标是“用户感知上的实时”，允许局部数据刷新

静态 UI 文案切换可以立即由前端 provider 完成。

对于依赖服务端 locale 输出的区域，例如 tag 名称：

- 切换语言后需要触发当前页数据刷新
- 可通过 `router.refresh()` 或等价机制完成

只要用户看到的是“切换后当前页很快完成更新”，就满足本 spec 的实时目标。这里不追求 WebSocket 级别或多标签强一致同步。

## Risks / Trade-offs

- **[Tag schema 增加双语字段]** → 需要 migration、seed 更新和旧数据迁移；但比单独翻译表更轻，适合当前阶段。
- **[前端与后端同时参与 locale]** → 需要明确边界；通过“UI copy 前端负责，领域字典后端负责”来降低混乱。
- **[切换后需要刷新服务端数据区]** → 不是纯客户端瞬时完成；但这仍符合“用户感知上的实时”，复杂度更可控。
- **[全站文案改造面较大]** → 容易漏翻；需要任务按页面和共享组件拆清楚，并在验收阶段做语言切换 smoke。
- **[UGC 保持原文]** → 某些页面可能出现“中文 UGC + 英文 UI”或反之；这是产品边界而非缺陷，需要在文档中明确。

## Migration Plan

1. 为 `tags` 增加双语字段并准备 migration
2. 更新 tag seed 逻辑，确保现有本地库能被重 seed 到双语数据
3. 后端读 locale cookie 并统一映射 tag `name`
4. 前端引入 locale provider、字典和切换组件
5. 按页面和共享组件逐步替换静态文案
6. 做中英切换 smoke，确认 UI 与 tag 同步切换

回滚策略：

- 前端 locale provider 可回退到固定默认语言
- 后端 tag 输出可暂时固定读取一个默认字段
- 新增双语列不影响旧单语读取的兼容性

## Open Questions

- 首页语言切换入口是否同时放进全局 Header，还是首页独立显式展示、Header 复用简化版组件？
- 默认语言最终采用固定英文，还是首次访问读取浏览器语言？当前设计倾向固定英文。
- 未来是否要把 team / competition 等领域显示名也走服务端本地化？spec7 先只强制 tags。
