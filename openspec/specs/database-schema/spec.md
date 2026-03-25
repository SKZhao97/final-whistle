## ADDED Requirements

### Requirement: Complete table schema for all core entities
The database SHALL have tables for all core entities defined in the technical design.

#### Scenario: Table creation
- **WHEN** migrations are run
- **THEN** all core tables SHALL be created: `users`, `teams`, `players`, `matches`, `match_players`, `check_ins`, `player_ratings`, `tags`, `checkin_tags`, `sessions`
- **AND** each table SHALL have appropriate columns for the entity's attributes

### Requirement: Proper data types and constraints
Each table SHALL use appropriate PostgreSQL data types and enforce data integrity constraints.

#### Scenario: Data type selection
- **WHEN** tables are created
- **THEN** primary keys SHALL use `BIGSERIAL` for auto-incrementing IDs
- **AND** timestamps SHALL use `TIMESTAMPTZ` for timezone-aware datetime storage
- **AND** string fields SHALL have appropriate length limits (`VARCHAR(n)`)

#### Scenario: Constraint enforcement
- **WHEN** data is inserted or updated
- **THEN** NOT NULL constraints SHALL be enforced where applicable
- **AND** foreign key constraints SHALL maintain referential integrity

### Requirement: Key indexes for performance
The database SHALL have appropriate indexes for common query patterns.

#### Scenario: Index creation
- **WHEN** migrations are run
- **THEN** indexes SHALL be created for:
  - `matches(kickoff_at)` for time-based queries
  - `check_ins(created_at)` for chronological access
  - `player_ratings(player_id)` for player performance analysis
  - `check_ins(user_id, match_id)` unique constraint for one check-in per user per match

### Requirement: Enum consistency
Enumerated values SHALL be consistently represented in the database schema.

#### Scenario: Enum representation
- **WHEN** storing enumerated values
- **THEN** `match.status` SHALL accept `SCHEDULED` and `FINISHED`
- **AND** `check_ins.watched_type` SHALL accept `FULL`, `PARTIAL`, and `HIGHLIGHTS`
- **AND** `check_ins.supporter_side` SHALL accept `HOME`, `AWAY`, and `NEUTRAL`

### Requirement: Migration tooling
Database schema changes SHALL be managed through version-controlled migrations.

#### Scenario: Migration execution
- **WHEN** running migrations
- **THEN** each migration SHALL be applied in order
- **AND** rollback mechanisms SHALL be available for development
