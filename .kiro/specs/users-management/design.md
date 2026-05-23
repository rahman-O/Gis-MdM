# Design Document: Users Management

## Overview

The Users Management feature adds a fully functional `/users` page to the HMDM frontend. It replaces the current placeholder with a table of all administrator accounts, a `Dialog`-based form (react-hook-form + zod) for creating and editing users, and an `AlertDialog` confirmation prompt for deletion. Role options are fetched from `GET /rest/private/roles` each time the form opens and rendered via a `Select` component. The password field is shown only in create mode.

All data flows through a dedicated `userService` module that wraps the shared `apiClient`. The UI is built exclusively from shadcn/ui components, following the same patterns established by the configurations-management feature.

### Key Backend API Findings

- **List users**: `GET /rest/private/users` — returns array of `User` objects
- **Create user**: `POST /rest/private/users` — body is `UserPayload`
- **Update user**: `PUT /rest/private/users/{id}` — body is `UserPayload`
- **Delete user**: `DELETE /rest/private/users/{id}`
- **List roles**: `GET /rest/private/roles` — returns array of `Role` objects
- **Response shape**: all responses are wrapped in `{ status: "OK" | "ERROR", data: T }` — unwrapped via the existing `unwrapHmdmData` helper from `src/services/hmdmEnvelope.ts`
- **Navigation**: Users entry already present in `navItems.ts` and route already registered in `App.tsx`

---

## Architecture

The feature follows the existing feature-slice pattern:

```
src/features/users/
  UsersPage.tsx        # Page component, route /users (replaces placeholder)
  userService.ts       # API calls via apiClient
  UserForm.tsx         # Dialog form (create + edit modes)
  types.ts             # TypeScript interfaces
```

All required shadcn/ui components are already scaffolded:
`Table`, `Dialog`, `Form`, `Input`, `Select`, `Button`, `AlertDialog`, `Skeleton`, `DropdownMenu`

### Data Flow

```mermaid
flowchart TD
    A[UsersPage] -->|mount| B[userService.getUsers]
    B -->|GET /rest/private/users| C[Backend]
    C -->|User[]| B
    B --> A

    A -->|openCreate| D[UserForm - create mode]
    D -->|form open| E[userService.getRoles]
    E -->|GET /rest/private/roles| C
    D -->|valid payload| F[userService.createUser]
    F -->|POST /rest/private/users| C

    A -->|openEdit user| G[UserForm - edit mode]
    G -->|form open| E
    G -->|valid payload + id| H[userService.updateUser]
    H -->|PUT /rest/private/users/:id| C

    A -->|openDelete user| I[DeleteDialog]
    I -->|user.id| J[userService.deleteUser]
    J -->|DELETE /rest/private/users/:id| C

    F & H & J -->|success| A
    A -->|refresh| B
```

---

## Components and Interfaces

### UsersPage

Top-level page component rendered at `/users`. Owns all state:

| State | Type | Purpose |
|---|---|---|
| `users` | `User[]` | Full list from API |
| `loading` | `boolean` | List fetch in flight |
| `error` | `string \| null` | List fetch error |
| `formMode` | `'create' \| 'edit' \| null` | Controls form dialog visibility |
| `selectedUser` | `User \| null` | User being edited |
| `userToDelete` | `User \| null` | User pending deletion |

The page fetches on mount and re-fetches after any successful create, update, or delete. A `refresh` callback is passed down to child components.

### UserForm

Rendered inside a shadcn/ui `Dialog`. Accepts:

| Prop | Type | Purpose |
|---|---|---|
| `mode` | `'create' \| 'edit'` | Controls title, submit action, and password field visibility |
| `initialData` | `User \| null` | Pre-populates fields in edit mode |
| `onSuccess` | `() => void` | Called after successful submit; triggers list refresh |
| `onClose` | `() => void` | Called on cancel or after success |

Uses `react-hook-form` with a `zod` schema for validation. Manages its own `submitting`, `submitError`, `roles`, and `rolesLoading` state. Calls `userService.getRoles()` on mount.

**Form fields:**

| Field | Component | Validation | Create | Edit |
|---|---|---|---|---|
| `login` | `Input` | Required, non-empty | ✓ | ✓ |
| `name` | `Input` | Required, non-empty | ✓ | ✓ |
| `email` | `Input` | Required, valid email pattern | ✓ | ✓ |
| `password` | `Input` (type=password) | Required, non-empty | ✓ | hidden |
| `roleId` | `Select` (Role_Select) | Required | ✓ | ✓ |
| `allDevicesAvailable` | `Checkbox` or `Switch` | Optional, defaults false | ✓ | ✓ |
| `customerStaff` | `Checkbox` or `Switch` | Optional, defaults false | ✓ | ✓ |

### DeleteDialog

