# Data Model: Phase 8 — Plugins

**Branch**: `009-complete-phase8-plugins` | **Date**: 2026-05-21

Legacy table/column names follow Java liquibase (camelCase columns preserved in Postgres).

## Core platform

### `plugins`

| Column | Type | Notes |
|--------|------|-------|
| id | serial PK | |
| identifier | varchar(50) UNIQUE | e.g. `audit`, `push`, `messaging` |
| name | text | Display name |
| description | text | |
| createTime | timestamp | |
| disabled | boolean | Global disable |
| javascriptModuleFile | varchar(200) | UI metadata (not served by Go) |
| functionsViewTemplate | varchar(200) | |
| settingsViewTemplate | varchar(200) | |
| nameLocalizationKey | varchar(200) | React may use key |
| settingsPermission | varchar(200) | |
| functionsPermission | varchar(200) | |
| deviceFunctionsPermission | varchar(200) | |
| enabledForDevice | boolean | Agent-facing plugins |

### `pluginsDisabled`

| Column | Type | Notes |
|--------|------|-------|
| pluginId | int FK → plugins | Composite uniqueness with customer |
| customerId | int FK → customers | Tenant opt-out |

**Relationships**: Many disabled rows per customer; replaces full set on `POST /private/disabled`.

## Audit (`plugin_audit_log`)

| Column | Type | Notes |
|--------|------|-------|
| id | serial PK | |
| createTime | bigint | Epoch ms |
| customerId | int | Nullable in legacy |
| userId | int | |
| login | varchar(100) | |
| action | varchar(100) | |
| payload | text | Request summary |
| ipAddress | varchar(500) | |
| errorCode | int | Default 0 |

**Search filter**: date range, userId, login, action, pagination (`pageNum`, `pageSize`).

## Messaging (`plugin_messaging_messages`)

| Column | Type | Notes |
|--------|------|-------|
| id | serial PK | |
| customerId | int FK | |
| deviceId | int FK | |
| ts | bigint | Send timestamp |
| message | varchar(5000) | Body |
| status | int | Delivery status enum (Java ints) |

**Send flow**: Insert row per targeted device + enqueue `MessageQueue` for agent poll.

## Device info

### `plugin_deviceinfo_settings`

| Column | Type | Notes |
|--------|------|-------|
| id | serial PK | |
| customerId | int UNIQUE | One row per tenant |
| dataPreservePeriod | int | Days, default 30 |

### `plugin_deviceinfo_deviceParams` (+ child tables)

Stores structured device parameters (wifi, gps, mobile, etc.) linked to `devices`. Dynamic telemetry uses related tables per liquibase (`deviceParams_device`, `_wifi`, `_gps`, `_mobile`, …).

### Dynamic info (runtime)

`DeviceDynamicInfo`: attribute key/value/time series tied to `deviceId` — persisted via DAO upsert on public PUT.

## Device log (Postgres)

### `plugin_devicelog_settings`

Per-customer settings (retention, enabled flags per Java).

### `plugin_devicelog_settings_rules`

Log collection rules (level, prefix, group scope).

### `plugin_devicelog_setting_rule_devices`

Rule-to-device associations.

### `plugin_devicelog_log`

Uploaded log lines: deviceId, customerId, timestamp, level, message, package, etc. (see `devicelog.postgres.changelog.xml`).

## Push plugin (Phase 7 existing)

### `plugin_push_messages`

History of plugin-originated pushes (Phase 7).

### `plugin_push_schedule`

| Column | Type | Notes |
|--------|------|-------|
| id | serial PK | |
| customerId | int | |
| scope | varchar | `device` \| `group` \| `configuration` |
| deviceId | int | When scope=device |
| groupId | int | When scope=group |
| configurationId | int | When scope=configuration |
| messageType | varchar | |
| payload | text | |
| schedule | varchar/cron | Execution **deferred** (no cron runner Phase 8) |

## Permissions (seed in 000010)

| Permission | Used by |
|------------|---------|
| `plugins_customer_access_management` | POST `/plugin/main/private/disabled` |
| `plugin_audit_access` | Audit search |
| `plugin_messaging_send` | Messaging send |
| `plugin_messaging_delete` | Messaging delete/purge |
| `plugin_deviceinfo_access` | Deviceinfo private routes |
| `plugin_devicelog_access` | Devicelog private routes |
| `plugin_push_send` / `plugin_push_delete` | Phase 7 push plugin (existing) |

## Cross-module dependencies

```text
plugins/platform ──reads──► plugins, pluginsDisabled
plugins/messaging ──writes──► plugin_messaging_messages
                 ──uses───► notifications.MessageQueue
                 ──uses───► devices (target resolution)
plugins/audit ──reads──► plugin_audit_log
plugins/deviceinfo ──reads/writes──► plugin_deviceinfo_* + devices
plugins/devicelog ──reads/writes──► plugin_devicelog_* + devices
plugins/push ──reads/writes──► plugin_push_schedule (Phase 8 completion)
```

## Validation rules

- Disabled plugin: mutating endpoints return permission denied or plugin-disabled error per Java.
- All private queries scoped by `principal.CustomerID` (and impersonation).
- Messaging/push send: at least one target (device list, group, or broadcast flag).
- Devicelog upload: device must exist and belong to customer.
- Plugin schedule task: resolve `deviceNumber` → `deviceId` when scope=device.
