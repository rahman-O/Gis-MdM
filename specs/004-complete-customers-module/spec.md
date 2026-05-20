# Feature Specification: Phase 3 — Customers Module Migration

**Feature Branch**: `004-complete-customers-module`

**Created**: 2026-05-20

**Status**: Draft

**Input**: Complete remaining **Phase 3** work for the Java→Go MDM migration by delivering the
**customers** module (`serverBackendGo/docs/MIGRATION.md`). Summary, settings, and hints are
already done; customers is the last Phase 3 module still at scaffold-only. Super administrators
must manage tenant customer accounts, search them from the control panel, and impersonate an
organization admin to support a tenant—without the Java WAR.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Search customer accounts (Priority: P1)

A platform super administrator opens the **Control Panel** and searches customer (tenant)
accounts by name or other criteria with paging so they can pick a tenant to support.

**Why this priority**: The React control panel (`ControlPanelPage`) calls customer search on
every load and filter; without it the page is empty.

**Independent Test**: Authenticated super admin submits a paginated search → receives a list
of customer records and total count in the standard success envelope.

**Acceptance Scenarios**:

1. **Given** a super admin with existing customer accounts, **When** they search with page
   size and optional text filter, **Then** matching customers are returned with pagination
   metadata (total count and current page of records).
2. **Given** a search with no matches, **When** submitted,
   **Then** the response is success with an empty record set and zero (or consistent) total.
3. **Given** no authenticated session, **When** search is requested,
   **Then** access is denied, consistent with other protected private routes.
4. **Given** a non–super-admin user, **When** they attempt customer search,
   **Then** access is denied (same permission model as legacy super-admin-only customer APIs).

---

### User Story 2 — Impersonate a tenant organization admin (Priority: P1)

A super administrator selects a customer and **impersonates** that tenant’s organization
administrator so they can troubleshoot the tenant’s MDM console as that admin would see it.

**Why this priority**: React control panel uses impersonation immediately after picking a row;
session/token hydration depends on a successful impersonation response.

**Independent Test**: Super admin requests impersonation for a customer with an org admin →
receives org-admin identity payload suitable for continuing in the web shell (session or token
per existing auth integration).

**Acceptance Scenarios**:

1. **Given** a super admin and a customer with an active org admin user, **When** they
   impersonate that customer, **Then** the operation succeeds and returns the org admin user
   profile fields needed by the client (without exposing the admin’s password).
2. **Given** a customer with no org admin, **When** impersonation is requested,
   **Then** the operation fails with a clear, localized error (not found admin), matching
   legacy behavior.
3. **Given** a non–super-admin user, **When** impersonation is requested,
   **Then** permission is denied.
4. **Given** an org admin who has not yet obtained an auth token, **When** impersonation runs,
   **Then** the system ensures a token exists so the client can continue (same as legacy).
5. **Given** an org admin with an active password-reset token, **When** impersonation is
   attempted, **Then** impersonation is blocked or fails securely (legacy disables
   impersonation in this state).

---

### User Story 3 — Create or update a customer account (Priority: P2)

A super administrator creates a new tenant customer or updates an existing one (name, email,
prefix, configuration defaults, status metadata) from administrative workflows that still rely
on the legacy customer save API.

**Why this priority**: Required for full `CustomerResource` parity and future admin UIs even if
the current React shell only uses search + impersonate.

**Independent Test**: PUT with new customer → success and optional initial admin credentials
returned on create; PUT with existing id → success and linked org admin profile updated when
email/name change.

**Acceptance Scenarios**:

1. **Given** valid new customer data with unique name and email, **When** saved without id,
   **Then** the customer is created and response includes initial administrator credentials
   when legacy does.
2. **Given** duplicate customer name or email, **When** save is attempted,
   **Then** duplicate-entity error with the same message keys as legacy (`error.duplicate.*`).
3. **Given** an existing customer id, **When** updated with valid fields,
   **Then** customer row and associated org admin login/name/email are updated consistently.
4. **Given** create/update with email conflicting with another tenant’s user,
   **Then** duplicate email error is returned.

