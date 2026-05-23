# Users API parity (`UserResource`)

**Java**: `backend/server/src/main/java/com/hmdm/rest/resource/UserResource.java`  
**Go**: `internal/modules/users/adapter/http/handler.go`  
**Base**: `/rest/private/users`

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| GET | `/current` | Done | Full user view, no password/authToken |
| PUT | `/details` | Done | Profile name/email |
| PUT | `/current` | Done | Self password change (MD5 hex) |
| GET | `/all` | Done | Requires `settings` or super admin |
| PUT | `/` | Done | Create/update; org admin or super admin |
| DELETE | `/other/:id` | Done | Org admin or super admin |
| GET | `/roles` | Done | Role dropdown list |
| GET | `/{id}` | Out of scope | Unused by React |
| GET | `/impersonate/:id` | Out of scope | Deferred |
| GET/PUT | `/superadmin/*` | Out of scope | Deferred |
