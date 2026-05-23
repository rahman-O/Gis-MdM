# API Contract: Device Sync (enrollment fixes) — `/rest/public/sync`

**Base path**: `/rest/public/sync`  
**Java reference**: `com.hmdm.rest.resource.SyncResource`  
**Consumers**: Headwind MDM Android launcher

**Delta from Phase 7 (`008`)**: configuration resolution on create-on-demand; file URLs must be downloadable.

---

## POST `/configuration/{deviceId}`

**Body**: `DeviceCreateOptions`

```json
{
  "customer": "optional customer name",
  "configuration": "<qrCodeKey NOT display name>",
  "groups": ["group-name-1"]
}
```

**Fix**: `configuration` MUST resolve via `configurations.qrcodekey` first, then optional fallback by `name`.

**Success**: Envelope `data` = `SyncResponse`; headers `X-Response-Signature`, `X-IP-Address`.

**Errors** (unchanged keys):

| Condition | Message key |
|-----------|-------------|
| Bad signature (secure mode) | `error.permission.denied` |
| Unknown device, no create | `error.notfound.device` |
| Duplicate enroll | `error.duplicate.device` |

**Post-condition**: New device row with `lastupdate = 0`, correct `configurationid`, `devicegroups` populated.

---

## GET `/configuration/{deviceId}`

Periodic sync; same `SyncResponse` shape. File `url` fields MUST use reachable `BASE_URL/files/...` paths.

---

## POST `/info`

Unchanged contract; UAT requires `lastupdate` bump and `devicestatuses` upsert when wired.

---

## POST `/applicationSettings/{deviceId}`

Unchanged; used after enroll for per-app keys.

---

## SyncResponse minimum agent fields

| Field | Required for UAT |
|-------|------------------|
| `deviceId` | yes |
| `configurationId` | yes |
| `applications[]` with `pkg`, `url`, `version` | yes if config has apps |
| `files[]` with `devicePath`, `url` | yes if config has files |
| `password` | if configuration password set (MD5 uppercase hex) |
