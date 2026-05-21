# Parity: Devices (`/rest/private/devices`)

**Java:** `com.hmdm.rest.resource.DeviceResource`  
**Go:** `internal/modules/devices/`

| Endpoint | Status | Notes |
|----------|--------|-------|
| `POST /search` | **Done** | `pageNum`, tenant scope, group access; configurations map; React-aligned filters (`status`, `mdmMode`, `androidVersion`, date windows, sort) |
| `GET /number/{number}` | **Done** | Single device view; nested `info` from `devices.info` + `infojson` |
| `PUT /` | **Done** | Create/update; bulk config via `ids` + `configurationId` |
| `DELETE /{id}` | **Done** | Requires `edit_devices` |
| `POST /deleteBulk` | **Done** | |
| `POST /groupBulk` | **Done** | `set` / clear groups |
| `POST /autocomplete` | **Done** | Up to 10 matches |
| `GET /{id}/applicationSettings` | **Done** | |
| `POST /{id}/applicationSettings` | **Done** | |
| `POST /{id}/applicationSettings/notify` | **Done** | Enqueues `appConfigUpdated` via `platform/push` (Phase 9) |
| `POST /{id}/description` | **Done** | Requires `edit_device_desc` |

## Partial

| Area | Note |
|------|------|
| Search enrichment | No nested apps/files in configuration map |
| `installationStatus` | Filter uses `infojson.applications[]`; Java also uses `deviceStatuses` table (not in Go schema yet) |
| `launcherVersion` | Matches `infojson.launcherVersion`; Java `mdm_device_launcher_version()` from `info` blob not replicated |
| Sort `INSTALLATIONS` / `FILES` | Require `deviceStatuses` table |
| Fast search | Exact `number` / `fastsearch` column match when `fastSearch: true` |
