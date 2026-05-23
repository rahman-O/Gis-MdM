# Data Model: إكمال تكامل React ↔ Go (014)

**Branch**: `014-complete-frontend-go-integration` | **Date**: 2026-05-21

**Note**: لا migrations جديدة في 014 — يعتمد على schema **013** (`000011`–`000017`). هذا المستند يصف **عقود DTO** بين React و Go.

---

## Existing tables (used, not created in 014)

| Table | Feature use |
|-------|-------------|
| `settings` | Tenant fields from `000015` |
| `configurations.settingsjson` | MDM policy keys |
| `configurationapplicationparameters` | `skipVersionCheck` |
| `configurationapplications` | `remove`, `longtap` |
| `devicestatuses` | Updated from sync (P2) |
| `usagestats` | Stats PUT (P2) |
| `uploadedfiles` + `icons` | Icon upload flow |

---

## Tenant Settings (extended DTO)

**Go**: `internal/modules/settings/domain/settings.go`  
**React**: `Settings` / `SettingsPayload` in [`frontend/src/features/settings/types.ts`](../../frontend/src/features/settings/types.ts)

| Field (JSON camelCase) | SQL column (`000015`) | Notes |
|------------------------|----------------------|--------|
| `newDeviceGroupId` | `newdevicegroupid` | FK `groups`, nullable |
| `phoneNumberFormat` | `phonenumberformat` | default `'+9 (999) 999-99-99'` |
| `customPropertyName1` | `custompropertyname1` | labels for device columns |
| `customPropertyName2` | `custompropertyname2` | |
| `customPropertyName3` | `custompropertyname3` | |
| `customMultiline1` | `custommultiline1` | boolean |
| `customMultiline2` | `custommultiline2` | |
| `customMultiline3` | `custommultiline3` | |
| `customSend1` | `customsend1` | |
| `customSend2` | `customsend2` | |
| `customSend3` | `customsend3` | |
| `desktopHeaderTemplate` | `desktopheadertemplate` | text |
| `sendDescription` | `senddescription` | boolean |

Existing fields (`language`, `createNewDevices`, design colors, …) unchanged.

---

## Configuration policy (JSON + SQL)

**Storage split** (unchanged from 013):

| Layer | Fields |
|-------|--------|
| SQL columns | `name`, `description`, `type`, `password`, colors, `qrcodekey`, `baseurl`, `mainappid`, `contentappid`, `defaultfilepath`, `permissive` |
| `settingsjson` | All MDM toggles in [`Configuration` type](../../frontend/src/features/configurations/types.ts): `gps`, `wifi`, `kioskMode`, `systemUpdateType`, … |

**Round-trip rule**: On PUT, non-column struct fields → merged into `settingsjson`; column fields → SQL UPDATE.

---

## ConfigurationApplication (extended)

| Field | Source |
|-------|--------|
| `skipVersionCheck` | `configurationapplicationparameters.skipversioncheck` |
| `remove` | `configurationapplications.remove` |
| `longTap` | `configurationapplications.longtap` |
| (existing) | `action`, `showIcon`, `usedVersionId`, … |

---

## DeviceInstallStatus (sync-derived)

| Field | Values (Java parity) |
|-------|---------------------|
| `applicationsstatus` | `SUCCESS`, `FAILURE`, `VERSION_MISMATCH`, … |
| `configfilesstatus` | `SUCCESS`, `OTHER`, … |

**Derivation (P2)**: From agent `info` payload in `POST /public/sync/info` — aggregate worst/best status per Java `DeviceStatusService` rules (simplified: any FAILURE → FAILURE; any VERSION_MISMATCH without FAILURE → VERSION_MISMATCH; else SUCCESS).

---

## UsageStats (new public DTO)

**Java**: `com.hmdm.persistence.domain.UsageStats`  
**Endpoint**: `PUT /rest/public/stats`

| Field | Type | DB column |
|-------|------|-----------|
| `instanceId` | string | `instanceid` |
| `webVersion` | string | `webversion` |
| `community` | bool | `community` |
| `devicesTotal` | int | `devicestotal` |
| `devicesOnline` | int | `devicesonline` |
| `cpuTotal` | int | `cputotal` |
| `cpuUsed` | int | `cpuused` |
| `ramTotal` | int | `ramtotal` |
| `ramUsed` | int | `ramused` |
| `scheme` | string | `scheme` |
| `arch` | string | `arch` |
| `os` | string | `os` |
| `ts` | date (optional) | `ts` default CURRENT_DATE |

**Unique**: `(ts, instanceid)` — upsert on conflict.

---

## Icon upload flow

1. `POST /private/icon-files` multipart `file` → `{ fileId, url? }`
2. `PUT /private/icons` `{ name, fileId }` → `IconRow`

---

## Validation rules

- Settings: `phoneNumberFormat` non-empty when provided; group FK must exist if `newDeviceGroupId` set.
- Configuration: reject PUT that would wipe `settingsjson` to `{}` if previous non-empty (merge only).
- Icon: PNG/square validation per existing `UploadIconFile` (backend).
- Stats: numeric fields ≥ 0; `instanceId` optional but required for upsert key when multiple instances.
