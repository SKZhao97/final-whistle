## Context

Final Whistle's v1 loop now supports public match browsing, session authentication, and check-in creation/editing. The authenticated user experience still lacks a personal profile page where users can see their recorded match statistics, aggregated ratings, and check-in history. The current `/me` page is a placeholder that only acknowledges sign-in status.

The backend has established patterns for handlers, services, repositories, DTOs, and database queries. The session‑auth spec already provides authentication middleware that can protect the new profile endpoints. The check‑in domain ensures the existence of recorded matches and player‑rating data that the profile page will aggregate.

## Goals / Non-Goals

**Goals:**
- Provide an authenticated profile summary endpoint (`GET /me/profile`) that returns the current user's base identity plus v1 profile statistics.
- Provide an authenticated user check‑in history endpoint (`GET /me/checkins`) that returns a paginated list of the user’s recorded matches with lightweight match context and check-in summary fields.
- Update the frontend `/me` page to display the user’s profile summary and check‑in history.
- Add backend aggregation queries to calculate per‑user statistics (average match ratings, most frequent tags, favorite teams).
- Follow the existing clean architecture pattern with handler, service, and repository layers.
- Add final release-quality checks for the authenticated v1 journey and `/me` page states.

**Non-Goals:**
- Creating new database tables (existing `users`, `check_ins`, `player_ratings`, `checkin_tags` tables hold all required data).
- Adding external dependencies beyond the current PostgreSQL + GORM stack.
- Building advanced analytics or real‑time dashboards beyond v1 scope.
- Supporting cross‑user profile viewing—only the signed‑in user can see their own aggregated data.
- Returning full check-in detail payloads or per-match aggregate panels inside the `/me/checkins` list.

## Decisions

### 1. Reuse the existing authentication middleware

The `session-auth` spec already provides the `middleware.ResolveCurrentUser` and `middleware.RequireAuth` functions. The new `/me/profile` and `/me/checkins` routes will be placed inside the same `protected` group as the existing check‑in endpoints.

Why this approach:
- Consistent authentication behavior across all protected endpoints.
- Minimal duplication of middleware setup.
- Single source of truth for user resolution.

Alternatives considered:
- Separate middleware for profile endpoints: rejected because it would duplicate logic already present.

### 2. Introduce dedicated profile and history DTOs

The profile API will return dedicated DTOs rather than extending the lightweight auth user shape. The summary response will include:
- `checkInCount` (number of recorded matches)
- `avgMatchRating` (average of the user’s match ratings)
- `favoriteTeamId` (team ID with most check‑ins)
- `mostUsedTagId` (tag ID most frequently used)
- `recentCheckInCount` (check‑ins in the last 30 days)

Why this approach:
- Keeps `/auth/me` lightweight and stable.
- Avoids coupling profile-specific aggregation fields to the auth contract.
- Makes `/me/profile` and `/me/checkins` explicit about the fields they own.

Alternatives considered:
- Extending `UserSummary`: rejected because it would leak profile aggregation concerns into unrelated auth responses.

### 3. Use the same pagination pattern as other list endpoints

The `GET /me/checkins` endpoint will accept `page` and `pageSize` query parameters, following the same structure as `GET /matches`. The response will include `items`, `page`, `pageSize`, and `total`, matching the existing `MatchListResponseDTO` pattern.

Why this approach:
- Frontend pagination logic stays consistent.
- Backend service and repository methods can follow the same signature as match‑list queries.

Alternatives considered:
- Unlimited‑length history endpoint: rejected because it could return many check‑ins and degrade performance for active users.

The history list will stay intentionally lightweight. Each item will include:
- match identity and display context
- the current user's check-in identity and summary fields
- timestamps needed for ordering and display

It will not include:
- full `playerRatings`
- full tag metadata beyond what the page needs
- match aggregate summary blocks already available on match detail pages

### 4. Keep the `/me` page a standalone page, not nested inside match‑detail

The `/me` page will be a top‑level route that presents the user’s personal statistics and history, separate from the match‑detail page where check‑ins are created/edited.

Why this approach:
- Clear separation between public browsing and personal data.
- Does not overload the match‑detail page with both public and private views.

Alternatives considered:
- Adding profile summary as a panel on match‑detail page: rejected because it would mix personal data with public‑focused UI, reducing clarity.

### 5. Treat release-quality as part of the final spec

This is the last planned v1 spec, so it will also close the remaining integration and quality work needed before calling the loop complete.

That includes:
- `/me` loading, empty, error, and signed-out states
- manual validation of the authenticated journey: `login -> /me -> /matches/:id -> create/edit check-in -> /me`
- re-running the standard backend/frontend validation commands after implementation

Why this approach:
- Prevents a final “cleanup-only” spec after the main product loop is already complete.
- Keeps the acceptance bar for v1 explicit instead of leaving final validation implicit.

## Risks / Trade‑offs

- [Aggregation query performance on large check‑in counts] → Use indexed columns (`user_id`, `match_id`) and keep history payloads lightweight.
- [Overloading `/me/checkins` with detail data] → Keep the list response focused on summary/history needs and defer full detail to match detail or future dedicated endpoints.
- [Missing profile page when auth or check‑in data changes] → Revalidate profile data on page load and after explicit user actions that return to `/me`.
- [No server‑side caching for profile statistics] → v1 uses real‑time calculations for simplicity; caching can be added later if needed.
