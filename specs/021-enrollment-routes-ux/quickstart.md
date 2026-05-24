# Quickstart: Enrollment Routes — Onboarding Gateway UX

**Feature**: `021-enrollment-routes-ux` | **Branch**: `021-enrollment-routes-ux`

## Prerequisites

- [017](../017-device-control-plane/quickstart.md) device tree + enrollment routes module
- `.env`: `MODULE_ENROLLMENT_ROUTES_ENABLED=true`, `MODULE_DEVICE_TREE_ENABLED=true`
- Migration through `000030` (bootstrap intent, `is_recommended`, QR view events)

```bash
cd serverBackendGo && make migrate
make dev
```

---

## Sprint 1 — Backend: profile-free API + intent

1. `POST /rest/private/enrollment-routes` without `profileVersionId` → `201`.
2. Body with `bootstrapIntent: "stable"`, `bootstrapApplicationId`, `targetNodeId` → `resolvedMainAppVersionId` populated.
3. `GET /rest/private/enrollment-routes/1` → no profile fields in JSON.
4. Stable with no recommended version → `400` `error.enrollment_route.stable_version_missing`.

**Pass**: `go test ./internal/modules/enrollment_routes/...`

---

## Sprint 2 — Options + impact APIs

1. `GET .../options/tree-nodes` → nodes with `placementKind`, `deviceCount`.
2. `GET .../options/bootstrap-apps` → versions with `isRecommended`.
3. `GET .../enrollment-routes/:id/impact` → three counters.
4. Hit `GET /rest/public/qr/{key}` → `domain_events` row `enrollment_route.qr_viewed`.

**Pass**: impact counts match seeded devices.

---

## Sprint 3 — Dialog hub (frontend)

1. Open **Enrollment routes** — list only; no full-page editor routes.
2. **New route** → dialog, Draft badge, Pending QR updates as you change fields (no network).
3. Save → Active badge, Active QR scannable, copy enabled.
4. No text containing “profile” or “policy” on screen.

**Pass**: SC-004 visual + network inspection.

---

## Sprint 4 — Tree picker + container warning

1. Pick inheritable (parent) folder → warning in picker; save allowed.
2. Overview shows persistent container warning until leaf selected or ack checkbox + save.
3. Pick leaf → warning clears.

**Pass**: User Story 3 scenarios.

---

## Sprint 5 — Edit, unsaved, delete

1. Edit saved route → Active + Unsaved; Pending QR on right; scan old QR still works.
2. Delete route with historical devices → impact shows 3 metrics → type route name → delete.
3. Devices remain; route gone from list.

**Pass**: FR-012, FR-013, FR-013a.

---

## Sprint 6 — Parity & i18n

1. Update `serverBackendGo/docs/parity/enrollment-routes-ux.md`.
2. AR/EN keys for new strings.
3. Remove onboarding gate “publish profile before route” on list **New** button.

**Pass**: quickstart sprints 1–6 green.

---

## Smoke curl

```bash
# Create gateway (no profile)
curl -s -b cookies.txt -H 'Content-Type: application/json' \
  -d '{"name":"Gate A","targetNodeId":2,"deviceIdentityMode":"imei","bootstrapIntent":"stable","bootstrapApplicationId":1}' \
  http://localhost:8080/rest/private/enrollment-routes | jq .

# Impact
curl -s -b cookies.txt http://localhost:8080/rest/private/enrollment-routes/1/impact | jq .
```
