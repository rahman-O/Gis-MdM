# Contract: Enrollment Contract Payload (QR)

**Feature**: `021-enrollment-routes-ux` | **Audience**: Admin UI (Pending preview), public QR agent, parity docs

## Principles

- Payload is **onboarding-only** — no profile/policy fields.
- **Pending** (client): visual + JSON preview; **not scannable** for production enrollment.
- **Active** (server): public URL uses real `qrcodekey`.

---

## Logical contract (canonical)

```json
{
  "routeId": 12,
  "targetNodeId": 45,
  "mainAppPackage": "com.example.mdm.agent",
  "mainAppVersion": "1.4.2",
  "mainAppVersionCode": 142,
  "deviceIdentityMode": "imei",
  "bootstrapFlags": {
    "create": true
  }
}
```

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| routeId | int | Active only | Omitted or `0` in Pending preview |
| targetNodeId | int | yes | `default_tree_node_id` |
| mainAppPackage | string | yes | From `applications.pkg` |
| mainAppVersion | string | yes | Human version string |
| mainAppVersionCode | int | recommended | Agent install selection |
| deviceIdentityMode | string | yes | `imei` \| `serial` \| `request` |
| bootstrapFlags.create | bool | yes | MUST be `true` for new device enrollment |
| qrcodeKey | string | Active only | Public scan key |

**Forbidden keys**: `profileId`, `profileVersionId`, `configurationId`, `policy*`.

---

## Public QR HTTP mapping (legacy parity)

Unchanged paths; query params carry enrollment semantics:

| Contract field | HTTP (public) |
|----------------|---------------|
| qrcodeKey | path `/rest/public/qr/{key}` |
| bootstrapFlags.create | `create=1` |
| deviceIdentityMode | `useId=imei` \| `useId=serial` \| (empty + request mode) |
| mainApp* | resolved server-side from `mainappid` on route row |

Admin **Active** preview uses same URLs as production once `qrcodekey` exists.

---

## Pending preview (client-only)

```json
{
  "routeId": 0,
  "targetNodeId": 45,
  "mainAppPackage": "com.example.mdm.agent",
  "mainAppVersion": "1.4.2",
  "mainAppVersionCode": 142,
  "deviceIdentityMode": "imei",
  "bootstrapFlags": { "create": true },
  "_preview": true
}
```

Rules:

- `_preview: true` MUST NOT appear in server-generated Active JSON.
- UI MUST disable copy/download of production URL while `_preview` or before first Save.
- QR image MAY render from local encoder with watermark **“Pending — save to activate”**.

---

## `GET /rest/private/enrollment-routes/:id/qr` (Active metadata)

Response (no profile fields):

```json
{
  "qrcodekey": "abc123xyz",
  "defaultDeviceIdMode": "imei",
  "resolvedMainAppVersionId": 88,
  "mainAppPackage": "com.example.mdm.agent",
  "mainAppVersion": "1.4.2",
  "mainAppVersionCode": 142,
  "targetNodeId": 45,
  "contract": { }
}
```

`contract` embeds logical fields above for admin display.

---

## Validation

- Save MUST fail if intent `stable` and no `is_recommended` version for package.
- Save MUST fail if `targetNodeId` invalid for tenant.
- Save MUST NOT require profile fields.
