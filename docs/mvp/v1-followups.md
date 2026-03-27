# V1 Follow-ups

## Language Consistency

- The current UI mixes English page copy with Chinese tag names from seed data.
- A future localization pass should make the active locale consistent across the whole product.
- Tag labels should come from locale-aware content instead of assuming one fixed language for all users.

## Release Quality

- Keep prioritizing interaction stability before visual redesign.
- If frontend and backend contracts change together, restart both dev servers before manual verification to avoid stale-process mismatches.
