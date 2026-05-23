# Feature Specification: Phase 7 — Agent Sync, Push, Notifications, Updates & QR

**Feature Branch**: `008-complete-phase7-sync-agent`

**Created**: 2026-05-21

**Status**: Draft

**Input**: Complete **Phase 7** of the Java→Go MDM migration per
`serverBackendGo/docs/MIGRATION.md`: implement **`sync`**, **`push`**, **`notifications`**,
**`updates`**, and **`qrcode`** with clean layered architecture and full API parity so enrolled
Android agents receive configuration and commands, administrators can enroll devices via QR and send
push messages from React, and the control panel can check for product updates—without the Java WAR.
Builds on Phases 1–6 (auth through files/icons/public API).

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Device enrollment and configuration sync (Priority: P1)

An Android MDM agent (launcher or secondary APK) enrolls using a device identifier, receives the
full device configuration payload (applications, files, settings, certificates metadata), and
periodically re-syncs when the administrator changes policies.

**Why this priority**: Without `/rest/public/sync`, agents cannot operate; this is the core MDM
runtime path and blocks production cutover.

**Independent Test**: POST or GET `/rest/public/sync/configuration/{deviceId}` for a seeded device
returns a `SyncResponse`-shaped payload; unknown device with valid create-on-demand options enrolls;
duplicate enrollment is blocked when configured.

**Acceptance Scenarios**:

1. **Given** a registered device, **When** GET configuration sync runs,
   **Then** response includes device record, configuration, applications, files, and settings
   expected by the agent.
2. **Given** POST enrollment with `DeviceCreateOptions`, **When** device does not exist and policy
   allows creation, **Then** device is created and configuration is returned.
3. **Given** `preventDuplicateEnrollment` enabled and device already enrolled, **When** enrollment
   POST runs, **Then** device-exists error matches legacy semantics.
4. **Given** `secureEnrollment` enabled, **When** request lacks valid `X-Request-Signature`,
   **Then** permission denied.
5. **Given** valid sync response, **When** response is produced,
   **Then** `X-Response-Signature` header is set per legacy (hash over payload).
6. **Given** device lookup by IMEI/serial when number unknown, **When** match exists,
   **Then** device is resolved and migration flags handled like Java.
7. **Given** IMEI/serial migration path, **When** sync completes,
   **Then** `X-IP-Address` and CPU arch headers are honored where applicable.

---

### User Story 2 — Device telemetry and per-app settings (Priority: P1)

An enrolled agent reports device info (battery, location, custom fields, IMEI changes) and
application-level settings so the console reflects live device state.

**Why this priority**: Operational visibility and compliance depend on `/info` and
`/applicationSettings` endpoints used on every agent heartbeat cycle.

**Independent Test**: POST `/rest/public/sync/info` updates stored device info; POST
`/rest/public/sync/applicationSettings/{deviceId}` persists settings; unknown device returns
device-not-found.

**Acceptance Scenarios**:

1. **Given** known device id, **When** POST info with telemetry JSON,
   **Then** device info and custom properties update; battery/location events may be recorded
   (event bus may be simplified stub with parity note).
2. **Given** single-customer mode and unknown device, **When** POST info,
   **Then** device may be created on demand per Java rules.
3. **Given** multi-tenant mode and unknown device, **When** POST info,
   **Then** creation is rejected.
4. **Given** application settings list, **When** POST applicationSettings,
   **Then** settings persist keyed by package and name.
5. **Given** unknown device number, **When** either endpoint runs,
   **Then** device-not-found envelope returns.

---

### User Story 3 — Agent notification delivery (Priority: P1)

An agent polls for pending push/command messages and receives them reliably while online, using
either the REST pull endpoint or long-polling servlet compatible with legacy agents.

**Why this priority**: Remote wipe, config refresh triggers, and messaging depend on notification
delivery infrastructure.

**Independent Test**: After admin queues a message, GET
`/rest/notifications/device/{deviceNumber}` returns pending messages; long-poll endpoint holds until
message available or timeout.

**Acceptance Scenarios**:

1. **Given** pending rows for device, **When** GET notifications by device number,
   **Then** plain push message list returns and messages are marked for delivery per Java.
2. **Given** device number not found (including old number fallback), **When** GET runs,
   **Then** device-not-found error.
