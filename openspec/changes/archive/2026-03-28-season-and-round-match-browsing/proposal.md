## Why

The current `/matches` page is a flat stream of cards, which already feels noisy with the seeded Premier League set and will become harder to browse as more seasons and fixtures are added. This change is needed now because match browsing is the main public entry path, and the next product phase depends on a clearer season-first structure before we add richer team-centric browsing.

## What Changes

- Reorganize public match browsing around a season selector with round groups instead of one undifferentiated list.
- Default `/matches` to the latest available season so the browsing surface stays focused as more seasons are added.
- Add a season-first browsing surface that lets users switch seasons and understand where each fixture belongs before opening match detail.
- Preserve the current match detail, team detail, and player detail destinations while changing how users discover matches from `/matches`.
- Introduce explicit information-architecture hooks for future team-based browsing without implementing full team-centric navigation in this change.
- Refresh match-list presentation so grouped fixtures feel scannable, football-native, and consistent with the current A-first visual direction.

## Capabilities

### New Capabilities
- `season-round-match-browsing`: Season-first public browsing that groups matches by round and defines the browsing structure for future extensions.

### Modified Capabilities
- `public-match-browsing`: The public match list and discovery experience changes from a flat list to a grouped season-and-round browsing surface.

## Impact

- Frontend `/matches` page information architecture and presentation
- Public match list API response shape and/or frontend grouping strategy
- Match DTOs and list metadata needed to render season and round groupings
- Seeded browsing experience and manual validation flows for public discovery
