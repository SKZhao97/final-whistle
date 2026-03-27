Findings

  1. The backend is not in a buildable state yet. backend/go.mod declares external modules, but backend/go.sum is missing, so go build ./... cannot resolve locked dependencies. On top of that, there are unused imports in backend/internal/middleware/
     middleware.go:4 and backend/internal/utils/logger.go:4, which would fail compilation even after dependency resolution. This blocks the most basic validation for spec1.
  2. The API server is binding to the wrong address string in backend/cmd/api/main.go:54. ":" + string(cfg.Server.Port) converts the integer port to a single rune, not "8080"-style decimal text. That means the server will not reliably listen on the
     configured port.
  3. The frontend toolchain is internally inconsistent and currently not green. frontend/package.json:14 uses next 16.2.1, while the project docs/spec target Next 15, and the lint script is still next lint in frontend/package.json:10, which now fails
     outright. npm run type-check also fails because frontend/tailwind.config.ts:4 is typed/configured in a way that does not match the installed Tailwind/Next setup. Right now the frontend foundation does not meet the “basic tooling works” bar.
  4. Production build is not reproducible in a restricted or offline environment because the root layout depends on Google font fetching at build time via frontend/src/app/layout.tsx:2. In this environment, npm run build failed trying to fetch Geist and
     Geist Mono. For a foundation spec, this is a real portability risk.
  5. spec1 boundary is leaking into later specs on the frontend. The API client hardcodes cookie credentials in frontend/src/lib/api/client.ts:69, and frontend/src/types/api.ts:35 already defines auth/check-in/profile request and response shapes that
     belong to later feature specs. That contradicts the narrowed spec1 boundary we agreed on.
  6. The database layer still exposes AutoMigrate in backend/internal/db/database.go:120, even though spec1 was explicitly tightened to “manual migrations as the single schema source of truth.” Leaving this helper around invites schema drift later.
  7. The reset path has a broken SQL fallback in backend/seed/reset.go:34. TRUNCATE TABLE ? CASCADE does not substitute table identifiers in PostgreSQL, so if the delete path fails, the fallback will not work as intended. This is a concrete bug in a dev
     utility.
  8. The panic recovery middleware leaks internal panic content back to clients in backend/internal/middleware/middleware.go:53 and backend/internal/middleware/middleware.go:65. Even in development, returning raw recovered panic strings in the API
     payload is a bad default for shared foundation code.
  9. The docs and implementation are already drifting. README.md:10 says “Next.js 15”, README.md:122 says shadcn/ui and React Hook Form are part of the stack, but frontend/package.json:13 does not actually include them. That makes the setup instructions
     misleading for the next implementation step.