3. **Given** agent long-polling, **When** GET `/rest/notification/polling/{deviceNumber}` with valid
   signature (if secure enrollment), **Then** connection completes with message JSON or timeout empty
   response per servlet behavior.
4. **Given** invalid polling signature when secure enrollment on, **When** poll runs,
   **Then** request rejected.
5. **Given** no external FCM credentials, **When** messages are queued,
   **Then** HTTP pull/polling path still delivers (no requirement for Google FCM in Phase 7).

---

### User Story 4 — Administrator sends push from console (Priority: P2)

A tenant administrator sends a push/command to one device, a group, or a broadcast from the React
devices UI without Java.

**Why this priority**: `frontend/src/features/push/pushService.ts` calls `POST /rest/private/push`;
this is the primary modern UI contract (distinct from legacy Push plugin screens).

**Independent Test**: User with `push_api` permission POSTs message; devices receive queued rows;
user without permission is denied.

**Acceptance Scenarios**:

1. **Given** `push_api` permission and valid payload (device numbers, groups, or broadcast),
   **When** POST `/rest/private/push`,
   **Then** push messages are queued for targeted devices.
2. **Given** missing permission, **When** POST runs,
   **Then** permission denied.
3. **Given** invalid device number in targeting list, **When** POST runs,
   **Then** error envelope explains failure without partial silent success (match Java).
4. **Given** group targeting, **When** POST runs,
   **Then** all devices in group for current tenant/user scope receive queued messages.

---

### User Story 5 — QR enrollment for configurations (Priority: P2)

An administrator opens the enrollment QR page for a configuration, displays a scannable PNG, and
optionally downloads provisioning JSON for offline tooling.

**Why this priority**: React `EnrollmentQrPage` and `enrollmentQrQuery.ts` depend on public QR
endpoints; enrollment UX is a common onboarding path.

**Independent Test**: GET `/rest/public/qr/{qrCodeKey}` returns PNG; GET
`/rest/public/qr/json/{qrCodeKey}` returns provisioning JSON; invalid key returns server error
consistent with Java.

**Acceptance Scenarios**:

1. **Given** configuration with `qrCodeKey`, **When** GET QR image with size/deviceId/create/group
   query params,
   **Then** PNG encodes enrollment payload with launcher URL and APK hash when main app configured.
2. **Given** same key, **When** GET json variant,
   **Then** JSON extras bundle matches Java `generateExtrasBundle` fields needed by Android setup.
3. **Given** unknown qr key, **When** either endpoint runs,
   **Then** error response (500 in Java — document same for parity).
4. **Given** localhost file URL in dev, **When** QR generated,
   **Then** loopback rewrite rules match Java for agent reachability.

---

### User Story 6 — Check and apply product updates (Priority: P2)

A super-administrator checks for Headwind MDM component updates (web panel, launcher, secondary APKs)
from the Updates UI and optionally downloads or applies updates.

**Why this priority**: `updatesService.ts` calls `GET /rest/private/update/check`; operations teams
use this for maintenance.

**Independent Test**: Super-admin GET check returns manifest entries; non-super-admin in multi-tenant
mode is rejected; POST download/update processes entries per request flags.

**Acceptance Scenarios**:

1. **Given** single-tenant or super-admin, **When** GET check runs,
   **Then** update manifest entries return with outdated/downloaded flags computed like Java.
2. **Given** non-super-admin in multi-tenant deployment, **When** GET check runs,
   **Then** security/permission error.
3. **Given** update request with download flags, **When** POST update runs,
   **Then** eligible APKs download to files directory and web manifest updated when applicable.
4. **Given** `update` flag in request, **When** mobile app entries outdated,
   **Then** configuration application versions may be upgraded per Java rules.
5. **Given** `sendStats` flag, **When** POST completes,
   **Then** stats may be sent to vendor URL (may be **partial** stub with parity note).

---

### User Story 7 — Push plugin administration (Priority: P3)

An administrator using the legacy Push plugin UI (or API clients) searches message history, sends
plugin-scoped messages, purges old messages, and manages scheduled push tasks.

**Why this priority**: Parity for `/rest/plugins/push/*` used by Angular plugin module; lower priority
than agent sync and React `/private/push`.

