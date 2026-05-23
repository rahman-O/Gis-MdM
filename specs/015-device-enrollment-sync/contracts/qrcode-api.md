# API Contract: QR Enrollment (fixes) — `/rest/public/qr`

**Base path**: `/rest/public/qr`  
**Auth**: None  
**Java reference**: `com.hmdm.rest.resource.QRCodeResource`  
**React reference**: `frontend/src/features/devices/enrollmentQrQuery.ts`

**Delta from Phase 7 (`008`)**: Go MUST return **byte-identical structure** to Java provisioning JSON (not a 3-field stub).

---

## GET `/{qrCodeKey}`

**Query**: `size`, `deviceId`, `create`, `useId`, `group` (repeatable) — all MUST affect generated content.

**Response**: `image/png` — QR encoding UTF-8 JSON object:

Required top-level keys (when configuration valid):

- `android.app.extra.PROVISIONING_DEVICE_ADMIN_COMPONENT_NAME`
- `android.app.extra.PROVISIONING_DEVICE_ADMIN_PACKAGE_DOWNLOAD_LOCATION`
- `android.app.extra.PROVISIONING_DEVICE_ADMIN_PACKAGE_CHECKSUM` (SHA-256 base64)
- `android.app.extra.PROVISIONING_LEAVE_ALL_SYSTEM_APPS_ENABLED`
- `android.app.extra.PROVISIONING_ADMIN_EXTRAS_BUNDLE` (nested JSON string/object per Java)
- Optional: WiFi SSID/password/security, `PROVISIONING_USE_MOBILE_DATA`, `PROVISIONING_SKIP_ENCRYPTION`, `qrParameters` expansion

**Nested admin extras** (`com.hmdm.*`):

| Key | When set |
|-----|----------|
| `com.hmdm.DEVICE_ID` | `deviceId` query non-empty |
| `com.hmdm.CONFIG` | `create=1` → configuration `qrCodeKey` |
| `com.hmdm.CUSTOMER` | `create=1` and multi-tenant |
| `com.hmdm.GROUP` | `create=1` and `group` query list |
| `com.hmdm.DEVICE_ID_USE` | `useId` = `imei` \| `serial` |
| `com.hmdm.BASE_URL` | Host from `BASE_URL` env (no `/rest` suffix) |
| `com.hmdm.SERVER_PROJECT` | Servlet context / deployment path |

**Errors**:

- `500` when key not found — log `configuration not found for key`
- `500` when main app missing URL — log explicit reason (not empty 200)

**Loopback rewrite**: `localhost` / `127.0.0.1` in APK URL → `BASE_URL` host (existing Go helper extended).

---

## GET `/json/{qrCodeKey}`

Same provisioning JSON string as PNG payload (without QR image). Query params identical except `size` ignored.

**Content-Type**: `application/json` body is the raw JSON text (Java compatibility).

---

## Repository fields required (`QRConfig`)

Extend port model to load: `launcherurl`, `wifi*`, `mobileenrollment`, `encryptdevice`, `qrparameters`, `adminextras`, `eventreceivingcomponent`, `apkhash` from `applicationversions`, application `pkg`.
