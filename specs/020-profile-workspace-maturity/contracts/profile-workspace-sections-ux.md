# Contract: Profile Workspace Sections (Maturity)

**Feature**: `020-profile-workspace-maturity` | **Audience**: Frontend, QA

Extends [019 profile-workspace-ux.md](../../019-profile-hub-ux/contracts/profile-workspace-ux.md).

---

## Global rules (unchanged)

- No nested dialogs  
- Publish from cockpit header only (workspace)  
- Close with unsaved guard in Editor  

---

## Assignments section

### Layout (top → bottom)

1. **Published policy bar** (read-only)  
   - Badge: `Published · v{N}`  
   - Compact strip: Kiosk on/off · Main app name · App count  
   - Button: **View full policy** → `section=editor&versionId={publishedId}&readOnly=1`  

2. **Overlap hint** (conditional)  
   - If parent + child both assigned: note on child row «Devices under this folder use this assignment (overrides parent)»  

3. **Assignment list** (018)  
   - Add folder / remove / version per row  

4. **Empty states**  
   - No published: amber callout + CTA «Publish a version» → Versions/Editor  
   - API error: message + Retry  

### MUST NOT

- Embed full configuration tabs in Assignments  
- Default version picker to draft  

---

## Editor section

- Version switcher: All versions from `GET /versions`  
- Default selection: current **draft** if exists; else prompt create/fork  
- Read-only mode when `readOnly=1` (from Assignments link)  
- Sticky footer: Save draft · Last saved time  
- Unsaved: block section change / workspace close  

---

## Versions section

- Table: version #, status, published date  
- Actions per row:  
  - Draft: Open in Editor · **Delete** (confirm)  
  - Published (historical, deletable): **Delete** (confirm, explains historical)  
  - Published (current): no delete  
- Fork from published → new draft (018)  

---

## Overview section

- Cards source: **published** `pinnedSettings` + rollout + assignments  
- Banner when `hasUnpublishedDraft`: «Draft has unpublished changes» + link Editor  
- If no published: cards from draft + label «Not published yet»  

---

## Publish impact sheet (replaces dialog in workspace)

**Trigger**: Header **Publish**

**Content**:

- Devices affected count  
- Table: folders to update (name, current vN → new vN+1, device count)  
- Routes affected (if any)  
- Primary: **Publish and update assignments**  
- Secondary: Cancel  

**On success**: close sheet · toast · bump `refreshGeneration` · navigate optional to Overview  

---

## Workspace refresh

After save / publish / assignment change / version delete:

- Refetch `GET /summary`  
- Invalidate Assignments list + Versions list  
- Overview cards update within 5s (SC-005)
