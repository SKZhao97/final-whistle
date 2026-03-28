# Spec9 Review: Season and Round Match Browsing

**Spec File:** `openspec/changes/season-and-round-match-browsing/specs/season-round-match-browsing/spec.md`
**Created:** 2026-03-28

## Summary

This spec introduces a season-first browsing experience for public matches, grouping fixtures by round within each season section. The spec is concise and focuses on user experience goals but lacks technical implementation details that are critical for proper development.

## Findings

### 1. API Response Structure Undefined
The spec does not specify whether the backend should return grouped data or if the frontend should handle grouping logic. Currently, `/matches` returns a flat paginated list (`MatchListItem[]`). For efficient grouping, the API likely needs to return a structured response like:
```typescript
{
  seasons: Array<{
    season: string;
    rounds: Array<{
      round: string;
      matches: MatchListItem[];
    }>;
  }>;
}
```
Without this specification, frontend and backend implementations may diverge.

### 2. Pagination Conflicts with Grouping
Current pagination (`page`, `pageSize`) works against season-round grouping. A user on page 2 might see partial rounds or missing seasons. The spec needs to clarify:
- Should we load all matches at once for complete grouping?
- If pagination is kept, how should it work with grouped data?
- What's the expected performance impact?

### 3. Missing Season Sorting Rules
The spec states "render one or more season sections" but doesn't define:
- Sort order (newest to oldest? alphabetical?)
- Whether to show all seasons or just those with matches in the current filter
- How to handle competitions with different season formats (calendar year vs academic year)

### 4. Round Display Edge Cases
From the seed data, rounds are optional (e.g., `Matchday 1`). The spec doesn't address:
- How to display matches without round information
- Whether to collapse "round-less" matches into a separate section
- Handling of localized round names (already supported via `localizedRound` in backend)

### 5. Frontend Component Reusability Concerns
The current `/matches/page.tsx` renders a flat grid of match cards. Moving to season-round grouping requires:
- New layout components for season headers and round sections
- Potential modifications to the existing match card component
- CSS adjustments for the nested grouping structure

### 6. Performance Implications
If frontend loads all matches for grouping (bypassing pagination):
- Initial load time increases with match count
- Memory usage grows linearly with data volume
- Mobile users on slow connections face degraded experience

### 7. Backward Compatibility Risk
Users may have bookmarks or expectations based on the current flat list view. The spec doesn't mention:
- URL structure changes (query parameters for season/round navigation)
- Browser history and deep linking
- Screen reader accessibility for the new grouped layout

### 8. Missing Empty State Definitions
No scenarios describe what happens when:
- A season has no matches (after filtering)
  - A round has no finished matches (only scheduled)
- The entire dataset is empty

### 9. Test Coverage Gaps
Implementation will need:
- Backend tests for new grouping queries or response structures
- Frontend unit tests for grouping logic utilities
- Integration tests verifying season-round navigation
- E2E tests covering the new browsing flow

### 10. Documentation Updates Required
- API documentation for any new endpoint or modified response
- Frontend component documentation for new grouping components
- User-facing help content explaining the new browsing model

## Recommendations

1. **Clarify Technical Approach**: Specify whether grouping happens server-side or client-side.
2. **Define API Contract**: Provide example request/response structures for the grouped view.
3. **Address Performance**: Consider server-side grouping with optional pagination (e.g., load one season at a time).
4. **Plan Migration**: Outline steps to transition from current flat view to grouped view without breaking user experience.
5. **Add Edge Case Scenarios**: Include handling for missing round data, empty states, and sorting edge cases.

## Overall Assessment

The spec captures the right user experience goal but lacks the technical depth needed for seamless implementation. The core concept aligns well with the project's football domain, but the missing details create implementation risk. Recommend refining with technical specifications before development begins.

**Priority:** Medium
**Complexity:** Medium (frontend restructuring, potential API changes)
**Dependencies:** None beyond existing match browsing infrastructure