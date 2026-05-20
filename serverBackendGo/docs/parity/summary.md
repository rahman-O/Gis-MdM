# Summary API parity

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| GET | `/rest/private/summary/devices` | `SummaryResource.getDeviceStats` | **Done** |

**Java:** `backend/server/src/main/java/com/hmdm/rest/resource/SummaryResource.java`

**Note:** Status summary (green/yellow/red), totals, and enrollment counts use real SQL when `devices` exist. Install-by-config charts remain simplified (no `devicestatuses` table in Phase 4).
