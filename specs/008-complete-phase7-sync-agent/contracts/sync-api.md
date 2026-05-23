# API Contract: Device Sync (`/rest/public/sync`)

**Base path**: `/rest/public/sync`  
**Auth**: None (public agent API)  
**Envelope**: Headwind standard for JSON endpoints

**Java reference**: `com.hmdm.rest.resource.SyncResource`  
**Consumers**: Android MDM launcher/agent

---

### POST `/configuration/{deviceId}`

Enroll or refresh device configuration.

**Body**: `DeviceCreateOptions` (configuration id, groups, customer hints — match Java JSON)

**Headers** (when `SECURE_ENROLLMENT=true`):
- `X-Request-Signature`: SHA1 hex of `HASH_SECRET + deviceId`
- Optional: `X-CPU-Arch`, client IP via proxy headers

**Response data**: `SyncResponse` (device, configuration, applications, files, settings, branding)

**Response headers**:
- `X-Response-Signature`: SHA1(`HASH_SECRET` + compact JSON of data)
- `X-IP-Address`: resolved client IP

**Errors**: `DEVICE_NOT_FOUND`, `DEVICE_EXISTS`, `PERMISSION_DENIED` (bad signature), `INTERNAL_ERROR`

---

### GET `/configuration/{deviceId}`

Same as POST without create body; used for periodic sync.

**Headers**: Same signature rules when secure enrollment enabled.

---

### POST `/info`

Update device telemetry.

**Body**: `DeviceInfo` (`deviceId`, `batteryLevel`, `location`, `imei`, `custom1`–`custom3`, …)

**Response**: `OK` empty data

**Errors**: `DEVICE_NOT_FOUND`; multi-tenant may reject unknown device creation

---

### POST `/applicationSettings/{deviceId}`

Persist per-application settings from agent.

**Body**: `SyncApplicationSetting[]` (`packageId`, `name`, `type`, `value`, `readonly`, `lastUpdate`)

**Response**: `OK`

**Errors**: `DEVICE_NOT_FOUND`
