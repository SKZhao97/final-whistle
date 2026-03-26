## MODIFIED Requirements

### Requirement: Check-in payload validation
The backend SHALL enforce the v1 check-in payload rules before any write succeeds.

#### Scenario: Validate score ranges
- **WHEN** a check-in create or update payload includes match, team, or player ratings outside the allowed `1-10` range
- **THEN** the API SHALL reject the request with a validation error

#### Scenario: Validate review and note lengths
- **WHEN** a payload includes `shortReview` longer than `280` characters or a player rating note longer than `80` characters
- **THEN** the API SHALL reject the request with a validation error

#### Scenario: Validate tag legality
- **WHEN** a payload includes tag IDs that do not exist or are not active
- **THEN** the API SHALL reject the request with a validation error

#### Scenario: Validate player eligibility
- **WHEN** a payload includes player IDs that do not belong to the target match through `match_players`
- **THEN** the API SHALL reject the request with a validation error

#### Scenario: Allow full match roster rating selection
- **WHEN** a payload includes any number of distinct player ratings whose player IDs all belong to the target match through `match_players`
- **THEN** the backend SHALL evaluate them against the normal rating and text validation rules
- **AND** it SHALL NOT reject the payload solely because more than five players are included
