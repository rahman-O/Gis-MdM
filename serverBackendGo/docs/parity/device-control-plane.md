# Parity: Device Control Plane

**Feature**: `017-device-control-plane` | **Status**: v1 implemented (2026-05-23)

## Device tree (`/rest/private/device-tree`)

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/device-tree` | Implemented | List nodes + rootId; seeds «All devices» root |
| POST | `/device-tree/nodes` | Implemented | Create folder under parent |
| PUT | `/device-tree/nodes/:id` | Implemented | Rename, reorder, move (cycle check) |
| POST | `/device-tree/nodes/:id/delete` | Implemented | Relocate devices in subtree, then delete |

## Devices (extensions)

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| POST | `/devices/search` | Extended | `treeNodeId`, `includeDescendants` filters |
| POST | `/devices/:id/move-tree` | Implemented | Body `{ "treeNodeId": N }` |

## Enrollment identity (US2)

| Area | Status | Notes |
|------|--------|-------|
| `devices.agent_id` | Implemented | UUID, unique, default on insert |
| `devices.enrollment_state` | Implemented | `pending` → `enrolled` on create → `active` on sync touch |
| `devices.enrollment_route_id` | Implemented | FK to `enrollment_routes.id` (backfill from configurations) |
| `configurations.default_tree_node_id` | Implemented | Used on `CreateOnDemand` |
| `configurations.default_device_id_mode` | Implemented | Default `imei`; QR defaults when `deviceId` empty |
| QR `create=1` | Implemented | Default for empty `deviceId` in Go + React |
| UI enrollment column | Implemented | Devices table |

## Profiles (US3 draft)

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/profiles` | Implemented | List + device/route counts |
| POST | `/profiles` | Implemented | Create profile + draft v1 |
| GET | `/profiles/:id` | Implemented | Meta; auto-forks draft from published when missing |
| GET | `/profiles/:id/versions/:versionId` | Implemented | Full editor payload (Configuration shape) |
| PUT | `/profiles/:id/versions/:versionId` | Implemented | Save draft only |
| `GET /profiles/:id/impact` | Implemented | ≥50 devices → `requiresConfirmDialog` |
| `POST /profiles/:id/versions/:id/publish` | Implemented | Writes artifact + domain event |
| Sync from `profile_version_artifacts` | Implemented | Fallback to legacy configurations |
| `domain_events` push worker | Implemented | Debounced poll, `ProfilePublished` |

## Enrollment routes (US4)

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/enrollment-routes` | Implemented | List routes |
| POST | `/enrollment-routes` | Implemented | Create; `default_device_id_mode` default `imei` |
| GET/PUT | `/enrollment-routes/:id` | Implemented | Binding editor |
| GET | `/enrollment-routes/:id/qr` | Implemented | QR metadata |
| GET | `/enrollment-routes/options/published-profile-versions` | Implemented | Profile version picker |
| QR resolve by `qrcodekey` | Implemented | `enrollment_routes` + profile version settings |

## Onboarding (US6)

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/onboarding/status` | Implemented | Tree / published profile / route flags |

Frontend: dashboard `OnboardingChecklist`, `/onboarding` wizard; enrollment route create guarded when no published profile.

## Configurations alias (transition)

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/configurations/:id` | Alias | When `profiles.legacy_configuration_id` matches, returns profile version payload |
| PUT | `/configurations` | Alias | Saves to profile draft when legacy id maps |

Nav: **Profiles** + **Enrollment routes**; `/configurations` redirects to `/profiles`.

## Migrations

| Migration | Purpose |
|-----------|---------|
| `000019`–`000020` | Device tree |
| `000021` | Enrollment identity columns |
| `000022`–`000023` | Profiles + version junctions + configuration backfill |
| `000024` | Enrollment routes + route backfill |
| `000025`–`000026` | Profile artifacts + domain events |
| `000027` | Devices without `tree_node_id` → customer root |

## Tests (2026-05-23)

```text
go test ./internal/modules/device_tree/... ./internal/modules/profiles/... ./internal/modules/enrollment_routes/... ./internal/modules/sync/...  → pass
```

## §20 gates

Recorded in [`specs/017-device-control-plane/gates/sprint-017-v1-gates.md`](../../../specs/017-device-control-plane/gates/sprint-017-v1-gates.md).

## Java reference

- `com.hmdm.rest.resource.QRCodeResource`
- `com.hmdm.rest.resource.SyncResource`

## Blueprint gates

See [DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md](../../../DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md) §20 before each sprint merge.
