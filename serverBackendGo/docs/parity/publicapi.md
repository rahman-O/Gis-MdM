# Parity: Public API (`PublicResource`)

**Go module**: `internal/modules/publicapi`  
**Base path**: `/rest/public`  
**Java**: `com.hmdm.rest.resource.PublicResource`

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/name` | **Done** | Rebranding JSON from env |
| GET | `/logo` | **Done** | File stream or redirect to `/images/logo.png` |
| POST | `/applications/upload` | **Done** | MD5(deviceId+HASH_SECRET); AppList utility |

**Out of scope**: `PublicFilesResource` (`/rest/public/files`) — deprecated in Java.
