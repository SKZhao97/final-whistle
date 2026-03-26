# Final Whistle

Final Whistle is a post-match recording product for football viewers. Users can record matches they've watched, rate matches, teams, and players, and build their personal football memory archive.

## Project Structure

```
final-whistle/
├── backend/          # Go + Gin backend API
├── frontend/         # Next.js frontend
├── docs/             # Documentation
└── openspec/         # OpenSpec specifications
```

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+

### 1. Database Setup

```bash
# Create a PostgreSQL database
createdb final_whistle

# Or using psql
psql -c "CREATE DATABASE final_whistle;"
```

### 2. Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Set environment variables
export DATABASE_URL="postgres://localhost:5432/final_whistle?sslmode=disable"
export PORT=8080
export ENV=development

# Apply SQL migrations
go run ./cmd/migrate/main.go

# Seed development data
go run ./cmd/seed/main.go -scope=all

# Start the server
go run ./cmd/api/main.go
```

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## Configuration

### Backend Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection URL | `postgres://postgres:postgres@localhost:5432/final_whistle?sslmode=disable` |
| `PORT` | Server port | `8080` |
| `ENV` | Environment (`development`, `production`) | `development` |
| `LOG_LEVEL` | Log level (`debug`, `info`, `warn`, `error`) | `info` |

### Frontend Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:8080` |

## Seed Data

The database includes seed data for development:

```bash
cd backend
go run ./cmd/seed/main.go -scope=all
```

Seed data includes:
- 6 football teams (Premier League clubs)
- 9 players with positions
- 5 matches with realistic scores
- 10 predefined emotion/impression tags
- 2 development users

## Development

### Backend Architecture

The backend follows a layered architecture:

- **Handlers**: HTTP request handling, validation, response formatting
- **Services**: Business logic and orchestration
- **Repositories**: Data access and persistence
- **Models**: Domain entities and database mappings

### Frontend Architecture

The frontend uses:

- **Next.js 16** with App Router
- **React 19**
- **TypeScript** with strict type checking
- **Tailwind CSS 4** for styling
- Shared app shell and type-safe API client foundation

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `GET /` - API info
- `GET /matches` - Public match list
- `GET /matches/:id` - Public match detail
- `GET /teams/:id` - Public team detail
- `GET /players/:id` - Public player detail

### Protected Endpoints (require authentication)

- `POST /auth/login`
- `POST /auth/logout`
- `GET /auth/me`
- `GET /matches/:id/my-checkin`
- `POST /matches/:id/checkin`
- `PUT /matches/:id/checkin`

## Database Schema

Key tables:

- `users` - Registered users
- `teams` - Football teams
- `players` - Football players
- `matches` - Football matches
- `check_ins` - User match records with ratings
- `player_ratings` - Individual player ratings
- `tags` - Emotion/impression tags
- `checkin_tags` - Many-to-many relationship between check-ins and tags
- `sessions` - Authentication sessions

## Testing

```bash
# Backend build
cd backend
go build ./...

# Frontend validation
cd frontend
npm run lint
npm run type-check
npm run build
```

## License

MIT
