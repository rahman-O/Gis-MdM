# Contract: Publish → Assignment Bump

**Feature**: `020-profile-workspace-maturity` | **Audience**: Backend + frontend publish flow

## Behavior

When admin publishes draft version `V_new` for profile `P`:

```text
[Confirm impact sheet]
        │
        ▼
PublishVersion(V_new)
        │
        ├── profiles.published_version_id ← V_new
        ├── draft pointer updated per 017 rules
        │
        └── FOR EACH row IN profile_tree_assignments WHERE profile_id = P
              SET profile_version_id = V_new
              Recompute devices in subtree → pending rollout
```

## Impact preview (before confirm)

Must list **every** active assignment row with:

- Folder name  
- Current assigned version number  
- Device count in subtree (same semantics as assignment impact)  

`requiresConfirmDialog = true` when `assignmentsToUpdate.length > 0` OR `deviceCount > 0` (same threshold policy as 018 publish).

## Idempotency

Publishing when assignments already on `V_new` (edge re-publish): no-op bump; still valid publish if new draft.

## Failure

If publish succeeds but bump fails: **transaction rollback** — publish must not commit without bump (single DB transaction).

## Events

- `ProfilePublished` with `{ versionNumber, assignmentsUpdated: N }`  
- Activity timeline: `profile.activity.published` + optional `profile.activity.assignmentsBumped`
