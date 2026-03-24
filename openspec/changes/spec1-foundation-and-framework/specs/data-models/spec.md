## ADDED Requirements

### Requirement: Complete GORM model definitions
The backend SHALL have Go struct definitions for all core entities using GORM tags.

#### Scenario: Model creation
- **WHEN** models are defined
- **THEN** there SHALL be Go structs for: `User`, `Team`, `Player`, `Match`, `MatchPlayer`, `CheckIn`, `PlayerRating`, `Tag`, `CheckInTag`, `Session`
- **AND** each struct SHALL have appropriate GORM tags for field mapping

### Requirement: Proper field mappings
Each model field SHALL map correctly to its corresponding database column.

#### Scenario: Field mapping
- **WHEN** models are used with GORM
- **THEN** exported Go struct field names SHALL use `PascalCase`
- **AND** database columns SHALL map to `snake_case`
- **AND** primary keys SHALL be properly tagged with `gorm:"primaryKey"`
- **AND** timestamps SHALL be automatically managed with `gorm:"autoCreateTime"` and `gorm:"autoUpdateTime"`

### Requirement: Relationship definitions
Models SHALL define proper relationships (belongs to, has many, many to many) according to the domain model.

#### Scenario: Relationship configuration
- **WHEN** models are defined
- **THEN** `CheckIn` SHALL belong to `User` and `Match`
- **AND** `CheckIn` SHALL have many `PlayerRating` and many `Tag` through `CheckInTag`
- **AND** `Match` SHALL belong to `home_team` and `away_team` (both referencing `Team`)
- **AND** `Player` SHALL belong to `Team`
- **AND** `MatchPlayer` SHALL belong to `Match` and `Player` and `Team`

### Requirement: Business enum types
Business enumerations SHALL be defined as Go types with proper validation.

#### Scenario: Enum type definition
- **WHEN** enums are defined
- **THEN** there SHALL be Go types for: `MatchStatus`, `WatchedType`, `SupporterSide`
- **AND** each enum SHALL implement database serialization/deserialization methods
- **AND** each enum SHALL have validation methods

### Requirement: Model validation
Models SHALL include validation logic for business rules.

#### Scenario: Validation implementation
- **WHEN** models are validated
- **THEN** `CheckIn` SHALL validate that `match_rating`, `home_team_rating`, `away_team_rating` are between 1-10
- **AND** `CheckIn` SHALL validate that `short_review` length does not exceed 280 characters
- **AND** `PlayerRating` SHALL validate that `rating` is between 1-10 and `note` does not exceed 80 characters

### Requirement: Repository interface skeletons
The backend SHALL define only the minimal shared repository structure needed to support foundational database access.

#### Scenario: Repository interfaces
- **WHEN** data access layer is defined
- **THEN** there SHALL be a repository package structure that later specs can extend
- **AND** foundation code SHALL avoid locking feature specs into premature CRUD interfaces
