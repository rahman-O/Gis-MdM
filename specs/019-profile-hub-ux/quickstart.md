# Quickstart: Profile Hub & Enrollment UX

**Feature**: `019-profile-hub-ux` | **Branch**: `019-profile-hub-ux`

## Prerequisites

- [017](../017-device-control-plane/quickstart.md) complete: tree, profiles, enrollment routes
- [018](../018-profile-rollout-ops/quickstart.md) complete: assignments, rollout, `000028` migration
- `.env`: `MODULE_PROFILES_ENABLED=true`, `MODULE_PROFILE_ROLLOUT_ENABLED=true`, `MODULE_ENROLLMENT_ROUTES_ENABLED=true`

```bash
cd serverBackendGo && make migrate   # through 000028+ (000029 index optional)
make dev   # API + frontend
```

---

## Sprint 1 — Enrollment route decoupling

1. Open **Enrollment routes** → create route **without** profile picker.
2. Set default tree folder + main app → save.
3. Confirm helper text mentions Profiles → Assignments.
4. API: `POST /rest/private/enrollment-routes` body without `profileVersionId` → `201`.

**Pass**: Route saves; QR works; no profile field in UI.

---

## Sprint 2 — Profile Workspace shell

1. Open **Profiles** list — see health chips and badges.
2. Click a profile → **Workspace** opens (large overlay, sidebar visible).
3. Header shows name, health, lifecycle, **Edit**, **Publish**, **Close**.
4. Default section **Overview** — cards only, no input fields.
5. Navigate sidebar: Assignments → Rollout → Versions → Activity → back Overview.
6. Press Esc → workspace closes; list scroll position preserved.

**Pass**: No nested popup dialogs during navigation (SC-007).

---

## Sprint 3 — Read vs Edit + Publish from header

1. From Overview, click **Edit** → Editor section + warning bar + sticky save.
2. Change a field → **Publish** in header enables (if draft valid).
3. Click **Publish** → **right sheet** shows impact (device/folder counts).
4. Confirm publish → return to Overview; published version updates.

**Pass**: Publish never only inside editor tabs.

---

## Sprint 4 — Assignments + tree preview + create wizard

1. Create new profile from list.
2. Workspace opens on **Assignments** with wizard hint.
3. Publish first if prompted → assign published version to folder.
4. **Tree preview** highlights selected folder path.
5. Overview card shows assignment count.

**Pass**: Create-to-assign &lt; 3 minutes.

---

## Sprint 5 — Activity + 5-second rule

1. Publish + assign → open **Activity** — see published + assigned events.
2. With test device failures, Overview rollout card shows failures; list badge `rollout_issues`.
3. **5-second test**: show Overview to colleague 5s → ask four questions (published? assigned? disabled? failures?) — ≥9/10 correct.

---

## Sprint 6 — Mobile smoke

1. Narrow viewport (&lt; 768px).
2. Open profile → full-screen sheet; sidebar via menu drawer.
3. Bottom bar: Edit / Publish / Close.

---

## API smoke

```bash
# Summary
curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE/rest/private/profiles/1/summary" | jq .

# Activity
curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE/rest/private/profiles/1/activity?limit=20" | jq .

# List with health
curl -s -H "Authorization: Bearer $TOKEN" \
  "$BASE/rest/private/profiles" | jq '.data[0] | {name, health, badges}'
```

---

## Regression

- Legacy URL `/profiles/1/edit` → redirects to workspace editor section.
- Device in assigned folder still syncs policy (018 resolver).
- Route without profile but with tree assignment → device gets policy.
