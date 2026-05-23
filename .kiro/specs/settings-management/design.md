# Design Document: Settings Management

## Overview

The Settings Management feature replaces the placeholder `SettingsPage` at `/settings` with a fully functional single-page form for viewing and updating global HMDM instance settings. The page fetches current settings on mount via `GET /rest/private/settings`, renders all fields in a single `react-hook-form` + `zod` form, and persists changes via `PUT /rest/private/settings`. Success and failure outcomes are communicated through shadcn/ui Toast notifications.

The `newDeviceConfigurationId` field is backed by a `Select` populated from `GET /rest/private/configurations`, fetched in parallel with the settings load. All data flows through a dedicated `settingsService` module that wraps the shared `apiClient`.

### Key Backend API Findings

- **Get settings**: `GET /rest/private/settings` — returns a single `Settings` object
- **Update settings**: `PUT /rest/private/settings` — body is the full settings payload; returns updated `Settings`
- **Configurations list**: `GET /rest/private/configurations` — returns `Configuration[]` used to populate the `newDeviceConfigurationId` select
- **Response shape**: all responses are wrapped in `{ status: "OK" | "ERROR", data: T }` — unwrapped via the existing `unwrapHmdmData` helper from `src/services/hmdmEnvelope.ts`
- **Navigation**: Settings route and nav entry already exist (`/settings` in `App.tsx`, `Settings` in `navItems.ts`)

---

## Architecture

The feature follows the existing feature-slice pattern:

```
src/features/settings/
  SettingsPage.tsx       # Page component, route /settings (replaces placeholder)
  settingsService.ts     # API calls via apiClient
  types.ts               # TypeScript interfaces

src/shared/ui/           # New shadcn/ui components to scaffold
  select.tsx             # npx shadcn@latest add select
  toast.tsx + toaster.tsx + use-toast.ts  # npx shadcn@latest add toast
```

Components already present: `Form`, `Input`, `Switch`, `Button`, `Skeleton`, `Label`.

### Data Flow

```mermaid
flowchart TD
    A[SettingsPage] -->|mount - parallel| B[settingsService.getSettings]
    A -->|mount - parallel| C[configurationService.getConfigurations]
    B -->|GET /rest/private/settings| D[Backend]
    C -->|GET /rest/private/configurations| D
    D -->|Settings| B
    D -->|Configuration[]| C
    B & C --> A

    A -->|form populated| E[SettingsForm]
    E -->|valid submit| F[settingsService.updateSettings]
    F -->|PUT /rest/private/settings| D
    D -->|Settings| F
    F -->|success| G[success Toast]
    F -->|failure| H[error Toast]
```

---

## Components and Interfaces

### SettingsPage

Top-level page component rendered at `/settings`. Owns all state:

| State | Type | Purpose |
|---|---|---|
| `settings` | `Settings \| null` | Current settings from API |
| `configurations` | `ConfigurationOption[]` | Options for newDeviceConfigurationId select |
| `loading` | `boolean` | Initial data fetch in flight |
| `error` | `string \| null` | Initial fetch error |
| `submitting` | `boolean` | PUT request in flight |

The page fetches settings and configurations in parallel on mount using `Promise.all`. On error it shows an error banner with a "Retry" button. While loading it renders `Skeleton` placeholders. Once data is available it renders the form pre-populated with the fetched values.

The `useToast` hook is called at the page level; `toast()` is invoked after the PUT resolves or rejects.

### SettingsForm (inline in SettingsPage)

The form is rendered directly inside `SettingsPage` (not a separate dialog). It uses `react-hook-form` with a `zod` resolver.

**Form fields:**

| Field | Component | Type | Validation |
|---|---|---|---|
| `customerName` | `Input` | text | Required, non-empty |
| `createNewDevices` | `Switch` | boolean | Optional, defaults false |
| `newDeviceConfigurationId` | `Select` | number \| null | Optional |
| `language` | `Select` | string | Optional |
| `passwordLength` | `Input` | number | Required, min 1 |
| `passwordStrength` | `Select` | number | Required |
| `sendDeviceInfoExpiryDays` | `Input` | number | Required, min 1 |
| `unsecureEnrollment` | `Switch` | boolean | Optional, defaults false |
| `deviceFastSearch` | `Switch` | boolean | Optional, defaults false |

**Language options** (fixed list matching HMDM backend supported locales):

| Value | Label |
|---|---|
| `en` | English |
| `ru` | Russian |
| `de` | German |
| `fr` | French |
| `es` | Spanish |
| `pt` | Portuguese |
| `zh` | Chinese |

**Password strength options** (fixed, per requirements):

| Value | Label |
|---|---|
| `0` | Any |
| `1` | Numeric |
| `2` | Alphabetic |
| `3` | Alphanumeric |

The submit button is disabled and shows a spinner while `submitting` is true. On successful PUT the form values are reset to the returned settings. On failure the form retains the submitted values.

