# Feature Specification: Phase 4 — Devices & Groups Module Migration

**Feature Branch**: `005-complete-phase-devices`

**Created**: 2026-05-20

**Status**: Draft

**Input**: Complete **Phase 4** of the Java→Go MDM migration per
`serverBackendGo/docs/MIGRATION.md`: full **`devices`** and **`groups`** modules so tenant
administrators can list, filter, create, edit, and delete enrolled devices and device groups in
the React MDM console without the Java WAR. Unblocks real dashboard statistics currently served
as empty placeholders by the `summary` module.

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Search and browse devices (Priority: P1)

A tenant administrator opens the **Devices** page, applies filters (text, group, configuration,
status, dates), and sees a paginated list of devices with configuration summary suitable for
the grid.

**Why this priority**: `POST /private/devices/search` is the primary Devices screen API;
without it the core MDM workflow is blocked.

**Independent Test**: Authenticated user with device permissions → POST search with page/size
→ receives device list view with total count and configuration metadata expected by React.

**Acceptance Scenarios**:

1. **Given** a tenant with enrolled devices, **When** they search with default paging,
   **Then** matching devices are returned in the legacy list shape (`configurations` +
   paginated `devices` / `items`).
2. **Given** filters (group, configuration id, status, date range), **When** applied,
   **Then** results respect tenant scope and filter criteria like Java.
3. **Given** no devices match, **When** search runs,
   **Then** success with empty list and zero total.
4. **Given** no authentication, **When** search is requested,
   **Then** access is denied.

---

### User Story 2 — Load device groups for filters and pickers (Priority: P1)

An administrator loads all device groups for the tenant to filter devices or assign bulk group
membership.

**Why this priority**: React `deviceService.getGroups()` and `groupService` call
`GET /private/groups/search` on Devices and Groups screens.

**Independent Test**: Authenticated user → GET groups search → array of id/name groups for
current customer.

**Acceptance Scenarios**:

1. **Given** existing groups for the tenant, **When** groups search is requested,
   **Then** all groups for that customer are returned.
2. **Given** user without `settings` permission where required for mutations,
   **When** read-only list is requested,
   **Then** list still succeeds for users who can view devices (read path matches Java).

---

### User Story 3 — Open, create, edit, and delete a device (Priority: P2)

An administrator views one device by number, creates a new device, updates fields (configuration,
groups, IMEI metadata), or deletes a device they manage.

**Why this priority**: Device detail and add/edit flows depend on GET by number, PUT save,
DELETE by id.

**Independent Test**: GET `/number/{n}` → device view; PUT without id → create; PUT with id →
update; DELETE → removed from search.

**Acceptance Scenarios**:

1. **Given** a valid device number in tenant scope, **When** GET by number,
   **Then** device details are returned or a not-found error consistent with Java.
2. **Given** `edit_devices` permission and valid payload, **When** PUT create,
   **Then** device is persisted and appears in search (subject to license/device limit rules).
3. **Given** duplicate device number in tenant, **When** create/update,
   **Then** duplicate-device error envelope matches legacy.
4. **Given** `edit_devices` permission, **When** DELETE by id,
   **Then** device is removed.
5. **Given** user without `edit_devices`, **When** mutate,
   **Then** permission denied.

---

### User Story 4 — Manage device groups (Priority: P2)

An administrator with settings permission creates, renames, or deletes device groups used to
organize devices.

**Why this priority**: Groups admin UI uses PUT `/private/groups` and DELETE `/private/groups/{id}`.

**Independent Test**: PUT create group → appears in search; DELETE empty group → success;
DELETE group with devices → `error.notempty.group`.

**Acceptance Scenarios**:

1. **Given** `settings` permission and unique name, **When** PUT without id,
   **Then** group is created; creator gains group access if user is group-scoped (Java behavior).
2. **Given** duplicate group name, **When** save,
   **Then** `error.duplicate.group`.
3. **Given** group with assigned devices, **When** delete,
   **Then** error not empty group; **Given** empty group, **Then** delete succeeds.

---

### User Story 5 — Bulk device operations (Priority: P2)

An administrator selects multiple devices and deletes them in bulk or assigns/clears group
membership in bulk.

**Why this priority**: React bulk actions use `POST /deleteBulk` and `POST /groupBulk`.

**Independent Test**: Select ids → bulk delete → devices gone; bulk set groups → membership updated.

**Acceptance Scenarios**:

1. **Given** `edit_devices` and valid id list, **When** delete bulk,
   **Then** all listed devices are removed.
2. **Given** `edit_devices` and group bulk request with action `set` or clear,
   **Then** device-group links update per Java semantics.

---

### User Story 6 — Device autocomplete and description (Priority: P3)

An administrator uses quick search autocomplete when filtering devices, or edits a device
description inline.

