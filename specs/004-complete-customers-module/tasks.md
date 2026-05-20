---
description: "Task list for Phase 3 Customers module migration"
---

# Tasks: Phase 3 — Customers Module Migration

**Input**: `specs/004-complete-customers-module/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Auth + JWT/Swagger Bearer (Phase 1–2); hints/settings/summary Phase 3 partial; Postgres via `./scripts/db-up.sh`

**Tests**: Included per FR-011 and User Story 7 (application + HTTP handler tests).

**Organization**: Tasks grouped by user story for independent delivery and verification.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Go backend: `serverBackendGo/internal/modules/customers/`
- Platform: `serverBackendGo/internal/platform/auth/`
- Migrations: `serverBackendGo/db/migrations/`
- Docs: `serverBackendGo/docs/parity/`, `specs/004-complete-customers-module/quickstart.md`
- Java reference: `backend/server/src/main/java/com/hmdm/rest/resource/CustomerResource.java`
- React: `frontend/src/features/customers/customersService.ts`

---

## Phase 1: Setup

**Purpose**: Confirm feature context and parity baseline before code changes.

- [x] T001 Verify feature context in `specs/004-complete-customers-module/spec.md` and `plan.md` against `serverBackendGo/docs/MIGRATION.md` Phase 3 (customers last open module)
- [x] T002 [P] Review Java `backend/server/src/main/java/com/hmdm/rest/resource/CustomerResource.java` and `backend/common/src/main/java/com/hmdm/persistence/mapper/CustomerMapper.xml` for SQL parity
- [x] T003 [P] Review React `frontend/src/features/customers/customersService.ts` and `ControlPanelPage.tsx` for required fields (`items`, `totalItemsCount`, impersonate payload)
- [x] T004 Run baseline `cd serverBackendGo && go build ./...` and note current `internal/modules/customers/module.go` scaffold state

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema extension, domain/ports, persistence skeleton, and super-admin gate required by all stories.

**⚠️ CRITICAL**: No customer endpoint work until migration applies and `RequireSuperAdmin` exists.

- [x] T005 Create `serverBackendGo/db/migrations/000005_customers_extend.up.sql` with `email`, `accounttype`, `customerstatus`, `registrationtime`, `expirytime`, `devicelimit`, `deviceconfigurationid` per `data-model.md`
- [x] T006 [P] Add `serverBackendGo/db/migrations/000005_customers_extend.down.sql` reversing column adds
- [x] T007 [P] Create `serverBackendGo/internal/modules/customers/domain/customer.go` with `Customer`, `SearchRequest`, `Paginated` types per `data-model.md`
- [x] T008 [P] Create `serverBackendGo/internal/modules/customers/domain/transliterate.go` porting `CustomerDAO.transliterate` from Java
- [x] T009 Define `serverBackendGo/internal/modules/customers/port/customer_repository.go` with `Search`, `Count`, `GetByID`, `GetForEdit`, `Insert`, `Update`, `Delete`, `PrefixUsed`, duplicate lookup methods
- [x] T010 Define `serverBackendGo/internal/modules/customers/port/user_lookup.go` with `FindOrgAdmin`, `EnsureAuthToken`, `UpdateOrgAdminMainDetails`, `FindByLoginOrEmail` stubs
- [x] T011 Implement `serverBackendGo/internal/modules/customers/adapter/persistence/postgres/customer_repo.go` search/count SQL mirroring `CustomerMapper.xml` (`master = false`, filters, OFFSET/LIMIT)
- [x] T012 [P] Implement `serverBackendGo/internal/platform/auth/superadmin.go` with `RequireSuperAdmin(c *gin.Context) (*Principal, bool)` returning `error.permission.denied` envelope
- [x] T013 Verify migration applies: restart `make dev` and confirm extended `customers` columns exist in Postgres

**Checkpoint**: Repo compiles; super-admin helper available — proceed to user stories.

---

## Phase 3: User Story 1 — Search customer accounts (Priority: P1) 🎯 MVP

**Goal**: `POST /rest/private/customers/search` returns paginated `items` + `totalItemsCount` for super admin.

**Independent Test**: Super-admin JWT → POST search with `currentPage`/`pageSize` → non-empty or empty `items` with total count.

### Tests for User Story 1

- [x] T014 [P] [US1] Add `serverBackendGo/internal/modules/customers/application/service_test.go` stub for `Search` pagination and empty filter
- [x] T015 [P] [US1] Add `serverBackendGo/internal/modules/customers/adapter/http/handler_test.go` for `POST /search` 403 without super admin and 200 with super-admin principal

### Implementation for User Story 1

- [x] T016 [US1] Implement `Search` in `serverBackendGo/internal/modules/customers/application/service.go` calling repo `Search` + `Count`
- [x] T017 [US1] Create `serverBackendGo/internal/modules/customers/adapter/http/handler.go` with `Search` handler, `RequireSuperAdmin`, Headwind `response.OK` paginated shape
- [x] T018 [US1] Register `POST /search` in `handler.go` `Register` on `groups.Private.Group("/customers")`
- [x] T019 [US1] Wire `serverBackendGo/internal/modules/customers/module.go`: repo → service → handler; require `deps.DB`; remove scaffold-only route group

**Checkpoint**: `specs/004-complete-customers-module/quickstart.md` §3 search curl succeeds.

---

## Phase 4: User Story 2 — Impersonate org admin (Priority: P1)

**Goal**: `GET /rest/private/customers/impersonate/{id}` returns `LoginUserPayload` without password for super admin.

**Independent Test**: Super-admin JWT → impersonate customer with org admin → `data.authToken` present; React control panel impersonate works.

### Tests for User Story 2

- [x] T020 [P] [US2] Extend `service_test.go` for `Impersonate` no org admin (`error.notfound.customer.admin`) and password-reset token blocked
- [x] T021 [P] [US2] Extend `handler_test.go` for `GET /impersonate/:id` 403 non–super-admin and 200 with token in body

### Implementation for User Story 2

- [x] T022 [US2] Implement `serverBackendGo/internal/modules/customers/adapter/persistence/postgres/user_lookup.go` — `FindOrgAdmin` (role id 2), `EnsureAuthToken` using `internal/platform/crypto`
- [x] T023 [US2] Implement `Impersonate` in `serverBackendGo/internal/modules/customers/application/service.go` mapping user to login payload (no password field)
- [x] T024 [US2] Add `Impersonate` handler and register `GET /impersonate/:id` in `serverBackendGo/internal/modules/customers/adapter/http/handler.go`

**Checkpoint**: `quickstart.md` §4 impersonate curl succeeds; Control Panel impersonate hydrates session.

---

## Phase 5: User Story 3 — Create or update customer (Priority: P2)

**Goal**: `PUT /rest/private/customers` creates tenant + org admin + `adminCredentials` or updates customer and syncs org admin.

**Independent Test**: PUT without `id` → `adminCredentials` in data; PUT with `id` → OK; duplicate name → `error.duplicate.customer.name`.

### Tests for User Story 3

- [x] T025 [P] [US3] Extend `service_test.go` for `Save` duplicate name/email and successful create returning credentials string
- [x] T026 [P] [US3] Extend `handler_test.go` for `PUT /` create/update envelope and duplicate error status

### Implementation for User Story 3

- [x] T027 [US3] Extend `customer_repo.go` with `Insert`, `Update`, `GetByName`, `GetByEmail`, `NameExists`, `EmailExists` duplicate helpers
- [x] T028 [US3] Extend `user_lookup.go` with cross-tenant login/email conflict checks per `CustomerResource.updateCustomer`
- [x] T029 [US3] Implement `Save` in `application/service.go` — create path: insert customer (`filesdir` UUID), insert org admin (role 2), transliterated login, generated password, return `adminCredentials`; update path: sync org admin main details
- [x] T030 [US3] Document **Partial** parity in code comment: defer default devices/config copy until Phase 4/5 (no `devices` table yet)
- [x] T031 [US3] Add `Save` handler and register `PUT /` in `serverBackendGo/internal/modules/customers/adapter/http/handler.go`

**Checkpoint**: Manual PUT create smoke returns `adminCredentials`; duplicate name returns ERROR envelope.

---

## Phase 6: User Story 4 — Delete customer account (Priority: P2)

**Goal**: `DELETE /rest/private/customers/{id}` removes customer for super admin.

**Independent Test**: Create or pick customer id → DELETE → search no longer lists it (cascade users).

### Tests for User Story 4

- [x] T032 [P] [US4] Extend `service_test.go` for `Delete` success and permission denied for non–super-admin

### Implementation for User Story 4

- [x] T033 [US4] Implement `Delete` in `customer_repo.go` and `application/service.go`
- [x] T034 [US4] Add `Delete` handler and register `DELETE /:id` in `serverBackendGo/internal/modules/customers/adapter/http/handler.go`

**Checkpoint**: DELETE smoke returns OK for seeded non-master customer.

---

## Phase 7: User Story 5 — Load customer for edit (Priority: P2)

**Goal**: `GET /rest/private/customers/{id}/edit` returns full customer DTO for update forms.

**Independent Test**: GET edit for known id → OK with customer fields; unknown id → ERROR (not 500 stack).

### Tests for User Story 5

- [x] T035 [P] [US5] Extend `service_test.go` for `GetForEdit` found vs not found

### Implementation for User Story 5

- [x] T036 [US5] Implement `GetForEdit` in `customer_repo.go` and `application/service.go`
- [x] T037 [US5] Add `GetForEdit` handler and register `GET /:id/edit` in `serverBackendGo/internal/modules/customers/adapter/http/handler.go`

**Checkpoint**: GET edit returns customer JSON matching legacy field names (camelCase JSON tags).

---

## Phase 8: User Story 6 — Validate prefix availability (Priority: P3)

**Goal**: `GET /rest/private/customers/prefix/{prefix}/used` returns boolean.

**Independent Test**: GET with existing seed prefix → `true`; unused prefix → `false`.

### Implementation for User Story 6

- [x] T038 [US6] Implement `PrefixUsed` in `customer_repo.go` and `application/service.go`
- [x] T039 [US6] Add `PrefixUsed` handler and register `GET /prefix/:prefix/used` in `serverBackendGo/internal/modules/customers/adapter/http/handler.go`

**Checkpoint**: `quickstart.md` §5 prefix check passes.

---

## Phase 9: User Story 7 — Verifiable API and regression safety (Priority: P2)

**Goal**: Swagger Bearer on all customer routes; module tests pass; parity doc complete.

**Independent Test**: `go test ./internal/modules/customers/...` green; Swagger Authorize → search + impersonate; parity checklist 100% in-scope endpoints Done.

### Tests for User Story 7

- [x] T040 [P] [US7] Run `go test ./internal/modules/customers/...` and fix failures across `application/` and `adapter/http/`
- [x] T041 [P] [US7] Add Swagger `@Security BearerAuth` annotations on all handlers in `serverBackendGo/internal/modules/customers/adapter/http/handler.go`

### Implementation for User Story 7

- [x] T042 [US7] Run `cd serverBackendGo && make swagger` and verify `/private/customers/*` appears in Swagger UI
- [x] T043 [US7] Create `serverBackendGo/docs/parity/customers.md` — endpoint matrix (Done/Partial/N/A) per `contracts/customers-api.md`
- [x] T044 [US7] Mark create side effects **Partial** (devices/config) in parity doc per `research.md` R5

**Checkpoint**: Module test suite and Swagger smoke complete.

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Close Phase 3 roadmap and validate end-to-end with React.

- [x] T045 [P] Update `serverBackendGo/docs/MIGRATION.md` Phase 3 row — customers **done** alongside hints/settings/summary
- [x] T046 [P] Update `serverBackendGo/docs/NEXT_STEPS.md` — customers **منجز**; next focus **devices** / Phase 4
- [x] T047 Run full `specs/004-complete-customers-module/quickstart.md` validation (search, impersonate, prefix, tests)
- [x] T048 Run `cd serverBackendGo && go build ./...` and `go test ./...` for regression
- [x] T049 [P] Manual E2E: React Control Panel search + impersonate against Go-only backend on `:8080`
- [x] T050 Seed or document dev data: at least one non-master `customers` row + org admin (`userroleid=2`) for impersonate smoke in `serverBackendGo/db/seed` or migration notes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — **BLOCKS** all user stories
- **US1 Search (Phase 3)**: After Foundational — **MVP**
- **US2 Impersonate (Phase 4)**: After Foundational; integrates with US1 repo/service but independently testable
- **US3 Save (Phase 5)**: After Foundational; uses shared repo + user_lookup
- **US4 Delete (Phase 6)**: After Foundational; can follow US3 or parallel if repo methods done
- **US5 Edit (Phase 7)**: After Foundational; parallel with US4
- **US6 Prefix (Phase 8)**: After Foundational; smallest increment
- **US7 Swagger/parity (Phase 9)**: After US1–US6 handlers exist
- **Polish (Phase 10)**: After desired stories complete (minimum US1+US2+US7 for control panel)

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|-------|
| US1 Search | Phase 2 | MVP alone unblocks empty control panel list |
| US2 Impersonate | Phase 2, US1 module wired | Shares `module.go` / handler file |
| US3 Save | Phase 2, user_lookup | Partial device parity documented |
| US4 Delete | Phase 2 | Independent |
| US5 Edit | Phase 2 | Independent |
| US6 Prefix | Phase 2 | Independent |
| US7 Tests/Swagger | US1–US6 routes | Cross-cutting verification |

### Parallel Opportunities

- **Phase 1**: T002, T003 in parallel
- **Phase 2**: T006–T008, T012 in parallel after T005
- **US1**: T014, T015 in parallel
- **US2**: T020, T021 in parallel; T022 parallel with T023 after T010
- **US3–US5**: test tasks [P] within each phase
- **Phase 10**: T045, T046, T050 in parallel

### Parallel Example: User Story 1

```bash
# Tests in parallel:
T014 service_test.go Search stubs
T015 handler_test.go POST /search auth cases

# Then sequential:
T016 application Search → T017–T019 HTTP + module wire
```

### Parallel Example: User Story 2

```bash
T020 service_test Impersonate cases
T021 handler_test GET /impersonate/:id
T022 user_lookup.go  # parallel with tests if interface defined in T010
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: US1 Search
4. Complete Phase 4: US2 Impersonate
5. **STOP and VALIDATE**: React Control Panel + `quickstart.md` §3–§4
6. Add US7 Swagger + parity doc (subset of Phase 9) for reviewability

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. US1 + US2 → Control Panel works (primary Phase 3 closure for UX)
3. US3 Save → admin create/update API
4. US4 + US5 + US6 → full CRUD/prefix parity
5. US7 + Polish → Phase 3 sign-off in `MIGRATION.md`

### Suggested MVP Scope

**Minimum shippable**: Phase 1 + 2 + **US1** + **US2** + T043 parity doc (partial) + T045–T047 smoke.

Delivers React control panel without Java for search and impersonate.

---

## Notes

- Mailchimp subscribe on create: **out of scope** (no-op)
- Deprecated `GET /search` endpoints: **not implemented**
- Default devices on customer create: **Partial** until Phase 4 `devices` module
- All private routes require super admin (`Principal.SuperAdmin`)
- Commit after each phase checkpoint; branch `004-complete-customers-module`
