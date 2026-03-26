## Context

Final Whistle already supports public match browsing, session auth, and backend check-in read/write APIs. The remaining gap in the core v1 loop is frontend integration: a signed-in user can browse a match and authenticate, but still cannot create or edit the one record that defines the product’s core value.

The match detail page is already the center of the browsing experience, so this change should attach check-in UI there rather than inventing a separate flow. The backend write contract is stable enough to consume directly: `GET /matches/:id/my-checkin` determines create vs edit mode, while `POST` and `PUT` persist the full form payload. One additional read contract change is required: the public match detail response must expose the match-scoped roster from `match_players` so the frontend can build a valid player-rating selector.

## Goals / Non-Goals

**Goals:**
- Add a check-in entry point on the match detail page that adapts to signed-out, signed-in-without-check-in, and signed-in-with-check-in states.
- Add a reusable check-in form UI that supports both create and edit flows using the existing backend DTO shape.
- Keep form validation and submit behavior consistent with backend rules so the user sees actionable feedback before and after submit.
- Refresh the match detail experience after successful submit so the user’s own record is immediately visible and the page feels closed-loop.
- Expose the full match-scoped player roster in public match detail so the check-in form can rate any player attached to the match.

**Non-Goals:**
- Adding new backend domain rules beyond the roster-read addition and removal of the five-player cap.
- Building profile/history pages or global “my check-ins” views.
- Adding E2E automation, visual polish passes, or broad release-hardening work that belongs to the final spec.
- Reworking the public match browsing information architecture outside the match detail page and the minimum navigation needed for check-in.

## Decisions

### 1. Keep check-in interaction anchored on the match detail page

The match detail page will remain the primary container for both community aggregates and the signed-in user’s own record state. The UI will provide a visible call to action when the user is signed in and eligible to record the match.

Why this approach:
- It matches the PRD flow: browse match, open detail, record reaction, return to the same detail page.
- It avoids splitting the experience across unrelated pages.
- It keeps the create/edit context tied to the match aggregates the user is reacting to.

Alternatives considered:
- A standalone top-level check-in page: rejected because it weakens the match-centered UX and adds navigation overhead.

### 2. Use one form component for both create and edit modes

The frontend will implement a single form model that can initialize from an existing `my-checkin` response or from sensible empty defaults. Submit behavior will switch between `POST` and `PUT` based on whether a current record exists.

Why this approach:
- The backend already uses a stable full-detail DTO shape that supports both modes.
- One form model minimizes drift between create and edit behavior.
- It reduces duplicated validation and rendering logic.

Alternatives considered:
- Separate create and edit forms: rejected because the fields are identical and the complexity would be duplicated without product benefit.

### 3. Treat `GET /matches/:id/my-checkin` as a page-level side read, not a global auth concern

The match detail page will fetch public match detail first, then conditionally fetch `my-checkin` only when auth state indicates a signed-in user. Absence of a user record remains a normal state rather than an error.

Why this approach:
- `my-checkin` is match-specific, not global session state.
- It avoids coupling general auth recovery with every match-detail render path.
- It preserves the public browsing page for signed-out users without blocking on protected requests.

Alternatives considered:
- Global preload of all current-user match records: rejected because it is unnecessary for v1 and would add cross-page data complexity.

### 4. Use match-scoped roster data as the source for player selection

The check-in form will build player-rating selectors from the `match_players` roster returned with public match detail. For v1, any player attached to `match_players` is considered rateable; the form will not attempt to distinguish actual appearance minutes from roster membership.

Why this approach:
- It unblocks the frontend without introducing a new appearance-tracking data model.
- It matches the backend validation rule that player IDs must belong to the match through `match_players`.
- It keeps future upgrades possible if later data sources add true appearance semantics.

Alternatives considered:
- Adding a new appearance-specific schema now: rejected because the project does not yet have a trustworthy appearance data source.

### 5. Mirror backend rules in lightweight frontend validation, but keep server as the source of truth

The form will apply immediate frontend checks for required enumerations, rating ranges, text limits, and duplicate player selection. Backend validation responses will still be surfaced as authoritative errors on submit.

Why this approach:
- It makes the form feel responsive and prevents obvious bad submits.
- It avoids promising that the frontend alone defines correctness.
- It matches the existing backend-centered domain design.

Alternatives considered:
- No client-side validation: rejected because the form would feel brittle and users would discover basic mistakes too late.

### 6. Refresh the detail UI after submit using server data

After a successful create or update, the page should use the returned check-in payload immediately and also refresh the relevant detail data path so the match page reflects the saved state instead of optimistic local assumptions.

Why this approach:
- It reduces the chance of stale page state after submit.
- It keeps the UI aligned with server-returned canonical data.
- It lets later aggregate changes appear without requiring a full manual browser refresh.

Alternatives considered:
- Local-only optimistic patching: rejected because the match detail page already has multiple public and private data sections and optimistic syncing would add unnecessary fragility.

## Risks / Trade-offs

- [Public detail and private check-in state can drift during navigation] → Keep `my-checkin` fetch scoped to the active match detail page and refresh after submit.
- [Frontend validation may diverge from backend rules] → Keep frontend checks intentionally lightweight and surface backend validation messages on failure.
- [Roster does not equal actual on-pitch appearances] → Document that v1 uses `match_players` as the rateable-player source and defer true appearance semantics to the data-ingestion phase.
- [The match detail page may become too crowded on mobile] → Treat the check-in form as a contained panel or route-level flow with clear collapsed/read modes.
- [Signed-out and signed-in states may create branching complexity] → Define explicit UI states for signed-out, signed-in-without-check-in, and signed-in-with-check-in behavior in the spec.

## Migration Plan

1. Extend the public match detail read contract to include the match-scoped player roster and remove the five-player write cap from check-in validation.
2. Extend the frontend API client/types to support roster-aware match detail plus `my-checkin`, create, and update calls.
3. Add match-detail UI state for loading the current user’s record when authenticated.
4. Implement the shared check-in form and wire create/update submit behavior.
5. Update the match detail page to show signed-in entry points, saved-record summaries, and post-submit refresh behavior.
6. Validate with frontend static checks and a manual browser flow against the local backend.

Rollback:
- Remove the check-in form UI and restore the match detail page to public-browsing-only behavior.
- The backend API remains available because this change consumes existing endpoints rather than changing server contracts.

## Open Questions

- Should the form live inline on the match detail page, in a drawer/modal, or on a nested route such as `/matches/:id/checkin`?
- After submit, should the UI stay in editable form mode or collapse back to a “my record” summary card by default?
