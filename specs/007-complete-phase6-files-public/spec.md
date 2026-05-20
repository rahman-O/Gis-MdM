# Feature Specification: Phase 6 — Files, Icons & Public API

**Feature Branch**: `007-complete-phase6-files-public`

**Created**: 2026-05-21

**Status**: Draft

**Input**: Complete **Phase 6** of the Java→Go MDM migration per
`serverBackendGo/docs/MIGRATION.md`: implement **`files`**, **`icons`**, and **`publicapi`**
modules with clean layered architecture, maintainable structure, and full API parity so tenant
administrators can manage the file library and launcher icons, upload APK/assets for applications,
and unauthenticated clients receive rebranding and AppList upload endpoints—without the Java WAR.
Builds on Phase 5 (`applications`, `configurations`, `configfiles`).

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Browse and manage the file library (Priority: P1)

A tenant administrator opens the **Files** page, sees all uploaded files for their organization,
searches by name, and removes files that are no longer needed.

**Why this priority**: `GET /private/web-ui-files/search` and `POST /remove` are the core Files
admin workflow; Phase 5 deferred binary library management to Phase 6.

**Independent Test**: User with `files` permission lists files; user with `edit_files` removes an
unused file; file in use returns legacy “file used” error.

**Acceptance Scenarios**:

1. **Given** uploaded files exist for the tenant, **When** GET search runs,
   **Then** file rows include metadata needed by the Files table (path, url, description, usage hints).
2. **Given** filter text, **When** GET search/{value} runs,
   **Then** filtered list returns.
3. **Given** file not referenced by configurations or icons, **When** POST remove,
   **Then** database row and on-disk file are removed (external URL files: DB only).
4. **Given** file still linked to configuration or icon, **When** POST remove,
   **Then** error envelope matches Java `FILE_USED` semantics.
5. **Given** user without `files` permission, **When** list/search is attempted,
   **Then** permission denied.
6. **Given** user without `edit_files` permission, **When** remove is attempted,
   **Then** permission denied.
7. **Given** unsafe file path in request, **When** remove or update runs,
   **Then** permission denied (path traversal guard).

---

### User Story 2 — Upload and commit files for applications and library (Priority: P1)

An administrator uploads an APK or other asset via multipart upload, receives a temporary server
path, then commits metadata so the file appears in the library and can be referenced when saving
applications or configurations.

**Why this priority**: Applications page (`webUiFilesService`) depends on
`POST /private/web-ui-files`, `POST /raw`, and `POST /update`; without these, APK version upload
is blocked.

**Independent Test**: Multipart upload returns `FileUploadResult`; commit via update creates
`UploadedFile` row; storage limit endpoint reflects tenant quota when configured.

**Acceptance Scenarios**:

1. **Given** `edit_files` permission, **When** multipart POST (default) with a file,
   **Then** response includes server temp path and name; APK uploads include parsed metadata when
   parseable (package, version, arch flags) or safe defaults when parser unavailable.
2. **Given** non-APK file, **When** POST raw,
   **Then** upload succeeds without APK parsing.
3. **Given** valid commit payload after upload, **When** POST update (create),
   **Then** file is moved from temp to tenant directory and persisted with public URL.
4. **Given** existing file id, **When** POST update (metadata or content refresh),
   **Then** changes persist per Java rules.
5. **Given** duplicate path in tenant, **When** create,
   **Then** duplicate error matches legacy.
6. **Given** tenant with size limit in multi-tenant mode, **When** upload would exceed limit,
   **Then** `error.size.limit.exceeded` style error returns.
7. **Given** external file URL payload, **When** POST update create external,
   **Then** row saved without on-disk move.

---

### User Story 3 — File usage and configuration links (Priority: P2)

An administrator inspects which applications reference a file URL, views storage usage, and
(optionally) updates which configurations use a given uploaded file.

**Why this priority**: Supports operational safety before delete; aligns with Java
`GET /apps/{url}`, `GET /limit`, `GET/POST /configurations`.

**Independent Test**: GET apps by URL returns application list; GET limit returns used/limit MB;
configuration link endpoints behave when UI invokes them.

**Acceptance Scenarios**:

