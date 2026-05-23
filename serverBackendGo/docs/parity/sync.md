# Parity: Device Sync (`/rest/public/sync`)

**Java**: `com.hmdm.rest.resource.SyncResource`  
**Status**: Phase 7 + **015** — **Done** (core paths); **Partial**: `SyncResponseHook` plugins

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/configuration/{deviceId}` | Done | SyncResponse + signatures; loads device applicationSettings |
| POST | `/configuration/{deviceId}` | Done | On-demand enroll; `configuration` resolves by **qrcodekey** first |
| POST | `/info` | Done | Telemetry; **014** upserts `devicestatuses` from payload (`applications`/`files` presence) |
| POST | `/applicationSettings/{deviceId}` | Done | Per-app settings |

**Env**: `SECURE_ENROLLMENT`, `PREVENT_DUPLICATE_ENROLLMENT`, `HASH_SECRET`, `REBRANDING_MOBILE_NAME`

## SyncResponse policy fields (016)

**Status**: **Done** — `ApplyConfigurationPolicy` in `internal/modules/sync/application/sync_configuration_mapper.go` maps `settingsjson` + `backgroundimageurl` onto `SyncResponse` (Java `SyncResource` parity subset).

| Area | Fields (representative) |
|------|-------------------------|
| Kiosk | `kioskMode`, `kioskHome`, `kioskRecents`, `kioskNotifications`, `kioskSystemInfo`, `kioskKeyguard`, `kioskLockButtons`, `kioskScreenOn`, `kioskExit` |
| Restrictions | `restrictions`, `allowedClasses`, `lockSafeSettings` |
| Connectivity | `gps`, `bluetooth`, `wifi`, `mobileData`, `usbStorage`, `showWifi` |
| Display / device | `orientation`, `brightness`, `timeout`, `volume`, `autoBrightness`, `disableScreenshots`, … |
| Updates | `downloadUpdates`, `systemUpdateType`, `systemUpdateFrom`/`To`, `appUpdateFrom`/`To`, `scheduleAppUpdate` |
| Apps / files | `device_sync_repo` — includes `showIcon`, `screenOrder`, `code`, launcher icon URL, and related flags (016) |

**Application settings merge**: configuration defaults + device overrides; `policyLocks` and readonly configuration settings skip device POST updates.
