# Final Whistle 数据库 Schema 草案

## 1. 说明

本文档定义 Final Whistle v1 的 PostgreSQL schema 草案，用于指导：

- migration 编写
- GORM model 设计
- repository 查询实现

设计原则：

- 优先满足 v1 主路径
- 保持结构化建模
- 仅在必要处为未来扩展预留字段

---

## 2. 通用约定

### 2.1 命名

- 表名：复数、`snake_case`
- 字段名：`snake_case`
- 主键：`id`
- 时间字段：
  - `created_at`
  - `updated_at`

### 2.2 时间类型

建议统一使用：

- `TIMESTAMPTZ`

### 2.3 主键

v1 可采用：

- `BIGSERIAL PRIMARY KEY`

---

## 3. 枚举建议

为减少数据库耦合，v1 可优先用 `VARCHAR` + 应用层校验。

涉及枚举值：

- `match.status`
  - `SCHEDULED`
  - `FINISHED`
- `check_ins.watched_type`
  - `FULL`
  - `PARTIAL`
  - `HIGHLIGHTS`
- `check_ins.supporter_side`
  - `HOME`
  - `AWAY`
  - `NEUTRAL`

---

## 4. 表结构

### 4.1 `users`

```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  avatar_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

索引：

- `UNIQUE(email)`

### 4.2 `teams`

```sql
CREATE TABLE teams (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  short_name VARCHAR(50),
  slug VARCHAR(100) NOT NULL UNIQUE,
  logo_url TEXT,
  external_source VARCHAR(50),
  external_id VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

索引：

```sql
CREATE INDEX idx_teams_external ON teams(external_source, external_id);
```

### 4.3 `players`

```sql
CREATE TABLE players (
  id BIGSERIAL PRIMARY KEY,
  team_id BIGINT NOT NULL REFERENCES teams(id),
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL UNIQUE,
  position VARCHAR(50),
  avatar_url TEXT,
  external_source VARCHAR(50),
  external_id VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

索引：

```sql
CREATE INDEX idx_players_team_id ON players(team_id);
CREATE INDEX idx_players_external ON players(external_source, external_id);
```

### 4.4 `matches`

```sql
CREATE TABLE matches (
  id BIGSERIAL PRIMARY KEY,
  competition VARCHAR(100) NOT NULL,
  season VARCHAR(50) NOT NULL,
  round VARCHAR(100),
  status VARCHAR(20) NOT NULL,
  kickoff_at TIMESTAMPTZ NOT NULL,
  home_team_id BIGINT NOT NULL REFERENCES teams(id),
  away_team_id BIGINT NOT NULL REFERENCES teams(id),
  home_score INT,
  away_score INT,
  venue VARCHAR(255),
  external_source VARCHAR(50),
  external_id VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_match_status CHECK (status IN ('SCHEDULED', 'FINISHED')),
  CONSTRAINT chk_match_teams CHECK (home_team_id <> away_team_id)
);
```

索引：

```sql
CREATE INDEX idx_matches_kickoff_at ON matches(kickoff_at DESC);
CREATE INDEX idx_matches_competition_season ON matches(competition, season);
CREATE INDEX idx_matches_home_team_id ON matches(home_team_id);
CREATE INDEX idx_matches_away_team_id ON matches(away_team_id);
CREATE INDEX idx_matches_external ON matches(external_source, external_id);
```

### 4.5 `match_players`

```sql
CREATE TABLE match_players (
  id BIGSERIAL PRIMARY KEY,
  match_id BIGINT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
  player_id BIGINT NOT NULL REFERENCES players(id),
  team_id BIGINT NOT NULL REFERENCES teams(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (match_id, player_id)
);
```

索引：

```sql
CREATE INDEX idx_match_players_match_id ON match_players(match_id);
CREATE INDEX idx_match_players_player_id ON match_players(player_id);
CREATE INDEX idx_match_players_team_id ON match_players(team_id);
```

### 4.6 `tags`

```sql
CREATE TABLE tags (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  slug VARCHAR(50) NOT NULL UNIQUE,
  sort_order INT NOT NULL DEFAULT 0,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

索引：

```sql
CREATE INDEX idx_tags_is_active_sort_order ON tags(is_active, sort_order);
```

### 4.7 `check_ins`

```sql
CREATE TABLE check_ins (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  match_id BIGINT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
  watched_type VARCHAR(20) NOT NULL,
  supporter_side VARCHAR(20) NOT NULL,
  match_rating INT NOT NULL,
  home_team_rating INT NOT NULL,
  away_team_rating INT NOT NULL,
  short_review VARCHAR(280),
  watched_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, match_id),
  CONSTRAINT chk_check_ins_watched_type CHECK (watched_type IN ('FULL', 'PARTIAL', 'HIGHLIGHTS')),
  CONSTRAINT chk_check_ins_supporter_side CHECK (supporter_side IN ('HOME', 'AWAY', 'NEUTRAL')),
  CONSTRAINT chk_check_ins_match_rating CHECK (match_rating BETWEEN 1 AND 10),
  CONSTRAINT chk_check_ins_home_team_rating CHECK (home_team_rating BETWEEN 1 AND 10),
  CONSTRAINT chk_check_ins_away_team_rating CHECK (away_team_rating BETWEEN 1 AND 10)
);
```

索引：

```sql
CREATE INDEX idx_check_ins_user_id ON check_ins(user_id);
CREATE INDEX idx_check_ins_match_id ON check_ins(match_id);
CREATE INDEX idx_check_ins_created_at ON check_ins(created_at DESC);
CREATE INDEX idx_check_ins_watched_at ON check_ins(watched_at DESC);
```

### 4.8 `player_ratings`

```sql
CREATE TABLE player_ratings (
  id BIGSERIAL PRIMARY KEY,
  check_in_id BIGINT NOT NULL REFERENCES check_ins(id) ON DELETE CASCADE,
  player_id BIGINT NOT NULL REFERENCES players(id),
  rating INT NOT NULL,
  note VARCHAR(80),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (check_in_id, player_id),
  CONSTRAINT chk_player_ratings_rating CHECK (rating BETWEEN 1 AND 10)
);
```

索引：

```sql
CREATE INDEX idx_player_ratings_check_in_id ON player_ratings(check_in_id);
CREATE INDEX idx_player_ratings_player_id ON player_ratings(player_id);
```

说明：

- “每条 check-in 最多 5 条球员评分”建议在应用层校验，不用数据库约束实现

### 4.9 `check_in_tags`

```sql
CREATE TABLE check_in_tags (
  id BIGSERIAL PRIMARY KEY,
  check_in_id BIGINT NOT NULL REFERENCES check_ins(id) ON DELETE CASCADE,
  tag_id BIGINT NOT NULL REFERENCES tags(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (check_in_id, tag_id)
);
```

索引：

```sql
CREATE INDEX idx_check_in_tags_check_in_id ON check_in_tags(check_in_id);
CREATE INDEX idx_check_in_tags_tag_id ON check_in_tags(tag_id);
```

### 4.10 `sessions`

```sql
CREATE TABLE sessions (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token VARCHAR(255) NOT NULL UNIQUE,
  expired_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

索引：

```sql
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expired_at ON sessions(expired_at);
```

---

## 5. 推荐查询支撑

### 5.1 比赛列表

依赖：

- `matches`
- `teams`
- `check_ins`

说明：

- 先基于比赛主表分页
- 聚合统计通过 join/subquery 补充

### 5.2 比赛详情聚合

依赖：

- `check_ins`
- `player_ratings`
- `check_in_tags`

说明：

- 直接实时聚合
- v1 不引入聚合快照表

### 5.3 用户主页

依赖：

- `check_ins`
- `player_ratings`
- `matches`

---

## 6. 不在 v1 实现的数据库设计

以下能力明确不进入 v1 schema：

- 所有表软删除 `deleted_at`
- 审计日志表
- 聚合快照表
- 事件表
- 限流表

---

## 7. migration 顺序建议

1. `users`
2. `teams`
3. `players`
4. `matches`
5. `match_players`
6. `tags`
7. `check_ins`
8. `player_ratings`
9. `check_in_tags`
10. `sessions`