---

## Data Models

### Frontend Types (`src/features/settings/types.ts`)

```typescript
// Mirrors Settings from backend
export interface Settings {
  id: number
  customerName: string
  createNewDevices: boolean
  newDeviceConfigurationId: number | null
  language: string
  passwordLength: number
  passwordStrength: number
  sendDeviceInfoExpiryDays: number
  unsecureEnrollment: boolean
  deviceFastSearch: boolean
}

// Body for PUT /rest/private/settings (id excluded)
export interface SettingsPayload {
  customerName: string
  createNewDevices: boolean
  newDeviceConfigurationId: number | null
  language: string
  passwordLength: number
  passwordStrength: number
  sendDeviceInfoExpiryDays: number
  unsecureEnrollment: boolean
  deviceFastSearch: boolean
}

// Lightweight option type for the configurations select
export interface ConfigurationOption {
  id: number
  name: string
}
```

### Zod Schema (inside SettingsPage)

```typescript
const settingsSchema = z.object({
  customerName: z.string().min(1, 'Customer name is required'),
  createNewDevices: z.boolean().default(false),
  newDeviceConfigurationId: z.number().nullable().optional(),
  language: z.string().min(1, 'Language is required'),
  passwordLength: z.number({ invalid_type_error: 'Password length is required' }).int().min(1, 'Must be at least 1'),
  passwordStrength: z.number({ invalid_type_error: 'Password strength is required' }),
  sendDeviceInfoExpiryDays: z.number({ invalid_type_error: 'Expiry days is required' }).int().min(1, 'Must be at least 1'),
  unsecureEnrollment: z.boolean().default(false),
  deviceFastSearch: z.boolean().default(false),
})
```

### API Response Wrapping

All responses use the `{ status: "OK" | "ERROR", data: T }` envelope. The `settingsService` calls `unwrapHmdmData` from `src/services/hmdmEnvelope.ts` to extract the inner payload, consistent with other services in the project.

### Null Coalescion on Load

When populating the form from the API response, boolean fields absent or null are coerced to `false`; numeric fields absent or null are coerced to `0`:

```typescript
form.reset({
  customerName: settings.customerName ?? '',
  createNewDevices: settings.createNewDevices ?? false,
  newDeviceConfigurationId: settings.newDeviceConfigurationId ?? null,
  language: settings.language ?? 'en',
  passwordLength: settings.passwordLength ?? 0,
  passwordStrength: settings.passwordStrength ?? 0,
  sendDeviceInfoExpiryDays: settings.sendDeviceInfoExpiryDays ?? 0,
  unsecureEnrollment: settings.unsecureEnrollment ?? false,
  deviceFastSearch: settings.deviceFastSearch ?? false,
})
```

---

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Form populated with any settings object

*For any* `Settings` object returned by the API, after the GET resolves, every form field must reflect the corresponding value from that settings object (with null booleans treated as `false` and null numerics treated as `0`).

**Validates: Requirements 1.4, 8.3, 8.4**

### Property 2: Configurations populate select options

*For any* array of `Configuration` objects returned by `GET /rest/private/configurations`, the `newDeviceConfigurationId` select must render exactly one option per configuration, with the configuration's `name` as the label and `id` as the value.

**Validates: Requirements 2.3**

### Property 3: Submit calls service with form values for any valid payload

*For any* valid `SettingsPayload` (non-empty `customerName`, `passwordLength >= 1`, `sendDeviceInfoExpiryDays >= 1`), submitting the form must call `settingsService.updateSettings` with exactly those values.

**Validates: Requirements 3.1, 7.2**

### Property 4: Save button re-enabled after any PUT outcome

*For any* PUT request outcome (success or failure), the "Save Settings" button must be re-enabled after the request completes — it must never remain permanently disabled.

**Validates: Requirements 3.5**

### Property 5: Validation rejects invalid required fields

*For any* form submission where `customerName` is empty or whitespace-only, or where `passwordLength` is less than 1 or absent, or where `sendDeviceInfoExpiryDays` is less than 1 or absent — the form must reject the submission, display a validation error adjacent to the invalid field, and make no call to `settingsService.updateSettings`.

**Validates: Requirements 4.1, 4.2, 4.3, 4.5**

### Property 6: Service routes to correct URLs

*For any* call to `settingsService.getSettings`, the underlying HTTP call must be `GET /rest/private/settings`. For any `SettingsPayload` passed to `settingsService.updateSettings`, the underlying HTTP call must be `PUT /rest/private/settings` with that payload as the request body.

**Validates: Requirements 7.1, 7.2**

### Property 7: Service error propagation

*For any* `settingsService` function, when the underlying `apiClient` call rejects (non-2xx response or `status: "ERROR"` envelope), the service function must re-throw the error rather than swallowing it.

**Validates: Requirements 7.4**

