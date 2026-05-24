# الخطوات التالية — بعد إغلاق Auth

مرجع التسلسل: [`MIGRATION.md`](MIGRATION.md) · حالة Auth: [`AUTH_COMPLETE.md`](AUTH_COMPLETE.md) · parity: [`parity/auth.md`](parity/auth.md)

---

## 014 — React ↔ Go integration (2026-05-21)

| Wave | Scope | Status |
|------|--------|--------|
| A (P1) | Settings `000015`, config `settingsjson`/CAP, icons upload UI | **منجز** |
| B (P2) | `sync/info` → `devicestatuses`, `stats` module | **منجز** |
| C (P3) | Updates apply, hints mark shown | **منجز** |
| Optional | Device list `infojson` columns, monthly enrollment chart | مفتوح |

Spec: [`specs/014-complete-frontend-go-integration/`](../specs/014-complete-frontend-go-integration/)

---

## 017 — Device control plane (2026-05-23)

| Area | Status |
|------|--------|
| Device tree + move/filter | **منجز** |
| Profiles (draft/publish/artifact) | **منجز** |
| Enrollment routes + QR | **منجز** |
| Sync from artifacts | **منجز** |
| Onboarding checklist/wizard | **منجز** |
| Nav: Profiles + Enrollment routes | **منجز** |

Spec: [`specs/017-device-control-plane/`](../specs/017-device-control-plane/) · Parity: [`parity/device-control-plane.md`](parity/device-control-plane.md) · Gates: [`gates/sprint-017-v1-gates.md`](../specs/017-device-control-plane/gates/sprint-017-v1-gates.md)

**Deploy:** `cd serverBackendGo && make migrate` (through `000027`), enable `MODULE_DEVICE_TREE_ENABLED`, `MODULE_PROFILES_ENABLED`, `MODULE_ENROLLMENT_ROUTES_ENABLED`.

---

## 018 — Profile rollout & operations (2026-05-23)

| Area | Status |
|------|--------|
| Tree assignment (nearest folder wins) | **منجز** |
| Version navigation + fork draft | **منجز** |
| Rollout status grid + recompute | **منجز** |
| Profile enable/disable | **منجز** |
| Sync effective profile (tree > route) | **منجز** |

Spec: [`specs/018-profile-rollout-ops/`](../specs/018-profile-rollout-ops/) · Parity: [`parity/profile-rollout-ops.md`](parity/profile-rollout-ops.md)

**Deploy:** `make migrate` through `000028`, `MODULE_PROFILE_ROLLOUT_ENABLED=true` (with 017 flags).

---

## 021 — Enrollment Routes UX (2026-06-01)

| Area | Status |
|------|--------|
| Dialog-based enrollment route CRUD | **منجز** |
| Bootstrap intent resolution (stable/specific/latest) | **منجز** |
| Target node picker with placement kind | **منجز** |
| Delete with multi-dimensional impact | **منجز** |
| Dual-column QR preview (Pending/Active) | **منجز** |
| Profile-free enrollment (no profile prerequisite) | **منجز** |

Spec: [`specs/021-enrollment-routes-ux/`](../specs/021-enrollment-routes-ux/) · Parity: [`parity/enrollment-routes-ux.md`](parity/enrollment-routes-ux.md)

**Deploy:** `make migrate` through `000030`, `MODULE_ENROLLMENT_ROUTES_ENABLED=true` (with 017 flags).

---

## منجز

| المرحلة | الوحدات | الحالة |
|---------|---------|--------|
| **1** | `auth` (session, JWT, options, RSA) | منجز |
| **1b** | `signup`, `passwordreset`, `twofactor`, `users/current` | منجز |

---

## الأولوية العالية (تسلسلي — ما بعد Auth)

هدف كل خطوة: **تسجيل دخول → شاشة تعمل في React** بدون Java.

| # | الوحدة | Endpoint رئيسي | لماذا الآن | الحالة |
|---|--------|----------------|------------|--------|
| **1** | `summary` | `GET /rest/private/summary/devices` | Dashboard بعد login | **منجز** |
| **2** | `settings` | `GET/POST /rest/private/settings/*` | تبويب الإعدادات + 2FA toggle | **منجز** |
| **2b** | `users/roles` | `GET /rest/private/users/roles` | أدوار في Settings | **منجز** |
| **3** | `hints` | `GET /rest/private/hints/*` | تلميحات الواجهة | **منجز** |
| **4** | `users` (إكمال) | `PUT /current`, `/details`, `GET /all` | الملف الشخصي + إدارة مستخدمين | **منجز** |
| **5** | `roles` | `/rest/private/roles/*` | إدارة الأدوار + صلاحيات | **منجز** |

---

## Phase 2 — مستخدمون وأدوار (منجز)

