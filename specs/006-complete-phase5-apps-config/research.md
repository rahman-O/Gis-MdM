# Research: Phase 5 — Applications, Configurations & Config Files

**Date**: 2026-05-20

## R1 — Migration strategy

**Decision**: Add `000007_applications_configurations_core.up.sql` with Liquibase-aligned subset:
`applications` (customerid, pkg, name, type, common flag, …), `applicationversions`,
`configurationapplications` (configurationid, applicationid, applicationversionid, action,
showicon, screenorder, …), `configurationfiles`, minimal `configurationapplicationsettings`;
`ALTER configurations` for editor columns (type, password, design, qrcodekey, …).

**Rationale**: Phase 4 created minimal `configurations` (id, name, customerid, permissive,
mainappid). Applications tables absent in `000001`–`000006`.

**Alternatives considered**: Import full Liquibase — rejected (too large); piggyback on
`000006` — rejected (already deployed).

## R2 — Configurations module refactor

**Decision**: Move Phase 4 list handler into full clean-architecture module; keep `GET /list`
implementation path stable for Devices page regression (SC-006).

**Rationale**: Phase 4 used `adapter/http/handler.go` + inline repo; Phase 5 needs transactional
save with child rows.

## R3 — Configuration save payload

**Decision**: Accept full `Configuration` JSON on `PUT /private/configurations` matching React
`configurationNormalize.ts` output (nested `applications`, `files`, `applicationSettings`).

**Rationale**: Java `ConfigurationResource.updateConfiguration` persists graph in one request;
React never uses `POST` — only `PUT`.

## R4 — Configuration type field

**Decision**: Map UI `ConfigurationKind` WORK/COMMON to backend `type` int: `0` = WORK, `1` =
COMMON (per `configurationService.configurationKindToType`).

**Rationale**: Frontend and Java use numeric type consistently.

## R5 — Application tenant vs common catalog

**Decision**: `applications.customerid` nullable or `common` boolean; tenant apps filtered by
`customerid`; admin endpoints return cross-tenant shared apps for super-admin only.

**Rationale**: `ApplicationResource` admin/search and `common` flag in Java domain.

## R6 — Application–configuration links

**Decision**: Implement both directions:
- Config side: `configurationapplications` rows on configuration save + upgrade endpoint
- App side: `GET/POST /applications/configurations` and version-level variants

**Rationale**: React uses both `configurationService` and `applicationService` link APIs.

## R7 — validatePkg and duplicate package

**Decision**: `PUT /validatePkg` returns array of conflicting `Application` rows for same pkg in
tenant scope (read-only check before save).

**Rationale**: `ApplicationResource.validatePkg` used by duplicate package dialog.

## R8 — Config file upload

**Decision**: `POST /private/config-files` writes to `{files.directory}/{customer.filesdir}/{filename}`
with UTF-8 filename workaround; returns `FileUploadResult` envelope; overwrite existing file
with warn (match Java).

**Rationale**: `ConfigurationFileResource.uploadConfigurationFile`; no separate DB table
required for upload result beyond optional `uploadedfiles` if present in schema.

## R9 — Push on configuration upgrade

**Decision**: `PUT /application/upgrade` returns OK; call no-op `port.PushNotifier` (reuse
pattern from devices Phase 4).

**Rationale**: FCM/APNs is Phase 7+; UI expects success.

## R10 — web-ui-files / APK upload

**Decision**: **Out of scope** — `applicationService` and `webUiFilesService` use
`/private/web-ui-files` (Phase 6 `files` module). Phase 5 persists `url`/`filePath` on
application/version when client supplies path from prior upload or manual entry.

**Rationale**: `MIGRATION.md` Phase 6 owns `files`; avoids half-implemented storage.

## R11 — Permissions seed

**Decision**: Add `applications` and `configurations` permission rows; link to org-admin role (2)
like `edit_devices` in `000006`.

**Rationale**: Java checks `hasPermission("applications")` and `hasPermission("configurations")`;
super-admin bypasses via `HasPermission` for super-admin principal.

## R12 — Delete configuration guard

**Decision**: Block delete when `devices.configurationid` references configuration (count &gt; 0)
→ `error.notempty.configuration` or Java-equivalent message.

**Rationale**: Prevent orphan devices; match Java ConfigurationDAO rules.
