# API Contract: Notifications (agent delivery)

**Java reference**: `com.hmdm.notification.rest.NotificationResource`, `LongPollingServlet`

---

## JAX-RS pull

**Base path**: `/rest/notifications`  
**Auth**: None

### GET `/device/{deviceNumber}`

Fetch pending push messages for device.

**Response data**: `PlainPushMessage[]` (`id`, `messageType`, `payload`, …)

**Side effect**: Marks messages delivered per Java `NotificationDAO`

**Errors**: `DEVICE_NOT_FOUND`

---

## Long polling servlet

**Base path**: `/rest/notification/polling`  
**Auth**: None (signature when secure enrollment)

### GET `/{deviceNumber}`

Hold connection up to `POLLING_TIMEOUT` (env, default ~60s) until message available.

**Headers** (when `SECURE_ENROLLMENT=true`):
- `X-Request-Signature`: SHA1(`HASH_SECRET + deviceNumber`)

**Response**: JSON message list or empty on timeout (match servlet)

**Errors**: 403/permission denied on bad signature
