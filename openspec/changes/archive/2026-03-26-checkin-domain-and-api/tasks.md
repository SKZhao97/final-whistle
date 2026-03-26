## 1. Backend contract and domain wiring

- [x] 1.1 Add check-in request and response DTOs for current-user read, create, and update flows
- [x] 1.2 Add repository methods for loading the current user's check-in detail, checking uniqueness/existence, validating match players, and validating active tags
- [x] 1.3 Add service-layer validation for finished-match-only writes, score ranges, text limits, duplicate player entries, player eligibility, and tag legality
- [x] 1.4 Add transactional service methods for `GET /matches/:id/my-checkin`, `POST /matches/:id/checkin`, and `PUT /matches/:id/checkin`
- [x] 1.5 Add authenticated handlers and route registration for the three check-in endpoints using the existing session-auth middleware

## 2. Transactional persistence and error behavior

- [x] 2.1 Implement check-in create persistence for `check_ins`, `player_ratings`, and `checkin_tags` inside one transaction
- [x] 2.2 Implement check-in update persistence with full replacement of child collections inside one transaction
- [x] 2.3 Return `data: null` for missing current-user check-ins without treating the match resource as missing
- [x] 2.4 Map duplicate-create, missing-update, non-finished-match, and invalid-payload cases to stable API errors
- [x] 2.5 Ensure create and update responses return the full check-in detail DTO needed by later edit flows

## 3. Validation, tests, and verification

- [x] 3.1 Add service and/or repository tests for player eligibility, active-tag validation, and duplicate player detection
- [x] 3.2 Add handler or integration tests for successful current-user read, create, and update flows
- [x] 3.3 Add handler or integration tests for unauthenticated access, non-finished matches, invalid scores/text, too many player ratings, invalid tags, and duplicate create
- [x] 3.4 Add transactional failure coverage to verify partial create/update work is rolled back
- [x] 3.5 Validate the backend with `go test ./...` and a local manual API smoke test against seeded data
