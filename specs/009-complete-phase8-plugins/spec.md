# Feature Specification: Phase 8 — Plugins Platform & Extension Modules

**Feature Branch**: `009-complete-phase8-plugins`

**Created**: 2026-05-21

**Status**: Draft

**Input**: Complete **Phase 8** of the Java→Go MDM migration per
`serverBackendGo/docs/MIGRATION.md`: implement remaining **`plugins/*`** modules with layered
architecture, API parity, and verifiable testing (unit tests, handler tests where valuable, and
documented HTTP smoke) so the React plugin settings screen and legacy plugin REST consumers work
without the Java WAR. Builds on Phases 1–7 (including partial `plugins/push` from Phase 7).

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Tenant plugin catalog and enablement (Priority: P1)

A tenant administrator opens **Plugin settings** in the React control panel, sees which plugins
are available and active for the organization, and enables or disables plugins per customer account.

**Why this priority**: `frontend/src/features/plugins/pluginService.ts` already calls
`/rest/plugin/main/private/*`; without platform parity the settings page fails after login.

**Independent Test**: Authenticated admin GETs active and available plugin lists; POST disabled
plugin IDs; lists reflect changes on subsequent GET.

**Acceptance Scenarios**:

1. **Given** an authenticated user with tenant context, **When** GET active plugins,
   **Then** response lists plugins not globally disabled and not in `pluginsDisabled` for that customer.
2. **Given** same user, **When** GET available plugins,
   **Then** response lists plugins enabled for the current build and permitted for the customer.
3. **Given** user with `plugins_customer_access_management`, **When** POST disabled plugin ID array,
   **Then** customer disabled rows persist and active list excludes those IDs.
4. **Given** user without `plugins_customer_access_management`, **When** POST disabled,
   **Then** permission denied envelope.
5. **Given** unauthenticated caller, **When** private platform endpoints are called,
   **Then** unauthorized.
6. **Given** any client, **When** GET public registered plugins,
   **Then** list of plugins registered in the deployment returns without auth (legacy contract).

---

### User Story 2 — Audit log search for compliance (Priority: P2)

A security or operations administrator searches historical user and API audit records filtered by
date, user, and action to investigate incidents.

**Why this priority**: Audit plugin is a standard enterprise MDM extension; search is the only
admin-facing REST surface in open-source audit.

**Independent Test**: User with `plugin_audit_access` POSTs search filter; paginated audit rows
return; user without permission is denied.

**Acceptance Scenarios**:

1. **Given** audit rows exist for the customer, **When** POST log search with date range and filters,
   **Then** paginated results match filter semantics from Java `AuditResource`.
2. **Given** no matching rows, **When** search runs,
   **Then** empty items with zero total count.
3. **Given** user lacking `plugin_audit_access`, **When** search runs,
   **Then** permission denied.
4. **Given** invalid date or pagination input, **When** search runs,
   **Then** validation or error envelope consistent with legacy.

---

### User Story 3 — Messaging plugin: send and manage device messages (Priority: P2)

An administrator composes messages to devices or groups, reviews message history, purges old
messages, and agents report delivery status via the public status endpoint.

**Why this priority**: Messaging is a distinct plugin from private push; operators rely on it for
SMS-style device communication in the legacy UI.

**Independent Test**: User with `plugin_messaging_send` sends message; history search returns row;
public status update succeeds; purge and delete respect `plugin_messaging_delete`.

**Acceptance Scenarios**:

1. **Given** send permission and valid targets (device numbers, groups, or broadcast),
   **When** POST private send,
   **Then** messages are queued for targeted devices (via shared notification delivery path).
2. **Given** send permission, **When** POST private search with pagination,
   **Then** message history for customer returns.
3. **Given** delete permission, **When** DELETE message by id,
   **Then** message removed for customer scope.
4. **Given** delete permission, **When** GET purge by days,
   **Then** messages older than threshold are removed.
5. **Given** agent callback, **When** GET public status update for message id,
   **Then** delivery status updates without admin session.
6. **Given** missing send permission, **When** send attempted,
   **Then** permission denied.

---

### User Story 4 — Device info plugin: telemetry and console views (Priority: P2)

Devices upload dynamic telemetry; administrators configure plugin settings, view per-device detail,
search devices by dynamic fields, and export results.

**Why this priority**: Device inventory enrichment beyond base `devices` module is required for
support and compliance workflows.

**Independent Test**: Agent PUTs public deviceinfo payload; admin GETs private detail and search;
settings GET/PUT round-trip.

**Acceptance Scenarios**:

1. **Given** known device number, **When** PUT public dynamic info list,
   **Then** data persisted per customer/device scope.
2. **Given** admin with plugin access, **When** GET private device detail by number,
   **Then** aggregated static and dynamic info returned.
3. **Given** admin, **When** GET private search device with query parameters,
   **Then** matching devices listed.
4. **Given** admin, **When** POST private dynamic search or export,
   **Then** filtered dynamic records or export payload per Java contract.
5. **Given** admin, **When** GET/PUT plugin settings private,
   **Then** settings persist for customer.
