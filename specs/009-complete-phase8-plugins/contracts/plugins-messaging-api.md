# API Contract: Messaging plugin

**Base path**: `/rest/plugins/messaging`  
**Auth**: Bearer JWT  
**Java reference**: `com.hmdm.plugins.messaging.rest.MessagingResource`

| Method | Path | Permission | Purpose |
|--------|------|------------|---------|
| POST | `/private/search` | authenticated | Paginated message history |
| POST | `/private/send` | `plugin_messaging_send` | Send to devices/groups/broadcast |
| DELETE | `/{id}` | `plugin_messaging_delete` | Delete message row |
| GET | `/private/purge/{days}` | `plugin_messaging_delete` | Purge older than N days |
| GET | `/public/status/{id}/{status}` | none | Agent delivery status callback |

## POST `/private/send`

**Body** (`SendRequest`):

```json
{
  "message": "string",
  "deviceNumbers": ["hmdm-001"],
  "groups": [],
  "broadcast": false
}
```

**Behavior**: Insert `plugin_messaging_messages` per device; enqueue agent notification (via shared `MessageQueue`).

## POST `/private/search`

**Body** (`MessageFilter`): pagination + text/date filters per Java.

**Response**: `PaginatedData<Message>`.

## GET `/public/status/{id}/{status}`

**Path params**: message id, numeric status code.  
**Auth**: none (legacy agent callback).
