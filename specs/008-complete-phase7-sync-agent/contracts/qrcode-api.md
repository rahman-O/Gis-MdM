# API Contract: QR Enrollment (`/rest/public/qr`)

**Base path**: `/rest/public/qr`  
**Auth**: None  
**Envelope**: N/A for PNG; JSON endpoint returns raw string or JSON body per Java

**Java reference**: `com.hmdm.rest.resource.QRCodeResource`  
**React reference**: `frontend/src/features/devices/enrollmentQrQuery.ts`, `EnrollmentQrPage.tsx`

---

### GET `/{qrCodeKey}`

Generate QR code PNG for configuration enrollment.

**Query params**:
- `size` (int, optional) — image dimensions
- `deviceId` (string, optional)
- `create` (string, optional) — create-on-demand flag
- `useId` (string, optional) — which identifier to embed
- `group` (repeatable) — group names for new device

**Response**: `image/png` stream (`Content-Type: application/octet-stream` in Java)

**Errors**: 500 when configuration key not found

---

### GET `/json/{qrCodeKey}`

Return provisioning extras JSON string for same parameters (except `size` ignored).

**Response**: JSON body (Android Device Admin extras bundle text)

**Errors**: 500 when key not found

**Notes**:
- Embeds launcher APK URL and SHA-256 when `configurations.mainAppId` resolves to version with URL
- `launcherUrl` on configuration overrides APK URL for closed networks
- Localhost URL rewrite for dev agents (match Java `rewriteLoopbackDownloadUrlForQr`)
