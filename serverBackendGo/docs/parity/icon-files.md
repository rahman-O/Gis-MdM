# Parity: Icon file upload (`IconFileResource`)

**Go module**: `internal/modules/icons` (handler on `/rest/private/icon-files`)  
**Java**: `com.hmdm.rest.resource.IconFileResource`

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| POST | `/` | **Done** | Multipart `file`; square image; 144×144 PNG; `uploadedfiles` row |

**Errors**: `error.icon.dimension.invalid` when width ≠ height.
