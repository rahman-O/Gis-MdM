# Data Model: Profile Workspace Maturity

**Feature**: `020-profile-workspace-maturity` | **Date**: 2026-05-23

## Overview

**No new tables.** Extends behavior on existing 017/018/019 schema and adds API DTOs + domain events.

**Depends on**: 019 (`ProfileSummary`, hub), 018 (`profile_tree_assignments`, rollout), 017 (`profile_versions`).

---

## Existing tables (unchanged)

| Table | 020 usage |
|-------|-----------|
| `profiles` | `published_version_id`, `draft_version_id` — delete must fix pointers |
| `profile_versions` | DELETE row when guards pass |
| `profile_tree_assignments` | Multi-row per profile; bulk UPDATE on publish |
| `devices` | Guard delete if `target_profile_version_id` references version |
| `domain_events` | `ProfileVersionDeleted`, `ProfileAssignmentsBumped` (optional) |

---

## Publish impact extension

### `PublishImpact` (extended)

| Field | Type | Description |
|-------|------|-------------|
| `deviceCount` | int | Devices that would receive new policy |
| `enrollmentRouteCount` | int | Legacy routes (if any) |
| `requiresConfirmDialog` | bool | Existing |
| `assignmentsToUpdate` | array | Folders with active assignment |
| `assignmentsToUpdate[].assignmentId` | int | |
| `assignmentsToUpdate[].treeNodeId` | int | |
| `assignmentsToUpdate[].treeNodeName` | string | |
| `assignmentsToUpdate[].currentVersionNumber` | int | |
| `assignmentsToUpdate[].deviceCount` | int | Subtree device count |

### Publish side effect (transactional)

On successful publish of version `V_new`:

1. Set `profiles.published_version_id = V_new`, clear/replace draft per 017 rules  
2. `UPDATE profile_tree_assignments SET profile_version_id = V_new WHERE profile_id = P`  
3. Mark affected devices `profile_rollout_status = pending`, update `target_profile_version_id`  
4. Insert `ProfilePublished` event (existing)  

---

## Version delete rules

### `VersionDeleteEligibility`

| Check | Error key |
|-------|-----------|
| Version not found | `error.notfound.profile` |
| Is current `published_version_id` | `error.profile.version.delete.activePublished` |
| Is referenced in `profile_tree_assignments` | `error.profile.version.delete.assigned` |
| Referenced in `devices.target_profile_version_id` | `error.profile.version.delete.devicesTarget` |
| Status = `draft` and is only draft | Allow; update `profiles.draft_version_id` if needed |

Historical published: allowed when all guards pass.

---

## Summary extension for Assignments

Add to `ProfileSummary` (or nested `publishedContext`):

| Field | Type |
|-------|------|
| `publishedContext.versionId` | int |
| `publishedContext.versionNumber` | int |
| `publishedContext.pinnedSettings` | same as Overview `pinnedSettings` |

When no published version: `publishedContext` null → UI shows CTA.

---

## Workspace URL state

| Query | Meaning |
|-------|---------|
| `open` | profile id |
| `section` | overview \| assignments \| … \| editor |
| `versionId` | optional editor target (published read-only or draft) |
| `readOnly` | `1` when opening from Assignments link |

---

## Frontend state (client-only)

| State | Scope |
|-------|-------|
| `editorDirty` | workspace context (019) |
| `refreshGeneration` | bump to refetch summary |
| `secondaryPanel` | `publish-impact` \| null (019) |

---

## Domain events (new/extended)

| Event | When |
|-------|------|
| `ProfileVersionDeleted` | Successful version delete |
| `ProfileAssignmentsBumped` | Optional; after publish bump (or fold into ProfilePublished params) |

---

## Resolver (unchanged)

Device effective version: **nearest assignment on tree path** wins over parent assignment (018). Documented in UI only for 020.
