# API Contract: Push messaging

---

## Private Push API (React)

**Base path**: `/rest/private/push`  
**Auth**: Bearer JWT + session  
**Envelope**: Headwind standard

**Java reference**: `com.hmdm.rest.resource.PushApiResource`  
**React reference**: `frontend/src/features/push/pushService.ts`

### POST `/`

Send push/command to devices.

**Permissions**: `push_api`

**Body** (`PushRequest` / React `PushPayload`):

```json
{
  "messageType": "string",
  "payload": "string",
  "deviceNumbers": ["hmdm-001"],
  "groups": [1],
  "broadcast": false
}
```

**Response data**: `OK` (empty or confirmation per Java)

**Errors**: `PERMISSION_DENIED`, device not found, empty targeting

---

## Push plugin (Angular legacy)

**Base path**: `/rest/plugins/push`  
**Auth**: Bearer JWT  
**Java reference**: `com.hmdm.plugins.push.rest.PushResource`

| Method | Path | Permission | Purpose |
|--------|------|------------|---------|
| POST | `/private/search` | authenticated | Paginated message history |
| POST | `/private/send` | `plugin_push_send` | Send to device/group/configuration scope |
| DELETE | `/private/{id}` | `plugin_push_delete` | Delete history row |
| GET | `/private/purge/{days}` | `plugin_push_delete` | Purge old messages |
| POST | `/private/searchTasks` | authenticated | List scheduled tasks |
| PUT | `/private/task` | `plugin_push_send` | Save schedule |
| DELETE | `/private/task/{id}` | `plugin_push_delete` | Delete schedule |

**Send body**: `PushSendRequest` (`scope`: `device` | `group` | `configuration`, `deviceNumber`, `groupId`, `messageType`, `payload`, …)

**Search body**: `PushMessageFilter` / `PushScheduleFilter` with pagination fields
