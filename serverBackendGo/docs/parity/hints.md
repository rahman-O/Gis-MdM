# Hints API parity (`HintResource`)

**Java**: `backend/server/src/main/java/com/hmdm/rest/resource/HintResource.java`  
**Go**: `internal/modules/hints/adapter/http/handler.go`  
**Base**: `/rest/private/hints`

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/history` | Done | `string[]` hint keys for current user |
| POST | `/history` | Done | Body: JSON string or plain hint key |
| POST | `/enable` | Done | Clears user hint history |
| POST | `/disable` | Done | Marks all `userhinttypes` keys as shown |

**Auth**: Session cookie or JWT Bearer (Swagger Authorize).

**Tables**: `userhints`, `userhinttypes` — migration `000004_hints_tables.up.sql`.
