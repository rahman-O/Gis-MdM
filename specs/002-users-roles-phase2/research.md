# Research: Phase 2 — Users & Roles API Migration

**Date**: 2026-05-20

## 1. Module boundaries (users vs auth)

**Decision**: Implement admin/profile use cases in `internal/modules/users/` with its
own `port.UserRepository` and `adapter/persistence/postgres`, not in `auth/application`.

**Rationale**: Constitution Principle I requires one bounded context per module.
Auth keeps login/JWT/signup; users owns `/rest/private/users/*` lifecycle.

**Alternatives considered**:
- Extend `auth/port.UserRepository` with 10 more methods — couples auth to admin CRUD.
- Single `internal/repository/user` package — valid later refactor; deferred to avoid
  large cross-module move in this slice.

## 2. Permission model

**Decision**: Add `internal/platform/auth/permissions.go` that loads the current user's
role permissions (from seed: `userrolepermissions` + `permissions.name`) once per
request after `Principal` is set, and exposes `HasPermission(name string) bool`.

**Rationale**: Java `SecurityContext.hasPermission("settings")` gates list/CRUD.
Without this, Go returns 200 for unauthorized callers and breaks parity.

**Alternatives considered**:
- Check only `userRole.superadmin` flag — insufficient for org admins with `settings`.
- Full SecurityContext port — overkill; map permission names to slice on Principal.

## 3. Password change payload

**Decision**: Accept `oldPassword` and `newPassword` as MD5 uppercase hex in JSON
(same as React `profileService` / Java `User` fields), normalize via
`shared/crypto.NormalizeLoginPassword`, verify with `crypto.PasswordMatch`.

**Rationale**: Matches existing auth module and frontend `updateCurrentPassword`.

**Alternatives considered**: Raw password in API — only for Swagger manual tests via
documented MD5 hex in quickstart.

## 4. User aggregate query

**Decision**: Reuse the same multi-row aggregate pattern as
`auth/adapter/persistence/postgres/user_repo.go` (`userDataSelect` joins) inside
`users/adapter/persistence/postgres` for `GetUserDetails`, `ListUsers`.

**Rationale**: Java `getUserDetails` returns groups, configurations, role, permissions.
React `normalizeUser` expects `userRole`, `groups`, `configurations`.

**Alternatives considered**: Separate queries per relation — simpler SQL but more
round trips; rejected for list endpoint performance.

## 5. Roles persistence

**Decision**: New `roles` module with tables already in `000001_init` (`userroles`,
`permissions`, `userrolepermissions`). Implement insert/update/delete for roles and
junction rows mirroring `UserRoleDAO`.

**Rationale**: Seed data has roles 1–3; custom roles need writable CRUD for Roles UI.

**Alternatives considered**: Read-only roles — fails FR-008 and React `roleService`.

## 6. Schema gaps

**Decision**: Add `000004_users_roles_parity.up.sql` only if implementation reveals
missing columns (e.g. `userroles.description`). Initial `000001` schema covers core
tables; verify against Java Liquibase before adding columns.

**Rationale**: Constitution V — minimal migrations.

## 7. Swagger

**Decision**: Annotate handlers in `users/adapter/http` and `roles/adapter/http`;
regenerate with existing `make swagger` target; document Bearer + session cookie in
`quickstart.md`.

**Rationale**: Spec FR-010 / user story 4; project already ships Swagger UI.

**Alternatives considered**: Postman collection only — less discoverable for team.

## 8. Out of scope confirmation

**Decision**: Defer `GET /impersonate/{id}`, `/superadmin/*` to a future spec.

**Rationale**: Not used by React paths in `frontend/src/features/users/userService.ts`.

**Alternatives considered**: Include for parity completeness — increases scope and
session security risk in first Go slice.
