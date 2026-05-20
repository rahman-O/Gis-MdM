# Summary API parity

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| GET | `/rest/private/summary/devices` | `SummaryResource.getDeviceStats` | Done (empty stats until `devices` migration) |

**Java:** `backend/server/src/main/java/com/hmdm/rest/resource/SummaryResource.java`

**Note:** Returns valid `SummaryResponse` shape with zeros when `devices` table is absent. Full counts will be implemented with the `devices` module (Phase 4).
