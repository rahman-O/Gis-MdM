# API Contract: Devices — search enrichment (012 P1)

**Base**: `/rest/private/devices`  
**Auth**: Bearer JWT + session  
**Java reference**: `com.hmdm.rest.resource.DeviceResource`, `com.hmdm.persistence.domain.DeviceSearchRequest`  
**Go module**: `internal/modules/devices/`  
**Parity**: `serverBackendGo/docs/parity/devices.md`

## POST `/search`

**Request body** (React-aligned; all optional except paging):

| Field | Type | Notes |
|-------|------|-------|
| pageNum | int | ≥1 |
| pageSize | int | 1–500 |
| value | string | ILIKE on number, description, imei, phone |
| groupId | int | Group filter |
| configurationId | int | Config filter |
| status | string | `green` \| `yellow` \| `red` (lastupdate bands) |
| androidVersion | string | `infojson` |
| launcherVersion | string | `infojson` |
| mdmMode | bool | `infojson` |
| kioskMode | bool | `infojson` |
| installationStatus | string | App install state filter |
| sortBy | string | `LAST_UPDATE`, `NUMBER`, … |
| sortDir | string | `asc` \| `desc` |
| dateFrom, dateTo | int64 | ms epoch on lastupdate |
| onlineEarlierMillis, onlineLaterMillis | int64 | Online window |
| enrollmentDateFrom, enrollmentDateTo | int64 | Enrollment window |
| imeiChanged | bool | If supported in Java |
| fastSearch | bool | Exact number match mode |

**Response** (`status: OK`):

```json
{
  "configurations": { "1": { "id": 1, "name": "Default" } },
  "devices": {
    "items": [ { "id": 1, "number": "hmdm-001", "statusCode": "green", ... } ],
    "totalItemsCount": 42
  }
}
```

## GET `/number/{number}`

**Response**: `DeviceView` with optional nested `info`:

```json
{
  "id": 1,
  "number": "hmdm-001",
  "statusCode": "green",
  "lastUpdate": 1710000000000,
  "info": {
    "batteryLevel": 85,
    "model": "Pixel",
    "androidVersion": "14",
    "mdmMode": true,
    "applications": [],
    "files": []
  }
}
```

## Unchanged endpoints

`PUT /`, `DELETE /{id}`, `POST /deleteBulk`, `POST /groupBulk`, `POST /autocomplete`, applicationSettings/*, `POST /{id}/description` — regression only.

## Errors

| Key | When |
|-----|------|
| error.permission.denied | Missing `edit_devices` / scope |
| error.internal.server | DB errors |
