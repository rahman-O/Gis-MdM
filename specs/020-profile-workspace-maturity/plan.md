# Implementation Plan: Profile Workspace Maturity

**Branch**: `020-profile-workspace-maturity` | **Date**: 2026-05-23 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/020-profile-workspace-maturity/spec.md`  
**Builds on**: [019-profile-hub-ux](../019-profile-hub-ux/plan.md) (workspace shell), [018-profile-rollout-ops](../018-profile-rollout-ops/plan.md) (assignments, versions, publish)

## Summary

Complete **Profile Workspace** as the sole admin path for profiles: remove standalone **Profile Editor** routes, fix **Assignments** with published-version context (compact summary + link to read-only Editor), implement **publish → bump all assignments** with impact preview, **safe version delete**, **Editor version switching** with unsaved guards, **Overview** synced to published (+ draft badge), and **ProfilePublishImpactSheet** (no nested dialogs).

**Approach**: Extend `profiles` application layer (`PublishService`, new `VersionDeleteService`, enriched `Impact`/`Summary`); refactor frontend `workspace/` sections; deprecate `ProfileEditorPage` routes; reuse 018 assignment/rollout APIs.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo`), TypeScript/React 18 + Vite (`frontend`)

**Primary Dependencies**: 019 workspace (`ProfileWorkspace`, `profileHubService`), 018 `profileRolloutService`, shadcn `Sheet` for publish impact, existing configuration tabs in embedded editor

**Storage**: PostgreSQL — reuse `profiles`, `profile_versions`, `profile_tree_assignments`, `devices`; **no new tables**; optional domain events for `ProfileVersionDeleted`

**Testing**: `go test` — `version_delete_test.go`, `publish_assignments_test.go`, extend `health_test.go`; frontend manual quickstart sprints 1–7

**Target Platform**: Admin web (desktop Dialog + mobile Sheet from 019)

**Project Type**: Web application (backend + frontend)

**Performance Goals**: Summary/assignments context &lt; 300ms p95; publish impact (incl. assignment list) &lt; 500ms p95; workspace section switch &lt; 200ms perceived (cached summary + invalidate on mutation)

**Constraints**: Headwind envelope; **no nested dialogs** (019/020); Overview always shows **published** pinned settings when published exists; child-over-parent resolver unchanged (018)

**Scale/Scope**: ~12 Go files touched; ~18 React files new/refactored; 0 migrations (events only); parity `profile-workspace-maturity.md`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | Extends `profiles` only; enrollment routes unchanged (019) |
| **II. Layered Clean** | ✅ | Delete/bump rules in `application/`; HTTP thin |
| **III. API Parity** | ✅ | New/extended private routes documented in contracts |
| **IV. Testable Delivery** | ✅ | Unit tests + [quickstart.md](./quickstart.md) |
| **V. Simplicity** | ✅ | No new tables; extend publish impact DTO |
| **VI. Security** | ✅ | Tenant + `edit_config`; delete guards in domain |
| **VII. Observability** | ✅ | `ProfileVersionDeleted` domain event |

**Post-design**: All gates ✅ — see [research.md](./research.md).

## Project Structure

### Documentation (this feature)

```text
specs/020-profile-workspace-maturity/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── profile-workspace-maturity-api.md
│   ├── profile-workspace-sections-ux.md
│   └── publish-assignment-bump.md
└── tasks.md             # (/speckit-tasks)
```

### Source Code

```text
serverBackendGo/internal/modules/profiles/
├── application/
│   ├── publish.go              # extend Impact + Publish: bump assignments
│   ├── version_delete.go       # delete rules FR-005
│   ├── hub.go                  # summary: publishedContext for assignments
│   └── *_test.go
├── domain/
│   ├── publish.go              # PublishImpact + AssignmentBumpPreview
│   └── version_delete.go
└── adapter/http/
    ├── handler.go              # DELETE version; extend impact response
    └── rollout_handlers.go     # (unchanged paths)

serverBackendGo/docs/parity/profile-workspace-maturity.md

frontend/src/features/profiles/
├── ProfileEditorPage.tsx       # embedded-only; extract ProfileEditorCore.tsx
├── ProfileEditRedirect.tsx     # extend version routes
├── workspace/
│   ├── ProfileWorkspace.tsx
│   ├── ProfilePublishImpactSheet.tsx   # replaces ProfilePublishDialog in workspace
│   ├── ProfileAssignmentsContext.tsx # published summary strip
│   ├── ProfileEditorSection.tsx      # version switcher + sticky save
│   ├── ProfileOverviewSection.tsx    # invalidate on workspace events
│   └── profileWorkspaceEvents.ts     # save/publish/assign → refresh summary
├── profileHubService.ts
└── ProfilesPage.tsx

frontend/src/app/App.tsx          # remove standalone editor routes; redirects only
```

**Structure Decision**: 020 is a **maturity layer** on 019 shell — avoid duplicating 018 panels; add thin context components and backend publish/assignment coordination.

## Implementation Phases (for `/speckit-tasks`)

### Phase A — Backend: publish bump + version delete

1. `PublishImpact` includes `assignmentsToUpdate[]`, `devicesAffected` (assignments + rollout)  
2. `Publish()` after success: `BumpAllAssignmentsToVersion(profileID, newVersionID)` + device pending  
3. `DELETE /profiles/:id/versions/:versionId` with delete guards  
4. Extend `GET /summary` with `publishedContext` block for Assignments strip  
5. Domain event `ProfileVersionDeleted`  
6. Unit tests + parity doc  

### Phase B — Workspace-only navigation (US1)

1. Remove `/profiles/:id/edit` and `/profiles/:id/versions/:vid/edit` as full pages  
2. Redirects to `?open=&section=editor&versionId=`  
3. Create profile → workspace (Overview or Assignments)  

### Phase C — Assignments context (US2)

1. `ProfileAssignmentsContext` — published version badge + pinned strip  
2. Link «View full policy» → Editor `readOnly` + `versionId=published`  
3. Empty/error states (no published, API fail)  

### Phase D — Editor + Versions (US3, US5)

1. `ProfileEditorSection` — version selector from 018 list; dirty guard  
2. `VersionsSection` — delete draft / superseded published with inline confirm  
3. Fork draft from published (018)  

### Phase E — Publish UX (US7, FR-011)

1. `ProfilePublishImpactSheet` in header flow  
2. Show assignment folders in impact sheet  
3. On success: emit workspace refresh → Overview + Assignments  

### Phase F — Overview sync (US6)

1. Overview reads summary; `hasUnpublishedDraft` banner only  
2. Subscribe to workspace refresh bus after save/publish/assign/delete  

### Phase G — Multi-folder polish (US4)

1. Assignment list shows overlap hint (parent/child)  
2. Verify multi PUT + DELETE (018 already supports)  

## Complexity Tracking

_No constitution violations._

## Appendix — Execution order

```text
019 + 018 baseline merged
  → Phase A (backend)
  → Phase B (routes)
  → Phase C + D (Assignments + Editor)  [P1]
  → Phase E (publish sheet)
  → Phase F + G (overview + polish)
```

## Appendix — Clarifications baked in

| Topic | Decision |
|-------|----------|
| Parent/child assignments | Allow both; nearest wins (018) |
| Publish | Auto-bump all assignments with confirm in impact sheet |
| Overview | Published cards + draft badge only |
| Delete v1 | Drafts + unused old published |
| Assignments UI | Compact summary + link to full read-only Editor |
