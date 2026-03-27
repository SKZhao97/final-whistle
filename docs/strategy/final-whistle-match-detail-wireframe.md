# Final Whistle 比赛详情页文本线框图

## 1. 目标

这份文档提供更接近页面排布的文本线框图，用来帮助后续 UI 重构时统一结构判断。

它不是最终视觉稿，而是信息层级和布局草图。

---

## 2. Desktop Wireframe

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│ FINAL WHISTLE                                                              │
│                                                                             │
│ Liverpool                                               2 — 1      Arsenal │
│ Premier League · Matchday 12 · Final · Anfield                        22:00 │
│                                                                             │
│ A match worth remembering. Leave your record while the feeling is fresh.   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MY MATCH RECORD                                                            │
│                                                                             │
│ You have not recorded this match yet.                                      │
│ Leave your final whistle and save it into your archive.                    │
│                                                                             │
│ [Leave Your Final Whistle]                                                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────┬─────────────────────────────────────┐
│ COMMUNITY PULSE                       │ PLAYER BOARD                        │
│                                       │                                     │
│ Avg Match Rating     8.4              │ Salah               8.8            │
│ Home Team Avg        8.1              │ Rice                8.1            │
│ Away Team Avg        7.7              │ Odegaard            7.9            │
│                                       │                                     │
│ Hot Tags                              │                                     │
│ Classic · Intense · Comeback          │                                     │
│                                       │                                     │
│ Recent Reactions                      │                                     │
│ “Title-race football at full volume”  │                                     │
│ “One swing of momentum after another” │                                     │
└───────────────────────────────────────┴─────────────────────────────────────┘
```

---

## 3. Editing State Wireframe

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│ MY MATCH RECORD                                                            │
│                                                                             │
│ Leave your final whistle for Liverpool vs Arsenal.                         │
│                                                                             │
│ Match Rating        [ 9 ]        Supporter Side      [ Liverpool v ]       │
│ Home Team Rating    [ 8 ]        Away Team Rating    [ 7 ]                 │
│                                                                             │
│ Tags                                                                        │
│ [Classic] [Intense] [Comeback] [Dominant] [Painful]                        │
│                                                                             │
│ Short Review                                                                │
│ ┌───────────────────────────────────────────────────────────────────────┐   │
│ │ This one had title-race pressure in every phase of the game.         │   │
│ └───────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│ Player Ratings                                                              │
│ Salah      [ 9 ]   “Relentless all night”                                   │
│ Rice       [ 8 ]   “Held the midfield shape together”                       │
│ [+ Add Player Rating]                                                       │
│                                                                             │
│ Watched Type      [ Full Match v ]                                          │
│ Watched At        [ 2026-03-27 22:14 ]                                      │
│                                                                             │
│ [Save Record]                                         [Cancel]              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 4. Saved State Wireframe

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│ SAVED RECORD                                                               │
│                                                                             │
│ Match Rating         9 / 10                                                 │
│ Supporter Side       Liverpool                                              │
│ Saved At             2026-03-27 22:14                                       │
│ Tags                 Classic · Intense · Comeback                           │
│                                                                             │
│ “This one had title-race pressure in every phase of the game.”             │
│                                                                             │
│ Key Players                                                                 │
│ Salah 9 · Rice 8 · Odegaard 7                                               │
│                                                                             │
│ [Edit Record]                                                               │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Mobile Wireframe

```text
┌──────────────────────────────┐
│ Liverpool 2 — 1 Arsenal      │
│ Premier League · Final       │
│ Matchday 12 · Anfield        │
└──────────────────────────────┘

┌──────────────────────────────┐
│ MY MATCH RECORD              │
│                              │
│ Leave your final whistle     │
│ for this match.              │
│                              │
│ [Leave Your Final Whistle]   │
└──────────────────────────────┘

┌──────────────────────────────┐
│ COMMUNITY PULSE              │
│ Avg Match Rating  8.4        │
│ Hot Tags                     │
│ Classic · Intense            │
│ Recent Reactions             │
└──────────────────────────────┘

┌──────────────────────────────┐
│ PLAYER BOARD                 │
│ Salah        8.8             │
│ Rice         8.1             │
└──────────────────────────────┘
```

---

## 6. 结构重点

这份线框图强调三件事：

1. Hero 是封面，不是信息堆叠区
2. My Match Record 是页面主角
3. Community Pulse 是补充层，不应压过记录区

---

## 7. 后续继续细化的方向

- Match Hero 的视觉节奏
- Saved State 的归档卡片语言
- Mobile 下 `My Match Record` 的展开方式
- Community Pulse 如何压缩得更轻
