## MODIFIED Requirements

### Requirement: Public match detail browsing
The system SHALL provide a public match detail endpoint and page with base match data, match-scoped roster data, and public aggregate summaries.

#### Scenario: View match detail
- **WHEN** a user opens a specific public match detail page
- **THEN** the API SHALL return the match base information for that match
- **AND** the response SHALL include the match-scoped roster needed for check-in player selection
- **AND** the response SHALL include aggregate summary fields needed by the detail page

#### Scenario: View match detail with no aggregate data
- **WHEN** a match has no check-ins or no aggregate data yet
- **THEN** the API SHALL still return the base match detail successfully
- **AND** the match-scoped roster SHALL still be present for that match
- **AND** aggregate score fields MAY be `null`
- **AND** aggregate list fields SHALL be empty arrays

#### Scenario: Match detail not found
- **WHEN** a client requests a match ID that does not exist
- **THEN** the API SHALL return `404 NOT_FOUND`
- **AND** the frontend SHALL render a not-found state for that route

### Requirement: Match detail public aggregates
The match detail response SHALL include the public aggregate information and roster data defined for v1 browsing.

#### Scenario: Match aggregate summary
- **WHEN** a match detail response is returned
- **THEN** it SHALL include match rating average, home team rating average, away team rating average, and check-in count
- **AND** average values SHALL follow the existing one-decimal aggregation rule

#### Scenario: Player rating leaderboard
- **WHEN** player ratings exist for a match
- **THEN** the match detail response SHALL include a public player rating summary list
- **AND** each item SHALL include the player identity and aggregate rating information needed for ranking display

#### Scenario: Recent public reviews
- **WHEN** public reviews exist for a match
- **THEN** the match detail response SHALL include a recent review list
- **AND** each review item SHALL include the public fields needed for display without requiring authentication

#### Scenario: Match roster for check-in selection
- **WHEN** a match detail response is returned
- **THEN** it SHALL include the players linked to that match through `match_players`
- **AND** each roster item SHALL include the player identity and team identity needed for frontend player-rating selection
