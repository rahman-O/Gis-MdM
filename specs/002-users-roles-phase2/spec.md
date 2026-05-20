# Feature Specification: Phase 2 — Users & Roles API Migration

**Feature Branch**: `002-users-roles-phase2`

**Created**: 2026-05-20

**Status**: Draft

**Input**: Complete Phase 2 from `serverBackendGo/docs/MIGRATION.md` and
`serverBackendGo/docs/NEXT_STEPS.md` after auth: migrate Java `UserResource` and
`UserRoleResource` to Go with parity, automated tests, and interactive API
documentation for hands-on verification.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Profile & session refresh (Priority: P1)

An authenticated administrator opens **Profile** or relies on post-login session
data. They view their full account (name, email, role, permissions context),
update profile fields, or change their password without using the legacy Java
server.

**Why this priority**: `GET /private/users/current` exists minimally today; Profile
and `authService.refreshSessionFromCurrentUser` need a complete, stable user
payload matching the legacy app.

**Independent Test**: Log in → open Profile → view details → update name/email →
change password with correct old password → session still valid.

**Acceptance Scenarios**:

1. **Given** a logged-in user, **When** they request current user details,
   **Then** the response includes id, login, name, email, role, and flags needed
   by the React shell (no password or auth token in the payload).
2. **Given** a logged-in user, **When** they update name and email via profile,
   **Then** changes persist and duplicate emails for another account are rejected
   with the same error semantics as the Java API.
3. **Given** a logged-in user, **When** they change password with a valid old
   password and non-empty new password, **Then** the password updates and
   subsequent login uses the new credentials.
4. **Given** a logged-in user, **When** they submit a wrong old password,
   **Then** the operation fails with `error.password.wrong` and no password change.

---

### User Story 2 — Tenant user management (Priority: P2)

An organization admin or super admin opens **Users** in the React app, lists
accounts, creates a new user, edits roles/groups/config access, and deletes another
user (not themselves).

**Why this priority**: Core MDM administration; React calls `/private/users/all`,
`PUT /private/users`, and `DELETE /private/users/other/:id`.

**Independent Test**: As admin with `settings` permission → list users → create
user with password → edit user → delete user.

**Acceptance Scenarios**:

1. **Given** a user with permission to manage settings, **When** they list all
   users, **Then** they receive all tenant users without passwords or auth tokens,
   and the current user is marked non-editable for self-delete semantics.
2. **Given** an org admin, **When** they create a user with unique login, email,
   role, and password, **Then** the account is created under the same customer
   tenant.
3. **Given** an org admin, **When** they update an existing user (role, groups,
   configurations, optional password), **Then** changes persist and duplicate
   login/email are rejected.
4. **Given** an org admin, **When** they delete another user's id,
   **Then** the account is removed and the API returns success.
5. **Given** a user without settings permission, **When** they call list/create/delete,
   **Then** access is denied consistent with Java (`error.permission.denied`).

---

### User Story 3 — Role catalog for Users & Roles screens (Priority: P2)

An administrator assigns roles when editing users or manages roles on the **Roles**
screen (list permissions, list roles, create/update/delete custom roles).

**Why this priority**: React uses `/private/users/roles` (partially implemented)
and `/private/roles/*` for the dedicated roles UI.

**Independent Test**: Load roles dropdown on Users page → open Roles admin → list
permissions → create role → delete role.

**Acceptance Scenarios**:

1. **Given** an authenticated user, **When** they request the role list for
   dropdowns, **Then** all assignable roles for the tenant are returned (id, name).
2. **Given** a user with role-management access, **When** they list permissions
   and all roles, **Then** data matches legacy shape for the Roles UI.
3. **Given** a user with role-management access, **When** they create or update a
   role by name, **Then** duplicate role names are rejected.
4. **Given** a user with role-management access, **When** they delete a role by id,
   **Then** the role is removed if allowed by legacy rules.

---

### User Story 4 — Verifiable API contract (Priority: P1)

A developer or QA validates Phase 2 without the Java WAR by using documented
interactive API docs and automated checks against a running database.

**Why this priority**: User explicitly requires Swagger and testing for real
migration confidence.

**Independent Test**: Start Go server → open API docs → execute sample requests for
each Phase 2 endpoint with admin session/JWT → run automated test suite for users
and roles modules.

**Acceptance Scenarios**:

1. **Given** the Go server with Swagger enabled, **When** a reviewer opens the
   docs UI, **Then** all Phase 2 endpoints are listed with request/response shapes.
2. **Given** documented smoke steps, **When** executed against Postgres with seed
   admin user, **Then** each P1–P2 endpoint returns the same status and envelope
   (`OK` / `ERROR`) behavior as Java for happy and primary error paths.
3. **Given** the module test suite, **When** run in CI or locally,
   **Then** permission checks, password validation, and duplicate detection are
   covered without regressions.

