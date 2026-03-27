## Why

Final Whistle has completed its v1 functional loop, but the product still feels like a collection of working screens instead of a distinct post-match recording experience. Now that sitewide i18n is in place, the next step is to reshape the core surfaces so users feel "I really captured this match" rather than "I filled out a football form."

## What Changes

- Reframe the match detail page around an A-first hierarchy: `Match Context -> My Match Record -> Community Pulse`.
- Redesign the `My Match Record` area across signed-out, empty, editing, and saved states so it feels like a primary recording surface instead of a secondary utility panel.
- Refresh the `/me` experience from a stats-and-list page into a personal football archive with stronger identity, pattern, and memory framing.
- Introduce a more modern football-oriented visual language built around balanced field-green accents, fixed team crests, and lighter archive/editorial surfaces instead of the current dark utility style.
- Align key football presentation elements across core surfaces, including league branding placement, fixed crest treatment, and more intentional information hierarchy.
- Keep business capabilities stable: no new auth, check-in rules, or data ingestion features are added in this change.

## Capabilities

### New Capabilities
- `a-first-experience`: Defines the cross-surface product direction, visual language, and archive-oriented interaction patterns for Final Whistle's next stage.

### Modified Capabilities
- `public-match-browsing`: Match detail presentation and hierarchy change to emphasize match context, fixed crest usage, and a secondary community pulse layer.
- `checkin-ui-and-integration`: The check-in experience changes from a generic form panel to a first-class "My Match Record" flow with redesigned states and expression-first structure.
- `user-profile-and-aggregation`: The `/me` page requirements change from basic profile/history display to an archive-oriented layout with stronger identity and memory framing.
- `frontend-framework`: Shared frontend requirements expand to cover the new visual language, reusable page sections, and football-oriented presentation primitives used by the refreshed experience.

## Impact

- Affected frontend surfaces include `/matches/[matchId]`, `/me`, shared layout/navigation elements, and the check-in components that power match recording.
- Shared styling tokens, presentation components, and asset treatment will need coordinated updates, including crest/logo placement and refreshed page sections.
- Existing APIs largely stay stable, but frontend rendering contracts must fully preserve the locale-aware labels introduced by `sitewide-i18n`.
