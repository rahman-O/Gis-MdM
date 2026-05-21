# API Contract: Icon file upload (Phase 9 P1)

**Base path**: `/rest/private/icon-files`  
**Auth**: Bearer JWT + tenant scope  
**Java reference**: `com.hmdm.rest.resource.IconFileResource`

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/` | Multipart upload icon image |

## POST `/`

**Consumes**: `multipart/form-data`  
**Field**: `file` (image)

**Validation**:

- Image must be readable; width == height or `error.icon.dimension.invalid`
- Resize to 144×144 PNG
- Store under `{FILES_DIRECTORY}/{customer.filesDir}/{uuid}.png`

**Persists**: `uploadedfiles` row (`customerid`, `filepath`, `uploadtime`, …)

**Response** (`OK`):

```json
{
  "status": "OK",
  "data": {
    "id": 1,
    "customerId": 1,
    "filePath": "uuid.png",
    "uploadTime": 1710000000000
  }
}
```

**Errors**: `error.icon.dimension.invalid`, `error.permission.denied`, storage quota errors (Phase E).

## Relationship to icons module

`PUT /rest/private/icons` continues metadata CRUD; clients may reference `uploadedfiles.id` or path from icon-files response per Java UI flow.