1. **Given** file URL, **When** GET apps/{url},
   **Then** applications using that URL return for tenant.
2. **Given** multi-tenant non-master customer with size limit, **When** GET limit,
   **Then** `sizeUsed` and `sizeLimit` reflect directory usage.
3. **Given** file id, **When** GET configurations/{id},
   **Then** configuration link rows return.
4. **Given** `edit_files` and valid link payload, **When** POST configurations,
   **Then** links update; devices may be notified (push delivery may be stubbed).

---

### User Story 4 — Manage launcher icons (Priority: P2)

An administrator lists custom icons, searches by name, creates or updates an icon linked to an
uploaded file, and deletes icons no longer needed.

**Why this priority**: Configuration editor and branding use `/private/icons`; React
`iconsService.ts` is the contract.

**Independent Test**: Search returns icons; PUT create/update returns icon; DELETE removes row;
delete requires `settings` permission per Java.

**Acceptance Scenarios**:

1. **Given** icons for tenant, **When** GET search or search/{value},
   **Then** icon list returns.
2. **Given** valid icon payload, **When** PUT without id (create) or with id (update),
   **Then** icon persisted and returned.
3. **Given** `settings` permission, **When** DELETE /{id},
   **Then** icon removed.
4. **Given** user without `settings` permission, **When** DELETE,
   **Then** permission denied.
5. **Given** icon references a file id, **When** file delete attempted elsewhere,
   **Then** file-used guard applies via icon linkage check.

---

### User Story 5 — Public rebranding and AppList upload (Priority: P2)

Unauthenticated clients load application branding on login/signup screens, and the AppList mobile
utility uploads applications to the server using device id + shared secret hash validation.

**Why this priority**: `PublicResource` at `/rest/public` serves login rebranding (`/name`, `/logo`)
and agent-side `/applications/upload` not used by React but required for parity and agents.

**Independent Test**: GET name returns JSON rebranding fields; GET logo returns image or redirect;
POST applications/upload with valid device hash succeeds; invalid hash rejected.

**Acceptance Scenarios**:

1. **Given** rebranding config, **When** GET /public/name,
   **Then** app name, vendor, signup, and terms links return.
2. **Given** logo file configured, **When** GET /public/logo,
   **Then** PNG stream or redirect to default logo matches Java.
3. **Given** valid device id and MD5(deviceId + secret) hash, **When** POST applications/upload,
   **Then** application row created for device customer; optional file stored under tenant files dir.
4. **Given** invalid or missing hash, **When** upload,
   **Then** error without creating application.
5. **Given** unknown device id, **When** upload,
   **Then** device-not-found error envelope.
6. **Given** duplicate package+version for customer, **When** upload,
   **Then** duplicate application error.

---

### User Story 6 — Verifiable API and regression safety (Priority: P2)

Developers exercise Phase 6 endpoints via Swagger and automated tests without Java.

**Independent Test**: Bearer auth smoke for private routes; public routes without auth; module
`go test` for permission, upload path safety, and icon CRUD.

**Acceptance Scenarios**:

1. **Given** regenerated Swagger, **When** browsing UI,
   **Then** Files, Icons, and Public API tags list in-scope endpoints.
2. **Given** module tests, **When** run locally/CI,
   **Then** permission guards, upload commit, and icon delete paths are covered.

---

### Edge Cases

- Cross-tenant access to uploaded file or icon by id → denied or not found.
- Commit upload with tmp path outside system temp directory → denied.
- Remove file with path containing `..` or absolute escape → denied.
- Master customer in multi-tenant mode → storage limit endpoint may return empty limits (Java behavior).
- APK with duplicate version code but different version name → upload blocked with legacy message.
- Split APK arch (armeabi/arm64) → exists/complete flags set per Java when version row exists.
- Push notify on file-configuration update → stub OK (no agent delivery required in Phase 6).
- Deprecated `PublicFilesResource` and broken Java `GET /web-ui-files/{filePath}` download → not
  required; agent downloads use servlet/static handler (documented partial if only admin API ships).
- Very large multipart uploads → respect configured max body size; friendly error on overflow.

## Requirements *(mandatory)*

