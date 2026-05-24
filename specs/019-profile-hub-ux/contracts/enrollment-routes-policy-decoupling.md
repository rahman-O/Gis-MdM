# Contract: Enrollment Routes — Policy Decoupling

**Feature**: `019-profile-hub-ux` | **Audience**: Admin UI + API consumers

---

## Principle

**Enrollment routes** handle *how devices enroll* (QR, default folder, device id mode, main app for provisioning).

**Profiles → Tree assignments** handle *what policy devices receive*.

---

## UI changes

### `EnrollmentRouteEditorPage`

**Remove**:

- Published profile version `<Select>`
- Dependency on `listPublishedProfileVersions()`
- Warning tied only to `profileEnabled` on selected version

**Add**:

- Helper callout (info): «Device policy is assigned from **Profiles → Assignments** to tree folders. This route only controls enrollment placement and QR.»
- Link button: «Open Profiles» (optional)

**Keep**:

- Name, description
- Default tree folder
- Device id mode
- Main app version id (QR provisioning)

---

## API changes

### `POST /rest/private/enrollment-routes`

**Body** (profile optional):

```json
{
  "name": "Warehouse QR",
  "description": null,
  "defaultTreeNodeId": 5,
  "defaultDeviceIdMode": "imei",
  "mainAppId": 12
}
```

`profileVersionId` — **optional**; if omitted, store `NULL` or leave unchanged on update.

### `PUT /rest/private/enrollment-routes/:id`

`profileVersionId` optional; updates ignore profile binding when not sent.

### `GET /rest/private/enrollment-routes/:id`

May still return `profileVersionId` / `profileVersionNumber` for legacy rows — UI does not display.

### Deprecated (UI only v1)

| Endpoint | Status |
|----------|--------|
| `GET /enrollment-routes/options/published-profile-versions` | Hidden from UI; parity marks deprecated |

---

## Sync behavior (unchanged from 018)

1. If device tree folder has assignment → use assigned `profile_version_id`
2. Else if route has legacy `profile_version_id` → use route (backward compat)
3. Else no profile artifact

Admin-facing docs MUST state: **new deployments should use tree assignment only**.

---

## Validation changes

| Rule | Before | After |
|------|--------|-------|
| Create route requires published profile | Yes | **No** |
| Create route requires tree node + main app | Yes | Yes |
| Save route with disabled profile warning | UI warning | Remove (N/A) |

---

## Migration / data

- No DDL required v1
- Optional backfill script (out of scope): copy route `profile_version_id` to tree assignment — manual ops only
