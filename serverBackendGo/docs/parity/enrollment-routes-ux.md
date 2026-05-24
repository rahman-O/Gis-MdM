# Parity: Enrollment Routes UX (021)

**Feature**: `021-enrollment-routes-ux` | **Status**: v1 implemented  
**Spec**: `specs/021-enrollment-routes-ux`  
**Admin API**: `/rest/private/enrollment-routes` — see `specs/021-enrollment-routes-ux/contracts/enrollment-routes-admin-api.md`

---

## Java Reference

The Go enrollment routes module replaces enrollment-related functionality that was previously embedded in the HMDM Java **Configuration** domain:

| Java class | Role | Go replacement |
|------------|------|----------------|
| `com.hmdm.rest.resource.ConfigurationResource` | Admin CRUD for configurations (which bundled enrollment QR settings) | `internal/modules/enrollment_routes/adapter/http/handler.go` |
| `com.hmdm.rest.resource.QRCodeResource` | Public QR code generation (`/rest/public/qr/{key}`) | `internal/modules/qrcode/` (public handler); enrollment route key resolution in `enrollment_routes` service |
| `com.hmdm.persistence.domain.Configuration` | Domain entity holding `qrCodeKey`, `mainAppId`, WiFi provisioning, admin extras | `internal/modules/enrollment_routes/domain/route.go` — `EnrollmentRouteDefinition` + `EnrollmentRouteRuntimeState` |
| `com.hmdm.persistence.ConfigurationDAO` | DB access for configurations including QR key lookup (`getConfigurationByQRCodeKey`) | `internal/modules/enrollment_routes/adapter/persistence/postgres/route_repo.go` |
| `com.hmdm.persistence.mapper.ConfigurationMapper` | MyBatis SQL mapping for `configurations` table | Go repo with raw SQL / sqlx |
| `com.hmdm.persistence.UnsecureDAO` | Public (unauthenticated) config lookup by QR key for device enrollment | `internal/modules/qrcode/` public handler + enrollment route repo |
| `com.hmdm.guice.module.PrivateRestModule` | Registers `ConfigurationResource` at `/private/configurations` | Gin router registration in `handler.go` at `/private/enrollment-routes` |

### Key behavioral differences from Java

| Concern | Java (HMDM) | Go (021) |
|---------|-------------|----------|
| Enrollment = Configuration | QR key, main app, and enrollment settings are fields on `Configuration` entity | Enrollment route is a **separate bounded context** — own module, own table, own API |
| Profile coupling | `Configuration` IS the policy; enrollment inherently tied to profile/policy | Enrollment route has **zero** profile/policy dependency; profile assignment is a separate system |
| QR payload | Built from `Configuration` fields (mainAppId, wifiSSID, adminExtras, etc.) | Built from `EnrollmentRouteDefinition` + resolved bootstrap app version |
| Bootstrap app resolution | `mainAppId` points directly to `applicationversions.id` (static) | Intent-based: `stable` / `specific` / `latest` — resolved at save time |
| Delete guard | `ConfigurationReferenceExistsException` blocks delete when devices assigned | Multi-dimensional impact check; delete **allowed** even with historical devices (SET NULL) |

---

## Breaking DTO Changes

### Fields removed from admin API response

The `EnrollmentRouteView` DTO returned by `GET /enrollment-routes` and `GET /enrollment-routes/:id` **no longer includes**:

| Removed field | Was in Java `Configuration` | Reason |
|---------------|----------------------------|--------|
| `profileId` | Implicit (configuration = profile) | Profile decoupled from enrollment (FR-003) |
| `profileVersionId` | `mainAppId` was version-coupled | Replaced by intent-based bootstrap resolution |
| `profileVersionNumber` | Derived from profile version | Not applicable — enrollment has no profile concept |
| `profile*` (any) | Various | Strict vocabulary rule: no profile/policy in enrollment domain |

### Fields removed from create/update request

| Removed field | Migration note |
|---------------|---------------|
| `profileVersionId` | Ignored if sent by old clients; not validated or stored in enrollment route |

### Fields added to admin API response (`EnrollmentRouteView`)

| New field | Type | Description |
|-----------|------|-------------|
| `targetNodeId` | int | Device tree node where enrolled devices are placed |
| `targetNodeName` | string | Human-readable node name |
| `targetNodePath` | string | Full breadcrumb path (e.g., `/All devices/Warehouse`) |
| `targetPlacementKind` | string | `locked` \| `inheritable` — computed from tree structure |
| `containerPlacementAcknowledged` | bool | Admin acknowledged inheritable (container) target |
| `bootstrapIntent` | string | `stable` \| `specific` \| `latest` |
| `bootstrapApplicationId` | int | Application identity for bootstrap |
| `bootstrapApplicationName` | string | Human-readable app name |
| `bootstrapVersionId` | int \| null | Pinned version (when intent = `specific`) |
| `resolvedMainAppVersionId` | int | Concrete version ID resolved from intent |
| `resolvedVersionLabel` | string | Human version string (e.g., `1.4.2`) |
| `deviceIdentityMode` | string | `imei` \| `serial` \| `request` |
| `status` | string | `draft` \| `active` |

