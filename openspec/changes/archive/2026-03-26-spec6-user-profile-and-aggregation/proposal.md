## Why

The product loop now supports public match browsing, session auth, and check-in creation/editing, but signed-in users cannot view their personal statistics, check-in history, or aggregated insights about their football viewing experience. The current `/me` page is only a placeholder, breaking the full v1 loop where users should be able to reflect on their recorded matches and see their personal trends.

## What Changes

- Add authenticated user profile summary API (`GET /me/profile`) that returns the current user's v1 profile statistics and base identity fields.
- Add authenticated user check-in history API (`GET /me/checkins`) that returns a paginated list of the user's recorded matches with lightweight match context and check-in summary fields.
- Add backend aggregation queries that calculate the minimal v1 profile statistics needed by the `/me` page.
- Extend the frontend `/me` page to display the user's profile summary and check-in history.
- Update API client and types to support the new profile and history endpoints.
- Add user-related backend service, repository, and handler layers following the existing clean architecture pattern.
- Add final release-quality validation for the authenticated user journey and `/me` page states.

## Capabilities

### New Capabilities
- `user-profile-and-release-quality`: Authenticated user profile summary, check-in history, and final release-quality validation.

### Modified Capabilities
<!-- Existing capabilities whose REQUIREMENTS are changing (not just implementation).
     Only list here if spec-level behavior changes. Each needs a delta spec file.
     Use existing spec names from openspec/specs/. Leave empty if no requirement changes. -->
- `session-auth`: The authenticated API surface expands to include `/me/profile` and `/me/checkins` endpoints, but the core authentication requirements remain unchanged.

## Impact

- Backend APIs: `GET /me/profile`, `GET /me/checkins` (both authenticated)
- Backend systems: new user handler, service, and repository layers; aggregation query logic; dedicated profile/history DTOs
- Frontend pages: `/me` page redesign with profile summary and check-in history display
- Frontend systems: API client updates, new types, loading/error state handling for profile data
- Database: existing tables (`users`, `check_ins`, `player_ratings`, `checkin_tags`) will be queried for aggregation; no new schema changes required
- Validation: final authenticated journey and page-state checks will be re-run before closing v1
- Dependencies: relies on completed check-in domain (`checkin-domain-and-api`), check-in UI integration (`checkin-ui-and-integration`), and session auth (`session-auth`) specs
