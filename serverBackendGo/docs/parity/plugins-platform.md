# Parity: Plugin platform

**Java**: `com.hmdm.plugin.rest.PluginResource`  
**Go**: `internal/modules/plugins/platform`

| Method | Path | Auth | Permission |
|--------|------|------|------------|
| GET | `/rest/plugin/main/private/available` | JWT | authenticated |
| GET | `/rest/plugin/main/private/active` | JWT | authenticated |
| GET | `/rest/plugin/main/public/registered` | none | — |
| POST | `/rest/plugin/main/private/disabled` | JWT | `plugins_customer_access_management` |

**Notes**: `ENABLED_PLUGINS` env filters build-time catalog. Disabled rows in `pluginsdisabled`.
