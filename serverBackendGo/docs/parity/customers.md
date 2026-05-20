# Customers API parity (`CustomerResource`)

**Java**: `com.hmdm.rest.resource.CustomerResource`  
**Go module**: `internal/modules/customers/`  
**Base path**: `/rest/private/customers`  
**Auth**: Super administrator only

| Method | Path | Status | Notes |
|--------|------|--------|-------|
| POST | `/search` | Done | Paginated `items` + `totalItemsCount` |
| GET | `/impersonate/{id}` | Done | Returns `UserView` JSON; no password |
| PUT | `/` | Partial | Create/update + org admin; **no** default devices/config copy (Phase 4/5) |
| DELETE | `/{id}` | Done | Non-master customers only |
| GET | `/{id}/edit` | Done | Customer by id |
| GET | `/prefix/{prefix}/used` | Done | Boolean in `data` |
| GET | `/search` | N/A | Deprecated in Java |
| GET | `/search/{value}` | N/A | Deprecated in Java |

**Out of scope**: Mailchimp subscribe on create.

**React**: `frontend/src/features/customers/customersService.ts` — search + impersonate.
