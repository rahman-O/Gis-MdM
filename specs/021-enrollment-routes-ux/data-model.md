# Data Model: Enrollment Routes UX

**Feature**: `021-enrollment-routes-ux` | **Date**: 2026-05-24

## Aggregate: EnrollmentRoute

Physical table: `enrollment_routes` (existing `000024`, extended `000030`).

Logical split:

```text
EnrollmentRoute (aggregate root)
├── EnrollmentRouteDefinition   (admin-edited)
└── EnrollmentRouteRuntimeState (system-managed)
```

---

## EnrollmentRouteDefinition

| Field | DB column | Type | Rules |
|-------|-----------|------|-------|
| id | `id` | int | PK |
| customerId | `customerid` | int | tenant scope |
| name | `name` | string(200) | unique per customer (ci) |
| description | `description` | text | optional |
| targetNodeId | `default_tree_node_id` | int FK → `device_tree_nodes` | required, valid for customer |
| deviceIdentityMode | `default_device_id_mode` | enum string | `imei` \| `serial` \| `request` |
| bootstrapIntent | `bootstrap_intent` | enum string | `stable` \| `specific` \| `latest` |
| bootstrapApplicationId | `bootstrap_application_id` | int FK → `applications` | required |
| bootstrapVersionId | `bootstrap_version_id` | int FK → `applicationversions` | required when intent=`specific` |
| resolvedMainAppVersionId | `mainappid` | int FK → `applicationversions` | written on Save after intent resolution |
| type | `type` | int | legacy, default 0 |
| containerPlacementAckAt | `container_placement_ack_at` | timestamptz | set when admin acks inheritable warning |

**Excluded from admin DTO**: `profile_version_id` (legacy compat only).

---

## EnrollmentRouteRuntimeState

| Field | DB column | Type | Rules |
|-------|-----------|------|-------|
| qrCodeKey | `qrcodekey` | string(200) | generated on first Save; unique globally (ci) |
| createdAt | `created_at` | timestamptz | |
| legacyConfigurationId | `legacy_configuration_id` | int | nullable, migration |

**Derived (not stored)**:
- `status`: `draft` (no row yet / no key) \| `active` (key present)
- Impact metrics via queries / `domain_events`

---

## EnrollmentContractPayload (QR logical)

Not a table — encoded in public QR JSON + client Pending preview.

| Field | Source |
|-------|--------|
| routeId | `enrollment_routes.id` |
| targetNodeId | `default_tree_node_id` |
| mainAppPackage | `applications.pkg` via version join |
| mainAppVersion | `applicationversions.version` / `versioncode` |
| deviceIdentityMode | `default_device_id_mode` |
| bootstrapFlags | `{ create: 1 }` + legacy-compatible query params |
| qrcodeKey | runtime only (Active) |

**MUST NOT include**: profile*, policy*.

---

## BootstrapAppIntent resolution

| Intent | Resolution on Save |
|--------|-------------------|
| `stable` | `applicationversions` where `applicationid = bootstrap_application_id` AND `is_recommended = true` |
| `latest` | max `versioncode` with installable `url` |
| `specific` | `bootstrap_version_id` |

Failure: no stable version → `error.enrollment_route.stable_version_missing`.

---

## TargetNodeSelection (read model)

From `GET /enrollment-routes/options/tree-nodes`:

| Field | Computation |
|-------|-------------|
| id | node id |
| name | node name |
| path | materialized path |
| placementKind | `locked` if no child nodes else `inheritable` |
| deviceCount | devices under node (existing tree query) |
| heavilyLoaded | `deviceCount >= threshold` |

---

## EnrollmentDeleteImpact (read model)

`GET /enrollment-routes/:id/impact`:

| Field | Definition |
|-------|------------|
| enrollingNowCount | devices with `enrollment_route_id`, `created_at` within 24h |
| historicalEnrolledCount | all devices with `enrollment_route_id` |
| activeQrScans7d | `domain_events` `enrollment_route.qr_viewed`, 7d window |

---

## State transitions (route lifecycle)

```text
[Create dialog] ──Save──► Active (qrcodekey assigned)
                │
                └──Cancel──► (no row)

Active ──Edit unsaved──► Active + Unsaved (UI only)
Active ──Save──► Active (resolved mainappid may change)
Active ──Delete──► removed (devices keep route reference historically optional SET NULL on delete — research: keep FK history via nullable, do not cascade delete devices)

```

**Note**: Deleting route: `DELETE` row; devices retain `enrollment_route_id` for audit OR set NULL — **Decision**: `ON DELETE SET NULL` on `devices.enrollment_route_id` (change FK in `000030`) so historical count still visible until reaper; impact API counts before delete.

---

## Migration `000030_enrollment_routes_ux` (planned)

- `bootstrap_intent`, `bootstrap_application_id`, `bootstrap_version_id`, `container_placement_ack_at`
- Backfill: `bootstrap_application_id` + intent `specific` from existing `mainappid` join
- `applicationversions.is_recommended`
- `domain_events` index optional: `(event_type, aggregate_id, created_at)`
- FK `devices.enrollment_route_id` → `ON DELETE SET NULL`

---

## Domain events

| event_type | aggregate_id | When |
|------------|--------------|------|
| `enrollment_route.qr_viewed` | route id string | public QR/JSON hit for route key |
| `enrollment_route.created` | route id | optional audit (v1 skip) |
