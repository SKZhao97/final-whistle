## Why

Final Whistle v1 requires a solid foundation before feature development can begin. This change establishes the minimum shared technical infrastructure: backend and frontend frameworks, database schema, data models, and base seed data. Without this foundation, subsequent specs (public match browsing, session auth, check-in workflow, profile, etc.) cannot be implemented with stable boundaries.

## What Changes

- **Backend Framework Setup**: Initialize Go + Gin project structure with proper configuration, middleware, minimal routing, and database connection
- **Frontend Framework Setup**: Initialize Next.js 15 project with TypeScript, Tailwind CSS, shadcn/ui, and app shell
- **Database Schema Design**: Create PostgreSQL migrations for all core tables with proper indexes and constraints
- **Data Models Definition**: Implement GORM models for all entities with proper relationships and business enums
- **Seed Data Scripts**: Create base seed data for teams, players, matches, match-player relations, tags, and optional dev users

## Capabilities

### New Capabilities
- **backend-framework**: Go + Gin backend framework setup with configuration, middleware, minimal routing, and database connection
- **frontend-framework**: Next.js 15 frontend framework setup with TypeScript, Tailwind CSS, shadcn/ui, and app shell
- **database-schema**: PostgreSQL database schema design with migrations for all core tables and indexes
- **data-models**: GORM model definitions for all entities with relationships and business enums
- **seed-data**: Base seed data scripts for development bootstrap

### Modified Capabilities
<!-- No existing capabilities to modify -->

## Impact

- **Code Impact**: Creates shared project structure for both frontend and backend
- **API Impact**: Establishes shared infrastructure for future API endpoints without implementing business features
- **Database Impact**: Creates all core tables with proper relationships and constraints
- **Development Workflow**: Sets up standardized development environment and tooling
- **Testing Impact**: Provides base data needed by subsequent feature specs
