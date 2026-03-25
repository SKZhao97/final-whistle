## Why

Final Whistle has completed the foundation layer, but users still cannot browse any meaningful football content. The next step is to expose the first user-visible read path so the product moves from "bootstrapped project" to "usable browsing experience."

## What Changes

- Add public read APIs for match list, match detail, team detail, and player detail
- Add match detail aggregation reads for summary ratings, check-in count, player rating leaderboard, and recent reviews
- Build public frontend pages for `/matches`, `/matches/[matchId]`, `/teams/[teamId]`, and `/players/[playerId]`
- Add loading, empty, not-found, and error states for the public browsing path
- Keep the entire change read-only: no login, no `my-checkin`, no create/update check-in flow

## Capabilities

### New Capabilities
- `public-match-browsing`: Public read-only match, team, and player browsing across backend APIs and frontend pages

### Modified Capabilities

## Impact

- **Backend impact**: Adds read handlers, services, repositories, DTOs, and aggregation queries
- **Frontend impact**: Replaces placeholder app-shell routes with real public pages and API integrations
- **API impact**: Introduces public endpoints for matches, teams, and players
- **Database impact**: Uses existing seed data and schema; no schema changes are required
