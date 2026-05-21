# API Contract: Devicelog plugin — gap endpoints (Phase 9 P1)

**Base**: `/rest/plugins/devicelog`  
**Java reference**: `com.hmdm.plugins.devicelog.rest.resource.DeviceLogResource`

## New endpoints

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/log/private/search/export` | Export log search results (file download) |
| GET | `/log/rules/{deviceNumber}` | Applied rules for device (agent/UI) |

## POST `/log/private/search/export`

**Body**: Same filter as `POST /log/private/search`.  
**Response**: Export file stream (CSV), matching Java columns.  
**Auth**: `plugin_devicelog` permission (align with search).

## GET `/log/rules/{deviceNumber}`

**Response**: List of `AppliedDeviceLogRule` JSON objects per Java.  
**Scope**: Tenant — device must belong to principal's customer.

## Existing (unchanged)

- Settings CRUD under `/devicelog-plugin-settings/private`
- `POST /log/private/search`
- `POST /rest/public/plugins/devicelog/log/list/{deviceNumber}` — agent upload
