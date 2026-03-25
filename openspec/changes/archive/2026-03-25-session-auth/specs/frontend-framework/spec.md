## MODIFIED Requirements

### Requirement: API client infrastructure
The frontend SHALL have a structured API client for communicating with the backend.

#### Scenario: API client initialization
- **WHEN** the frontend needs to make API calls
- **THEN** there SHALL be a centralized API client with proper error handling
- **AND** it SHALL support cookie-based authenticated requests for later protected APIs

#### Scenario: Type-safe API calls
- **WHEN** making API requests
- **THEN** the API client SHALL use TypeScript types for shared response envelopes and foundational requests
- **AND** endpoint-specific request and response types MAY be introduced in later feature specs

#### Scenario: Auth state recovery
- **WHEN** the frontend needs to restore login state
- **THEN** it SHALL be able to call the backend current-user endpoint through the shared API client
- **AND** the shared client pattern SHALL support signed-in and unauthorized responses without exposing the session token to client code