6. **Given** unknown device on public upload, **When** policy rejects unknown devices,
   **Then** device-not-found or equivalent error.

---

### User Story 5 — Device log plugin: rules, upload, and search (Priority: P2)

Administrators configure log collection rules; devices upload log batches; administrators search and
export logs from the console.

**Why this priority**: Remote troubleshooting depends on devicelog; agents call public upload paths.

**Independent Test**: PUT settings and rule; device POSTs log list; admin search returns entries;
export endpoint returns downloadable result shape.

**Acceptance Scenarios**:

1. **Given** admin, **When** GET/PUT devicelog plugin settings,
   **Then** customer settings round-trip.
2. **Given** admin, **When** PUT new rule or DELETE rule by id,
   **Then** rules list reflects change.
3. **Given** device with applicable rules, **When** GET rules by device number,
   **Then** active rules for device returned.
4. **Given** device, **When** POST log list upload,
   **Then** log records stored for customer/device.
5. **Given** admin, **When** POST private search or search/export,
   **Then** filtered logs returned per Java pagination and filters.
6. **Given** storage backend configured for Postgres (default deployment),
   **When** logs are written,
   **Then** persistence uses Postgres tables from legacy schema (not alternate backends in Phase 8).

---

### User Story 6 — Complete Push plugin schedule tasks (Priority: P3)

An administrator manages scheduled push campaigns: search tasks, create/update tasks, and delete
tasks—extending Phase 7 partial push plugin implementation.

**Why this priority**: Phase 7 delivered message search/send/delete/purge but deferred schedule
task CRUD documented in parity notes.

**Independent Test**: Admin with push permissions uses searchTasks, PUT task, DELETE task; rows
appear in `plugin_push_schedule`.

**Acceptance Scenarios**:

1. **Given** `plugin_push_send`, **When** POST private searchTasks with filter,
   **Then** paginated schedule tasks return.
2. **Given** send permission, **When** PUT private task,
   **Then** task created or updated.
3. **Given** delete permission, **When** DELETE private task by id,
   **Then** task removed.
4. **Given** Phase 7 message endpoints unchanged, **When** regression smoke runs,
   **Then** search/send/delete/purge still pass.

---

### Edge Cases

- Customer has all plugins disabled via `pluginsDisabled` — active list empty; available may still list catalog.
- Plugin globally `disabled=TRUE` in `plugins` table — never appears in active/available lists.
- Build-time plugin identifier not in enabled list — excluded from registered/active responses (matches Java `PluginList`).
- Messaging/push send with empty target sets — validation error, no silent broadcast unless explicitly requested.
- Device log upload exceeds size limits — reject with clear error; do not corrupt partial batches.
- Concurrent disabled-plugin POST — last write wins per customer; cache invalidated for plugin status.
- Multi-tenant isolation — no cross-customer reads on search, messages, audit, or device info/logs.
- Super-admin impersonation — respects impersonated customer id on all plugin queries.

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

| Module | Java reference(s) | REST prefix | Parity doc |
|--------|---------------------|-------------|------------|
| `plugins/platform` | `PluginResource`, `PluginDAO` | `/rest/plugin/main` | `docs/parity/plugins-platform.md` |
| `plugins/audit` | `AuditResource` | `/rest/plugins/audit` | `docs/parity/plugins-audit.md` |
| `plugins/messaging` | `MessagingResource` | `/rest/plugins/messaging` | `docs/parity/plugins-messaging.md` |
| `plugins/deviceinfo` | `DeviceInfoResource`, `DeviceInfoPluginSettingsResource` | `/rest/plugins/deviceinfo` | `docs/parity/plugins-deviceinfo.md` |
| `plugins/devicelog` | `DeviceLogResource`, `DeviceLogPluginSettingsResource` | `/rest/plugins/devicelog` | `docs/parity/plugins-devicelog.md` |
| `plugins/push` (completion) | `PushResource` (task endpoints) | `/rest/plugins/push` | extend `docs/parity/push.md` |

- **Layer ownership**: each subtree under `internal/modules/plugins/<name>/` with `domain/`, `port/`,
  `application/`, `adapter/http/`, `adapter/persistence/postgres/`.
- **No scaffold-only registration**: modules MUST register real routes when `MODULE_*_ENABLED` flags
  are true; remove log-only scaffold registration for implemented modules.
- **Shared delivery**: messaging send and push plugin send SHOULD enqueue via the Phase 7
  notifications queue (`port` injection or shared application service), not duplicate queue logic.
- **Migration**: new `000010_plugins_core` (or next sequence) for any missing tables, permissions,
  and plugin seed rows required for dev/smoke (audit log, messaging messages, deviceinfo/devicelog
  tables per Java liquibase names).
- **Feature flags**: per-module env toggles aligned with existing `MODULE_PLUGINS_*` or new names in
  `.env.example`.

### Functional Requirements

- **FR-001**: System MUST expose `/rest/plugin/main/private/available`, `/private/active`,
  `/public/registered`, and `/private/disabled` with legacy JSON envelopes and permission checks.
