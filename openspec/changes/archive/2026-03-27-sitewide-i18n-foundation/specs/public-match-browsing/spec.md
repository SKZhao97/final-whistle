## MODIFIED Requirements

### Requirement: Match detail public aggregates
The match detail response SHALL include the public aggregate information and roster data defined for v1 browsing.

#### Scenario: Locale-aware public tag labels
- **WHEN** a match detail response includes recent public reviews with tags
- **THEN** each tag item SHALL expose the display `name` for the current locale
- **AND** the response shape SHALL remain stable for frontend consumers
