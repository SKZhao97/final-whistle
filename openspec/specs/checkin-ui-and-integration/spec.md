## ADDED Requirements

### Requirement: Match detail check-in entry states
The frontend SHALL present match-detail check-in actions based on the current user’s authentication state and existing record state.

#### Scenario: Signed-out user views match detail
- **WHEN** a signed-out user opens a match detail page
- **THEN** the page SHALL continue to render the public match detail content
- **AND** it SHALL NOT attempt an authenticated check-in fetch
- **AND** it SHALL present a login-oriented entry instead of an active check-in editor

#### Scenario: Signed-in user without a check-in views match detail
- **WHEN** a signed-in user opens a finished match detail page and `GET /matches/:id/my-checkin` returns `data: null`
- **THEN** the page SHALL present an entry point for creating a new check-in for that match

#### Scenario: Signed-in user with an existing check-in views match detail
- **WHEN** a signed-in user opens a finished match detail page and `GET /matches/:id/my-checkin` returns an existing record
- **THEN** the page SHALL present that user’s saved check-in summary state
- **AND** it SHALL expose an edit path for updating the same record

### Requirement: Check-in form create and edit flow
The frontend SHALL provide one check-in form experience that supports both create and edit behavior using the existing backend API contract.

#### Scenario: Open create flow
- **WHEN** a signed-in user starts a new check-in from the match detail page
- **THEN** the frontend SHALL open the check-in form with create-mode defaults for the current match
- **AND** it SHALL use the match detail roster to provide the available player-rating choices

#### Scenario: Open edit flow
- **WHEN** a signed-in user chooses to edit an existing match record
- **THEN** the frontend SHALL open the same check-in form prefilled from the current-user check-in payload

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
After a successful create or update, the match detail experience SHALL reflect the saved current-user record without requiring a manual browser refresh.

#### Scenario: Refresh after create
- **WHEN** a create request succeeds
- **THEN** the page SHALL show the newly saved user record
- **AND** it SHALL refresh the relevant match-detail state needed to keep the page consistent with backend data

#### Scenario: Refresh after update
- **WHEN** an update request succeeds
- **THEN** the page SHALL show the updated user record
- **AND** it SHALL NOT continue displaying stale pre-edit values

### Requirement: Non-finished match handling in the UI
The frontend SHALL not present an active check-in submission path for matches that are not finished.

#### Scenario: Scheduled match detail
- **WHEN** a user views a match detail page for a match whose status is not `FINISHED`
- **THEN** the page SHALL not present an active create or edit submission path
- **AND** it SHALL communicate that check-ins are only available after the match is finished
