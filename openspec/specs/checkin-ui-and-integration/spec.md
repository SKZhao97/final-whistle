## ADDED Requirements

### Requirement: Match detail check-in entry states
The frontend SHALL present the match-detail recording surface as `My Match Record`, with state-specific entry behavior based on authentication and existing record status.

#### Scenario: Signed-out user views match detail
- **WHEN** a signed-out user opens a match detail page
- **THEN** the page SHALL continue to render the public match detail content
- **AND** it SHALL NOT attempt an authenticated check-in fetch
- **AND** it SHALL present `My Match Record` as a login-oriented invitation to save ratings, tags, and a short review into the user’s archive

#### Scenario: Signed-in user without a check-in views match detail
- **WHEN** a signed-in user opens a finished match detail page and `GET /matches/:id/my-checkin` returns `data: null`
- **THEN** the page SHALL present `My Match Record` as the primary action surface for starting a new record for that match
- **AND** the entry state SHALL emphasize creating the user’s own saved match record rather than browsing community content

#### Scenario: Signed-in user with an existing check-in views match detail
- **WHEN** a signed-in user opens a finished match detail page and `GET /matches/:id/my-checkin` returns an existing record
- **THEN** the page SHALL present the saved record as an archive-oriented state within `My Match Record`
- **AND** it SHALL expose an edit path for updating the same record

### Requirement: Check-in form create and edit flow
The frontend SHALL provide one `My Match Record` editing experience for create and edit behavior using the existing backend API contract, ordered around post-match expression instead of raw storage order.

#### Scenario: Open create flow
- **WHEN** a signed-in user starts a new check-in from the match detail page
- **THEN** the frontend SHALL open the check-in form with create-mode defaults for the current match
- **AND** it SHALL use the match detail roster to provide the available player-rating choices

#### Scenario: Open edit flow
- **WHEN** a signed-in user chooses to edit an existing match record
- **THEN** the frontend SHALL open the same check-in form prefilled from the current-user check-in payload

#### Scenario: Editing flow follows expression-first structure
- **WHEN** the user is in the record editing state
- **THEN** the form SHALL group and order the inputs around match rating, allegiance/team ratings, tags, short review, player ratings, and then viewing metadata
- **AND** it SHALL avoid presenting the experience as a generic administrative form block

#### Scenario: Submit create flow
- **WHEN** the user submits a valid create-mode form
- **THEN** the frontend SHALL call `POST /matches/:id/checkin`
- **AND** it SHALL use the response payload to transition the page into the saved-record state

#### Scenario: Submit edit flow
- **WHEN** the user submits a valid edit-mode form
- **THEN** the frontend SHALL call `PUT /matches/:id/checkin`
- **AND** it SHALL use the response payload to update the saved-record state for that match

### Requirement: Frontend validation and submit feedback
The check-in form SHALL provide immediate validation for basic input rules and clear submit-state feedback.

#### Scenario: Validate required and bounded fields before submit
- **WHEN** the user enters invalid ratings, exceeds text limits, or selects duplicate players
- **THEN** the frontend SHALL prevent submit
- **AND** it SHALL present actionable validation feedback in the form

#### Scenario: Limit player selection to the match roster
- **WHEN** the user adds or edits player ratings in the form
- **THEN** the frontend SHALL only offer players that belong to the current match roster returned by the match detail response

#### Scenario: Surface backend validation failures
- **WHEN** the backend rejects a create or update request with a validation or business-rule error
- **THEN** the frontend SHALL keep the user in the check-in flow
- **AND** it SHALL surface the returned error in a form-visible way rather than silently failing

#### Scenario: Show submit progress
- **WHEN** the user submits the form
- **THEN** the frontend SHALL show a submitting state that prevents accidental duplicate submission until the request settles

### Requirement: Match detail state refresh after successful submit
After a successful create or update, the match detail experience SHALL transition into a saved record state that feels owned and archived without requiring a manual browser refresh.

#### Scenario: Refresh after create
- **WHEN** a create request succeeds
- **THEN** the page SHALL show the newly saved user record within `My Match Record`
- **AND** it SHALL refresh the relevant match-detail state needed to keep the page consistent with backend data
- **AND** the resulting state SHALL read as a saved archive entry rather than a still-open form

#### Scenario: Refresh after update
- **WHEN** an update request succeeds
- **THEN** the page SHALL show the updated user record in the saved state
- **AND** it SHALL NOT continue displaying stale pre-edit values

### Requirement: Non-finished match handling in the UI
The frontend SHALL not present an active check-in submission path for matches that are not finished.

#### Scenario: Scheduled match detail
- **WHEN** a user views a match detail page for a match whose status is not `FINISHED`
- **THEN** the page SHALL not present an active create or edit submission path
- **AND** it SHALL communicate that check-ins are only available after the match is finished
