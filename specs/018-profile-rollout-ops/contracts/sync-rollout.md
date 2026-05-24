# Contract: Sync & Rollout Status Integration

**Feature**: `018-profile-rollout-ops` | **Audience**: `sync` module implementers

## Effective profile resolution (per sync request)

**Input**: `deviceId`, `customerId`, `treeNodeId`, `enrollmentRouteId`

**Output**: `EffectiveProfile` `{ profileId, profileVersionId, artifactHash, source, enabled }`

**Order** (see [data-model.md](../data-model.md)):

1. Nearest `profile_tree_assignments` on ancestor path (deepest first) if profile enabled.
2. Else `enrollment_routes.profile_version_id` if route set and profile enabled.
3. Else legacy `configurationId` / artifact fallback from 017.

## After successful sync configuration delivery

1. Set `devices.applied_profile_version_id` to delivered `profileVersionId` when agent acknowledges config revision (hash match or version id in sync response path).
2. Call `RolloutStatusService.Recompute(deviceId)`.

## After device info ingest (`infojson` / applications status)

When device POSTs telemetry containing `applications[].status`:

1. Recompute rollout status (partial/failed/installed) without requiring full sync.
2. Update `profile_rollout_updated_at`.

## Push queue interaction

On assignment create/update or profile enable:

- Insert `domain_events` type `ProfileAssignmentChanged` (or reuse `ProfilePublished` with payload `{ profileId, treeNodeId }`).
- Worker batches push to affected devices (017 pattern).

## No agent contract change (v1)

- Status inference uses existing Headwind app status strings in device info.
- Document expected strings in parity doc during implementation (`OK`, `FAILURE`, `VERSION_MISMATCH` or project equivalents from Java).