### Property 8: Null-safe rendering for any settings with null fields

*For any* `Settings` object where boolean fields (`createNewDevices`, `unsecureEnrollment`, `deviceFastSearch`) are absent or null, or where numeric fields (`passwordLength`, `passwordStrength`, `sendDeviceInfoExpiryDays`) are absent or null, rendering the `SettingsPage` must not throw a runtime error. Null booleans must render as `false`; null numerics must render as `0`.

**Validates: Requirements 8.3, 8.4**

---

## Error Handling

| Scenario | Behavior |
|---|---|
| GET settings fails | Show error banner with message and "Retry" button; form hidden |
| GET configurations fails | `newDeviceConfigurationId` select shows empty options; form still renders |
| PUT fails | Error Toast with title "Failed to save settings" and destructive style; form retains submitted values; button re-enabled |
| PUT succeeds | Success Toast with title "Settings saved"; form reset to returned values |
| 401 on any call | `apiClient` interceptor clears token and redirects to `/login` |
| Null/absent fields in GET response | Coerced to safe defaults (`false` for booleans, `0` for numerics) |

All error messages use the error's `message` field if available, falling back to a generic string.

---

## Testing Strategy

### Dual Testing Approach

Both unit tests and property-based tests are required and complementary:
- Unit tests cover specific examples, integration points, and error states
- Property tests verify universal correctness across randomized inputs

### Unit Tests (Vitest + React Testing Library)

Focus areas:
- `SettingsPage` shows skeleton while loading, error state with retry on GET failure
- `SettingsPage` renders all 10 form fields once data loads
- `SettingsPage` calls `settingsService.getSettings` on mount
- `SettingsPage` disables submit button while submitting, re-enables after completion
- `SettingsPage` shows success toast on PUT success with title "Settings saved"
- `SettingsPage` shows error toast on PUT failure with title "Failed to save settings"
- `SettingsPage` validation: empty `customerName` shows error and prevents submit
- `SettingsPage` validation: `passwordLength` of 0 shows error and prevents submit
- `settingsService.getSettings` calls `GET /rest/private/settings`
- `settingsService.updateSettings` calls `PUT /rest/private/settings` with payload
- `settingsService` error propagation: mock a 500 response and assert the promise rejects
- Navigation: "Settings" entry present in `NAV_ITEMS` (already verified — no change needed)

### Property-Based Tests (fast-check)

Each property test runs a minimum of 100 iterations. Each test is tagged with a comment referencing the design property.

```
// Feature: settings-management, Property 1: Form populated with any settings object
// Feature: settings-management, Property 2: Configurations populate select options
// Feature: settings-management, Property 3: Submit calls service with form values for any valid payload
// Feature: settings-management, Property 4: Save button re-enabled after any PUT outcome
// Feature: settings-management, Property 5: Validation rejects invalid required fields
// Feature: settings-management, Property 6: Service routes to correct URLs
// Feature: settings-management, Property 7: Service error propagation
// Feature: settings-management, Property 8: Null-safe rendering for any settings with null fields
```

**Generators needed**:
- `arbitrarySettings()` — `fc.record({ id: fc.integer({ min: 1 }), customerName: fc.string({ minLength: 1 }), createNewDevices: fc.boolean(), newDeviceConfigurationId: fc.oneof(fc.integer({ min: 1 }), fc.constant(null)), language: fc.constantFrom('en', 'ru', 'de', 'fr', 'es', 'pt', 'zh'), passwordLength: fc.integer({ min: 1 }), passwordStrength: fc.integer({ min: 0, max: 3 }), sendDeviceInfoExpiryDays: fc.integer({ min: 1 }), unsecureEnrollment: fc.boolean(), deviceFastSearch: fc.boolean() })`
- `arbitrarySettingsPayload()` — same as `arbitrarySettings()` but without `id`
- `arbitrarySettingsWithNulls()` — same as `arbitrarySettings()` but with boolean and numeric fields forced to `null` for Property 8
- `arbitraryConfigurationOption()` — `fc.record({ id: fc.integer({ min: 1 }), name: fc.string({ minLength: 1 }) })`
- `arbitraryInvalidCustomerName()` — `fc.stringOf(fc.constantFrom(' ', '\t', '\n'))` for Property 5
- `arbitraryInvalidPositiveInt()` — `fc.oneof(fc.constant(0), fc.integer({ max: 0 }))` for Property 5

### shadcn/ui Components to Scaffold

The following components must be added to `src/shared/ui/` before implementation:

| Component | Command | Status |
|---|---|---|
| `select.tsx` | `npx shadcn@latest add select` | New |
| `toast.tsx` + `toaster.tsx` + `use-toast.ts` | `npx shadcn@latest add toast` | New |

Run both commands inside `frontend/`. Components already present: `Form`, `Input`, `Button`, `Skeleton`, `Switch` (via radix), `Label`.
