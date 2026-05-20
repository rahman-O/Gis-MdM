# API Contract: Public API (`/rest/public`)

**Base path**: `/rest/public`  
**Auth**: None (public routes)  
**Envelope**: JSON routes use Headwind standard; `/logo` returns raw image or redirect

**Java reference**: `com.hmdm.rest.resource.PublicResource`  
**Legacy UI**: Angular `rebranding.service.js` (`GET rest/public/name`, `GET rest/public/logo`)

**Out of scope**: `PublicFilesResource` at `/rest/public/files` (deprecated)

---

### GET `/name`

Rebranding metadata for login/signup screens.

**Response data**:

```json
{
  "appName": "string",
  "vendorName": "string",
  "vendorLink": "string",
  "signupLink": "string",
  "termsLink": "string"
}
```

(Java `NameResponse`; year may be client-side)

---

### GET `/logo`

Rebranded logo image.

**Response**: `image/png` stream when `REBRANDING_LOGO` file exists; otherwise HTTP 302/redirect
to default `../images/logo.png` or 404 per Java

**Headers**: `Cache-Control: no-cache`

---

### POST `/applications/upload`

AppList utility upload (multipart).

**Form fields**:
- `file` — optional APK/binary
- `app` — JSON string (`UploadAppRequest`)

**Required JSON fields**: `deviceId`, `hash`, `name`, `pkg`, `version`; when file present also
`localPath`, `fileName`

**Hash**: `MD5(deviceId + HASH_SECRET)` case-insensitive match

**Success**: OK envelope

**Errors**: `error.params.missing`, `Invalid hash`, device not found, duplicate application,
permission denied on unsafe paths

**Side effects**: Inserts `applications` row for device's customer; optional file under tenant files dir
