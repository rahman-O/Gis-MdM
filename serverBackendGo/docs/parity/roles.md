# Roles API parity (`UserRoleResource`)

**Java**: `backend/server/src/main/java/com/hmdm/rest/resource/UserRoleResource.java`  
**Go**: `internal/modules/roles/adapter/http/handler.go`  
**Base**: `/rest/private/roles`

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/permissions` | Done | Non-superadmin permissions |
| GET | `/all` | Done | Roles list; org admin role hidden from edit list |
| PUT | `/` | Done | Create/update role + permission links |
| DELETE | `/:id` | Done | Cannot delete org admin role (id 2) |

**Related**: `GET /rest/private/users/roles` — same role ids/names for Users dropdown (`users` module).
