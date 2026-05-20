# Feature Specification: Phase 5 — Applications, Configurations & Config Files

**Feature Branch**: `006-complete-phase5-apps-config`

**Created**: 2026-05-20

**Status**: Draft

**Input**: Complete **Phase 5** of the Java→Go MDM migration per
`serverBackendGo/docs/MIGRATION.md`: full **`applications`**, **`configurations`**, and
**`configfiles`** modules so tenant administrators can manage the application catalog and device
policy profiles (configurations) in the React MDM console without the Java WAR. Builds on Phase 4
(`devices`, `groups`, and read-only `GET /configurations/list`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Browse and open configurations (Priority: P1)

A tenant administrator opens the **Configurations** page, sees all policy profiles for their
organization, searches by name, and opens one configuration for editing.

**Why this priority**: `GET /private/configurations/search` and `GET /private/configurations/{id}`
are the entry points for the Configurations UI; Phase 4 only implemented `GET /list` for device
dropdowns.

**Independent Test**: User with `configurations` permission → search returns full configuration
rows; GET by id returns editor-ready object.

**Acceptance Scenarios**:

1. **Given** configurations exist for the tenant, **When** GET search runs,
   **Then** all tenant configurations are returned with fields needed by the list table.
2. **Given** a valid configuration id, **When** GET by id runs,
   **Then** full configuration detail is returned (name, type, design, applications, files metadata).
3. **Given** user without `configurations` permission, **When** search or detail is requested,
   **Then** permission denied matches Java.
4. **Given** unknown configuration id, **When** GET by id runs,
   **Then** not-found error envelope consistent with legacy.

---

### User Story 2 — Browse application catalog (Priority: P1)

An administrator opens the **Applications** page, lists Android and web applications for the
tenant, searches by name or package, and opens an application to manage versions.

**Why this priority**: `GET /private/applications/search` drives the Applications grid; without
it the app management workflow is blocked.

**Independent Test**: User with `applications` permission → search returns applications;
GET `/{id}` and GET `/{id}/versions` return detail for version management page.

**Acceptance Scenarios**:

1. **Given** applications for the tenant, **When** applications search runs,
   **Then** list matches Java tenant scope.
2. **Given** filter text, **When** search by value runs,
   **Then** filtered results return.
3. **Given** valid application id, **When** GET application and versions,
   **Then** version list is returned for the versions page.
4. **Given** no `applications` permission, **When** catalog access is attempted,
   **Then** permission denied.

---

### User Story 3 — Create, update, copy, and delete configurations (Priority: P2)

An administrator creates a new configuration (COMMON or WORK type), updates settings across
tabs, copies an existing profile, or deletes one no longer needed.

**Why this priority**: Configuration editor uses `PUT /private/configurations`, `PUT /copy`,
`DELETE /{id}`.

**Independent Test**: PUT create → appears in search; PUT update → persisted; copy → new row;
DELETE → removed; duplicate name → legacy error.

**Acceptance Scenarios**:

1. **Given** `configurations` permission and valid payload, **When** PUT without id (create),
   **Then** configuration is saved and returned.
2. **Given** existing configuration, **When** PUT with id (update),
   **Then** changes persist including linked applications and file definitions in payload.
3. **Given** copy request with source id and new name, **When** PUT copy,
   **Then** duplicate configuration exists for tenant.
4. **Given** configuration not assigned to devices (or Java delete rules), **When** DELETE,
   **Then** success; **Given** delete blocked by business rules,
   **Then** appropriate error envelope.
5. **Given** duplicate configuration name in tenant, **When** save,
   **Then** duplicate error matches legacy.

---

### User Story 4 — Manage Android and web applications (Priority: P2)

An administrator adds or edits an Android app (package, versions) or a web application URL app,
validates package uniqueness, and deletes apps or versions.

**Why this priority**: Applications page uses `PUT /android`, `PUT /web`, `PUT /versions`,
`DELETE /{id}`, `DELETE /versions/{id}`, `PUT /validatePkg`.

**Independent Test**: Save Android app → in search; add version → listed; delete version →
gone; validatePkg returns conflicts when duplicate package exists.

**Acceptance Scenarios**:

1. **Given** `applications` permission, **When** PUT android or web with valid payload,
   **Then** application is persisted.
2. **Given** new version payload, **When** PUT versions,
   **Then** version is stored and linked to application.
3. **Given** duplicate package in tenant, **When** validatePkg,
   **Then** conflicting applications are returned for UI warning.
4. **Given** delete application or version, **When** DELETE,
   **Then** removed per Java cascade rules.

---

### User Story 5 — Assign applications to configurations (Priority: P2)

An administrator configures which applications belong to a policy profile, upgrades app
versions on a configuration, and loads the application picker lists.

**Why this priority**: Configuration editor tabs call
`GET /configurations/applications/{id}`, `GET /configurations/applications`,
`PUT /configurations/application/upgrade`.

**Independent Test**: Load configuration applications → list; upgrade version → configuration
reflects new version; autocomplete for configuration names works.

**Acceptance Scenarios**:

1. **Given** configuration id, **When** GET applications for configuration,
   **Then** linked apps with version/status fields expected by React normalize layer.
2. **Given** upgrade payload (configuration + application + version),
   **When** PUT application/upgrade,
   **Then** configuration application version updates (push notify may be stubbed).
3. **Given** GET all applications for picker, **When** configurations/applications or fallback
   applications search,
   **Then** picker is populated for the editor.

---

### User Story 6 — Link configurations to applications (Priority: P2)

From the Applications UI, an administrator sees which configurations use an app (or version)
and updates those links in bulk.

**Why this priority**: `GET /applications/configurations/{id}`,
`POST /applications/configurations`, `GET /applications/version/{id}/configurations`,
`POST /applications/version/configurations`.

**Independent Test**: Open link dialog → configurations listed; save links → persisted and
visible on reload.

**Acceptance Scenarios**:

1. **Given** application id, **When** GET configurations for application,
   **Then** link rows with selection state return.
2. **Given** link update payload, **When** POST configurations,
   **Then** junction data updated for tenant.
3. **Given** application version id, **When** version configuration endpoints used,
   **Then** same behavior at version granularity.

---

### User Story 7 — Super-admin shared application catalog (Priority: P3)

A super administrator manages the cross-tenant **shared** application catalog and can merge
duplicate packages into a common app.

**Why this priority**: Admin Applications page uses `GET /applications/admin/search`,
`GET /applications/admin/common/{id}`.

**Independent Test**: Super-admin → admin search returns shared catalog; turn into common →
success envelope.

**Acceptance Scenarios**:

1. **Given** super-admin principal, **When** admin search runs,
   **Then** shared/common applications list returns.
2. **Given** non-super-admin, **When** admin endpoints called,
   **Then** permission denied.
3. **Given** application id eligible for merge, **When** GET admin/common/{id},
   **Then** operation completes per Java rules.

---

### User Story 8 — Upload configuration file assets (Priority: P3)

An administrator uploads a file used in a configuration (e.g. certificate, payload) through
the server upload endpoint so the configuration save can reference the stored path.

**Why this priority**: `ConfigurationFileResource` at `/private/config-files` supports binary
upload for configuration editor; may be invoked when file upload UX is wired.

**Independent Test**: Authenticated upload → file metadata returned; invalid type/size →
error; tenant isolation on storage path.

**Acceptance Scenarios**:

1. **Given** valid multipart upload, **When** POST config-files,
   **Then** upload result includes path/url fields expected by Java.
2. **Given** user without appropriate permission, **When** upload attempted,
   **Then** denied.
3. **Given** storage limit exceeded (if enforced), **When** upload,
   **Then** error envelope matches legacy semantics.

---

### User Story 9 — Configuration autocomplete and name list (Priority: P3)

Pickers across the UI load configuration names via list and autocomplete without loading full
configuration bodies.

**Why this priority**: Extends Phase 4 `GET /list`; adds `POST /autocomplete` and
`GET /search/{value}` for parity.

**Independent Test**: GET list → id/name; POST autocomplete with filter → suggestions; search by
value → subset.

**Acceptance Scenarios**:

1. **Given** tenant configurations, **When** GET list (Phase 4 endpoint),
   **Then** still works unchanged.
2. **Given** filter string, **When** autocomplete POST,
   **Then** lookup items return.
3. **Given** search value path, **When** GET search/{value},
   **Then** filtered configurations return for users with permission.

---

### User Story 10 — Verifiable API and regression safety (Priority: P2)

Developers exercise Phase 5 endpoints via Swagger and automated tests without Java.

**Independent Test**: Bearer auth smoke paths; `go test` for applications/configurations modules.

**Acceptance Scenarios**:

1. **Given** regenerated Swagger, **When** browsing UI,
   **Then** Applications and Configurations tags list in-scope endpoints.
2. **Given** module tests, **When** run locally/CI,
   **Then** permission, CRUD, and link operations are covered.

---

### Edge Cases

- Cross-tenant access to configuration or application by id → denied or not found.
- Delete configuration still used by devices → error or blocked per Java (`error.notempty` or equivalent).
- Delete application referenced by configurations → cascade or block per Java rules.
- Configuration save with empty application list → valid.
- Super-admin operating on tenant data vs shared catalog → correct scope per endpoint.
- Push notification on configuration upgrade → stub OK (no agent delivery required in Phase 5).
- Large APK/binary upload via `/private/web-ui-files` → **out of scope** (Phase 6 `files` module);
  Phase 5 may persist application metadata referencing pre-uploaded paths when provided.
- Device search configuration enrichment (nested apps in device grid) → optional partial extension
  of Phase 4 devices search.

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Modules**: Replace scaffolds in `internal/modules/applications/`,
  `internal/modules/configurations/` (extend beyond list-only), `internal/modules/configfiles/`.
- **Phase**: 5 in `MIGRATION.md`; marks Phase 5 **done** when modules ship for React in-scope APIs.
- **Java reference**: `ApplicationResource`, `ConfigurationResource`, `ConfigurationFileResource`;
  DAOs `ApplicationDAO`, `ConfigurationDAO`, related mappers.
- **REST bases**:
  - `/rest/private/applications` — ApplicationResource paths used by React
  - `/rest/private/configurations` — full ConfigurationResource (including existing `/list`)
  - `/rest/private/config-files` — ConfigurationFileResource upload
- **Parity docs**: `docs/parity/applications.md`, `docs/parity/configurations.md`,
  `docs/parity/configfiles.md` (or config-files).
- **Layers**: domain/port/application/adapter per module; permissions `applications`,
  `configurations` via `platformauth.HasPermission`.
- **Migrations**: Additive SQL for `applications`, `applicationversions`,
  `configurationapplications`, `configurationfiles`, `configurationapplicationsettings`, and
  extended `configurations` columns beyond Phase 4 minimal schema.
- **Auth**: JWT/session on private routes; strict tenant scoping; super-admin rules for admin/* paths.

### Functional Requirements — Configurations

- **FR-C01**: System MUST expose `GET /rest/private/configurations/search` with `configurations`
  permission.
- **FR-C02**: System MUST expose `GET /rest/private/configurations/{id}` for editor load.
- **FR-C03**: System MUST expose `PUT /rest/private/configurations` for create and update (React
  uses single PUT endpoint).
- **FR-C04**: System MUST expose `DELETE /rest/private/configurations/{id}` with business rules.
- **FR-C05**: System MUST expose `PUT /rest/private/configurations/copy`.
- **FR-C06**: System MUST expose `GET /rest/private/configurations/applications`,
  `GET /rest/private/configurations/applications/{id}`, and
  `PUT /rest/private/configurations/application/upgrade`.
- **FR-C07**: System MUST expose `POST /rest/private/configurations/autocomplete` and
  `GET /rest/private/configurations/search/{value}` for parity.
- **FR-C08**: System MUST retain `GET /rest/private/configurations/list` (Phase 4) behavior.

### Functional Requirements — Applications

- **FR-A01**: System MUST expose `GET /rest/private/applications/search` and
  `GET /rest/private/applications/search/{value}` with `applications` permission.
- **FR-A02**: System MUST expose `POST /rest/private/applications/autocomplete`.
- **FR-A03**: System MUST expose `GET /rest/private/applications/{id}` and
  `GET /rest/private/applications/{id}/versions`.
- **FR-A04**: System MUST expose `PUT /rest/private/applications/android`,
  `PUT /rest/private/applications/web`, `PUT /rest/private/applications/versions`.
- **FR-A05**: System MUST expose `DELETE /rest/private/applications/{id}` and
  `DELETE /rest/private/applications/versions/{id}`.
- **FR-A06**: System MUST expose `PUT /rest/private/applications/validatePkg`.
- **FR-A07**: System MUST expose application–configuration link endpoints:
  `GET/POST .../configurations`, `GET/POST .../version/.../configurations`.
- **FR-A08**: System MUST expose super-admin `GET /rest/private/applications/admin/search`,
  `GET .../admin/search/{value}`, `GET .../admin/common/{id}` for super-admin only.

### Functional Requirements — Config files

- **FR-F01**: System MUST expose `POST /rest/private/config-files` multipart upload returning
  legacy upload result shape.

### Functional Requirements — Cross-cutting

- **FR-X01**: System MUST add database migrations for applications/configurations junction tables.
- **FR-X02**: All responses MUST use Headwind envelope (`status`, `message`, `data`).
- **FR-X03**: React Configurations and Applications pages MUST work against Go-only backend for
  calls in `configurationService.ts` and `applicationService.ts` (except endpoints explicitly
  deferred to Phase 6).
- **FR-X04**: Swagger MUST document Phase 5 endpoints with Bearer auth after `make swagger`.
- **FR-X05**: Automated tests MUST cover configuration permission, application search, and one
  CRUD path per module.
- **FR-X06**: Phase 5 row in `MIGRATION.md` moves from pending to **done** when parity criteria met.

### Key Entities

- **Application**: Tenant or shared catalog app (name, package, type android/web, versions).
- **Application version**: Version string, APK/url reference, link to application.
- **Configuration**: Policy profile (type COMMON/WORK, design, MDM settings, QR-related fields).
- **Configuration application link**: Which app version is deployed on a configuration.
- **Configuration file**: Path/url metadata attached to a configuration.
- **Uploaded file**: Server-stored artifact referenced by configuration or application save.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An administrator can open the Configurations list and load one configuration for
  editing in under 5 seconds on a seeded tenant database.
- **SC-002**: An administrator can open the Applications list and drill into versions without
  calling Java.
- **SC-003**: Create, update, and delete flows for configurations and applications complete
  successfully for a tenant admin with the correct permissions.
- **SC-004**: 100% of in-scope `ConfigurationResource` and `ApplicationResource` endpoints used
  by React are marked Done (or Partial with notes) in parity docs.
- **SC-005**: Phase 5 row in `MIGRATION.md` moves from pending to **done**.
- **SC-006**: Devices page configuration dropdown continues to work via `GET /configurations/list`
  without regression.

## Assumptions

- Phases 1–4 (auth, users/roles, customers/settings, devices/groups) remain available.
- React `configurationService.ts` and `applicationService.ts` are the authoritative API contract.
- Binary uploads for large APKs via `/private/web-ui-files` remain **Phase 6**; Phase 5 accepts
  metadata-only saves when file path is already known.
- Push notification delivery on configuration upgrade is stubbed; endpoint returns success.
- QR code generation and base URL injection follow Java `baseUrl` behavior where applicable.
- Super-admin shared catalog endpoints are only required for super-admin principals.

## Dependencies

- **Requires**: Phase 4 `devices` and minimal `configurations` table/list endpoint.
- **Blocks**: Rich device search enrichment (optional), customer default devices on create (Phase 3 partial).
- **Enables**: Phase 6 files/icons (binary storage), Phase 7 agent sync consuming configurations.

## Out of Scope (Phase 5)

- Phase 6–8 modules (`files`, `icons`, `publicapi`, `sync`, `push`, `plugins`, …) except noting
  `/private/web-ui-files` dependency for APK upload UX.
- Public provisioning / QR enrollment APIs.
- Real FCM/APNs push infrastructure.
- Plugin-specific application hooks.
- Full maps/geofencing and plugin audit for configuration changes.
