## ADDED Requirements

### Requirement: Archive-oriented `/me` page structure
The authenticated `/me` experience SHALL present the user’s content as a personal football archive rather than as a generic stats dashboard.

#### Scenario: `/me` emphasizes archive framing
- **WHEN** a signed-in user opens `/me`
- **THEN** the page SHALL present a clear identity layer followed by archive-oriented content layers such as patterns, archive/history, and memory-oriented highlights
- **AND** it SHALL avoid presenting the page primarily as a generic administrative or analytics dashboard

#### Scenario: Empty archive state still preserves archive meaning
- **WHEN** a signed-in user with no check-ins opens `/me`
- **THEN** the page SHALL communicate that this surface is the user’s football archive
- **AND** it SHALL guide the user back toward recording matches rather than only showing an empty statistics shell

## MODIFIED Requirements

### Requirement: User check‑in history
The system SHALL provide an authenticated paginated endpoint for retrieving current user’s recorded matches with match context.

#### Scenario: Locale-aware tag labels in history
- **WHEN** a user’s check‑in history is loaded
- **THEN** tag items in the history response SHALL expose the display `name` for the current locale
- **AND** UGC fields such as `shortReview` SHALL remain unchanged

#### Scenario: History supports archive-style recall
- **WHEN** a user’s check-in history is rendered in the `/me` experience
- **THEN** each item SHALL provide enough localized match context and user-record context to support archive-style recall of that match
- **AND** the frontend SHALL be able to present each item as a distinct saved record rather than a bare database row
