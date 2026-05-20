# Parity: Push messaging

**Java**: `PushApiResource`, `com.hmdm.plugins.push.rest.PushResource`  
**Status**: Phase 7 — **Done** (API); **Partial**: plugin schedule cron

## Private API (`/rest/private/push`)

| Method | Path | Permission | Status |
|--------|------|------------|--------|
| POST | `/` | `push_api` | Done |

**React**: `frontend/src/features/push/pushService.ts`

## Plugin (`/rest/plugins/push`)

| Method | Path | Status |
|--------|------|--------|
| POST | `/private/search` | Done |
| POST | `/private/send` | Done |
| DELETE | `/private/{id}` | Done |
| GET | `/private/purge/{days}` | Done |
| POST | `/private/searchTasks` | Partial — not implemented |
| PUT | `/private/task` | Partial — not implemented |
| DELETE | `/private/task/{id}` | Partial — not implemented |
