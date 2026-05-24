# Implementation Plan: Enrollment Routes — Controlled Onboarding Gateway

**Branch**: `021-enrollment-routes-ux` | **Date**: 2026-05-24 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `specs/021-enrollment-routes-ux/spec.md`  
**Builds on**: [017-device-control-plane](../017-device-control-plane/plan.md), [019-profile-hub-ux](../019-profile-hub-ux/plan.md) (policy decoupling — enrollment UI completes gateway model)

## Summary

Turn **Enrollment routes** into a **policy-free onboarding gateway**: admin **list + dialog** (dual-column config + live QR), **client Pending QR** / **server Active QR**, **bootstrap app intent** (stable / specific / latest), **tree picker** with placement context, **multi-metric delete impact**, and backend **Definition/Runtime** domain split on existing `enrollment_routes` table.

**Approach**: Extend `enrollment_routes` module (migration `000030`, new option/impact handlers, strip profile from DTOs); instrument public QR for scan metrics; refactor `frontend/src/features/enrollment-routes/` to dialog hub; remove profile gate on create.

## Technical Context

**Language/Version**: Go 1.22+ (`serverBackendGo`), TypeScript/React 18 + Vite (`frontend`)

**Primary Dependencies**: shadcn `Dialog`/`Sheet`, existing `enrollmentQrQuery` + `EnrollmentQrExperience` patterns, `deviceTreeService`, apps catalog queries

**Storage**: PostgreSQL — `enrollment_routes` extended; `applicationversions.is_recommended`; `domain_events` for QR views

**Testing**: `go test` `enrollment_routes/application` (intent resolution, validation); manual quickstart sprints 1–6

**Target Platform**: Admin web; responsive sheet on mobile

**Performance Goals**: Options endpoints &lt; 300ms p95; Pending QR client render &lt; 100ms; list &lt; 500ms for 100 routes

**Constraints**: Headwind envelope; **no profile fields** in UI DTOs; **no nested dialogs**; `qrcodekey` immutable v1; public `/rest/public/qr` parity preserved

**Scale/Scope**: ~12 Go files; ~10–14 React files (new dialog, picker, remove editor page routes); 1 migration; parity `enrollment-routes-ux.md`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Gate | Pass? | Notes |
|------|-------|-------|
| **I. Module-First** | ✅ | `enrollment_routes` + thin touch `qrcode`/public handler for events |
| **II. Layered Clean** | ✅ | Intent resolution in `application/`; SQL in `adapter/persistence` |
| **III. API Parity** | ✅ | Public QR paths unchanged; admin DTO breaking change documented |
| **IV. Testable Delivery** | ✅ | Intent + validation tests; quickstart |
| **V. Simplicity** | ✅ | Single table aggregate; computed tree placement kind |
| **VI. Security** | ✅ | Tenant scope; existing config permissions |
| **VII. Observability** | ✅ | `enrollment_route.qr_viewed` events; env threshold for tree load |

**Post-design**: All gates ✅ — see [research.md](./research.md).

## Project Structure

### Documentation (this feature)

```text
specs/021-enrollment-routes-ux/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── enrollment-contract-payload.md
│   ├── enrollment-route-dialog-ux.md
│   └── enrollment-routes-admin-api.md
└── tasks.md             # (/speckit-tasks)
```

### Source Code

```text
serverBackendGo/
├── db/migrations/
│   ├── 000030_enrollment_routes_ux.up.sql
│   └── 000030_enrollment_routes_ux.down.sql
├── internal/modules/enrollment_routes/
│   ├── domain/
│   │   ├── route.go              # Definition + Runtime structs; EnrollmentRouteView
│   │   ├── bootstrap_intent.go
│   │   └── impact.go
│   ├── application/
│   │   ├── service.go            # remove profile required validation
│   │   ├── intent_resolve.go     # stable/latest/specific
│   │   ├── impact.go
│   │   └── service_test.go
│   ├── adapter/http/
│   │   └── handler.go            # options + impact routes
│   └── adapter/persistence/postgres/
│       └── route_repo.go
├── internal/modules/qrcode/        # or public adapter: emit qr_viewed event
└── docs/parity/enrollment-routes-ux.md

frontend/src/features/enrollment-routes/
├── EnrollmentRouteListPage.tsx     # list + open dialog
├── EnrollmentRouteDialog.tsx         # state machine shell
├── EnrollmentRouteDialogHeader.tsx
├── EnrollmentRouteForm.tsx           # left column
├── EnrollmentRouteQrColumn.tsx       # Pending vs Active
├── TargetNodePicker.tsx              # inline tree panel
├── BootstrapAppPicker.tsx
├── DeleteRouteConfirm.tsx            # inline steps
├── buildEnrollmentContractPreview.ts # client Pending payload
├── enrollmentRouteService.ts         # v2 types, no profile*
└── (remove or redirect EnrollmentRouteEditorPage.tsx)

frontend/src/app/App.tsx              # drop /enrollment-routes/:id editor routes
```

**Structure Decision**: Web app — extend existing modules; no new top-level module.

## Phase 0: Research

Completed → [research.md](./research.md). All Technical Context items resolved (no NEEDS CLARIFICATION).

## Phase 1: Design & Contracts

| Artifact | Path |
|----------|------|
| Data model | [data-model.md](./data-model.md) |
| QR contract | [contracts/enrollment-contract-payload.md](./contracts/enrollment-contract-payload.md) |
| Dialog UX | [contracts/enrollment-route-dialog-ux.md](./contracts/enrollment-route-dialog-ux.md) |
| Admin API | [contracts/enrollment-routes-admin-api.md](./contracts/enrollment-routes-admin-api.md) |
| Quickstart | [quickstart.md](./quickstart.md) |

## Phase 2: Implementation outline (for `/speckit-tasks`)

1. **Migration 000030** — bootstrap columns, `is_recommended`, FK SET NULL, backfill
2. **Backend** — domain split, intent resolver, view DTO mapper, impact + options handlers, drop profile validation
3. **QR telemetry** — `enrollment_route.qr_viewed` on public GET
4. **Frontend service** — types without profile; new API methods
5. **Dialog hub** — state machine per contract; remove editor routes
6. **Pickers** — tree + bootstrap intent
7. **i18n + parity doc** — AR/EN; remove profile onboarding gate on list
8. **Tests + quickstart**

## Complexity Tracking

| Item | Justification |
|------|----------------|
| `applicationversions.is_recommended` | Required for stable ≠ latest (clarification); minimal column |
| `domain_events` QR views | Required for delete impact dimension 3; reuses existing table |
| Admin DTO breaking change | Intentional — strips profile leakage (FR-003) |

## Dependencies & Risks

| Risk | Mitigation |
|------|------------|
| No recommended version seeded | Migration seed or admin doc; clear error on save |
| Pending QR differs from Active | Document contract; Active always server after save |
| Legacy sync still reads `profile_version_id` | Column untouched; out of enrollment module scope |

## Agent Context

Active feature for implementers: branch `021-enrollment-routes-ux`, spec/plan/contracts above. Do **not** add profile/policy copy to enrollment UI.
