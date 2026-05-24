# Parity: Profile Workspace Maturity (020)

**Feature**: `020-profile-workspace-maturity` | **Status**: Implemented (verify via quickstart)

## API extensions

### `GET /private/profiles/:id/summary`

Adds `publishedContext` (`versionId`, `versionNumber`, `status`, `pinnedSettings`) or `null` when unpublished.

### `GET /private/profiles/:id/impact`

Adds `assignmentsToUpdate[]` with folder rows. `requiresConfirmDialog` is true when device count ≥ 50 **or** any assignment will bump.

### `POST /private/profiles/:id/versions/:versionId/publish`

Transactional publish + bump all `profile_tree_assignments` to the new version and mark subtree devices pending.

Response adds `assignmentsUpdated`.

### `DELETE /private/profiles/:id/versions/:versionId`

Deletes draft or unused historical versions.

| Error key | When |
|-----------|------|
| `error.profile.version.delete.activePublished` | Current published version |
| `error.profile.version.delete.assigned` | Referenced in `profile_tree_assignments` |
| `error.profile.version.delete.devicesTarget` | Devices still target this version |

Side effect: `ProfileVersionDeleted` domain event.

## Assignments (018 + 020 UI)

- Multiple folders per profile allowed (parent + child).
- **Resolver (unchanged)**: nearest assignment on the device tree path wins; child folder overrides parent.
- UI shows «Child override» on nested rows and an informational banner when overlap exists.

## Frontend routes

| Legacy | Workspace URL |
|--------|----------------|
| `/profiles/:id/edit` | `/profiles?open=:id&section=editor` |
| `/profiles/:id/versions/:vid/edit` | `/profiles?open=:id&section=editor&versionId=:vid` |
| Read-only from Assignments | `&readOnly=1` |

Publish: cockpit header → **Publish impact sheet** (not nested dialog). Editor publish hidden in workspace.

## Manual checks (quickstart)

```bash
export BASE=http://localhost:8080/rest
# Session cookie from browser devtools after login

# 1 — Legacy redirect (browser): /profiles/2/edit → workspace editor

# 2 — Summary published context
curl -s -b "$COOKIE" "$BASE/private/profiles/1/summary" | jq '.data.publishedContext'

# 3 — Impact preview with assignments
curl -s -b "$COOKIE" "$BASE/private/profiles/1/impact" | jq '.data | {deviceCount, assignmentsToUpdate}'

# 4 — Publish (draft version id from meta)
curl -s -b "$COOKIE" -X POST "$BASE/private/profiles/1/versions/DRAFT_ID/publish" \
  -H 'Content-Type: application/json' -d '{"confirmImpact":true}' | jq '.data.assignmentsUpdated'

# 5 — Delete unused draft
curl -s -b "$COOKIE" -X DELETE "$BASE/private/profiles/1/versions/UNUSED_ID"

# 6 — Block delete active published (expect ERROR)
curl -s -b "$COOKIE" -X DELETE "$BASE/private/profiles/1/versions/PUBLISHED_ID"
```

## UI smoke (dev server)

1. Open `/profiles?open=PROFILE_ID&section=assignments` — published bar + folder list.
2. **View full policy** → editor read-only with `readOnly=1`.
3. Header **Publish** → side sheet with folder table → confirm → Overview updates without full page reload.
4. **Versions** — delete draft; block delete on active published.
5. Assign parent + child folders — overlap hint visible.
