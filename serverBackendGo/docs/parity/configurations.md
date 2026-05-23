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
| Full Liquibase column parity | Extended UI fields in `settingsjson`; legacy SQL columns → JSON via `000017` on Java dumps |
| `configurationapplicationparameters` | **Done** — `000013`; upsert on save when `skipVersionCheck` set |
| `configurationapplications.remove` / `longtap` | **Done** — `000016`; returned on GET apps list (014) |
| MDM policy flatten on GET/PUT | **Done** (014) — `settingsjson` merged via `Configuration.Policy`; `ParseConfigurationBody` / `ConfigurationResponseMap` |
| `policyLocks` in `settingsjson` | **Done** (016) — map of field key → `true`; allowlist in `domain/policy_locks.go`; round-trip on GET/PUT |
| Readonly app settings | **Done** (016) — configuration-level readonly settings merged into sync; device cannot override locked keys via sync POST |
