# Parity: Device Sync (`/rest/public/sync`)

**Java**: `com.hmdm.rest.resource.SyncResource`  
**Status**: Phase 7 — **Done** (core paths); **Partial**: `SyncResponseHook` plugins, full settings merge

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/configuration/{deviceId}` | Done | SyncResponse + signatures |
| POST | `/configuration/{deviceId}` | Done | On-demand enroll |
| POST | `/info` | Done | Telemetry |
| POST | `/applicationSettings/{deviceId}` | Done | Per-app settings |

**Env**: `SECURE_ENROLLMENT`, `PREVENT_DUPLICATE_ENROLLMENT`, `HASH_SECRET`, `REBRANDING_MOBILE_NAME`
