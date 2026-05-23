# Feature Specification: Phase 3 — Hints Module Migration

**Feature Branch**: `003-complete-hints-module`

**Created**: 2026-05-20

**Status**: Draft

**Input**: Complete the **hints** module for Phase 3 of the Java→Go MDM migration
(`serverBackendGo/docs/MIGRATION.md`, `NEXT_STEPS.md` priority #3) so authenticated
users can load, reset, and permanently dismiss UI tutorial hints without the Java
WAR.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Restore hint history after login (Priority: P1)

An authenticated user opens the MDM web shell. The client loads which tutorial
hint keys were already shown so guided tours do not repeat unnecessarily.

**Why this priority**: Legacy Angular `hintService` and the React **Hints** admin
page both call `GET /private/hints/history` on startup or refresh.

**Independent Test**: Log in → call history endpoint → receive a list of hint key
strings (possibly empty for a new user).

**Acceptance Scenarios**:

1. **Given** a logged-in user with no prior hints, **When** they request hint
   history, **Then** the response is `OK` with an empty list (or empty array).
2. **Given** a user who previously marked hints as shown, **When** they request
   history, **Then** all stored hint keys for that user are returned as strings.
3. **Given** no authenticated session, **When** history is requested,
   **Then** access is denied (forbidden), matching protected private routes.

---

### User Story 2 — Record a hint as shown during a tour (Priority: P1)

While using in-app guided steps, each completed step is persisted so the same hint
is not offered again on the next visit.

**Why this priority**: Java `HintResource.markHintAsShown` and legacy `hintService`
POST each `hintKey` when the user advances a step.

**Independent Test**: Authenticated POST with a hint key → history list includes
that key on subsequent GET.

**Acceptance Scenarios**:

1. **Given** a logged-in user and a valid hint key, **When** they mark the hint
   shown, **Then** the operation succeeds with `OK` and the key appears in history.
2. **Given** the same hint key posted twice, **When** marked again,
   **Then** behavior matches Java (no duplicate rows / idempotent success).
3. **Given** an empty or missing hint key, **When** posted,
   **Then** the request fails with a clear error envelope (not a server crash).

---

### User Story 3 — Re-enable tutorials from Settings (Priority: P2)

An administrator with access to the Hints screen chooses **Enable hints** to clear
their personal hint history so tutorials can run again.

**Why this priority**: React `hintsService.enableHints()` and Java Settings UI
depend on `POST /private/hints/enable`.

**Independent Test**: Mark several hints → enable → history is empty again.

**Acceptance Scenarios**:

1. **Given** a user with existing hint history, **When** they enable hints,
   **Then** history is cleared and subsequent GET returns an empty list.
2. **Given** a user with no history, **When** they enable hints,
   **Then** the operation still returns `OK`.

---

### User Story 4 — Disable all tutorials (Priority: P2)

An administrator chooses **Disable hints** so every known tutorial is treated as
already shown (no more popups).

**Why this priority**: Java `disableHints` clears history then inserts all keys from
`userHintTypes`; React exposes the same button.

**Independent Test**: Enable hints → disable → GET history contains all catalog keys.

**Acceptance Scenarios**:

1. **Given** a logged-in user, **When** they disable hints,
   **Then** hint history includes every hint key defined in the catalog table.
2. **Given** hints were already partially shown, **When** disable runs,
   **Then** the result matches Java (full catalog marked shown, not a partial list).

---

### User Story 5 — Verifiable API for migration sign-off (Priority: P2)

A developer validates Phase 3 using interactive API docs and automated tests
without the Java server.

**Why this priority**: Same confidence pattern as Phase 2 (users/roles).

**Independent Test**: Swagger Authorize with JWT → exercise all four endpoints →
`go test` for hints module passes.

**Acceptance Scenarios**:

1. **Given** Swagger with Bearer auth, **When** all hints endpoints are executed
   with an admin token, **Then** responses match Java envelope and status behavior.
2. **Given** the hints module test suite, **When** run locally,
   **Then** history, mark-shown, enable, and disable flows are covered.

---

### Edge Cases

- Concurrent mark-shown for the same user/key → no duplicate constraint violations.
- User deleted (cascade) → hint rows removed with user FK.
- `userHintTypes` empty in dev DB → disable still succeeds (empty or zero keys).
- Very long `hintKey` (>100 chars) → rejected or truncated per Java column limit.
- Internal errors → `ERROR` envelope with `error.internal.server`, not raw stack traces.

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Module**: `internal/modules/hints/` (replace current scaffold).
- **Java reference**: `com.hmdm.rest.resource.HintResource`, `UserDAO` hint methods.
- **REST base**: `/rest/private/hints` — paths must match legacy:
  - `GET /history`
  - `POST /history` (body: hint key string)
  - `POST /enable`
  - `POST /disable`
- **Parity doc**: `serverBackendGo/docs/parity/hints.md`.
- **Layers**: `domain` (hint key, catalog), `port` (repository), `application`
  (use cases), `adapter/http`, `adapter/persistence/postgres`.
- **Migration**: Add `userHints` + `userHintTypes` tables to Go migrations if absent
  from `000001_init` (Liquibase changeSets 19.09.19).
- **Auth**: Same as other private routes — session cookie and/or JWT Bearer;
  Swagger `@Security BearerAuth` on all hints handlers.

### Functional Requirements

- **FR-001**: System MUST expose `GET /rest/private/hints/history` returning
  `OK` with `data` as `string[]` of hint keys for the current user.
- **FR-002**: System MUST expose `POST /rest/private/hints/history` accepting a
  hint key and persisting it for the current user.
- **FR-003**: System MUST expose `POST /rest/private/hints/enable` clearing all
  hint history rows for the current user.
- **FR-004**: System MUST expose `POST /rest/private/hints/disable` clearing history
  then recording every key from `userHintTypes` for the current user.
- **FR-005**: All endpoints MUST require an authenticated principal (403/401
  consistent with existing Go middleware).
- **FR-006**: Responses MUST use the Headwind envelope (`status`, `message`, `data`).
- **FR-007**: System MUST include automated tests for application logic and at least
  one HTTP handler test with authenticated context.
- **FR-008**: System MUST document endpoints in Swagger with Bearer authorization.
- **FR-009**: React Hints page (`/hints`) MUST work against Go-only backend for
  history, enable, and disable without code changes.

### Key Entities

- **UserHint**: A record that user X has seen hint key Y (`userId`, `hintKey`,
  optional `created` timestamp).
- **UserHintType**: Catalog of all known hint keys (`hintKey` unique) used by disable.
- **Hint history (API view)**: List of hint key strings for one user (not full rows).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: After login, a user can load hint history in under 2 seconds on a
  seeded dev database.
- **SC-002**: Marking a hint shown persists across logout/login for the same account.
- **SC-003**: Enable and disable behave indistinguishably from Java for the same
  sequence of actions (verified by parity smoke checklist).
- **SC-004**: React Hints screen completes all three actions without Java backend.
- **SC-005**: 100% of in-scope HintResource endpoints are listed as Done in parity docs.

## Assumptions

- Phase 1 auth and JWT/Swagger Authorize from Phase 2 remain available.
- Hint keys are opaque strings defined by the frontend/legacy UI (e.g. `hint.step.1`);
  no new hint content is authored in this migration.
- `POST /history` body format matches Java (plain string or JSON string per legacy
  client); implementation will follow `HintResource` and Angular `hint.service.js`.
- Settings UI tour integration in React may still be partial; this feature delivers
  API parity and the dedicated Hints admin page first.
- No admin cross-user hint management (only current user), matching Java.

## Dependencies

- Completed: auth, users/current, settings (partial Phase 3).
- Blocks: none for devices/groups phase.
- Database: `userHints`, `userHintTypes` tables and seed hint keys from Liquibase.
