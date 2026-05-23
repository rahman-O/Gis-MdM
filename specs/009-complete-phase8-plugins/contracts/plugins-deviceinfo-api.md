# API Contract: Device info plugin

**Java references**:
- `DeviceInfoResource` → `/plugins/deviceinfo/deviceinfo`
- `DeviceInfoPluginSettingsResource` → `/plugins/deviceinfo/deviceinfo-plugin-settings`

**Auth**: JWT on private routes; public upload unauthenticated (device number in path).

## Settings — `/rest/plugins/deviceinfo/deviceinfo-plugin-settings`

| Method | Path | Permission | Purpose |
|--------|------|------------|---------|
| GET | `/private` | `plugin_deviceinfo_access` | Load customer settings |
| PUT | `/private` | `plugin_deviceinfo_access` | Save customer settings |

## Device info — `/rest/plugins/deviceinfo/deviceinfo`

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| PUT | `/public/{deviceNumber}` | public | Upload dynamic info list |
| GET | `/private/{deviceNumber}` | JWT + permission | Device detail aggregate |
| GET | `/private/search/device` | JWT + permission | Search devices (query params) |
| POST | `/private/search/dynamic` | JWT + permission | Paginated dynamic info search |
| POST | `/private/export` | JWT + permission | Export dynamic info (file/stream) |

**Public PUT body**: `List<DeviceDynamicInfo>` JSON array.

**Plugin disabled guard**: reject when `deviceinfo` disabled for customer (cache/DB).
