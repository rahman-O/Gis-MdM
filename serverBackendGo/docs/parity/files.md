# Parity: Files (`FilesResource`)

**Go module**: `internal/modules/files`  
**Base path**: `/rest/private/web-ui-files`  
**Java**: `com.hmdm.rest.resource.FilesResource`

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/search` | **Done** | `files` permission |
| GET | `/search/{value}` | **Done** | Filtered list |
| POST | `/remove` | **Done** | `edit_files`; `error.used.file` |
| POST | `/update` | **Done** | Create / external / update |
| POST | `/` | **Done** | Multipart; APK metadata best-effort |
| POST | `/raw` | **Done** | Multipart without APK parse |
| GET | `/limit` | **Done** | Multi-tenant quota |
| GET | `/apps/{url}` | **Done** | Applications by file URL |
| GET | `/configurations/{id}` | **Done** | Configuration link rows |
| POST | `/configurations` | **Done** | Push notify stubbed |
| GET | `/files/*` (agent download) | **Done** (015) | `internal/app/app.go` + `platform/storage/static_files.go`; maps `FILES_DIRECTORY` |

**Permissions**: `files`, `edit_files` (see `platform/auth/permissions.go`).
