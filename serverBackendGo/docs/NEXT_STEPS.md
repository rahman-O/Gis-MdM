# الخطوات التالية — بعد إغلاق Auth

مرجع التسلسل: [`MIGRATION.md`](MIGRATION.md) · حالة Auth: [`AUTH_COMPLETE.md`](AUTH_COMPLETE.md) · parity: [`parity/auth.md`](parity/auth.md)

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

## Phase 5–8 (لاحقاً)

| Phase | وحدات |
|-------|--------|
| **5** | `applications`, `configurations`, `configfiles` |
| **6** | `files`, `icons`, `publicapi` |
| **7** | `sync`, `push`, `notifications`, `updates`, `qrcode` |
| **8** | `plugins/*` |

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
