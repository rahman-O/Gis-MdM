# Contract: Enrollment Route Dialog UX (State Machine)

**Feature**: `021-enrollment-routes-ux` | **Audience**: Frontend implementers

## Shell

- **Route**: `/enrollment-routes` list only (remove `/enrollment-routes/new`, `/enrollment-routes/:id` full pages from default nav).
- **Component**: `EnrollmentRouteDialog` — shadcn `Dialog` (≥md) / `Sheet` (<md).
- **No nested dialogs** — delete confirm = inline step inside dialog body/footer area.

---

## States

| State ID | Entry | Header badges | Left column | Right column | Footer |
|----------|-------|---------------|-------------|--------------|--------|
| `LIST` | default | — | table | — | New route |
| `DIALOG_CREATE` | New route | Draft | edit fields | Pending QR | Save, Cancel |
| `DIALOG_OVERVIEW` | row click | Active | read-only | Active QR | Edit, Delete, Close |
| `DIALOG_EDIT` | Edit from overview | Active + Unsaved* | edit fields | Pending QR | Save, Cancel |
| `DELETE_STEP1` | Delete | Active | impact summary (3 metrics) | Active QR dimmed | Continue, Back |
| `DELETE_STEP2` | impact any > 0 | Active | typed name + metrics | — | Confirm delete, Back |
| `DELETE_CONFIRM_ZERO` | all metrics 0 | Active | short summary | — | Confirm, Back |

\*Unsaved only when `form ≠ lastSavedSnapshot`.

---

## Transitions

```text
LIST --New--> DIALOG_CREATE
LIST --row--> DIALOG_OVERVIEW

DIALOG_CREATE --Save OK--> DIALOG_OVERVIEW
DIALOG_CREATE --Cancel--> LIST (confirm if dirty)

DIALOG_OVERVIEW --Edit--> DIALOG_EDIT
DIALOG_OVERVIEW --Close--> LIST
DIALOG_OVERVIEW --Delete--> DELETE_STEP1

DIALOG_EDIT --Save OK--> DIALOG_OVERVIEW (refresh Active QR)
DIALOG_EDIT --Cancel--> DIALOG_OVERVIEW (confirm if dirty)

DELETE_STEP1 --all zero--> DELETE_CONFIRM_ZERO
DELETE_STEP1 --any > 0--> DELETE_STEP2
DELETE_STEP1 --Back--> DIALOG_OVERVIEW

DELETE_STEP2 --name match + Confirm--> LIST
DELETE_CONFIRM_ZERO --Confirm--> LIST
```

---

## Header badge rules

| Badge | Condition |
|-------|-----------|
| Draft | `DIALOG_CREATE` only (no persisted route) |
| Active | persisted route (`routeId > 0`) |
| Unsaved changes | `DIALOG_EDIT` && dirty form |

---

## QR column rules

| Dialog mode | QR column |
|-------------|-----------|
| CREATE, EDIT | Pending — client `buildEnrollmentContractPreview(form)` |
| OVERVIEW | Active — `getEnrollmentRouteQr` + public image URL |
| DELETE_* | Active dimmed + disclaimer |

While EDIT with unsaved changes: **production QR remains scannable** (last saved key) — optional small link “Open last active QR” in footer hint; right column shows Pending only (per spec FR-006b).

---

## Target node picker (sub-flow)

| Step | UI |
|------|-----|
| Open picker | Tree panel overlay **inside** dialog (not new Dialog) |
| Hover/select node | breadcrumb + deviceCount + heavilyLoaded + container warning |
| Confirm inheritable | allow + checkbox optional ack stored on Save |
| Invalid node | block confirm |

---

## Bootstrap app picker

1. Select application (package).
2. Select intent: Stable (default) | Specific | Latest.
3. If Specific → version dropdown.
4. Show resolved version line under picker.

---

## Copy / i18n keys (prefix `enrollmentRoute.`)

- `status.draft`, `status.active`, `status.unsaved`
- `qr.pending`, `qr.active`, `qr.saveToActivate`
- `tree.containerWarning`, `tree.heavilyLoaded`
- `delete.impact.enrolling`, `delete.impact.historical`, `delete.impact.scans`
- `delete.typeRouteName`

**Forbidden** copy keys: `profile`, `policy`, `assignment`.

---

## Accessibility

- Focus trap in dialog; Esc → close with unsaved confirm.
- QR image: `alt` includes route name + Pending/Active state.
