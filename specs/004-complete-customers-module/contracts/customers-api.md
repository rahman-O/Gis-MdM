# API Contract: Customers (`/rest/private/customers`)

**Base path**: `/rest/private/customers`  
**Auth**: Session cookie and/or `Authorization: Bearer <jwt>`  
**Authorization**: Super administrator only (`Principal.SuperAdmin`)  
**Envelope**: `{ "status": "OK"|"ERROR", "message"?: string, "data"?: T }`

**Java reference**: `com.hmdm.rest.resource.CustomerResource`

---

### POST `/search`

Paginated customer search (control panel).

**Body** (`CustomerSearchRequest`):

```json
{
  "currentPage": 1,
  "pageSize": 100,
  "searchValue": "",
  "sortValue": null,
  "sortDirection": null,
  "accountType": null,
  "customerStatus": null
}
```

**Response data** (`PaginatedData<Customer>`):

```json
{
  "items": [{ "id": 2, "name": "Acme", "email": "a@acme.com" }],
  "totalItemsCount": 1
}
```

**Errors**: `403` / `error.permission.denied` if not super admin

---

### GET `/impersonate/{id}`

Assume org-admin identity for customer `{id}`.

**Response data**: `LoginUserPayload` (no password)

- Includes `authToken` when available or after server generates one
- `superAdmin`, `singleCustomer`, `userRole` as in login

**Errors**:

| message | When |
|---------|------|
| `error.permission.denied` | Not super admin |
| `error.notfound.customer.admin` | No org admin for customer |
| `error.internal.server` | Unexpected failure |

**Note**: Blocked when org admin has active `passwordResetToken` (legacy rule).

---

### PUT `/`

Create (no `id`) or update customer.

**Body**: `Customer` JSON (legacy fields; `id` null on create)

**Create success data**:

```json
{ "adminCredentials": "login/generatedPassword" }
```

**Update success**: `status: OK` (no data)

**Errors**: `error.duplicate.customer.name`, `error.duplicate.email`

---

### DELETE `/{id}`

Remove customer by id.

**Success**: `status: OK`

---

### GET `/{id}/edit`

Customer record for update UI (`findByIdForUpdate` parity).

**Success data**: full `Customer` object

---

### GET `/prefix/{prefix}/used`

**Success data**: `boolean` — `true` if prefix already assigned

---

## React consumers

| Client | Endpoints |
|--------|-----------|
| `frontend/src/features/customers/customersService.ts` | POST `/search`, GET `/impersonate/{id}` |
| `ControlPanelPage.tsx` | search + impersonate flow |

---

## Out of scope

| Item | Reason |
|------|--------|
| `GET /search`, `GET /search/{value}` | Deprecated in Java |
| Mailchimp subscribe | Spec assumption |
| Default devices on create | Until Phase 4 `devices` schema |
| Config/design copy on create | Until Phase 5 configurations |
