# Implementation Plan: Profile Rollout & Operations

**Branch**: `018-profile-rollout-ops` | **Date**: 2026-05-23 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/018-profile-rollout-ops/spec.md`  
**Builds on**: [017-device-control-plane](../017-device-control-plane/plan.md) (tree, profiles, publish, routes, sync artifacts)

## Summary

Extend the control plane so **published profile versions can be assigned to device-tree folders**, admins can **navigate profile versions** in the editor, and **per-device rollout status** (pending / installed / partial / failed) is visible and kept fresh from sync and device telemetry. **Enable/disable** at profile level pauses policy push without deleting versions or assignments.

**Approach**: Extend `profiles` module (assignments, resolver, rollout status, enable/disable); add migration `000028`; hook `sync` and `devices` info ingest to recompute status; extend React `features/profiles` with version switcher, assignment panel, and rollout grid. No Android agent changes in v1.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo`), TypeScript/React 18 (`frontend`)

**Primary Dependencies**: Existing `profiles`, `device_tree`, `devices`, `sync`, `push` (domain events), 016 configuration mapper inside profile compile

**Storage**: PostgreSQL — `profile_tree_assignments`, `profiles.enabled`, rollout columns on `devices` ([data-model.md](./data-model.md))

**Testing**: `go test` on `profiles/application` (resolver, assignment impact, status recompute); sync resolver tests; frontend smoke per [quickstart.md](./quickstart.md)

**Target Platform**: Go API `:8080`; Vite admin; Headwind agent unchanged

**Project Type**: Web application (backend + frontend)

**Performance Goals**: Assignment impact count &lt; 2s for 5k devices (indexed `tree_node_id` + `path`); rollout list page &lt; 500ms p95; status recompute batch &lt; 100ms per device

**Constraints**: Headwind JSON envelope; layered modules; tree assignment nearest-wins; ≥50 device confirm (017 pattern); no staged % rollout in v1

**Scale/Scope**: ~25–35 Go files (mostly under `profiles/`); ~8–12 React files; 1 migration; parity `profile-rollout-ops.md`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Extends `profiles`; integrates via `port` into `sync`/`devices`; MIGRATION/NEXT_STEPS row |
| **II. Layered Clean** | ✅ | Resolver/rollout in `application/`; SQL in `adapter/persistence` |
| **III. API Parity** | ✅ | New private routes documented; sync public paths unchanged |
| **IV. Testable Delivery** | ✅ | Unit tests + quickstart sprints |
| **V. Simplicity** | ✅ | Device columns for status; single assignment per tree node |
| **VI. Security** | ✅ | Tenant + config permissions on all new routes |
| **VII. Observability** | ✅ | `MODULE_PROFILE_ROLLOUT_ENABLED`; stable error keys |

**Post-design**: All gates ✅. Resolved in [research.md](./research.md).

## Project Structure

### Documentation (this feature)

```text
specs/018-profile-rollout-ops/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── profile-rollout-api.md
│   ├── profile-editor-versions-ux.md
│   └── sync-rollout.md
└── tasks.md              # (/speckit-tasks)
```

### Source Code

```text
serverBackendGo/
├── db/migrations/
│   └── 000028_profile_rollout_ops.up.sql
├── internal/modules/profiles/
│   ├── domain/assignment.go, rollout.go
│   ├── application/
│   │   ├── assignment.go
│   │   ├── resolver.go          # effective profile for device
│   │   ├── rollout_status.go    # recompute pending/installed/...
│   │   └── enable.go
│   ├── port/resolver.go
│   └── adapter/
│       ├── http/assignment_handler.go, rollout_handler.go
│       └── persistence/postgres/assignment_repo.go, rollout_repo.go
├── internal/modules/sync/application/
│   └── effective_profile.go     # calls profiles resolver
├── internal/modules/devices/application/
│   └── info_rollout_hook.go     # recompute on info ingest
└── docs/parity/profile-rollout-ops.md

frontend/src/features/profiles/
├── ProfileVersionSelect.tsx
├── ProfileTreeAssignmentPanel.tsx
├── ProfileRolloutStatusPanel.tsx
├── ProfileDisableBanner.tsx
└── profileRolloutService.ts
```

**Structure Decision**: No new top-level module; rollout is a vertical slice on `profiles` with narrow hooks in `sync` and `devices` via interfaces.

## Implementation Phases (for `/speckit-tasks`)

### Phase 1 — Schema & domain (P1)

1. Migration `000028`: `profiles.enabled`, `profile_tree_assignments`, device rollout columns + indexes.
2. Domain types + repository methods (list versions, assignments CRUD).
3. `MODULE_PROFILE_ROLLOUT_ENABLED` in config.

### Phase 2 — Version navigation API + UI (P1) — US2

1. `GET /profiles/:id/versions`, `POST .../fork-draft`.
2. `ProfileVersionSelect` + read-only published view + unsaved guard.
3. Route `/profiles/:id/versions/:versionId/edit`.

### Phase 3 — Tree assignment (P1) — US1

1. Assignment service + impact count (subtree devices via `path`).
2. `PUT/GET/DELETE` assignments APIs.
3. `ProfileTreeAssignmentPanel` + confirm ≥50.

### Phase 4 — Effective profile + sync (P1)

1. `resolver.go` + sync integration (tree &gt; route).
2. Update `target_profile_version_id` on assign/publish.
3. Push events on assignment change.

### Phase 5 — Rollout status (P1) — US3

1. `rollout_status.go` recompute from version + device info apps.
2. `GET .../rollout/devices`, `POST .../recompute`.
3. `ProfileRolloutStatusPanel` with polling.

### Phase 6 — Enable / disable (P2) — US4

1. `POST disable` / `enable` APIs.
2. Block assignment when disabled; route editor warning.
3. List badge «Disabled».

### Phase 7 — Best practices UX (P3) — US5

1. Hints in editor + assignment flow (FR-012).
2. Parity doc + quickstart validation.

**Out of scope v1**: canary %, agent changes, automatic assignment bump on publish (manual re-assign).

## Complexity Tracking

> No constitution violations requiring justification.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | — |

## Artifacts Generated (Phase 0–1)

| Artifact | Path |
|----------|------|
| Research | [research.md](./research.md) |
| Data model | [data-model.md](./data-model.md) |
| Contracts | [contracts/](./contracts/) |
| Quickstart | [quickstart.md](./quickstart.md) |

## Next Step

Run **`/speckit-tasks`** to generate `tasks.md`, then implement Phase 1 after 017 is merged/stable on the branch.
