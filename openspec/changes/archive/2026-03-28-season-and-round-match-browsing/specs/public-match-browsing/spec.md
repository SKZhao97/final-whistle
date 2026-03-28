## MODIFIED Requirements

### Requirement: Public match list browsing
The system SHALL provide a public matches list page that supports season-first browsing and routes users into public match detail pages.

#### Scenario: View grouped public match list
- **WHEN** a user opens the public matches page
- **THEN** the system SHALL return and render public match list data needed to organize fixtures by season and round
- **AND** the page SHALL default to the latest available season
- **AND** the page SHALL display locale-aware competition, season, round, kickoff, team, score/status, and aggregate summary information for each fixture

#### Scenario: Switch the selected season
- **WHEN** the public matches page contains fixtures from more than one season
- **THEN** the page SHALL provide a season-selection control
- **AND** selecting a season SHALL update the grouped round and fixture display for that season

#### Scenario: Grouping strategy remains implementation-aligned
- **WHEN** the grouped public matches page is implemented
- **THEN** the system SHALL use a single documented grouping strategy for season and round organization
- **AND** the frontend and backend SHALL NOT diverge on whether the list is grouped client-side or server-side for this change

#### Scenario: Match list remains readable as fixtures grow
- **WHEN** the public matches page contains fixtures from multiple rounds or seasons
- **THEN** the browsing surface SHALL group fixtures into season and round sections instead of a single flat list
- **AND** each fixture card SHALL remain scannable through team identity, score/status, kickoff context, and summary metadata

#### Scenario: Preserve coherent round groups
- **WHEN** grouped public browsing is rendered
- **THEN** the system SHALL prefer complete round groups over flat pagination slices
- **AND** it SHALL avoid showing partial rounds as if they were complete grouped sections

#### Scenario: Open public match detail from grouped list
- **WHEN** a user selects a fixture from any season-round group
- **THEN** the system SHALL navigate to the existing public match detail page for that fixture
- **AND** the grouped list SHALL remain the primary public discovery entry point for match detail pages