| الوحدة | Java | REST | parity |
|--------|------|------|--------|
| `users` | `UserResource` | `/rest/private/users` | [`parity/users.md`](parity/users.md) |
| `roles` | `UserRoleResource` | `/rest/private/roles`, `/users/roles` | [`parity/roles.md`](parity/roles.md) |

---

## Phase 3 — عملاء وإعدادات (منجز)

| الوحدة | REST | parity |
|--------|------|--------|
| `customers` | `/rest/private/customers` | [`parity/customers.md`](parity/customers.md) — **منجز** |
| `settings` | `/rest/private/settings/*` | منجز |
| `hints` | `/rest/private/hints/*` | [`parity/hints.md`](parity/hints.md) |
| `summary` | `/rest/private/summary/devices` | منجز |

---

## Phase 4 — أجهزة ومجموعات (منجز)

| الوحدة | REST | parity |
|--------|------|--------|
| `devices` | `/rest/private/devices` | [`parity/devices.md`](parity/devices.md) |
| `groups` | `/rest/private/groups` | [`parity/groups.md`](parity/groups.md) |
| `configurations` | `GET /rest/private/configurations/list` | قائمة للواجهة (CRUD كامل في Phase 5) |

Migration: `000006_devices_groups_core`. تشغيل: `make dev` ثم `make swagger`.

---

## Phase 5 — تطبيقات وإعدادات (منجز)

| الوحدة | REST | parity |
|--------|------|--------|
| `applications` | `/rest/private/applications` | [`parity/applications.md`](parity/applications.md) |
| `configurations` | `/rest/private/configurations` | [`parity/configurations.md`](parity/configurations.md) — يشمل `GET /list` (Phase 4) |
| `configfiles` | `POST /rest/private/config-files` | [`parity/configfiles.md`](parity/configfiles.md) |

Migration: `000007_applications_configurations_core`. صلاحيات: `applications`, `configurations`.

---

## Phase 6 — ملفات وأيقونات وواجهة عامة (منجز)

| الوحدة | REST | parity |
|--------|------|--------|
| `files` | `/rest/private/web-ui-files` | [`parity/files.md`](parity/files.md) |
| `icons` | `/rest/private/icons` | [`parity/icons.md`](parity/icons.md) |
| `publicapi` | `/rest/public` | [`parity/publicapi.md`](parity/publicapi.md) |

Migration: `000008_files_icons_core`. صلاحيات: `files`, `edit_files`.

---

## Phase 7 — وكيل وتزامن (منجز)

| الوحدة | REST | parity |
|--------|------|--------|
| `sync` | `/rest/public/sync` | [`parity/sync.md`](parity/sync.md) |
| `notifications` | `/rest/notifications`, `/rest/notification/polling` | [`parity/notifications.md`](parity/notifications.md) |
| `push` | `/rest/private/push`, `/rest/plugins/push` | [`parity/push.md`](parity/push.md) |
| `updates` | `/rest/private/update` | [`parity/updates.md`](parity/updates.md) |
| `qrcode` | `/rest/public/qr` | [`parity/qrcode.md`](parity/qrcode.md) |

Migration: `000009_agent_push_notifications`.

---

## Phase 8 — plugins (منجز)

| الوحدة | REST | parity |
|--------|------|--------|
| `plugins/platform` | `/rest/plugin/main` | [`parity/plugins-platform.md`](parity/plugins-platform.md) |
| `plugins/audit` | `/rest/plugins/audit` | [`parity/plugins-audit.md`](parity/plugins-audit.md) |
| `plugins/messaging` | `/rest/plugins/messaging` | [`parity/plugins-messaging.md`](parity/plugins-messaging.md) |
| `plugins/deviceinfo` | `/rest/plugins/deviceinfo` | [`parity/plugins-deviceinfo.md`](parity/plugins-deviceinfo.md) |
| `plugins/devicelog` | `/rest/plugins/devicelog` | [`parity/plugins-devicelog.md`](parity/plugins-devicelog.md) |
| `plugins/push` (schedule) | `/rest/plugins/push/private/task*` | [`parity/push.md`](parity/push.md) |

Migration: `000010_plugins_core`.

---

## checklist لكل وحدة جديدة

1. `domain/` — نماذج من Java DTO
2. `port/` — واجهات repository
3. `application/` — use cases
4. `adapter/http/` — Gin + `routes.go`
5. `adapter/persistence/postgres/` — SQL
6. `docs/parity/<name>.md`
7. اختبار: `go test ./internal/modules/<name>/...`
8. تحديث هذا الملف (عمود **الحالة**)

---

## تشغيل وتحقق

```bash
cd serverBackendGo
./scripts/db-up.sh
make dev
# frontend
cd ../frontend && npm run dev
```

تحقق يدوي: login → `/dashboard` → Network بدون 404 على endpoints المرحلة الحالية.

---

## Git / دمج

- الفرع الحالي: `feature/config-editor-density-backend-fixes`
- بعد كل مرحلة: commit + push
- PR إلى `main` عند استقرار Dashboard + Settings
