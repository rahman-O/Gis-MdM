# API Contract: Roles (`/rest/private/roles`)

**Base path**: `/rest/private/roles`  
**Auth**: Session cookie and/or JWT Bearer  
**Envelope**: Headwind standard

## Endpoints

### GET `/permissions`

**Response data**: `Permission[]` — `{ id, name, description?, superadmin? }`

**Auth**: role management access (super admin or `settings` permission).

---

### GET `/all`

**Response data**: `UserRole[]` with embedded or linked permission ids/names per Java shape.

**Auth**: same as above.

---

### PUT `/`

Create or update role.

**Body**:
```json
{
  "id": null,
  "name": "Custom Role",
  "description": "Optional",
  "permissions": [{ "id": 1 }]
}
```

**Success**: `status: OK`

**Errors**: `error.duplicate.role` (DUPLICATE_ENTITY message), `error.permission.denied`

---

### DELETE `/{id}`

Delete role by id.

**Success**: `status: OK`

**Errors**: `error.permission.denied`; referential integrity errors as ERROR message

---

## Related endpoint (users module)

`GET /rest/private/users/roles` — same role list for Users screen dropdown; may
return subset or identical list; MUST stay consistent on `id` and `name`.
