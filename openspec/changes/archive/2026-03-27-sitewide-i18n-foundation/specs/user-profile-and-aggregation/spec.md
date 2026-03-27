## MODIFIED Requirements

### Requirement: User check‑in history
The system SHALL provide an authenticated paginated endpoint for retrieving current user’s recorded matches with match context.

#### Scenario: Locale-aware tag labels in history
- **WHEN** a user’s check‑in history is loaded
- **THEN** tag items in the history response SHALL expose the display `name` for the current locale
- **AND** UGC fields such as `shortReview` SHALL remain unchanged
