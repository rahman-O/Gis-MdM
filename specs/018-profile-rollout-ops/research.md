# Research: Profile Rollout & Operations

**Feature**: `018-profile-rollout-ops` | **Date**: 2026-05-23

## R1 — Tree assignment vs enrollment route

**Decision**: **Dual source** with explicit precedence at sync time:
1. **Enrollment route** (`devices.enrollment_route_id` → `enrollment_routes.profile_version_id`) defines policy for newly enrolled devices and remains the QR binding.
2. **Tree assignment** (`profile_tree_assignments` per `device_tree_node`) defines policy for devices **in that subtree** when an assignment exists; **nearest node wins** (deepest assignment on path from device `tree_node_id` to root).

**Rationale**: Spec requires applying profiles to the tree containing devices without removing route semantics from 017. Nearest-wins matches folder-level operations in enterprise MDM.

**Alternatives considered**:
- Tree-only (rejected — breaks route-follow and QR binding).
- Route overrides tree always (rejected — contradicts user request for tree application).

## R2 — Persisting rollout status

**Decision**: Add **device-level columns** (`target_profile_version_id`, `applied_profile_version_id`, `profile_rollout_status`, `profile_rollout_reason`, `profile_rollout_updated_at`) plus optional detail in existing `devices.info` telemetry. Recompute status on: assignment change, publish of assigned version, sync/info ingest, push notify completion.

**Rationale**: Device list and rollout panel need fast filters without joining large event tables. Headwind already reports per-app install status in `DeviceInfoView.applications[].status`.

**Alternatives considered**:
- Event-sourced rollout only (rejected — heavier queries for admin grid).
- Derive live on every GET without cache (rejected — violates SC-008 refresh budget at scale).

## R3 — Status enum mapping (v1, no agent changes)

**Decision**: Map to five admin-visible states:

| Status | Meaning | Primary signals |
|--------|---------|-----------------|
| `pending` | Waiting for sync/push | Target set; `applied` missing or stale; push queued |
| `installing` | Optional; merge with `pending` if no signal | Only if future agent reports in-progress |
| `installed` | Target matches applied; apps OK | `applied_profile_version_id == target` + no failed/mismatch apps in profile |
| `partial` | Version applied but apps/files incomplete | Any profile app `status` in info = failure/mismatch |
| `failed` | Cannot reach target | Sync error, version mismatch after grace, all profile apps failed |

**Rationale**: Spec FR-006/FR-013; reuses 016 install summary patterns and `DeviceApplication.Status` strings from agent info JSON.

**Alternatives considered**:
- Binary success/fail only (rejected — misses partial installs).

## R4 — Profile enable/disable

**Decision**: `profiles.enabled` boolean (default `true`). When `false`: block new tree assignments; sync resolver skips pushing that profile’s artifact (device keeps last applied state); enrollment route save returns warning if bound profile disabled; re-enable sets affected devices to `pending`.

**Rationale**: Spec assumption — no forced wipe on disable in v1.

## R5 — Version navigation in editor

**Decision**: New `GET /profiles/:id/versions` list endpoint; editor loads version by id; published versions read-only; **Fork draft from version** reuses existing `EnsureDraft` / `forkDraftFromVersion` in profile repo (017).

**Rationale**: 017 already forks draft from published when opening meta; 018 adds explicit version picker and fork-from-old-version action.

## R6 — Assignment impact threshold

**Decision**: Reuse **≥50 devices** confirm dialog pattern from 017 publish/impact (`CountImpact` extended to subtree device count for assignment target nodes).

**Rationale**: Spec FR-011 consistency.

## R7 — Sync loader extension

**Decision**: Extend `sync` `route_resolver` / artifact loader with `ResolveEffectiveProfileVersion(device)` calling `profiles/application` resolver (tree + route + enabled flag). Set `devices.target_profile_version_id` on assignment and after resolve.

**Rationale**: Constitution module boundaries — resolution logic in `profiles/application`, sync consumes port interface.

## R8 — Module flag

**Decision**: `MODULE_PROFILE_ROLLOUT_ENABLED` (default `true` when `MODULE_PROFILES_ENABLED=true`); no separate frontend flag.

**Rationale**: Feature is extension of profiles, not a fourth top-level module.
