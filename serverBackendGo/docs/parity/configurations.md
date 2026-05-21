# Parity: Configurations (`/rest/private/configurations`)

**Java:** `com.hmdm.rest.resource.ConfigurationResource`  
**Go:** `internal/modules/configurations/`

| Endpoint | Status | Notes |
|----------|--------|-------|
| `GET /search` | **Done** | Requires `configurations` permission |
| `GET /search/{value}` | **Done** | Name/description filter |
| `GET /list` | **Done** | Phase 4 regression — auth only, id/name |
| `POST /autocomplete` | **Done** | Lookup items |
| `GET /{id}` | **Done** | Editor payload + nested apps/files/settings |
| `PUT /` | **Done** | Create/update; `settingsjson` for extended fields |
| `DELETE /{id}` | **Done** | `error.notempty.configuration` when devices assigned |
| `PUT /copy` | **Done** | Copies row + child tables |
| `GET /applications` | **Done** | Picker catalog |
| `GET /applications/{id}` | **Done** | Linked apps for configuration |
| `PUT /application/upgrade` | **Done** | Bumps junction to latest version |

## Partial

| Area | Note |
|------|------|
| Push notify on save | **Done** — `platform/push` enqueues `configUpdated` (Phase 9) |
| Full Liquibase column parity | Extended UI fields stored in `settingsjson` JSONB |
