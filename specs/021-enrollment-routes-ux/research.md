# Research: Enrollment Routes — Controlled Onboarding Gateway

**Feature**: `021-enrollment-routes-ux` | **Date**: 2026-05-24

## R1 — Definition vs RuntimeState storage

**Decision**: Single aggregate `enrollment_routes` row with **two domain value objects** — `EnrollmentRouteDefinition` and `EnrollmentRouteRuntimeState` — mapped in `domain/` and `application/`; no second table in v1.

**Rationale**: Existing migration `000024` already holds both config (`name`, `target_node`, `mainappid`) and runtime (`qrcodekey`, `created_at`). Splitting tables adds join cost without immediate benefit. FR-014 satisfied at domain layer; physical split deferred until metrics volume requires it.

**Alternatives considered**:
- `enrollment_route_runtime` table — rejected for v1 (YAGNI).
- Full CQRS event store — rejected (overkill).

## R2 — Bootstrap app intent persistence

**Decision**: Add columns on `enrollment_routes`:
- `bootstrap_intent` VARCHAR(20) NOT NULL DEFAULT `'stable'` — enum: `stable` | `specific` | `latest`
- `bootstrap_application_id` INT REFERENCES `applications(id)` — package identity
- `bootstrap_version_id` INT REFERENCES `applicationversions(id)` — pinned or last-resolved version
- Keep `mainappid` as **resolved** `applicationversions.id` for legacy QR/Android parity (same as today).

**Rationale**: Headwind agents read `mainAppId` as application **version** id. Intent is admin metadata; resolution writes concrete version id to `mainappid` on each Save.

**Alternatives considered**:
- Store only intent + resolve at QR read time — rejected (non-deterministic audits, harder support).

## R3 — Stable channel resolution (no `is_recommended` today)

**Decision**: Migration `000030` adds `applicationversions.is_recommended BOOLEAN NOT NULL DEFAULT FALSE` with partial unique index: at most one recommended version per `applicationid`. Admin apps UI may set flag later; v1 seed via SQL for demo app optional.

Resolution:
- **stable** → version where `is_recommended = true` for selected `bootstrap_application_id`
- **latest** → highest `versioncode` among versions with non-empty `url` for that app
- **specific** → `bootstrap_version_id` fixed

**Rationale**: Clarification 2026-05-24 requires stable ≠ latest; DB had no channel flag.

**Alternatives considered**:
- Use `applications.latestversion` column — rejected (ambiguous vs “published stable”).

## R4 — Tree node Locked vs Inheritable

**Decision**: **Computed** at API read time, not stored:
- `placementKind: "locked"` — node has **no child folders** in `device_tree_nodes` (leaf folder).
- `placementKind: "inheritable"` — node has ≥1 child folder.

Device count: reuse `TreeNode.DeviceCount` (devices directly under node; plan documents this in picker).

Heavily loaded threshold: env `ENROLLMENT_TREE_HEAVY_DEVICE_THRESHOLD` default **500** (informational warning only).

**Rationale**: No schema change on `device_tree_nodes`; matches admin mental model (parent = container).

## R5 — Pending QR (client preview)

**Decision**: Reuse `frontend/src/features/devices/enrollmentQrQuery.ts` + local QR renderer:
- Pending: placeholder key `pending-preview` or empty; PNG/QR lib encodes JSON preview client-side **or** shows static sample with watermark “Pending”.
- Active: existing `loadQrImageObjectUrl` + `/public/qr/{qrcodekey}` with `create=1` and `deviceIdUseMode`.

Contract JSON for admin preview (documented in `enrollment-contract-payload.md`) mirrors public JSON shape minus server-only fields.

**Rationale**: Clarification — no server preview endpoint.

## R6 — Delete impact metrics

**Decision**:

| Metric | Source |
|--------|--------|
| `historicalEnrolledCount` | `COUNT(*) FROM devices WHERE enrollment_route_id = $id` |
| `enrollingNowCount` | `COUNT(*) FROM devices WHERE enrollment_route_id = $id AND created_at > NOW() - INTERVAL '24 hours'` |
| `activeQrScans7d` | `COUNT(*) FROM domain_events WHERE event_type = 'enrollment_route.qr_viewed' AND aggregate_id = $id::text AND created_at > NOW() - INTERVAL '7 days'` |

Emit `enrollment_route.qr_viewed` on successful `GET /rest/public/qr/:key` and `GET /rest/public/qr/json/:key` when key resolves to an enrollment route (extend `qrcode` or public handler).

**Rationale**: `domain_events` exists (`000026`); devices already store `enrollment_route_id`.

## R7 — Profile-free admin API

**Decision**: New response DTO `EnrollmentRouteView` strips `profileId`, `profileVersionId`, `profileVersionNumber`. Remove `validateBinding` requirement for `profileVersionId`. Deprecate `GET .../options/published-profile-versions` (410 or keep for admin tools only, not used by UI).

Legacy column `profile_version_id` remains for sync resolver; never mapped to UI DTO.

**Rationale**: FR-003 / strict vocabulary.

## R8 — Container placement acknowledgment

**Decision**: Column `container_placement_ack_at TIMESTAMPTZ NULL` on `enrollment_routes`. Set when admin checks “أبقِ هذا المجلد” on inheritable target. Cleared when target changes to locked node.

**Rationale**: FR-009a persistent warning until explicit ack.

## R9 — QR key rotation

**Decision**: `qrcodekey` **immutable** after first save unless admin uses future “Rotate QR” action (out of scope v1). Definition edits update provisioning JSON params but **not** the key.

**Rationale**: FR-015; avoids breaking printed QR sheets.

## R10 — Frontend shell

**Decision**: Replace `EnrollmentRouteEditorPage` route with list-only hub + `EnrollmentRouteDialog` (shadcn `Dialog` / `Sheet` md breakpoint). State machine in `enrollment-route-dialog-ux.md`.

**Rationale**: FR-006, FR-007.
