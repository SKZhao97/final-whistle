## MODIFIED Requirements

### Requirement: Current-user check-in retrieval and write responses
The backend SHALL return the current user’s check-in record in a locale-consistent response shape.

#### Scenario: Locale-aware tag labels in check-in detail
- **WHEN** `GET /matches/:id/my-checkin`, `POST /matches/:id/checkin`, or `PUT /matches/:id/checkin` returns tag data
- **THEN** each tag item SHALL expose the display `name` for the current locale
- **AND** the rest of the check-in DTO shape SHALL remain unchanged
