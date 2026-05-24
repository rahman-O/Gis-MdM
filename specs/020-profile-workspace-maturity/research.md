# Research: Profile Workspace Maturity

**Feature**: `020-profile-workspace-maturity` | **Date**: 2026-05-23

## R1 — Standalone Profile Editor removal

**Decision**: Keep `ProfileEditorPage` logic as **`ProfileEditorSection`** (or `ProfileEditorCore`) mounted only inside workspace `section=editor`; remove standalone routes from `App.tsx`.

**Rationale**: FR-001 requires single path; embedded mode already exists from 019 partial work; avoids duplicating configuration tabs.

**Alternatives considered**:
- Delete `ProfileEditorPage` entirely and inline — rejected (large refactor risk).
- Keep routes as hidden fallback — rejected (violates SC-001).

---

## R2 — Assignments published context

**Decision**: Extend `GET /profiles/:id/summary` with `publishedContext`:

```json
{
  "publishedVersionId": 42,
  "publishedVersionNumber": 2,
  "pinnedSettings": { "kioskMode": true, "mainAppName": "...", "appCount": 7 }
}
```

Frontend `ProfileAssignmentsContext` renders strip + deep-link `setSection('editor')` with `versionId` query param.

**Rationale**: One round-trip for Assignments header; aligns with Overview pinned data (same source: published version row).

**Alternatives considered**:
- Separate `GET /assignments/context` — rejected (extra latency).
- Full editor embedded in Assignments — rejected per clarification C.

---

## R3 — Publish bumps all assignments

**Decision**: Extend `PublishService.Impact` to return `assignmentsToUpdate: [{ assignmentId, treeNodeId, treeNodeName, deviceCount, currentVersionNumber }]` and `totalDevicesOnAssignments`. On `Publish` with `confirmImpact: true`, after `PublishVersion`:

1. `UPDATE profile_tree_assignments SET profile_version_id = $new WHERE profile_id = $id AND customerid = $cid`
2. Recompute device targets for affected subtrees (reuse `AssignmentService` / rollout recompute from 018)
3. Emit `ProfilePublished` + optional `ProfileAssignmentsBumped` payload

**Rationale**: Clarification B; matches MDM expectation that publish propagates to assigned folders.

**Alternatives considered**:
- Manual per-folder update — rejected by spec.
- Silent bump without preview — rejected by clarification.

---

## R4 — Safe version delete

**Decision**: New `VersionDeleteService` with rules:

| Version state | Allow delete when |
|---------------|-------------------|
| `draft` | Always (not `profiles.draft_version_id` pointer fixup required) |
| `published` | Never if `profiles.published_version_id = versionId` |
| `published` historical | Allowed if no row in `profile_tree_assignments` for that `profile_version_id` AND no `devices.target_profile_version_id` |

Soft-delete vs hard-delete: **hard delete** version row + cascade join tables for that version only (match Java if exists; else hard delete draft rows only in v1).

**Rationale**: FR-005 + clarification B.

**Alternatives considered**:
- Archive flag column — rejected (YAGNI for v1).

---

## R5 — Overview sync without full reload

**Decision**: Lightweight **workspace event bus** in `profileWorkspaceEvents.ts`:

```ts
export const workspaceRefreshKey = ['profile-workspace', profileId] as const
// invalidate on: saveDraft, publish, putAssignment, deleteAssignment, deleteVersion
```

Sections call `getProfileSummary(profileId)` on invalidate; Overview never shows draft body in cards.

**Rationale**: FR-006 + clarification A; React Query or simple `useState` generation counter — prefer counter + `useEffect` to avoid new dependency.

**Alternatives considered**:
- Polling summary every 5s — rejected.
- WebSocket — out of scope.

---

## R6 — Publish impact UI pattern

**Decision**: Replace `ProfilePublishDialog` usage in workspace with **`ProfilePublishImpactSheet`** (shadcn `Sheet` side="right", ~400px). Editor may keep inline publish disabled when workspace header owns publish.

**Rationale**: FR-008 / 019 anti-pattern; clarification B needs room for assignment table.

**Alternatives considered**:
- Extend existing dialog — rejected (nested dialog + cramped).

---

## R7 — Version route redirect

**Decision**: `/profiles/:id/versions/:versionId/edit` → `/profiles?open=:id&section=editor&versionId=:versionId`

**Rationale**: Bookmarks and old links; US1.

---

## R8 — Overlap UI hint

**Decision**: When listing assignments, detect path prefix overlap (`treePath` from 018 list); show badge «Child override» on child row if ancestor also assigned same profile.

**Rationale**: Clarification A; informational only — resolver unchanged.

**Alternatives considered**:
- Block on PUT — rejected (clarification A allows both).
