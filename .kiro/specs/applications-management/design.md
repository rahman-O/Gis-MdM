# Design Document: Applications Management

## Overview

The Applications Management feature adds a fully functional `/applications` page to the HMDM frontend. It provides a table of all managed Android applications with create, edit, and delete operations. Create and edit use a `Dialog`-based form (react-hook-form + zod) that supports both manual entry and APK file upload. Delete uses an `AlertDialog` confirmation prompt.

All data flows through a dedicated `applicationService` module that wraps the shared `apiClient`. The UI is built exclusively from shadcn/ui components, following the same patterns established by the devices-management and configurations-management features.

### Key Backend API Findings

- **List**: `GET /rest/private/applications` — returns array of `Application` objects
- **Create**: `POST /rest/private/applications` — body is `ApplicationPayload`
- **Update**: `PUT /rest/private/applications/{id}` — body is `ApplicationPayload`
- **Delete**: `DELETE /rest/private/applications/{id}`
- **APK Upload**: `POST /rest/private/web-ui-files` — `multipart/form-data` with the APK file; returns `{ url: string }` (and optionally `name`, `pkg`, `version`)
- **Response shape**: all responses are wrapped in `{ status: "OK" | "ERROR", data: T }` — unwrapped via the existing `unwrapHmdmData` helper from `src/services/hmdmEnvelope.ts`

---

## Architecture

The feature follows the existing feature-slice pattern:

```
src/features/applications/
  ApplicationsPage.tsx     # Page component, route /applications
  applicationService.ts    # API calls via apiClient
  ApplicationForm.tsx      # Dialog form (create + edit, APK upload)
  types.ts                 # TypeScript interfaces
```

All shadcn/ui components required (`Table`, `Dialog`, `Form`, `Input`, `Switch`, `Button`, `AlertDialog`, `Skeleton`, `DropdownMenu`, `Badge`) are already scaffolded from the devices-management and configurations-management features.

### Data Flow

```mermaid
flowchart TD
    A[ApplicationsPage] -->|mount| B[applicationService.getApplications]
    B -->|GET /rest/private/applications| C[Backend]
    C -->|Application[]| B
    B --> A

    A -->|openCreate| D[ApplicationForm - create mode]
    D -->|file selected| E[applicationService.uploadApk]
    E -->|POST /rest/private/web-ui-files| C
    C -->|ApkUploadResponse| E
    E -->|url + metadata| D
    D -->|valid payload| F[applicationService.createApplication]
    F -->|POST /rest/private/applications| C

    A -->|openEdit app| G[ApplicationForm - edit mode]
    G -->|valid payload + id| H[applicationService.updateApplication]
    H -->|PUT /rest/private/applications/:id| C

    A -->|openDelete app| I[DeleteDialog]
    I -->|app.id| J[applicationService.deleteApplication]
    J -->|DELETE /rest/private/applications/:id| C

    F & H & J -->|success| A
    A -->|refresh| B
```

---

## Components and Interfaces

### ApplicationsPage

Top-level page component rendered at `/applications`. Owns all state:

| State | Type | Purpose |
|---|---|---|
| `applications` | `Application[]` | Full list from API |
| `loading` | `boolean` | List fetch in flight |
| `error` | `string \| null` | List fetch error |
| `formMode` | `'create' \| 'edit' \| null` | Controls form dialog visibility |
| `selectedApp` | `Application \| null` | Application being edited |
| `appToDelete` | `Application \| null` | Application pending deletion |

The page fetches on mount and re-fetches after any successful create, update, or delete.

### ApplicationForm

Rendered inside a shadcn/ui `Dialog`. Accepts:

| Prop | Type | Purpose |
|---|---|---|
| `mode` | `'create' \| 'edit'` | Controls title and submit action |
| `initialData` | `Application \| null` | Pre-populates fields in edit mode |
| `onSuccess` | `() => void` | Called after successful submit; triggers list refresh |
| `onClose` | `() => void` | Called on cancel or after success |

