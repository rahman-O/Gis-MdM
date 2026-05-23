---
description: "Task list for Phase 5 Applications, Configurations & Config Files migration"
---

# Tasks: Phase 5 — Applications, Configurations & Config Files

**Input**: `specs/006-complete-phase5-apps-config/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Phases 1–4 complete; Postgres via `./scripts/db-up.sh`

**Tests**: Included per FR-X05, FR-X06, and User Story 10.

**Organization**: Tasks grouped by user story for independent delivery and verification.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Applications: `serverBackendGo/internal/modules/applications/`
- Configurations: `serverBackendGo/internal/modules/configurations/`
- Config files: `serverBackendGo/internal/modules/configfiles/`
- Platform: `serverBackendGo/internal/platform/auth/`
- Migrations: `serverBackendGo/db/migrations/`
- Docs: `serverBackendGo/docs/parity/`

---

## Phase 1: Setup

**Purpose**: Confirm Phase 5 context and Java/React parity baseline.

- [X] T001 Verify feature context in `specs/006-complete-phase5-apps-config/spec.md` against `serverBackendGo/docs/MIGRATION.md` Phase 5 pending row
- [X] T002 [P] Review Java `backend/server/src/main/java/com/hmdm/rest/resource/ConfigurationResource.java` and `ApplicationResource.java` for endpoint checklist
- [X] T003 [P] Review React `frontend/src/features/configurations/configurationService.ts` and `frontend/src/features/applications/services/applicationService.ts` for required JSON shapes
- [X] T004 Run baseline `cd serverBackendGo && go build ./...` and note `applications`/`configfiles` scaffolds and Phase 4 `configurations` list-only state

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema, permissions, and module skeletons required by all user stories.

**⚠️ CRITICAL**: No configuration/application endpoint work until migration `000007` applies.

- [X] T005 Create `serverBackendGo/db/migrations/000007_applications_configurations_core.up.sql` per `data-model.md` (`applications`, `applicationversions`, `configurationapplications`, `configurationfiles`, `configurationapplicationsettings`, extended `configurations` columns)
- [X] T006 [P] Create `serverBackendGo/db/migrations/000007_applications_configurations_core.down.sql`
- [X] T007 Add dev seed in `000007`: sample application + version, link to default configuration, extend default configuration columns in `serverBackendGo/db/migrations/000007_applications_configurations_core.up.sql`
- [X] T008 Extend `serverBackendGo/internal/platform/auth/permissions.go` with `PermApplications`, `PermConfigurations` constants and helper methods
- [X] T009 [P] Add permission tests in `serverBackendGo/internal/platform/auth/permissions_test.go` for applications and configurations
- [X] T010 [P] Seed `permissions` rows and `userrolepermissions` for `applications` and `configurations` (role 2 org-admin) in `000007` migration
- [X] T011 [P] Create `serverBackendGo/internal/modules/configurations/domain/configuration.go` with `Configuration`, `LookupItem`, copy/upgrade payloads per `contracts/configurations-api.md`
- [X] T012 [P] Create `serverBackendGo/internal/modules/applications/domain/application.go` with `Application`, `ApplicationVersion`, link payloads per `contracts/applications-api.md`
- [X] T013 Define `serverBackendGo/internal/modules/configurations/port/repository.go` and `push.go` (no-op notifier)
- [X] T014 Define `serverBackendGo/internal/modules/applications/port/repository.go`
- [X] T015 Verify migration: `make dev` and confirm `applications`, `applicationversions`, `configurationapplications` tables exist; `GET /configurations/list` still OK
- [X] T016 [P] Move Phase 4 list logic from `serverBackendGo/internal/modules/configurations/adapter/http/handler.go` into `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go` `ListByCustomer` (preserve behavior)

**Checkpoint**: Schema + permissions + domain ports ready — proceed to user stories.

---

## Phase 3: User Story 1 — Browse and open configurations (Priority: P1) 🎯 MVP (part 1)

**Goal**: `GET /rest/private/configurations/search` and `GET /{id}` for Configurations page.

**Independent Test**: JWT with `configurations` permission → search returns rows; GET id returns editor payload.

### Tests for User Story 1

- [X] T017 [P] [US1] Add `serverBackendGo/internal/modules/configurations/application/service_test.go` for search and get-by-id permission denied
- [X] T018 [P] [US1] Add `serverBackendGo/internal/modules/configurations/adapter/http/handler_test.go` for GET `/search` and GET `/:id` 200 with principal

### Implementation for User Story 1

- [X] T019 [US1] Implement `Search`, `GetByID`, `ListNames` in `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go`
- [X] T020 [US1] Implement `Search`, `GetByID` in `serverBackendGo/internal/modules/configurations/application/service.go` with `configurations` permission checks
- [X] T021 [US1] Refactor `serverBackendGo/internal/modules/configurations/adapter/http/handler.go` — register `GET /search`, `GET /:id`, retain `GET /list` delegating to service
- [X] T022 [US1] Update `serverBackendGo/internal/modules/configurations/module.go` — wire repo → service → handler (replace list-only wiring)

**Checkpoint**: Configurations list + detail load in Swagger/React.

---

## Phase 4: User Story 2 — Browse application catalog (Priority: P1) 🎯 MVP (part 2)

**Goal**: `GET /rest/private/applications/search` and version list for Applications page.

**Independent Test**: JWT with `applications` permission → search returns apps; GET `/{id}/versions` returns versions.

### Tests for User Story 2

- [X] T023 [P] [US2] Add `serverBackendGo/internal/modules/applications/application/service_test.go` for search and permission denied
- [X] T024 [P] [US2] Add `serverBackendGo/internal/modules/applications/adapter/http/handler_test.go` for GET `/search` and GET `/:id/versions`

### Implementation for User Story 2

- [X] T025 [US2] Implement `serverBackendGo/internal/modules/applications/adapter/persistence/postgres/application_repo.go` — `Search`, `SearchByValue`, `GetByID`, `ListVersions`
- [X] T026 [US2] Implement `Search`, `GetByID`, `ListVersions` in `serverBackendGo/internal/modules/applications/application/service.go`
- [X] T027 [US2] Create `serverBackendGo/internal/modules/applications/adapter/http/handler.go` with search and detail routes
- [X] T028 [US2] Wire `serverBackendGo/internal/modules/applications/module.go` — repo → service → handler on `/applications`

**Checkpoint**: Applications page list + versions drill-down works.

---

## Phase 5: User Story 9 — Configuration autocomplete and name list (Priority: P3)

**Goal**: `POST /autocomplete` and `GET /search/{value}`; retain `GET /list`.

**Independent Test**: POST autocomplete → suggestions; GET search/{value} → filtered list.

### Implementation for User Story 9

- [X] T029 [US9] Add `SearchByValue`, `Autocomplete` to `config_repo.go` and `application/service.go` in `serverBackendGo/internal/modules/configurations/`
- [X] T030 [US9] Register `GET /search/:value` and `POST /autocomplete` in `serverBackendGo/internal/modules/configurations/adapter/http/handler.go`

**Checkpoint**: Configuration pickers and filtered search work.

---

## Phase 6: User Story 3 — Configuration CRUD (Priority: P2)

**Goal**: `PUT /`, `DELETE /{id}`, `PUT /copy` with nested apps/files persistence.

**Independent Test**: PUT create → in search; PUT update → persisted; DELETE when no devices; copy → new name.

### Tests for User Story 3

- [X] T031 [P] [US3] Extend `configurations/application/service_test.go` for duplicate name and delete-not-empty
- [X] T032 [P] [US3] Extend `configurations/adapter/http/handler_test.go` for PUT and DELETE 403 without permission

### Implementation for User Story 3

- [X] T033 [US3] Implement `Insert`, `Update`, `Delete`, `Copy`, `CountDevicesUsing` in `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go` (transactional child rows)
- [X] T034 [US3] Implement `Save`, `Delete`, `Copy` in `serverBackendGo/internal/modules/configurations/application/service.go`
- [X] T035 [US3] Add `PUT /`, `DELETE /:id`, `PUT /copy` handlers with Swagger in `serverBackendGo/internal/modules/configurations/adapter/http/handler.go`

**Checkpoint**: Configuration create/edit/delete/copy smoke passes.

---

## Phase 7: User Story 5 — Assign applications to configurations (Priority: P2)

**Goal**: Configuration applications tab + upgrade endpoint.

**Independent Test**: GET `/applications/{configId}` → linked apps; PUT `/application/upgrade` → version updated.

### Implementation for User Story 5

- [X] T036 [US5] Add `ListConfigurationApplications`, `ListAllApplicationsForPicker`, `UpgradeApplication` in `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go`
- [X] T037 [US5] Implement application tab use cases in `serverBackendGo/internal/modules/configurations/application/service.go` (inject applications read port or shared repo)
- [X] T038 [US5] Register `GET /applications`, `GET /applications/:id`, `PUT /application/upgrade` in `serverBackendGo/internal/modules/configurations/adapter/http/handler.go`

**Checkpoint**: Configuration editor Applications tab loads and upgrade returns OK.

---

## Phase 8: User Story 4 — Manage Android and web applications (Priority: P2)

**Goal**: Full application CRUD and validatePkg.

**Independent Test**: PUT android → in search; PUT versions → listed; validatePkg returns conflicts; DELETE works.

### Tests for User Story 4

- [X] T039 [P] [US4] Extend `applications/application/service_test.go` for validatePkg and delete
- [X] T040 [P] [US4] Extend `applications/adapter/http/handler_test.go` for PUT `/android` forbidden without permission

### Implementation for User Story 4

- [X] T041 [US4] Extend `application_repo.go` with `SaveAndroid`, `SaveWeb`, `SaveVersion`, `DeleteApp`, `DeleteVersion`, `ValidatePkg`, `FindByPkg`
- [X] T042 [US4] Implement mutation use cases in `serverBackendGo/internal/modules/applications/application/service.go`
- [X] T043 [US4] Register `PUT /android`, `PUT /web`, `PUT /versions`, `DELETE /:id`, `DELETE /versions/:id`, `PUT /validatePkg`, `POST /autocomplete` in `serverBackendGo/internal/modules/applications/adapter/http/handler.go`

**Checkpoint**: Application create/edit/delete and versions page mutations work.

---

## Phase 9: User Story 6 — Link configurations to applications (Priority: P2)

**Goal**: Application ↔ configuration link dialogs from Applications UI.

**Independent Test**: GET `/applications/configurations/{id}` → links; POST `/applications/configurations` → persisted.

### Implementation for User Story 6

- [X] T044 [US6] Add link query/update methods in `serverBackendGo/internal/modules/applications/adapter/persistence/postgres/application_repo.go` for app and version junctions
- [X] T045 [US6] Implement `GetAppConfigurations`, `UpdateAppConfigurations`, version variants in `serverBackendGo/internal/modules/applications/application/service.go`
- [X] T046 [US6] Register `GET /configurations/:id`, `POST /configurations`, `GET /version/:versionId/configurations`, `POST /version/configurations` in `serverBackendGo/internal/modules/applications/adapter/http/handler.go`

**Checkpoint**: Link configurations dialog saves and reloads correctly.

---

## Phase 10: User Story 7 — Super-admin shared catalog (Priority: P3)

**Goal**: Admin applications page endpoints.

**Independent Test**: Super-admin → admin search; non-super-admin → 403.

### Implementation for User Story 7

- [X] T047 [US7] Add `AdminSearch`, `TurnIntoCommon` in `serverBackendGo/internal/modules/applications/adapter/persistence/postgres/application_repo.go`
- [X] T048 [US7] Implement admin use cases with `RequireSuperAdmin` in `serverBackendGo/internal/modules/applications/application/service.go`
- [X] T049 [US7] Register `GET /admin/search`, `GET /admin/search/:value`, `GET /admin/common/:id` in `serverBackendGo/internal/modules/applications/adapter/http/handler.go`

**Checkpoint**: Admin Applications page loads for super-admin.

---

## Phase 11: User Story 8 — Upload configuration file assets (Priority: P3)

**Goal**: `POST /rest/private/config-files` multipart upload.

**Independent Test**: POST small file → `FileUploadResult` with path/url; writes under customer filesdir.

### Implementation for User Story 8

- [X] T050 [US8] Create `serverBackendGo/internal/modules/configfiles/adapter/http/handler.go` — multipart parse, disk write under `FILES_DIRECTORY` + `customer.filesdir`
- [X] T051 [US8] Wire `serverBackendGo/internal/modules/configfiles/module.go` on `config-files` route group `/rest/private/config-files`
- [X] T052 [P] [US8] Add `serverBackendGo/internal/modules/configfiles/adapter/http/handler_test.go` for upload 200 with principal (temp dir)

**Checkpoint**: Config file upload returns legacy-shaped result.

---

## Phase 12: User Story 10 — Verifiable API (Priority: P2)

**Goal**: Swagger, tests green, parity docs.

**Independent Test**: `go test ./internal/modules/applications/... ./internal/modules/configurations/...`; Swagger shows tags.

### Tests for User Story 10

- [X] T053 [P] [US10] Run `go test ./internal/modules/applications/... ./internal/modules/configurations/... ./internal/modules/configfiles/...` and fix failures
- [X] T054 [P] [US10] Run `cd serverBackendGo && make swagger` and verify Applications/Configurations/ConfigFiles paths in `internal/platform/httpx/swagger/swagger.json`

### Implementation for User Story 10

- [X] T055 [US10] Create `serverBackendGo/docs/parity/configurations.md` per `contracts/configurations-api.md`
- [X] T056 [US10] Create `serverBackendGo/docs/parity/applications.md` per `contracts/applications-api.md`
- [X] T057 [US10] Create `serverBackendGo/docs/parity/configfiles.md` per `contracts/configfiles-api.md`

**Checkpoint**: Module tests pass; Swagger complete.

---

## Phase 13: Polish & Cross-Cutting Concerns

**Purpose**: Close Phase 5 in roadmap and validate React E2E.

- [X] T058 [P] Update `serverBackendGo/docs/MIGRATION.md` — Phase 5 **done**; Swagger table for Applications/Configurations
- [X] T059 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` — Phase 5 منجز; next Phase 6 files/icons
- [X] T060 Run full `specs/006-complete-phase5-apps-config/quickstart.md` validation end-to-end
- [X] T061 Run `cd serverBackendGo && go build ./...` and `go test ./...` for regression
- [X] T062 [P] Manual E2E: React Configurations + Applications pages against Go-only `:8080`
- [X] T063 [P] Regression: Devices page `GET /configurations/list` still works (SC-006)
- [X] T064 [P] Document **Partial** items (web-ui-files APK upload, push notify, optional device search enrichment) in parity docs

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)** → **Foundational (Phase 2)** → user stories
- **US1** (config read) and **US2** (app read) after Phase 2 — parallelizable after T016
- **US9** after US1 repo exists
- **US3** before **US5** (save graph before upgrade tab)
- **US2** before **US5** picker (applications must exist)
- **US4** after US2 repo
- **US6** after US4 + US3 junction tables
- **US7** after US4
- **US8** independent after Phase 2
- **US10** after handlers exist
- **Polish** last

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|-------|
| US1 | Phase 2 | Config browse |
| US2 | Phase 2 | App browse |
| US9 | US1 | Autocomplete on config repo |
| US3 | US1 | Config mutations |
| US5 | US3, US2 | Apps on configuration |
| US4 | US2 | App mutations |
| US6 | US3, US4 | Bidirectional links |
| US7 | US4 | Admin catalog |
| US8 | Phase 2 | File upload |
| US10 | All handlers | Swagger + parity |

