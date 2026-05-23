# API Contract: Platform hardening (012 P2)

**Scope**: Cross-cutting behavior + customer bootstrap  
**Java references**: `AuditFilter`, `SyncResource` hooks, `CustomerResource`  
**Go**: `platform/audit`, `platform/synchooks`, `modules/customers`, `modules/sync`

## Audit auto-capture

**Trigger**: Any successful `POST`/`PUT`/`DELETE` on `/rest/private/*` (configurable list).

**Storage**: `plugin_audit_log` — same schema as manual audit search.

**Not exposed as new REST** — observable via existing:

`POST /rest/plugins/audit/private/log/search`

## Sync response hooks

**Trigger**: `GET`/`POST` `/rest/public/sync/configuration/{deviceId}` and related sync paths.

**Behavior**: After core `SyncResponse` built, merge JSON keys from each registered plugin hook.

**Registration**: Plugins call `synchooks.Register(name, hook)` in `module.Register`.

## Customer bootstrap

**Endpoint**: `PUT /rest/private/customers` (create — no id in body).

**Behavior**: After customer + org admin created, copy default configuration (and optional template device) per Java.

**Response**: Unchanged Headwind envelope.

## Configurations / applications fields (012 P3 partial)

**Endpoint**: `PUT /rest/private/configurations`

**Critical fields** for QR/enrollment (must persist if sent):

- `qrCodeKey`, `baseUrl`, `mainAppId`, `contentAppId`, design colors, `defaultFilePath`

Document supported subset in parity after audit of handler.
