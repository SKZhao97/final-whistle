## ADDED Requirements

### Requirement: Session-based login
The system SHALL support development login using the local user table and a database-backed cookie session.

#### Scenario: Login with existing dev user
- **WHEN** a client submits `POST /auth/login` with a valid development login payload for an existing user
- **THEN** the API SHALL return success with the user summary needed by the frontend
- **AND** the backend SHALL create a session record linked to that user
- **AND** the response SHALL set an `HTTP-only` session cookie

#### Scenario: Login with new dev user
- **WHEN** a client submits `POST /auth/login` in development mode with a development login payload for an email that does not yet exist
- **THEN** the backend SHALL create a local user record for development use
- **AND** it SHALL create a session record for that new user
- **AND** it SHALL return the same success shape as an existing-user login

#### Scenario: Unknown user outside development mode
- **WHEN** a client submits `POST /auth/login` for an email that does not yet exist while auto-create behavior is disabled
- **THEN** the API SHALL reject the request with an authentication or validation error
- **AND** it SHALL NOT create a user or a session

#### Scenario: Invalid login payload
- **WHEN** a client submits `POST /auth/login` with missing or invalid required fields
- **THEN** the API SHALL reject the request with a validation error response
- **AND** it SHALL NOT create a user or a session

### Requirement: Session termination
The system SHALL support explicit logout by invalidating the current session and clearing the session cookie.

#### Scenario: Logout with active session
- **WHEN** an authenticated client submits `POST /auth/logout`
- **THEN** the backend SHALL delete or invalidate the current session record
- **AND** the response SHALL clear the session cookie
- **AND** the API SHALL return a success response

#### Scenario: Logout without active session
- **WHEN** a client submits `POST /auth/logout` without a valid active session
- **THEN** the API SHALL still return a success response
- **AND** the response SHALL clear any session cookie value on the client

### Requirement: Session cookie behavior
The system SHALL issue and clear session cookies in a way that supports local development and later deployment safely.

#### Scenario: Issue session cookie on login
- **WHEN** login succeeds
- **THEN** the response SHALL set a session cookie with `HttpOnly`
- **AND** the cookie SHALL use `Path=/`
- **AND** the cookie SHALL use `SameSite=Lax`
- **AND** the `Secure` flag SHALL be environment-sensitive so local HTTP development remains usable

#### Scenario: Clear session cookie on logout
- **WHEN** logout succeeds or logout is called without a valid active session
- **THEN** the response SHALL clear the session cookie using the same cookie scope needed to remove the client value

### Requirement: Current-user lookup
The system SHALL expose the current authenticated user through `GET /auth/me`.

#### Scenario: Resolve current user from valid session
- **WHEN** a client calls `GET /auth/me` with a valid unexpired session cookie
- **THEN** the API SHALL return the current user summary
- **AND** it SHALL NOT expose the session token in the JSON payload

#### Scenario: Missing or invalid session on current-user lookup
- **WHEN** a client calls `GET /auth/me` without a valid active session
- **THEN** the API SHALL return `401 UNAUTHORIZED`
- **AND** it SHALL NOT return user data

#### Scenario: Expired session on current-user lookup
- **WHEN** a client calls `GET /auth/me` with an expired session cookie
- **THEN** the backend SHALL treat the session as unauthorized
- **AND** it SHALL NOT authenticate the request
- **AND** it MAY clear the stale cookie in the response

### Requirement: Auth middleware for protected routes
The backend SHALL provide reusable auth middleware for later protected APIs.

#### Scenario: Attach current user to request context
- **WHEN** a request includes a valid active session cookie
- **THEN** the auth layer SHALL resolve the associated user
- **AND** it SHALL make that user available in request context for handlers and services

#### Scenario: Reject protected request without valid session
- **WHEN** a protected API route is accessed without a valid active session
- **THEN** the middleware SHALL return `401 UNAUTHORIZED`
- **AND** the protected handler SHALL NOT run

### Requirement: Frontend auth state recovery
The frontend SHALL restore and expose auth state using the established session APIs.

#### Scenario: Restore signed-in state
- **WHEN** the frontend initializes and the browser already has a valid session cookie
- **THEN** it SHALL call `GET /auth/me`
- **AND** it SHALL populate the app’s current-user state from the API response

#### Scenario: Restore signed-out state
- **WHEN** the frontend initializes without a valid session
- **THEN** `GET /auth/me` SHALL resolve to unauthorized
- **AND** the frontend SHALL treat the user as signed out without rendering a stale authenticated state

#### Scenario: Minimal login entry
- **WHEN** a user needs to start a development login flow from the frontend
- **THEN** the application SHALL provide at least one minimal login entry point
- **AND** that entry point SHALL allow submitting the login payload and observing signed-in or signed-out state

#### Scenario: Authenticated API request from frontend
- **WHEN** the frontend calls an authenticated API after login
- **THEN** the request path SHALL include the session cookie automatically
- **AND** the shared client flow SHALL support unauthorized handling for later protected pages
