# API Contract: Deviceinfo plugin — gap endpoints (Phase 9 P1)

**Base**: `/rest/plugins/deviceinfo`  
**Auth**: Bearer JWT (private); public paths unchanged  
**Java reference**: `com.hmdm.plugins.deviceinfo.rest.DeviceInfoResource`, `DeviceInfoPluginSettingsResource`

## New / completed endpoints

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/deviceinfo/private/search/device` | Device-scoped search (filter by deviceNumber, date range) |
| POST | `/deviceinfo/private/export` | Export telemetry (CSV or Java-equivalent stream) |
| GET | `/deviceinfo-plugin-settings/device/{deviceNumber}` | Per-device plugin settings overlay |

## POST `/deviceinfo/private/search/device`

**Body**: Filter DTO matching Java (`DeviceInfoFilter` / device search request).  
**Response**: Paginated or list payload in Headwind envelope.

## POST `/deviceinfo/private/export`

**Body**: Export filter (devices, period).  
**Response**: `application/octet-stream` or attachment JSON per Java.  
**Performance**: Stream rows; cap default page size in service layer.

## GET `.../deviceinfo-plugin-settings/device/{deviceNumber}`

**Response**: Settings merged with device-specific overrides (Java parity).

## Existing (unchanged)

- `PUT /deviceinfo/public/{deviceNumber}` — agent telemetry
- `POST /deviceinfo/private/search/dynamic`
- `GET /deviceinfo/private/{deviceNumber}`
