# Final Whistle 比赛详情页 UI 概念草图

## 1. 目标

这份文档提供一个非实现层的 UI 概念稿，用来帮助判断比赛详情页应该长成什么气质。

它服务于 A-first 路线：

- 比赛详情页首先是赛后记录页
- 我的记录是主角
- 社区脉搏是补充层

---

## 2. 视觉基调

### 推荐方向

- `Record Book`
- `Editorial Layout`
- `Afterglow of the Final Whistle`

### 不推荐方向

- 黑色占主导的内容平台风格
- 体育资讯门户页风格
- 社区 feed / 热门讨论流风格

### 配色概念

```text
Base
- Warm Ivory / Paper White
- Soft Stone

Text
- Ink Black
- Slate Gray

Accents
- Deep Burgundy
- Field Green
- Brushed Silver
```

简化理解：

```text
底色像纸
文字像墨
强调像赛场余温
```

---

## 3. 页面结构草图

```text
┌──────────────────────────────────────────────────────────────┐
│ FINAL WHISTLE                                                │
│                                                              │
│  Liverpool                           2 — 1           Arsenal │
│  Premier League · Matchday 12 · Final · Anfield             │
│                                                              │
│  A match worth remembering. Leave your record below.         │
└──────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────┐
│ MY MATCH RECORD                                              │
│                                                              │
│  [Leave Your Final Whistle]                                  │
│                                                              │
│  Match Rating     Supporter Side     Team Ratings            │
│  Tags             Short Review       Player Ratings          │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ “One of the most chaotic title-race matches this year”│  │
│  └────────────────────────────────────────────────────────┘  │
│                                                              │
│  [Save Record]                               [Cancel]        │
└──────────────────────────────────────────────────────────────┘

┌──────────────────────────────┬───────────────────────────────┐
│ COMMUNITY PULSE             │ PLAYER BOARD                  │
│                              │                               │
│ Average match rating  8.4    │ 1. Salah          8.8         │
│ Hot tags: Classic / Tense    │ 2. Rice           8.1         │
│ Recent reactions             │ 3. Odegaard       7.9         │
└──────────────────────────────┴───────────────────────────────┘
```

---

## 4. 重点：My Match Record 区块

这个区块不是普通卡片，而应该有“被认真书写”的感觉。

### Empty State

```text
┌──────────────────────────────────────────────────────────────┐
│ MY MATCH RECORD                                              │
│                                                              │
│ You have not recorded this match yet.                        │
│ Save your own final whistle while the memory is still fresh. │
│                                                              │
│ [Leave Your Final Whistle]                                   │
└──────────────────────────────────────────────────────────────┘
```

### Saved State

```text
┌──────────────────────────────────────────────────────────────┐
│ SAVED RECORD                                                 │
│                                                              │
│ Match Rating     9/10                                        │
│ Supporter Side   Liverpool                                   │
│ Saved At         2026-03-27 22:14                            │
│ Tags             Classic · Intense · Comeback                │
│                                                              │
│ “This one had title-race pressure in every phase.”           │
│                                                              │
│ Key Players                                                 │
│ Salah 9 · Rice 8 · Odegaard 7                               │
│                                                              │
│ [Edit Record]                                                │
└──────────────────────────────────────────────────────────────┘
```

这个状态要给人的感觉是：

> 这不是临时表单结果，而是已经进入我的档案的一条记录。

---

## 5. 视觉层级建议

```text
Hero
  ↓
My Match Record
  ↓
Community Pulse
  ↓
Player Board
```

其中：

- Hero 是封面
- My Match Record 是正文
- Community Pulse 是附注
- Player Board 是补充说明

如果 Community 和 Player Board 的视觉存在感接近或超过 My Match Record，说明层级有问题。

---

## 6. 黑色为什么不适合做主底色

当前如果整页过黑，会把这页带向：

- 比赛直播总结页
- 内容平台详情页
- 气氛很重但不够细腻的界面

而 Final Whistle 想要的不是“压迫的赛事感”，而是“赛后留下记录的分量”。

所以更建议：

- 用浅底来放大排版与留白
- 用深色只强调关键动作和记录完成状态
- 让页面更像“被写下的一页”而不是“被刷过的一屏”

---

## 7. 可以继续往下拆的方向

如果这份 UI 概念方向成立，后面可以继续拆：

1. Match Hero 的标题与比分排版方案
2. My Match Record 的四种状态细稿
3. Saved State 的归档感组件语言
4. `/me` 页如何延续同一套视觉体系
