# Parity: QR Enrollment (`/rest/public/qr`)

**Java**: `com.hmdm.rest.resource.QRCodeResource`  
**Status**: Phase 7 + **015** — **Done** (full Device Owner provisioning bundle)

| Method | Path | Status |
|--------|------|--------|
| GET | `/{qrCodeKey}` | Done — PNG with checksum, WiFi, admin extras |
| GET | `/json/{qrCodeKey}` | Done — same JSON as PNG |

**015 fixes**: Query params `deviceId`, `create`, `useId`, `group`; SHA-256 APK checksum; `settingsjson` MDM fields; loopback URL rewrite.

**React**: `enrollmentQrQuery.ts`, `EnrollmentQrPage.tsx`
