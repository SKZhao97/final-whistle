## 1. Backend DTOs and Types

- [x] 1.1 Add dedicated profile response DTOs with base user identity plus checkInCount, avgMatchRating, favoriteTeamId, mostUsedTagId, and recentCheckInCount fields
- [x] 1.2 Add `UserCheckInHistoryItemDTO` with match context and the current user's check-in summary fields needed for history display
- [x] 1.3 Add `UserCheckInHistoryResponseDTO` with pagination fields (items, page, pageSize, total)
- [x] 1.4 Keep existing auth/user summary DTOs unchanged unless a shared field is truly required by both contracts

## 2. Backend Repository Layer

- [x] 2.1 Create `user_repository.go` with `UserRepository` interface
- [x] 2.2 Implement `GetUserProfileSummary(userID uint)` method for aggregation queries
- [x] 2.3 Implement `GetUserCheckInHistory(userID uint, page, pageSize int)` method for paginated history
- [x] 2.4 Reuse existing tag lookup support where needed for profile aggregation instead of duplicating repository methods

## 3. Backend Service Layer

- [x] 3.1 Create `user_service.go` with `UserService` interface
- [x] 3.2 Implement `GetProfileSummary(userID uint)` method that calls repository
- [x] 3.3 Implement `GetCheckInHistory(userID uint, page, pageSize int)` method that calls repository
- [x] 3.4 Add error handling and business logic validation

## 4. Backend Handler Layer

- [x] 4.1 Create `user_handler.go` with `UserHandler` struct
- [x] 4.2 Implement `GetProfile` handler for `/me/profile` endpoint
- [x] 4.3 Implement `GetCheckInHistory` handler for `/me/checkins` endpoint
- [x] 4.4 Add route registration in `cmd/api/main.go` within protected group

## 5. Frontend Types and API Client

- [x] 5.1 Add `UserProfileSummary` interface to `frontend/src/types/api.ts`
- [x] 5.2 Add `UserCheckInHistoryItem` and `UserCheckInHistoryResponse` interfaces
- [x] 5.3 Add `usersApi` object to `frontend/src/lib/api/client.ts` with `profile()` and `checkins()` methods
- [x] 5.4 Keep auth-specific frontend user types separate from profile-specific response types

## 6. Frontend Page Implementation

- [x] 6.1 Update `/frontend/src/app/me/page.tsx` to fetch and display profile summary
- [x] 6.2 Add check-in history table or list with pagination controls
- [x] 6.3 Add loading, error, and empty state handling for profile data
- [x] 6.4 Add UI components for statistics display (match counts, average ratings, etc.)

## 7. Testing

- [x] 7.1 Add unit tests for `user_service.go`
- [x] 7.2 Add unit tests for `user_handler.go`
- [x] 7.3 Add integration tests for `/me/profile` and `/me/checkins` endpoints
- [x] 7.4 Add frontend component tests for profile page updates

## 8. Release Quality

- [x] 8.1 Verify all endpoints return consistent error responses
- [x] 8.2 Test pagination behavior with various page sizes
- [x] 8.3 Ensure profile statistics update after new check-in creation
- [x] 8.4 Validate the authenticated v1 journey end-to-end: `login -> /me -> /matches/:id -> create or edit check-in -> /me`
- [x] 8.5 Run existing backend/frontend validation commands to confirm no regressions
