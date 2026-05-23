# Design Document: Configurations Management

## Overview

The Configurations Management feature adds a fully functional `/configurations` page to the HMDM frontend. It provides a table of all MDM policy configurations with create, edit, and delete operations. Create and edit use a `Dialog`-based form (react-hook-form + zod). Delete uses an `AlertDialog` confirmation prompt. A Device Count column shows how many devices are assigned to each configuration.

All data flows through a dedicated `configurationService` module that wraps the shared `apiClient`. The UI is built exclusively from shadcn/ui components, following the same patterns established by the devices-management feature.

### Key Backend API Findings

- **List**: `GET /rest/private/configurations` — returns array of `Configuration` objects
- **Create**: `POST /rest/private/configurations` — body is `ConfigurationPayload`
- **Update**: `PUT /rest/private/configurations/{id}` — body is `ConfigurationPayload`
- **Delete**: `DELETE /rest/private/configurations/{id}`
- **Response shape**: all responses are wrapped in `{ status: "OK" | "ERROR", data: T }` — unwrapped via the existing `unwrapHmdmData` helper from `src/services/hmdmEnvelope.ts`
- **Device count**: derived from a `deviceCount` field on the `Configuration` object (or falls back to `0` if absent)

---

## Architecture

The feature follows the existing feature-slice pattern:

```
src/features/configurations/
  ConfigurationsPage.tsx     # Page component, route /configurations
  configurationService.ts    # API calls via apiClient
  ConfigurationForm.tsx      # Dialog form (create + edit)
  types.ts                   # TypeScript interfaces

src/shared/ui/               # New shadcn/ui components to scaffold
  alert-dialog.tsx           # Reused from devices-management pattern
  dropdown-menu.tsx          # Reused from devices-management pattern
  select.tsx                 # New — needed for Type_Selector
  dialog.tsx                 # New — needed for Configuration_Form
```

### Data Flow

```mermaid
flowchart TD
    A[ConfigurationsPage] -->|mount| B[configurationService.getConfigurations]
    B -->|GET /rest/private/configurations| C[Backend]
    C -->|Configuration[]| B
    B --> A

    A -->|openCreate| D[ConfigurationForm - create mode]
    D -->|valid payload| E[configurationService.createConfiguration]
    E -->|POST /rest/private/configurations| C

    A -->|openEdit config| F[ConfigurationForm - edit mode]
    F -->|valid payload + id| G[configurationService.updateConfiguration]
    G -->|PUT /rest/private/configurations/:id| C

    A -->|openDelete config| H[DeleteDialog]
    H -->|config.id| I[configurationService.deleteConfiguration]
    I -->|DELETE /rest/private/configurations/:id| C

    E & G & I -->|success| A
    A -->|refresh| B
```

---

## Components and Interfaces

### ConfigurationsPage

Top-level page component rendered at `/configurations`. Owns all state:

| State | Type | Purpose |
|---|---|---|
| `configurations` | `Configuration[]` | Full list from API |
| `loading` | `boolean` | List fetch in flight |
| `error` | `string \| null` | List fetch error |
| `formMode` | `'create' \| 'edit' \| null` | Controls form dialog visibility |
| `selectedConfig` | `Configuration \| null` | Config being edited |
| `configToDelete` | `Configuration \| null` | Config pending deletion |

The page fetches on mount and re-fetches after any successful create, update, or delete. A `refresh` callback is passed down to child components.

### ConfigurationForm

Rendered inside a shadcn/ui `Dialog`. Accepts:

| Prop | Type | Purpose |
|---|---|---|
| `mode` | `'create' \| 'edit'` | Controls title and submit action |
| `initialData` | `Configuration \| null` | Pre-populates fields in edit mode |
| `onSuccess` | `() => void` | Called after successful submit; triggers list refresh |
| `onClose` | `() => void` | Called on cancel or after success |

Uses `react-hook-form` with a `zod` schema for validation. Manages its own `submitting` and `submitError` state.

**Form fields:**

| Field | Component | Validation |
|---|---|---|
| `name` | `Input` | Required, non-empty |
| `type` | `Select` (COMMON / WORK) | Required |
| `description` | `Input` (or `Textarea`) | Optional |

### DeleteDialog

Wraps `AlertDialog`. Receives `config: Configuration | null`, `onConfirm`, `onCancel`. Manages `deleting` and `deleteError` state internally. Shows the configuration name in the dialog body.

### RowActionsMenu

`DropdownMenu` with two items: "Edit" and "Delete". Rendered in the last column of each table row.

---

