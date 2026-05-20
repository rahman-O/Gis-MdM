# API Contract: Push plugin — schedule tasks (Phase 8 completion)

**Base path**: `/rest/plugins/push`  
**Auth**: Bearer JWT  
**Java reference**: `com.hmdm.plugins.push.rest.PushResource` (task section)  
**Note**: Message endpoints implemented in Phase 7 — this contract covers **remaining** task CRUD only.

| Method | Path | Permission (Java) | Purpose |
|--------|------|-----------------|---------|
| POST | `/private/searchTasks` | authenticated | Paginated schedule list |
| PUT | `/private/task` | `plugin_push_delete` | Create/update schedule row |
| DELETE | `/private/task/{id}` | `plugin_push_delete` | Delete schedule |

## POST `/private/searchTasks`

**Body** (`PushScheduleFilter`): pagination + scope filters.

**Response**: `PaginatedData<PluginPushSchedule>`.

## PUT `/private/task`

**Body** (`PluginPushSchedule`): id optional; scope `device`|`group`|`configuration`; `deviceNumber` resolved to `deviceId` when scope=device.

**Response**: `OK`

## DELETE `/private/task/{id}`

**Response**: `OK`

**Partial**: Cron execution of due schedules not implemented (storage + admin API only).
