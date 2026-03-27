## MODIFIED Requirements

### Requirement: Predefined tag dictionary
The system SHALL include a standard set of emotion and impression tags.

#### Scenario: Bilingual tag dictionary
- **WHEN** tags are seeded
- **THEN** each predefined tag SHALL have a unique slug, sort order, and both English and Chinese display names
- **AND** reseeding SHALL update existing local tag records to the current bilingual dictionary values
