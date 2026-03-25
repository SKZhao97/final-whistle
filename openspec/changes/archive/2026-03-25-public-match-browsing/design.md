## Context

Spec1 established the shared project foundation: backend framework, frontend app shell, schema, models, and seed data. The system already contains enough seeded teams, players, matches, tags, and match-player relationships to support a meaningful read-only browsing experience. The next change should expose this data through stable public APIs and user-facing pages without pulling in authentication or write flows too early.

## Goals / Non-Goals

**Goals:**
- Implement public read APIs for matches, match detail, teams, and players
- Surface match-level aggregation data needed by the public browsing flow
- Replace placeholder frontend routes with real public pages backed by the API
- Define clean empty/error/not-found behavior for read-only routes
- Keep the change implementation-ready for later auth and check-in specs

**Non-Goals:**
- Implement login, logout, session recovery, or protected routes
- Implement `my-checkin`, create check-in, or update check-in
- Add new schema changes, background jobs, or caching
- Add advanced filtering beyond the agreed v1 public list query shape
- Implement full profile pages or social/community interaction

## Decisions

### 1. Keep the change strictly read-only
**Decision**: Only public GET endpoints and read-only frontend pages are included.
**Rationale**: This keeps spec2 focused and prevents auth/check-in complexity from leaking in before the public browsing experience is stable.
**Alternatives considered**:
- Add `GET /matches/:id/my-checkin` now (rejected: introduces auth boundary too early)
- Add partial check-in form scaffolding now (rejected: spec boundary leak)

### 2. Compose match detail from base match data plus real-time aggregates
**Decision**: `GET /matches/:id` should assemble one response from match base data, aggregate score summaries, recent reviews, and player rating summaries.
**Rationale**: This aligns with the v1 technical design and gives the frontend a single detail payload for the public page.
**Alternatives considered**:
- Separate aggregate endpoints (rejected: more frontend orchestration with no clear v1 benefit)
- Precomputed aggregation tables (rejected: unnecessary complexity for current scale)

### 3. Use empty-state friendly DTOs for public detail views
**Decision**: Detail endpoints should always return stable structures, with `null` or empty arrays for missing aggregate data.
**Rationale**: Public pages need to render cleanly before any real check-ins exist, especially in early development and demo environments.
**Alternatives considered**:
- Omit aggregate fields entirely when empty (rejected: creates avoidable frontend branching)

### 4. Keep frontend data fetching page-oriented
**Decision**: Each public route fetches its own page data through the shared API client and renders route-specific loading/error/empty states.
**Rationale**: This matches the existing frontend foundation and avoids introducing broader client state patterns before they are needed.
**Alternatives considered**:
- Centralized cross-route state management (rejected: premature)
- Server Actions for read APIs (rejected: the project already treats backend HTTP APIs as the source of truth)

## Risks / Trade-offs

- **Risk**: Match detail aggregation queries become too coupled to future check-in behavior
  **Mitigation**: Keep read DTOs focused on public display needs and avoid encoding write-flow assumptions

- **Risk**: Empty aggregation data may hide integration mistakes during development
  **Mitigation**: Distinguish between `404`, real empty aggregate states, and internal errors in handler/service code

- **Risk**: Team and player pages become oversized in this spec
  **Mitigation**: Limit them to concise v1 summaries: base info, recent related matches, and rating summary data

- **Risk**: Public pages drift from the technical design if route-level state handling is inconsistent
  **Mitigation**: Define one consistent pattern for loading, empty, error, and not-found handling across all public pages
