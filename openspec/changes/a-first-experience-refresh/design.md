## Context

Final Whistle already supports the full v1 product loop, but the core pages still present themselves like functional product screens rather than a distinct post-match recording product. The product direction is now explicitly A-first: the match detail page should feel like a place to leave a personal record, and `/me` should feel like a football archive instead of a generic account page.

This change is cross-cutting because it spans shared frontend layout, the match-detail surface, the check-in component tree, and the `/me` page. It also needs to preserve the sitewide i18n foundation, meaning any refreshed hierarchy or visual system must work with the locale-aware labels already shipped for UI copy, tags, teams, and competitions.

## Goals / Non-Goals

**Goals:**
- Rebuild the match detail page around the A-first hierarchy: `Match Context -> My Match Record -> Community Pulse`.
- Turn the current check-in area into a primary `My Match Record` surface with clear signed-out, empty, editing, and saved states.
- Reframe `/me` as a personal football archive with stronger identity, pattern, and memory structure.
- Introduce a shared visual language that feels more modern and football-native without reverting to a heavy dark sports-app style.
- Codify crest placement, league-logo placement, balanced field-green accents, and archive/editorial surface treatment so implementation can be shared across pages.

**Non-Goals:**
- Adding new business capabilities, endpoints, auth flows, or check-in rules.
- Building a community feed, ranking system, or new social interactions.
- Expanding i18n scope beyond ensuring the refreshed UI fully preserves the localized surface already defined by `sitewide-i18n`.
- Introducing a new design-system package or external component library.

## Decisions

### 1. Use a cross-surface A-first experience layer plus focused delta specs
We will add a new `a-first-experience` capability for the shared product direction and modify the existing match-detail, check-in, profile, and frontend-framework specs for page-specific behavior.

Why:
- The change is bigger than a one-page redesign; it defines the new primary product posture.
- Existing capabilities still own their concrete behavior and should remain the source of truth for those surfaces.

Alternative considered:
- Capture everything only as deltas in existing specs. Rejected because the shared A-first hierarchy and visual language would become fragmented across several files.

### 2. Keep APIs stable and concentrate the change in presentation and interaction structure
The refresh will reuse the existing API surface wherever possible. Match detail, check-in, and profile pages will be reorganized visually and behaviorally without adding new domain endpoints.

Why:
- The value of spec8 is expression, hierarchy, and product clarity, not backend expansion.
- This keeps implementation risk lower and prevents the change from collapsing into another backend-heavy milestone.

Alternative considered:
- Add new aggregate/profile endpoints for richer archive modules immediately. Rejected because it would widen scope and delay the core experience refresh.

### 3. Use a modern football presentation language built on balanced field-green accents, warm light surfaces, and fixed crest treatment
The refreshed UI will not keep the current heavy dark styling. It will instead use warm light surfaces, deeper green accents, fixed team crests in match hero structures, and secondary league-logo placement near competition text.

Why:
- The current dark utility style weakens both the archive feeling and football identity.
- Pure editorial paper styling is not football-native enough; pure pitch styling becomes too sports-template-like.
- The balanced direction preserves record-book clarity while making the product unmistakably football.

Alternative considered:
- Fully dark redesign. Rejected because it repeats the current weakness.
- Strong paper/archive aesthetic with minimal football cues. Rejected because it underplays football identity.

### 4. Treat `My Match Record` as a state machine, not a generic form panel
The match record surface will be designed around four states: signed out, signed in with no record, editing, and saved. The saved state is a first-class destination and must look like an archived record rather than a collapsed form result.

Why:
- The product promise is “I really captured this match,” which is not fulfilled by a utility-like form shell.
- The current product needs a more intentional transition from action to ownership.

Alternative considered:
- Preserve the existing form shell and only restyle it. Rejected because the hierarchy and completion feeling would remain weak.

### 5. Reorder the editing experience by expression flow instead of storage order
Within the record form, fields will be grouped and ordered around the natural post-match expression sequence: match rating, allegiance, team ratings, tags, short review, player ratings, then meta fields like watched type/time.

Why:
- Users remember and express a match emotionally before they think about metadata.
- This order better supports the A-first framing and reduces “admin form” feel.

Alternative considered:
- Keep current field order and only restyle inputs. Rejected because it preserves a data-entry feel.

### 6. Structure `/me` as archive layers rather than a dashboard
The `/me` surface will be organized into identity, patterns, archive, and memory/highlight layers rather than a statistics block plus history list.

Why:
- `/me` is the value-realization page for an A-first product.
- Dashboard-like layouts would dilute the archive/product identity and make the page feel replaceable.

Alternative considered:
- Keep current page structure and only improve styling. Rejected because the product meaning of the page would remain too weak.

## Risks / Trade-offs

- [Risk] The visual refresh could become too broad and turn into an unbounded redesign. → Mitigation: keep the scope to match detail, `/me`, shared presentation primitives, and the record flow; exclude new features and secondary pages.
- [Risk] A stronger visual language could conflict with the freshly shipped i18n system. → Mitigation: require locale-safe primitives and keep all refreshed components built on the existing locale provider and localized data contracts.
- [Risk] Saved-state/archive treatment could over-index on style and fail to improve clarity. → Mitigation: keep the saved state structurally explicit with readable rating/tag/review sections and a clear edit path.
- [Risk] Fixed crest and league-branding treatment could produce uneven layouts across localized labels and future data coverage. → Mitigation: define crest/logo placement as layout primitives with secondary branding weight and test both English and Chinese labels.
