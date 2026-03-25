## 1. Backend Read APIs

- [x] 1.1 Add match read DTOs for list and detail responses, including aggregate summary shapes
- [x] 1.2 Implement repository queries for public match listing with supported filters and pagination
- [x] 1.3 Implement repository/service queries for match detail base data, aggregate summaries, player rating leaderboard, and recent reviews
- [x] 1.4 Implement `GET /matches` and `GET /matches/:id` handlers with standard success and error envelopes
- [x] 1.5 Implement repository/service/handler flow for `GET /teams/:id`
- [x] 1.6 Implement repository/service/handler flow for `GET /players/:id`
- [x] 1.7 Add handler and service tests covering success, empty, and not-found cases for the public read APIs

## 2. Frontend Public Pages

- [x] 2.1 Replace the placeholder `/matches` page with a public match list page backed by the API client
- [x] 2.2 Add the `/matches/[matchId]` route and render match base information plus aggregate sections
- [x] 2.3 Add the `/teams/[teamId]` route for public team detail display
- [x] 2.4 Add the `/players/[playerId]` route for public player detail display
- [x] 2.5 Add route-level loading, empty, error, and not-found handling for the public browsing pages
- [x] 2.6 Add route navigation between match, team, and player pages using the shared app shell

## 3. Integration & Validation

- [x] 3.1 Extend frontend API client helpers and page-level types for the new public read endpoints
- [x] 3.2 Verify the seeded database supports the new public browsing pages end-to-end
- [x] 3.3 Validate backend build and frontend lint/type-check/build after the public browsing implementation
- [x] 3.4 Manually verify the public browsing path: `/matches` -> match detail -> related team/player pages
