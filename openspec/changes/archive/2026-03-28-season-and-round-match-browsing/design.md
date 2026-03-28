## Context

Final Whistle currently exposes public match browsing as a single flat list. That was sufficient for the first seeded slice, but it no longer matches the product direction: users need browsing that feels closer to how football fixtures are actually remembered, namely by season and round. The browsing layer also needs to support the refreshed visual language from `a-first-experience-refresh` and coexist with the sitewide i18n foundation.

The change is primarily frontend-facing, but it crosses DTO design, public list response handling, page information architecture, and future extensibility. We also already know that a later change will add team-centric browsing, so this design should avoid painting the public browsing model into a corner.

## Goals / Non-Goals

**Goals:**
- Present `/matches` as a season-first browsing experience that defaults to the latest season and groups fixtures by round within the selected season.
- Keep the match detail route as the destination for each fixture while making the list page easier to scan.
- Preserve locale-aware competition, round, team, and status labels throughout the grouped browsing experience.
- Create a match-list structure that can later support team-centric entry points without redesigning the page model again.
- Strengthen football recognition on the list page through grouped headers, team crests, and more match-sheet-like cards.

**Non-Goals:**
- Implement full team-based browsing, filtering, or a new `/teams` browsing index.
- Introduce a new data-ingestion pipeline or change the underlying source-of-truth for matches.
- Redesign match detail behavior or check-in rules.
- Add pagination/search/filter systems beyond what is needed to support season-round grouping.

## Decisions

### 1. Group matches in the frontend from stable list data
The public list API already returns each match with competition, season, round, teams, kickoff time, and aggregates. For this change, the frontend will group the returned fixtures by season and then by round using stable locale-aware display values already provided by the API.

Why this over a brand-new grouped endpoint:
- It minimizes backend contract churn for the first iteration.
- It keeps the change scoped to browsing architecture rather than introducing a parallel list response format.
- It still allows a later grouped API if real data volume makes server-side grouping worthwhile.

Alternative considered:
- Return nested grouped sections directly from the backend. Rejected for this change because the current seeded scale does not require it and it would create more migration complexity before the grouping model itself is proven.

### 2. Replace flat pagination semantics with grouped browsing loading
The current flat `page/pageSize` semantics are not compatible with season-round grouping because they can split a round or season across pages. For this change, the public `/matches` experience will load the full seeded season-round browsing dataset needed for coherent grouping instead of preserving the current flat pagination behavior.

Why this over preserving the current pagination model:
- The user experience goal is coherent football browsing, and partial round groups directly undermine it.
- The current seeded data size is small enough that loading the grouped set is acceptable.
- It avoids introducing a more complex partial-group paging system before the grouped model is validated.

Alternative considered:
- Keep `page/pageSize` and accept partial round groups. Rejected because it produces confusing browsing sections and breaks the core information architecture.
- Build a new paginated grouped API in the same change. Rejected because it adds unnecessary backend complexity before the grouped browsing model itself is proven.

### 3. Sort seasons and rounds by football-first recency
Season sections will be ordered from newest to oldest, and rounds inside a season will be ordered from earliest to latest so users can read a season as a progressing fixture journey. If a match is missing round information, it will be grouped under an explicit fallback bucket rather than silently mixed into an adjacent round.

Why this over alphabetical or undefined ordering:
- Season recency is the most natural top-level browsing mental model.
- Round progression inside a season is easier to scan in ascending order.
- An explicit fallback bucket prevents silent mis-grouping when round data is absent.

Alternative considered:
- Alphabetical season sorting. Rejected because it does not reflect how users browse football competitions.
- Newest rounds first within each season. Rejected for the first version because ascending round order better supports “read the season” scanning on a grouped page.

### 4. Introduce a dedicated browsing capability for season-round organization
This change will not only tweak the existing list page; it formalizes season-and-round browsing as its own capability because the information architecture matters beyond one component. That lets future team-based browsing build on the same capability instead of overloading `public-match-browsing` with every discovery concern.

Why this over only modifying `public-match-browsing`:
- The new grouped browsing model is conceptually larger than a visual refresh.
- It creates a cleaner place to describe future-compatible browsing requirements.
- It keeps detail-page behavior and list-page organization from being conflated.

### 5. Design the list page around a latest-season default plus season switching
The initial browsing structure will be:
- season switcher at the top of `/matches`
- latest season selected by default
- round groups inside the selected season
- ordered fixture cards inside each round

Why this over rendering every season section in one long page:
- It keeps the default browsing surface focused on the most relevant current fixtures.
- It prevents the page from becoming excessively long as more seasons are added.
- It still preserves a football-native browsing model while making season changes explicit and user-controlled.

Alternative considered:
- Render all seasons in one long page. Rejected because it becomes noisy as soon as multiple seasons are available and pushes too much historical data into the first screen.

### 6. Leave explicit extension hooks for future team-centric browsing
This change will not implement team-based browsing, but it will:
- keep grouped data helpers generic enough to regroup by another dimension later
- avoid copy that implies round grouping is the only permanent browsing mode
- preserve consistent match-card composition that can be reused in team-based contexts

Alternative considered:
- Start team-based browsing now. Rejected because it would dilute the current problem statement and turn one structural cleanup into a much larger discovery-system change.

## Risks / Trade-offs

- **[Risk]** Frontend-only grouping may eventually become awkward with larger real-data volume.  
  **Mitigation:** Keep grouping helpers isolated so the page can switch to a grouped API later without rewriting the layout.

- **[Risk]** Loading the full grouped browsing dataset could become expensive once real-data volume grows.  
  **Mitigation:** Treat this as a seeded-scale solution and keep the API/list abstraction ready for a later grouped or season-scoped backend endpoint.

- **[Risk]** Round labels may vary in format across data sources or locales.  
  **Mitigation:** Treat round text as display data from the API and centralize sorting/label grouping assumptions in one helper layer.

- **[Risk]** Missing round values could produce awkward or invisible fixture placement.  
  **Mitigation:** Define an explicit fallback round bucket and cover it in frontend grouping tests and empty-state rendering.

- **[Risk]** A season switcher may feel unnecessary with only one seeded season.  
  **Mitigation:** Keep the switcher compact and allow it to collapse into a single static label when only one season exists.

- **[Risk]** This change could drift into team-centric browsing prematurely.  
  **Mitigation:** Explicitly keep team browsing out of scope and only add extension-friendly structure, not new team entry flows.
