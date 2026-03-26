## Context

Final Whistle now has public browsing and session auth in place, but the product loop still cannot record the user’s actual match reaction. The next backend slice is the check-in write domain: a signed-in user must be able to load their own existing check-in for a match, create one if none exists, and update it later without breaking the one-user-one-match constraint.

The schema foundation is already present. `check_ins`, `player_ratings`, `checkin_tags`, `match_players`, and `tags` tables exist, and `session-auth` already provides the authenticated user context. This change does not need new product-facing read aggregation; it needs a strict backend write contract that later UI work can consume safely.

## Goals / Non-Goals

**Goals:**
- Establish `GET /matches/:id/my-checkin`, `POST /matches/:id/checkin`, and `PUT /matches/:id/checkin` as the backend contract for the current user’s match record.
- Enforce the core domain rules on the server: signed-in access, finished-match-only write access, one check-in per user per match, valid tags, valid player references, and score/text limits.
- Persist `check_ins`, `player_ratings`, and `checkin_tags` transactionally so create and update flows remain consistent.
- Return a stable full check-in DTO that later frontend form integration can reuse without redesigning the API.

**Non-Goals:**
- Building the full check-in form UI, match-detail integration, or success/error UX polish on the frontend.
- Updating public match aggregate endpoints beyond what naturally results from the underlying stored data.
- Profile history endpoints such as `GET /me/checkins` or user summary calculations.
- Background aggregation, caching, moderation, or audit-history features.

## Decisions

### 1. Treat `CheckIn` as the write aggregate root

The backend will model `CheckIn` as the single authoritative record per `(user_id, match_id)`. `player_ratings` and `checkin_tags` will always be loaded and written through that parent record, never as standalone write endpoints.

Why this approach:
- It matches the existing domain design and unique constraint model.
- It keeps validation and transactional writes centered on one aggregate boundary.
- It avoids exposing partial write paths that could leave tags or player ratings detached from the parent record.

Alternatives considered:
- Separate write endpoints for ratings or tags: rejected because they complicate consistency and do not match v1 workflow needs.

### 2. Use create-vs-update semantics instead of upsert endpoints

`POST /matches/:id/checkin` will require that the current user does not already have a check-in for that match. `PUT /matches/:id/checkin` will require that the check-in already exists.

Why this approach:
- It keeps API behavior explicit and easy to reason about.
- It prevents silent overwrites when a client accidentally uses the wrong method.
- It matches the current product semantics of “create once, then edit.”

Alternatives considered:
- One upsert endpoint: rejected because it hides important domain distinctions and weakens error handling.

### 3. Enforce all check-in business rules in the service layer before persistence

The service will validate:
- authenticated user presence
- match existence
- match status is `FINISHED` for create and update
- uniqueness for create, existence for update
- score range `1-10`
- `shortReview` length `<= 280`
- player rating count `<= 5`
- player note length `<= 80`
- every rated player belongs to the target match through `match_players`
- every submitted tag exists and is active

Why this approach:
- Database constraints cover only part of the domain.
- These rules need clean API errors, not only raw persistence failures.
- It keeps future UI integration predictable because invalid input is rejected consistently in one place.

Alternatives considered:
- Relying mostly on database constraints: rejected because the error surface is too coarse and misses relationship validations like player eligibility.

### 4. Use full replacement semantics for update child collections

On `PUT`, the backend will update the parent `check_ins` row and replace `player_ratings` and `checkin_tags` wholesale inside one transaction.

Why this approach:
- The client submits the full current form state, so replacement semantics are simpler and less error-prone than diff-based patching.
- It avoids stale child rows surviving after an edit.
- It keeps update logic deterministic for v1.

Alternatives considered:
- Child-level diff updates: rejected because they add complexity without helping the v1 UX.

### 5. Return `data: null` for missing `my-checkin`

`GET /matches/:id/my-checkin` will be a protected endpoint that returns success with `data: null` when the current user has no record for that match, rather than `404`.

Why this approach:
- The absence of the current user’s check-in is not an error condition for the match resource.
- It gives later UI integration a clean way to branch between “create” and “edit” states.

Alternatives considered:
- Returning `404`: rejected because it conflates “match missing” with “user record absent.”

## Risks / Trade-offs

- [Validation drift between API and later UI] → Keep the server as the source of truth and return structured validation errors for invalid payloads.
- [Transactional update logic may accidentally leave child rows stale] → Replace child collections inside one transaction and test update behavior explicitly.
- [Public aggregate endpoints may read partially updated data if writes are not atomic] → Wrap create and update flows in a single database transaction.
- [Player eligibility checks may become expensive] → Validate against `match_players` in one batched query per request rather than per-player lookups.

## Migration Plan

1. Add DTOs, repositories, services, and handlers for the authenticated check-in read/write contract.
2. Reuse `session-auth` middleware to require a current user on all three endpoints.
3. Implement transactional create and update flows on top of existing tables.
4. Validate against seeded matches, tags, and match-player relationships before UI integration begins.

Rollback:
- Remove the route registrations and stop exposing the endpoints.
- The underlying schema remains unchanged because this change activates existing tables rather than introducing new ones.

## Open Questions

- Should duplicate rated-player entries in a single payload be rejected explicitly, or should the service collapse them before validation?
- Should `watchedAt` be allowed to be any timestamp, or should v1 constrain it to be at or after the match kickoff time?
