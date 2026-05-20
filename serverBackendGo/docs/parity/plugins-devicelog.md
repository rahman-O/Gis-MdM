# Parity: Device log plugin

**Java**: `DeviceLogResource`, `DeviceLogPluginSettingsResource`  
**Go**: `internal/modules/plugins/devicelog`

| Method | Path |
|--------|------|
| GET/PUT | `/rest/plugins/devicelog/devicelog-plugin-settings/private` |
| PUT/DELETE | `/rest/plugins/devicelog/devicelog-plugin-settings/private/rule` |
| POST | `/rest/plugins/devicelog/log/private/search` |
| POST | `/rest/public/plugins/devicelog/log/list/{deviceNumber}` |

**Storage**: Postgres tables `plugin_devicelog_*`.
