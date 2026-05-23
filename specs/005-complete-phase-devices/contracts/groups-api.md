# API Contract: Groups (`/rest/private/groups`)

**Base path**: `/rest/private/groups`  
**Auth**: Bearer / session  
**Envelope**: Headwind standard

**Java reference**: `com.hmdm.rest.resource.GroupResource`

---

### GET `/search`

All groups for current customer.

**Response data**: `Group[]` or `{ id, name }[]`

---

### GET `/search/{value}`

Filter groups by name (legacy parity).

---

### POST `/autocomplete`

**Body**: filter string

**Response data**: `LookupItem[]`

---

### PUT `/`

Create (`id` null) or update group.

**Permissions**: `settings`

**Errors**: `error.duplicate.group`, `error.permission.denied`

---

### DELETE `/{id}`

**Permissions**: `settings`

**Errors**: `error.notempty.group` when devices still assigned

---

## React consumers

- `frontend/src/features/devices/deviceService.ts` — `GET /search`
- `frontend/src/features/groups/groupService.ts` — CRUD
