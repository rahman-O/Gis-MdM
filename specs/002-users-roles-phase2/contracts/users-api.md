# API Contract: Users (`/rest/private/users`)

**Base path**: `/rest/private/users`  
**Auth**: Session cookie and/or `Authorization: Bearer <jwt>`  
**Envelope**: `{ "status": "OK"|"ERROR", "message"?: string, "data"?: T }`

## Endpoints

### GET `/current`

Returns full user for authenticated principal.

**Response data**: User object (no `password`, no `authToken`).

**Errors**: `OK` with `data: null` if unauthenticated (match Java).

---

### PUT `/details`

Update own `name` and `email`.

**Body**:
```json
{ "id": 1, "name": "Admin", "email": "admin@localhost" }
```

**Success**: `status: OK`, `data` updated user (optional).

**Errors**: `error.duplicate.email`, `error.user.not.found`

---

### PUT `/current`

Change own password.

**Body**:
```json
{
  "id": 1,
  "login": "admin",
  "oldPassword": "<MD5_HEX>",
  "newPassword": "<MD5_HEX>"
}
```

**Success**: `status: OK`, message `success.operation.completed`

**Errors**: `error.password.wrong`, `error.password.empty`, `error.permission.denied` (wrong id)

---

### GET `/all?filter={optional}`

List tenant users (admin).

**Auth**: requires `settings` permission (or super admin / org admin).

**Response data**: `User[]` with `editable` flag; no secrets.

**Errors**: `error.permission.denied`

---

### PUT `/`

Create or update user.

**Body** (create — no `id`):
```json
{
  "login": "user1",
  "name": "User One",
  "email": "u1@test.local",
  "newPassword": "<MD5_HEX>",
  "userRole": { "id": 2 },
  "allDevicesAvailable": true,
  "allConfigAvailable": true
}
```

**Body** (update — with `id`): same fields; `newPassword` optional.

**Success**: `status: OK`

**Errors**: `error.duplicate.login`, `error.duplicate.email`, `error.password.empty`, `error.permission.denied`

---

### DELETE `/other/{id}`

Delete user by id (not self).

**Success**: `status: OK`

**Errors**: `error.permission.denied`, delete constraint errors as message string

---

### GET `/roles`

Assignable roles for dropdowns.

**Response data**: `{ "id": number, "name": string }[]`

---

## Out of scope (this contract)

- `GET /{id}` — unused by React
- `GET /impersonate/{id}`
- `/superadmin/*`
