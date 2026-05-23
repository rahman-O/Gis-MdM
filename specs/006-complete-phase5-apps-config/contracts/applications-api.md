# API Contract: Applications (`/rest/private/applications`)

**Base path**: `/rest/private/applications`  
**Auth**: Session and/or `Authorization: Bearer <jwt>`  
**Envelope**: Headwind standard

**Java reference**: `com.hmdm.rest.resource.ApplicationResource`  
**React reference**: `frontend/src/features/applications/services/applicationService.ts`

---

### GET `/search`

Tenant application catalog.

**Permissions**: `applications`

**Response data**: `Application[]`

---

### GET `/search/{value}`

Filter by name/pkg.

**Permissions**: `applications`

---

### POST `/autocomplete`

**Body**: JSON string filter

**Response data**: `{ id, name }[]`

---

### GET `/{id}`

Single application.

**Response data**: `Application`

---

### GET `/{id}/versions`

**Response data**: `ApplicationVersion[]`

---

### PUT `/android`

Create/update Android application.

**Body**: `Application`

**Permissions**: `applications`

---

### PUT `/web`

Create/update web application.

**Body**: `Application`

---

### PUT `/versions`

Create/update application version.

**Body**: `ApplicationVersion`

---

### DELETE `/{id}`

Delete application (cascade rules per Java).

---

### DELETE `/versions/{id}`

Delete version.

---

### PUT `/validatePkg`

**Body**: `{ id?, name?, pkg }`

**Response data**: `Application[]` — conflicting apps with same package

---

### GET `/configurations/{id}`

Configurations linked to application `id`.

**Response data**: `ApplicationConfigurationLink[]`

---

### POST `/configurations`

**Body**: `LinkConfigurationsToAppRequest` — `{ applicationId, configurations: [{ configurationId, action, selected, ... }] }`

Update junction rows.

---

### GET `/version/{versionId}/configurations`

Version-level configuration links.

---

### POST `/version/configurations`

**Body**: `LinkConfigurationsToAppVersionRequest`

---

### GET `/admin/search` (super-admin)

Shared/common applications catalog.

**Permissions**: super-admin only

---

### GET `/admin/search/{value}` (super-admin)

---

### GET `/admin/common/{id}` (super-admin)

Turn application into shared/common app (legacy GET mutation).

---

## Out of scope (Phase 5)

| Endpoint | Phase |
|----------|-------|
| `POST /private/web-ui-files` | 6 — APK/binary upload for `url` fields |

Phase 5 accepts `url` / `filePath` on PUT when already known.

## Partial parity (v1)

| Area | Note |
|------|------|
| Icon upload (`iconId`) | May defer until `icons` module |
| Plugin hooks on app save | Out of scope |
