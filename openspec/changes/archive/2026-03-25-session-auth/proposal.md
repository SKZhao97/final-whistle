## Why

The project already supports public match browsing, but all user-specific paths still lack a stable identity layer. A minimal session-based authentication spec is needed now so later check-in and profile work can depend on a clear login model instead of re-deciding auth behavior inside each feature.

## What Changes

- Add session-based authentication APIs for login, logout, and current-user lookup.
- Add backend session persistence, cookie issuance, and protected-route middleware for later authenticated features.
- Add a frontend auth client flow and app-level login state recovery so protected pages can identify signed-in users.
- Define the v1 development login behavior explicitly, using internal seeded users instead of external OAuth providers and limiting dev-user auto-create behavior to development mode.

## Capabilities

### New Capabilities
- `session-auth`: Minimal cookie-session authentication for Final Whistle v1, including login, logout, current-user lookup, and auth state recovery.

### Modified Capabilities
- `frontend-framework`: Add the requirement that the shared frontend foundation can restore auth state and route authenticated requests through the established API client pattern.

## Impact

- Backend APIs: `POST /auth/login`, `POST /auth/logout`, `GET /auth/me`
- Backend systems: session storage, auth middleware, seeded dev-user login flow
- Frontend systems: auth API helpers, auth state bootstrap, protected navigation behavior
- Dependencies: existing `users` and `sessions` tables become active runtime infrastructure for later authenticated features
