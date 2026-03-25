## 1. Backend Auth APIs

- [x] 1.1 Add auth and session DTOs for login, logout, current-user, and unauthorized responses
- [x] 1.2 Implement session repository support for create, lookup, delete, and expiry-aware session resolution
- [x] 1.3 Implement auth service logic for dev login, logout, and current-user lookup using the existing users and sessions tables, with development-only auto-create behavior
- [x] 1.4 Add auth middleware that resolves the current user from the session cookie and exposes a strict guard for protected routes
- [x] 1.5 Implement `POST /auth/login`, `POST /auth/logout`, and `GET /auth/me` handlers with environment-sensitive cookie issuance, clearing, and stale-session handling
- [x] 1.6 Register auth routes and wire the new middleware into the backend router for later protected modules
- [x] 1.7 Add backend tests covering login success, auto-created dev user flow, invalid payload, unauthorized `/auth/me`, logout, and expired-session handling

## 2. Frontend Auth State

- [x] 2.1 Extend the shared API client to support cookie-based auth requests and auth endpoint helpers
- [x] 2.2 Add frontend auth request and response types for login and current-user state
- [x] 2.3 Implement frontend auth-state recovery using `GET /auth/me` during app initialization
- [x] 2.4 Add a minimal development login entry point and signed-in/signed-out state handling in the shared app shell without locking the implementation to a specific state library
- [x] 2.5 Add unauthorized handling paths so later protected pages can reuse the shared auth-state behavior

## 3. Integration & Validation

- [x] 3.1 Verify seeded dev users and optional auto-create login behavior against the local database
- [x] 3.2 Manually validate the browser flow: login -> refresh with restored session -> logout -> signed-out state
- [x] 3.3 Validate backend `go test ./...` and `go build ./...` after auth implementation
- [x] 3.4 Validate frontend lint, type-check, and build after auth-state integration
