# Parity: Messaging plugin

**Java**: `com.hmdm.plugins.messaging.rest.MessagingResource`  
**Go**: `internal/modules/plugins/messaging`

| Method | Path | Permission |
|--------|------|------------|
| POST | `/rest/plugins/messaging/private/search` | authenticated |
| POST | `/rest/plugins/messaging/private/send` | `plugin_messaging_send` |
| DELETE | `/rest/plugins/messaging/{id}` | `plugin_messaging_delete` |
| GET | `/rest/plugins/messaging/private/purge/{days}` | `plugin_messaging_delete` |
| GET | `/rest/plugins/messaging/public/status/{id}/{status}` | public |

**Delivery**: Enqueues `textMessage` on shared notifications queue.