---

### User Story 4 — Delete a customer account (Priority: P2)

A super administrator removes a customer account that is no longer needed.

**Why this priority**: Legacy exposes delete; parity and data lifecycle require it.

**Independent Test**: DELETE by customer id → success; subsequent search no longer lists it
(subject to legacy cascade rules).

**Acceptance Scenarios**:

1. **Given** an existing customer id, **When** delete is requested by a super admin,
   **Then** the operation returns success.
2. **Given** invalid or unknown id, **When** delete is requested,
   **Then** behavior matches legacy (success no-op or not-found per existing DAO semantics).
3. **Given** unauthorized user, **When** delete is requested,
   **Then** permission denied.

---

### User Story 5 — Load customer details for editing (Priority: P2)

A super administrator opens the edit form for one customer and loads full details including
fields needed for update (prefix, configuration, status, contact fields).

**Why this priority**: Legacy `GET .../edit` supports admin maintenance flows.

**Independent Test**: GET edit by id → customer DTO matching legacy shape for update screens.

**Acceptance Scenarios**:

1. **Given** a valid customer id, **When** edit details are requested,
   **Then** full customer record for update is returned.
2. **Given** unknown id, **When** requested,
   **Then** appropriate error envelope (not internal stack trace).

---

### User Story 6 — Validate device number prefix availability (Priority: P3)

While configuring a customer, the administrator checks whether a proposed device-number
prefix is already used by another tenant.

**Why this priority**: Legacy prefix validation prevents collisions for auto-generated devices.

**Independent Test**: GET prefix check → boolean indicating whether prefix is already taken.

**Acceptance Scenarios**:

1. **Given** an unused prefix, **When** checked,
   **Then** response indicates prefix is not used.
2. **Given** a prefix assigned to another customer, **When** checked,
   **Then** response indicates prefix is used.

---

### User Story 7 — Verifiable API and regression safety (Priority: P2)

A developer validates Phase 3 customer endpoints via interactive API docs and automated tests
without the Java server.

**Why this priority**: Same quality bar as Phase 2 (users/roles) and Phase 3 hints.

**Independent Test**: Authorize as super admin → exercise in-scope endpoints → module tests pass.

**Acceptance Scenarios**:

1. **Given** API documentation with bearer authorization, **When** search and impersonation
   are executed with a super-admin token, **Then** envelopes and status codes match legacy.
2. **Given** the customers module test suite, **When** run locally,
   **Then** search, impersonation, and core CRUD/prefix flows are covered at application or
   HTTP layer.

---

### Edge Cases

- Empty search text with large page size → returns first page without error.
- Impersonation invalidates prior browser session server-side (legacy session swap).
- Create customer with blank email → allowed where legacy allows; duplicate checks skipped for email.
- Concurrent create with same name → one succeeds, other gets duplicate error.
- Customer delete with dependent data → follows legacy cascade/delete rules (no partial orphan state).
- Internal failures → standard error envelope, not raw exception text to clients.

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Module**: `internal/modules/customers/` (replace current scaffold).
- **Phase**: 3 — completes `customers` row in `MIGRATION.md` Phase 3.
- **Java reference**: `com.hmdm.rest.resource.CustomerResource`, `CustomerDAO`, related user lookups.
- **REST base**: `/rest/private/customers` — in-scope paths must match legacy:
  - `POST /search` (paginated; primary client path)
  - `PUT /` (create when `id` null, update when `id` set)
  - `DELETE /{id}`
  - `GET /{id}/edit`
  - `GET /prefix/{prefix}/used`
  - `GET /impersonate/{id}`
- **Out of scope (deprecated in Java)**: `GET /search`, `GET /search/{value}` — not required for React control panel.
- **Parity doc**: `serverBackendGo/docs/parity/customers.md`.
- **Layers**: `domain` (Customer aggregate fields), `port` (repository), `application` (use cases),
  `adapter/http`, `adapter/persistence/postgres`.
