# Data Model: Device Control Plane

**Feature**: `017-device-control-plane` | **Date**: 2026-05-23

## New tables

### `device_tree_nodes`

| Column | Type | Notes |
|--------|------|-------|
| `id` | SERIAL PK | |
| `customerid` | INT FK → customers | tenant scope |
| `parent_id` | INT FK → device_tree_nodes NULL | NULL = virtual root per customer (or single root row) |
| `name` | VARCHAR(200) | unique per (customerid, parent_id, lower(name)) |
| `sort_order` | INT DEFAULT 0 | sibling order |
| `path` | TEXT NOT NULL | e.g. `/1/4/9/` |
| `depth` | INT NOT NULL DEFAULT 0 | |
| `created_at` | TIMESTAMPTZ | |

**Rules**: No cycles on move; delete requires `target_node_id` relocation (API).

### `profiles`

| Column | Type | Notes |
|--------|------|-------|
| `id` | SERIAL PK | may equal legacy `configurations.id` after migration |
| `customerid` | INT | |
| `name` | VARCHAR | unique per customer |
| `description` | TEXT | |
| `created_at` | TIMESTAMPTZ | |

### `profile_versions`

| Column | Type | Notes |
|--------|------|-------|
| `id` | SERIAL PK | |
| `profile_id` | INT FK → profiles | |
| `version_number` | INT | monotonic per profile |
| `status` | VARCHAR(20) | `draft` \| `published` \| `archived` |
| `settingsjson` | JSONB | policy map + `policyLocks` |
| scalars | … | password, colors, permissive, mainappid, contentappid — same split as configurations |
| `published_at` | TIMESTAMPTZ NULL | |
| `published_by` | INT NULL | user id |

**Junction** (moved from configuration_*):

- `profile_applications` (was `configurationapplications`)
- `profile_files` (was `configurationfiles`)
- `profile_application_settings` (was `configurationapplicationsettings`)
- `profile_application_parameters` (was `configurationapplicationparameters`)

### `profile_version_artifacts`

| Column | Type | Notes |
|--------|------|-------|
| `profile_version_id` | INT PK FK → profile_versions | |
| `artifact_json` | JSONB NOT NULL | pre-built agent sync payload body |
| `artifact_hash` | VARCHAR(64) | SHA-256 for revision compare |
| `compiled_at` | TIMESTAMPTZ | |

### `enrollment_routes`

| Column | Type | Notes |
|--------|------|-------|
| `id` | SERIAL PK | may map 1:1 from configurations.id |
| `customerid` | INT | |
| `name` | VARCHAR | |
| `description` | TEXT | |
| `qrcodekey` | VARCHAR(200) UNIQUE | |
| `mainappid` | INT | QR eligibility |
| `profile_version_id` | INT FK → profile_versions | must be published |
| `default_tree_node_id` | INT FK → device_tree_nodes | |
| `default_device_id_mode` | VARCHAR(20) DEFAULT 'imei' | `imei` \| `serial` \| `request` |
| `type` | INT | WORK/COMMON legacy |

### `domain_events` (Phase 4)

| Column | Type | Notes |
|--------|------|-------|
| `id` | BIGSERIAL PK | |
| `event_type` | VARCHAR(64) | e.g. `ProfilePublished` |
| `aggregate_id` | VARCHAR(64) | profile id |
| `payload` | JSONB | |
| `created_at` | TIMESTAMPTZ | |
| `processed_at` | TIMESTAMPTZ NULL | |

### `device_events` (Phase 4, optional v1 minimal)

| Column | Type | Notes |
|--------|------|-------|
| `id` | BIGSERIAL PK | |
| `device_id` | INT FK | |
| `event_type` | VARCHAR(64) | enrolled, moved_tree, profile_applied, … |
| `payload` | JSONB | |
| `created_at` | TIMESTAMPTZ | |

## Extended `devices`

| Column | Type | Notes |
|--------|------|-------|
| `agent_id` | UUID NOT NULL UNIQUE | immutable internal identity |
| `tree_node_id` | INT FK → device_tree_nodes | required after migration |
| `enrollment_route_id` | INT FK → enrollment_routes | replaces semantic of configurationid for binding |
| `enrollment_state` | VARCHAR(20) | pending, enrolled, active, stale, archived |
| `configurationid` | INT | **retained** for agent parity / migration; sync may still expose as `configurationId` |

**Note**: During transition, `configurationid` may mirror `enrollment_route_id` or legacy id until cutover complete.

## Relationships

```text
profiles 1──* profile_versions 1──1 profile_version_artifacts
profiles 1──* profile_applications (via version or profile_id per design — prefer version-scoped junction on publish)

enrollment_routes *──1 profile_versions (published)
enrollment_routes *──1 device_tree_nodes (default placement)

devices *──1 enrollment_routes
devices *──1 device_tree_nodes
devices *──1 customers

device_tree_nodes *──1 device_tree_nodes (parent)
```

**Junction design choice**: On publish, snapshot apps/files/settings into version-scoped rows OR embed in `artifact_json` only for sync; junction tables copied at publish time from draft (research R3).

## State transitions

### Profile version

```text
draft → (publish) → published
published → (archive optional) → archived
new draft forked from published (edit)
```

### Device enrollment_state

```text
(pending) → first successful POST enroll → enrolled
enrolled → POST /sync/info → active
active → (job: no heartbeat) → stale
any → admin action → archived
```

### Enrollment route

```text
created (draft binding) → valid when profile_version published + tree node set → QR enabled
```

## Validation rules

| Rule | Layer |
|------|-------|
| Unique `qrcodekey` | enrollment_routes |
| Cannot publish profile without compile success | profiles/application |
| Cannot save route without published `profile_version_id` | enrollment_routes |
| Tree move: no cycle | device_tree/application |
| Delete node: `targetNodeId` required if devices in subtree | device_tree/application |
| Publish confirm if devices ≥ 50 | profiles/application + frontend |
| Tenant filter on all queries | all repos |
| IMEI default on new route | enrollment_routes default + qrcode builder |

## Sync resolution (runtime)

```text
device.enrollment_route_id
  → enrollment_routes.profile_version_id  (route follow — may differ from device row cache until sync)
  → profile_version_artifacts.artifact_json
  → SyncResponse (+ configurationId = route.id or legacy mapping)
```

## Migration mapping

| Legacy | Target |
|--------|--------|
| `configurations` row | `profiles` + `profile_versions` v1 published + artifact |
| same row | `enrollment_routes` (qrcodekey, mainappid, tree → root) |
| `devices.configurationid` | `devices.enrollment_route_id` + retain `configurationid` |
| — | `devices.agent_id` = gen_random_uuid() |
| — | `devices.tree_node_id` = customer root |
