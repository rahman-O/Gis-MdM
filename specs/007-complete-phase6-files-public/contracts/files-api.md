# API Contract: Files (`/rest/private/web-ui-files`)

**Base path**: `/rest/private/web-ui-files`  
**Auth**: Session and/or `Authorization: Bearer <jwt>`  
**Envelope**: Headwind standard

**Java reference**: `com.hmdm.rest.resource.FilesResource`  
**React reference**: `frontend/src/features/files/filesService.ts`,  
`frontend/src/features/applications/services/webUiFilesService.ts`

---

### GET `/search`

List all uploaded files for current tenant.

**Permissions**: `files`

**Response data**: `FileView[]` (includes `usedByConfigurations`, `usedByIcons` when available)

---

### GET `/search/{value}`

Filter by filepath, description, or external URL (ILIKE).

**Permissions**: `files`

**Response data**: `FileView[]`

---

### POST `/remove`

Delete file metadata and on-disk file when not external.

**Permissions**: `edit_files`

**Body**: `FileView` or `UploadedFile` subset (`id`, `filePath`, `external`, …)

**Errors**: `FILE_USED` when referenced; permission denied on unsafe path

---

### POST `/update`

Create or update uploaded file record (commit after multipart upload).

**Permissions**: `edit_files`

**Body**: `UploadedFile` — `id` null = create; `external` = external URL create;
`tmpPath` set = replace file content

**Response data**: `UploadedFile` on create; empty OK on update

**Errors**: `error.file.save`, `FILE_EXISTS`, permission denied

---

### POST `/`

Multipart upload with optional APK parsing.

**Permissions**: `edit_files`

**Form**: `file` (binary)

**Response data**: `FileUploadResult` (`name`, `serverPath`, `fileDetails?`, `exists?`, `complete?`, `application?`)

**Errors**: `error.size.limit.exceeded`, duplicate version code message for APK conflicts

---

### POST `/raw`

Multipart upload without APK parsing.

**Permissions**: `edit_files`

**Form**: `file`

**Response data**: `FileUploadResult`

---

### GET `/limit`

Tenant storage usage (multi-tenant non-master with `sizeLimit` > 0).

**Permissions**: authenticated (Java does not check `files` on this route)

**Response data**: `{ sizeUsed, sizeLimit }` (MB)

---

### GET `/apps/{url}`

Applications referencing file URL (URL-encoded path parameter).

**Permissions**: `files`

**Response data**: `Application[]`

---

### GET `/configurations/{id}`

Configuration links for uploaded file id.

**Permissions**: `files`

**Response data**: `FileConfigurationLink[]` or `ApplicationConfigurationLink[]` per Java shape

---

### POST `/configurations`

Bulk update configuration ↔ file links.

**Permissions**: `edit_files`

**Body**: `{ configurations: FileConfigurationLink[] }`

**Response**: OK; push notify **stubbed**

---

### Out of scope / partial

| Route | Status |
|-------|--------|
| `GET /{filePath}` (octet download) | **Out** — broken/unused in Java |
| `GET /files/*` (agent download) | **Partial** — see `DownloadFilesServlet`; platform route optional |
