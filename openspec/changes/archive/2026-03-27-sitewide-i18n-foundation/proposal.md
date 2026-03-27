## Why

Final Whistle 现在已经有完整的 v1 闭环，但产品表面仍存在明显的语言不一致问题，例如 UI 主文案以英文为主，而标签等领域内容仍以中文展示。既然下一阶段要继续打磨产品体验，就需要先建立一套可用、统一、可持续的双语基础能力，避免后续 UI 重构再返工一遍文案与领域显示层。

## What Changes

- 引入全站英语 / 汉语切换能力，并在首页提供明确的语言选择入口。
- 建立前端 locale 状态、文案字典和共享翻译访问方式，使全站 UI 文案、按钮、表单标签、空态、错误态和导航文案可随当前语言切换。
- 将 tag 显示名改为由服务端按当前 locale 返回，而不是前端根据 slug 临时映射。
- 明确 UGC 不参与翻译，用户短评、用户名等内容保持原样显示。
- 为后续 UI 重构建立稳定的语言层，让 Spec8 可以专注做 A-first 体验而不是再次处理全站文案系统。

## Capabilities

### New Capabilities
- `sitewide-i18n`: 全站双语基础能力，包括语言切换、locale 持久化、UI 文案翻译、以及 UGC 非翻译边界。

### Modified Capabilities
- `frontend-framework`: 前端基础设施需要支持全局 locale state、翻译字典访问和语言持久化。
- `database-schema`: 标签字典需要存储英语和汉语显示名，以支持服务端本地化输出。
- `data-models`: Tag 相关模型需要支持双语显示名字段和 locale-aware 映射。
- `seed-data`: 预置 tag 字典需要提供中英双语内容，而不是单语内容。
- `public-match-browsing`: 含 tag 的公开比赛详情响应需要根据当前 locale 返回对应显示名。
- `checkin-domain-and-api`: 当前用户 check-in 读取与写入响应中的 tag 显示名需要根据当前 locale 返回。
- `user-profile-and-aggregation`: `/me/checkins` 等历史响应中的 tag 显示名需要根据当前 locale 返回。

## Impact

- 前端应用壳、Header / 首页入口、页面文案组织方式、共享 UI helper
- 后端 tag schema、model、seed、DTO 映射和 locale 读取逻辑
- 公开比赛详情、check-in 相关响应、用户历史响应中的 tag 展示字段
- 后续 UI 重构、产品体验文案和多页面一致性
