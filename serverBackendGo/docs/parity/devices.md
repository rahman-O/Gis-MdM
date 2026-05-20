# Parity: Devices (`/rest/private/devices`)

**Java:** `com.hmdm.rest.resource.DeviceResource`  
**Go:** `internal/modules/devices/`

| Endpoint | Status | Notes |
|----------|--------|-------|
| `POST /search` | **Done** | `pageNum`, tenant scope, group access; configurations map |
| `GET /number/{number}` | **Done** | Single device view |
| `PUT /` | **Done** | Create/update; bulk config via `ids` + `configurationId` |
| `DELETE /{id}` | **Done** | Requires `edit_devices` |
| `POST /deleteBulk` | **Done** | |
| `POST /groupBulk` | **Done** | `set` / clear groups |
| `POST /autocomplete` | **Done** | Up to 10 matches |
| `GET /{id}/applicationSettings` | **Done** | |
| `POST /{id}/applicationSettings` | **Done** | |
| `POST /{id}/applicationSettings/notify` | **Done** | Push stub (no FCM) |
| `POST /{id}/description` | **Done** | Requires `edit_device_desc` |

## Partial

| Area | Note |
|------|------|
| Search enrichment | No nested apps/files in configuration map |
| Advanced filters | `mdmMode`, `launcherVersion`, `deviceStatuses` deferred |
| `infojson` telemetry | Minimal columns only |
| Fast search exact match | Basic ILIKE text search |
