# Implementation Plan: Phase 2 — Users & Roles API Migration

**Branch**: `002-users-roles-phase2` | **Date**: 2026-05-20 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/002-users-roles-phase2/spec.md`

## Summary

Complete **Phase 2** of the Java→Go migration by implementing full parity for
`UserResource` and `UserRoleResource` in `serverBackendGo`: profile (`current`,
`details`, password change), tenant user CRUD, and roles admin APIs. Reuse existing
auth password hashing and session/JWT middleware. Add `platform/auth` permission
resolution, dedicated `users` and `roles` clean-architecture modules, SQL
migrations where Liquibase fields are missing, Swagger annotations for manual QA,
and `application/` + HTTP tests per constitution v1.0.0.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo/go.mod`)

**Primary Dependencies**: Gin, `lib/pq`, `gin-contrib/sessions`, `golang-jwt/jwt/v5`,
`swaggo/gin-swagger` (existing in project)

**Storage**: PostgreSQL (`DATABASE_URL`); migrations in `serverBackendGo/db/migrations/`

**Testing**: `go test ./internal/modules/users/... ./internal/modules/roles/...`;
handler tests with `httptest` + session cookie; optional table-driven permission tests

**Target Platform**: Linux/macOS dev server (`:8080`); Docker Postgres via `scripts/db-up.sh`

**Project Type**: Web service (Headwind MDM REST API) + React consumer (`frontend/`)

**Performance Goals**: Admin list endpoints &lt; 500ms p95 on seed DB (&lt; 100 users);
no N+1 beyond Java parity (aggregated user query already joins groups/configs)

**Constraints**: REST path and envelope parity; tenant scoping by `customerId`;
no `impersonate` / `superadmin/*` in this slice

**Scale/Scope**: ~10 endpoints users + ~4 endpoints roles; 2 modules; 1 migration tranche

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*
*Reference: `.specify/memory/constitution.md` (Gis-MdM v1.0.0)*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | `users` + `roles` modules; Phase 2 in MIGRATION.md |
| **II. Layered Clean** | ✅ | domain/port/application/adapter per module |
| **III. API Parity** | ✅ | contracts + `docs/parity/users.md`, `roles.md` |
| **IV. Testable Delivery** | ✅ | quickstart smoke + `go test` + Swagger |
| **V. Simplicity** | ✅ | Reuse `shared/crypto`; extend auth repo only via shared queries if needed |
| **VI. Security** | ✅ | Permission checker + customer scoping on all mutations |
| **VII. Observability** | ✅ | Legacy `message` keys; structured slog on errors |

**Post-design**: All gates remain ✅. No Complexity Tracking entries required.

## Project Structure

### Documentation (this feature)

```text
specs/002-users-roles-phase2/
├── plan.md              # This file
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1 smoke + Swagger
├── contracts/
│   ├── users-api.md
│   └── roles-api.md
└── tasks.md             # (/speckit-tasks — not created by plan)
```

### Source Code (repository root)

```text
serverBackendGo/
├── db/migrations/
│   └── 000004_users_roles_parity.up.sql   # only if columns/tables missing vs Java
├── docs/parity/
│   ├── users.md                           # update all endpoints Done
│   └── roles.md                           # new
├── internal/platform/auth/
│   ├── context.go                           # extend Principal with permissions
│   └── permissions.go                     # NEW: HasPermission("settings"), IsSuperAdmin
├── internal/modules/users/
│   ├── module.go
│   ├── domain/                              # User, UserPayload, RoleRef
│   ├── port/                                # UserRepository (admin + profile)
│   ├── application/                         # profile, list, upsert, delete
│   └── adapter/
│       ├── http/                            # handlers + swagger comments
│       └── persistence/postgres/            # user_repo.go (admin queries)
├── internal/modules/roles/
│   ├── module.go
│   ├── domain/
│   ├── port/
│   ├── application/
│   └── adapter/http/ + adapter/persistence/postgres/
└── internal/modules/auth/
    └── adapter/persistence/postgres/        # keep auth UserRepository; users may delegate FindByID
```

**Structure Decision**: Extend existing `users` scaffold into a full module; replace
`roles` scaffold. Avoid duplicating `userDataSelect` in three places: extract shared
SQL fragment to `internal/modules/users/adapter/persistence/postgres/queries.go` OR
have `users` repository call the same query builder used by auth (single package
import from `users/adapter/persistence` only — auth imports users persistence is
forbidden upward; preferred: **move** shared user read SQL to `internal/repository/user`
only if both need it — YAGNI: **copy once** into users repo, auth keeps its repo until
a later refactor).

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Platform permissions (blocking)

1. After auth middleware resolves `Principal`, load role + permission names from DB.
2. Store on `Principal` or `context` for `HasPermission("settings")`, `IsSuperAdmin()`.
3. Unit-test permission matrix (super admin role id=1 from seed).

### Phase B — Users module

| Endpoint | Handler | Application use case |
|----------|---------|----------------------|
| `GET /current` | enhance | `GetCurrentUser` — full `User` JSON (reuse aggregate query) |
| `PUT /details` | new | `UpdateProfile` |
| `PUT /current` | new | `ChangePassword` (old/new MD5 fields per Java) |
| `GET /all` | new | `ListUsers` (+ optional `?filter=`) |
| `PUT /` | new | `UpsertUser` create/update |
| `DELETE /other/:id` | new | `DeleteUser` |
| `GET /roles` | move/enhance | keep `ListRoles` in users or delegate to roles port |

Authorization: `settings` permission OR super admin OR org-admin heuristic matching
`UserDAO.isOrgAdmin` (role id 2 for org admin in seed).

### Phase C — Roles module

| Endpoint | Use case |
|----------|----------|
| `GET /permissions` | `ListPermissions` |
| `GET /all` | `ListRoles` with permission ids |
| `PUT /` | `UpsertRole` |
| `DELETE /:id` | `DeleteRole` |

`hasAccess()` in Java → super admin OR `settings` permission.

### Phase D — Swagger & docs

1. Add `// @Summary` etc. on each new handler (follow `auth/adapter/http`).
2. Run `make swagger` from `serverBackendGo/`.
3. Fill `docs/parity/users.md`, `docs/parity/roles.md`.
4. Update `docs/NEXT_STEPS.md` Phase 2 rows to **منجز** after implementation.

### Phase E — Tests

- `application/*_test.go` with stub repos: password change, duplicate email, permission denied.
- `adapter/http/handler_test.go`: login session → `GET /all` OK for admin.
- Document commands in `quickstart.md`.

## Complexity Tracking

> No constitution violations requiring justification.
