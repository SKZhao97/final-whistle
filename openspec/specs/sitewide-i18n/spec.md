## ADDED Requirements

### Requirement: Sitewide locale selection
The system SHALL allow users to actively choose between English and Chinese from the product surface, including the homepage.

#### Scenario: Homepage language switcher
- **WHEN** a user opens the homepage
- **THEN** the page SHALL expose a visible language selection control
- **AND** that control SHALL allow switching between English and Chinese without leaving the page

#### Scenario: Persist chosen locale
- **WHEN** a user changes the language
- **THEN** the system SHALL persist that choice for future page loads in the same browser
- **AND** later navigations SHALL continue to use the chosen locale until the user changes it again

### Requirement: Sitewide UI copy localization
The frontend SHALL render application UI copy in the currently selected locale across the shipped product surface.

#### Scenario: Localize static UI copy
- **WHEN** the current locale is English or Chinese
- **THEN** page titles, buttons, navigation labels, form labels, helper text, and call-to-action text SHALL render in that locale

#### Scenario: Localize page states
- **WHEN** the current locale changes
- **THEN** loading states, empty states, signed-out states, and error states SHALL render in that locale

#### Scenario: Do not translate user-generated content
- **WHEN** a page renders user short reviews, usernames, or other UGC
- **THEN** the system SHALL display that content as originally stored
- **AND** it SHALL NOT attempt machine translation or locale-specific rewriting

### Requirement: Immediate user-visible locale switching
The system SHALL make locale changes visible immediately from the user’s perspective.

#### Scenario: Switch current page language
- **WHEN** a user switches the locale while viewing a page
- **THEN** the visible UI copy on that page SHALL update during the same interaction flow
- **AND** the user SHALL NOT need to restart the application or reseed data to observe the change

#### Scenario: Refresh locale-aware data blocks
- **WHEN** a page includes server-backed data whose display labels depend on locale
- **THEN** the page SHALL refresh the relevant data needed to reflect the new locale
- **AND** the overall experience SHALL still feel immediate to the user
