## ADDED Requirements

### Requirement: Public match browsing is organized by season and round
The system SHALL present public match browsing as a season-first experience where the latest available season is selected by default and fixtures are grouped by round within the selected season.

#### Scenario: Browse the latest season by default
- **WHEN** a user opens the public `/matches` page
- **THEN** the page SHALL select the latest available season from the public match data by default
- **AND** the selected season SHALL contain grouped rounds rather than a single undifferentiated fixture stream

#### Scenario: Switch seasons explicitly
- **WHEN** the public `/matches` page contains fixtures from more than one season
- **THEN** the page SHALL provide a season-selection control
- **AND** choosing a season SHALL update the grouped rounds and fixture cards to that season without changing the destination route shape

#### Scenario: Sort seasons and rounds predictably
- **WHEN** the public matches page contains fixtures from multiple seasons or rounds
- **THEN** season sections SHALL be ordered from newest to oldest
- **AND** round groups within a season SHALL be ordered in a stable football-readable progression

#### Scenario: Group a match without round information
- **WHEN** a public match has no round value available
- **THEN** the system SHALL place that fixture into an explicit fallback round group instead of silently merging it into another round
- **AND** the fallback group label SHALL remain locale-aware for the browsing surface

#### Scenario: Browse fixtures within a round
- **WHEN** the selected season contains matches for one or more rounds
- **THEN** the page SHALL render a round header for each round group
- **AND** each round group SHALL list the fixtures that belong to that round in a stable order

### Requirement: Grouped match browsing remains a match-detail discovery surface
The system SHALL use the grouped browsing page to guide users into match detail without replacing match detail as the main public fixture destination.

#### Scenario: Open a fixture from a grouped round
- **WHEN** a user selects a fixture card inside a season-round group
- **THEN** the system SHALL navigate to the existing public match detail page for that fixture
- **AND** the grouped browsing structure SHALL not require a new match-detail route shape

### Requirement: Season-round browsing leaves room for future team-based discovery
The system SHALL structure grouped match browsing so future team-based browsing can reuse the same fixture-card and grouping foundations.

#### Scenario: Future browsing extensions reuse grouped list primitives
- **WHEN** a later change introduces team-based browsing or regrouping
- **THEN** the current season-round browsing model SHALL expose reusable grouping and fixture presentation concepts
- **AND** the current change SHALL NOT hard-code copy or structure that prevents another grouping dimension from being introduced later

### Requirement: Grouped browsing prefers coherent sections over flat pagination
The system SHALL favor complete season-round sections for the public grouped browsing view instead of preserving the current flat pagination semantics.

#### Scenario: Render complete grouped browsing data
- **WHEN** the public `/matches` page is rendered for the grouped browsing experience
- **THEN** it SHALL load the data required to render the season selector and complete round groups for the selected season
- **AND** it SHALL NOT split a round across flat list pages in the grouped browsing UI

#### Scenario: Render grouped empty states
- **WHEN** there are no public matches available for grouped browsing
- **THEN** the system SHALL render an empty state for the `/matches` page
- **AND** it SHALL avoid rendering empty season or round shells without explaining the absence of fixtures
