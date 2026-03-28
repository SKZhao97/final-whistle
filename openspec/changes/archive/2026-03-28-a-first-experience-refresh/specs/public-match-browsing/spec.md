## ADDED Requirements

### Requirement: Match detail A-first hierarchy
The public match detail page SHALL present the match as a recording-first surface with a clear hero, a primary personal-record section, and a secondary community layer.

#### Scenario: Match detail prioritizes context before interaction
- **WHEN** any user opens a public match detail page
- **THEN** the first major content block SHALL be a match-context hero containing the participating teams, score or status, competition context, and fixed crest treatment
- **AND** the page SHALL place the community-derived content after the `My Match Record` section rather than before it

#### Scenario: Match detail uses secondary league branding
- **WHEN** the match-context hero displays competition information
- **THEN** league branding SHALL appear adjacent to the competition label as a secondary identifier
- **AND** it SHALL NOT displace team identity, score/status, or the record surface as the hero’s dominant visual elements
