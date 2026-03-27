# Final Whistle API DTO 草案

## 1. 说明

本文档定义 Final Whistle v1 的 API DTO 草案，用于指导：

- 后端 `dto/` 目录实现
- 前端 `types/api.ts` 设计
- OpenAPI 草案编写

约定：

- API 字段使用 `camelCase`
- 所有响应包裹在统一结构中
- 列表接口按统一分页格式返回

---

## 2. 通用响应 DTO

### 2.1 SuccessResponse

```ts
type SuccessResponse<T> = {
  success: true;
  data: T;
};
```

### 2.2 ErrorResponse

```ts
type ErrorResponse = {
  success: false;
  error: {
    code:
      | "UNAUTHORIZED"
      | "FORBIDDEN"
      | "NOT_FOUND"
      | "VALIDATION_ERROR"
      | "CONFLICT"
      | "INTERNAL_ERROR";
    message: string;
    details?: Record<string, unknown>;
  };
};
```

### 2.3 PaginatedData

```ts
type PaginatedData<T> = {
  items: T[];
  page: number;
  pageSize: number;
  total: number;
};
```

---

## 3. 通用基础 DTO

### 3.1 UserSummaryDTO

```ts
type UserSummaryDTO = {
  id: number;
  name: string;
  email: string;
  avatarUrl: string | null;
};
```

### 3.2 TeamSummaryDTO

```ts
type TeamSummaryDTO = {
  id: number;
  name: string;
  shortName: string | null;
  slug: string;
  logoUrl: string | null;
};
```

### 3.3 PlayerSummaryDTO

```ts
type PlayerSummaryDTO = {
  id: number;
  name: string;
  slug: string;
  position: string | null;
  avatarUrl: string | null;
  team: TeamSummaryDTO;
};
```

### 3.4 TagDTO

```ts
type TagDTO = {
  id: number;
  name: string;
  slug: string;
};
```

### 3.5 MatchScoreDTO

```ts
type MatchScoreDTO = {
  homeScore: number | null;
  awayScore: number | null;
};
```

### 3.6 MatchBaseDTO

```ts
type MatchBaseDTO = {
  id: number;
  competition: string;
  season: string;
  round: string | null;
  status: "SCHEDULED" | "FINISHED";
  kickoffAt: string;
  venue: string | null;
  homeTeam: TeamSummaryDTO;
  awayTeam: TeamSummaryDTO;
  score: MatchScoreDTO;
};
```

---

## 4. Auth DTO

### 4.1 POST `/auth/login`

Request:

```ts
type LoginRequestDTO = {
  email: string;
  name: string;
};
```

Response:

```ts
type LoginResponseDTO = {
  user: UserSummaryDTO;
};
```

说明：

- session 通过 Cookie 下发，不在 JSON 中返回 token

### 4.2 POST `/auth/logout`

Response:

```ts
type LogoutResponseDTO = {
  ok: true;
};
```

### 4.3 GET `/auth/me`

Response:

```ts
type MeResponseDTO = {
  user: UserSummaryDTO;
};
```

---

## 5. Match DTO

### 5.1 GET `/matches`

Query:

```ts
type ListMatchesQueryDTO = {
  competition?: string;
  season?: string;
  page?: number;
  pageSize?: number;
};
```

Response item:

```ts
type MatchListItemDTO = MatchBaseDTO & {
  aggregate: {
    averageMatchRating: number | null;
    checkInCount: number;
  };
};
```

Response:

```ts
type ListMatchesResponseDTO = PaginatedData<MatchListItemDTO>;
```

### 5.2 GET `/matches/:id`

```ts
type MatchPlayerRankingDTO = {
  player: PlayerSummaryDTO;
  averageRating: number;
  ratingCount: number;
};
```

```ts
type MatchReviewItemDTO = {
  id: number;
  user: UserSummaryDTO;
  matchRating: number;
  shortReview: string | null;
  tags: TagDTO[];
  createdAt: string;
};
```

```ts
type MatchDetailResponseDTO = {
  match: MatchBaseDTO;
  aggregate: {
    averageMatchRating: number | null;
    averageHomeTeamRating: number | null;
    averageAwayTeamRating: number | null;
    checkInCount: number;
    lowSample: boolean;
  };
  topPlayers: MatchPlayerRankingDTO[];
  recentReviews: MatchReviewItemDTO[];
};
```

### 5.3 GET `/matches/:id/my-checkin`

Response:

```ts
type MyCheckInResponseDTO = CheckInDetailDTO | null;
```

---

## 6. CheckIn DTO

### 6.1 基础 DTO

```ts
type PlayerRatingInputDTO = {
  playerId: number;
  rating: number;
  note?: string;
};
```

