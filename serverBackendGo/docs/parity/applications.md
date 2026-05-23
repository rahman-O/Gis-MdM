# Parity: Applications (`/rest/private/applications`)

**Java:** `com.hmdm.rest.resource.ApplicationResource`  
**Go:** `internal/modules/applications/`

| Endpoint | Status | Notes |
|----------|--------|-------|
| `GET /search` | **Done** | Tenant + common apps |
| `GET /search/{value}` | **Done** | Name/pkg filter |
| `POST /autocomplete` | **Done** | Lookup items |
| `GET /{id}` | **Done** | Single application |
| `GET /{id}/versions` | **Done** | Version list |
| `PUT /android` | **Done** | Create/update Android app |
| `PUT /web` | **Done** | Create/update web app |
| `PUT /versions` | **Done** | Create/update version |
| `DELETE /{id}` | **Done** | Tenant-scoped delete |
| `DELETE /versions/{id}` | **Done** | Version delete |
| `PUT /validatePkg` | **Done** | Returns conflicting apps |
| `GET /configurations/{id}` | **Done** | App ↔ configuration links |
| `POST /configurations` | **Done** | Bulk link update |
| `GET /version/{versionId}/configurations` | **Done** | Version-level links |
| `POST /version/configurations` | **Done** | Version link update |
| `GET /admin/search` | **Done** | Super-admin only |
| `GET /admin/search/{value}` | **Done** | Super-admin filter |
| `GET /admin/common/{id}` | **Done** | Mark app as shared/common |

## Partial / out of scope (Phase 5)

| Area | Note |
|------|------|
| `POST /private/web-ui-files` | **Done** (Phase 6 `files` module) |
| Icon upload (`iconId`) | Phase 6 `icons` module |
| Plugin hooks on app save | Not migrated |