- **FR-002**: System MUST persist per-customer disabled plugins in `pluginsDisabled` (or equivalent
  legacy table/columns) and honor them on active plugin resolution.
- **FR-003**: System MUST filter plugin catalog by build-enabled identifiers (configurable list
  defaulting to open-source plugins: audit, push, messaging, deviceinfo, devicelog).
- **FR-004**: System MUST implement `POST /rest/plugins/audit/private/log/search` with pagination
  and filters equivalent to Java audit search.
- **FR-005**: System MUST implement messaging private search, send, delete, purge, and public status
  endpoints under `/rest/plugins/messaging` with permissions `plugin_messaging_send` and
  `plugin_messaging_delete`.
- **FR-006**: System MUST implement deviceinfo public upload, private detail/search/export, and
  plugin settings endpoints under `/rest/plugins/deviceinfo`.
- **FR-007**: System MUST implement devicelog settings, rules CRUD, public rules lookup, log upload,
  and private search/export under `/rest/plugins/devicelog` using Postgres persistence.
- **FR-008**: System MUST complete push plugin schedule task endpoints (`searchTasks`, PUT task,
  DELETE task) on `/rest/plugins/push/private/*`.
- **FR-009**: System MUST enforce tenant scoping and permission names matching Java on every private
  plugin endpoint.
- **FR-010**: System MUST provide `docs/parity/*.md` endpoint tables and update `MIGRATION.md` /
  `NEXT_STEPS.md` marking Phase 8 **done** when acceptance tests pass.
- **FR-011**: System MUST include `go test` coverage for non-trivial `application/` logic per
  module (filter building, target resolution, rule application, plugin list filtering).
- **FR-012**: System MUST document a `quickstart.md` smoke sequence (curl or script) covering
  platform, audit search, messaging send, deviceinfo round-trip, devicelog search, and push task CRUD.
- **FR-013**: Automatic servlet-level audit capture (wrapping all REST responses) MAY be deferred
  and documented as **partial** in parity; manual or seeded audit rows suffice for search parity.
- **FR-014**: Plugin Angular static assets (`javascriptModuleFile`, HTML templates) are OUT of scope;
  React uses platform settings only; legacy Angular plugin UIs remain on Java until separately migrated.
- **FR-015**: `xtra` marketing plugin and premium-only extensions are OUT of scope for Phase 8.

### Key Entities

- **Plugin**: Catalog entry (identifier, name, localization key, permissions, disabled flag, UI metadata).
- **DisabledPlugin**: Association of plugin id and customer id for tenant opt-out.
- **AuditLogRecord**: User action audit row (timestamp, user, request, payload summary, customer).
- **PluginMessage**: Messaging plugin outbound/history row with delivery status.
- **DeviceDynamicInfo**: Key/value or structured telemetry tied to device and customer.
- **DeviceInfoPluginSettings**: Per-customer configuration for device info collection.
- **DeviceLogRule**: Rule defining what logs to collect per device/group.
- **DeviceLogRecord**: Stored log line/batch from device upload.
- **DeviceLogPluginSettings**: Per-customer devicelog configuration.
- **PluginPushSchedule**: Scheduled push task (Phase 7 table; CRUD completed in Phase 8).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Administrator can load Plugin settings in React and save disabled plugins without
  Java backend (zero 404/500 on `/rest/plugin/main/private/*` in smoke).
- **SC-002**: 100% of Phase 8 private plugin endpoints listed in parity docs respond with correct
  auth/permission behavior (OK or documented permission denied), not scaffold silence.
- **SC-003**: Documented smoke script completes in under 10 minutes on a fresh `make migrate` dev DB.
- **SC-004**: `go test ./internal/modules/plugins/...` and affected shared packages pass in CI-local run.
- **SC-005**: At least one independent user story per plugin family (P1–P3) is demonstrable via smoke
  without manual DB edits beyond documented seed/migration.
- **SC-006**: Cross-tenant data leak tests (negative cases in smoke or unit tests) show no foreign
  customer rows in search results.

## Assumptions

- Legacy Postgres schema from Java liquibase is the source of truth; Go migrations only add missing
  objects/seeds for local dev, not rename tables.
- Phase 7 notifications queue remains the delivery mechanism for plugin-originated device messages.
- Open-source plugin set is sufficient for MVP; additional commercial plugins are excluded.
- Device log persistence uses Postgres module path only (matches default `context-docker.xml`).
- Build-enabled plugin identifiers are configured via environment variable with sensible defaults.
- Existing `plugins/push` message endpoints from Phase 7 remain stable; Phase 8 adds schedule tasks only.
- Audit automatic request interception is optional deferral; search API is required for Phase 8 done.

## Dependencies

- Phases 1–7 complete (auth, devices, groups, configurations, notifications queue, push permissions).
- Postgres with `plugins`, `pluginsDisabled`, and plugin-specific tables populated by migration or legacy import.
- React `PluginSettingsPage` for platform UI validation; other plugin UIs may still target Java during transition.
