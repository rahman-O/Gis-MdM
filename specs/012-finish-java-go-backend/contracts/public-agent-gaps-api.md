# API Contract: Public & agent gaps (012 P3)

**Auth**: Public routes unauthenticated unless noted  
**Java references**: `StatsResource`, `VideosResource`, `UpdateResource`, `SummaryResource`

## Stats (NEW module)

| Method | Path | Body | Response |
|--------|------|------|----------|
| PUT | `/rest/public/stats` | Usage stats JSON (Java `UsageStats` shape) | `{ status: OK }` |

**Go module**: `internal/modules/stats/`

## Videos (conditional)

| Method | Path | Notes |
|--------|------|-------|
| POST | `/rest/public/videos/{fileName}` | Upload training video |
| GET | `/rest/public/videos/{fileName}` | Download |

**Go module**: `internal/modules/videos/` **or** ⊘ documented in parity if unused.

**Env**: `MODULE_VIDEOS_ENABLED`, `VIDEO_DIRECTORY`

## Updates (extend)

| Method | Path | 012 change |
|--------|------|------------|
| GET | `/rest/private/update/check` | Regression |
| POST | `/rest/private/update` | Download remote APK; apply; optional `sendStats` |

## Summary (extend)

| Method | Path | 012 change |
|--------|------|------------|
| GET | `/rest/private/summary/devices` | Populate chart arrays from DB when data exists |

## Agent paths (regression + files)

| Method | Path | Module |
|--------|------|--------|
| GET/POST | `/rest/public/sync/*` | sync + hooks |
| GET | `/rest/notifications/device/{deviceNumber}` | notifications |
| GET | `/rest/notification/polling/{deviceNumber}` | notifications |

See [files-static-api.md](./files-static-api.md) for `/files/*`.
