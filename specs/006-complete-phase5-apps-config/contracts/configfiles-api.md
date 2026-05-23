# API Contract: Configuration Files (`/rest/private/config-files`)

**Base path**: `/rest/private/config-files`  
**Auth**: Session and/or `Authorization: Bearer <jwt>`  
**Envelope**: Headwind standard

**Java reference**: `com.hmdm.rest.resource.ConfigurationFileResource`

---

### POST `/`

Upload a file for use in configuration editor (certificates, payloads, etc.).

**Content-Type**: `multipart/form-data`

**Form field**: `file` — binary upload

**Permissions**: authenticated user with valid `customerId` (same as other private routes)

**Response data** (`FileUploadResult`):

```json
{
  "fileName": "cert.pem",
  "path": "relative/path/under/customer/filesdir",
  "url": "https://server/base/customer/files/cert.pem"
}
```

(Exact field names must match Java `com.hmdm.rest.json.FileUploadResult` and React expectations.)

**Behavior**:

- Write under `{FILES_DIRECTORY}/{customer.filesdir}/{fileName}`
- Create customer subdirectory if missing
- Overwrite existing file (Java logs warn, returns success)

**Errors**: permission denied, disk failure, optional storage limit exceeded

---

## Relationship to configuration save

Uploaded paths are referenced in `Configuration.files[]` on subsequent
`PUT /private/configurations` — not stored only by this endpoint.

## Partial parity (v1)

| Area | Note |
|------|------|
| Storage quota enforcement | Implement if `customers` limit columns exist; else document Partial |
| Virus scanning | Out of scope |
