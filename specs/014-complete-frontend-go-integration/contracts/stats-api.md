# API Contract: Stats (014 — new module)

**Base**: `/rest/public/stats`  
**Method**: `PUT`  
**Auth**: none (public agent/heartbeat — same as Java `StatsResource`)

**Java**: `com.hmdm.rest.resource.StatsResource`  
**Table**: `usagestats` (migration `000014`, prerequisite 013)

**Module**: `internal/modules/stats`  
**Config**: `MODULE_STATS_ENABLED=true` (default on in dev)

---

## PUT `/rest/public/stats`

**Content-Type**: `application/json`

**Body** (`UsageStats`):

| Field | Type | Required |
|-------|------|----------|
| `instanceId` | string | recommended (upsert key) |
| `webVersion` | string | no |
| `community` | boolean | no |
| `devicesTotal` | number | no |
| `devicesOnline` | number | no |
| `cpuTotal` | number | no |
| `cpuUsed` | number | no |
| `ramTotal` | number | no |
| `ramUsed` | number | no |
| `scheme` | string | no |
| `arch` | string | no |
| `os` | string | no |
| `ts` | string (date) | no — default server date (today) |

**Behavior**:

- Upsert on `(ts, instanceid)` — update counters on conflict.
- `200` + `{ "status": "OK" }` on success.

**Errors**:

- Malformed JSON → ERROR
- Module disabled → 404 or ERROR per platform convention

---

## Example

```bash
curl -X PUT http://localhost:8080/rest/public/stats \
  -H "Content-Type: application/json" \
  -d '{
    "instanceId": "mdm-go-dev",
    "webVersion": "1.0.0",
    "community": true,
    "devicesTotal": 100,
    "devicesOnline": 12,
    "cpuTotal": 4,
    "cpuUsed": 1,
    "ramTotal": 8192,
    "ramUsed": 2048,
    "scheme": "https",
    "arch": "amd64",
    "os": "linux"
  }'
```

---

## Out of scope

- Admin UI for stats (no React page in 014)
- Aggregation dashboards
