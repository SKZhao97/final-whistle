## ADDED Requirements

### Requirement: User profile summary
The system SHALL provide an authenticated endpoint for retrieving current user profile summary data.

#### Scenario: Load profile summary
- **WHEN** an authenticated user calls `GET /me/profile`
- **THEN** the response SHALL include the current user's base identity fields needed by the `/me` page
- **AND** it SHALL include: `checkInCount`, `avgMatchRating`, `favoriteTeamId`, `mostUsedTagId`, and `recentCheckInCount`
- **AND** it SHALL use a dedicated profile response shape rather than reusing the lightweight auth user DTO

#### Scenario: Profile statistics are aggregated in real‑time
- **WHEN** the user has created or updated recorded matches
- **THEN** a fresh `GET /me/profile` response SHALL reflect the latest committed check‑ins
- **AND** the backend SHALL calculate the summary without requiring precomputed cache tables

### Requirement: User check‑in history
The system SHALL provide an authenticated paginated endpoint for retrieving current user’s recorded matches with match context.

#### Scenario: Load paginated check‑in history
- **WHEN** an authenticated user calls `GET /me/checkins` with valid `page` and `pageSize` parameters
- **THEN** the response SHALL include: `items`, `page`, `pageSize`, and `total`
- **AND** each item SHALL contain: match ID, match context, and the current user's check-in summary fields needed for history display

#### Scenario: Match context is included with check‑in records
- **WHEN** a user’s check‑in history is loaded
- **THEN** each item SHALL include: match competition, season, round, kickoff time, home team, away team, and final scores
- **AND** it SHALL NOT require per-match aggregate summary blocks that are already available on match detail pages

#### Scenario: History list stays lightweight
- **WHEN** a user’s check‑in history is returned
- **THEN** the response SHALL NOT include full `playerRatings` collections for every list item
- **AND** it SHALL remain usable as a profile/history list rather than a full check-in detail feed

### Requirement: Backend user‑profile service layer
The backend SHALL implement user‑profile‑specific service, repository, and handler layers matching the existing clean‑architecture pattern.

#### Scenario: Add user‑profile handler
- **WHEN** the user‑profile handler is added
- **THEN** it SHALL follow the same pattern as existing handlers (`auth_handler.go`, `checkin_handler.go`, etc.)
- **AND** it SHALL be registered in the protected route group along with check‑in endpoints

#### Scenario: Add user‑profile service
- **WHEN** the user‑profile service is added
- **THEN** it SHALL follow the same interface style as existing services (`AuthService`, `CheckInService`, etc.)
- **AND** when called with a valid user ID it SHALL return the corresponding profile summary or check-in history DTOs

#### Scenario: Add user‑profile repository
- **WHEN** the user‑profile repository is added
- **WHEN** called with the user’s ID
- **THEN** it SHALL retrieve the user’s profile data and aggregate statistics from the database

### Requirement: Release-quality validation
The final v1 spec SHALL include validation of the authenticated profile journey and `/me` page states.

#### Scenario: Signed-out profile page
- **WHEN** an unauthenticated user opens `/me`
- **THEN** the page SHALL render a signed-out state with a clear route to log in

#### Scenario: Profile page loading, empty, and error states
- **WHEN** the `/me` page is waiting for data, has no check-in history yet, or a profile request fails
- **THEN** the page SHALL render explicit loading, empty, and error states instead of a blank or misleading layout

#### Scenario: Final authenticated journey validation
- **WHEN** the implementation is considered complete
- **THEN** the team SHALL verify the authenticated v1 journey end-to-end
- **AND** it SHALL include `login -> /me -> /matches/:id -> create or edit check-in -> /me`