## Data Models

### Frontend Types (`src/features/configurations/types.ts`)

```typescript
export interface ConfigurationFile {
  id: number
}

export interface ConfigurationApplication {
  id: number
}

// Mirrors Configuration from backend
export interface Configuration {
  id: number
  name: string
  description: string | null
  type: 'COMMON' | 'WORK'
  applicationId: number | null
  deviceCount?: number | null      // count of assigned devices
  files: ConfigurationFile[]
  applications: ConfigurationApplication[]
}

// Body for POST and PUT
export interface ConfigurationPayload {
  name: string
  description?: string
  type: 'COMMON' | 'WORK'
}
```

### API Response Wrapping

All responses use the `{ status: "OK" | "ERROR", data: T }` envelope. The `configurationService` calls `unwrapHmdmData` from `src/services/hmdmEnvelope.ts` to extract the inner payload, consistent with other services in the project.

### Zod Schema (inside ConfigurationForm)

```typescript
const configurationSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  type: z.enum(['COMMON', 'WORK'], { required_error: 'Type is required' }),
  description: z.string().optional(),
})
```

---

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Table columns rendered for any configuration list

*For any* non-empty array of `Configuration` objects returned by the API, the rendered table must contain exactly the columns: Name, Type, Description, and Device Count.

**Validates: Requirements 1.3, 6.1**

### Property 2: Form submit calls correct endpoint with form values

*For any* valid `ConfigurationPayload` (non-empty name, valid type), submitting the form in create mode must call `configurationService.createConfiguration` with exactly those values; submitting in edit mode must call `configurationService.updateConfiguration` with the configuration's `id` and exactly those values.

**Validates: Requirements 2.3, 3.2**

### Property 3: Cancel makes no API call

*For any* dialog state (form open in create or edit mode, or delete dialog open), clicking Cancel must not trigger any call to `configurationService.createConfiguration`, `configurationService.updateConfiguration`, or `configurationService.deleteConfiguration`.

**Validates: Requirements 2.7, 4.6**

### Property 4: Edit mode pre-populates fields for any configuration

*For any* `Configuration` object, opening the `ConfigurationForm` in edit mode must pre-populate the `name`, `description`, and `type` fields with that configuration's current values.

**Validates: Requirements 3.1**

### Property 5: Delete dialog shows configuration name for any configuration

*For any* `Configuration` object, opening the `DeleteDialog` must render that configuration's `name` in the dialog body.

**Validates: Requirements 4.1**

### Property 6: Delete confirm calls DELETE with correct id

*For any* `Configuration` object, when the user confirms deletion, `configurationService.deleteConfiguration` must be called with exactly that configuration's numeric `id`.

**Validates: Requirements 4.2**

### Property 7: Required field validation rejects empty or whitespace inputs

*For any* form submission where `name` is empty or composed entirely of whitespace, or where `type` is not one of `'COMMON'` or `'WORK'`, the form must reject the submission and display a validation error — no API call must be made.

**Validates: Requirements 5.1, 5.2**

### Property 8: Optional description allows submission

*For any* valid `name` and `type`, submitting the form with an empty or absent `description` must succeed (not be rejected by validation).

**Validates: Requirements 5.3**

### Property 9: Device count displayed correctly for any configuration

*For any* `Configuration` object with a `deviceCount` value (including `0`), the rendered table row must display that exact count in the Device Count column.

**Validates: Requirements 6.2, 6.3**

### Property 10: Service routes to correct URL for any operation

*For any* configuration `id` and `ConfigurationPayload`, each service function must call the correct HTTP method and URL:
- `getConfigurations()` → `GET /rest/private/configurations`
- `getConfiguration(id)` → `GET /rest/private/configurations/{id}`
- `createConfiguration(data)` → `POST /rest/private/configurations`
- `updateConfiguration(id, data)` → `PUT /rest/private/configurations/{id}`
- `deleteConfiguration(id)` → `DELETE /rest/private/configurations/{id}`

**Validates: Requirements 8.1, 8.2, 8.3, 8.4, 8.5**

### Property 11: Service error propagation

*For any* `configurationService` function, when the underlying `apiClient` call rejects (non-2xx response or `status: "ERROR"` envelope), the service function must re-throw the error rather than swallowing it.

**Validates: Requirements 8.7**

### Property 12: Null-safe rendering for any configuration with null optional fields

*For any* `Configuration` object where `description`, `deviceCount`, `applicationId`, `files`, or `applications` are absent or null, rendering the table row and form must not throw a runtime error.

**Validates: Requirements 9.5**

