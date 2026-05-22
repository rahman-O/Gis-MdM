# Parity: Stats (`/rest/public/stats`)

**Status**: Done (014)

| Method | Path | Java | Go |
|--------|------|------|-----|
| PUT | `/rest/public/stats` | `StatsResource.saveStats` | `internal/modules/stats` |

**Module flag**: `MODULE_STATS_ENABLED=true`

**Table**: `usagestats` (migration `000014`)

**Body**: `UsageStats` JSON (`instanceId`, `devicesTotal`, `devicesOnline`, …) — upsert on `(ts, instanceid)`.