### Fields added to create/update request

| New field | Required | Description |
|-----------|----------|-------------|
| `targetNodeId` | yes | Target tree node |
| `deviceIdentityMode` | yes (default `imei`) | Identity mode |
| `bootstrapIntent` | yes | `stable` \| `specific` \| `latest` |
| `bootstrapApplicationId` | yes | Bootstrap app |
| `bootstrapVersionId` | yes if intent=`specific` | Pinned version |
| `acknowledgeContainerPlacement` | yes if target is inheritable | Container ack |

---

## Endpoint Parity

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/enrollment-routes` | **Done** | List `EnrollmentRouteView[]` — no profile fields |
| POST | `/enrollment-routes` | **Done** | Create; `profileVersionId` ignored if sent |
| GET | `/enrollment-routes/:id` | **Done** | Detail view |
| PUT | `/enrollment-routes/:id` | **Done** | Partial update; `qrcodekey` immutable |
| DELETE | `/enrollment-routes/:id` | **Done** | Hard delete; devices FK SET NULL |
| GET | `/enrollment-routes/:id/qr` | **Done** | Active QR metadata with resolved package/version |
| GET | `/enrollment-routes/:id/impact` | **New** | Three-dimensional delete impact metrics |
| GET | `/enrollment-routes/options/tree-nodes` | **New** | Tree nodes with `placementKind`, `deviceCount`, `heavilyLoaded` |
| GET | `/enrollment-routes/options/bootstrap-apps` | **New** | Apps with versions, `isRecommended`, `isLatest` flags |
| GET | `/enrollment-routes/options/published-profile-versions` | **Deprecated** | Returns empty `[]`; not used by new UI |

---

## Migration

`000030_enrollment_routes_ux` adds:

| Change | Table | Details |
|--------|-------|---------|
| `bootstrap_intent` | `enrollment_routes` | VARCHAR(20) NOT NULL DEFAULT `'stable'` |
| `bootstrap_application_id` | `enrollment_routes` | FK → `applications(id)` |
| `bootstrap_version_id` | `enrollment_routes` | FK → `applicationversions(id)` |
| `container_placement_ack_at` | `enrollment_routes` | TIMESTAMPTZ NULL |
| `is_recommended` | `applicationversions` | BOOLEAN NOT NULL DEFAULT FALSE; partial unique per app |
| Devices FK | `devices.enrollment_route_id` | ON DELETE SET NULL (was RESTRICT) |
| Backfill | `enrollment_routes` | Populates `bootstrap_application_id` from existing `mainappid` → `applicationversions.applicationid` |

---

## Config

| Variable | Default | Purpose |
|----------|---------|---------|
| `ENROLLMENT_TREE_HEAVY_DEVICE_THRESHOLD` | 500 | Tree option `heavilyLoaded` flag threshold |

---

## Domain Events

| Event | Emitted when |
|-------|-------------|
| `enrollment_route.qr_viewed` | Public QR/JSON hit (`GET /rest/public/qr/{key}` or `/json/{key}`) when key resolves to enrollment route |

Used by `GET /:id/impact` → `activeQrScans7d` metric.

---

## Client Migration Guide

### For API consumers upgrading from 017 enrollment routes

1. **Remove `profileVersionId` from request bodies** — field is ignored; sending it will not cause errors but has no effect.
2. **Add required fields to create/update**: `targetNodeId`, `bootstrapIntent`, `bootstrapApplicationId`, `deviceIdentityMode`.
3. **Update response parsing**: `EnrollmentRouteView` shape has changed (see "Fields added" above); old profile-related fields are absent.
4. **Profile version picker**: `GET /options/published-profile-versions` returns empty array — remove any UI dependency on this endpoint.
5. **Delete flow**: Use `GET /:id/impact` before delete to show impact dimensions; delete no longer blocked by device references.
6. **QR metadata**: `GET /:id/qr` response now includes `mainAppPackage`, `mainAppVersion`, `mainAppVersionCode`, `targetNodeId` — no profile fields.

### Backward compatibility

- Legacy `profile_version_id` column remains in DB for sync resolver compatibility; **never** exposed in admin API DTO.
- Public QR paths (`/rest/public/qr/{key}`) unchanged — existing printed QR codes continue to work.
- `qrcodekey` is immutable after first save; no rotation in v1.

---

## Tests

```bash
go test ./internal/modules/enrollment_routes/... → pass
```

Covers: intent resolution (stable/specific/latest), impact metrics, tree options placement kind, create/update without profile fields.
