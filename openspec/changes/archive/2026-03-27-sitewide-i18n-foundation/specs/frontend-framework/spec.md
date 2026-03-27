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
