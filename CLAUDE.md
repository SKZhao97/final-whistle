# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Final Whistle is a post-match recording product for football viewers. It helps users record matches they've watched, express ratings and feelings about matches, teams, and players, and build personal football viewing archives.

**Core Value**: Helps users record "what this match meant to you" rather than "telling you what happened in the match."

**v1 Scope**: Complete one closed loop: `login → browse matches → match detail → create/edit CheckIn → view aggregation → view personal profile`

## Tech Stack

- **Frontend**: Next.js 15, TypeScript, App Router, Tailwind CSS, shadcn/ui, React Hook Form, Zod
- **Backend**: Go, Gin, GORM, go-playground/validator
- **Database**: PostgreSQL
- **Testing**: Frontend: Vitest, Playwright; Backend: go test, httptest

## Architecture

The system follows a clean separation between frontend and backend:

```
User Browser → Next.js Frontend → Go Gin API → PostgreSQL
```

- **Frontend responsibilities**: Page routing, UI rendering, form interactions, API calls, login state display, loading/empty/error/unauthorized state handling
- **Backend responsibilities**: Cookie Session authentication, parameter validation, business rule validation, transaction control, data persistence, aggregation queries, unified error responses

## Directory Structure

### Frontend (`frontend/`)
```
frontend/
  src/
    app/
      page.tsx
      login/
        page.tsx
      matches/
        page.tsx
        [matchId]/
          page.tsx
      teams/
        [teamId]/
          page.tsx
      players/
        [playerId]/
          page.tsx
      me/
        page.tsx
    components/
      ui/
      layout/
      auth/
      matches/
      checkins/
      profile/
      teams/
      players/
    lib/
      api/
        client.ts
        auth.ts
        matches.ts
        checkins.ts
        teams.ts
        players.ts
        users.ts
      validations/
      utils/
    types/
      api.ts
      domain.ts
```

### Backend (`backend/`)
```
backend/
  cmd/
    api/
      main.go
  internal/
    config/
    db/
    middleware/
    router/
    handler/
      auth_handler.go
      match_handler.go
      checkin_handler.go
      team_handler.go
      player_handler.go
      user_handler.go
    service/
      auth_service.go
      match_service.go
      checkin_service.go
      team_service.go
      player_service.go
      user_service.go
    repository/
      user_repository.go
      match_repository.go
      checkin_repository.go
      team_repository.go
      player_repository.go
      tag_repository.go
      session_repository.go
    dto/
      auth_dto.go
      match_dto.go
      checkin_dto.go
      team_dto.go
      player_dto.go
      user_dto.go
      common_dto.go
    model/
      user.go
      team.go
      player.go
      match.go
      match_player.go
      checkin.go
      player_rating.go
      tag.go
      checkin_tag.go
      session.go
    utils/
  migrations/
  seed/
  go.mod
```

## Key API Endpoints

### Auth
- `POST /auth/login` - Dev login (email + name)
- `POST /auth/logout`
- `GET /auth/me` - Get current user

### Match
- `GET /matches` - List with filtering (competition, season)
- `GET /matches/:id` - Detail with aggregation data

### CheckIn
- `GET /matches/:id/my-checkin` - Get user's check-in for match
- `POST /matches/:id/checkin` - Create check-in
- `PUT /matches/:id/checkin` - Update check-in

### Team / Player
- `GET /teams/:id` - Team detail with recent matches
- `GET /players/:id` - Player detail with rating summary

### User Profile
- `GET /me/profile` - User profile summary
- `GET /me/checkins` - User's check-in history

## Data Model

**Core entities**:
- `User`, `Team`, `Player`, `Match`, `MatchPlayer`
- `CheckIn` (aggregate root), `PlayerRating`, `Tag`, `CheckInTag`, `Session`

**Key constraints**:
- Unique index: `check_ins(user_id, match_id)` - One check-in per user per match
- `match.status` must be `FINISHED` to create/update check-in
- Player ratings: max 5 players per check-in, must belong to the match
- All ratings: 1-10 integers

**Enums**:
- `match.status`: `SCHEDULED`, `FINISHED`
- `check_ins.watched_type`: `FULL`, `PARTIAL`, `HIGHLIGHTS`
- `check_ins.supporter_side`: `HOME`, `AWAY`, `NEUTRAL`

## Development Setup

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+

### Expected Commands

**Backend**:
```bash
# Run migrations
go run cmd/migrate/main.go

# Run seed data
go run cmd/seed/main.go

# Start development server
go run cmd/api/main.go

# Run tests
go test ./...

# Run specific test
go test -v ./internal/service -run TestAuthService_Login
```

**Frontend**:
```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm test

# Run E2E tests
npm run test:e2e

# Build for production
npm run build
```

## Development Phases (v1)

Follow this order to ensure proper dependencies:

1. **Phase 1: Foundation** - Database schema, migrations, seed data, basic models
2. **Phase 2: Public APIs** - `GET /matches`, `GET /matches/:id`, `GET /teams/:id`, `GET /players/:id`
3. **Phase 3: Auth & CheckIn Write** - Login, logout, create/update check-ins
4. **Phase 4: User Profile & Aggregation** - Personal profile, check-in history, aggregation queries
5. **Phase 5: Frontend Integration** - Page connections, loading states, error handling, E2E tests

## Business Rules

### CheckIn Creation
- User can only have one check-in per match
- Match must be `FINISHED`
- Player ratings: 0-5 players, each 1-10 rating
- Tags: from predefined dictionary only
- `shortReview`: optional, max 280 characters

### Aggregation Display
- Show "sample too small" when check-in count < 3
- All averages: 1 decimal place
- Real-time calculation (no caching in v1)

### Authentication
- HTTP-only Cookie Session
- Dev login only in v1 (no OAuth)
- Session stored in database

## API Conventions

### Success Response
```json
{
  "success": true,
  "data": {}
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "invalid request body",
    "details": {}
  }
}
```

### Error Codes
- `UNAUTHORIZED` - Authentication required/invalid
- `FORBIDDEN` - Insufficient permissions
- `NOT_FOUND` - Resource doesn't exist
- `VALIDATION_ERROR` - Invalid request parameters
- `CONFLICT` - Resource conflict (e.g., duplicate check-in)
- `INTERNAL_ERROR` - Server-side error

## Testing Strategy

### Backend Focus
- `AuthService.Login` - Dev login flow
- `CheckInService.CreateCheckIn` - Business rule validation
- `CheckInService.UpdateCheckIn` - Update with constraints
- `MatchService.GetMatchDetail` - Aggregation queries
- `UserService.GetProfileSummary` - User statistics

### Frontend E2E
Cover the main path:
1. Login
2. Browse match list
3. Open match detail
4. Create CheckIn
5. Edit CheckIn
6. Open personal profile

## Important Notes

- **v1 uses internal seed data only** - No external API integration
- **Aggregation is real-time** - No caching or pre-computation
- **Focus on the closed loop** - Avoid scope creep
- **CheckIn is the aggregate root** - Structure all related data around it
- **Follow the phase order** - Each phase depends on the previous one

## Project Context

This project follows a **spec-driven development approach** with OpenSpec. Always refer to the design documents in `docs/` for detailed specifications before implementation.