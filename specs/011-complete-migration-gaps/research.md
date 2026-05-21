# Research: Phase 9 — Complete Migration Gaps

**Branch**: `011-complete-migration-gaps` | **Date**: 2026-05-21

## R1 — Push delivery: FCM vs polling vs MQTT

**Decision**: Implement **database queue + agent long-poll** as primary path (already in Go Phase 7); do **not** add MQTT client in Phase 9 v1.

**Rationale**:

- Go `notifications.QueueRepository` already writes `pushmessages` / `pendingpushes` — same tables Java `PushSenderPolling` uses.
- Phase 7 quickstart proves `messageType: configUpdated` works for agents.
- Java `PushService.send()` calls both MQTT and polling; agents choose transport via device configuration — polling path is sufficient for React + agent smoke without ActiveMQ dependency.

**Alternatives considered**:

| Alternative | Rejected because |
|-------------|------------------|
| Firebase Admin SDK | Not used in core Java server module; adds ops burden |
| MQTT in Go | Requires `mqtt.server.uri` infra; duplicate of polling for MVP |
| HTTP callback to devices | Not Headwind model |

**Follow-up**: Document `MQTT_SERVER_URI` optional Phase 9b if production requires instant delivery without poll interval.

---

## R2 — Shared notifier placement

**Decision**: `internal/platform/push` with interface `Notifier` injected into module `Register()`.

**Rationale**: Constitution allows `internal/platform/` for cross-cutting concerns; mirrors Java `com.hmdm.notification.PushService` singleton.

**Alternatives**: Embed in `notifications` module only — rejected because configurations/devices shouldn’t import notifications HTTP adapters.

---

## R3 — Schedule worker interval

**Decision**: Default **60 seconds** tick (`PUSH_SCHEDULE_INTERVAL_SEC`), matching Java `submitRepeatableTask(..., 1, 1, TimeUnit.MINUTES)`.

**Rationale**: Direct parity with `PushScheduleTaskModule`.

---

## R4 — Icon file upload

**Decision**: Implement on **`icons`** module (same team as `IconResource`) at path `/rest/private/icon-files` (separate route group registration).

**Rationale**: Java separates `IconResource` vs `IconFileResource` but shares `uploadedfiles` + `files.directory`; Go migration `000008` already has `uploadedfiles` table.

**Validation**: Square image check + Scalr 144px — use Go `image` + `github.com/disintegration/imaging` or stdlib resize (evaluate in tasks).

---

## R5 — Deviceinfo / devicelog export format

**Decision**: Match Java response **Content-Type** and column order from `DeviceInfoResource.export` / `DeviceLogResource` search export (read Java during implementation).

**Rationale**: Angular/React plugin UIs parse legacy format.

**Schema**: Audit liquibase in `backend/plugins/deviceinfo` for GPS/WiFi tables; add migration only for tables absent in dev DB.

---

## R6 — Audit auto-capture

**Decision**: Gin middleware on `RouteGroups.Private` recording method, path, user, customer, status code; async insert to `plugin_audit_log`.

**Rationale**: Java `AuditFilter` is servlet-level; middleware is idiomatic Go equivalent.

**Exclusions**: `GET` health, `/swagger/*`, static assets.

---

## R7 — SyncResponseHook

**Decision**: Registry interface in `platform/synchooks`; plugins call `Register(hook)` from `module.Register`; `sync` application merges `hook.Extend(response)` after core build.

**Rationale**: Mirrors Guice multibind `Set<SyncResponseHook>` without DI framework.

---

## R8 — Stats & videos

**Decision**:

- **stats**: New small module `internal/modules/stats` — table `usagestats` (Java `UsageStatsMapper`).
- **videos**: New module `videos` OR extend `publicapi` with `VIDEO_DIRECTORY` env — prefer **separate module** for parity clarity.

**Rationale**: Keeps `publicapi` focused on branding/upload; videos are optional training content.

---

## R9 — Out of scope confirmation

**Decision**: Mailchimp, xtra Angular UI, `PublicFilesResource`, user impersonate/superadmin remain **out of scope** for Phase 9.

**Rationale**: Spec §Out of Scope; React doesn’t call those paths per parity docs.
