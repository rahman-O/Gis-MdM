# Implementation Plan: Device Control Plane

**Branch**: `017-device-control-plane` | **Date**: 2026-05-23 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/017-device-control-plane/spec.md`  
**Blueprint**: [DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md](../../DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md)

## Summary

Transform the MDM admin platform from a monolithic **Configuration** into a **Device Control Plane**: **device tree** (where devices live), **Profile** (device behavior with draft/publish versions and compiled sync artifacts), and **Enrollment route** (QR + tree placement + profile version binding). **v1 is server + web only** — no Android agent changes; existing Headwind agent remains the QA target.

**Approach**: Add modules `device_tree`, `profiles`, `enrollment_routes`; extend `devices`, `sync`, `qrcode`, `push`; migrate `configurations` data; replace admin nav with Profiles + Enrollment routes; enforce clarified UX (delete-tree dialog, IMEI default QR, ≥50 publish confirm, route-follow on profile version update).

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo`), TypeScript/React 18 (`frontend`)

**Primary Dependencies**: Gin, `lib/pq`, golang-migrate SQL; existing `configurations`, `sync`, `devices`, `qrcode`, `push` modules

**Storage**: PostgreSQL — new tables per [data-model.md](./data-model.md); retain `configurations` during transition with backfill

**Testing**: `go test` on `application/` (compile, publish impact, tree path); golden sync artifact tests; [quickstart.md](./quickstart.md) UAT with Headwind agent

**Target Platform**: Go API `:8080`; Vite admin; Headwind launcher (unchanged binary)

**Project Type**: Web application (backend + frontend)

**Performance Goals**: Tree load &lt; 500ms for 500 nodes; publish compile &lt; 5s typical profile; sync GET from artifact &lt; 200ms p95

**Constraints**: REST parity envelope; `/rest/public/sync/*` paths unchanged; layered modules; no agent fork (§0 blueprint)

**Scale/Scope**: ~60–90 Go files across 3 new modules + extensions; ~15–25 frontend files; 4–6 SQL migrations; parity doc `device-control-plane.md`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | `device_tree`, `profiles`, `enrollment_routes` + phased extensions; MIGRATION/NEXT_STEPS rows |
| **II. Layered Clean** | ✅ | Compile/publish in `profiles/application`; tree rules in `device_tree/application` |
| **III. API Parity** | ✅ | Public sync/QR paths preserved; private aliases documented in contracts |
| **IV. Testable Delivery** | ✅ | Golden artifact + tree unit tests + quickstart |
| **V. Simplicity** | ✅ | Materialized path v1; domain_events only when push batching needed |
| **VI. Security** | ✅ | JWT on `/rest/private/*`; customerId scoping on all repos |
| **VII. Observability** | ✅ | `ProfilePublished` events; structured logs on compile failures |

**Post-design**: All gates ✅. Unknowns resolved in [research.md](./research.md).

## Project Structure

### Documentation (this feature)

```text
specs/017-device-control-plane/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── device-tree-api.md
│   ├── profiles-api.md
│   ├── enrollment-routes-api.md
│   ├── sync-artifact.md
│   └── frontend-control-plane-ux.md
└── tasks.md              # (/speckit-tasks — not created here)
```

### Source Code

```text
serverBackendGo/
├── db/migrations/
│   ├── 0000xx_device_tree.up.sql
│   ├── 0000xx_profiles_versions_artifacts.up.sql
│   ├── 0000xx_enrollment_routes.up.sql
│   ├── 0000xx_devices_control_plane_columns.up.sql
│   └── 0000xx_backfill_configurations_to_profiles.up.sql
├── internal/modules/device_tree/
│   ├── module.go
│   ├── domain/
│   ├── port/
│   ├── application/
│   └── adapter/http/, adapter/persistence/postgres/
├── internal/modules/profiles/
│   ├── application/compile.go, publish.go, impact.go
│   └── ...
├── internal/modules/enrollment_routes/
│   └── ...
├── internal/modules/devices/          # extend repo + move-tree
├── internal/modules/sync/           # artifact loader
├── internal/modules/qrcode/         # create=1, imei default
├── internal/modules/push/           # event consumer (phase 4)
├── internal/platform/push/          # optional outbox worker
└── docs/parity/device-control-plane.md

frontend/src/
├── features/device-tree/
├── features/profiles/               # migrate from configurations/
├── features/enrollment-routes/
├── features/devices/                # tree sidebar integration
├── features/onboarding/             # P2
└── navItems.ts

DEVICE-TREE-POLICY-PROFILE-ANALYSIS.md   # blueprint §20 gates
```

**Structure Decision**: New modules for tree/profile/route; reuse 016 sync mapper inside `profiles/application` compile step; configurations module shrinks to alias handlers then retires.

## Implementation Phases (for `/speckit-tasks`)

Aligned with blueprint §18.10 and spec priorities.

### Phase 1 — Device tree (P1) — US1, FR-001, FR-001a

1. Migration `device_tree_nodes` + customer root seed.
2. Module `device_tree`: CRUD, cycle detection, path/depth recompute.
3. Extend `devices`: `tree_node_id`, list filter by node.
4. Frontend: tree sidebar + move device + delete-with-relocation dialog.
5. Parity section in `device-control-plane.md`.

### Phase 2 — Enrollment auto + identity (P1) — US2, FR-002, FR-002a, FR-013

1. `devices.agent_id` (UUID), `enrollment_state`, `enrollment_route_id` columns.
2. `qrcode` + enroll handler: always `create=1`; default `imei` from route.
3. On enroll success: create device row, assign tree node from route.
4. Frontend: enrollment QR panel locked to create-on-demand.

### Phase 3 — Profiles + publish (P1) — US3, US5, FR-003–007, FR-006

1. Tables `profiles`, `profile_versions`, `profile_version_artifacts`.
2. Module `profiles`: draft save, publish, compile (reuse 016 mapper logic).
3. Impact API + ≥50 confirm dialog.
4. Frontend: Profile editor + usage panel + publish flow.

### Phase 4 — Enrollment routes (P1) — US4, FR-008, FR-008a

1. Table `enrollment_routes`; migrate from `configurations`.
2. Module `enrollment_routes`: binding validation, QR metadata.
3. Route-follow: sync reads route’s current `profile_version_id`.
4. Frontend: route editor (no policy tabs) + rename labels.

### Phase 5 — Sync from artifact (P1) — FR-007, FR-009, NFR-004

1. `sync` loads `artifact_json` instead of live junction assembly.
2. Keep `configurationId` in response = route id.
3. Golden tests vs legacy configuration sample (SC-005).

### Phase 6 — Migration + nav cutover (P1) — FR-010, FR-010a

1. Backfill job: configuration → profile v1 + route + devices links.
2. Replace nav; redirect `/configurations` → `/profiles` (router).
3. API alias `/configurations` one release window.

### Phase 7 — Domain events + push batching (P1 tail) — FR-014

1. `domain_events` + worker → batched `pushmessages`.
2. Admin publish returns immediately.

### Phase 8 — Onboarding wizard (P2) — US6, FR-012

1. Checklist component + guided wizard routes.

**Out of scope v1**: agent app, advanced rollout, profile inheritance (spec).

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

Run **`/speckit-tasks`** to generate dependency-ordered `tasks.md`, then implement Phase 1 behind blueprint §20 Go/No-Go gates.
