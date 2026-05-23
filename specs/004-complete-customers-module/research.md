# Research: Phase 3 — Customers Module

**Date**: 2026-05-20

## R1 — Authorization model

**Decision**: All customer endpoints require authenticated principal with `SuperAdmin == true`
(mirror `CustomerDAO` / `CustomerResource.impersonate` checks). Add
`platformauth.RequireSuperAdmin(c *gin.Context) bool` helper used by handlers; return
`ERROR` + `error.permission.denied` on failure.

**Rationale**: Java throws `SecurityException.onAdminDataAccessViolation` for non–super-admin;
React control panel is super-admin tooling only.

**Alternatives considered**: `settings` permission — not used in `CustomerResource`.

## R2 — Paginated search response shape

**Decision**: Return Headwind `PaginatedData` JSON: `{ "items": [Customer...], "totalItemsCount": N }`.

**Rationale**: Java `PaginatedData` serializes `items` + `totalItemsCount`; React
`unwrapCustomerRows` reads `items` first.

**SQL**: Port `CustomerMapper.xml` `searchCustomers` + `countAllCustomers` (filter
`master = false`, ILIKE on name/description, optional accountType/customerStatus, OFFSET/LIMIT).

## R3 — Schema gap vs Java Liquibase

**Decision**: Add migration `000005_customers_extend.up.sql` with `ALTER TABLE ... ADD COLUMN IF NOT EXISTS`
for columns used by search/edit: `email`, `accounttype`, `customerstatus`, `registrationtime`,
`expirytime`, `devicelimit`, `deviceconfigurationid` (and any missing from minimal `000001_init`).

**Rationale**: `000001_init` customers table is auth-minimal; search filters reference
`accountType` / `customerStatus` in MyBatis.

**Alternatives considered**: Rewrite `000001` — rejected (breaks applied dev DBs).

## R4 — Impersonation semantics (Go + React)

**Decision**:

1. Load org admin (`userroleid = 2`, `OrgAdminRoleID`) for target `customerId`.
2. If `authtoken` empty, generate token and persist (mirror `UnsecureDAO.setUserNewPasswordUnsecure` path using existing `crypto.GenerateAuthToken()`).
3. Reject if `passwordresettoken` is non-empty (legacy security rule).
4. Return `LoginUserPayload`-compatible JSON (`authToken`, `superAdmin`, `singleCustomer`, `userRole`, no password).
5. For cookie sessions: optionally invalidate prior session and store new principal — **JWT path is primary** for React `hydrateSessionAfterImpersonation`.

**Rationale**: Control panel uses bearer token after impersonate; Java session swap is equivalent outcome.

**Alternatives considered**: Reuse deferred `users/impersonate` — out of scope per `parity/users.md`.

## R5 — Customer create side effects (devices/config)

**Decision**: **Defer** default device creation (`hmdm-001`…), configuration copy, design-settings copy,
and plugin `CustomerCreatedEvent` until Phase 4/5 when `devices` / `configurations` tables exist in
Go migrations. Implement create as: customer row + unique `filesdir` + org admin user +
`adminCredentials` string (`login/password`).

**Rationale**: `serverBackendGo` has no `devices` table yet; inserting devices would fail or require
premature Phase 4 scope.

**Parity doc**: Mark create side effects as **Partial** until devices module lands; search + impersonate **Done**.

**Alternatives considered**: Add minimal `devices` table in `000005` — rejected (violates module-first phases).

## R6 — Mailchimp on create

**Decision**: No-op stub; customer create succeeds without external API (per spec assumption).

## R7 — Module boundaries

**Decision**: `customers` module owns `customer_repo.go` and `user_lookup.go` (org admin queries)
in `adapter/persistence/postgres`. Does not import `users` or `auth` application packages.

**Rationale**: Constitution II — explicit SQL at adapter edge; duplicate small queries acceptable vs
circular module imports.

## R8 — Duplicate validation messages

**Decision**: Use existing httpx duplicate helpers with keys `error.duplicate.customer.name`,
`error.duplicate.email` matching Java `Response.DUPLICATE_ENTITY`.
