# API Contract: Audit plugin

**Base path**: `/rest/plugins/audit`  
**Auth**: Bearer JWT on `/rest/plugins` group  
**Java reference**: `com.hmdm.plugins.audit.rest.AuditResource`

## POST `/private/log/search`

**Permission**: `plugin_audit_access`  
**Body** (`AuditLogFilter`):

```json
{
  "pageNum": 1,
  "pageSize": 50,
  "dateFrom": 0,
  "dateTo": 0,
  "userId": null,
  "login": null,
  "action": null
}
```

**Response data**: `PaginatedData<AuditLogRecord>` — `{ items, totalItemsCount }` or legacy list+count shape matching Java `PaginatedData`.

**Errors**: `error.permission.denied`, validation, `error.internal.server`

**Partial**: Automatic request/response audit capture (servlet filter) not implemented in Go Phase 8.
