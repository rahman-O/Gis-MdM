# API Contract: Public static files — `/files/*`

**Purpose**: Serve tenant files and launcher APKs referenced by QR provisioning and `SyncResponse.files[].url`.

**Java reference**: Local file access in `QRCodeResource.calculateApkHash` and servlet static `/files/` mapping (Tomcat); `PublicFilesResource` deprecated.

**Go placement**: `internal/app` router or `internal/platform/storage/http` — not a new business module.

---

## GET `/files/{customerFilesDir}/{relativePath...}`

**Auth**: None (public read, same as Java agent download)

**Behavior**:

- Map path under `FILES_DIRECTORY` env (default `/var/lib/hmdm/files` or `./data/files` in dev).
- `customerFilesDir` matches `customers.filesdir` (e.g. `customer-1`).
- Reject path traversal (`..`).
- `Content-Type` by extension; `Cache-Control` reasonable for APK (Java used no-cache on QR only).

**Success**: `200` file stream

**Errors**:

- `404` missing file — sync/QR should log; admin sees broken URL in quickstart diagnostics

---

## URL builder alignment

`storage.BuildPublicURL(baseURL, filesDir, relPath)` MUST produce paths this handler serves.

Example: `http://192.168.1.10:8080/files/customer-1/apk/launcher.apk`

---

## Out of scope

- Upload via this path (uploads remain `/rest/private/web-ui-files` and related).
- Authentication on static files (match Java public read).
