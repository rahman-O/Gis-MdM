# Parity: Push messaging

**Java**: `PushApiResource`, `com.hmdm.plugins.push.rest.PushResource`  
**Status**: Phase 7–8 — **Done** (messages + schedule CRUD); **Partial**: schedule cron runner

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
| POST | `/private/searchTasks` | Done (Phase 8) |
| PUT | `/private/task` | Done — permission `plugin_push_delete` per Java |
| DELETE | `/private/task/{id}` | Done |
