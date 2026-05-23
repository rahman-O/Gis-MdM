---
description: "Task list for Phase 2 Users & Roles API migration"
---

# Tasks: Phase 2 — Users & Roles API Migration

**Input**: `specs/002-users-roles-phase2/` (plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md)

**Prerequisites**: Auth Phase 1 complete; `serverBackendGo` builds; Postgres via `./scripts/db-up.sh`

**Tests**: Included per FR-011 and spec User Story 4 (application + HTTP handler tests).

**Organization**: Tasks grouped by user story for independent delivery and verification.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Parallelizable (different files, no dependency on incomplete tasks in same phase)
- **[USn]**: User story from spec.md

## Path Conventions

- Go backend: `serverBackendGo/internal/modules/users/`, `serverBackendGo/internal/modules/roles/`
- Platform: `serverBackendGo/internal/platform/auth/`
- Migrations: `serverBackendGo/db/migrations/`
- Docs: `serverBackendGo/docs/parity/`, `specs/002-users-roles-phase2/quickstart.md`

---

## Phase 1: Setup

**Purpose**: Confirm feature branch context and migration baseline before code changes.

- [x] T001 Verify feature context in `specs/002-users-roles-phase2/spec.md` and `plan.md` against `serverBackendGo/docs/MIGRATION.md` Phase 2 scope
- [x] T002 [P] Review Java references `backend/server/src/main/java/com/hmdm/rest/resource/UserResource.java` and `UserRoleResource.java` for endpoint parity checklist
- [x] T003 [P] Review React consumers `frontend/src/features/profile/profileService.ts` and `frontend/src/features/users/userService.ts` for required JSON fields
- [x] T004 Run baseline `cd serverBackendGo && go build ./...` and record current `internal/modules/users` / `roles` scaffold state in task notes
- [x] T005 [P] Compare `000001_init.up.sql` to Java Liquibase for missing user/role columns; draft `serverBackendGo/db/migrations/000004_users_roles_parity.up.sql` only if gaps found

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Permission resolution required before admin user/role mutations (matches Java `SecurityContext`).

**⚠️ CRITICAL**: User Story 2 and 3 admin endpoints MUST NOT ship before this phase completes.

- [x] T006 Extend `serverBackendGo/internal/platform/auth/context.go` `Principal` with permission names and `SuperAdmin` flag fields
- [x] T007 Implement `serverBackendGo/internal/platform/auth/permissions.go` with `LoadPermissions(ctx, userID)` and `HasPermission(name string) bool`
- [x] T008 Wire permission load in `serverBackendGo/internal/platform/httpx/middleware/jwt.go` and `session_auth.go` after principal resolution
- [x] T009 [P] Add `serverBackendGo/internal/platform/auth/permissions_test.go` for seed admin (`settings` / superadmin) vs user without settings
- [x] T010 [P] Add org-admin helper in `serverBackendGo/internal/platform/auth/permissions.go` matching Java `UserDAO.isOrgAdmin` (role id 2 heuristic from seed)

**Checkpoint**: Authenticated requests expose `HasPermission("settings")` — proceed to user stories.

---

## Phase 3: User Story 1 — Profile & session refresh (Priority: P1) 🎯 MVP

**Goal**: Full `GET /current`, profile update, and password change for React Profile and session refresh.

**Independent Test**: Login → `GET /rest/private/users/current` → `PUT /details` → `PUT /current` with valid passwords → re-login succeeds.

### Tests for User Story 1

- [x] T011 [P] [US1] Add `serverBackendGo/internal/modules/users/application/profile_test.go` stubs for `UpdateProfile` duplicate email and `ChangePassword` wrong old password
- [x] T012 [P] [US1] Add `serverBackendGo/internal/modules/users/adapter/http/handler_profile_test.go` for `GET /current` and `PUT /details` with session cookie

### Implementation for User Story 1

- [x] T013 [P] [US1] Create `serverBackendGo/internal/modules/users/domain/user.go` and `user_detail.go` DTOs per `data-model.md` and React `BackendUser` shape
- [x] T014 [P] [US1] Define `serverBackendGo/internal/modules/users/port/repository.go` with `GetByID`, `UpdateMainDetails`, `UpdatePassword` for self-service
- [x] T015 [US1] Implement `serverBackendGo/internal/modules/users/adapter/persistence/postgres/user_repo.go` with `userDataSelect` aggregate query (per `research.md`)
- [x] T016 [US1] Implement `serverBackendGo/internal/modules/users/application/profile.go` — `GetCurrentUser`, `UpdateProfile`, `ChangePassword` using `shared/crypto`
- [x] T017 [US1] Refactor `serverBackendGo/internal/modules/users/adapter/http/handler.go` — enhance `Current`; add `UpdateDetails`, `ChangePassword` handlers
- [x] T018 [US1] Register routes in `serverBackendGo/internal/modules/users/adapter/http/routes.go` for `PUT /details` and `PUT /current`
- [x] T019 [US1] Update `serverBackendGo/internal/modules/users/module.go` to use users `application.Service` instead of `authapp.UsersService` for profile paths
- [x] T020 [US1] Deprecate or delegate `serverBackendGo/internal/modules/auth/application/users.go` `CurrentUser` to users module to avoid duplicate logic