**Why this priority**: Secondary UX; autocomplete POST and description POST are used but not
blocking first list load.

**Independent Test**: POST autocomplete with partial text → up to 10 suggestions; POST
description → persisted on device.

**Acceptance Scenarios**:

1. **Given** partial device number/name, **When** autocomplete,
   **Then** up to 10 matching lookup items returned.
2. **Given** `edit_device_desc` permission, **When** description saved,
   **Then** description updates; **Given** without permission, **Then** denied.

---

### User Story 7 — Per-device application settings (Priority: P3)

An administrator views and saves per-device application settings overrides and optionally
notifies the device.

**Why this priority**: Device detail advanced panel; lower traffic than list/search.

**Independent Test**: GET applicationSettings → list; POST save → OK; POST notify → OK (push
may be stubbed).

**Acceptance Scenarios**:

1. **Given** device id in tenant, **When** GET application settings,
   **Then** settings list returned (possibly empty).
2. **Given** valid settings payload, **When** POST save,
   **Then** persisted successfully.
3. **Given** notify requested, **When** POST notify,
   **Then** success (actual push delivery out of scope; no error to client).

---

### User Story 8 — Real dashboard device statistics (Priority: P2)

After devices exist in the database, the dashboard shows accurate online/offline/enrollment
counts instead of empty placeholders.

**Why this priority**: Phase 3 `summary` module is partial until device tables exist.

**Independent Test**: Seed devices → GET `/private/summary/devices` → non-zero/status breakdown
matching data.

**Acceptance Scenarios**:

1. **Given** devices with varied last-update timestamps, **When** dashboard stats requested,
   **Then** counts reflect tenant scope and legacy status rules.
2. **Given** no devices, **When** stats requested,
   **Then** zeroed summary still returns valid structure.

---

### User Story 9 — Configuration list for device UI (Priority: P2)

The Devices screen configuration filter and create/edit dropdown load available configurations
(id + name) for the tenant.

**Why this priority**: React calls `GET /private/configurations/list` from the devices page;
without it filters break even if devices search works.

**Independent Test**: Authenticated GET configurations list → id/name array for tenant.

**Acceptance Scenarios**:

1. **Given** configurations seeded for customer, **When** list requested,
   **Then** minimal lookup list returned.
2. **Given** no configurations, **When** list requested,
   **Then** empty array success (not 404).

---

### User Story 10 — Verifiable API and regression safety (Priority: P2)

Developers exercise all Phase 4 endpoints via Swagger and automated tests without Java.

**Independent Test**: Bearer auth → smoke critical paths; `go test` for devices/groups modules.

**Acceptance Scenarios**:

1. **Given** regenerated Swagger, **When** browsing UI,
   **Then** Devices and Groups tags list all in-scope endpoints with Bearer auth.
2. **Given** module tests, **When** run in CI/local,
   **Then** search, CRUD, and permission checks covered.

---

### Edge Cases

- Device search across customers → only current tenant (`customerId` from principal).
- Device limit / license: create rejected when limit exceeded (match Java settings).
- Bulk operations with empty id list → success no-op.
- Group delete with devices → `error.notempty.group`.
- Configuration missing for device row → device omitted or handled like Java filter in list view.
- Configuration enrichment heavy (apps/files) → may be simplified in v1 with documented partial parity.
- Push notifications for config/app changes → stub OK, no agent delivery required in Phase 4.

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Modules**: `internal/modules/devices/`, `internal/modules/groups/` (replace scaffolds);
  extend `internal/modules/summary/` for real stats; minimal read in
  `internal/modules/configurations/` for list-only (or dedicated thin adapter).
- **Phase**: 4 in `MIGRATION.md`; marks Phase 4 **done** when both modules ship.
- **Java reference**: `DeviceResource`, `GroupResource`, `DeviceDAO`, `GroupDAO`,
  `SummaryResource` / device stats queries.
- **REST bases**:
  - `/rest/private/devices` — all `DeviceResource` in-scope paths
  - `/rest/private/groups` — all `GroupResource` in-scope paths
  - `/rest/private/summary/devices` — upgrade existing handler/repo
  - `/rest/private/configurations/list` — minimal read for React devices page
- **Parity docs**: `docs/parity/devices.md`, `docs/parity/groups.md`; update `summary.md`.
- **Layers**: domain/port/application/adapter per module; shared permission checks via
  `platformauth` (`edit_devices`, `edit_device_desc`, `settings`).
- **Migrations**: New SQL migrations for `devices`, `groups`, `devicegroups` (or equivalent
  junction), `configurations` minimal schema if absent — idempotent additive files.
- **Auth**: JWT/session on private routes; tenant scoping mandatory.

### Functional Requirements — Devices