Wraps `AlertDialog`. Receives `user: User | null`, `onConfirm`, `onCancel`. Manages `deleting` and `deleteError` state internally. Shows the user's login and name in the dialog body.

### RowActionsMenu

`DropdownMenu` with two items: "Edit" and "Delete". Rendered in the last column of each table row.

---

## Data Models

### Frontend Types (`src/features/users/types.ts`)

```typescript
export interface Role {
  id: number
  name: string
}

// Mirrors User from backend
export interface User {
  id: number
  login: string
  name: string
  email: string
  role: Role | null
  allDevicesAvailable: boolean
  customerStaff: boolean
}

// Body for POST and PUT
export interface UserPayload {
  login: string
  name: string
  email: string
  password?: string       // required in create mode, omitted in edit mode
  roleId: number
  allDevicesAvailable: boolean
  customerStaff: boolean
}
```

### API Response Wrapping

All responses use the `{ status: "OK" | "ERROR", data: T }` envelope. The `userService` calls `unwrapHmdmData` from `src/services/hmdmEnvelope.ts` to extract the inner payload, consistent with other services in the project.

### Zod Schema (inside UserForm)

```typescript
const baseSchema = z.object({
  login: z.string().min(1, 'Login is required'),
  name: z.string().min(1, 'Name is required'),
  email: z.string().regex(/[^@]+@[^.]+\..+/, 'Invalid email address'),
  roleId: z.number({ required_error: 'Role is required' }).int().positive('Role is required'),
  allDevicesAvailable: z.boolean().default(false),
  customerStaff: z.boolean().default(false),
})

const createSchema = baseSchema.extend({
  password: z.string().min(1, 'Password is required'),
})

const editSchema = baseSchema
```

---

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Table columns rendered for any user list

*For any* non-empty array of `User` objects returned by the API, the rendered table must contain exactly the columns: Login, Name, Email, Role, Status.

**Validates: Requirements 1.3**

### Property 2: Form submit calls correct endpoint with form values

*For any* valid `UserPayload` (non-empty login, name, email, valid email format, non-empty password in create mode, valid roleId), submitting the form in create mode must call `userService.createUser` with exactly those values; submitting in edit mode must call `userService.updateUser` with the user's `id` and exactly those values.

**Validates: Requirements 2.4, 3.3**

### Property 3: Cancel makes no API call

*For any* dialog state (form open in create or edit mode, or delete dialog open), clicking Cancel must not trigger any call to `userService.createUser`, `userService.updateUser`, or `userService.deleteUser`.

**Validates: Requirements 2.8, 4.6**

### Property 4: Edit mode pre-populates fields for any user

*For any* `User` object, opening the `UserForm` in edit mode must pre-populate the `login`, `name`, `email`, `roleId`, `allDevicesAvailable`, and `customerStaff` fields with that user's current values.

**Validates: Requirements 3.1**

### Property 5: Delete dialog shows user login and name for any user

*For any* `User` object, opening the `DeleteDialog` must render that user's `login` and `name` in the dialog body.

**Validates: Requirements 4.1**

### Property 6: Delete confirm calls DELETE with correct id

*For any* `User` object, when the user confirms deletion, `userService.deleteUser` must be called with exactly that user's numeric `id`.

**Validates: Requirements 4.2**

### Property 7: Role select populates and pre-selects correctly

*For any* array of `Role` objects returned by `getRoles`, the rendered `Role_Select` must contain one option per role with the role's `name` as the label and `id` as the value. In edit mode, for any `User`, the option matching `user.role.id` must be pre-selected.

**Validates: Requirements 5.3, 5.4**

### Property 8: Required field validation rejects invalid inputs

*For any* form submission where `login` or `name` is empty or whitespace-only, where `email` does not match the pattern `[^@]+@[^.]+\..+`, where `roleId` is absent, or (in create mode) where `password` is empty — the form must reject the submission and display a validation error; no API call must be made.

**Validates: Requirements 5.5, 6.1, 6.2, 6.3, 6.4**

### Property 9: Service routes to correct URL for any operation

*For any* user `id` and `UserPayload`, each service function must call the correct HTTP method and URL:
- `getUsers()` → `GET /rest/private/users`
- `createUser(data)` → `POST /rest/private/users`
- `updateUser(id, data)` → `PUT /rest/private/users/{id}`
- `deleteUser(id)` → `DELETE /rest/private/users/{id}`
- `getRoles()` → `GET /rest/private/roles`

**Validates: Requirements 8.1, 8.2, 8.3, 8.4, 8.5**

### Property 10: Service error propagation

*For any* `userService` function, when the underlying `apiClient` call rejects (non-2xx response or `status: "ERROR"` envelope), the service function must re-throw the error rather than swallowing it.

**Validates: Requirements 8.7**

