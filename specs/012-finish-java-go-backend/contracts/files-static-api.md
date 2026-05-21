# API Contract: Agent file download (012 P2)

**Path**: `GET /files/{customerFilesDir}/{filePath...}`  
**Java reference**: `FilesResource.downloadFile`, public file servlet  
**Go**: `internal/app` static route + `platform/storage.LocalStore`

## Behavior

- Serve files from `FILES_DIRECTORY` on disk.
- Path must match customer `filesdir` + relative path stored in configuration/file records.
- Content-Type from extension; 404 if missing.
- No Headwind JSON envelope (binary response).

## Security

- Public read (same as Java agent download URLs embedded in sync).
- Optional future: signed URLs — out of scope v1.

## Related private API (regression)

| Method | Path |
|--------|------|
| POST | `/rest/private/config-files` |
| POST | `/rest/private/web-ui-files` |
| GET | `/rest/private/web-ui-files/limit` |

## Quota (012)

Before upload:

- Sum sizes in `uploadedfiles` for customer vs `customers.sizeLimit`.
- Error: `error.size.limit` or Java-equivalent message key.