- **FR-D01**: System MUST expose `POST /rest/private/devices/search` returning legacy
  `DeviceListView`-compatible JSON for React.
- **FR-D02**: System MUST expose `GET /rest/private/devices/number/{number}` for device detail.
- **FR-D03**: System MUST expose `PUT /rest/private/devices` for create, update, and bulk
  configuration update (`ids` + `configurationId`) per Java.
- **FR-D04**: System MUST expose `DELETE /rest/private/devices/{id}` with `edit_devices` check.
- **FR-D05**: System MUST expose `POST /rest/private/devices/deleteBulk` and
  `POST /rest/private/devices/groupBulk`.
- **FR-D06**: System MUST expose `POST /rest/private/devices/autocomplete`.
- **FR-D07**: System MUST expose `GET/POST /rest/private/devices/{id}/applicationSettings` and
  `POST .../notify` (notify may no-op push).
- **FR-D08**: System MUST expose `POST /rest/private/devices/{id}/description` with
  `edit_device_desc` permission.
- **FR-D09**: Device mutations MUST enforce tenant isolation and permission names matching Java.
- **FR-D10**: Duplicate device number MUST return legacy device-exists error semantics.

### Functional Requirements — Groups

- **FR-G01**: System MUST expose `GET /rest/private/groups/search` (and optional
  `GET /search/{value}` for parity).
- **FR-G02**: System MUST expose `POST /rest/private/groups/autocomplete`.
- **FR-G03**: System MUST expose `PUT /rest/private/groups` with `settings` permission and
  duplicate name handling.
- **FR-G04**: System MUST expose `DELETE /rest/private/groups/{id}` with not-empty guard.

### Functional Requirements — Cross-cutting

- **FR-X01**: System MUST add database schema migrations required for devices/groups (and
  minimal configurations if missing).
- **FR-X02**: `GET /rest/private/summary/devices` MUST return real aggregates when device data exists.
- **FR-X03**: System MUST expose `GET /rest/private/configurations/list` as minimal id/name
  list for the tenant (read-only).
- **FR-X04**: All responses MUST use Headwind envelope.
- **FR-X05**: React Devices and Groups pages MUST function against Go-only backend for
  in-scope calls in `deviceService.ts` and `groupService.ts`.
- **FR-X06**: Swagger MUST document Phase 4 endpoints with Bearer auth after `make swagger`.
- **FR-X07**: Automated tests MUST cover devices search, device CRUD permission, groups CRUD.

### Key Entities

- **Device**: Tenant-owned enrolled endpoint (number, configuration, status fields, groups,
  description, application settings overrides, last update/enrollment metadata).
- **Device group**: Named collection; many-to-many with devices.
- **Device search request**: Pagination, filters (group, configuration, status, dates, text, etc.).
- **Device list view**: Configurations map + paginated device views for UI grid.
- **Configuration (minimal)**: Id and name for dropdowns (full CRUD remains Phase 5).
- **Device statistics**: Dashboard counters derived from device records.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An administrator can load the Devices page list (search) in under 5 seconds on a
  seeded database with 500 devices.
- **SC-002**: Create, update, delete device and group flows complete without Java for the same
  tenant admin role.
- **SC-003**: Dashboard summary shows non-placeholder stats when devices are seeded.
- **SC-004**: 100% of in-scope `DeviceResource` and `GroupResource` endpoints used by React are
  marked Done (or Partial with notes) in parity docs.
- **SC-005**: Phase 4 row in `MIGRATION.md` moves from pending to **done**.
- **SC-006**: Configuration filter dropdown on Devices page loads without calling Java.

## Assumptions

- Phases 1–3 (auth, users/roles, customers/settings/hints/summary scaffold) remain available.
- React remains the primary UI validator; Android agent sync is Phase 7 (out of scope).
- Push notification delivery (FCM/APNs) is stubbed in Phase 4; endpoints return success.
- Full `ConfigurationResource` CRUD and application catalog remain **Phase 5**; Phase 4 only
  delivers list-read needed by Devices UI.
- Device search configuration enrichment (embedded apps/files in list) may ship as **partial**
  in v1 (core device fields + config name first).
- Deprecated group `GET /search/{value}` included for parity if low cost.

## Dependencies

- **Blocks**: Phase 5 applications/configurations full migration (easier after devices exist).
- **Requires**: Postgres migrations; permission constants aligned with Java role permissions.
- **Enables**: Customer create default devices (Phase 3 partial), maps page device coordinates.

## Out of Scope (Phase 4)

- Phase 5–8 modules (applications, files, sync, plugins, …).
- Public device enrollment / QR provisioning APIs (later phases).
- Real push notification infrastructure.
- Full configuration editor and config file management.
- Plugin audit hooks for device/group mutations.
