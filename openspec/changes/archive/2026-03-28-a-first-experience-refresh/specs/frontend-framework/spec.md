## ADDED Requirements

### Requirement: Shared football presentation primitives
The frontend SHALL provide reusable presentation primitives for the refreshed football experience so core pages can share hierarchy, branding treatment, and visual language.

#### Scenario: Match hero primitives are reusable
- **WHEN** match detail and related football surfaces render match-context hero sections
- **THEN** they SHALL be able to use shared primitives for fixed team crest treatment, secondary league-brand placement, and balanced field-accent styling

#### Scenario: Archive surfaces share a common visual system
- **WHEN** the match record surface and `/me` archive sections are rendered
- **THEN** they SHALL be able to use shared section, card, and typography patterns that support archive/editorial presentation
- **AND** those primitives SHALL remain compatible with the sitewide locale and localized text lengths
