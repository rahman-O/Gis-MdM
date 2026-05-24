# Data Model: Profile Hub & Enrollment UX

**Feature**: `019-profile-hub-ux` | **Date**: 2026-05-23

## Overview

Primarily a **presentation and API aggregation** feature. Persistent schema changes are **minimal** (optional index + new domain events). Most entities are **read models** for the workspace UI.

**Depends on**: 017 (`profiles`, `enrollment_routes`, `device_tree`, `domain_events`), 018 (`profile_tree_assignments`, rollout columns, `profiles.enabled`).

---

## Existing tables (unchanged in v1)

| Table | Role |
|-------|------|
| `profiles` | Identity, `enabled`, `draft_version_id`, `published_version_id` |
| `profile_versions` | Draft/published payloads |
| `profile_tree_assignments` | Policy target per tree folder |
| `enrollment_routes` | `profile_version_id` **deprecated in UI**; column retained |
| `devices` | Rollout status columns (018) |
| `domain_events` | Activity timeline source |

### Optional migration `000029_profile_hub_activity_idx` (recommended)

```sql
CREATE INDEX IF NOT EXISTS domain_events_aggregate_type_idx
    ON domain_events (aggregate_id, event_type, created_at DESC);
```

No new business tables required.

---

## Domain / API read models

### ProfileHealth

| Field | Type | Description |
|-------|------|-------------|
| `health` | enum | `healthy` \| `warning` \| `error` \| `draft_only` |
| `reasons` | string[] | Machine keys for UI (`no_assignment`, `rollout_failures`, …) |

### ProfileListItem (extended)

Extends 017 list row:

| Field | Type |
|-------|------|
| `enabled` | bool |
| `health` | ProfileHealth.health |
| `badges` | string[] | `no_assignment`, `disabled`, `draft_changes`, `rollout_issues`, `stale` |
| `assignmentCount` | int |
| `rolloutFailureCount` | int |
| `publishedVersion` | int? |
| `draftVersionId` | int? |

### ProfileSummary (workspace cockpit)

| Field | Type |
|-------|------|
| `id`, `name`, `description` | |
| `enabled` | bool |
| `health` | ProfileHealth |
| `lifecycle` | `draft` \| `published` \| `disabled` |
| `publishedVersionId`, `publishedVersionNumber` | int? |
| `draftVersionId` | int? |
| `assignmentCount`, `assignedFolderNames` | |
| `rollout` | `{ pending, installed, partial, failed, total }` |
| `pinnedSettings` | `{ kioskMode, mainAppName, appCount, lastPublishedAt }` |
| `hasUnpublishedDraft` | bool |
| `canPublish` | bool |

### ProfileActivityEvent

| Field | Type |
|-------|------|
| `id` | int64 |
| `eventType` | string |
| `summary` | string (localized template key + params) |
| `occurredAt` | ISO8601 |
| `actorUserId` | int? (from payload if present) |
| `metadata` | object |

### EnrollmentRoute (API write model change)

| Field | Change |
|-------|--------|
| `profileVersionId` | **Optional** on create/update; omitted in UI |

---

## Event types (activity)

| event_type | When |
|------------|------|
| `ProfilePublished` | Existing (017/018) |
| `ProfileAssignmentChanged` | Existing (018) |
| `ProfileEnabled` | New — enable service |
| `ProfileDisabled` | New — enable service |

`aggregate_id`: `"{profileId}"` or `"profile:{profileId}"` — normalize reader to accept both.

---

## UI state (client-only)

### ProfileWorkspaceState

| Field | Type |
|-------|------|
| `profileId` | number |
| `section` | `overview` \| `assignments` \| `rollout` \| `versions` \| `editor` \| `activity` |
| `editorDirty` | boolean |
| `secondaryPanel` | `null` \| `publish-impact` |
| `createWizardStep` | `null` \| `publish-first` \| `assign` |

Not persisted server-side.

---

## Validation rules

| Rule | Layer |
|------|--------|
| Enrollment route create without `profileVersionId` | Allowed |
| Workspace publish without draft changes | `canPublish: false` |
| Assignment requires published version | 018 (unchanged) |
| Health `error` when disabled + devices still targeted | Computed |

---

## State transitions (health)

```text
draft_only → warning (first publish)
warning → healthy (assign + stable rollout)
healthy → warning (remove all assignments OR new draft)
healthy → error (rollout failures threshold > 0)
any → error/warning when disabled (policy off)
```

Threshold: **≥1** failed device in profile rollout scope → `rollout_issues` badge; **≥1** failed → health contributes to `error` if also enabled and assigned.
