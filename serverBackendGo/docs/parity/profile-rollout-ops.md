# Parity: Profile Rollout & Operations

**Feature**: `018-profile-rollout-ops` | **Status**: Implemented (2026-05-23)

## Tree assignments

| Method | Path | Status |
|--------|------|--------|
| GET | `/profiles/:id/assignments` | Implemented |
| GET | `/profiles/:id/assignments/impact?treeNodeId=` | Implemented |
| PUT | `/profiles/:id/assignments` | Implemented |
| DELETE | `/profiles/:id/assignments/:assignmentId` | Implemented |

## Version navigation

| Method | Path | Status |
|--------|------|--------|
| GET | `/profiles/:id/versions` | Implemented |
| POST | `/profiles/:id/versions/:versionId/fork-draft` | Implemented |

## Rollout status

| Method | Path | Status |
|--------|------|--------|
| GET | `/profiles/:id/rollout/devices` | Implemented |
| POST | `/profiles/:id/rollout/recompute` | Implemented |

## Profile enable/disable

| Method | Path | Status |
|--------|------|--------|
| POST | `/profiles/:id/disable` | Implemented |
| POST | `/profiles/:id/enable` | Implemented |

## Effective profile (sync)

- Tree assignment wins over enrollment route when both exist (`profiles/application/resolver.go`).
- Sync loads artifact by resolved `profile_version_id`.
- Device columns: `target_profile_version_id`, `applied_profile_version_id`, `profile_rollout_status`, `profile_rollout_reason`.

## Rollout status values

| Status | Meaning |
|--------|---------|
| `pending` | Target set; sync not confirmed |
| `installed` | Applied version matches target; apps OK |
| `partial` | Version applied; app install issues in device info |
| `failed` | Reserved for explicit failure paths |

## Config

- `MODULE_PROFILE_ROLLOUT_ENABLED` (default `true` when profiles enabled)
- Migration `000028_profile_rollout_ops`

## Tests

```bash
go test ./internal/modules/profiles/application/...
go build ./...
```
