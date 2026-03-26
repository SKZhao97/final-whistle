## Why

The product loop is still missing its core user action: recording a unique post-match check-in. Public browsing and session auth are already in place, so the next step is to establish the backend write contract and domain rules for creating, updating, and retrieving a user’s check-in for a finished match.

## What Changes

- Add authenticated check-in read and write APIs for `GET /matches/:id/my-checkin`, `POST /matches/:id/checkin`, and `PUT /matches/:id/checkin`.
- Add backend domain validation for finished-match enforcement, one-check-in-per-user-per-match, player eligibility, tag legality, score ranges, and payload size limits.
- Add transactional persistence for `check_ins`, `player_ratings`, and `checkin_tags`, including replacement behavior for updates.
- Add backend response DTOs that return the current user’s full check-in record so later UI integration can reuse the same contract.

## Capabilities

### New Capabilities
- `checkin-domain-and-api`: Authenticated domain rules and backend APIs for reading, creating, and updating a user’s unique match check-in.

### Modified Capabilities

## Impact

- Backend APIs: `GET /matches/:id/my-checkin`, `POST /matches/:id/checkin`, `PUT /matches/:id/checkin`
- Backend systems: auth-protected handlers, transactional write flow, check-in repositories, validation logic, and DTO mapping
- Database tables already activated by this change: `check_ins`, `player_ratings`, `checkin_tags`, `match_players`, `tags`
- Later dependencies: match detail UI integration, personal profile history, and aggregate readbacks all depend on this domain contract being stable
