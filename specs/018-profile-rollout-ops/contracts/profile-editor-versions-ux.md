# Contract: Profile Editor — Version Navigation UX

**Feature**: `018-profile-rollout-ops` | **Audience**: React admin

## Header controls

| Control | Behavior |
|---------|----------|
| Version dropdown | Lists all versions from `GET /profiles/:id/versions`; shows `#N · Draft` or `#N · Published · date` |
| Status badge | Draft = editable; Published = read-only shell |
| «Fork draft from this version» | Visible on published rows; calls `POST .../fork-draft` then navigates to new draft |
| Best-practice hint | Static line: «Edit draft → Publish → Assign to tree folder → Monitor rollout» (FR-012) |

## Unsaved draft guard

When `dirty === true` and user changes version in dropdown or navigates away:

1. Modal: **Save** | **Discard** | **Cancel**
2. Save → `PUT /profiles/:id/versions/:draftId` then switch
3. Discard → reload selected version without save

## Routes

| Path | Mode |
|------|------|
| `/profiles/:profileId/edit` | Latest draft or auto-fork (017 behavior) |
| `/profiles/:profileId/versions/:versionId/edit` | Explicit version; read-only if published |

## Read-only published view

- All tabs visible (Restrictions, MDM, Apps, Design, Files)
- Inputs disabled; Save hidden; Publish hidden
- Banner: «Viewing published version #N — create draft to edit»

## Assignment entry point (same page)

- Tab or side panel **«Tree assignment»** (FR-001):
  - Tree picker (reuse `DeviceTreeSidebar` node picker modal)
  - Published version selector (default: latest published)
  - Impact preview + confirm ≥50
  - List current assignments with remove action

## Rollout tab

- Sub-route or tab **«Rollout status»**:
  - Filters: folder, status
  - Table from `GET .../rollout/devices`
  - Auto-refresh every 60s + manual Refresh (SC-008)
  - Status chips: pending (amber), installed (green), partial (orange), failed (red)

## i18n keys (suggested)

- `profile.versions.select`, `profile.versions.forkDraft`, `profile.versions.unsavedWarning`
- `profile.rollout.status.pending`, `.installed`, `.partial`, `.failed`
- `profile.assignment.title`, `profile.assignment.confirm.title`
- `profile.disabled.banner`
