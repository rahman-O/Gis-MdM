# Data Model: Phase 9 — Migration Gap Completion

**Branch**: `011-complete-migration-gaps` | **Date**: 2026-05-21

Entities below are **new or extended** for gap closure. Existing Phase 1–8 tables unchanged unless noted.

## Push (reuse + behavior)

### `pushmessages` / `pendingpushes` (existing)

| Field | Usage in Phase 9 |
|-------|------------------|
| messagetype | `configUpdated`, `appConfigUpdated`, plugin types from schedule |
| deviceid | Target device |
| payload | Optional JSON string |

**Flow**: `platform/push.Notifier` → `MessageQueue.Enqueue` → agent polls `GET /rest/notifications/device/{deviceNumber}`.

### `plugin_push_schedule` (existing)

| Column | Notes |
|--------|-------|
| id | PK |
| customerId | Tenant |
| scope | `device` \| `group` \| `configuration` |
| deviceId, groupId, configurationId | Scope targets |
| messageType, payload | Enqueued per device |
| scheduledTime | Due timestamp |
| processed | Worker marks after send |

**State**: pending → processed (add column if Java uses status flag — verify mapper).

---

## `uploadedfiles` (existing, Phase 8 migration)

Used by **icon-files** upload.

| Column | Notes |
|--------|-------|
| id | Returned to client |
| customerid | Tenant |
| filepath | Filename under customer files dir |
| uploadtime | ms epoch |
| devicepath, external, replacevariables | Legacy fields |

**Validation**: Customer `sizeLimit` sum of file sizes (Phase E).

---

## `usagestats` (new migration `000011`)

| Column | Type | Notes |
|--------|------|-------|
| id | serial PK | `usagestats_id_seq` |
| *(fields from Java UsageStats domain)* | | Match `com.hmdm.persistence.domain.UsageStats` |

**Action**: Copy column list from Java domain during implementation; typical: customerId, deviceId, metric keys, timestamp.

---

## Audit log (existing `plugin_audit_log`)

Middleware writes same shape as manual search rows:

| Column | Source |
|--------|--------|
| userId, login | Principal |
| customerId | Principal |
| action | `HTTP_METHOD path` or mapped code |
| payload | Truncated body/query summary |
| ipAddress | `X-Forwarded-For` or client IP |
| createTime | ms epoch |

---

## Deviceinfo export (optional tables)

If export requires GPS/WiFi child tables not in dev DB, migration adds:

- `plugin_deviceinfo_deviceParams_*` per Java liquibase `deviceinfo.changelog.xml`

**Relationship**: `deviceParams` → device by `deviceNumber` / `deviceId`.

---

## Devicelog export

Uses existing `plugin_devicelog_*` tables from Phase 8.

**Export**: Query by filter → CSV/stream; no new tables unless `rules` cache table missing.

---

## Tenant bootstrap (customers)

Logical entities created on customer `PUT` (create):

| Entity | Notes |
|--------|-------|
| Customer | Row in `customers` |
| Default configuration | Copy from template customer id=1 or SQL seed |
| Optional default device | Java copies demo device — mirror `CustomerResource` |

**No new tables** — copy operations on existing `configurations`, `devices`.

---

## Sync hook extension (runtime only)

Not persisted — `SyncResponse` JSON map extended with plugin-provided keys per device sync.

---

## Videos (filesystem + optional metadata)

| Store | Notes |
|-------|-------|
| Filesystem | `{VIDEO_DIRECTORY}/{fileName}` |
| DB | None in Java `VideosResource` — file-only |

---

## Environment keys (configuration entity)

| Env | Module |
|-----|--------|
| `PUSH_SCHEDULE_INTERVAL_SEC` | scheduler |
| `MODULE_PUSH_NOTIFIER_ENABLED` | wire real vs noop notifier |
| `VIDEO_DIRECTORY` | videos |
| `MODULE_STATS_ENABLED` | stats |
| `MODULE_VIDEOS_ENABLED` | videos |