### Property 11: Null-safe rendering for any user with null optional fields

*For any* `User` object where `role` is null, or where `allDevicesAvailable` or `customerStaff` are absent or null, rendering the table row and opening the form must not throw a runtime error. Null booleans must be treated as `false`; null role must render an empty cell.

**Validates: Requirements 9.4, 9.5**

---

## Error Handling

| Scenario | Behavior |
|---|---|
| List fetch fails | Show error banner with message and "Retry" button; table hidden |
| Create fails | Show error inside `UserForm` dialog; dialog stays open; submit button re-enabled |
| Update fails | Show error inside `UserForm` dialog; dialog stays open; submit button re-enabled |
| Delete fails | Show error inside `DeleteDialog`; dialog stays open; confirm button re-enabled |
| Roles fetch fails | Show error inside `UserForm`; Role_Select disabled with error message |
| 401 on any call | `apiClient` interceptor clears token and redirects to `/login` |
| Empty list | Show empty-state message: "No users found" |

All error messages use the error's `message` field if available, falling back to a generic string.

---

## Testing Strategy

### Dual Testing Approach

Both unit tests and property-based tests are required and complementary:
- Unit tests cover specific examples, integration points, and error states
- Property tests verify universal correctness across randomized inputs

### Unit Tests (Vitest + React Testing Library)

Focus areas:
- `UsersPage` shows skeleton while loading, error state on failure, empty state on empty list
- `UsersPage` renders "Add User" button
- `UserForm` in create mode has all fields empty and shows password field on open
- `UserForm` in edit mode pre-populates fields from `initialData` and hides password field
- `UserForm` disables submit button while submitting, shows error on failure, closes on success
- `UserForm` calls `getRoles` on mount; Role_Select is disabled while roles are loading
- `DeleteDialog` shows user login and name, disables confirm while deleting, shows error on failure, closes on success
- `userService.getUsers` calls `GET /rest/private/users`
- `userService` error propagation: mock a 500 response and assert the promise rejects
- Navigation: "Users" entry already present in `NAV_ITEMS` (verify, no change needed)

### Property-Based Tests (fast-check)

Each property test runs a minimum of 100 iterations. Each test is tagged with a comment referencing the design property.

```
// Feature: users-management, Property 1: Table columns rendered for any user list
// Feature: users-management, Property 2: Form submit calls correct endpoint with form values
// Feature: users-management, Property 3: Cancel makes no API call
// Feature: users-management, Property 4: Edit mode pre-populates fields for any user
// Feature: users-management, Property 5: Delete dialog shows user login and name for any user
// Feature: users-management, Property 6: Delete confirm calls DELETE with correct id
// Feature: users-management, Property 7: Role select populates and pre-selects correctly
// Feature: users-management, Property 8: Required field validation rejects invalid inputs
// Feature: users-management, Property 9: Service routes to correct URL for any operation
// Feature: users-management, Property 10: Service error propagation
// Feature: users-management, Property 11: Null-safe rendering for any user with null optional fields
```

**Generators needed**:
- `arbitraryRole()` — `fc.record({ id: fc.integer({ min: 1 }), name: fc.string({ minLength: 1 }) })`
- `arbitraryUser()` — `fc.record({ id: fc.integer({ min: 1 }), login: fc.string({ minLength: 1 }), name: fc.string({ minLength: 1 }), email: fc.emailAddress(), role: fc.oneof(arbitraryRole(), fc.constant(null)), allDevicesAvailable: fc.boolean(), customerStaff: fc.boolean() })`
- `arbitraryUserPayload()` — `fc.record({ login: fc.string({ minLength: 1 }), name: fc.string({ minLength: 1 }), email: fc.emailAddress(), password: fc.option(fc.string({ minLength: 1 })), roleId: fc.integer({ min: 1 }), allDevicesAvailable: fc.boolean(), customerStaff: fc.boolean() })`
- `arbitraryEmptyOrWhitespace()` — `fc.stringOf(fc.constantFrom(' ', '\t', '\n'))` for Property 8
- `arbitraryUserWithNulls()` — same as `arbitraryUser()` but with `role: null`, `allDevicesAvailable: null`, `customerStaff: null` for Property 11

### shadcn/ui Components Status

All required components are already scaffolded in `src/shared/ui/`:

| Component | Status |
|---|---|
| `table.tsx` | Already present |
| `dialog.tsx` | Needs scaffolding |
| `form.tsx` | Already present |
| `input.tsx` | Already present |
| `select.tsx` | Needs scaffolding |
| `button.tsx` | Already present |
| `alert-dialog.tsx` | Already present |
| `skeleton.tsx` | Already present |
| `dropdown-menu.tsx` | Already present |

Run `npx shadcn@latest add dialog select` inside `frontend/` to add the two missing components.