Uses `react-hook-form` with a `zod` schema for validation. Manages its own `submitting`, `submitError`, `uploading`, and `uploadError` state.

**Form fields:**

| Field | Component | Validation |
|---|---|---|
| `name` | `Input` | Required, non-empty |
| `pkg` | `Input` | Required, non-empty |
| `version` | `Input` | Required, non-empty |
| `url` | `Input` | Required, non-empty |
| APK file | `<input type="file" accept=".apk">` | Optional; `.apk` extension only |
| `system` | `Switch` | Optional, defaults to `false` |
| `showIcon` | `Switch` | Optional, defaults to `false` |
| `runAfterInstall` | `Switch` | Optional, defaults to `false` |
| `runAtBoot` | `Switch` | Optional, defaults to `false` |

When a file is selected, `applicationService.uploadApk` is called immediately. On success the `url` field is populated and `name`/`pkg`/`version` are auto-filled if the response includes them.

### DeleteDialog

Wraps `AlertDialog`. Receives `app: Application | null`, `onConfirm`, `onCancel`. Manages `deleting` and `deleteError` state internally. Shows the application name and package in the dialog body.

### SystemBadge

Inline component. Renders a `Badge` with variant `default` and label "System" when `system` is `true`; renders nothing (or a muted dash) when `false`.

### RowActionsMenu

`DropdownMenu` with two items: "Edit" and "Delete". Rendered in the last column of each table row.

---

## Data Models

### Frontend Types (`src/features/applications/types.ts`)

```typescript
// Mirrors Application from backend
export interface Application {
  id: number
  name: string
  pkg: string
  version: string
  url: string
  system: boolean
  showIcon: boolean
  runAfterInstall: boolean
  runAtBoot: boolean
}

// Body for POST and PUT
export interface ApplicationPayload {
  name: string
  pkg: string
  version: string
  url: string
  system: boolean
  showIcon: boolean
  runAfterInstall: boolean
  runAtBoot: boolean
}

// Response from POST /rest/private/web-ui-files
export interface ApkUploadResponse {
  url: string
  name?: string
  pkg?: string
  version?: string
}
```

### API Response Wrapping

All responses use the `{ status: "OK" | "ERROR", data: T }` envelope. The `applicationService` calls `unwrapHmdmData` from `src/services/hmdmEnvelope.ts` to extract the inner payload, consistent with other services in the project.

### Zod Schema (inside ApplicationForm)

```typescript
const applicationSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  pkg: z.string().min(1, 'Package name is required'),
  version: z.string().min(1, 'Version is required'),
  url: z.string().min(1, 'URL is required'),
  system: z.boolean().default(false),
  showIcon: z.boolean().default(false),
  runAfterInstall: z.boolean().default(false),
  runAtBoot: z.boolean().default(false),
})
```

---

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Table columns rendered for any application list

*For any* non-empty array of `Application` objects returned by the API, the rendered table must contain exactly the columns: Name, Package, Version, and System.

**Validates: Requirements 1.3**

### Property 2: System flag indicator rendered correctly for any application

*For any* `Application` object, the System column must render a positive indicator (badge or icon) when `system` is `true` and a negative or empty indicator when `system` is `false`.

**Validates: Requirements 1.4**

### Property 3: Form submit calls correct endpoint based on mode

*For any* valid `ApplicationPayload` (non-empty name, pkg, version, url), submitting the form in create mode must call `applicationService.createApplication` with exactly those values; submitting in edit mode must call `applicationService.updateApplication` with the application's `id` and exactly those values.

**Validates: Requirements 2.6, 3.2**

### Property 4: Cancel makes no API call

*For any* dialog state (form open in create or edit mode, or delete dialog open), clicking Cancel must not trigger any call to `applicationService.createApplication`, `applicationService.updateApplication`, or `applicationService.deleteApplication`.

