## 1. Frontend API and state plumbing

- [x] 1.1 Extend frontend check-in and match-detail types to cover roster-aware match detail, current-user check-in detail, create payload, and update payload
- [x] 1.2 Add authenticated API client helpers for `GET /matches/:id/my-checkin`, `POST /matches/:id/checkin`, and `PUT /matches/:id/checkin`
- [x] 1.3 Add match-detail page state that conditionally fetches `my-checkin` only when the user is signed in
- [x] 1.4 Update the backend/public match detail contract to return the match roster needed for player-rating selection
- [x] 1.5 Remove the backend five-player-rating cap while preserving duplicate-player and roster-membership validation

## 2. Match detail integration

- [x] 2.1 Add signed-out, signed-in-without-check-in, and signed-in-with-check-in UI states to the match detail page
- [x] 2.2 Add a saved-record summary view on the match detail page for users who already have a check-in
- [x] 2.3 Add non-finished match handling so scheduled matches do not expose an active check-in submission path

## 3. Check-in form flow

- [x] 3.1 Implement a reusable check-in form component that supports both create and edit modes
- [x] 3.2 Initialize the form with empty defaults for create mode and backend data for edit mode
- [x] 3.3 Add client-side validation for ratings, text limits, duplicate player selection, and roster-bounded player selection
- [x] 3.4 Wire create-mode submit to `POST /matches/:id/checkin` with loading and error handling
- [x] 3.5 Wire edit-mode submit to `PUT /matches/:id/checkin` with loading and error handling

## 4. Refresh and validation

- [x] 4.1 Refresh the match-detail UI after successful create or update so the saved record state is immediately visible
- [x] 4.2 Add frontend tests or component-level coverage for entry-state branching and create/edit submit behavior where practical
- [x] 4.3 Validate with `npm run lint`, `npm run type-check`, `npm run build`, and a manual browser flow against the local backend
