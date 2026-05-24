# Contract: Profile Rollout API

**Feature**: `018-profile-rollout-ops` | **Audience**: Admin React client

**Base**: `/rest/private` | **Envelope**: Headwind `{ status, message?, data? }`

**Auth**: JWT/session; scoped by `customerId` from principal. Permissions: same as profiles (`edit_config` / `add_config`).

---

## Profile versions (editor navigation)

### `GET /profiles/:profileId/versions`

List all versions for version switcher (FR-004).

**Response `data`**:

```json
[
  {
    "versionId": 45,
    "versionNumber": 3,
    "status": "draft",
    "publishedAt": null,
    "createdAt": "2026-05-20T10:00:00Z"
  },
  {
    "versionId": 42,
    "versionNumber": 2,
    "status": "published",
    "publishedAt": "2026-05-18T09:00:00Z",
    "createdAt": "2026-05-15T08:00:00Z"
  }
]
```

### `POST /profiles/:profileId/versions/:versionId/fork-draft`

Create new draft from a **published** version (FR-005). Returns new draft meta (same shape as `GET /profiles/:id`).

**Errors**: `error.profile.version.notPublished`, `error.profile.draft.exists` (if single-draft policy — return existing draft id).

---

## Tree assignments

### `GET /profiles/:profileId/assignments`

List tree assignments for this profile.

**Response item**:

```json
{
  "assignmentId": 7,
  "treeNodeId": 12,
  "treeNodeName": "Riyadh branch",
  "treePath": "/1/12/",
  "profileVersionId": 42,
  "versionNumber": 2,
  "deviceCount": 34,
  "createdAt": "2026-05-22T12:00:00Z"
}
```

### `GET /profiles/:profileId/assignments/impact`

Preview before assign/reassign (FR-011).

**Query**: `treeNodeId`, `profileVersionId`

**Response**:

```json
{
  "deviceCount": 120,
  "requiresConfirmDialog": true,
  "folderName": "Riyadh branch"
}
```

### `PUT /profiles/:profileId/assignments`

Create or replace assignment for one tree node.

**Body**:

```json
{
  "treeNodeId": 12,
  "profileVersionId": 42,
  "confirmImpact": true
}
```

**Behavior**:
- Rejects if version not `published` or profile `enabled = false`.
- Replaces existing assignment on same `treeNodeId` (unique per folder).
- Sets affected devices `profile_rollout_status = pending`, updates `target_profile_version_id`.
- Enqueues push notify (reuse 017 domain event pattern) when `confirmImpact` or count &lt; 50.

**Response `data`**: assignment object + `{ "affectedDevices": 120 }`

### `DELETE /profiles/:profileId/assignments/:assignmentId`

Remove assignment from folder; devices in subtree fall back to route-only or no target; recompute statuses.

---

## Rollout status dashboard

### `GET /profiles/:profileId/rollout/devices`

Paginated device rollout grid (FR-006–008).

**Query**: `treeNodeId?`, `status?` (`pending|partial|installed|failed`), `page`, `pageSize`

**Response item**:

```json
{
  "deviceId": 1001,
  "deviceName": "Tablet-01",
  "treeNodeId": 15,
  "treeNodeName": "Riyadh / Floor 2",
  "targetVersionId": 42,
  "targetVersionNumber": 2,
  "appliedVersionId": 40,
  "appliedVersionNumber": 1,
  "status": "partial",
  "reason": "App com.example.kiosk: version mismatch",
  "lastUpdate": 1716460000000,
  "resolutionSource": "tree"
}
```

### `POST /profiles/:profileId/rollout/recompute`

Optional manual refresh (admin «Refresh status»); re-runs status engine for profile’s affected devices.

---

## Enable / disable profile

### `POST /profiles/:profileId/disable`

Sets `profiles.enabled = false`. Blocks new assignments. Warn-only on existing enrollment routes.

**Response**: `{ "enabled": false, "devicesMarkedPending": 120 }`

### `POST /profiles/:profileId/enable`

Sets `enabled = true`; marks targeted devices `pending` for re-push.

---

## Sync integration (internal, not HTTP)

`profiles/port.EffectiveProfileResolver` implemented by `profiles/application/resolver.go` and called from `sync/application` when building artifact reference.

**Sync response** (existing public contract, extended fields optional):

```json
{
  "profileId": 10,
  "profileVersionId": 42,
  "profileRevision": "sha256:..."
}
```

Device columns updated during sync/info ingest — see [sync-rollout.md](./sync-rollout.md).

---

## Error keys (stable)

| Key | When |
|-----|------|
| `error.profile.disabled` | Assignment or publish while disabled |
| `error.profile.version.notPublished` | Assign draft version |
| `error.profile.assignment.nodeNotFound` | Invalid tree node |
| `error.profile.assignment.confirmRequired` | ≥50 devices without `confirmImpact` |
| `error.permission.denied` | Missing config permission |
