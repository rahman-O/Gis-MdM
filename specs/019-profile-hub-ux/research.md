# Research: Profile Hub & Enrollment UX

**Feature**: `019-profile-hub-ux` | **Date**: 2026-05-23

## R1 — Workspace shell: Dialog vs Sheet

**Decision**: Desktop uses **Radix `Dialog`** with custom layout classes (`max-w-[min(1400px,96vw)]`, `h-[min(900px,94vh)]`, light overlay `bg-black/40`). Mobile (`< md`) uses **`Sheet`** full-screen from bottom with same inner layout.

**Rationale**: Spec requires near-fullscreen workspace, not centered 400px popup. Project already has `dialog.tsx` and `sheet.tsx` (shadcn). Single `ProfileWorkspace` component switches primitive by breakpoint.

**Alternatives considered**:
- Dedicated route `/profiles/:id` — rejected; spec requires list context preserved.
- Only Sheet from right — rejected; spec asks for centered workspace with side padding on desktop.

---

## R2 — Secondary UI: no nested dialogs

**Decision**: **Impact Preview**, assignment confirm, and rollout filters use **`Sheet` side="right"`** or **inline collapsible panels** inside workspace content area — same z-index stack, `modal={false}` on nested sheets where needed, or single workspace-level `secondaryPanel` state.

**Rationale**: FR-005 / SC-007 forbid dialog-on-dialog cognitive maze.

**Alternatives considered**:
- `AlertDialog` for publish — rejected for rich impact content; use side panel with primary/secondary actions.

---

## R3 — Policy source: decouple enrollment route from profile

**Decision**:
- **UI**: Remove profile version picker from `EnrollmentRouteEditorPage`; show helper text linking to Profiles → Assignments.
- **API**: `profileVersionId` becomes **optional** on `POST/PUT /enrollment-routes`; validation no longer requires published profile version. Existing rows keep column value for backward compatibility; **sync resolver** (018) already prefers tree assignment over route.
- **Deprecate** `GET /enrollment-routes/options/published-profile-versions` in parity doc (endpoint may return 410 or empty with deprecation header in later phase; v1: hide from UI only).

**Rationale**: Single admin mental model — tree is policy source; routes are enrollment plumbing only.

**Alternatives considered**:
- DB migration dropping `profile_version_id` — out of scope v1.
- Auto-copy route profile to tree assignment — rejected; explicit admin assignment is clearer.

---

## R4 — Profile Health computation

**Decision**: Server-side **`ProfileHealthService`** in `profiles/application` computes health on read (list + summary), no new table.

| Health | Rule (priority order) |
|--------|------------------------|
| `draft_only` | No `published_version_id` |
| `error` | `enabled=false` AND devices with pending target OR rollout `failed` count > 0 |
| `warning` | No tree assignments OR unpublished draft exists OR `stale` (published_at &lt; now − 30 days) |
| `healthy` | Published + ≥1 assignment + enabled + zero failed rollout in assigned subtree |

**Stale threshold**: **30 days** (configurable later via env `PROFILE_STALE_PUBLISH_DAYS`).

**List badges**: Derived flags on same query (`hasAssignment`, `hasDraftChanges`, `rolloutFailureCount`, etc.).

**Rationale**: Avoid stale client-side logic; list and workspace header stay consistent.

**Alternatives considered**:
- Materialized health column — rejected (YAGNI); recompute on read with indexed queries is enough for v1 scale.

---

## R5 — Overview summary API

**Decision**: Add `GET /profiles/:id/summary` returning cockpit + card payloads in one round-trip (meta + assignment counts + rollout aggregates + pinned settings snippet).

**Rationale**: 5-second rule needs one fast load; avoids 4 parallel calls on workspace open.

**Alternatives considered**:
- Extend `GET /profiles/:id` only — acceptable fallback; summary endpoint keeps meta stable and adds aggregates.

---

## R6 — Activity timeline

**Decision**: `GET /profiles/:id/activity?limit=50` reads `domain_events` where `aggregate_id` IN (`strconv(profileId)`, `profile:{id}`) and `event_type` IN (`ProfilePublished`, `ProfileAssignmentChanged`). Add `ProfileEnabled` / `ProfileDisabled` events in enable service (small addition).

Rollout failures: include synthetic activity rows from latest failed devices query (cap 10) merged into timeline by timestamp.

**Rationale**: Reuses 017 outbox table; no new schema.

**Alternatives considered**:
- Separate `profile_activity` table — rejected for v1.

---

## R7 — Frontend state architecture

**Decision**:
- `ProfileWorkspaceContext` — `profileId`, `section`, `mode: 'read' | 'edit'`, `dirty`, `secondaryPanel: null | 'publish-impact' | ...`
- URL query sync: `/profiles?open=12&section=rollout` for deep links and legacy redirect from `/profiles/12/edit`
- Reuse 018 panels as **section components** mounted lazily

**Rationale**: One state owner prevents prop drilling; query params restore workspace after refresh.

---

## R8 — Create flow wizard

**Decision**: After `POST /profiles`, open workspace with `section=assignments` and `wizardStep: 'assign' | 'publish-first'` (client-only state). No separate modal wizard.

**Rationale**: FR-016 — wizard inside workspace, not nested.

---

## R9 — Miniature tree preview

**Decision**: Reuse `getDeviceTree()` + compact recursive `TreePreview` component (read-only, max depth 4, scroll). Highlight selected `treeNodeId` with `ring-2 ring-primary`.

**Rationale**: Data already available; no new API.

---

## R10 — i18n & design tokens

**Decision**: Keys under `profile.workspace.*`, `profile.health.*`, `profile.badges.*` in `en.json` / `ar.json`. Lifecycle colors via Tailwind semantic classes: `amber` draft, `emerald` published, `muted`/`destructive` disabled.

**Rationale**: Matches 018 pattern; supports RTL for sidebar layout (`dir` aware).
