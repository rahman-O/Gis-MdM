# API Contract: Icons (`/rest/private/icons`)

**Base path**: `/rest/private/icons`  
**Auth**: Session and/or `Authorization: Bearer <jwt>`  
**Envelope**: Headwind standard

**Java reference**: `com.hmdm.rest.resource.IconResource`  
**React reference**: `frontend/src/features/icons/iconsService.ts`

---

### GET `/search`

List all icons for current tenant.

**Permissions**: none in Java (authenticated private route)

**Response data**: `Icon[]` with `fileName` joined from uploaded file

---

### GET `/search/{value}`

Filter icons by name (ILIKE).

**Response data**: `Icon[]`

---

### PUT `/`

Create or update icon.

**Body**: `Icon` — `id` null/absent → insert; present → update

**Permissions**: none in Java for PUT

**Response data**: `Icon` (saved row)

**Errors**: `INTERNAL_ERROR` on failure

---

### DELETE `/{id}`

Remove icon by id.

**Permissions**: `settings`

**Response**: OK

**Errors**: permission denied without `settings`
