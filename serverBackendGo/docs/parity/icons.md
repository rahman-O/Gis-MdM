# Parity: Icons (`IconResource`)

**Go module**: `internal/modules/icons`  
**Base path**: `/rest/private/icons`  
**Java**: `com.hmdm.rest.resource.IconResource`

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/search` | **Done** | Tenant-scoped list |
| GET | `/search/{value}` | **Done** | Name filter |
| PUT | `/` | **Done** | Create/update |
| DELETE | `/{id}` | **Done** | Requires `settings` permission |