**Independent Test**: Plugin enabled; POST search returns paginated messages; send requires
`plugin_push_send`; purge and task CRUD behave per Java.

**Acceptance Scenarios**:

1. **Given** plugin tables migrated, **When** POST `/rest/plugins/push/private/search`,
   **Then** paginated plugin push messages return for tenant filter.
2. **Given** `plugin_push_send`, **When** POST `/rest/plugins/push/private/send`,
   **Then** messages queue for device/group/configuration scope.
3. **Given** message id, **When** DELETE `/rest/plugins/push/private/{id}`,
   **Then** row removed for tenant.
4. **Given** days parameter, **When** GET purge,
   **Then** old messages purged count returned.
5. **Given** schedule APIs, **When** searchTasks/saveTask/deleteTask used,
   **Then** CRUD matches Java plugin behavior.

---

### User Story 8 — Verifiable API and regression safety (Priority: P2)

Developers exercise Phase 7 endpoints via Swagger and automated tests without Java.

**Independent Test**: Public sync smoke without JWT; private push/update with Bearer; module
`go test` for hash/signature, permission guards, and QR key resolution.

**Acceptance Scenarios**:

1. **Given** regenerated Swagger, **When** browsing UI,
   **Then** Sync, Notifications, Push, Updates, and QR tags list in-scope endpoints.
2. **Given** module tests, **When** run locally/CI,
   **Then** sync signature, push permission, and QR unknown-key paths are covered.

---

### Edge Cases

- Cross-tenant device access via sync or notifications → not found or denied.
- Enrollment signature mismatch under `secureEnrollment` → permission denied on sync and polling.
- Device migration (`oldNumber`) during sync/info → migration completed when new id used.
- Sync response hooks from Java plugins → not executed unless hook interface stubbed (document
  **partial**).
- Push message flood to large groups → bounded batch insert; same practical limits as Java
  (pageSize 1M device search).
- Long-poll timeout with no messages → empty or timeout response per servlet contract.
- Update manifest download failure → entry marked update-disabled with reason.
- QR configuration without main app version → error or degraded payload per Java.
- External manifest URL unreachable on update check → ERROR envelope.
- MQTT-only push mode in configuration → HTTP polling still required in Phase 7 (agent compatibility).

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Modules**: Replace scaffolds in `internal/modules/sync/`, `internal/modules/push/`,
  `internal/modules/notifications/`, `internal/modules/updates/`, `internal/modules/qrcode/`, and
  implement `internal/modules/plugins/push/` (plugin subtree already registered in `app/modules.go`).
- **Phase**: 7 in `MIGRATION.md`; marks Phase 7 **done** when agent + admin contracts in scope ship.
- **Java reference**:
  - `SyncResource` — `/rest/public/sync`
  - `NotificationResource` + `LongPollingServlet` — `/rest/notifications`, `/rest/notification/polling`
  - `PushApiResource` — `/rest/private/push`
  - `PushResource` (plugin) — `/rest/plugins/push`
  - `UpdateResource` — `/rest/private/update`
  - `QRCodeResource` — `/rest/public/qr`
  - DAOs: `UnsecureDAO`, `DeviceDAO`, `NotificationDAO`, `PushDAO`, `PushScheduleDAO`, `ApplicationDAO`
- **REST bases** (must match legacy paths):
  - `/rest/public/sync`
  - `/rest/notifications` (JAX-RS mount)
  - `/rest/notification/polling` (servlet long-poll)
  - `/rest/private/push`
  - `/rest/plugins/push`
  - `/rest/private/update`
  - `/rest/public/qr`
- **Parity docs**: `docs/parity/sync.md`, `notifications.md`, `push.md`, `updates.md`, `qrcode.md`,
  and `docs/parity/plugins-push.md` (or section in `push.md`).
- **Layers**: `domain/`, `port/`, `application/`, `adapter/http/`, `adapter/persistence/postgres/`
  per module; long-polling may use `adapter/http/polling.go` or platform servlet-style handler.
- **Permissions**: `push_api` (PushApiResource); `plugin_push_send`, `plugin_push_delete` (plugin);
  updates check restricted to super-admin in multi-tenant mode.
