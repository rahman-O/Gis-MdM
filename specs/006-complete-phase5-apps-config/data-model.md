# Data Model: Phase 5 — Applications, Configurations & Config Files

**Date**: 2026-05-20  
**Storage**: PostgreSQL (`000007_applications_configurations_core.up.sql`)

## Entity: Application

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| pkg | varchar | NOT NULL; unique per tenant (validatePkg) |
| name | varchar | NOT NULL |
| customerid | int | FK → customers; NULL for shared/common catalog |
| type | varchar/int | android / web / intent discriminator |
| common | boolean | shared across tenants when true |
| showicon, system, … | various | per Liquibase subset |

**Relationships**: 1:N `applicationversions`; M:N `configurations` via `configurationapplications`.

## Entity: ApplicationVersion

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| applicationid | int | FK → applications |
| version | varchar | display version |
| versioncode | int | Android versionCode |
| url, urlarmeabi, urlarm64 | varchar | APK paths (may reference uploaded files) |
| filePath | varchar | optional server path |

## Entity: Configuration (extended)

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| name | varchar | unique per customer |
| description | text | optional |
| customerid | int | FK (from Phase 4) |
| type | int | 0 WORK, 1 COMMON |
| password, design fields | varchar/text | MDM launcher UI |
| mainappid, contentappid | int | launcher apps |
| qrcodekey, baseurl | varchar | provisioning display |
| permissive | boolean | from Phase 4 |

**Relationships**: 1:N `configurationapplications`, `configurationfiles`, `configurationapplicationsettings`; referenced by `devices.configurationid`.

## Entity: ConfigurationApplication (junction)

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| configurationid | int | FK |
| applicationid | int | FK |
| applicationversionid | int | FK nullable |
| action | int | install/remove/run rules |
| showicon | boolean | |
| screenorder, keycode, bottom | int/bool | launcher layout |

## Entity: ConfigurationFile

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| configurationid | int | FK |
| path | varchar | device path |
| externalurl, url | varchar | source URLs |
| remove | boolean | mark for removal on save |

## Entity: ConfigurationApplicationSetting

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| configurationid | int | FK |
| applicationid | int | optional |
| name, type, value | varchar | per-app override on configuration |

## DTO: Configuration (API)

Mirrors React `Configuration` in `frontend/src/features/configurations/types.ts` — large
optional surface; list view uses `id`, `name`, `description`, `type`, `deviceCount`.

## DTO: Application (API)

Mirrors React `Application` / `ApplicationVersion` in `applications/model/types.ts`.

## DTO: Link payloads

- `LinkConfigurationsToAppRequest` — app ↔ configs bulk update
- `LinkConfigurationsToAppVersionRequest` — version-level links
- `UpgradeConfigurationApplicationPayload` — `{ configurationId, applicationId }`

## Validation rules

| Rule | Error |
|------|-------|
| Duplicate configuration name | `error.duplicate.configuration` |
| Configuration in use by devices | `error.notempty.configuration` (or Java equivalent) |
| Missing `configurations` permission | `error.permission.denied` |
| Missing `applications` permission | `error.permission.denied` |
| Non-super-admin on admin routes | `error.permission.denied` |
| Duplicate package (validatePkg) | return conflicting apps in data (not always ERROR) |

## Relationships

```text
customers (1) ──< applications (many, tenant-scoped)
applications (1) ──< applicationversions (many)
customers (1) ──< configurations (many)
configurations (many) ──< configurationapplications >── applications
configurations (1) ──< configurationfiles (many)
devices (many) ──> configurations (1)
```

## Deferred schema (Phase 6+)

- `uploadedfiles` full parity
- `icons` table for `iconId` on applications
- `web-ui-files` storage quotas (`sizeLimit` / `sizeUsed`)
