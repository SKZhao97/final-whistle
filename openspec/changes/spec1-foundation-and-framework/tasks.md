## 1. Backend Framework Setup

- [x] 1.1 Initialize Go module with `go mod init final-whistle/backend`
- [x] 1.2 Create project directory structure: `cmd/`, `internal/`, `migrations/`, `seed/`
- [x] 1.3 Add dependencies: Gin, GORM, PostgreSQL driver, configuration libraries
- [x] 1.4 Implement configuration management with environment variable support
- [x] 1.5 Create main entry point in `cmd/api/main.go` with Gin server setup
- [x] 1.6 Implement health check endpoint at `/health`
- [x] 1.7 Add middleware: request logging, error recovery, CORS
- [x] 1.8 Create database connection pool configuration
- [x] 1.9 Set up structured logging with different log levels
- [x] 1.10 Create error handling framework with standard response format
- [x] 1.11 Set up base router skeleton with infrastructure routes only (`/health`) and route groups for future modules without placeholder business handlers

## 2. Frontend Framework Setup

- [x] 2.1 Initialize Next.js 15 project with TypeScript and App Router
- [x] 2.2 Set up Tailwind CSS with basic configuration
- [x] 2.3 Install and configure shadcn/ui component library
- [x] 2.4 Create project directory structure: `src/app/`, `src/components/`, `src/lib/`, `src/types/`
- [x] 2.5 Set up TypeScript with strict mode and path aliases
- [x] 2.6 Create basic layout components (header, footer, navigation)
- [x] 2.7 Implement base API client configuration for backend communication
- [x] 2.8 Create shared frontend type conventions and minimal response envelope types
- [x] 2.9 Set up app shell and only the minimal bootstrap routes needed for foundation validation
- [x] 2.10 Configure development scripts and build process
- [x] 2.11 Set up linting and formatting (ESLint, Prettier)

## 3. Database Schema Design

- [x] 3.1 Create migration for `users` table with all required fields
- [x] 3.2 Create migration for `teams` table with v1-required fields
- [x] 3.3 Create migration for `players` table with team relationship
- [x] 3.4 Create migration for `matches` table with v1-required fields and team relationships
- [x] 3.5 Create migration for `match_players` table for match participation
- [x] 3.6 Create migration for `check_ins` table with all rating fields
- [x] 3.7 Create migration for `player_ratings` table with check-in relationship
- [x] 3.8 Create migration for `tags` table with predefined dictionary
- [x] 3.9 Create migration for `checkin_tags` many-to-many relationship
- [x] 3.10 Create migration for `sessions` table for authentication
- [x] 3.11 Add all required indexes for performance
- [x] 3.12 Add unique constraint for `check_ins(user_id, match_id)`

## 4. GORM Model Definitions

- [x] 4.1 Create `User` model with GORM tags and validation
- [x] 4.2 Create `Team` model with v1-required fields
- [x] 4.3 Create `Player` model with team relationship
- [x] 4.4 Create `Match` model with home/away team relationships
- [x] 4.5 Create `MatchPlayer` model for participation tracking
- [x] 4.6 Create `CheckIn` model as aggregate root with all fields
- [x] 4.7 Create `PlayerRating` model with check-in relationship
- [x] 4.8 Create `Tag` model with predefined values
- [x] 4.9 Create `CheckInTag` model for many-to-many relationship
- [x] 4.10 Create `Session` model for authentication
- [x] 4.11 Define enum types: `MatchStatus`, `WatchedType`, `SupporterSide`
- [x] 4.12 Implement model validation methods for business rules
- [x] 4.13 Create repository package structure and minimal shared DB access abstractions needed by foundation

## 5. Seed Data Implementation

- [x] 5.1 Create base seed script for teams with realistic football clubs
- [x] 5.2 Create seed script for players with positions and team assignments
- [x] 5.3 Create seed script for matches with realistic scores and timestamps
- [x] 5.4 Create seed script for match-player relationships
- [x] 5.5 Create seed script for predefined tag dictionary
- [x] 5.6 Create optional dev user seed for local development bootstrap
- [x] 5.7 Defer sample check-ins to the dedicated check-in spec
- [x] 5.8 Implement modular seed system with individual script execution
- [x] 5.9 Add database reset functionality for development
- [x] 5.10 Validate seed data integrity and relationships

## 6. Integration & Validation

- [x] 6.1 Verify backend can connect to database and run migrations
- [x] 6.2 Verify frontend development server starts correctly
- [x] 6.3 Test health check endpoint from frontend
- [x] 6.4 Verify all seed scripts run successfully
- [x] 6.5 Test basic model queries through GORM
- [x] 6.6 Verify frontend can reach backend health endpoint
- [x] 6.7 Verify the app shell renders without runtime errors
- [x] 6.8 Create basic README with setup instructions
- [x] 6.9 Document configuration options for different environments
- [x] 6.10 Verify the complete foundation is ready for feature development
