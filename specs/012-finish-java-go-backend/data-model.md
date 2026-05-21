# Data Model: إكمال نقل الباكند Java → Go (012)

**Branch**: `012-finish-java-go-backend` | **Date**: 2026-05-21  
**Extends**: [011 data-model](../011-complete-migration-gaps/data-model.md)

## Devices (extended)

### `devices` table (existing)

| Column | 012 usage |
|--------|-----------|
| `infojson` | JSONB/text telemetry; parse to `DeviceInfoView` on read |
| `lastupdate` | Status bands + sort `LAST_UPDATE` |
| `configurationid`, `customerid` | Filters + tenant scope |

### `SearchRequest` (domain extension)

| Field | Type | Filter behavior |
|-------|------|-----------------|
| pageNum, pageSize | int | Existing |
| value, groupId, configurationId, fastSearch | optional | Existing |
| status | string | Maps to lastupdate CASE |
| androidVersion, launcherVersion | string | `infojson ->>` |
| mdmMode, kioskMode | bool | `infojson ->>` |
| installationStatus | string | `infojson` apps array or status field |
| sortBy, sortDir | string | SQL ORDER BY |
| dateFrom, dateTo, … | int64 | `lastupdate` or enroll columns |

### `DeviceView` / `DeviceInfoView` (response)

Nested `info` object for React:

| Field | Source in infojson |
|-------|-------------------|
| batteryLevel, model, androidVersion | top-level keys |
| applications[], files[] | optional arrays |
| mdmMode, kioskMode, launcherVersion | keys per Java export mapper |

---

## Push (unchanged — regression)

See 011 data-model: `pushmessages`, `pendingpushes`, `plugin_push_schedule`.

---

## `uploadedfiles` + quota (extended)

| Rule | Implementation |
|------|----------------|
| Insert on config-files / web-ui-files upload | customerId, filepath, size |
| Quota check | `SUM(size)` vs `customers.sizeLimit` before accept |

---

## `usagestats` (new migration `000011`)

Align columns with Java `com.hmdm.persistence.domain.UsageStats` during implementation.

Typical keys: customerId, deviceId, metric identifiers, timestamp (verify Java DAO).

---

## Audit log (middleware writes)

Table: `plugin_audit_log` (existing Phase 8).

| Column | Source |
|--------|--------|
| userId, customerId | Principal |
| action | HTTP method + path template |
| details | Truncated request summary |
| createTime | epoch ms |

---

## Plugin deviceinfo / devicelog

No new tables unless export requires missing GPS/WiFi child tables (verify liquibase `deviceinfo.changelog.xml`).

**Export**: Read-only queries across `plugin_deviceinfo_*`, `plugin_devicelog_*`.

**Rules endpoint**: `plugin_devicelog_rule` per deviceNumber (existing schema).

---

## Tenant bootstrap (logical)

On customer **create**:

| Created | Source |
|---------|--------|
| Default configuration | Template row / Java copy logic |
| Optional default device | Java `CustomerResource` |
| Admin user | Existing customers service |

---

## Static files (no new table)

URL pattern: `{BASE_URL}/files/{customerFilesDir}/{relativePath}`  
Served from `FILES_DIRECTORY` on disk.

---

## Sync response extension

`SyncResponse` map extended at runtime by registered hooks (opaque JSON keys per plugin).

---

## Videos (optional module)

Filesystem: `VIDEO_DIRECTORY/{fileName}`  
No DB table if mirroring Java `VideosResource` file-only storage.
