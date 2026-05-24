# Implementation Plan: Profile Hub & Enrollment UX

**Branch**: `019-profile-hub-ux` | **Date**: 2026-05-23 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/019-profile-hub-ux/spec.md`  
**Builds on**: [017-device-control-plane](../017-device-control-plane/plan.md), [018-profile-rollout-ops](../018-profile-rollout-ops/plan.md)

## Summary

Deliver an **enterprise-grade Profile Control Plane UX**: decouple **policy** (tree assignments) from **enrollment routes**, replace full-page profile editor with a **layered Profile Workspace** (cockpit header, overview cards, sidebar sections), add **health/badges** on the list, **summary + activity APIs**, and reuse 018 rollout/assignment backends inside section panels.

**Approach**: Thin backend extensions in `profiles` (health, summary, activity) + `enrollment_routes` (optional `profileVersionId`); large frontend refactor under `features/profiles/workspace/`; deprecate profile picker in enrollment UI. No Android changes.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo`), TypeScript/React 18 + Vite (`frontend`)

**Primary Dependencies**: shadcn `Dialog` / `Sheet`, 018 `profileRolloutService`, `deviceTreeService`, existing profile publish/impact APIs

**Storage**: PostgreSQL — reuse 017/018 tables; optional index `000029_profile_hub_activity_idx` on `domain_events`

**Testing**: `go test` on `profiles/application` (health computation); frontend component tests optional; quickstart sprints 1–6

**Target Platform**: Admin web; responsive mobile sheet layout

**Project Type**: Web application (backend + frontend)

**Performance Goals**: `GET /profiles/:id/summary` &lt; 300ms p95; list with health &lt; 500ms p95 for 200 profiles; workspace open perceived &lt; 1s (single summary call)

**Constraints**: Headwind envelope; **no nested dialogs**; tree-only policy in admin UI; stale = 30 days without publish

**Scale/Scope**: ~15 Go files (mostly `profiles/`); ~15–20 new/refactored React files; 0–1 migration; parity `profile-hub-ux.md`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Extends `profiles`, `enrollment_routes`; no new module |
| **II. Layered Clean** | ✅ | Health/summary/activity in `application/`; SQL in adapters |
| **III. API Parity** | ✅ | New private routes documented; enrollment change backward compatible |
| **IV. Testable Delivery** | ✅ | Health unit tests + quickstart |
| **V. Simplicity** | ✅ | Computed health, no materialized table |
| **VI. Security** | ✅ | Tenant + config permissions |
| **VII. Observability** | ✅ | Optional `PROFILE_STALE_PUBLISH_DAYS`; domain events for enable/disable |

**Post-design**: All gates ✅ — see [research.md](./research.md).

## Project Structure

### Documentation (this feature)

```text
specs/019-profile-hub-ux/
├── plan.md              # This file
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/
│   ├── profile-workspace-ux.md
│   ├── profile-hub-api.md
│   └── enrollment-routes-policy-decoupling.md
└── tasks.md             # (/speckit-tasks)
```

### Source Code

```text
serverBackendGo/
├── db/migrations/
│   └── 000029_profile_hub_activity_idx.up.sql   # optional index
├── internal/modules/profiles/
│   ├── application/
│   │   ├── health.go              # ProfileHealthService
│   │   ├── summary.go             # workspace summary aggregator
│   │   └── activity.go            # domain_events reader
│   └── adapter/http/
│       └── hub_handlers.go        # GET summary, GET activity; extend list
├── internal/modules/enrollment_routes/
│   └── application/service.go     # optional profileVersionId; validateBinding
└── docs/parity/profile-hub-ux.md

frontend/src/features/profiles/
├── ProfilesPage.tsx               # Control Radar badges + open workspace
├── workspace/
│   ├── ProfileWorkspace.tsx       # shell: Dialog/Sheet + context
│   ├── ProfileCockpitHeader.tsx
│   ├── ProfileWorkspaceSidebar.tsx
│   ├── ProfileOverviewCards.tsx
│   ├── ProfilePublishImpactSheet.tsx
│   ├── TreePreview.tsx
│   ├── sections/
│   │   ├── AssignmentsSection.tsx   # wraps 018 panel
│   │   ├── RolloutSection.tsx
│   │   ├── VersionsSection.tsx
│   │   ├── EditorSection.tsx
│   │   └── ActivitySection.tsx
│   └── profileWorkspaceState.ts
├── profileHubService.ts           # summary, activity
└── (deprecate direct ProfileEditorPage route default)

frontend/src/features/enrollment-routes/
└── EnrollmentRouteEditorPage.tsx  # remove profile picker

frontend/src/app/App.tsx           # legacy redirect ?open=&section=
```

**Structure Decision**: New `workspace/` subtree keeps 018 components reusable; enrollment change isolated to one page + Go validation.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Backend foundation

1. `health.go` + tests — rules from [research.md](./research.md) R4  
2. Extend `List` repo query with badge flags  
3. `GET /profiles/:id/summary`  
4. `GET /profiles/:id/activity` + enable/disable domain events  
5. Optional migration `000029`  
6. `enrollment_routes`: optional `profileVersionId`, update tests  

### Phase B — Workspace shell

1. `ProfileWorkspace` + context + URL `?open=&section=`  
2. `ProfileCockpitHeader` (Edit, Publish, Close)  
3. Sidebar navigation  
4. Legacy route redirect  

### Phase C — Overview + list radar

1. `ProfileOverviewCards` + summary fetch  
2. `ProfilesPage` health/badges  
3. `profileHubService.ts`  

### Phase D — Sections (018 integration)

1. Mount 018 panels in sections  
2. `TreePreview` in Assignments  
3. `ProfilePublishImpactSheet` (header publish)  

### Phase E — Enrollment + create wizard

1. Strip profile from enrollment editor  
2. Create profile → workspace wizard on Assignments  

### Phase F — Mobile + polish

1. Responsive Sheet layout  
2. i18n keys  
3. Parity doc + quickstart validation  

## Complexity Tracking

_No constitution violations._

## Appendix — Recommended execution order

```text
Phase A (backend) → Phase B (shell) → Phase C (overview/list)
    → Phase D (sections) ∥ Phase E (enrollment)
    → Phase F (mobile/i18n)
```

Depends on **018** APIs being available in target environment.

## Appendix — UX architecture reference

See [contracts/profile-workspace-ux.md](./contracts/profile-workspace-ux.md) for layered workspace, anti-patterns, and 5-second QA rule.
