# Data Model: Phase 3 — Customers

**Date**: 2026-05-20  
**Storage**: PostgreSQL (`000001_init` + `000005_customers_extend.up.sql`)

## Entity: Customer (tenant)

| Field | Type | Rules |
|-------|------|-------|
| id | serial | PK |
| name | varchar(100) | NOT NULL, UNIQUE (case-insensitive check) |
| description | text | optional |
| master | boolean | `true` = platform master tenant; excluded from search |
| filesdir | varchar(200) | UNIQUE; generated UUID dir name on create |
| prefix | varchar(20+) | NOT NULL; unique across tenants (`isPrefixUsed`) |
| email | varchar(50) | optional; unique when set |
| lastlogintime | bigint | optional epoch ms |
| registrationtime | bigint | set on create |
| accounttype | int | default 0; filter in search |
| customerstatus | varchar(100) | optional; `customer.new` matches NULL |
| expirytime | bigint | optional |
| devicelimit | int | default 3 |
| deviceconfigurationid | int | optional default config for new devices |

**Search scope**: `master = false` only (non-master tenants).

## Entity: Organization admin (User)

Not owned by customers module but referenced for impersonation and create:

| Field | Use |
|-------|-----|
| id, login, name, email | Returned on impersonate |
| customerid | FK to customer |
| userroleid | `2` = org admin (`OrgAdminRoleID`) |
| authtoken | Generated if missing on impersonate |
| passwordresettoken | Blocks impersonation when set |
| password | Never returned in API |

## DTO: CustomerSearchRequest

| Field | Type | Notes |
|-------|------|-------|
| currentPage | int | 1-based |
| pageSize | int | LIMIT |
| searchValue | string | optional ILIKE filter |
| sortValue | string | `registrationTime`, `lastLoginTime`, `expiryTime`, or default `name` |
| sortDirection | string | `desc` for NULLS LAST on time fields |
| accountType | int | optional filter |
| customerStatus | string | optional filter |

## DTO: PaginatedCustomers (API)

| Field | Type |
|-------|------|
| items | Customer[] |
| totalItemsCount | long |

## DTO: Impersonation response

Maps to React `LoginUserPayload`: user identity + `authToken` + role flags; password omitted.

## DTO: Create customer result

| Field | Type |
|-------|------|
| adminCredentials | string | `login/password` combined (legacy map key) |

## Relationships

```text
customers (1) ──< users (many)
customers.master=false → listed in super-admin search
users.userroleid=2 → org admin per customer (first match)
```

## Validation rules

| Rule | Error key |
|------|-----------|
| Duplicate customer name | `error.duplicate.customer.name` |
| Duplicate customer email | `error.duplicate.email` |
| Email used by another user's login path | `error.duplicate.email` |
| Prefix already used | boolean `true` on prefix check endpoint |
| Non–super-admin access | `error.permission.denied` |
| Impersonate without org admin | `error.notfound.customer.admin` |

## Application operations

| Use case | DB / behavior |
|----------|----------------|
| Search | SELECT page + COUNT with filters |
| GetForEdit | SELECT customer by id (update view) |
| Save (create) | INSERT customer + INSERT org admin + return credentials |
| Save (update) | UPDATE customer + UPDATE org admin main details |
| Delete | DELETE customer by id (cascade users per FK) |
| PrefixUsed | EXISTS query on `prefix` |
| Impersonate | SELECT org admin; ensure token; map to login payload |

## Deferred (Phase 4/5)

- Insert 3 default devices per new customer
- Copy configurations / design settings
- Plugin customer-created events
