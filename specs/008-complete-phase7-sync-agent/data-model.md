# Data Model: Phase 7 — Agent Sync, Push, Notifications, Updates & QR

**Branch**: `008-complete-phase7-sync-agent` | **Date**: 2026-05-21

## Overview

Phase 7 adds agent messaging tables and relies on existing Phase 4–6 entities (`devices`,
`configurations`, `applications`, `applicationversions`, `configurationfiles`, `uploadedfiles`,
`customers`, `groups`). Plugin push history uses separate tables from the core notification queue.

---

## New / migrated tables (000009)

### `pushmessages`

| Column | Type | Notes |
|--------|------|-------|
| id | SERIAL PK | |
| messagetype | VARCHAR(50) NOT NULL | Agent command type |
| deviceid | INT NOT NULL FK → devices(id) ON DELETE CASCADE | |
| payload | TEXT | JSON or plain string |

Java: `pushMessages` (folded to lowercase in Postgres).

### `pendingpushes`

| Column | Type | Notes |
|--------|------|-------|
| id | SERIAL PK | |
| messageid | INT NOT NULL UNIQUE FK → pushmessages(id) | |
| status | INT NOT NULL DEFAULT 0 | 0=pending, delivered states per Java |
| createtime | BIGINT NOT NULL | |
| sendtime | BIGINT | Set on delivery |

### `plugin_push_messages`

| Column | Type | Notes |
|--------|------|-------|
| id | SERIAL PK | |
| customerid | INT NOT NULL FK → customers | |
| deviceid | INT NOT NULL FK → devices | |
| ts | BIGINT NOT NULL | |
| messagetype | VARCHAR(255) | |
| payload | TEXT | |

### `plugin_push_schedule`

| Column | Type | Notes |
|--------|------|-------|
| id | SERIAL PK | |
| customerid | INT NOT NULL | |
| deviceid | INT DEFAULT 0 | |
| groupid | INT DEFAULT 0 | |
| configurationid | INT DEFAULT 0 | |
| scope | VARCHAR(255) | device / group / configuration |
| messagetype | VARCHAR(255) | |
| payload | TEXT | |
| comment | TEXT | |
| min, hour, day, weekday, month | VARCHAR(1024) | Cron-like fields |
| minbit, hourbit, daybit, weekdaybit, monthbit | BIT | Java BIT columns → BYTEA or simplified TEXT in Go migration |

**Note**: Schedule bit columns may use `BYTEA` placeholders in Go migration if BIT type complicates `lib/pq`;
parity doc marks schedule **partial** if cron execution deferred.

### Permissions (seed)

| name | Used by |
|------|---------|
| push_api | `POST /rest/private/push` |
| plugin_push_send | Plugin send |
| plugin_push_delete | Plugin delete/purge |

Assign to role id 2 (org admin) like prior phases.

---

## Existing entities (read/write)

### Device (`devices`)

- **Sync reads**: number, oldnumber, imei, serial, customerid, configurationid, info JSON, lastupdate,
  custom1–3, groups.
- **Sync writes**: info JSON, lastupdate, ip, imei migration timestamp, custom fields, application settings
  child rows (`deviceapplicationsettings` or equivalent per Java mapper).

### Configuration (`configurations`)

- **Sync reads**: full policy for assigned device; password; mainappid; launcherurl; adminextras; qrCodeKey.
- **QR reads**: by `qrCodeKey` (column `qrcodekey`).

### Application / version

- **Sync reads**: apps linked to configuration; version URLs, hashes, split flags.
- **Updates**: version comparison and upgrade paths.

### Customer (`customers`)

- **Sync**: single-customer vs multi-tenant creation rules; master flag for updates check.

---

## Domain DTOs (by module)

### sync

| DTO | Purpose |
|-----|---------|
| DeviceCreateOptions | POST enrollment body |
| DeviceInfo | POST /info telemetry |
| SyncApplicationSetting | POST /applicationSettings items |
| SyncResponse | Agent configuration payload (nested device, configuration, applications, files, settings) |

### notifications

| DTO | Purpose |
|-----|---------|
| PlainPushMessage | Agent-facing message (`messageType`, `payload`, ids) |

### push

| DTO | Purpose |
|-----|---------|
| PushRequest | React `PushPayload` / Java `PushRequest` |
| PushSendRequest | Plugin send scope |

### updates

| DTO | Purpose |
|-----|---------|
| UpdateEntry | Manifest row + flags |
| UpdateRequest | POST body with `updates[]`, `update`, `sendStats` |

### qrcode

| DTO | Purpose |
|-----|---------|
| QRQuery | deviceId, create, useId, group[], size |
| ProvisioningExtras | JSON string for Android |

---

## Relationships

```text
customers 1──* devices *──1 configurations
devices 1──* pushmessages 1──1 pendingpushes
customers 1──* plugin_push_messages *──1 devices
plugin_push_schedule *── customers (optional device/group/configuration scope)
configurations 1──* configurationapplications *──* applicationversions
```

---

## State transitions

### Pending push delivery

1. Admin POST `/private/push` → insert `pushmessages` + `pendingpushes` (status=0).
2. Agent GET `/notifications/device/{n}` or long-poll → return messages, mark delivered (update status/sendtime).
3. Plugin send also inserts `plugin_push_messages` for audit (in addition to queue).

### Device enrollment

1. Unknown device + create options → insert device, assign configuration/groups.
2. `preventDuplicateEnrollment` + `lastupdate>0` → DEVICE_EXISTS error.
3. Migration: match oldnumber or imei/serial → return config without duplicate create.

---

## Validation rules

- Sync signatures required when `SECURE_ENROLLMENT=true`.
- Push API requires `push_api`; plugin send requires `plugin_push_send`.
- Update check: super-admin only when multiple customers exist.
- QR: configuration must exist for key; main app version required for PNG with APK hash (Java behavior).
- Push payload and messageType non-empty for send endpoints.
