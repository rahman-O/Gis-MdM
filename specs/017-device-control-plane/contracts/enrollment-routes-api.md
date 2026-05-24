# Contract: Enrollment Routes API

**Feature**: `017-device-control-plane` | **Audience**: Admin React client

**Product label**: «مسار التسجيل» / **Enrollment route** (never «Policy» in UI)

## Endpoints

### `GET /enrollment-routes`

List routes with binding summary.

```json
{
  "id": 5,
  "name": "Store front QR",
  "qrcodekey": "abc123",
  "profileId": 10,
  "profileVersionId": 46,
  "profileVersionNumber": 4,
  "defaultTreeNodeId": 2,
  "defaultTreeNodeName": "Stores",
  "defaultDeviceIdMode": "imei",
  "mainAppId": 1
}
```

### `GET /enrollment-routes/:id`

Detail for editor (no restrictions/apps tabs — binding fields only).

### `POST /enrollment-routes`

**Body**:

```json
{
  "name": "Store front QR",
  "description": "",
  "profileVersionId": 46,
  "defaultTreeNodeId": 2,
  "defaultDeviceIdMode": "imei",
  "mainAppId": 1
}
```

**Validation**: `profileVersionId` MUST reference `status=published`.

### `PUT /enrollment-routes/:id`

Update binding. Changing `profileVersionId` triggers route-follow on next device sync (FR-008a).

### `GET /enrollment-routes/:id/qr`

Returns provisioning payload metadata (same as legacy QR resource).

**Guarantees**:

- `create=1` always in generated QR JSON
- `deviceIdUseMode` from `defaultDeviceIdMode` (default `imei`)

### Public (unchanged paths)

- `GET /rest/public/qr/:key` — resolve route by `qrcodekey`
- `POST /rest/public/sync/info` — enroll with `create=1`

## Editor UX contract

Fields on route editor:

| Field | Required |
|-------|----------|
| Name | yes |
| Published profile version (picker) | yes |
| Default tree folder | yes |
| Default device id mode | yes (default imei) |
| Main app | yes for QR |

**Must NOT** show: Restrictions, Applications, Design, Files tabs (see [frontend-control-plane-ux.md](./frontend-control-plane-ux.md)).

**Inline help**: «Device restrictions and apps are configured in the Profile, not here.»

## Alias

`GET/POST /configurations` list may return enrollment routes during transition with `_legacyType: "configuration"` stripped in UI.