**Checkpoint**: Profile smoke from `specs/002-users-roles-phase2/quickstart.md` §2 Users (current, details, password) passes.

---

## Phase 4: User Story 2 — Tenant user management (Priority: P2)

**Goal**: List, create, update, and delete tenant users for React Users screen.

**Independent Test**: Admin session → `GET /all` → `PUT /` create → `PUT /` update → `DELETE /other/:id`.

### Tests for User Story 2

- [x] T021 [P] [US2] Add `serverBackendGo/internal/modules/users/application/admin_test.go` for `ListUsers` editable flag, `UpsertUser` duplicate login, `DeleteUser` permission denied
- [x] T022 [P] [US2] Extend `serverBackendGo/internal/modules/users/adapter/http/handler_test.go` for `GET /all` success and forbidden without settings permission

### Implementation for User Story 2

- [x] T023 [US2] Extend `serverBackendGo/internal/modules/users/port/repository.go` with `ListByCustomer`, `Upsert`, `Delete`, group/config link updates
- [x] T024 [US2] Implement list/upsert/delete SQL in `serverBackendGo/internal/modules/users/adapter/persistence/postgres/user_repo.go` with `customerId` scoping
- [x] T025 [US2] Implement `serverBackendGo/internal/modules/users/application/admin.go` — `ListUsers`, `UpsertUser`, `DeleteUser` with authorization checks
- [x] T026 [US2] Add handlers `ListAll`, `Upsert`, `DeleteOther` in `serverBackendGo/internal/modules/users/adapter/http/handler.go`
- [x] T027 [US2] Register `GET /all`, `PUT /`, `DELETE /other/:id` in `serverBackendGo/internal/modules/users/adapter/http/routes.go`
- [x] T028 [US2] Strip `password` and `authToken` from list/detail responses; set `editable: false` for current user id per `contracts/users-api.md`

**Checkpoint**: React Users page loads list and CRUD without 403/404 when Go-only backend is used.

---

## Phase 5: User Story 3 — Role catalog (Priority: P2)

**Goal**: Roles dropdown and Roles admin screen (`/private/users/roles` + `/private/roles/*`).

**Independent Test**: `GET /users/roles` → `GET /roles/permissions` → `GET /roles/all` → `PUT /roles` → `DELETE /roles/:id`.

### Tests for User Story 3

- [x] T029 [P] [US3] Add `serverBackendGo/internal/modules/roles/application/service_test.go` for duplicate role name and permission denied
- [x] T030 [P] [US3] Add `serverBackendGo/internal/modules/roles/adapter/http/handler_test.go` for `GET /permissions` with admin session

### Implementation for User Story 3

- [x] T031 [P] [US3] Create `serverBackendGo/internal/modules/roles/domain/role.go` and `permission.go` per `data-model.md`
- [x] T032 [P] [US3] Define `serverBackendGo/internal/modules/roles/port/repository.go` for permissions, roles CRUD
- [x] T033 [US3] Implement `serverBackendGo/internal/modules/roles/adapter/persistence/postgres/role_repo.go`
- [x] T034 [US3] Implement `serverBackendGo/internal/modules/roles/application/service.go` — `ListPermissions`, `ListRoles`, `UpsertRole`, `DeleteRole`
- [x] T035 [US3] Create `serverBackendGo/internal/modules/roles/adapter/http/handler.go` and `routes.go` per `contracts/roles-api.md`
- [x] T036 [US3] Replace scaffold in `serverBackendGo/internal/modules/roles/module.go` with full module wiring
- [x] T037 [US3] Align `serverBackendGo/internal/modules/users/application/roles.go` `ListRoles` with roles repo or shared query for consistent id/name
- [x] T038 [US3] Ensure `GET /rest/private/users/roles` returns same role ids as `GET /rest/private/roles/all` for seed data

**Checkpoint**: Roles UI and Users role dropdown work against Go backend only.

---

## Phase 6: User Story 4 — Verifiable API contract (Priority: P1)

**Goal**: Swagger coverage, parity docs, and full module test pass for migration sign-off.

**Independent Test**: `make swagger` → UI lists Phase 2 routes → `go test ./internal/modules/users/... ./internal/modules/roles/...` green → `quickstart.md` full run.

### Implementation for User Story 4

