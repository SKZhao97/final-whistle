## ADDED Requirements

### Requirement: Public match list browsing
The system SHALL provide a public match list endpoint and page so users can browse available football matches without authentication.

#### Scenario: Browse match list
- **WHEN** a user opens the public matches page
- **THEN** the frontend SHALL request a public match list API
- **AND** the API SHALL return seeded matches with basic match information needed for list display

#### Scenario: Filter public match list
- **WHEN** a client supplies supported list query parameters such as competition, season, page, or pageSize
- **THEN** the API SHALL apply those filters and pagination rules
- **AND** the response SHALL follow the shared paginated envelope

#### Scenario: Empty public match list
- **WHEN** no matches satisfy the current list query
- **THEN** the API SHALL return an empty `items` array
- **AND** the frontend SHALL render a clear empty state instead of failing

### Requirement: Public match detail browsing
The system SHALL provide a public match detail endpoint and page with base match data and public aggregate summaries.

#### Scenario: View match detail
- **WHEN** a user opens a specific public match detail page
- **THEN** the API SHALL return the match base information for that match
- **AND** the response SHALL include aggregate summary fields needed by the detail page

#### Scenario: View match detail with no aggregate data
- **WHEN** a match has no check-ins or no aggregate data yet
- **THEN** the API SHALL still return the base match detail successfully
- **AND** aggregate score fields MAY be `null`
- **AND** aggregate list fields SHALL be empty arrays

#### Scenario: Match detail not found
- **WHEN** a client requests a match ID that does not exist
- **THEN** the API SHALL return `404 NOT_FOUND`
- **AND** the frontend SHALL render a not-found state for that route

### Requirement: Match detail public aggregates
The match detail response SHALL include the public aggregate information defined for v1 browsing.

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

### Requirement: Public team and player detail browsing
The system SHALL provide public team and player detail endpoints and pages as read-only destinations from match browsing.

#### Scenario: View team detail
- **WHEN** a user opens a team detail page
- **THEN** the API SHALL return the team base information
- **AND** it SHALL include recent related matches and rating summary data needed for the v1 team page

#### Scenario: View player detail
- **WHEN** a user opens a player detail page
- **THEN** the API SHALL return the player base information
- **AND** it SHALL include recent rated matches and rating summary data needed for the v1 player page

#### Scenario: Team or player detail not found
- **WHEN** a client requests a team or player that does not exist
- **THEN** the API SHALL return `404 NOT_FOUND`
- **AND** the frontend SHALL render a route-level not-found state

### Requirement: Public page state handling
The frontend SHALL handle loading, empty, error, and not-found states consistently across public browsing pages.

#### Scenario: Loading state
- **WHEN** a public page is waiting for API data
- **THEN** the page SHALL render a loading state appropriate to that route

#### Scenario: API error state
- **WHEN** a public API request fails with a non-not-found error
- **THEN** the frontend SHALL render an error state for the route
- **AND** it SHALL avoid rendering incomplete or misleading content

#### Scenario: Navigation across public entities
- **WHEN** a user navigates from a match page to a related team or player page
- **THEN** the frontend SHALL support direct route-based navigation between those public pages
- **AND** no authentication step SHALL be required
