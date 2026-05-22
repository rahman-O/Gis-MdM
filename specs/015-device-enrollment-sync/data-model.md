# Data Model: Device Enrollment & Sync

**Branch**: `015-device-enrollment-sync` | **Date**: 2026-05-21

Logical entities and persistence touchpoints. Schema changes are **not required** for P0 fixes; optional migration only if new indexes needed for QR key lookup performance.

---

## Enrollment QR key

| Field / concept | Source | Notes |
|-----------------|--------|-------|
| `qrcodekey` | `configurations.qrcodekey` | Public lookup key for `/rest/public/qr/{key}` |
| `mainappid` | `configurations.mainappid` | FK → `applicationversions.id` for launcher APK |
| `launcherurl` | `configurations.launcherurl` | Overrides version URL in closed networks |
| `wifiSSID`, `wifiPassword`, `wifiSecurityType` | `configurations` | Embedded in provisioning JSON |
| `mobileenrollment`, `encryptdevice`, `qrparameters`, `adminextras` | `configurations` | Android extras |
| `eventreceivingcomponent` | `configurations` | Device admin receiver class |

**Validation rules**:

- QR generation MUST fail clearly when `mainappid` is null or version has empty `url` (unless `launcherurl` set).
- Key MUST be unique per customer (existing DB constraint assumed).

---

## Device (enrollment lifecycle)

| Field | Table | Notes |
|-------|-------|-------|
| `number` | `devices.number` | Agent device id (path param + `com.hmdm.DEVICE_ID`) |
| `configurationid` | `devices.configurationid` | Set on create-on-demand from QR key resolution |
| `customerid` | `devices.customerid` | From configuration or default customer |
| `lastupdate` | `devices.lastupdate` | `0` until first successful sync; used with `PREVENT_DUPLICATE_ENROLLMENT` |
| `enrolltime` | `devices.enrolltime` | Set on insert |
| `oldnumber` | `devices.oldnumber` | IMEI/serial migration path |
| `info` | `devices.info` | JSON telemetry from `POST /sync/info` |

**State transitions**:

```text
[no row] --QR scan + create=1 + POST sync--> [enrolled, lastupdate=0]
[enrolled, lastupdate=0] --first successful sync--> [active, lastupdate>0]
[active] --PREVENT_DUPLICATE + re-enroll--> rejected (DEVICE_EXISTS)
[unknown number] --POST info, single-tenant--> optional create-on-demand
```

---

## Device create options (agent POST body)

| JSON field | Maps to |
|------------|---------|
| `configuration` | **QR code key** (Java), not display name — fix lookup |
| `customer` | Customer name when multi-tenant + create on demand |
| `groups` | Group **names** → `devicegroups` insert |

---

## Sync configuration payload (`SyncResponse`)

| Section | Source tables |
|---------|----------------|
| Device identity | `devices` |
| Configuration meta | `configurations` (password hash, colors, permissive, `settingsjson` subset) |
| Applications | `configurationapplications` → `applications` / `applicationversions` |
| Files | `configurationfiles` → public URLs via `BuildPublicURL` |
| Application settings | `devicesapplicationsettings` (if loaded on sync) |

---

## Provisioning JSON (transient, not stored)

Built at request time from configuration + query params:

- Outer: Android Device Owner keys (`PROVISIONING_DEVICE_ADMIN_*`, WiFi, encryption flags).
- Inner `PROVISIONING_ADMIN_EXTRAS_BUNDLE`: `com.hmdm.BASE_URL`, `com.hmdm.SERVER_PROJECT`, `com.hmdm.DEVICE_ID`, `com.hmdm.CONFIG`, `com.hmdm.CUSTOMER`, `com.hmdm.GROUP`, `com.hmdm.DEVICE_ID_USE`.

---

## File download URLs

| Pattern | Served by |
|---------|-----------|
| `{BASE_URL}/files/{customerFilesDir}/{relativePath}` | **New** static handler (P0) |
| External URL on configuration file | Direct HTTP to third party |

---

## Device status (console “online”)

| Field | Table | Updated by |
|-------|-------|------------|
| `applicationsstatus`, `filesstatus` | `devicestatuses` | `POST /sync/info` upsert (014) |

Used for dashboard/search “installation” signals; enrollment UAT should verify row exists after first `info` post.
