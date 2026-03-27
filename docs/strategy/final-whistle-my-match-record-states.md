# Final Whistle `My Match Record` 四状态细化

## 1. 目标

这份文档只聚焦比赛详情页中的核心区块：

> `My Match Record`

目标是把它的四种状态拆清楚：

1. 未登录
2. 已登录但未记录
3. 正在编辑
4. 已保存

这四种状态共同定义 Final Whistle 的产品气质。

---

## 2. 总原则

`My Match Record` 不是一个普通表单容器，而是比赛详情页的主操作区。

它需要同时做到：

- 让用户知道这里是主动作
- 让用户在赛后自然进入记录状态
- 让记录完成后具备“被归档”的感觉

所以这四种状态必须满足：

- 语义不同
- 气质统一
- 视觉层级一致

---

## 3. 状态一：未登录

## 3.1 用户心智

当前用户还不能记录，但已经进入了“我可能要留下判断”的场景。

所以未登录态不应该只是权限阻断，而应该是轻量引导。

## 3.2 目标

- 不打断赛后情绪
- 不把用户推离比赛详情页
- 明确告诉用户登录的价值是“留下这场比赛”

## 3.3 推荐结构

```text
MY MATCH RECORD

Sign in to leave your own final whistle for this match.
Save your ratings, tags, and short review into your archive.

[Go to Dev Login]
```

## 3.4 设计要求

- CTA 清晰，但不夸张
- 语气更像邀请，不像权限警告
- 文案必须强调“记录价值”，而不是“登录本身”

## 3.5 不推荐

- 过强的提示框感
- 红色警示态
- “You must sign in” 这种过硬文案

---

## 4. 状态二：已登录但未记录

## 4.1 用户心智

这是最关键的状态。

用户已经具备记录能力，但是否愿意开始写，就取决于这里的入口设计。

## 4.2 目标

- 让用户明确知道记录是这一页的主动作
- 降低开始记录的心理门槛
- 传递“这是加入我的档案”的感觉

## 4.3 推荐结构

```text
MY MATCH RECORD

You have not recorded this match yet.
Leave your final whistle while the memory is still fresh.

[Leave Your Final Whistle]
```

## 4.4 交互要求

- 主按钮必须明显
- 上方文案要有赛后情绪承接
- 区块本身要明显比社区块更重要

## 4.5 文案方向

优先考虑：

- Leave Your Final Whistle
- Save Your Match Record
- Add This Match to Your Archive

不建议长期停留在：

- Record This Match

因为它更像工具动作，不够有产品辨识度。

---

## 5. 状态三：正在编辑

## 5.1 用户心智

用户已经决定留下记录。

此时页面的任务不是提供更多信息，而是帮助用户流畅表达。

## 5.2 目标

- 让记录过程像表达，而不是录入
- 让字段顺序符合赛后思路
- 保持页面安静，减少干扰

## 5.3 表单顺序建议

推荐顺序：

1. Match Rating
2. Supporter Side
3. Team Ratings
4. Tags
5. Short Review
6. Player Ratings
7. Watched Type / Watched At

## 5.4 视觉要求

- 区块应有明显边界，但不要像后台表单
- 第一屏尽量就看到记录动作和核心输入
- `Short Review` 需要更好的书写空间
- `Player Ratings` 是延展表达区，不应抢最前

## 5.5 动效要求

- 从 empty state 进入 editing 应有平滑过渡
- 提交按钮的反馈要稳定，不要太碎
- 整体更像“开始写下记录”，不是“展开一个管理表单”

---

## 6. 状态四：已保存

## 6.1 用户心智

用户最需要感受到的是：

> 这条记录已经属于我了。

## 6.2 目标

- 提供明显的完成感
- 强化“归档”语义
- 让用户愿意回看

## 6.3 推荐结构

```text
SAVED RECORD

Match Rating     9/10
Supporter Side   Liverpool
Saved At         2026-03-27 22:14
Tags             Classic · Intense

"This one had title-race pressure in every phase."

Key Players
Salah 9 · Rice 8 · Odegaard 7

[Edit Record]
```

## 6.4 视觉要求

- 比编辑态更稳定
- 比普通信息卡更有重量
- 要像一条进入档案的赛后记录
- 不建议只靠 toast 证明保存成功

## 6.5 文案建议

可以出现一类很轻的确认语言：

- Saved to Your Archive
- Your Final Whistle Is Saved

但不必过度强调技术成功，而要强调“归档完成”。

---

## 7. 四状态之间的关系

```text
未登录
  ↓
获得记录资格

已登录未记录
  ↓
开始记录动作

正在编辑
  ↓
提交并完成表达

已保存
  ↓
回看 / 继续编辑 / 进入个人档案
```

这意味着 `My Match Record` 不只是一个 UI 组件，而是比赛详情页的核心状态机。

---

## 8. 判断这块设计是否成功的标准

### 成功时，用户会有这些感受

- 我很快知道这页的主动作是什么
- 我愿意开始记录
- 我记录时没有被太多内容打断
- 我提交后觉得这条内容被认真保存了

### 失败时，用户会有这些感受

- 我不知道应该先看什么
- 这像在填后台表单
- 保存后没有分量
- 社区块比我的记录更像主角