- **Auth**: Private routes; super-admin checks for customer management and impersonation;
  same session/JWT middleware as Phase 2; document Bearer auth in API docs for tests.

### Functional Requirements

- **FR-001**: System MUST expose paginated customer search via `POST /rest/private/customers/search`
  accepting page, page size, optional search text, sort, account type, and customer status filters;
  response `data` MUST match legacy paginated shape (`items`/`records` + `totalItemsCount` or
  equivalent consumed by React `unwrapCustomerRows`).
- **FR-002**: System MUST restrict customer search and mutations to users with super-admin
  privileges (consistent with legacy `SecurityContext` checks on impersonation and implied
  admin-only customer management).
- **FR-003**: System MUST expose `GET /rest/private/customers/impersonate/{id}` for super admins,
  returning org-admin user payload for client session hydration; MUST NOT return password fields.
- **FR-004**: System MUST expose `PUT /rest/private/customers` for create/update with duplicate
  name/email validation and org-admin sync on update per legacy rules.
- **FR-005**: On customer create, when legacy returns initial admin credentials, response MUST
  include `adminCredentials` in `data` map for the creating administrator.
- **FR-006**: System MUST expose `DELETE /rest/private/customers/{id}` for super admins.
- **FR-007**: System MUST expose `GET /rest/private/customers/{id}/edit` returning customer
  details for update forms.
- **FR-008**: System MUST expose `GET /rest/private/customers/prefix/{prefix}/used` returning
  whether the prefix is already assigned.
- **FR-009**: All responses MUST use the Headwind envelope (`status`, `message`, `data`).
- **FR-010**: React Control Panel MUST work against Go-only backend for search and impersonation
  without frontend code changes.
- **FR-011**: System MUST include automated tests for application logic and HTTP handlers for
  at least search and impersonation.
- **FR-012**: In-scope endpoints MUST be documented in interactive API docs with bearer authorization.

### Key Entities

- **Customer (tenant account)**: Organization boundary in MDM—name, email, description, device
  number prefix, default device configuration, account type, customer status, registration and
  last-login metadata, master flag, files directory.
- **Customer search request**: Pagination and filter criteria (page, size, text, sort, account
  type, status).
- **Paginated customers (API view)**: Page of customer rows plus total item count for UI tables.
- **Organization admin user**: Primary admin user for a customer; target of impersonation;
  linked by customer id.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A super admin can load the control panel customer list (search) in under 3 seconds
  on a seeded development database with at least 100 customers.
- **SC-002**: Impersonation from the control panel allows the super admin to reach the tenant
  dashboard in one action without restarting the Java backend.
- **SC-003**: Create/update/delete/prefix-check behaviors are indistinguishable from Java for the
  same inputs (verified via parity smoke checklist).
- **SC-004**: 100% of in-scope `CustomerResource` endpoints (excluding deprecated GET search) are
  marked Done in parity documentation.
- **SC-005**: Phase 3 migration roadmap lists all four Phase 3 modules (customers, settings,
  hints, summary) as complete after this feature ships.

## Assumptions

- Phase 1 auth and Phase 2 users/roles (permissions, principal enrichment) remain available.
- Existing PostgreSQL schema already contains `customers` and related user tables from legacy
  Liquibase/migrations; no greenfield schema redesign in this feature.
- **Mailchimp** marketing subscribe on customer create is **out of scope**; customer creation
  still succeeds without external newsletter integration.
- Deprecated Java GET search endpoints are not used by the current React app and need not be
  reimplemented unless a parity audit later requires them.
- Only super administrators manage customers and impersonate; tenant org admins do not access
  these endpoints.
- Impersonation session semantics follow existing Go auth adapter (cookie session and/or JWT
  token fields returned to React `hydrateSessionAfterImpersonation`).

## Dependencies

- **Completed**: auth (Phase 1), users/roles (Phase 2), hints/settings/summary (Phase 3 partial).
- **Blocks**: None for Phase 4 devices/groups beyond shared user/customer FK data already present.
- **Uses**: User lookup/update for org admin sync and impersonation token generation.
