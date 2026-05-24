# Contract: Enrollment Routes Admin API (v2 — gateway)

**Feature**: `021-enrollment-routes-ux` | **Base path**: `/rest/private/enrollment-routes`  
**Replaces for UI**: [017 enrollment-routes-api.md](../../017-device-control-plane/contracts/enrollment-routes-api.md) profile-centric editor fields

**Envelope**: Headwind `{ status, message?, data? }`

---

## DTO: EnrollmentRouteView (list & detail)

**MUST NOT include**: `profileId`, `profileVersionId`, `profileVersionNumber`, `profile*`.

```json
{
  "id": 5,
  "name": "Warehouse gate",
  "description": "Front desk",
  "qrcodekey": "abc123",
  "targetNodeId": 12,
  "targetNodeName": "Warehouse",
  "targetNodePath": "/All devices/Warehouse",
  "targetPlacementKind": "inheritable",
  "containerPlacementAcknowledged": false,
  "deviceIdentityMode": "imei",
  "bootstrapIntent": "stable",
  "bootstrapApplicationId": 3,
  "bootstrapApplicationName": "MDM Agent",
  "bootstrapVersionId": null,
  "resolvedMainAppVersionId": 88,
  "resolvedVersionLabel": "1.4.2",
  "status": "active",
  "type": 0
}
```

| Field | Notes |
|-------|-------|
| status | `draft` \| `active` — server: active if `qrcodekey` set |
| containerPlacementAcknowledged | `container_placement_ack_at != null` |
| targetPlacementKind | `locked` \| `inheritable` |

---

## `GET /enrollment-routes`

List `EnrollmentRouteView[]`.

---

## `GET /enrollment-routes/:id`

Detail `EnrollmentRouteView`.

---

## `POST /enrollment-routes`

```json
{
  "name": "Warehouse gate",
  "description": null,
  "targetNodeId": 12,
  "deviceIdentityMode": "imei",
  "bootstrapIntent": "stable",
  "bootstrapApplicationId": 3,
  "bootstrapVersionId": null,
  "acknowledgeContainerPlacement": false
}
```

| Field | Required |
|-------|----------|
| name | yes |
| targetNodeId | yes |
| deviceIdentityMode | yes (default imei) |
| bootstrapIntent | yes |
| bootstrapApplicationId | yes |
| bootstrapVersionId | yes if intent=`specific` |
| acknowledgeContainerPlacement | yes if target is inheritable (or save fails with `error.enrollment_route.container_ack_required`) |

**Removed**: `profileVersionId` — ignored if sent by old clients.

**Errors**:
- `error.enrollment_route.stable_version_missing`
- `error.enrollment_route.tree_node_required`
- `error.enrollment_route.main_app_required`
- `error.duplicate.enrollment_route`

---

## `PUT /enrollment-routes/:id`

Partial update — same body fields as POST (all optional).

`qrcodekey` immutable.

---

## `DELETE /enrollment-routes/:id`

Hard delete route row. Devices: `enrollment_route_id` SET NULL (migration `000030`).

---

## `GET /enrollment-routes/:id/qr`

Active QR metadata — see [enrollment-contract-payload.md](./enrollment-contract-payload.md).

---

## `GET /enrollment-routes/:id/impact`

```json
{
  "enrollingNowCount": 2,
  "historicalEnrolledCount": 150,
  "activeQrScans7d": 42
}
```

Used by delete flow.

---

## `GET /enrollment-routes/options/tree-nodes`

```json
[
  {
    "id": 12,
    "name": "Warehouse",
    "path": "/All devices/Warehouse",
    "parentId": 1,
    "placementKind": "inheritable",
    "deviceCount": 120,
    "heavilyLoaded": false
  }
]
```

---

## `GET /enrollment-routes/options/bootstrap-apps`

```json
[
  {
    "applicationId": 3,
    "name": "MDM Agent",
    "package": "com.headwind.mdm",
    "versions": [
      {
        "versionId": 88,
        "version": "1.4.2",
        "versionCode": 142,
        "isRecommended": true,
        "isLatest": true
      }
    ]
  }
]
```

---

## Deprecated

| Endpoint | Action |
|----------|--------|
| `GET /enrollment-routes/options/published-profile-versions` | Return `410` or empty `[]`; document in parity — **not used by UI** |

---

## Permissions

Unchanged: `CanManageConfigurations()` (legacy `add_config` / `edit_config` / `del_config` on frontend).
