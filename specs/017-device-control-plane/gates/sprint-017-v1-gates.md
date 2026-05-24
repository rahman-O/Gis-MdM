# §20 Go/No-Go — Feature 017 (device control plane v1)

**Date**: 2026-05-23  
**Branch**: `017-device-control-plane`  
**Scope**: US1–US6 + Phase 9 polish (migrations through `000027`)

## Gate summary (blueprint §20.13)

| Transition | Key criterion | Result | Notes |
|------------|---------------|--------|-------|
| 0 → 1 | Staging + enroll baseline | **Go** | Existing 015 QR + Go sync |
| 1 → 2 | Tree + path + agent_id | **Go** | `device_tree` module, migrations 000019–000020 |
| 2 → 3 | QR + enroll → DB + tree ≤60s | **Go** | `enrollment_routes`, `tree_node_id` on create |
| 3 → 4 | Stable agent_id; safe number change | **Go** | `agent_id` UUID; enrollment_state |
| 4 → 5 | Publish → compile → sync artifact | **Go** | `profile_version_artifacts`, sync loader |
| 5 → 6 | enrollment_route + profile_version_id | **Go** | Routes bind published versions |
| 6 → 7 | Push via events; no storm | **Go** | `domain_events` worker (debounced) |
| 7 → 8 | Onboarding wizard | **Go** | Dashboard checklist + `/onboarding` |
| 8 → 9 | Staged rollout | **Deferred** | Out of 017 scope |
| 9 → 10 | Typed config extraction | **Deferred** | Future backlog |

## Automated checks (2026-05-23)

```text
go test ./internal/modules/device_tree/...     → ok (application)
go test ./internal/modules/profiles/...        → ok (application)
go test ./internal/modules/enrollment_routes/... → ok (application)
go test ./internal/modules/sync/...            → ok (application, postgres)
go build ./...                                 → ok
```

## Manual smoke (operator)

Run after `make migrate` through `000027`:

1. Device tree: create folder, move device, delete folder with relocate.
2. Profile: draft save, publish, impact dialog when ≥50 devices.
3. Enrollment route: create (requires published profile), QR panel.
4. Sync: enrolled device receives artifact payload (`profileId` / `profileVersionId` when set).
5. Migration smoke SQL in [quickstart.md](../quickstart.md) § Migration smoke.

## Known gaps / follow-ups

- Legacy `GET/PUT /configurations/:id` delegates via `profiles.legacy_configuration_id` only (no route alias).
- Full §20.2 per-transition narrative reports optional for future sprints.
- Android agent unchanged (v1 constraint).
