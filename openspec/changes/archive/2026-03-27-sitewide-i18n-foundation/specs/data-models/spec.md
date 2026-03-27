## MODIFIED Requirements

### Requirement: Proper field mappings
Each model field SHALL map correctly to its corresponding database column.

#### Scenario: Localized tag model fields
- **WHEN** the `Tag` model is used with GORM
- **THEN** it SHALL include fields for both English and Chinese display names
- **AND** database columns SHALL map consistently to those localized values

### Requirement: Repository and DTO mapping support
The backend SHALL provide model and mapping support for locale-aware tag output.

#### Scenario: Locale-aware tag projection
- **WHEN** the backend maps tag entities into API responses
- **THEN** it SHALL choose the correct display name for the current locale
- **AND** it SHALL continue to expose a stable DTO shape to frontend consumers