- **Config**: Reuse `HASH_SECRET`, `SECURE_ENROLLMENT`, `PREVENT_DUPLICATE_ENROLLMENT`, `BASE_URL`,
  `FILES_DIRECTORY`, rebranding mobile/vendor names, `POLLING_TIMEOUT`, proxy IP headers for sync.
- **Migrations**: Additive SQL for `pushmessages`, `plugin_push_messages`, `plugin_push_schedule`
  (names per Java Liquibase) if absent in Postgres; seed permissions for role 2 where applicable.
- **Shared**: Reuse `internal/shared/crypto` for enrollment signatures; reuse `platform/storage` for
  update APK downloads and QR APK hash calculation.

### Functional Requirements — Sync (public)

- **FR-SY01**: System MUST expose `POST` and `GET` `/rest/public/sync/configuration/{deviceId}` with
  enrollment, migration, and on-demand device creation semantics.
- **FR-SY02**: System MUST expose `POST /rest/public/sync/info` for device telemetry updates.
- **FR-SY03**: System MUST expose `POST /rest/public/sync/applicationSettings/{deviceId}` for per-app
  settings persistence.
- **FR-SY04**: System MUST honor `secureEnrollment` request signature validation using
  `hashSecret + deviceId` (legacy `CryptoUtil.checkRequestSignature`).
- **FR-SY05**: System MUST attach `X-Response-Signature` on configuration sync responses.
- **FR-SY06**: System MUST build `SyncResponse` payload equivalent to Java (configuration, apps,
  files, settings, customer branding fields) for agent consumption.

### Functional Requirements — Notifications (public/agent)

- **FR-NT01**: System MUST expose `GET /rest/notifications/device/{deviceNumber}` returning pending
  messages as plain push message DTOs.
- **FR-NT02**: System MUST expose long-polling at `/rest/notification/polling/{deviceNumber}` with
  async timeout behavior compatible with agents.
- **FR-NT03**: System MUST resolve device by current or old number before delivery.
- **FR-NT04**: System MUST queue and mark messages delivered through a port abstraction shared with
  push modules (no duplicate persistence logic in handlers).

### Functional Requirements — Push (private API)

- **FR-PU01**: System MUST expose `POST /rest/private/push` accepting React `PushPayload` shape
  (messageType, payload, deviceNumbers, groups, broadcast).
- **FR-PU02**: System MUST require `push_api` permission for private push API.
- **FR-PU03**: System MUST expand group/broadcast targeting to device list using tenant-scoped device
  search equivalent to Java `DeviceSearchRequest`.

### Functional Requirements — Push plugin

- **FR-PP01**: System MUST expose plugin routes under `/rest/plugins/push/private/*` for search,
  send, delete, purge, and schedule task CRUD per `PushResource`.
- **FR-PP02**: System MUST enforce `plugin_push_send` and `plugin_push_delete` permissions on
  mutating plugin operations.
- **FR-PP03**: System MUST scope all plugin push records to current customer id.

### Functional Requirements — Updates (private)

- **FR-UP01**: System MUST expose `GET /rest/private/update/check` with super-admin guard in
  multi-tenant mode.
- **FR-UP02**: System MUST expose `POST /rest/private/update` for download/apply flows on update
  entries (web, launcher, secondary APK packages).
- **FR-UP03**: System MUST compare installed versions using files directory and application DAO data
  like Java `UpdateResource`.

### Functional Requirements — QR (public)

- **FR-QR01**: System MUST expose `GET /rest/public/qr/{id}` returning PNG QR image.
- **FR-QR02**: System MUST expose `GET /rest/public/qr/json/{id}` returning provisioning JSON.
- **FR-QR03**: System MUST resolve configuration by `qrCodeKey` and embed launcher APK URL and
  SHA-256 when main application version is configured.
- **FR-QR04**: Query parameters `deviceId`, `create`, `useId`, `group`, `size` MUST be supported
  per Java.

### Functional Requirements — Cross-cutting

- **FR-X01**: System MUST add permission constants and tests for `push_api`, `plugin_push_send`,
  `plugin_push_delete` in `platform/auth`.
- **FR-X02**: All JSON responses MUST use Headwind envelope (`status`, `message`, `data`).
- **FR-X03**: Android agent enrollment and periodic sync MUST work against Go-only backend for
  in-scope sync and notification paths.
