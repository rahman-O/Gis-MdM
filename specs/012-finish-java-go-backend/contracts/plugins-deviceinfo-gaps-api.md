# API Contract: Deviceinfo plugin — gap endpoints (012 P1)

**Base**: `/rest/plugins/deviceinfo`  
**Auth**: Bearer JWT (private)  
**Java reference**: `DeviceInfoResource`, `DeviceInfoPluginSettingsResource`  
**Go module**: `internal/modules/plugins/deviceinfo/`  
**Prior art**: [011 contract](../011-complete-migration-gaps/contracts/plugins-deviceinfo-gaps-api.md)

## New endpoints (012)

| Method | Path | Status in Go before 012 |
|--------|------|-------------------------|
| POST | `/deviceinfo/private/search/device` | ❌ Missing |
| POST | `/deviceinfo/private/export` | ❌ Missing |
| GET | `/deviceinfo-plugin-settings/device/{deviceNumber}` | ❌ Missing |

## POST `/deviceinfo/private/search/device`

**Body**: Device-scoped filter (deviceNumber, date range, dynamic fields).  
**Response**: Headwind envelope with device telemetry rows.

## POST `/deviceinfo/private/export`

**Body**: Export filter (devices, period).  
**Response**: `application/octet-stream` attachment (CSV) matching Java column order.  
**Limits**: Stream; max rows configurable in service.

## GET `/deviceinfo-plugin-settings/device/{deviceNumber}`

**Response**: Merged global + per-device plugin settings JSON.

## Existing (regression)

- GET/PUT `/deviceinfo-plugin-settings/private`
- GET `/deviceinfo/private/{deviceNumber}`
- POST `/deviceinfo/private/search/dynamic`
- PUT `/deviceinfo/public/{deviceNumber}` (agent)
