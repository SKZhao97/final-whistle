## Why

The backend check-in contract is now in place, but the product loop still breaks on the main match-detail path because signed-in users cannot create or edit their own record from the UI. This change makes the core “record this match” flow usable end to end.

## What Changes

- Add signed-in check-in UI behavior to the match detail page, including loading the current user’s existing check-in state.
- Add a dedicated check-in form flow for creating and editing a match record using the existing authenticated backend APIs.
- Add frontend validation, submit states, and success/error handling for check-in create and update actions.
- Add match-detail UI integration so the page reflects the user’s current record after successful submission without requiring a manual refresh.

## Capabilities

### New Capabilities
- `checkin-ui-and-integration`: Frontend check-in form behavior and match-detail integration for creating, editing, and reviewing the signed-in user’s own match record.

### Modified Capabilities

## Impact

- Frontend pages: match detail page, check-in entry flow, check-in form state, and signed-in match-detail presentation
- Frontend systems: authenticated API client usage, auth-aware page state, form validation, loading/error/success feedback
- Backend APIs consumed by this change: `GET /matches/:id/my-checkin`, `POST /matches/:id/checkin`, `PUT /matches/:id/checkin`
- Later dependencies: profile history and release-quality work depend on this UI flow being stable
