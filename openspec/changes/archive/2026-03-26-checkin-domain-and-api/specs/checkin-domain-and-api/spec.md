## ADDED Requirements

### Requirement: Protected current-user check-in lookup
The system SHALL provide a protected endpoint for loading the current user’s full check-in record for a specific match.

#### Scenario: Load existing current-user check-in
- **WHEN** an authenticated user calls `GET /matches/:id/my-checkin` for a match where they already have a check-in
- **THEN** the API SHALL return the full current-user check-in detail payload for that match
- **AND** the payload SHALL include tags and player ratings needed for later edit flows

#### Scenario: Load missing current-user check-in
- **WHEN** an authenticated user calls `GET /matches/:id/my-checkin` for a match where they have no check-in
- **THEN** the API SHALL return success with `data: null`
- **AND** it SHALL NOT treat the absence of that user record as a not-found error

#### Scenario: Reject unauthenticated current-user lookup
- **WHEN** an unauthenticated client calls `GET /matches/:id/my-checkin`
- **THEN** the API SHALL return `401 UNAUTHORIZED`

### Requirement: Check-in creation
The system SHALL allow an authenticated user to create one unique check-in for a finished match.

#### Scenario: Create check-in successfully
- **WHEN** an authenticated user submits `POST /matches/:id/checkin` for a finished match with a valid payload and no existing check-in
- **THEN** the backend SHALL create the parent `check_ins` row
- **AND** it SHALL create the associated `player_ratings` and `checkin_tags` rows
- **AND** it SHALL return the full created check-in detail payload

#### Scenario: Reject create for non-finished match
- **WHEN** an authenticated user submits `POST /matches/:id/checkin` for a match whose status is not `FINISHED`
- **THEN** the API SHALL reject the request with a validation or business-rule error
- **AND** it SHALL NOT persist any check-in data

#### Scenario: Reject duplicate create
- **WHEN** an authenticated user submits `POST /matches/:id/checkin` for a match where they already have a check-in
- **THEN** the API SHALL reject the request with a conflict-style error
- **AND** it SHALL NOT create a second check-in for the same `(user, match)` pair

### Requirement: Check-in update
The system SHALL allow an authenticated user to update their existing check-in for a finished match.

#### Scenario: Update check-in successfully
- **WHEN** an authenticated user submits `PUT /matches/:id/checkin` for a finished match where they already have a check-in
- **THEN** the backend SHALL update the parent `check_ins` row
- **AND** it SHALL replace the existing `player_ratings` and `checkin_tags` rows inside the same transaction
- **AND** it SHALL return the full updated check-in detail payload

#### Scenario: Reject update when check-in is missing
- **WHEN** an authenticated user submits `PUT /matches/:id/checkin` for a match where they do not yet have a check-in
- **THEN** the API SHALL reject the request
- **AND** it SHALL NOT create a new check-in implicitly

#### Scenario: Reject update for non-finished match
- **WHEN** an authenticated user submits `PUT /matches/:id/checkin` for a match whose status is not `FINISHED`
- **THEN** the API SHALL reject the request with a validation or business-rule error

### Requirement: Check-in payload validation
The backend SHALL enforce the v1 check-in payload rules before any write succeeds.

#### Scenario: Validate score ranges
- **WHEN** a check-in create or update payload includes match, team, or player ratings outside the allowed `1-10` range
- **THEN** the API SHALL reject the request with a validation error

#### Scenario: Validate review and note lengths
- **WHEN** a payload includes `shortReview` longer than `280` characters or a player rating note longer than `80` characters
- **THEN** the API SHALL reject the request with a validation error

#### Scenario: Validate player rating count
- **WHEN** a payload includes more than `5` player ratings
- **THEN** the API SHALL reject the request with a validation error

#### Scenario: Validate tag legality
- **WHEN** a payload includes tag IDs that do not exist or are not active
- **THEN** the API SHALL reject the request with a validation error

#### Scenario: Validate player eligibility
- **WHEN** a payload includes player IDs that do not belong to the target match through `match_players`
- **THEN** the API SHALL reject the request with a validation error

### Requirement: Transactional check-in persistence
The backend SHALL persist check-in writes atomically across parent and child tables.

#### Scenario: Atomic create
- **WHEN** any step fails while creating `check_ins`, `player_ratings`, or `checkin_tags`
- **THEN** the backend SHALL roll back the entire create transaction
- **AND** it SHALL NOT leave partial rows persisted

#### Scenario: Atomic update
- **WHEN** any step fails while updating the parent check-in or replacing child rows
- **THEN** the backend SHALL roll back the entire update transaction
- **AND** it SHALL preserve the previously committed state

### Requirement: Stable check-in detail DTO
The backend SHALL return a consistent full-detail check-in payload for current-user read, create, and update operations.

#### Scenario: Full-detail response shape
- **WHEN** a current-user check-in payload is returned
- **THEN** it SHALL include the check-in identity, match identity, watched type, supporter side, rating fields, optional short review, watched timestamp, tags, player ratings, and timestamps
- **AND** later frontend edit flows SHALL be able to use the same shape without requiring additional child-resource calls