- [x] T039 [P] [US4] Add Swagger `// @Summary` comments on all new handlers in `serverBackendGo/internal/modules/users/adapter/http/handler.go`
- [x] T040 [P] [US4] Add Swagger comments on `serverBackendGo/internal/modules/roles/adapter/http/handler.go`
- [x] T041 [US4] Run `make swagger` in `serverBackendGo/` and commit regenerated `internal/platform/httpx/swagger/docs.go` if changed
- [x] T042 [P] [US4] Create `serverBackendGo/docs/parity/users.md` marking all in-scope endpoints Done with Java class references
- [x] T043 [P] [US4] Create `serverBackendGo/docs/parity/roles.md` marking all in-scope endpoints Done
- [x] T044 [US4] Update `serverBackendGo/docs/NEXT_STEPS.md` Phase 2 rows (#4 users complete, #5 roles) to **منجز**
- [x] T045 [US4] Run full smoke in `specs/002-users-roles-phase2/quickstart.md` and fix any failures
- [x] T046 [US4] Run `cd serverBackendGo && go test ./internal/modules/users/... ./internal/modules/roles/... ./internal/platform/auth/... -v` and `go build ./...`

**Checkpoint**: Phase 2 migration sign-off criteria SC-001 through SC-005 from spec.md satisfied.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Hardening and documentation consistency across modules.

- [x] T047 [P] Remove dead scaffold files `serverBackendGo/internal/modules/users/domain/doc.go` and `roles/domain/doc.go` or replace with package docs
- [x] T048 Verify no upward imports from `users`/`roles` domain into `adapter` violated (constitution II) via `go build` and import review
- [x] T049 [P] Ensure all Phase 2 error messages use legacy keys (`error.permission.denied`, `error.duplicate.email`, etc.) in `serverBackendGo/internal/platform/httpx/response/envelope.go` usages
- [ ] T050 Manual E2E: `frontend` login → Profile → Users → Roles with Vite proxy to `:8080`; note defects as follow-up issues if any

---

## Dependencies & Execution Order

### Phase Dependencies

```text
Phase 1 Setup
    ↓
Phase 2 Foundational (permissions) — BLOCKS US2, US3 admin writes
    ↓
Phase 3 US1 (Profile) — can start after Phase 2
    ↓
Phase 4 US2 (Users admin) — requires Phase 2
    ↓
Phase 5 US3 (Roles) — requires Phase 2; parallel with US4 implementation after US2 repo patterns exist
    ↓
Phase 6 US4 (Swagger/docs/tests sign-off)
    ↓
Phase 7 Polish
```

### User Story Dependencies

| Story | Depends on | Can parallel with |
|-------|------------|-------------------|
| US1 | Phase 2 | — |
| US2 | Phase 2, US1 repo patterns (T015) | US3 after T033 started |
| US3 | Phase 2 | US2 (different module) after foundational |
| US4 | US1–US3 handlers exist | Polish tasks [P] |

### Within Each User Story

- Tests written alongside or immediately after application layer (T011–T012 before T016–T018, etc.)
- domain + port before application
- application before HTTP handlers
- routes registered last in each story

---

## Parallel Example: User Story 3

```bash
# After Phase 2 complete, launch in parallel:
T031 domain/role.go
T032 port/repository.go
T039 Swagger users handlers  # US4 partial early
```

---

## Parallel Example: User Story 1 + Foundational tests

```bash
# After T006–T008:
T009 permissions_test.go
T011 profile_test.go
T013 domain/user.go
```

---

## Implementation Strategy

### MVP First (User Story 1 + Foundational)

1. Complete Phase 1 + Phase 2 (permissions).
2. Complete Phase 3 (US1): Profile + `GET /current` full payload.
3. **STOP and VALIDATE**: `quickstart.md` profile section + React Profile page.
4. Demo without Java WAR for auth + profile path.

### Incremental Delivery

1. US1 → US2 (Users admin) → US3 (Roles) → US4 (Swagger/parity/sign-off).
2. Each story independently testable per spec Independent Test criteria.
3. Do not implement `impersonate` or `superadmin/*` (out of scope).

### Suggested MVP Scope

- **Minimum**: Phase 2 + Phase 3 (US1) — 20 tasks through T020.
- **Full Phase 2**: All phases through T046.

---

## Task Summary

| Phase | Task IDs | Count |
|-------|----------|-------|
| Setup | T001–T005 | 5 |
| Foundational | T006–T010 | 5 |
| US1 Profile | T011–T020 | 10 |
| US2 Users admin | T021–T028 | 8 |
| US3 Roles | T029–T038 | 10 |
| US4 Verify contract | T039–T046 | 8 |
| Polish | T047–T050 | 4 |
| **Total** | T001–T050 | **50** |

**Parallel opportunities**: 24 tasks marked `[P]`

**Independent test criteria**: See each phase **Checkpoint** and spec.md User Story sections.
