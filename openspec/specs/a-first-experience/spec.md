## ADDED Requirements

### Requirement: A-first core surface hierarchy
The product SHALL present its primary authenticated and public football surfaces around an A-first hierarchy that emphasizes personal post-match recording before community reaction.

#### Scenario: Match detail uses A-first ordering
- **WHEN** a user opens a match detail page
- **THEN** the page SHALL present `Match Context` before `My Match Record`
- **AND** it SHALL present `Community Pulse` after the record surface rather than as the page’s primary content block

#### Scenario: Profile page reinforces archive value
- **WHEN** a signed-in user opens `/me`
- **THEN** the page SHALL frame the user’s content as a personal football archive
- **AND** it SHALL prioritize identity, patterns, archive, and memory framing over a generic dashboard layout

### Requirement: Modern football visual language
The frontend SHALL apply a shared visual language for the refreshed experience that combines modern football identity with archive-oriented clarity.

#### Scenario: Match surfaces use football presentation primitives
- **WHEN** the refreshed match detail page is rendered
- **THEN** it SHALL use fixed team crests as part of the primary match-hero structure
- **AND** league branding SHALL appear as secondary presentation near competition text instead of replacing the hero’s main hierarchy

#### Scenario: Refreshed pages avoid heavy dark utility styling
- **WHEN** the refreshed match detail and archive surfaces are rendered
- **THEN** they SHALL use lighter archive/editorial surfaces with balanced field-green accents
- **AND** they SHALL NOT rely on a dominant heavy-dark theme as the defining product style
