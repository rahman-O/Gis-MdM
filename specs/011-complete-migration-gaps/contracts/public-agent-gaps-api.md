# API Contract: Public & agent gaps (Phase 9 P3)

## Stats

**Java**: `com.hmdm.rest.resource.StatsResource`  
**Base**: `/rest/public/stats`

| Method | Path | Body | Response |
|--------|------|------|----------|
| PUT | `/` | `UsageStats` JSON | `OK` |

**Auth**: Public (match Java).  
**Persistence**: `usagestats` table.

---

## Videos

**Java**: `com.hmdm.rest.resource.VideosResource`  
**Base**: `/rest/public/videos`

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/` | Multipart upload training video |
| GET | `/{fileName}` | Download stream `application/octet-stream` |

**Env**: `VIDEO_DIRECTORY`, `BASE_URL` for path in upload response.

---

## Updates (extend existing module)

**Java**: `com.hmdm.rest.resource.UpdateResource`

| Method | Path | Phase 9 completion |
|--------|------|-------------------|
| GET | `/private/update/check` | Done |
| POST | `/private/update` | Remote APK download + apply flags |
| — | `sendStats` in body | Persist usage stats when true |

---

## Agent file download

**Java**: Files servlet / static under customer files dir

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/files/{customerFilesDir}/{filePath}` or documented legacy path | Agent downloads configuration files |

**Implementation**: Gin static or dedicated handler; must match paths agents use (verify Java `web.xml` / `FilesResource`).

---

## Summary (extend)

| Method | Path | Enhancement |
|--------|------|-------------|
| GET | `/rest/private/summary/devices` | Charts from `devicestatuses` when data exists |
