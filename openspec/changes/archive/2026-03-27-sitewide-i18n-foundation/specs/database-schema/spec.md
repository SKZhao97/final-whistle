## MODIFIED Requirements

### Requirement: Complete table schema for all core entities
The database SHALL have tables for all core entities defined in the technical design.

#### Scenario: Localized tag dictionary storage
- **WHEN** migrations are run for the tag dictionary
- **THEN** the `tags` table SHALL store both English and Chinese display names for each predefined tag
- **AND** the schema SHALL continue to preserve unique slug-based identity for each tag
