## Why

The backend check-in contract is now in place, but the product loop still breaks on the main match-detail path because signed-in users cannot create or edit their own record from the UI. This change makes the core “record this match” flow usable end to end.

## What Changes

- Add signed-in check-in UI behavior to the match detail page, including loading the current user’s existing check-in state.
- Add a dedicated check-in form flow for creating and editing a match record using the existing authenticated backend APIs.
- Add frontend validation, submit states, and success/error handling for check-in create and update actions.
- Add match-detail UI integration so the page reflects the user’s current record after successful submission without requiring a manual refresh.
- Extend the public match detail contract to include the match-scoped player roster needed for check-in selection.
- Remove the v1 backend rule that limited player ratings to five entries so users can rate any player in the match roster.

## Capabilities

### New Capabilities
- `checkin-ui-and-integration`: Frontend check-in form behavior and match-detail integration for creating, editing, and reviewing the signed-in user’s own match record.

### Modified Capabilities
- `public-match-browsing`: Match detail responses now include the match-scoped player roster needed for check-in selection.
- `checkin-domain-and-api`: Check-in validation no longer limits player ratings to five entries for a match.

## Impact

- Frontend pages: match detail page, check-in entry flow, check-in form state, and signed-in match-detail presentation
- Frontend systems: authenticated API client usage, auth-aware page state, form validation, loading/error/success feedback
- Backend read contract: public match detail now needs to return the roster of players that belong to the match
- Backend APIs consumed by this change: `GET /matches/:id/my-checkin`, `POST /matches/:id/checkin`, `PUT /matches/:id/checkin`
- Backend validation: player-rating payload rules change from capped selection to roster-bounded selection
- Later dependencies: profile history and release-quality work depend on this UI flow being stable
