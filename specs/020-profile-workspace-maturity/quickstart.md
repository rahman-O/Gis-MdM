# Quickstart: Profile Workspace Maturity

**Feature**: `020-profile-workspace-maturity` | **Branch**: `020-profile-workspace-maturity`

## Prerequisites

- [019](../019-profile-hub-ux/quickstart.md) Sprints 1–2 complete (workspace shell, list health)  
- [018](../018-profile-rollout-ops/quickstart.md) complete (`000028`, assignments API)  
- Migrations through `000029` (optional index) applied  

```bash
cd serverBackendGo && make migrate
make dev
```

---

## Sprint 1 — Workspace-only routes

1. Open `/profiles/2/edit` → redirects to `/profiles?open=2&section=editor`.  
2. No full-page editor outside workspace overlay.  
3. Create profile → workspace opens (not editor page).

**Pass**: SC-001 path check.

---

## Sprint 2 — Assignments published context

1. Profile with published v2 → Workspace → **Assignments**.  
2. See `Published · v2` + kiosk/apps summary.  
3. Click **View full policy** → Editor read-only on v2.  
4. Profile with no publish → CTA, not blank screen.

**Pass**: SC-002, SC-003 for Assignments.

---

## Sprint 3 — Multi-folder assignment

1. Assign v2 to folder A and folder B.  
2. List shows both; device counts sane.  
3. Remove A only → B remains.

**Pass**: SC-006.

---

## Sprint 4 — Parent + child assignments

1. Assign profile to parent folder.  
2. Assign same profile to child folder (different version optional).  
3. UI shows both; child row has override hint.  
4. Device under child gets child version (018 resolver).

**Pass**: Clarification A.

---

## Sprint 5 — Publish bumps assignments

1. Assignments on v1; create draft from v1; edit; publish.  
2. Impact sheet lists folders to update.  
3. Confirm → assignments show v2; devices pending.  
4. Overview shows v2 within 5s.

**Pass**: FR-011, SC-005.

---

## Sprint 6 — Version delete

1. Create extra draft → delete from Versions → gone.  
2. Try delete current published while assigned → blocked with message.  
3. Superseded published (no assignment, no devices) → delete allowed.

**Pass**: SC-004 pattern for delete flow.

---

## Sprint 7 — Editor version switch

1. Open Editor → switch to published read-only.  
2. Fork draft from published → edit → unsaved warning on switch.  
3. Save → «Saved» indicator.

**Pass**: US3 acceptance.

---

## API smoke

```bash
# Summary publishedContext
curl -s -b cookies.txt "$BASE/rest/private/profiles/2/summary" | jq '.data.publishedContext'

# Impact with assignments
curl -s -b cookies.txt "$BASE/rest/private/profiles/2/impact" | jq '.data.assignmentsToUpdate'

# Delete draft version
curl -s -X DELETE -b cookies.txt "$BASE/rest/private/profiles/2/versions/99"
```

---

## Regression

- Enrollment routes still work without profile binding (019).  
- Workspace layout: no empty content area (flex + error states).  
- No nested publish dialog inside workspace.
