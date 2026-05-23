# Parity: Configuration Files (`/rest/private/config-files`)

**Java:** `com.hmdm.rest.resource.ConfigurationFileResource`  
**Go:** `internal/modules/configfiles/`

| Endpoint | Status | Notes |
|----------|--------|-------|
| `POST /` | **Done** | Multipart field `file`; writes under `FILES_DIRECTORY/{filesdir}/` |

**Response data:** `{ customerId, filePath, url, name }` (Headwind envelope).

## Partial

| Area | Note |
|------|------|
| `uploadedfiles` DB row + checksum | v1 writes disk only; paths referenced on `PUT /configurations` |
| Storage quota (`sizeLimit`) | Not enforced until customer limit columns used |
