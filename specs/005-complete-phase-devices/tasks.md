---
description: "Task list for Phase 4 Devices & Groups migration"
---

# Tasks: Phase 4 — Devices & Groups Module Migration

**Input**: `specs/005-complete-phase-devices/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Phases 1–3 complete; Postgres via `./scripts/db-up.sh`

**Tests**: Included per FR-X07, FR-X06, and User Story 10.

**Organization**: Tasks grouped by user story for independent delivery and verification.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Devices: `serverBackendGo/internal/modules/devices/`
- Groups: `serverBackendGo/internal/modules/groups/`
- Configurations: `serverBackendGo/internal/modules/configurations/`
- Summary: `serverBackendGo/internal/modules/summary/`
- Platform: `serverBackendGo/internal/platform/auth/`
- Migrations: `serverBackendGo/db/migrations/`
- Docs: `serverBackendGo/docs/parity/`

---

## Phase 1: Setup

**Purpose**: Confirm Phase 4 context and Java/React parity baseline.

- [x] T001 Verify feature context in `specs/005-complete-phase-devices/spec.md` against `serverBackendGo/docs/MIGRATION.md` Phase 4 pending row
- [x] T002 [P] Review Java `backend/server/src/main/java/com/hmdm/rest/resource/DeviceResource.java` and `GroupResource.java` for endpoint checklist
- [x] T003 [P] Review React `frontend/src/features/devices/deviceService.ts` and `frontend/src/features/groups/groupService.ts` for required JSON shapes
- [x] T004 Run baseline `cd serverBackendGo && go build ./...` and note `devices`/`groups`/`configurations` scaffold state

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema, permissions, and shared ports required by all user stories.

**⚠️ CRITICAL**: No device/group endpoint work until migrations apply successfully.

- [x] T005 Create `serverBackendGo/db/migrations/000006_devices_groups_core.up.sql` with `groups`, `devices`, `devicegroups`, `configurations`, optional `deviceapplicationsettings`, `userdevicegroupsaccess` per `data-model.md`
- [x] T006 [P] Create `serverBackendGo/db/migrations/000006_devices_groups_core.down.sql`
- [x] T007 Add dev seed in migration or `db/seed`: default group, default configuration, sample devices for customer id 1 in `serverBackendGo/db/migrations/000006_devices_groups_core.up.sql`
- [x] T008 Extend `serverBackendGo/internal/platform/auth/permissions.go` with `edit_devices`, `edit_device_desc` helpers and permission name constants
- [x] T009 [P] Add `serverBackendGo/internal/platform/auth/permissions_test.go` cases for new permission helpers
- [x] T010 [P] Create `serverBackendGo/internal/modules/devices/domain/device.go` with `Device`, `SearchRequest`, `DeviceListView`, `DeviceView`, bulk payloads per `contracts/devices-api.md`
- [x] T011 [P] Create `serverBackendGo/internal/modules/groups/domain/group.go` with `Group`, `LookupItem` types
- [x] T012 Define `serverBackendGo/internal/modules/devices/port/repository.go` and `push.go` (no-op notifier interface)
- [x] T013 Define `serverBackendGo/internal/modules/groups/port/repository.go`
- [x] T014 Verify migration: restart `make dev` and confirm `devices`, `groups`, `devicegroups`, `configurations` tables exist
- [x] T015 [P] Seed `permissions` rows for `edit_devices` and `edit_device_desc` if missing from `000001_init` (org-admin/super-admin role links)

**Checkpoint**: Tables present; permission helpers compile — proceed to user stories.

---

## Phase 3: User Story 2 — Load device groups (Priority: P1)

**Goal**: `GET /rest/private/groups/search` returns tenant groups for Devices filters.

**Independent Test**: JWT → GET groups/search → non-empty id/name array for seeded tenant.

### Tests for User Story 2

- [x] T016 [P] [US2] Add `serverBackendGo/internal/modules/groups/application/service_test.go` for list by customer and empty tenant
- [x] T017 [P] [US2] Add `serverBackendGo/internal/modules/groups/adapter/http/handler_test.go` for GET `/search` 200 with principal

### Implementation for User Story 2

- [x] T018 [US2] Implement `serverBackendGo/internal/modules/groups/adapter/persistence/postgres/group_repo.go` — `ListByCustomer`, `ListByValue`, `CountDevicesInGroup`
- [x] T019 [US2] Implement `List`, `SearchByValue` in `serverBackendGo/internal/modules/groups/application/service.go`
- [x] T020 [US2] Create `serverBackendGo/internal/modules/groups/adapter/http/handler.go` with `GET /search` and `GET /search/:value`
- [x] T021 [US2] Wire `serverBackendGo/internal/modules/groups/module.go` — repo → service → handler on `groups.Private.Group("/groups")`

**Checkpoint**: `quickstart.md` §3 groups smoke passes.

---

## Phase 4: User Story 1 — Search and browse devices (Priority: P1) 🎯 MVP

**Goal**: `POST /rest/private/devices/search` returns `DeviceListResponse` for React Devices grid.

**Independent Test**: POST search with `pageNum`/`pageSize` → `data.devices.items` + `data.configurations` map.

### Tests for User Story 1

- [x] T022 [P] [US1] Add `serverBackendGo/internal/modules/devices/application/service_test.go` stub for `Search` pagination and tenant scope
- [x] T023 [P] [US1] Add `serverBackendGo/internal/modules/devices/adapter/http/handler_test.go` for POST `/search` 200 envelope shape

### Implementation for User Story 1

- [x] T024 [US1] Implement `Search` + `Count` SQL in `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_repo.go` mirroring `DeviceMapper.xml` (customerId, group access, `pageNum`, status color)
- [x] T025 [US1] Implement `Search` in `serverBackendGo/internal/modules/devices/application/service.go` building `DeviceListView` (`configurations` map + `devices.items` + `totalItemsCount`)
- [x] T026 [US1] Create `serverBackendGo/internal/modules/devices/adapter/http/handler.go` with `Search` handler and `@Router` + `@Security BearerAuth`
- [x] T027 [US1] Register `POST /search` in `serverBackendGo/internal/modules/devices/adapter/http/handler.go` `Register` method
- [x] T028 [US1] Wire `serverBackendGo/internal/modules/devices/module.go` — repo → service → handler; require `deps.DB`

**Checkpoint**: `quickstart.md` §5 device search returns list structure.

---

## Phase 5: User Story 9 — Configuration list (Priority: P2)

**Goal**: `GET /rest/private/configurations/list` for Devices configuration dropdown.

**Independent Test**: GET list → `[{id,name}]` for tenant.

### Implementation for User Story 9

- [x] T029 [P] [US9] Implement `ListByCustomer` in configurations postgres repo or inline query in `serverBackendGo/internal/modules/configurations/adapter/persistence/postgres/config_repo.go`
- [x] T030 [US9] Add `GET /list` handler in `serverBackendGo/internal/modules/configurations/adapter/http/handler.go` with Swagger annotations
- [x] T031 [US9] Update `serverBackendGo/internal/modules/configurations/module.go` to wire list handler (replace scaffold)

**Checkpoint**: `quickstart.md` §4 configurations list succeeds.

---

## Phase 6: User Story 8 — Real dashboard statistics (Priority: P2)

**Goal**: `GET /rest/private/summary/devices` returns real counts when devices exist.

**Independent Test**: After seed devices → summary shows non-zero status breakdown.

### Implementation for User Story 8

- [x] T032 [US8] Port summary count SQL into `serverBackendGo/internal/modules/summary/adapter/persistence/postgres/summary_repo.go` `GetDeviceStats` (replace `EmptyDeviceStats` stub)
- [x] T033 [P] [US8] Add `serverBackendGo/internal/modules/summary/application/service_test.go` case for non-empty stats with seeded data
- [x] T034 [US8] Update `serverBackendGo/docs/parity/summary.md` — mark stats **Done** when SQL live

**Checkpoint**: `quickstart.md` §6 summary shows real numbers.

---

## Phase 7: User Story 4 — Manage device groups (Priority: P2)

**Goal**: PUT/DELETE groups with `settings` permission and duplicate/not-empty rules.

**Independent Test**: PUT create group → appears in search; DELETE empty group OK; DELETE used group → `error.notempty.group`.

### Tests for User Story 4

- [x] T035 [P] [US4] Extend `groups/application/service_test.go` for duplicate name and not-empty delete
- [x] T036 [P] [US4] Extend `groups/adapter/http/handler_test.go` for PUT forbidden without settings permission

### Implementation for User Story 4

- [x] T037 [US4] Extend `group_repo.go` with `Insert`, `Update`, `Delete`, `GetByName`, `AssignGroupToUserOnCreate` (creator access)
- [x] T038 [US4] Implement `Save`, `Delete` in `serverBackendGo/internal/modules/groups/application/service.go`
- [x] T039 [US4] Add `PUT /`, `DELETE /:id`, `POST /autocomplete` handlers in `serverBackendGo/internal/modules/groups/adapter/http/handler.go`

**Checkpoint**: Groups admin CRUD smoke via curl or React groups page.

---

## Phase 8: User Story 3 — Device CRUD (Priority: P2)

**Goal**: GET by number, PUT create/update, DELETE device with `edit_devices` permission.

**Independent Test**: GET number → view; PUT create → in search; DELETE → gone.

### Tests for User Story 3

- [x] T040 [P] [US3] Extend `devices/application/service_test.go` for duplicate number and permission denied on delete
- [x] T041 [P] [US3] Extend `devices/adapter/http/handler_test.go` for GET `/number/:n` and DELETE 403 without `edit_devices`

### Implementation for User Story 3

- [x] T042 [US3] Extend `device_repo.go` with `GetByNumber`, `GetByID`, `Insert`, `Update`, `Delete`, duplicate checks, device limit query
- [x] T043 [US3] Implement `GetByNumber`, `Save`, `Delete` in `serverBackendGo/internal/modules/devices/application/service.go` including bulk config update via `ids` on PUT
- [x] T044 [US3] Add `GetByNumber`, `Save`, `Delete` handlers and routes in `serverBackendGo/internal/modules/devices/adapter/http/handler.go`

**Checkpoint**: Single-device create/edit/delete works in Swagger.

---

## Phase 9: User Story 5 — Bulk device operations (Priority: P2)

**Goal**: `POST /deleteBulk` and `POST /groupBulk` for multi-select actions.

**Independent Test**: Bulk delete removes rows; groupBulk updates `devicegroups` links.

### Tests for User Story 5

- [x] T045 [P] [US5] Extend `devices/application/service_test.go` for `DeleteBulk` and `GroupBulk` set/clear actions

### Implementation for User Story 5

- [x] T046 [US5] Add `DeleteBulk`, `UpdateGroupBulk` repo methods in `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/device_repo.go`
- [x] T047 [US5] Implement bulk use cases in `serverBackendGo/internal/modules/devices/application/service.go`
- [x] T048 [US5] Register `POST /deleteBulk` and `POST /groupBulk` in `serverBackendGo/internal/modules/devices/adapter/http/handler.go`

**Checkpoint**: Bulk operations return OK in Swagger smoke.

---

## Phase 10: User Story 6 — Autocomplete and description (Priority: P3)

**Goal**: Device autocomplete and description POST with `edit_device_desc`.

**Independent Test**: POST autocomplete → suggestions; POST description → persisted.

### Implementation for User Story 6

- [x] T049 [US6] Add `Autocomplete`, `UpdateDescription` to `device_repo.go` and `application/service.go`
- [x] T050 [US6] Add `POST /autocomplete` and `POST /:id/description` handlers in `serverBackendGo/internal/modules/devices/adapter/http/handler.go`

**Checkpoint**: Autocomplete returns up to 10 matches.

---

## Phase 11: User Story 7 — Per-device application settings (Priority: P3)

**Goal**: GET/POST application settings and notify (push stub).

**Independent Test**: GET settings array; POST save OK; notify OK.

### Implementation for User Story 7

- [x] T051 [P] [US7] Implement no-op `serverBackendGo/internal/modules/devices/adapter/persistence/postgres/push_noop.go` or `port/push_noop.go` satisfying `port.PushNotifier`
- [x] T052 [US7] Add app settings repo methods in `device_repo.go` (or `device_app_settings_repo.go`)
- [x] T053 [US7] Implement `GetAppSettings`, `SaveAppSettings`, `NotifyAppSettings` in `application/service.go`
- [x] T054 [US7] Register `GET/POST /:id/applicationSettings` and `POST .../notify` in `serverBackendGo/internal/modules/devices/adapter/http/handler.go`

**Checkpoint**: App settings endpoints return OK.

---

## Phase 12: User Story 10 — Verifiable API (Priority: P2)

**Goal**: Swagger tags, module tests green, parity docs complete.

**Independent Test**: `go test ./internal/modules/devices/... ./internal/modules/groups/...`; Swagger shows Devices/Groups tags.

### Tests for User Story 10

- [x] T055 [P] [US10] Run `go test ./internal/modules/devices/... ./internal/modules/groups/... ./internal/modules/summary/...` and fix failures
- [x] T056 [P] [US10] Add missing Swagger `@Router` / `@Tags` on all new handlers (devices, groups, configurations list)

### Implementation for User Story 10

- [x] T057 [US10] Run `cd serverBackendGo && make swagger` and verify `/private/devices/*`, `/private/groups/*`, `/private/configurations/list` in `internal/platform/httpx/swagger/swagger.json`
- [x] T058 [US10] Create `serverBackendGo/docs/parity/devices.md` with Done/Partial matrix per `contracts/devices-api.md`
- [x] T059 [US10] Create `serverBackendGo/docs/parity/groups.md` with Done matrix per `contracts/groups-api.md`

**Checkpoint**: Module tests pass; Swagger complete.

---

## Phase 13: Polish & Cross-Cutting Concerns

**Purpose**: Close Phase 4 in roadmap and validate React E2E.

- [x] T060 [P] Update `serverBackendGo/docs/MIGRATION.md` — Phase 4 status **done**; Swagger table adds Devices/Groups/Configurations
- [x] T061 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` — Phase 4 منجز; next Phase 5 applications/configurations
- [x] T062 Run full `specs/005-complete-phase-devices/quickstart.md` validation end-to-end
- [x] T063 Run `cd serverBackendGo && go build ./...` and `go test ./...` for regression
- [x] T064 [P] Manual E2E: React Devices page (list, filters, create/edit) and Groups page against Go-only `:8080`
- [x] T065 Document **Partial** search enrichment (apps/files in configuration map) in `serverBackendGo/docs/parity/devices.md` if not fully implemented
- [x] T066 [P] Enable Phase 3 customer create default devices follow-up note in `serverBackendGo/docs/parity/customers.md` (unblocks after Phase 4 schema)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)** → **Foundational (Phase 2)** → user stories
- **US2 (groups list)** before **US1** recommended (filter dropdown) but US1 can proceed after Phase 2 with null groupId
- **US9 (config list)** before React Devices full UX (parallel after US1 repo exists)
- **US8 (summary)** after migration T014
- **US3–US7** after US1 search repo exists
- **US10** after handlers exist
- **Polish** last

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|-------|
| US2 | Phase 2 | Groups list |
| US1 | Phase 2, US2 optional | Device search MVP |
| US9 | Phase 2 | Configurations table |
| US8 | Phase 2 | devices table |
| US4 | US2 repo | Group mutations |
| US3 | US1 repo | Device mutations |
| US5 | US3 | Bulk uses delete/update |
| US6–US7 | US3 | Same device repo |
| US10 | All handlers | Swagger + tests |

### Parallel Opportunities

- Phase 1: T002, T003
- Phase 2: T006–T011, T015 parallel after T005
- Per-story test tasks marked [P]
- Polish: T060, T061, T064, T066

### Parallel Example: User Story 1

```bash
T022 service_test.go Search stubs
T023 handler_test.go POST /search
# then T024–T028 sequential repo → service → http
```

---

## Implementation Strategy

### MVP First (US2 + US1 + US9)

1. Phase 1–2: Setup + Foundational (migrations critical)
2. Phase 3: US2 groups search
3. Phase 4: US1 device search
4. Phase 5: US9 configurations list
5. **STOP and VALIDATE**: React Devices page loads list + dropdowns
6. Phase 6: US8 summary real stats

### Incremental Delivery

1. Foundation → US2 → US1 → US9 → US8 (usable Devices + Dashboard)
2. US4 + US3 (full CRUD)
3. US5 bulk → US6 → US7
4. US10 + Polish → Phase 4 sign-off

### Suggested MVP Scope

**Minimum**: Phases 1–2 + **US2** + **US1** + **US9** + T057–T058 parity subset + T060 MIGRATION update.

Delivers Devices page list and filters without Java.

---

## Notes

- React uses `pageNum` not `currentPage` for device search.
- Push notify returns OK without FCM (Phase 7).
- Full `ConfigurationResource` CRUD is Phase 5.
- Branch: `005-complete-phase-devices`
