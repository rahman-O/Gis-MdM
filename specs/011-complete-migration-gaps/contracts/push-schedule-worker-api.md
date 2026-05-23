# API Contract: Push schedule worker (Phase 9 P0)

**Java reference**: `com.hmdm.plugins.push.guice.module.PushScheduleTaskModule`

## Runtime worker

| Property | Value |
|----------|-------|
| Interval | `PUSH_SCHEDULE_INTERVAL_SEC` (default 60) |
| Enabled | `MODULE_PLUGINS_PUSH_ENABLED` && plugin `push` in `ENABLED_PLUGINS` |
| Entry | `app.scheduler` → `plugins/push/application.ScheduleRunner` |

## Processing algorithm

1. `SELECT` due rows from `plugin_push_schedule` (match Java `findMatchingTime()`).
2. For each task, resolve device list by `scope`:
   - `device` → single `deviceId`
   - `group` → devices in group
   - `configuration` → devices with `configurationId`
3. Enqueue `messageType` + `payload` per device via `MessageQueue`.
4. Mark task processed / update timestamp (match Java post-send state).

## Admin API (existing — no path change)

| Method | Path | Status after Phase 9 |
|--------|------|----------------------|
| POST | `/rest/plugins/push/private/searchTasks` | Done |
| PUT | `/rest/plugins/push/private/task` | Done |
| DELETE | `/rest/plugins/push/private/task/{id}` | Done |

**Partial flag removed** from `docs/parity/push.md` when worker ships.

## Smoke

1. `PUT /rest/plugins/push/private/task` with `scheduledTime` = now + 90s.
2. Wait 2 minutes.
3. `GET /rest/notifications/device/{deviceNumber}` shows new pending message OR `pushmessages` row exists.
