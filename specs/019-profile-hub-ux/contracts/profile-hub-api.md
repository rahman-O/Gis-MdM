# Contract: Profile Hub API

**Feature**: `019-profile-hub-ux` | **Audience**: Admin React client

**Base**: `/rest/private` | **Envelope**: Headwind `{ status, message?, data? }`

**Auth**: JWT/session; `customerId` from principal. Permissions: same as profiles (`edit_config`, `add_config`).

**Builds on**: [018 profile-rollout-api.md](../../018-profile-rollout-ops/contracts/profile-rollout-api.md) — assignment, rollout, versions, enable/disable unchanged.

---

## Extended list

### `GET /profiles`

**Change**: Each item includes health and badges.

**Response `data[]` item** (additions):

```json
{
  "id": 1,
  "name": "Kiosk Retail",
  "enabled": true,
  "health": "warning",
  "badges": ["no_assignment", "draft_changes"],
  "assignmentCount": 0,
  "rolloutFailureCount": 0,
  "publishedVersion": 2,
  "draftVersionId": 48,
  "deviceCount": 0,
  "enrollmentRouteCount": 1
}
```

---

## Workspace summary (new)

### `GET /profiles/:id/summary`

Single payload for cockpit + overview cards.

**Response `data`**:

```json
{
  "id": 1,
  "name": "Kiosk Retail",
  "description": "",
  "enabled": true,
  "health": "warning",
  "healthReasons": ["no_assignment"],
  "lifecycle": "published",
  "publishedVersionId": 42,
  "publishedVersionNumber": 2,
  "draftVersionId": 48,
  "hasUnpublishedDraft": true,
  "canPublish": true,
  "assignmentCount": 0,
  "assignedFolders": [],
  "rollout": {
    "pending": 0,
    "installed": 0,
    "partial": 0,
    "failed": 0,
    "total": 0
  },
  "pinnedSettings": {
    "kioskMode": true,
    "mainAppName": "MDM Agent",
    "appCount": 7,
    "lastPublishedAt": "2026-05-22T14:00:00Z"
  }
}
```

**Errors**: `404` if profile not found; `error.permission.denied`.

---

## Activity timeline (new)

### `GET /profiles/:id/activity`

**Query**: `limit` (default 50, max 100)

**Response `data`**:

```json
{
  "items": [
    {
      "id": 901,
      "eventType": "ProfilePublished",
      "summaryKey": "profile.activity.published",
      "summaryParams": { "versionNumber": 2 },
      "occurredAt": "2026-05-22T14:00:00Z",
      "actorUserId": 3
    },
    {
      "id": 880,
      "eventType": "ProfileAssignmentChanged",
      "summaryKey": "profile.activity.assigned",
      "summaryParams": { "folderName": "Baghdad/Sales", "versionNumber": 2 },
      "occurredAt": "2026-05-21T09:00:00Z"
    }
  ]
}
```

**Note**: `summaryKey` resolved to localized string on client via i18n.

---

## Meta extension

### `GET /profiles/:id`

Existing meta response **adds**:

```json
{
  "enabled": true,
  "health": "healthy",
  "badges": []
}
```

_(Optional: clients prefer `/summary` for workspace open.)_

---

## Publish impact (reuse 018)

Workspace header **Publish** uses existing:

- `GET /profiles/:id/impact` (or publish impact from 017/018)
- `POST /profiles/:id/versions/:versionId/publish` with `confirmImpact`

Shown in **side sheet**, not modal.

---

## Error keys (new)

| Key | When |
|-----|------|
| `error.profile.summary.notfound` | Invalid profile id |
| `error.enrollment_route.profile_not_required` | Client sent profileVersionId on create (ignored; optional log) |
