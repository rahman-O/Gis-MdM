# Data Model: Phase 4 — Devices & Groups

**Date**: 2026-05-20  
**Storage**: PostgreSQL (`000006_devices_groups_core.up.sql` + extensions)

## Entity: Group

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| name | varchar(100) | NOT NULL per customer |
| customerid | int | FK → customers; tenant scope |

**Uniqueness**: name unique per customer (application check).

## Entity: Device

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| number | varchar(100) | NOT NULL; unique per customer |
| description | text | optional |
| lastupdate | bigint | epoch ms; drives status color |
| configurationid | int | FK → configurations |
| customerid | int | FK → customers |
| info / infojson | text/json | device telemetry (optional v1) |
| imei, phone, model, … | varchar | from later Liquibase columns as needed |
| enrolltime, publicip | bigint/varchar | filters |
| custom1–custom3, oldnumber | varchar | metadata |
| fastsearch | boolean | search optimization flag |

**Status color (computed)**: green/yellow/red from `lastupdate` vs now (2h/4h thresholds per Java).

## Entity: DeviceGroup (junction)

| Field | Type | Rules |
|-------|------|-------|
| deviceid | int | FK → devices ON DELETE CASCADE |
| groupid | int | FK → groups ON DELETE CASCADE |

**Cardinality**: many-to-many devices ↔ groups.

## Entity: UserDeviceGroupsAccess

| Field | Type | Rules |
|-------|------|-------|
| userid | int | FK → users |
| groupid | int | FK → groups |

**Purpose**: Restrict device visibility when `users.alldevicesavailable = false`.

## Entity: Configuration (minimal Phase 4)

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| name | varchar | NOT NULL |
| customerid | int | tenant |
| permissive | boolean | optional in list view |
| mainappid | int | optional; launcher version join deferred |

## Entity: DeviceApplicationSetting (optional table)

| Field | Type | Rules |
|-------|------|-------|
| deviceid | int | |
| applicationpkg | varchar | |
| name, type, value | varchar | per-app override |

## DTO: DeviceSearchRequest

| Field | Notes |
|-------|-------|
| pageNum, pageSize | 1-based page |
| value | text search |
| groupId, configurationId | filters |
| status | online filter code |
| sortBy, sortDir | column sort |
| dateFrom, dateTo, enrollmentDateFrom/To | ranges |
| mdmMode, kioskMode, launcherVersion, … | optional filters |

## DTO: DeviceListView (API)

| Part | Content |
|------|---------|
| configurations | map[int]ConfigurationView |
| devices.items | []DeviceView |
| devices.totalItemsCount | int64 |

## DTO: GroupBulkRequest

| Field | Notes |
|-------|-------|
| ids | device ids |
| action | `set` or clear |
| groups | []LookupItem |

## Validation rules

| Rule | Error |
|------|-------|
| Duplicate device number | device exists envelope |
| Duplicate group name | `error.duplicate.group` |
| Delete non-empty group | `error.notempty.group` |
| No edit_devices on mutate | `error.permission.denied` |
| Device limit exceeded | generic ERROR (Java) |

## Relationships

```text
customers (1) ──< groups (many)
customers (1) ──< devices (many)
customers (1) ──< configurations (many)
devices (many) ──< devicegroups >── groups (many)
users (many) ──< userdevicegroupsaccess >── groups (many)
devices (1) ──< deviceapplicationsettings (many)
```

## Deferred schema (Phase 5+)

- `applications`, `applicationversions`, `configurationapplications`
- Full `infojson` parsing parity
- `deviceStatuses` installation tracking table
