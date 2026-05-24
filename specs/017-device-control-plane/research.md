# Research: Device Control Plane

**Feature**: `017-device-control-plane` | **Date**: 2026-05-23

## R1 — Module boundaries

| Decision | Three new bounded contexts + extend existing modules |
|----------|-----------------------------------------------------|
| **Rationale** | Constitution I: one module per cohesive capability; avoids god-module `configurations`. |
| **Alternatives** | Single mega-module (rejected — violates maintainability NFR). |

| Module | Responsibility |
|--------|----------------|
| `device_tree` | CRUD tree nodes, move validation, path/depth maintenance |
| `profiles` | Profile identity, draft/publish versions, compile artifact, usage counts |
| `enrollment_routes` | Binding: QR key, tree placement, profile version, default device id mode |
| `devices` (extend) | `agent_id`, `tree_node_id`, `enrollment_route_id`, `enrollment_state` |
| `sync` (extend) | Load policy from `profile_version_artifacts` by route-follow semantics |
| `qrcode` (extend) | Force `create=1`; default `useId=imei` from route |
| `push` / `platform/push` | Async fanout via domain event outbox (Phase 4) |

## R2 — Tree storage: path + depth

| Decision | Materialized `path` (TEXT) + `depth` (INT) on `device_tree_nodes` |
|----------|---------------------------------------------------------------------|
| **Rationale** | Fast descendants/breadcrumbs without deep recursive CTE on every request; simpler than `ltree` extension dependency for v1. |
| **Alternatives** | `ltree` + GiST (deferred — can migrate later). Parent-only queries (rejected — slow at scale). |

**Maintenance**: On insert/move, recompute `path` for node and descendants in one transaction.

## R3 — Profile versioning & compile pipeline

| Decision | `profiles` (identity) + `profile_versions` (draft/published) + `profile_version_artifacts` (immutable JSON) |
|----------|-------------------------------------------------------------------------------------------------------------|
| **Rationale** | FR-005/007; prevents overwrite + enables rollback; sync reads artifact only (NFR-004). |
| **Alternatives** | Duplicate-on-save only (rejected — spec clarifications prefer publish workflow). In-place edit (rejected — incident risk). |

**Publish flow**: Edit draft → validate → compile to artifact (Headwind-shaped `SyncResponse` subset + apps/files junction snapshot) → emit `ProfilePublished` event.

**Route follow**: Device effective version = `enrollment_routes.profile_version_id` at sync time (FR-008a); updating route updates all devices on next sync.

## R4 — Migration from `configurations`

| Decision | Phased SQL migration: (1) add tables/columns (2) backfill v1 (3) switch sync read path behind flag (4) deprecate configurations UI |
|----------|------------------------------------------------------------------------------------------------------------------------------|
| **Rationale** | NFR-003 rollback; staging parity SC-005. |
| **Alternatives** | Big-bang cutover without flag (rejected). Dual-write forever (rejected — complexity). |

**Backfill**: Each `configurations` row → `profiles` + `profile_versions` v1 published + artifact compiled + `enrollment_routes` row + link devices.

**API alias**: `GET/PUT /rest/private/configurations/*` delegates to enrollment_routes or profiles during transition window (max one release).

## R5 — Enrollment QR defaults

| Decision | All route QR: `create=1` always; default `deviceIdUseMode=imei` per route (overridable to serial/request) |
|----------|--------------------------------------------------------------------------------------------------|
| **Rationale** | Spec clarifications; server+web only; Headwind supports `com.hmdm.DEVICE_ID_USE`. |
| **Alternatives** | Random device id (deferred — requires agent phase). Optional create checkbox (rejected — FR-002). |

## R6 — Delete tree node with devices

| Decision | API `POST /device-tree/nodes/:id/delete` with body `{ targetNodeId }` — moves all devices in subtree then deletes node |
|----------|------------------------------------------------------------------------------------------------------------------|
| **Rationale** | Spec clarification B; explicit UX contract. |
| **Alternatives** | Block delete only (rejected by product). Auto-move to parent (rejected — surprising). |

## R7 — Publish impact threshold

| Decision | Mandatory confirm dialog when affected device count ≥ **50** |
|----------|---------------------------------------------------------------|
| **Rationale** | Spec clarification; aligns SC-004. |
| **Alternatives** | Always confirm (too noisy); 10 devices (too aggressive for enterprise). |

## R8 — Push / notify storm prevention

| Decision | `domain_events` outbox table + worker enqueues `pushmessages` in batches (debounce 2s per profile publish) |
|----------|-----------------------------------------------------------------------------------------------------------|
| **Rationale** | FR-014; blueprint §17.10. |
| **Alternatives** | Synchronous notify in handler (rejected). |

## R9 — Frontend navigation

| Decision | Replace `Configurations` nav with `Profiles` + `Enrollment routes`; Devices layout = tree + table |
|----------|-----------------------------------------------------------------------------------------------|
| **Rationale** | FR-010a; immediate replacement per clarification. |
| **Alternatives** | Dual nav (rejected). |

## R10 — Agent compatibility (no mobile changes)

| Decision | Keep `configurationId` in `SyncResponse`; add optional `profileRevision`, `profileVersionId` |
|----------|-------------------------------------------------------------------------------------------|
| **Rationale** | FR-009; agent ignores unknown fields. |
| **Alternatives** | Break JSON shape (rejected). |

## R11 — Identity: `agent_id`

| Decision | `UUID` column on `devices`, immutable; internal FKs use `devices.id` / `agent_id` |
|----------|----------------------------------------------------------------------------------|
| **Rationale** | FR-013; `number` remains mutable with `oldnumber` for agent migration. |
| **Alternatives** | number as sole PK (rejected — rename breaks audit). |

## R12 — Testing strategy

| Decision | Golden tests: compile artifact parity vs legacy `BuildSyncResponse`; tree path recompute; publish impact counts |
|----------|---------------------------------------------------------------------------------------------------------------|
| **Rationale** | Constitution IV + NFR-002. |
