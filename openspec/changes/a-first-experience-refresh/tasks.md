## 1. Shared Visual Foundation

- [x] 1.1 Define the refreshed experience tokens for balanced field-green accents, lighter archive surfaces, and typography hierarchy in the frontend styling layer
- [x] 1.2 Add reusable presentation primitives for fixed team crests, secondary league-brand placement, and archive/editorial section shells
- [x] 1.3 Ensure the refreshed shared primitives remain locale-safe for both English and Chinese labels

## 2. Match Context Hero

- [x] 2.1 Rebuild the match detail hero so team identity, score/status, competition context, and venue/time follow the A-first hierarchy
- [x] 2.2 Integrate fixed team crests into the primary hero structure and place league branding as a secondary element near competition text
- [x] 2.3 Validate the refreshed hero on desktop and mobile breakpoints with localized labels

## 3. My Match Record Experience

- [x] 3.1 Refactor the match-detail record section into a first-class `My Match Record` surface with distinct signed-out, empty, editing, and saved states
- [x] 3.2 Reorder the editing flow around expression-first grouping: match rating, allegiance/team ratings, tags, short review, player ratings, then viewing metadata
- [x] 3.3 Redesign the saved state so it reads as an archived personal record with a clear edit path and stronger completion feeling
- [x] 3.4 Preserve all existing create/edit/non-finished behavior while adapting the interaction copy and hierarchy to the new A-first framing

## 4. Community Pulse and Supporting Match Content

- [x] 4.1 Reposition community aggregates, hot tags, recent reactions, and player board content beneath the record surface as a secondary layer
- [x] 4.2 Adjust spacing, visual weight, and section headings so community content supports rather than competes with `My Match Record`
- [x] 4.3 Verify the refreshed match detail still works cleanly for signed-out, signed-in, and already-recorded scenarios

## 5. `/me` Archive Refresh

- [x] 5.1 Rework `/me` into an archive-oriented layout with identity, patterns, archive/history, and memory-oriented framing
- [x] 5.2 Redesign history items as saved-record entries with stronger match recall and localized shell copy
- [x] 5.3 Improve empty/loading/error states so `/me` still communicates archive value when the user has little or no history

## 6. Verification and Polish

- [x] 6.1 Update or add frontend tests for the refreshed record-state logic and archive-oriented rendering helpers
- [x] 6.2 Run frontend validation (`npm run lint`, `npm run type-check`, `npm run build`) after the refresh is implemented
- [x] 6.3 Perform bilingual manual smoke checks on the refreshed match detail and `/me` flows, including saved-state rendering and crest/league-brand placement
