## MODIFIED Requirements

### Requirement: Shared frontend application state
The frontend SHALL provide shared application-level state for both authentication recovery and current locale handling.

#### Scenario: Initialize locale state
- **WHEN** the frontend application initializes
- **THEN** it SHALL resolve the current locale from persisted user preference or the product default
- **AND** it SHALL make that locale available through a shared app-level mechanism

#### Scenario: Update locale state
- **WHEN** the user changes language from the supported switcher
- **THEN** the shared frontend state SHALL update immediately
- **AND** components that consume translated UI copy SHALL re-render in the new locale

### Requirement: API client infrastructure
The frontend SHALL have a structured API client for communicating with the backend.

#### Scenario: Locale-aware requests
- **WHEN** the frontend needs locale-aware backend responses
- **THEN** the shared client pattern SHALL preserve the user’s persisted locale context with requests
- **AND** it SHALL support follow-up refresh of current-page data after a locale change

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
