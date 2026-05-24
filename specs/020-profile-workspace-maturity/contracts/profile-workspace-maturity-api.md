# Contract: Profile Workspace Maturity API

**Feature**: `020-profile-workspace-maturity` | **Audience**: Admin React client

**Base**: `/rest/private` | **Envelope**: Headwind `{ status, message?, data? }`

**Builds on**: [019 profile-hub-api.md](../../019-profile-hub-ux/contracts/profile-hub-api.md), [018 profile-rollout-api.md](../../018-profile-rollout-ops/contracts/profile-rollout-api.md)

---

## Extended summary (Assignments context)

### `GET /profiles/:id/summary`

**Addition** to 019 response — explicit block for Assignments strip (may duplicate top-level pinned fields for convenience):

```json
{
  "publishedContext": {
    "versionId": 42,
    "versionNumber": 2,
    "status": "published",
    "pinnedSettings": {
      "kioskMode": true,
      "mainAppName": "MDM Agent",
      "appCount": 7,
      "lastPublishedAt": "2026-05-22T14:00:00Z"
    }
  },
  "hasUnpublishedDraft": true
}
```

If no published version: `publishedContext: null`.

---

## Extended publish impact

### `GET /profiles/:id/impact`

**When** draft exists and will be published (or query `?versionId=` draft id):

**Extended `data`**:

```json
{
  "deviceCount": 120,
  "enrollmentRouteCount": 0,
  "requiresConfirmDialog": true,
  "assignmentsToUpdate": [
    {
      "assignmentId": 7,
      "treeNodeId": 12,
      "treeNodeName": "Riyadh branch",
      "currentVersionNumber": 1,
      "deviceCount": 34
    }
  ]
}
```

### `POST /profiles/:id/versions/:versionId/publish`

**Unchanged path**; **extended behavior** when `confirmImpact: true`:

1. Publish draft → new published version  
2. Update **all** `profile_tree_assignments.profile_version_id` to new version  
3. Recompute rollout pending for affected devices  

**Response** (add fields):

```json
{
  "publishedVersionId": 48,
  "versionNumber": 3,
  "artifactHash": "...",
  "affectedDevices": 120,
  "affectedRoutes": 0,
  "assignmentsUpdated": 2
}
```

---

## Version delete (new)

### `DELETE /profiles/:id/versions/:versionId`

**Auth**: `edit_config`

**Success**: `200` `{ "status": "OK" }`

**Errors**:

| Key | When |
|-----|------|
| `error.profile.version.delete.activePublished` | Version is current published |
| `error.profile.version.delete.assigned` | Referenced in `profile_tree_assignments` |
| `error.profile.version.delete.devicesTarget` | Devices still target this version |
| `error.notfound.profile` | Invalid ids |

**Side effect**: `ProfileVersionDeleted` domain event.

---

## Unchanged (reuse 018)

- `GET/PUT/DELETE /profiles/:id/assignments*`  
- `GET /profiles/:id/versions`  
- `POST /profiles/:id/versions/:id/fork-draft`  
- `GET/PUT /profiles/:id/versions/:versionId` (editor payload)

---

## Frontend routes (deprecate)

| Legacy | Replacement |
|--------|-------------|
| `/profiles/:id/edit` | `/profiles?open=:id&section=editor` |
| `/profiles/:id/versions/:vid/edit` | `/profiles?open=:id&section=editor&versionId=:vid` |
