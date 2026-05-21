# API Contract: Devicelog plugin — gap endpoints (012 P1)

**Base**: `/rest/plugins/devicelog`  
**Auth**: Bearer JWT (private); agent upload public  
**Java reference**: `DeviceLogResource`, `DeviceLogPluginSettingsResource`  
**Go module**: `internal/modules/plugins/devicelog/`  
**Prior art**: [011 contract](../011-complete-migration-gaps/contracts/plugins-devicelog-gaps-api.md)

## New endpoints (012)

| Method | Path | Status in Go before 012 |
|--------|------|-------------------------|
| POST | `/log/private/search/export` | ❌ Missing |
| GET | `/log/rules/{deviceNumber}` | ❌ Missing |

## POST `/log/private/search/export`

**Body**: Same filter shape as `POST /log/private/search`.  
**Response**: CSV/stream export (not JSON envelope).

## GET `/log/rules/{deviceNumber}`

**Response**: Headwind envelope with array of active rules for device.

## Existing (regression)

- GET/PUT `/devicelog-plugin-settings/private`
- PUT/DELETE `/devicelog-plugin-settings/private/rule`
- POST `/log/private/search`
- POST `/rest/plugins/devicelog/log/list/{deviceNumber}` (agent upload)
