# Quickstart: Profile Rollout & Operations

**Feature**: `018-profile-rollout-ops` | **Prerequisite**: [017 quickstart](../017-device-control-plane/quickstart.md) complete (tree, profiles, publish, routes)

## Environment

```bash
# serverBackendGo/.env
MODULE_PROFILES_ENABLED=true
MODULE_PROFILE_ROLLOUT_ENABLED=true
```

## Migrate

```bash
cd serverBackendGo
make db-up
make migrate   # through 000028_profile_rollout_ops
make dev
cd ../frontend && npm run dev
```

## Sprint 1 — Version navigation (US2)

1. Open `/profiles/:id/edit` — version dropdown shows draft + published history.
2. Select published v1 — editor read-only.
3. «Fork draft from v1» — new draft opens editable.
4. Switch version with unsaved draft — confirm dialog appears.

## Sprint 2 — Tree assignment (US1)

1. Publish profile v2.
2. In Profile editor → **Tree assignment** — pick folder «Test branch» with 3 devices.
3. Confirm impact if ≥50 devices.
4. Verify `GET /profiles/:id/assignments` lists folder + v2.
5. Devices show `pending` in rollout tab.

## Sprint 3 — Rollout status (US3)

1. Trigger sync on test device (agent or simulate).
2. Rollout tab: device moves to `installed` or `partial`/`failed` with reason.
3. Remove one app from device / force mismatch — status becomes `partial`.
4. Manual **Refresh status** updates within 2 minutes.

## Sprint 4 — Disable / enable (US4)

1. **Disable** profile — new sync does not advance target from disabled profile.
2. **Enable** — devices return to `pending`; next sync applies policy.

## Overlap test (edge case)

1. Assign v2 on parent folder; assign v3 on child folder.
2. Device in child gets v3 (nearest wins).
3. Device in parent-only subfolder gets v2.

## Go tests

```bash
cd serverBackendGo
go test ./internal/modules/profiles/application/... -run 'Resolver|Rollout|Assignment'
go test ./internal/modules/sync/application/... -run 'EffectiveProfile'
```

## Parity

Update [`serverBackendGo/docs/parity/profile-rollout-ops.md`](../../serverBackendGo/docs/parity/profile-rollout-ops.md) after each sprint.