---

### Edge Cases

- Empty or duplicate email on profile update → `error.duplicate.email` or clear
  validation per Java.
- Create user without password → `error.password.empty`.
- Non-admin attempts user list or role CRUD → permission denied.
- Update password for another user's id via `/current` → denied.
- Filter query on `GET /all` (optional) returns filtered subset when supported.
- Super-admin-only endpoints (`impersonate`, `superadmin/*`) remain out of scope
  unless explicitly added in a follow-up slice.

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Phase**: 2 per `serverBackendGo/docs/MIGRATION.md` (`users`, `roles`).
- **Modules**: `internal/modules/users/` (extend beyond `current` + `roles` list),
  `internal/modules/roles/` (new implementation replacing scaffold).
- **Java references**:
  - `backend/server/.../UserResource.java`
  - `backend/server/.../UserRoleResource.java`
  - DAOs: `UserDAO`, `UserRoleDAO` in `backend/common/.../persistence/`
- **REST prefixes** (unchanged):
  - `/rest/private/users/*`
  - `/rest/private/roles/*`
- **Parity docs**: `serverBackendGo/docs/parity/users.md`,
  `serverBackendGo/docs/parity/roles.md` (create/update).
- **Layers**: User/Role aggregates in `domain/`; repositories in `port/` +
  `adapter/persistence/postgres/`; use cases in `application/`; HTTP in
  `adapter/http/`.
- **Dependencies**: Auth module (session/JWT, principal, password hashing) MUST
  remain the single source for credential rules.

### Functional Requirements

- **FR-001**: System MUST expose `GET /rest/private/users/current` returning full
  user details for the authenticated principal (password and auth token omitted).
- **FR-002**: System MUST expose `PUT /rest/private/users/details` for the
  authenticated user to update own name and email with uniqueness rules aligned
  with Java.
- **FR-003**: System MUST expose `PUT /rest/private/users/current` for password
  change with old/new password validation and legacy error messages.
- **FR-004**: System MUST expose `GET /rest/private/users/all` for authorized
  admins with optional filter, excluding secrets from listed users.
- **FR-005**: System MUST expose `PUT /rest/private/users` to create and update
  tenant users including role, group, and configuration assignments per payload.
- **FR-006**: System MUST expose `DELETE /rest/private/users/other/{id}` for
  authorized admins to delete another user.
- **FR-007**: System MUST expose `GET /rest/private/users/roles` listing assignable
  roles (id, name, and fields required by React).
- **FR-008**: System MUST expose `GET /rest/private/roles/permissions`,
  `GET /rest/private/roles/all`, `PUT /rest/private/roles`, and
  `DELETE /rest/private/roles/{id}` with Java-equivalent authorization and errors.
- **FR-009**: All responses MUST use the Headwind envelope
  (`status`, `message`, `data`) and HTTP status codes matching legacy behavior for
  the covered endpoints.
- **FR-010**: System MUST provide interactive API documentation covering every
  Phase 2 endpoint with enough detail to execute authenticated test calls.
- **FR-011**: System MUST include automated tests for application-layer rules
  (permissions, password match, duplicates) and handler smoke tests where
  practical.

### Key Entities

- **User**: Tenant account (login, name, email, role, customer, device/config
  access flags, group and configuration assignments).
- **UserRole**: Named role with permissions used for authorization and UI assignment.
- **Permission**: Named capability (e.g. `settings`, `superadmin`) linked to roles.
- **Group / Configuration access**: Optional assignments restricting device or
  config visibility per user.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of React-used Phase 2 endpoints (users list/CRUD/profile/password
  and roles CRUD/list) return successful responses for the documented admin smoke
  path when the Go server replaces Java.
- **SC-002**: Profile and Users screens load without console network errors for
  Phase 2 calls when pointed at the Go backend only.
- **SC-003**: Interactive API documentation lists at least every endpoint in FR-001
  through FR-008 and allows a reviewer to complete the smoke script in under
  30 minutes.
- **SC-004**: Automated tests for users and roles modules pass in a clean checkout
  (`go test` for those packages) with no manual steps.
- **SC-005**: Parity documents for `users` and `roles` are published and marked Done
  for all in-scope endpoints; `NEXT_STEPS.md` Phase 2 rows updated.

## Assumptions

- Phase 1 auth (session, JWT, password hashing) remains implemented and is not
  re-specified here.
- PostgreSQL schema is extended via migrations only as needed for user/role fields
  (groups, configurations, role permissions) consistent with legacy Liquibase.
- Permission checks mirror Java `SecurityContext` / `settings` permission for
  admin operations until a fuller RBAC port exists.
- Super-admin impersonation and `superadmin/*` user endpoints are deferred to a
  later slice to keep Phase 2 deliverable.
- `GET /private/users/roles` already exists in Go and may be extended, not
  replaced, if behavior already matches React expectations.
