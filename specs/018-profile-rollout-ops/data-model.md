# Data Model: Profile Rollout & Operations

**Feature**: `018-profile-rollout-ops` | **Date**: 2026-05-23 | **Depends on**: [017 data-model](../017-device-control-plane/data-model.md)

## Schema changes (migration `000028_profile_rollout_ops`)

### `profiles` (alter)

| Column | Type | Notes |
|--------|------|-------|
| `enabled` | BOOLEAN NOT NULL DEFAULT TRUE | FR-009/010; disable stops new push |

### `profile_tree_assignments` (new)

| Column | Type | Notes |
|--------|------|-------|
| `id` | SERIAL PK | |
| `customerid` | INT NOT NULL FK → customers | tenant scope |
| `profile_id` | INT NOT NULL FK → profiles | denormalized for FK checks |
| `profile_version_id` | INT NOT NULL FK → profile_versions | must be `published` |
| `tree_node_id` | INT NOT NULL FK → device_tree_nodes | folder receiving assignment |
| `created_at` | TIMESTAMPTZ NOT NULL DEFAULT now() | |
| `created_by` | INT NULL FK → users | audit |

**Constraints**:
- `UNIQUE (customerid, tree_node_id)` — one active assignment per folder.
- `CHECK` via application: `profile_version_id` references version with `status = 'published'` and matching `profile_id`.

**Rules**:
- Assigning to node N affects all devices where `devices.tree_node_id` is N or descendant (path prefix match on `device_tree_nodes.path`).
- Child assignment overrides parent (nearest wins).

### `devices` (alter)

| Column | Type | Notes |
|--------|------|-------|
| `target_profile_version_id` | INT NULL FK → profile_versions | resolved target after assignment/sync resolve |
| `applied_profile_version_id` | INT NULL FK → profile_versions | last version confirmed on device |
| `profile_rollout_status` | VARCHAR(20) NULL | `pending` \| `installing` \| `installed` \| `partial` \| `failed` |
| `profile_rollout_reason` | TEXT NULL | admin-facing summary |
| `profile_rollout_updated_at` | TIMESTAMPTZ NULL | last status recompute |

**Indexes**:
- `(customerid, profile_rollout_status)` for rollout dashboard filters.
- `(target_profile_version_id)` for version impact queries.

## Domain entities

### `ProfileTreeAssignment`

Links a published version to a tree folder; includes impact preview counts (transient).

### `DeviceRolloutState`

View model per device: device id, name, tree path, target version, applied version, status, reason, last sync time.

### `EffectiveProfileResolution`

Result of resolver: `profileVersionId`, `source` (`tree` \| `route` \| `none`), `profileId`, `enabled`.

## State transitions (`profile_rollout_status`)

```text
(pending) ──sync success + version match + apps OK──► (installed)
(pending) ──sync success + version match + app issues──► (partial)
(pending) ──sync error / mismatch timeout──► (failed)
(installed|partial|failed) ──assignment or publish change──► (pending)
(any) ──profile disabled──► (pending) frozen for push; display may show last applied
```

**Publish side-effect**: When a published version is referenced by any `profile_tree_assignments.profile_version_id`, publish of a **new** version does not auto-retarget assignments — admin must re-assign or use optional «update assignment to new version» action (v1: re-assign same folder to new version id).

## Resolution algorithm (application layer)

For device `D` with `tree_node_id = T`:

1. Walk ancestors of `T` (including T) ordered by **depth DESC**.
2. First `profile_tree_assignments` row on that node where `profiles.enabled = true` → candidate **tree version**.
3. If `D.enrollment_route_id` set, load route’s `profile_version_id` → **route version** (if profile enabled).
4. **Effective target** for sync:
   - If tree candidate exists → use tree candidate (spec: tree assignment updates devices in folder).
   - Else if route version exists → use route version.
   - Else → no profile target (`null`).
5. Write `devices.target_profile_version_id` when changed; set `profile_rollout_status = pending` if target changed.

## Status recompute (application layer)

Inputs: `target_profile_version_id`, `applied_profile_version_id`, `profile_version_artifacts.artifact_hash`, device `info.applications` / `info.files` for packages/paths in target version.

1. No target → clear status or `pending` with reason «No assignment».
2. Target ≠ applied (or applied null) and last sync &lt; grace → `pending`.
3. Target ≠ applied after grace (e.g. 15 min) → `failed` («Version mismatch»).
4. Target = applied → compare each profile app in version to device info app status:
   - All OK → `installed`.
   - Some mismatch/fail → `partial` with reason listing packages.
   - All fail → `failed`.

Invoke recompute from: `sync` post-handler, `devices` info update, assignment service, profile enable/disable, publish hook (when version id referenced).

## API surface (summary)

See [contracts/profile-rollout-api.md](./contracts/profile-rollout-api.md).

## Parity doc

`serverBackendGo/docs/parity/profile-rollout-ops.md` (created during implementation).

## Java reference

- `com.hmdm.persistence.domain.Device` / device search list columns for install status
- `DeviceApplicationsStatus` patterns in REST device list (install status column)
- No new public agent endpoints in v1
