# Data Model: Phase 2 — Users & Roles

**Date**: 2026-05-20  
**Storage**: PostgreSQL (existing Headwind MDM schema + optional `000004` patch)

## Entity: User

| Field | Type | Rules |
|-------|------|-------|
| id | int64 | PK |
| login | string | Unique per tenant; max ~30 chars |
| email | string | Optional; unique among users when set |
| name | string | Display name |
| password | string | SHA1(MD5+salt) hash; never in API responses |
| customerId | int | FK → customers; all queries scoped |
| userRoleId | int | FK → userroles |
| allDevicesAvailable | bool | If true, groups access ignored |
| allConfigAvailable | bool | If true, config access ignored |
| authToken | string | Omitted in list/detail responses |
| passwordReset | bool | Admin-created reset flow |
| twoFactorSecret / twoFactorAccepted | optional | Cleared on admin password reset |

**Relations**:
- Many-to-many **groups** via `userdevicegroupsaccess`
- Many-to-many **configurations** via `userconfigurationaccess`
- **UserRole** via `userroleid`

**Validation**:
- Create: `newPassword` required (MD5 hex)
- Update password: `oldPassword` must match stored hash
- Duplicate login/email → `error.duplicate.login` / `error.duplicate.email`

## Entity: UserRole

| Field | Type | Rules |
|-------|------|-------|
| id | int | PK |
| name | string | Unique |
| description | string | Optional (if column exists) |
| superadmin | bool | Seed role 1 |
| permissions | Permission[] | Via `userrolepermissions` |

## Entity: Permission

| Field | Type | Rules |
|-------|------|-------|
| id | int | PK |
| name | string | e.g. `settings`, `superadmin` |
| superadmin | bool | Flag on permission row |

## Entity: Group / Configuration (lookup)

Referenced only as `{ id, name }` on user payloads for React; tables `groups`,
`configurations` exist in `000001_init`.

## Authorization rules (application layer)

| Action | Allowed when |
|--------|----------------|
| GET /users/current, PUT /details, PUT /current | Authenticated self |
| GET /users/all, PUT /users, DELETE /other/:id | `settings` permission OR super admin OR org admin |
| GET /users/roles | Authenticated (legacy: any logged-in user) |
| GET/PUT/DELETE /roles/* | `userRoleDAO.hasAccess()` equivalent: super admin or settings |

## State transitions

**User lifecycle (admin)**:
1. Create → insert user + hash password + optional group/config links
2. Update → patch details; optional password rotation
3. Delete → remove user row (cascade per FK rules)

**Role lifecycle**:
1. Create role → insert `userroles` + `userrolepermissions`
2. Update → replace permission links
3. Delete → remove role if not referenced (or block if users reference — match Java)

## API DTO shapes (JSON)

Responses use Headwind envelope. User detail mirrors Java `User` + `UserView` fields
needed by React (`userRole`, `groups`, `configurations`, `editable` on list items).

List users: set `editable: false` for current user id; strip `password`, `authToken`.
