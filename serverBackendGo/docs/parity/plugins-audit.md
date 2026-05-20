# Parity: Audit plugin

**Java**: `com.hmdm.plugins.audit.rest.AuditResource`  
**Go**: `internal/modules/plugins/audit`

| Method | Path | Permission |
|--------|------|------------|
| POST | `/rest/plugins/audit/private/log/search` | `plugin_audit_access` |

**Partial**: Servlet `AuditFilter` auto-logging not ported; search API + seeded rows only.