**Validates: Requirements 2.10, 4.6**

### Property 5: Edit mode pre-populates all fields for any application

*For any* `Application` object, opening the `ApplicationForm` in edit mode must pre-populate the `name`, `pkg`, `version`, `url`, `system`, `showIcon`, `runAfterInstall`, and `runAtBoot` fields with that application's current values.

**Validates: Requirements 3.1**

### Property 6: Delete dialog shows application name and package for any application

*For any* `Application` object, opening the `DeleteDialog` must render that application's `name` and `pkg` in the dialog body.

**Validates: Requirements 4.1**

### Property 7: Delete confirm calls DELETE with correct id

*For any* `Application` object, when the user confirms deletion, `applicationService.deleteApplication` must be called with exactly that application's numeric `id`.

**Validates: Requirements 4.2**

### Property 8: Required field validation rejects empty or whitespace inputs

*For any* form submission where `name`, `pkg`, `version`, or `url` is empty or composed entirely of whitespace, the form must reject the submission and display a validation error — no API call must be made.

**Validates: Requirements 5.1, 5.2, 5.3, 5.4**

### Property 9: Validation errors displayed adjacent to invalid fields

*For any* invalid field submission, the rendered form must display a `FormMessage` error adjacent to each invalid field.

**Validates: Requirements 5.6**

### Property 10: Non-APK file rejected before upload starts

*For any* file whose name does not end in `.apk`, selecting it in the file input must display a validation error and must not trigger a call to `applicationService.uploadApk`.

**Validates: Requirements 6.2**

### Property 11: APK upload sets url field to returned URL

*For any* `ApkUploadResponse`, after a successful APK upload the `url` form field must be set to exactly the `url` value from the response.

**Validates: Requirements 2.5, 6.5**

### Property 12: Service routes to correct URL for any operation

*For any* application `id` and `ApplicationPayload`, each service function must call the correct HTTP method and URL:
- `getApplications()` → `GET /rest/private/applications`
- `createApplication(data)` → `POST /rest/private/applications`
- `updateApplication(id, data)` → `PUT /rest/private/applications/{id}`
- `deleteApplication(id)` → `DELETE /rest/private/applications/{id}`
- `uploadApk(file)` → `POST /rest/private/web-ui-files` with `Content-Type: multipart/form-data`

**Validates: Requirements 8.1, 8.2, 8.3, 8.4, 8.5**

### Property 13: Service error propagation

*For any* `applicationService` function, when the underlying `apiClient` call rejects (non-2xx response or `status: "ERROR"` envelope), the service function must re-throw the error rather than swallowing it.

**Validates: Requirements 8.7**

### Property 14: Null-safe boolean field handling

*For any* `Application` object where boolean fields (`system`, `showIcon`, `runAfterInstall`, `runAtBoot`) are absent or null in the API response, rendering the table row and form must not throw a runtime error, and those fields must be treated as `false`.

**Validates: Requirements 9.4**

---

## Error Handling

| Scenario | Behavior |
|---|---|
| List fetch fails | Show error banner with message and "Retry" button; table hidden |
| APK upload fails | Show error on file input; allow user to retry by selecting file again |
| Create fails | Show error inside `ApplicationForm` dialog; dialog stays open; submit button re-enabled |
| Update fails | Show error inside `ApplicationForm` dialog; dialog stays open; submit button re-enabled |
| Delete fails | Show error inside `DeleteDialog`; dialog stays open; confirm button re-enabled |
| 401 on any call | `apiClient` interceptor clears token and redirects to `/login` |
| Empty list | Show empty-state message: "No applications found" |

All error messages use the error's `message` field if available, falling back to a generic string.

---

## Testing Strategy

### Dual Testing Approach

Both unit tests and property-based tests are required and complementary:
- Unit tests cover specific examples, integration points, and error states
- Property tests verify universal correctness across randomized inputs

### Unit Tests (Vitest + React Testing Library)

