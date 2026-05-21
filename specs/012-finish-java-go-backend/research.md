# Research: إكمال نقل الباكند Java → Go (012)

**Branch**: `012-finish-java-go-backend` | **Date**: 2026-05-21  
**Builds on**: [011 research](../011-complete-migration-gaps/research.md)

## R0 — Baseline from 011 (already implemented)

**Decision**: Treat Phase 9 P0 as **complete** in codebase: `internal/platform/push`, `app/scheduler.go`, `POST /rest/private/icon-files`.

**Rationale**: `JAVA-GO-MIGRATION-STATUS.md` and smoke on branch `011-complete-migration-gaps` / current tree.

**012 scope**: P1–P3 only + regression on P0.

---

## R1 — Device search filters (Java `DeviceSearchRequest`)

**Decision**: Extend Go `domain.SearchRequest` and SQL in `device_repo.go` to match Java `DeviceMapper.getAllDevices` / `countAllDevices` filter set used by React.

**Fields to implement** (from frontend `deviceService.ts` + Java domain):

- `status` — derive from `lastupdate` bands (green/yellow/red) same as list `statuscode` CASE
- `groupId`, `configurationId`, `value`, `fastSearch` — already partial
- `androidVersion`, `launcherVersion`, `mdmMode`, `kioskMode`, `installationStatus` — filter via `infojson` JSON operators (`->>`, `->`) where columns exist
- `sortBy`, `sortDir` — allow `LAST_UPDATE`, `NUMBER`; default `lower(number)`
- Date ranges: `dateFrom`, `dateTo`, `onlineEarlierMillis`, `onlineLaterMillis`, `enrollmentDateFrom`, `enrollmentDateTo`, `imeiChanged`

**Rationale**: FR-001; `FRONTEND-GO-BACKEND-INTEGRATION.md` documents ignored filters today.

**Alternatives**:

| Alternative | Rejected because |
|-------------|------------------|
| Strip filters in React | Breaks parity with Java WAR |
| Post-filter in Go app layer | Wrong pagination totals |

---

## R2 — Device detail telemetry (`infojson`)

**Decision**: On `GetByNumber`, unmarshal `devices.infojson` into nested `info` matching React `DeviceInfoView` (batteryLevel, model, androidVersion, applications[], files[] subset).

**Rationale**: FR-002; Java stores telemetry in JSON column.

**Alternatives**: Separate plugin deviceinfo call from UI — rejected; React uses `/private/devices/number/{n}` only.

---

## R3 — Plugin export formats

**Decision**: Copy Content-Type and column order from Java `DeviceInfoResource.export` and `DeviceLogResource` private search export during implementation (same as 011 R5).

**Rationale**: Legacy plugin UIs expect CSV/stream shape.

---

## R4 — Audit middleware

**Decision**: `internal/platform/audit` Gin middleware on `RouteGroups.Private`; async insert to `plugin_audit_log`.

**Rationale**: 011 R6; FR-005.

**Exclusions**: `/swagger/*`, `/rest/public/auth/login`, health.

---

## R5 — SyncResponseHook registry

**Decision**: `platform/synchooks.Registry`; plugins register in `module.Register`; `sync` application calls `Extend(response)` after building core payload.

**Rationale**: 011 R7; FR-006.

---

## R6 — Customer bootstrap

**Decision**: On `customers` create (`PUT` without id), copy template configuration + optional default device from Java `CustomerResource` logic (read method during tasks).

**Rationale**: FR-007; parity `customers.md` Partial.

---

## R7 — File quota and static serving

**Decision**:

- Quota: sum `uploadedfiles` sizes vs `customers.sizeLimit` on upload paths in `configfiles` and `files`.
- Static: register `GET /files/*filepath` on engine using `platform/storage.LocalStore` + `FILES_DIRECTORY` (parity Java servlet).

**Rationale**: FR-008, FR-009; agents use URLs from sync response.

**Alternatives**: Nginx-only static — rejected for dev parity and single-binary deploy story.

---

## R8 — Stats module

**Decision**: New module `internal/modules/stats` at `PUT /rest/public/stats`; map body to `usagestats` table (migration `000011` if table missing).

**Rationale**: FR-010; only missing core REST resource.

**Reference**: `com.hmdm.rest.resource.StatsResource`, `UsageStats` domain.

---

## R9 — Videos module

**Decision**: **Conditional** — implement `videos` module only if product confirms training videos in use; otherwise document ⊘ in `parity/videos.md` and skip module (FR-011).

**Default for planning**: scaffold module behind `MODULE_VIDEOS_ENABLED=false` with parity ⊘ note until UAT confirms need.

**Rationale**: Low React usage; full Java `VideosResource` is file storage under video directory.

---

## R10 — Updates APK + sendStats

**Decision**: Extend `updates` service to download remote APK when manifest says outdated; call stats port after apply (interface to `stats` module).

**Rationale**: FR-012; links to R8.

---

## R11 — Summary charts

**Decision**: Extend `summary` repo queries to populate `statusSummary`, install summaries, and per-config arrays when `devicestatuses` / related tables have rows; keep empty shape when not.

**Rationale**: FR-013; dashboard already consumes structure.

---

## R12 — MQTT / WebSocket

**Decision**: **Out of scope** (v1). Polling + DB queue sufficient per 011 R1 and spec Assumptions.

**Rationale**: Spec Out of Scope; no React WebSocket usage.

---

## Resolved Clarifications

All technical unknowns resolved; no NEEDS CLARIFICATION remaining for planning.