### Parallel Opportunities

- Phase 1: T002, T003
- Phase 2: T006–T014, T016 parallel after T005
- US1 tests T017–T018 parallel; US2 tests T023–T024 parallel after Phase 2
- US10: T053–T054, T055–T057 [P]
- Polish: T058, T059, T062–T064

### Parallel Example: MVP (US1 + US2)

```bash
# After Phase 2 complete:
# Track A: T017–T022 configurations read
# Track B: T023–T028 applications read
# Then quickstart §3 + §4
```

---

## Implementation Strategy

### MVP First (US1 + US2)

1. Phase 1–2: Setup + Foundational (migration critical)
2. Phase 3: US1 configuration search + detail
3. Phase 4: US2 application search + versions
4. **STOP and VALIDATE**: React Configurations list opens; Applications list opens
5. Continue US3 → US4 → US5/US6 → US8 → US10 → Polish

### Incremental Delivery

1. Foundation → US1 + US2 (browse both pages)
2. US9 + US3 + US5 (full Configurations editor)
3. US4 + US6 (full Applications + links)
4. US7 + US8 + US10 + Polish

### Suggested MVP Scope

**Minimum**: Phases 1–2 + **US1** + **US2** + T055–T057 subset + T058 + T063 regression.

Delivers usable Configurations and Applications list/detail without Java.

---

## Notes

- Refactor Phase 4 `configurations` list handler into full module without breaking `GET /list`.
- React uses `PUT` for configuration create/update, not `POST`.
- `/private/web-ui-files` is Phase 6 — document in parity Partial.
- Branch: `006-complete-phase5-apps-config`