---

## Error Handling

| Scenario | Behavior |
|---|---|
| List fetch fails | Show error banner with message and "Retry" button; table hidden |
| Create fails | Show error inside `ConfigurationForm` dialog; dialog stays open; submit button re-enabled |
| Update fails | Show error inside `ConfigurationForm` dialog; dialog stays open; submit button re-enabled |
| Delete fails | Show error inside `DeleteDialog`; dialog stays open; confirm button re-enabled |
| 401 on any call | `apiClient` interceptor clears token and redirects to `/login` |
| Empty list | Show empty-state message: "No configurations found" |

All error messages use the error's `message` field if available, falling back to a generic string.

---

## Testing Strategy

### Dual Testing Approach

Both unit tests and property-based tests are required and complementary:
- Unit tests cover specific examples, integration points, and error states
- Property tests verify universal correctness across randomized inputs

### Unit Tests (Vitest + React Testing Library)

Focus areas:
- `ConfigurationsPage` shows skeleton while loading, error state on failure, empty state on empty list
- `ConfigurationsPage` renders "New Configuration" button
- `ConfigurationForm` in create mode has all fields empty on open
- `ConfigurationForm` in edit mode pre-populates fields from `initialData`
- `ConfigurationForm` disables submit button while submitting, shows error on failure, closes on success
- `DeleteDialog` shows configuration name, disables confirm while deleting, shows error on failure, closes on success
- `Type_Selector` renders exactly two options: COMMON and WORK
- `configurationService.getConfigurations` calls `GET /rest/private/configurations`
- `configurationService` error propagation: mock a 500 response and assert the promise rejects
- Navigation: "Configurations" entry appears after "Devices" in `NAV_ITEMS`

### Property-Based Tests (fast-check)

Each property test runs a minimum of 100 iterations. Each test is tagged with a comment referencing the design property.

```
// Feature: configurations-management, Property 1: Table columns rendered for any configuration list
// Feature: configurations-management, Property 2: Form submit calls correct endpoint with form values
// Feature: configurations-management, Property 3: Cancel makes no API call
// Feature: configurations-management, Property 4: Edit mode pre-populates fields for any configuration
// Feature: configurations-management, Property 5: Delete dialog shows configuration name for any configuration
// Feature: configurations-management, Property 6: Delete confirm calls DELETE with correct id
// Feature: configurations-management, Property 7: Required field validation rejects empty or whitespace inputs
// Feature: configurations-management, Property 8: Optional description allows submission
// Feature: configurations-management, Property 9: Device count displayed correctly for any configuration
// Feature: configurations-management, Property 10: Service routes to correct URL for any operation
// Feature: configurations-management, Property 11: Service error propagation
// Feature: configurations-management, Property 12: Null-safe rendering for any configuration with null optional fields
```

**Generators needed**:
- `arbitraryConfiguration()` — `fc.record({ id: fc.integer({ min: 1 }), name: fc.string({ minLength: 1 }), description: fc.option(fc.string()), type: fc.oneof(fc.constant('COMMON'), fc.constant('WORK')), deviceCount: fc.option(fc.integer({ min: 0 })), applicationId: fc.option(fc.integer({ min: 1 })), files: fc.array(fc.record({ id: fc.integer() })), applications: fc.array(fc.record({ id: fc.integer() })) })`
- `arbitraryConfigurationPayload()` — `fc.record({ name: fc.string({ minLength: 1 }), type: fc.oneof(fc.constant('COMMON'), fc.constant('WORK')), description: fc.option(fc.string()) })`
- `arbitraryEmptyOrWhitespace()` — `fc.stringOf(fc.constantFrom(' ', '\t', '\n'))` for Property 7
- `fc.integer({ min: 1 })` for id-based properties (10, 6)
- `arbitraryConfigurationWithNulls()` — same as `arbitraryConfiguration()` but with all optional fields forced to `null` or `[]` for Property 12

### shadcn/ui Components to Scaffold

The following components must be added to `src/shared/ui/` before implementation:

| Component | Source | Status |
|---|---|---|
| `alert-dialog.tsx` | `npx shadcn@latest add alert-dialog` | New |
| `dropdown-menu.tsx` | `npx shadcn@latest add dropdown-menu` | New |
| `select.tsx` | `npx shadcn@latest add select` | New |
| `dialog.tsx` | `npx shadcn@latest add dialog` | New |

Components already present: `table.tsx`, `skeleton.tsx`, `badge.tsx`, `button.tsx`, `input.tsx`, `form.tsx`, `label.tsx`.
