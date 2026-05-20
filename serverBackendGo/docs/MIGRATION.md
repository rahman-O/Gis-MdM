# Java → Go migration roadmap

Gradual migration from [`backend/`](../../backend/) (Java/JAX-RS) to [`serverBackendGo/`](../).

**Rule:** implement one module at a time; keep REST paths identical so [`frontend/`](../../frontend/) and Android agents keep working.

## Phases

| Phase | Modules | Java sources (indicative) | REST prefix |
|-------|---------|----------------------------|-------------|
| **1** | `auth` | `AuthResource`, `JWTAuthResource` | `/rest/public/auth`, `/rest/public/jwt` |
| **1b** | `signup`, `passwordreset` | `SignupResource`, `PasswordResetResource` | `/rest/public/signup`, `/rest/public/passwordReset` |
| **2** | `users`, `roles` | `UserResource`, `UserRoleResource` | `/rest/private/users`, `/rest/private/roles` |
| **3** | `customers`, `settings`, `hints`, `summary` | matching `*Resource.java` | `/rest/private/...` |
| **4** | `devices`, `groups` | `DeviceResource`, `GroupResource` | `/rest/private/devices`, `groups` |
| **5** | `applications`, `configurations`, `configfiles` | app/config resources | `/rest/private/...` |
| **6** | `files`, `icons`, `publicapi` | `FilesResource`, `PublicResource`, … | private + `/rest/public` |
| **7** | `sync`, `push`, `notifications`, `updates`, `qrcode` | agent/sync/push | public + private |
| **8** | `plugins/*` | `backend/plugins/*` | `/rest/plugins`, `/rest/plugin/main` |

## Per-module checklist

For each module under `internal/modules/<name>/`:

1. **domain** — entities from `backend/common/.../persistence/domain/`
2. **port** — repository interfaces (MyBatis DAO equivalents)
3. **application** — use cases (one file per endpoint group)
4. **adapter/http** — Gin handlers + `routes.go`
5. **adapter/persistence/postgres** — SQL implementations
6. **docs/parity/<name>.md** — endpoint parity table
7. Enable module flag in `.env` if using feature flags
8. Integration test against running Postgres (legacy schema)

## Auth (current focus)

See [parity/auth.md](parity/auth.md).

Enable: `MODULE_AUTH_ENABLED=true`

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