Focus areas:
- `ApplicationsPage` shows skeleton while loading, error state on failure, empty state on empty list
- `ApplicationsPage` renders "Add Application" button
- `ApplicationForm` in create mode has all fields empty and switches defaulting to false on open
- `ApplicationForm` in edit mode pre-populates all fields from `initialData`
- `ApplicationForm` disables submit button while submitting, shows error on failure, closes on success
- `ApplicationForm` file input accepts only `.apk` files
- `DeleteDialog` shows application name and package, disables confirm while deleting, shows error on failure, closes on success
- `applicationService.getApplications` calls `GET /rest/private/applications`
- `applicationService.uploadApk` sends `multipart/form-data` to `POST /rest/private/web-ui-files`
- `applicationService` error propagation: mock a 500 response and assert the promise rejects
- Navigation: "Applications" entry appears in `NAV_ITEMS` with the `Package` icon

### Property-Based Tests (fast-check)

Each property test runs a minimum of 100 iterations. Each test is tagged with a comment referencing the design property.

```
// Feature: applications-management, Property 1: Table columns rendered for any application list
// Feature: applications-management, Property 2: System flag indicator rendered correctly for any application
// Feature: applications-management, Property 3: Form submit calls correct endpoint based on mode
// Feature: applications-management, Property 4: Cancel makes no API call
// Feature: applications-management, Property 5: Edit mode pre-populates all fields for any application
// Feature: applications-management, Property 6: Delete dialog shows application name and package for any application
// Feature: applications-management, Property 7: Delete confirm calls DELETE with correct id
// Feature: applications-management, Property 8: Required field validation rejects empty or whitespace inputs
// Feature: applications-management, Property 9: Validation errors displayed adjacent to invalid fields
// Feature: applications-management, Property 10: Non-APK file rejected before upload starts
// Feature: applications-management, Property 11: APK upload sets url field to returned URL
// Feature: applications-management, Property 12: Service routes to correct URL for any operation
// Feature: applications-management, Property 13: Service error propagation
// Feature: applications-management, Property 14: Null-safe boolean field handling
```

**Generators needed**:
- `arbitraryApplication()` — `fc.record({ id: fc.integer({ min: 1 }), name: fc.string({ minLength: 1 }), pkg: fc.string({ minLength: 1 }), version: fc.string({ minLength: 1 }), url: fc.string({ minLength: 1 }), system: fc.boolean(), showIcon: fc.boolean(), runAfterInstall: fc.boolean(), runAtBoot: fc.boolean() })`
- `arbitraryApplicationPayload()` — same as above without `id`
- `arbitraryEmptyOrWhitespace()` — `fc.stringOf(fc.constantFrom(' ', '\t', '\n'))` for Property 8
- `arbitraryNonApkFile()` — `fc.string({ minLength: 1 }).filter(s => !s.endsWith('.apk'))` for Property 10
- `arbitraryApkUploadResponse()` — `fc.record({ url: fc.string({ minLength: 1 }), name: fc.option(fc.string()), pkg: fc.option(fc.string()), version: fc.option(fc.string()) })` for Property 11
- `arbitraryApplicationWithNullBooleans()` — same as `arbitraryApplication()` but with boolean fields forced to `null` for Property 14

### shadcn/ui Components Required

All required components are already present from devices-management and configurations-management:

| Component | Status |
|---|---|
| `table.tsx` | Already scaffolded |
| `dialog.tsx` | Already scaffolded |
| `alert-dialog.tsx` | Already scaffolded |
| `dropdown-menu.tsx` | Already scaffolded |
| `skeleton.tsx` | Already scaffolded |
| `badge.tsx` | Already scaffolded |
| `button.tsx` | Already scaffolded |
| `input.tsx` | Already scaffolded |
| `form.tsx` | Already scaffolded |
| `switch.tsx` | Needs scaffolding: `npx shadcn@latest add switch` |
