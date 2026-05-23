# Research: Configuration–Device Sync & Admin UX

**Feature**: `016-config-sync-ux` | **Date**: 2026-05-23

## R1 — Root cause: device does not receive full configuration policy

**Decision**: Extend `BuildSyncResponse` to map all agent-relevant fields from `configurations` columns + `settingsjson` into `domain.SyncResponse`, matching Java `SyncResponse` and `ConfigurationService` assembly.

**Rationale**: Go `SyncResponse` today exposes only `password`, colors, `permissive`, `pushOptions`, `requestUpdates`, `applications`, `files`, and per-device `deviceapplicationsettings`. Java includes `kioskMode`, `restrictions`, GPS/Wi‑Fi/bluetooth flags, kiosk sub-options, `systemUpdateType`, `appPermissions`, etc. (see `backend/server/src/main/java/com/hmdm/rest/json/SyncResponse.java`). Agents ignore policy that never appears in JSON — explains “settings saved in UI but not on device.”

**Alternatives considered**:
- *Frontend-only fix* — rejected; agent reads sync API, not admin UI.
- *Second sync endpoint* — rejected; breaks parity and launcher expectations.

---

## R2 — Configuration save round-trip (editor ↔ DB)

**Decision**: Keep split persistence (SQL columns + `settingsjson` policy map) per 014; add golden tests for critical MDM keys (`kioskMode`, `restrictions`, `mainAppId`, `applications[]`).

**Rationale**: `configuration_json.go` already merges policy; gaps are missing keys in mapper or dropped nested arrays on save.

**Alternatives considered**:
- *Single JSON column only* — rejected; breaks legacy schema and migrations.

---

## R3 — “Restricted at configuration” semantics

**Decision**: Phase 1 locks use two mechanisms already in Headwind:
1. **Policy field locks** stored in `settingsjson` under `policyLocks` (map of field key → `true`) for scalar MDM fields (`mainAppId`, `kioskMode`, `restrictions`, …).
2. **Application setting readonly** via existing `configurationapplicationsettings.readonly` (and `variable` where used), merged into sync `applicationSettings` with `readonly: true`; device POST to `/sync/applicationSettings` must not upsert rows that match a readonly configuration default.

**Rationale**: Java already has `readonly` on configuration application settings (`ConfigurationMapper.xml`). Explicit `policyLocks` avoids overloading unrelated columns and is easy to render in React (lock icon per field).

**Alternatives considered**:
- *New `configuration_field_locks` table* — deferred; JSON map is enough for v1.
- *Block entire device settings tab* — rejected; too coarse for admins.

---

## R4 — Notify devices after configuration save

**Decision**: Wire `configurations` module to real `push.PushNotifier` (already stubbed in `module.go`) so `NotifyConfigurationChanged` enqueues `configUpdated` for devices on that configuration.

**Rationale**: `Service.Save` already calls `push.NotifyConfigurationChanged`; noop notifier means devices wait for polling only.

**Alternatives considered**:
- *Manual “Notify devices” only* — kept as optional UI later; auto-notify on save is default Headwind behavior.

---

## R5 — Admin UX: tabs and copy reduction

**Decision**: Consolidate `ConfigurationEditorPage` to tab components only (`ConfigurationCommonTab`, dedicated `ConfigurationMdmTab`, `ConfigurationRestrictionsTab`, existing design/apps/files tabs). Remove inline duplicate MDM block and long `CardDescription` paragraphs; use `aria-label` + single-line hints.

**Rationale**: Page already imports tab components but still embeds a partial MDM section (line ~285 “MDM block phase 1”), causing duplicate/confusing UX.

**Alternatives considered**:
- *New route per tab* — rejected; breaks deep links; single editor with tab state is enough.

---

## R6 — BASE_URL and application URLs

**Decision**: No change in this feature beyond documenting dependency on stable `BASE_URL` (015); sync must rewrite app/file URLs with `storage.BuildPublicURL` when stored URLs are localhost or stale tunnel hosts.

**Rationale**: Enrollment already addressed; sync failures from bad APK URLs are operational, not configuration schema.

---

## R7 — Parity verification method

**Decision**: Add contract tests: fixture configuration in DB → `GET /rest/public/sync/configuration/{device}` (signed) compared to Java sample JSON for same seed (field subset: applications, kioskMode, restrictions, files).

**Rationale**: Meets spec SC-001/SC-005 without manual-only QA.

**Alternatives considered**:
- *Manual Android only* — kept in quickstart as UAT, not sole gate.