### Constitution Constraints *(Gis-MdM backend features)*

- **Modules**: Replace scaffolds in `internal/modules/files/`, `internal/modules/icons/`,
  `internal/modules/publicapi/`.
- **Phase**: 6 in `MIGRATION.md`; marks Phase 6 **done** when modules ship for React and public
  contracts in scope.
- **Java reference**: `FilesResource`, `IconResource`, `PublicResource`; DAOs `UploadedFileDAO`,
  `IconDAO`, `ConfigurationFileDAO`, `ApplicationDAO`, `CustomerDAO`, `UnsecureDAO`; utilities
  `FileUtil`, `APKFileAnalyzer` (APK parse may be **partial** with parity notes).
- **REST bases**:
  - `/rest/private/web-ui-files` — FilesResource
  - `/rest/private/icons` — IconResource
  - `/rest/public` — PublicResource (`/name`, `/logo`, `/applications/upload`)
- **Parity docs**: `docs/parity/files.md`, `docs/parity/icons.md`, `docs/parity/publicapi.md`.
- **Layers**: `domain/`, `port/`, `application/`, `adapter/http/`, `adapter/persistence/postgres/`
  per module; shared file I/O via `internal/platform/` or dedicated port (no SQL in handlers).
- **Permissions**: `files`, `edit_files` (FilesResource); icon DELETE uses `settings` (Java);
  icon search/PUT follow Java (no extra permission on read/write in legacy).
- **Config**: `FILES_DIRECTORY`, `BASE_URL`, `HASH_SECRET`, rebranding env vars mirror Java
  `context.xml` / `.env.example`.
- **Migrations**: Additive SQL for `uploadedfiles`, `icons`, and indexes if absent in Postgres
  legacy schema; reuse existing `customers.filesdir`, `customers.sizelimit` columns.

### Functional Requirements — Files (`web-ui-files`)

- **FR-FL01**: System MUST expose `GET /rest/private/web-ui-files/search` and
  `GET /rest/private/web-ui-files/search/{value}` with `files` permission.
- **FR-FL02**: System MUST expose `POST /rest/private/web-ui-files/remove` with `edit_files`
  permission and path-safety checks.
- **FR-FL03**: System MUST expose `POST /rest/private/web-ui-files/update` for create, external
  create, and update flows.
- **FR-FL04**: System MUST expose multipart `POST /rest/private/web-ui-files` and
  `POST /rest/private/web-ui-files/raw` returning legacy `FileUploadResult` shape.
- **FR-FL05**: System MUST expose `GET /rest/private/web-ui-files/limit` for tenant storage usage.
- **FR-FL06**: System MUST expose `GET /rest/private/web-ui-files/apps/{url}` (URL-encoded).
- **FR-FL07**: System MUST expose `GET /rest/private/web-ui-files/configurations/{id}` and
  `POST /rest/private/web-ui-files/configurations` with tenant and user configuration scope checks.
- **FR-FL08**: System MUST store tenant files under `FILES_DIRECTORY` + customer `filesDir`
  subdirectory matching Java URL generation (`{baseUrl}/files/...`).
- **FR-FL09**: System MUST enforce `configurationFileDAO` / `iconDAO` “file in use” checks on delete.

### Functional Requirements — Icons

- **FR-IC01**: System MUST expose `GET /rest/private/icons/search` and
  `GET /rest/private/icons/search/{value}`.
- **FR-IC02**: System MUST expose `PUT /rest/private/icons` for create (no id) and update (with id).
- **FR-IC03**: System MUST expose `DELETE /rest/private/icons/{id}` requiring `settings` permission.
- **FR-IC04**: System MUST scope icons to current customer tenant.

### Functional Requirements — Public API

- **FR-PU01**: System MUST expose unauthenticated `GET /rest/public/name` returning rebranding JSON.
- **FR-PU02**: System MUST expose unauthenticated `GET /rest/public/logo` as image stream or redirect.
- **FR-PU03**: System MUST expose `POST /rest/public/applications/upload` (multipart: file + `app`
  JSON) with device hash validation using configured secret.
- **FR-PU04**: Public routes MUST NOT require session/JWT; rate limiting may reuse platform middleware
  where configured.

