## ADDED Requirements

### Requirement: Comprehensive seed data coverage
The seed data SHALL provide complete test data for all core entities with realistic relationships.

#### Scenario: Data completeness
- **WHEN** seed data is loaded
- **THEN** there SHALL be data for: teams, players, matches, match-player relationships, tags, and optional dev users
- **AND** all foundational relationships SHALL be properly established

### Requirement: Realistic domain data
Seed data SHALL reflect realistic football scenarios and relationships.

#### Scenario: Domain realism
- **WHEN** examining seed data
- **THEN** teams SHALL represent actual football clubs (or realistic analogues)
- **AND** matches SHALL have plausible scores and timestamps
- **AND** player positions SHALL be appropriate for their roles

### Requirement: Modular seed scripts
Seed data SHALL be organized into modular scripts that can be run independently.

#### Scenario: Script modularity
- **WHEN** running seed scripts
- **THEN** there SHALL be separate scripts for: base data (teams, players, matches), tags, and optional users
- **AND** scripts SHALL be runnable individually or together

### Requirement: Test data for development scenarios
Seed data SHALL support common development and testing scenarios.

#### Scenario: Development support
- **WHEN** developing features
- **THEN** there SHALL be enough seeded matches and roster data to support later auth and check-in specs
- **AND** feature-specific sample check-ins MAY be introduced in later specs

### Requirement: Predefined tag dictionary
The system SHALL include a standard set of emotion and impression tags.

#### Scenario: Tag dictionary
- **WHEN** tags are seeded
- **THEN** there SHALL be predefined tags: `热血`, `无聊`, `窒息`, `经典`, `离谱`, `可惜`, `统治力`, `折磨`, `逆转`, `宿命感`
- **AND** each tag SHALL have a unique slug and sort order

### Requirement: Reset and re-seed capability
The seed process SHALL support resetting the database and reloading data.

#### Scenario: Reset functionality
- **WHEN** resetting the database
- **THEN** existing data SHALL be cleared (except system tables)
- **AND** base seed scripts SHALL recreate foundational data
- **AND** the process SHALL be idempotent (can be run multiple times with same result)

### Requirement: Seed data validation
Seed data SHALL pass basic validation checks for data integrity.

#### Scenario: Data validation
- **WHEN** seed data is loaded
- **THEN** all foreign key relationships SHALL be valid
- **AND** all enum values SHALL be within defined ranges
- **AND** all text fields SHALL be within length limits
