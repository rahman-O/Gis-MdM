# API Contract: Device log plugin

**Java references**:
- `DeviceLogResource` → `/plugins/devicelog/log`
- `DeviceLogPluginSettingsResource` → `/plugins/devicelog/devicelog-plugin-settings`

## Settings — `/rest/plugins/devicelog/devicelog-plugin-settings`

| Method | Path | Permission |
|--------|------|------------|
| GET | `/private` | `plugin_devicelog_access` |
| PUT | `/private` | `plugin_devicelog_access` |
| PUT | `/private/rule` | `plugin_devicelog_access` |
| DELETE | `/private/rule/{id}` | `plugin_devicelog_access` |

## Logs — `/rest/plugins/devicelog/log`

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| POST | `/private/search` | JWT + permission | Paginated log search |
| POST | `/private/search/export` | JWT + permission | Export search results |
| POST | `/list/{deviceNumber}` | public | Device uploads log batch JSON |
| GET | `/rules/{deviceNumber}` | public | Active rules for device |

**Upload body**: list of log records (see Java `UploadedDeviceLogRecord`).

**Storage**: Postgres tables `plugin_devicelog_*` only in Phase 8.