```ts
type PlayerRatingDTO = {
  id: number;
  player: PlayerSummaryDTO;
  rating: number;
  note: string | null;
};
```

```ts
type CheckInDetailDTO = {
  id: number;
  matchId: number;
  watchedType: "FULL" | "PARTIAL" | "HIGHLIGHTS";
  supporterSide: "HOME" | "AWAY" | "NEUTRAL";
  matchRating: number;
  homeTeamRating: number;
  awayTeamRating: number;
  shortReview: string | null;
  watchedAt: string;
  tags: TagDTO[];
  playerRatings: PlayerRatingDTO[];
  createdAt: string;
  updatedAt: string;
};
```

### 6.2 POST `/matches/:id/checkin`

Request:

```ts
type CreateCheckInRequestDTO = {
  watchedType: "FULL" | "PARTIAL" | "HIGHLIGHTS";
  supporterSide: "HOME" | "AWAY" | "NEUTRAL";
  matchRating: number;
  homeTeamRating: number;
  awayTeamRating: number;
  shortReview?: string;
  watchedAt: string;
  tags: number[];
  playerRatings: PlayerRatingInputDTO[];
};
```

Response:

```ts
type CreateCheckInResponseDTO = CheckInDetailDTO;
```

### 6.3 PUT `/matches/:id/checkin`

Request:

```ts
type UpdateCheckInRequestDTO = CreateCheckInRequestDTO;
```

Response:

```ts
type UpdateCheckInResponseDTO = CheckInDetailDTO;
```

---

## 7. Team DTO

### 7.1 GET `/teams/:id`

```ts
type TeamRelatedMatchDTO = MatchBaseDTO & {
  aggregate: {
    averageMatchRating: number | null;
    checkInCount: number;
  };
};
```

```ts
type TeamDetailResponseDTO = {
  team: TeamSummaryDTO;
  summary: {
    averageTeamRating: number | null;
    ratingCount: number;
  };
  recentMatches: TeamRelatedMatchDTO[];
};
```

---

## 8. Player DTO

### 8.1 GET `/players/:id`

```ts
type PlayerRatedMatchDTO = MatchBaseDTO & {
  averageRating: number;
  ratingCount: number;
};
```

```ts
type PlayerDetailResponseDTO = {
  player: PlayerSummaryDTO;
  summary: {
    averageRating: number | null;
    ratingCount: number;
  };
  recentRatedMatches: PlayerRatedMatchDTO[];
};
```

---

## 9. User DTO

### 9.1 GET `/me/profile`

```ts
type ProfileRecentCheckInDTO = {
  id: number;
  match: MatchBaseDTO;
  matchRating: number;
  shortReview: string | null;
  watchedAt: string;
  createdAt: string;
};
```

```ts
type FavoriteTeamDTO = {
  team: TeamSummaryDTO;
  checkInCount: number;
};
```

```ts
type FavoritePlayerDTO = {
  player: PlayerSummaryDTO;
  ratingCount: number;
};
```

```ts
type ProfileResponseDTO = {
  user: UserSummaryDTO;
  stats: {
    totalCheckIns: number;
    averageMatchRating: number | null;
  };
  recentCheckIns: ProfileRecentCheckInDTO[];
  recentReviews: MatchReviewItemDTO[];
  favoriteTeams: FavoriteTeamDTO[];
  favoritePlayers: FavoritePlayerDTO[];
};
```

### 9.2 GET `/me/checkins`

```ts
type MyCheckInListItemDTO = {
  id: number;
  match: MatchBaseDTO;
  matchRating: number;
  homeTeamRating: number;
  awayTeamRating: number;
  shortReview: string | null;
  watchedAt: string;
  createdAt: string;
  updatedAt: string;
};
```

```ts
type MyCheckInsResponseDTO = PaginatedData<MyCheckInListItemDTO>;
```

---

## 10. 最低校验规则

### 10.1 LoginRequestDTO

- `email` 必填，合法邮箱
- `name` 必填，长度 `1-100`

### 10.2 Create/UpdateCheckInRequestDTO

- `watchedType` 必填
- `supporterSide` 必填
- `matchRating` 必填，范围 `1-10`
- `homeTeamRating` 必填，范围 `1-10`
- `awayTeamRating` 必填，范围 `1-10`
- `shortReview` 可选，最大 `280`
- `watchedAt` 必填，合法时间
- `tags` 可为空数组
- `playerRatings` 可为空数组，最大 `5`
- `playerRatings[].rating` 范围 `1-10`
- `playerRatings[].note` 最大 `80`

---

## 11. v1 不进入 DTO 的内容

- `requestId`
- `pagination.totalPages`
- 限流相关字段
- 缓存命中信息
- 后台管理字段
