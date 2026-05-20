# Parity: Device info plugin

**Java**: `DeviceInfoResource`, `DeviceInfoPluginSettingsResource`  
**Go**: `internal/modules/plugins/deviceinfo`

| Method | Path |
|--------|------|
| GET/PUT | `/rest/plugins/deviceinfo/deviceinfo-plugin-settings/private` |
| PUT | `/rest/public/plugins/deviceinfo/deviceinfo/public/{deviceNumber}` |
| GET | `/rest/plugins/deviceinfo/deviceinfo/private/{deviceNumber}` |
| POST | `/rest/plugins/deviceinfo/deviceinfo/private/search/dynamic` |

**Partial**: Full multi-table GPS/WiFi export simplified to core `deviceparams` + battery telemetry.
