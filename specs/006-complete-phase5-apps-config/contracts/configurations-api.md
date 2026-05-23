# API Contract: Configurations (`/rest/private/configurations`)

**Base path**: `/rest/private/configurations`  
**Auth**: Session and/or `Authorization: Bearer <jwt>`  
**Envelope**: `{ "status": "OK"|"ERROR", "message"?: string, "data"?: T }`

**Java reference**: `com.hmdm.rest.resource.ConfigurationResource`  
**React reference**: `frontend/src/features/configurations/configurationService.ts`

---

### GET `/search`

Full configuration list for Configurations page.

**Permissions**: `configurations`

**Response data**: `Configuration[]` (includes fields for table: name, type, description, deviceCount)

---

### GET `/search/{value}`

Filter configurations by name (parity).

**Permissions**: `configurations`

---

### GET `/list`

Minimal id/name list (implemented Phase 4 — must not regress).

**Permissions**: authenticated (no `configurations` required in Java)

**Response data**: `LookupItem[]` — `{ id, name }`

---

### POST `/autocomplete`

**Body**: JSON string filter value (same pattern as groups/devices)

**Response data**: `LookupItem[]`

---

### GET `/{id}`

Load configuration for editor.

**Response data**: `Configuration` with nested `applications`, `files`, `applicationSettings`, design fields

**Errors**: not found, permission denied

---

### PUT `/`

Create (`id` null/omitted) or update (`id` set).

**Permissions**: `configurations`

**Body**: Full `Configuration` object (React `buildCreateConfigurationBody` / `mergeConfigurationForUpdate`)

**Response data**: Saved `Configuration`

**Errors**: duplicate name, permission denied

---

### DELETE `/{id}`

**Permissions**: `configurations`

**Errors**: not empty (devices still assigned), not found

---

### PUT `/copy`

**Body**: `{ "id": number, "name": string, "description"?: string }`

**Response**: OK

---

### GET `/applications`

All applications available for configuration editor picker (may return application-shaped rows).

---

### GET `/applications/{id}`

Applications linked to configuration `id` for Applications tab.

---

### PUT `/application/upgrade`

**Body**: `{ "configurationId": number, "applicationId": number }` (version resolved server-side to latest)

**Response data**: Updated configuration fragment or OK

**Note**: Push notify stubbed in Phase 5

---

## Type mapping

| UI `ConfigurationKind` | Backend `type` |
|--------------------------|----------------|
| WORK | 0 |
| COMMON | 1 |

## Partial parity (v1)

| Area | Note |
|------|------|
| QR code image generation | May return `qrCodeKey` + `baseUrl` only |
| Push after upgrade | No-op |
| Every Liquibase configuration column | Ship columns used by React editor first |
