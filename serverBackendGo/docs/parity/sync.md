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