### Functional Requirements — Cross-cutting

- **FR-X01**: System MUST add permissions constants and tests for `files`, `edit_files` in
  `platform/auth`.
- **FR-X02**: All JSON responses MUST use Headwind envelope (`status`, `message`, `data`).
- **FR-X03**: React Files page, Applications upload flows, and Icons UI MUST work against Go-only
  backend for calls in `filesService.ts`, `webUiFilesService.ts`, and `iconsService.ts`.
- **FR-X04**: Swagger MUST document Phase 6 endpoints after `make swagger`.
- **FR-X05**: Automated tests MUST cover file permission denial, safe-path rejection, icon delete
  permission, and public hash rejection.
- **FR-X06**: Phase 6 row in `MIGRATION.md` moves from pending to **done** when parity criteria met.
- **FR-X07**: Module wiring MUST enable feature flags `MODULE_FILES_ENABLED`, `MODULE_ICONS_ENABLED`,
  `MODULE_PUBLICAPI_ENABLED` (default true in dev) consistent with other modules.

### Key Entities

- **Uploaded file**: Tenant-owned file metadata (path, name, url, external flag, customer id,
  optional subdir, upload time).
- **File view**: API list shape with usage flags (`usedByConfigurations`, `usedByIcons`).
- **Icon**: Named launcher icon referencing an uploaded file id.
- **File configuration link**: Association between uploaded file and configuration with optional notify flag.
- **File upload result**: Temp path, APK details, duplicate/exists flags for application upload UX.
- **Public upload request**: Device-scoped application metadata + optional binary for AppList utility.
- **Rebranding**: App name, vendor, links, logo path for unauthenticated branding endpoints.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An administrator can open the Files list and see tenant files in under 5 seconds on a
  seeded database.
- **SC-002**: An administrator can upload an APK through the Applications flow (multipart + commit)
  and save the application without calling Java.
- **SC-003**: An administrator can list, save, and delete icons from the Icons UI without Java.
- **SC-004**: 100% of in-scope `FilesResource`, `IconResource`, and `PublicResource` endpoints used
  by React or documented agent flows are marked Done (or Partial with notes) in parity docs.
- **SC-005**: Phase 6 row in `MIGRATION.md` moves from **pending** to **done**.
- **SC-006**: Phase 5 application/configuration flows that depended on deferred file upload no longer
  return 404 on `web-ui-files` endpoints.

## Assumptions

- Phases 1–5 (auth through applications/configurations) remain available.
- Legacy Postgres schema includes or will gain `uploadedfiles` and `icons` tables compatible with
  Java MyBatis mappers.
- `FILES_DIRECTORY` is writable in dev (`./data/files` via `scripts/dev.sh`).
- Full APK binary analysis matches Java `APKFileAnalyzer` where feasible; otherwise **partial**
  parity documents which fields are best-effort (package, version, versionCode).
- Push notification on file-configuration update is stubbed; endpoint returns success.
- `DownloadFilesServlet` / static `/files/*` download for enrolled devices may ship as platform
  middleware in Phase 6 or be documented **partial** if only admin REST API is in scope.
- React does not yet call `/rest/public/name` or `/logo`; endpoints are still required for parity
  and future login rebranding in React.
- `PublicFilesResource` (`/rest/public/files`) remains **out of scope** (deprecated in Java).

## Dependencies

- **Requires**: Phase 5 `applications`, `configurations`, `configfiles` (file references in saves).
- **Requires**: Phase 3 `customers` (per-tenant `filesDir`, `sizeLimit`).
- **Blocks**: Phase 7 agent `sync` consuming file URLs; QR enrollment flows using public assets.
- **Enables**: Complete application version upload UX; configuration file assets; launcher icon picker.

## Out of Scope (Phase 6)

- Phase 7–8 modules (`sync`, `push`, `notifications`, `updates`, `qrcode`, `plugins/*`).
- Deprecated `PublicFilesResource` REST download.
- Full FCM/APNs push delivery.
- Plugin-specific file hooks.
- QR code generation (Phase 7).
- Replacing `DownloadFilesServlet` unless required for same-origin file URLs in dev (may be follow-up).
