# Parity: Groups (`/rest/private/groups`)

**Java:** `com.hmdm.rest.resource.GroupResource`  
**Go:** `internal/modules/groups/`

| Endpoint | Status | Notes |
|----------|--------|-------|
| `GET /search` | **Done** | All groups for tenant |
| `GET /search/{value}` | **Done** | Name filter |
| `POST /autocomplete` | **Done** | Lookup items |
| `PUT /` | **Done** | Requires `settings`; duplicate name check |
| `DELETE /{id}` | **Done** | `error.notempty.group` when devices assigned |
