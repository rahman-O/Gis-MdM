# Java → Go migration roadmap

Gradual migration from [`backend/`](../../backend/) (Java/JAX-RS) to [`serverBackendGo/`](../).

**Rule:** implement one module at a time; keep REST paths identical so [`frontend/`](../../frontend/) and Android agents keep working.

**Governance:** See [`.specify/memory/constitution.md`](../../.specify/memory/constitution.md) (Gis-MdM v1.0.0) for architecture principles and quality gates.

## Status legend

| Mark | Meaning |
|------|---------|
| **done** | Go module implemented; parity doc updated; React shell paths work |
| **partial** | Routes exist; some behavior deferred (noted in parity doc) |
| **pending** | Scaffold or not started |

## Progress summary

| Phase | Status |
|-------|--------|
| **1** — auth | **done** |
| **1b** — signup, passwordreset (+ twofactor) | **done** |
| **2** — users, roles | **done** |
| **3** — customers, settings, hints, summary | **done** |
| **4** — devices, groups, configurations list | **done** |
| **5** — applications, configurations, configfiles | **done** |
| **6** — files, icons, publicapi | **done** |
| **7** — sync, push, notifications, updates, qrcode | **done** |
| **8** — plugins platform, audit, messaging, deviceinfo, devicelog | **done** |
| **9** — migration gap closure (push notifier, schedule cron, icon-files, …) | **partial** |

**Next:** Phase 9 remaining gaps (plugin exports, audit middleware, stats/videos, …) — see [`JAVA-GO-MIGRATION-GAP-ANALYSIS.md`](../../JAVA-GO-MIGRATION-GAP-ANALYSIS.md)

---

## Phases

| Phase | Status | Modules | Java sources (indicative) | REST prefix / parity |
|-------|--------|---------|----------------------------|----------------------|
| **1** | **done** | `auth` | `AuthResource`, `JWTAuthResource` | `/rest/public/auth`, `/rest/public/jwt` — [auth](parity/auth.md), [AUTH_COMPLETE](AUTH_COMPLETE.md) |
| **1b** | **done** | `signup`, `passwordreset`, `twofactor` | `SignupResource`, `PasswordResetResource`, 2FA webapp | `/rest/public/signup`, `/rest/public/passwordReset`, `/rest/private/twofactor` — [auth](parity/auth.md) |
| **2** | **done** | `users`, `roles` | `UserResource`, `UserRoleResource` | `/rest/private/users`, `/rest/private/roles` — [users](parity/users.md), [roles](parity/roles.md) |
| **3** | **done** | `customers`, `settings`, `hints`, `summary` | `CustomerResource`, `SettingsResource`, `HintResource`, `SummaryResource` | `/rest/private/...` — [customers](parity/customers.md), [hints](parity/hints.md), [settings](parity/settings.md), [summary](parity/summary.md) |
| **4** | **done** | `devices`, `groups`, `configurations` (list) | `DeviceResource`, `GroupResource`, `ConfigurationResource` (list) | [devices](parity/devices.md), [groups](parity/groups.md) — `GET /configurations/list` |
| **5** | **done** | `applications`, `configurations`, `configfiles` | `ApplicationResource`, `ConfigurationResource`, `ConfigurationFileResource` | [applications](parity/applications.md), [configurations](parity/configurations.md), [configfiles](parity/configfiles.md) |
| **6** | **done** | `files`, `icons`, `publicapi` | `FilesResource`, `IconResource`, `PublicResource` | [files](parity/files.md), [icons](parity/icons.md), [publicapi](parity/publicapi.md) |
| **7** | **done** | `sync`, `push`, `notifications`, `updates`, `qrcode` | agent/sync/push | [sync](parity/sync.md), [notifications](parity/notifications.md), [push](parity/push.md), [updates](parity/updates.md), [qrcode](parity/qrcode.md) |
| **8** | **done** | `plugins/*` | `backend/plugins/*` | `/rest/plugins`, `/rest/plugin/main` — [plugins-platform](parity/plugins-platform.md), [plugins-audit](parity/plugins-audit.md), [plugins-messaging](parity/plugins-messaging.md), [plugins-deviceinfo](parity/plugins-deviceinfo.md), [plugins-devicelog](parity/plugins-devicelog.md), [push](parity/push.md) |

### Phase 1–1b detail (done)

| Module | Parity | Notes |
|--------|--------|-------|
| `auth` | [parity/auth.md](parity/auth.md) | Session + JWT login/logout/options |
| `signup` | [parity/auth.md](parity/auth.md) § Signup | Simplified (no Mailchimp) |
| `passwordreset` | [parity/auth.md](parity/auth.md) § Password reset | Full public reset flow |
| `twofactor` | [parity/auth.md](parity/auth.md) § Two-factor | TOTP; blocks private routes until verified |

### Phase 3 detail (done)

| Module | Parity | Notes |
|--------|--------|-------|
| `customers` | [parity/customers.md](parity/customers.md) | Control panel search + impersonate; create **partial** (no default devices yet) |
| `settings` | [parity/settings.md](parity/settings.md) | Settings UI |
| `hints` | [parity/hints.md](parity/hints.md) | Hint history / enable / disable |
| `summary` | [parity/summary.md](parity/summary.md) | Real device status counts (Phase 4) |

### Phase 4 detail (done)

| Module | Parity | Notes |
|--------|--------|-------|
| `devices` | [parity/devices.md](parity/devices.md) | Search, CRUD, bulk, app settings; search enrichment **partial** |
| `groups` | [parity/groups.md](parity/groups.md) | List, CRUD, autocomplete |
| `configurations` | — | `GET /list` only (full CRUD in Phase 5) |

---

## Per-module checklist

For each module under `internal/modules/<name>/`:

1. **domain** — entities from `backend/common/.../persistence/domain/`
2. **port** — repository interfaces (MyBatis DAO equivalents)
3. **application** — use cases (one file per endpoint group)
4. **adapter/http** — Gin handlers + `routes.go`
5. **adapter/persistence/postgres** — SQL implementations
6. **docs/parity/<name>.md` — endpoint parity table
7. Enable module flag in `.env` if using feature flags
8. Integration test against running Postgres (legacy schema)

---

## Swagger UI coverage

Swagger lists only routes with `// @Router` comments in handlers. After adding or changing endpoints:

```bash
cd serverBackendGo && make swagger   # then restart make dev
```

| Tag in Swagger | Phase | Notes |
|----------------|-------|-------|
| Authentication | 1 | auth + jwt |
| Signup | 1b | public |
| PasswordReset | 1b | public |
| TwoFactor | 1b | private, Bearer |
| Users, Roles | 2 | private |
| Customers, Hints, Settings, Summary | 3 | private |
| Devices, Groups, Configurations | 4 | private device/group routes + `GET /configurations/list` |
| Applications, … | 5+ | not implemented yet |

---

## Post-auth sequence (high priority)

See [NEXT_STEPS.md](NEXT_STEPS.md). **Current focus:** Phase 4 — **devices** / **groups**.

---

## References

- Legacy flat Go experiment: [`backend-go/`](../../backend-go/) (do not copy blindly; use Java as source of truth)
- DAO layer: `backend/common/src/main/java/com/hmdm/persistence/`
- REST resources: `backend/server/src/main/java/com/hmdm/rest/resource/`

## Module layout

```
internal/modules/<name>/
├── module.go
├── domain/
├── port/
├── application/
└── adapter/
    ├── http/
    └── persistence/postgres/
```