- **FR-X04**: React enrollment QR page and push-from-devices MUST work without Java.
- **FR-X05**: Swagger MUST document Phase 7 endpoints after `make swagger`.
- **FR-X06**: Phase 7 row in `MIGRATION.md` moves from pending to **done** when parity criteria met.
- **FR-X07**: Module wiring MUST use feature flags `MODULE_SYNC_ENABLED`, `MODULE_PUSH_ENABLED`,
  `MODULE_NOTIFICATIONS_ENABLED`, `MODULE_UPDATES_ENABLED`, `MODULE_QRCODE_ENABLED` (default true in dev).
- **FR-X08**: Scaffold modules MUST NOT register placeholder route groups that return fake success;
  routes register only when handlers are implemented.

### Key Entities

- **Sync response**: Agent-facing bundle of device, configuration, applications, files, settings.
- **Device create options**: Enrollment payload for on-demand device/configuration assignment.
- **Device info**: Telemetry document stored on device row (battery, location, IMEI, customs).
- **Application setting**: Per-device per-package key/value with type and readonly flag.
- **Push message**: Queued command/message row keyed by device id, type, payload, delivery state.
- **Plain push message**: Agent API view of pending delivery items.
- **Plugin push message**: Tenant-scoped plugin history record with filter metadata.
- **Plugin push schedule**: Scheduled task for recurring plugin pushes.
- **Update entry**: Manifest row describing package, version, url, outdated/downloaded flags.
- **Update request**: Batch command to download, apply, or send stats for entries.
- **QR provisioning bundle**: Encoded enrollment JSON/extras for Android Device Owner setup.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An enrolled test agent (or curl simulation) receives a complete configuration payload
  from sync in under 3 seconds on a seeded database.
- **SC-002**: An administrator can open the enrollment QR page and load a valid PNG for a
  configuration with `qrCodeKey` without calling Java.
- **SC-003**: An administrator with `push_api` can send a push from the React devices UI and the
  target device can retrieve it via notifications GET within one polling cycle.
- **SC-004**: Super-admin can run Updates check and see manifest entries when vendor manifest URL is
  reachable (or documented stub returns empty list in offline dev).
- **SC-005**: 100% of in-scope endpoints listed in Java resources above are marked Done (or Partial
  with notes) in parity docs.
- **SC-006**: Phase 7 row in `MIGRATION.md` moves from **pending** to **done**.
- **SC-007**: Phase 6 file URLs referenced in sync payloads resolve using existing storage URL rules.

## Assumptions

- Phases 1–6 remain deployed and stable (devices, configurations, applications, files).
- Legacy Postgres schema exists or will gain push/notification/plugin tables via migration `000009`
  (exact number assigned at plan time).
- Agents continue to use HTTP sync and polling; native FCM delivery is **out of scope** for Phase 7.
- `PushService` background sender may be simplified to DB queue + polling wakeup (parity doc notes
  if push latency differs from Java).
- Java `SyncResponseHook` plugin extensions are not ported; hooks are no-op unless explicitly added
  later in Phase 8.
- Device location/battery **events** may update persistence without full event bus (partial noted).
- Update manifest fetch uses vendor URL from Java `UpdateSettings` with customer domain substitution.
- React Updates page may only call `check` initially; POST update can ship with parity even if UI
  defers apply button.
- Long-polling servlet path uses singular `notification` per Java (`/rest/notification/polling`).

## Dependencies

- **Requires**: Phase 4 `devices`, Phase 5 `configurations`/`applications`, Phase 6 `files` URLs in
  sync payloads.
- **Requires**: Phase 1 auth for private routes; `HASH_SECRET` from Phase 6 public API config.
- **Blocks**: Production agent cutover from Java; reliable remote commands from console.
- **Enables**: Phase 8 plugin platform with push plugin already partially migrated.

## Out of Scope (Phase 7)

- Phase 8 generic `plugins/*` framework (except `plugins/push` endpoints in this phase).
- Google FCM / APNs native push gateways.
- Full replacement of Java `EventService` streaming to external analytics.
- `DownloadFilesServlet` static file serving (Phase 6 partial; may remain platform follow-up).
- MQTT broker implementation (configuration may reference mqttWorker; HTTP polling still required).
- Angular-only Push plugin UI rewrite in React (API parity only).
- QR code styling/branding beyond Java parity.
