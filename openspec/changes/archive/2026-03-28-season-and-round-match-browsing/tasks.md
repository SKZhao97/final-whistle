## 1. Grouping Foundation

- [x] 1.1 Review the current public match list response and confirm the season/round fields needed for grouped browsing are already available
- [x] 1.2 Lock the implementation approach for grouping so frontend and backend use one documented strategy for this change
- [x] 1.3 Replace the current flat pagination assumption with a loading strategy that can render complete season-round groups
- [x] 1.4 Add frontend grouping utilities that transform flat match list data into season sections and round groups with stable season/round ordering
- [x] 1.5 Add utility coverage for grouping, fallback round buckets, and ordering behavior

## 2. Match List Information Architecture

- [x] 2.1 Rebuild `/matches` around a latest-season default, season switcher, and round headers instead of a flat card stream
- [x] 2.2 Update fixture cards so team crests, score/status, kickoff context, and aggregate metadata remain scannable inside grouped browsing
- [x] 2.3 Add explicit empty-state handling for fully empty grouped browsing and any fallback round grouping
- [x] 2.4 Ensure grouped browsing copy and structure remain locale-aware and consistent with the current A-first visual language

## 3. Forward-Compatible Browsing Structure

- [x] 3.1 Keep fixture-card and grouping primitives reusable for a future team-based browsing mode
- [x] 3.2 Avoid route or copy assumptions that would block later regrouping by team
- [x] 3.3 Document or encode the browsing hooks needed for a later team-centric change without implementing the feature now

## 4. Validation

- [x] 4.1 Verify the public `/matches` page with seeded data shows season sections and round groupings correctly in English
- [x] 4.2 Verify the same grouped browsing flow in Chinese with locale-aware season/round/team labels
- [x] 4.3 Verify fallback handling for matches without round data and grouped empty-state behavior
- [x] 4.4 Verify grouped fixture cards still navigate correctly into existing public match detail pages
