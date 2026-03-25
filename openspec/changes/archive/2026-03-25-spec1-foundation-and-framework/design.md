## Context

Final Whistle is a post-match recording product for football viewers. The v1 goal is to complete a closed loop: `login → browse matches → record match → view aggregation → view personal profile`. This change establishes the foundational technical infrastructure required before any feature development can begin. The project uses a modern tech stack: Go + Gin backend, Next.js 15 frontend, and PostgreSQL database.

## Goals / Non-Goals

**Goals:**
1. Establish a production-ready backend framework with proper configuration, middleware, and database connectivity
2. Set up a modern frontend framework with TypeScript, Tailwind CSS, and API client infrastructure
3. Design and implement a complete database schema for all core entities with proper relationships
4. Create GORM models that accurately represent the domain entities and business rules
5. Provide base seed data for development across foundational entities

**Non-Goals:**
1. Implement any business logic or feature functionality (authentication, match browsing, check-in workflows, profile queries, etc.)
2. Create production deployment configurations or CI/CD pipelines
3. Implement advanced performance optimizations or caching strategies
4. Create comprehensive test suites (beyond basic connectivity tests)
5. Implement security hardening beyond basic configuration

## Decisions

### 1. Backend Architecture Pattern
**Decision**: Use layered architecture (Handler → Service → Repository) for clear separation of concerns
**Rationale**: This pattern provides clear boundaries between API layer, business logic, and data access, making the codebase maintainable and testable
**Alternatives Considered**:
- Monolithic handlers with direct database access (rejected: poor separation of concerns)
- Clean Architecture/Hexagonal (rejected: overkill for v1 scope)

### 2. Database Migration Tool
**Decision**: Use version-controlled manual migrations as the single schema source of truth
**Rationale**: Explicit migrations keep schema evolution reviewable, deterministic, and aligned with later specs
**Alternatives Considered**:
- Use GORM AutoMigrate (rejected: schema drift is harder to review and control)

### 3. Configuration Management
**Decision**: Use environment variables with sensible local defaults
**Rationale**: This keeps startup simple and predictable while remaining deployment-friendly
**Alternatives Considered**:
- Config files with fallback (rejected: adds another configuration source to maintain in v1)

### 4. Frontend State Management
**Decision**: Use React state + useEffect for v1, no external state management library
**Rationale**: v1 has minimal complex state needs; keeping it simple reduces complexity
**Alternatives Considered**:
- Zustand (rejected: premature optimization)
- Redux (rejected: overkill for v1)

### 5. Seed Data Strategy
**Decision**: Create modular base seed scripts that can be run independently or together
**Rationale**: This supports development bootstrap now without forcing check-in or profile-specific data too early
**Alternatives Considered**:
- Single monolithic seed file (rejected: inflexible)
- Database dumps (rejected: not version-controlled)

## Risks / Trade-offs

**Risk**: Over-engineering the foundation vs. moving fast
- **Mitigation**: Keep spec1 limited to shared infrastructure and base data; defer feature-shaped placeholders to later specs

**Risk**: Database schema changes causing breaking changes in later specs
- **Mitigation**: Thoroughly review schema against all planned v1 features; use migrations for safe evolution

**Risk**: Frontend/backend integration issues due to framework mismatches
- **Mitigation**: Establish clear API contracts early; create shared TypeScript/Go types where possible

**Risk**: Seed data becoming outdated as features evolve
- **Mitigation**: Limit spec1 seed data to foundational entities; add feature-specific samples in later specs

**Risk**: Configuration complexity slowing down development
- **Mitigation**: Provide sensible defaults; document configuration clearly in README
