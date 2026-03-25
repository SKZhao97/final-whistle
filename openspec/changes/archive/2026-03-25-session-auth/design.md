## Context

Final Whistle v1 has finished the public browsing slice, but every remaining user-facing feature depends on a stable authenticated identity. The project has already committed to `HTTP-only Cookie Session + dev login`, and the `users` plus `sessions` tables already exist from the foundation work. This change turns those dormant tables into an active auth layer without introducing OAuth, JWT, or any third-party identity dependency.

The current backend is a Gin JSON API with shared response envelopes and middleware support. The current frontend is a Next.js App Router application with a centralized API client and public pages. The auth design needs to fit those existing patterns and stay narrow enough that check-in and profile specs can build on it without reworking the contract.

## Goals / Non-Goals

**Goals:**
- Establish the v1 login contract with `POST /auth/login`, `POST /auth/logout`, and `GET /auth/me`.
- Persist sessions in PostgreSQL and issue an `HTTP-only` cookie that the browser automatically sends on later requests.
- Add backend middleware that resolves the current user from the session cookie and can protect later routes.
- Add frontend auth helpers and app-level auth recovery so the UI can determine whether a user is signed in.
- Keep the login flow development-friendly by using seeded internal users and optional auto-create behavior for local development.

**Non-Goals:**
- GitHub OAuth, password-based auth, magic links, JWT, refresh tokens, or multi-provider auth.
- Authorization rules for check-in and profile resources beyond “must be signed in.”
- Full login page UX polish, redirect orchestration, or final protected-route UX for every page.
- Session analytics, device management, or background cleanup jobs beyond basic expiry handling.

## Decisions

### 1. Use database-backed opaque sessions with an `HTTP-only` cookie

The backend will create a random session token, store it in the `sessions` table with `user_id` and `expired_at`, and send that token in an `HTTP-only` cookie. Each authenticated request will look up the session by token, reject expired sessions, and load the current user.

Why this approach:
- It matches the existing v1 architecture decision and avoids introducing a second auth model.
- The database already contains the required `sessions` table, so this reuses existing infrastructure.
- Opaque session tokens keep user identity and session metadata server-side, which is simpler than JWT for v1.

Alternatives considered:
- JWT in cookies: rejected because it adds token validation and revocation complexity without solving a v1 problem.
- In-memory session store: rejected because it would break across restarts and diverge from the existing schema.

### 2. Keep login as `dev login` against the local user table

`POST /auth/login` will accept the minimal development login payload and resolve a user from the `users` table. If the submitted email is not present, the backend MAY create a local dev user using the submitted name and email, but that behavior is limited to development mode.

Why this approach:
- It matches the project’s documented v1 scope.
- It keeps auth unblocked by external providers.
- It works naturally with existing seeded users like `demo@final-whistle.test`.

Alternatives considered:
- Hardcoded single demo user only: rejected because it is too rigid for local testing.
- GitHub OAuth now: rejected because it increases scope and operational cost before the product loop is proven.

### 3. Represent auth state through `/auth/me`, not by decoding cookies in the frontend

The frontend will restore login state by calling `GET /auth/me` and treating the backend as the source of truth. The browser will send the session cookie automatically, and the frontend will only consume the returned user summary or an unauthorized response.

Why this approach:
- It keeps session handling entirely server-owned.
- It fits the current API client pattern and avoids exposing token details to the UI.
- It gives later specs a single consistent auth bootstrap path.

Alternatives considered:
- Reading user data from a client-stored token: rejected because it breaks the `HTTP-only` session model.
- Duplicating auth state in multiple places: rejected because it increases inconsistency risk.

### 4. Add middleware for optional and required auth separately

The backend auth layer will expose:
- a resolver that attaches the current user to request context when a valid session exists
- a guard for protected endpoints that returns `401 UNAUTHORIZED` when no valid session exists

Why this approach:
- `/auth/me` and future protected APIs need shared session parsing.
- Some routes may only need “current user if present” behavior, while later write routes need strict enforcement.

Alternatives considered:
- Only a strict middleware: rejected because it is less flexible for mixed public/auth-aware routes.

### 5. Use conservative cookie defaults for local development

The session cookie will be `HttpOnly`, `Path=/`, and use `SameSite=Lax`. `Secure` can remain environment-sensitive so local HTTP development still works. Logout will clear both the database session and the cookie.

Why this approach:
- It is safe enough for v1 local development while remaining deployable later.
- It avoids accidental breakage in non-HTTPS local environments.

Alternatives considered:
- Always `Secure=true`: rejected because it would complicate local development on plain HTTP.

### 6. Keep frontend login UX minimal and implementation-agnostic in this spec

This spec requires a usable development login entry point and auth-state recovery, but it does not require a final standalone login page or a specific client-side state library. The implementation may use a dedicated route, a simple shell-level entry, or another minimal UI that supports login, logout, and current-user restoration.

Why this approach:
- It keeps the auth spec focused on capability rather than UI polish.
- It avoids locking the project into React Context, Zustand, or any other state solution before that trade-off matters.
- It leaves room for a later UI-focused spec to refine the login experience without changing the auth contract.

Alternatives considered:
- Requiring a full `/login` page now: rejected because it couples auth infrastructure to unnecessary UI scope.
- Mandating a specific frontend state library: rejected because the current requirement is behavior, not implementation choice.

## Risks / Trade-offs

- [Session expiry handling can become inconsistent] → Store expiry in the database and enforce it centrally in middleware and `/auth/me`.
- [Auto-creating dev users may create noisy local data] → Keep the behavior explicit in the spec and constrain it to the development login contract.
- [Cookie behavior may diverge across local and deployed environments] → Keep `Secure` environment-sensitive while making `HttpOnly`, `Path=/`, and `SameSite=Lax` part of the contract.
- [Frontend auth bootstrap can create extra request overhead] → Limit bootstrap to `/auth/me` and keep the returned user summary small.
- [Protected route behavior may evolve in later specs] → Keep this spec focused on auth state resolution and basic unauthorized handling, not final UX flows.

## Migration Plan

1. Implement auth DTOs, repositories, services, handlers, and middleware on top of the existing `users` and `sessions` tables.
2. Add frontend auth helpers and login-state bootstrap using the existing API client.
3. Validate login, logout, and `/auth/me` behavior locally with seeded users, session expiry checks, and environment-appropriate cookie settings.
4. Use the new auth middleware as the dependency for later protected check-in and profile APIs.

Rollback:
- Remove the auth routes from the router.
- Stop issuing session cookies.
- Leave the `sessions` table in place; it is harmless foundational schema and does not require rollback.

## Open Questions

- Should development-mode auto-created users be limited to a local test email pattern such as `*.test`, or is unrestricted local creation acceptable for v1?
