## 1. Locale Foundation

- [x] 1.1 Add a shared locale provider and locale persistence mechanism for the frontend app shell
- [x] 1.2 Add a homepage-visible language switcher and reuse it in the global header where appropriate
- [x] 1.3 Define the initial English and Chinese translation dictionaries for shared UI copy
- [x] 1.4 Add shared frontend helpers for localized enum labels and common translated text access

## 2. Backend Localized Tag Data

- [x] 2.1 Add bilingual tag columns to the database schema through a version-controlled migration
- [x] 2.2 Update the `Tag` model and related DTO mapping to support localized display names
- [x] 2.3 Add backend locale resolution from persisted user preference for request handling
- [x] 2.4 Update tag seeding to populate and refresh both English and Chinese display names

## 3. Locale-Aware API Responses

- [x] 3.1 Update public match-detail tag output to return the locale-appropriate tag `name`
- [x] 3.2 Update current-user check-in read/write responses to return the locale-appropriate tag `name`
- [x] 3.3 Update user profile history responses to return the locale-appropriate tag `name`
- [x] 3.4 Add backend tests covering locale-sensitive tag output behavior

## 4. Sitewide Frontend Localization

- [x] 4.1 Localize homepage, navigation, and shared layout copy
- [x] 4.2 Localize auth-related pages and auth state messaging
- [x] 4.3 Localize match list, match detail, team/player pages, and public browsing states
- [x] 4.4 Localize check-in UI states, form labels, validation messages, and saved-record messaging
- [x] 4.5 Localize `/me` profile page copy, history labels, and loading/empty/error states

## 5. Locale Switching Experience

- [x] 5.1 Ensure changing language updates current-page frontend copy immediately
- [x] 5.2 Ensure pages with locale-aware server data refresh the relevant sections after language change
- [x] 5.3 Verify UGC remains untranslated while surrounding UI switches languages

## 6. Validation

- [x] 6.1 Run backend `go test ./...` and `go build ./...`
- [x] 6.2 Run frontend `npm run lint`, `npm run type-check`, and `npm run build`
- [x] 6.3 Manually verify homepage language switching between English and Chinese
- [x] 6.4 Manually verify `login -> /matches/:id -> check-in -> /me` in both locales with consistent UI copy and localized tag labels
