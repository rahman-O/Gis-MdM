# API Contract: Platform hardening (Phase 9 P2)

Covers non-REST or cross-cutting behaviors.

## Audit middleware

**Java reference**: `com.hmdm.plugins.audit.guice.servlet.AuditFilter`

| Aspect | Contract |
|--------|----------|
| Scope | All `/rest/private/*` except documented exclusions |
| Storage | INSERT `plugin_audit_log` |
| Fields | userId, login, customerId, action, payload (truncated), ipAddress, createTime |
| Failure | Must not block request if audit insert fails (log error) |

**Existing search API unchanged**: `POST /rest/plugins/audit/private/log/search`

## Sync hook registry

**Java reference**: `com.hmdm.rest.json.SyncResponseHook`, `SyncResource`

| Aspect | Contract |
|--------|----------|
| Registration | Plugins register hook at module init |
| Invocation | After core sync response built for `GET/POST .../sync/configuration/{deviceId}` |
| Merge | Shallow merge plugin JSON keys into `data` object |
| Order | Registration order; conflicts last-wins (document in parity) |

## Customers bootstrap

**Java reference**: `com.hmdm.rest.resource.CustomerResource` (create branch)

| Method | Path | Added behavior on create |
|--------|------|--------------------------|
| PUT | `/rest/private/customers` | Copy default configuration + optional seed devices from template |

**Response**: Unchanged envelope; side effect visible in DB.

## Files quota

| Module | Check |
|--------|-------|
| `configfiles` | Reject upload when customer storage exceeded |
| `files`, `icon-files` | Same `sizeLimit` logic |

**Error key**: Match Java storage error message key from `FilesResource`.

## Devices search enrichment

| Method | Path | Added query/body fields |
|--------|------|-------------------------|
| POST | `/rest/private/devices/search` | `mdmMode`, `launcherVersion`, `deviceStatuses` filters |

**Response**: Optional nested apps/files in configuration map when `enrich=true` or Java-default behavior.
