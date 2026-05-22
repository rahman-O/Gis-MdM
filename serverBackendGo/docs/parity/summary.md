# Summary API parity

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| GET | `/rest/private/summary/devices` | `SummaryResource.getDeviceStats` | **Done** |

**Java:** `backend/server/src/main/java/com/hmdm/rest/resource/SummaryResource.java`

**Note:** Status summary, totals, enrollment, `installSummary`, and per-config app status charts use `devicestatuses` (`000011`). Monthly enrollment series still simplified.
