# Implementation Plan: Configurations Management

## Overview

Implement the `/configurations` page following the feature-slice pattern established by devices-management. Tasks proceed from scaffolding shared UI components through types, service, form, page, routing, and navigation — each step building on the last.

## Tasks

- [x] 1. Scaffold required shadcn/ui components
  - Run `npx shadcn@latest add alert-dialog dropdown-menu select dialog` inside `frontend/` to generate the four components into `src/shared/ui/`
  - Verify `alert-dialog.tsx`, `dropdown-menu.tsx`, `select.tsx`, and `dialog.tsx` exist in `src/shared/ui/`
  - _Requirements: 2.2, 3.1, 4.1_
  - **Done:** `dialog.tsx`, `select.tsx`, `textarea.tsx` added manually (alert-dialog + dropdown-menu already existed); matches shadcn patterns.

- [x] 2. Define TypeScript types
  - [x] 2.1 Create `src/features/configurations/types.ts`
    - Define `ConfigurationFile`, `ConfigurationApplication`, `Configuration`, and `ConfigurationPayload` interfaces exactly as specified in the design
    - `Configuration.deviceCount` is `number | null | undefined`; optional fields default to `null` or `[]` at runtime
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5_

  - [x]* 2.2 Write property test for null-safe field access (Property 12)
    - **Property 12: Null-safe field access**
    - **Validates: Requirements 9.5**
    - Use `fc.record` with nullable optional fields; assert no runtime error when rendering a row with all nulls

- [x] 3. Implement the configuration service
  - [x] 3.1 Create `src/features/configurations/configurationService.ts`
    - Implement `getConfigurations()`, `getConfiguration(id)`, `createConfiguration(data)`, `updateConfiguration(id, data)`, `deleteConfiguration(id)` using the shared `apiClient` and `unwrapHmdmData` helper
    - All non-2xx responses or `status: "ERROR"` envelopes must propagate as thrown errors
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7_

  - [x]* 3.2 Write property test for service URL routing (Property 10)
    - **Property 10: Service routes to correct URL for any operation**
    - **Validates: Requirements 8.1, 8.2, 8.3, 8.4, 8.5**
    - Use `fc.integer({ min: 1 })` for ids and `arbitraryConfigurationPayload()` for payloads; mock `apiClient` and assert method + URL for each function; run 100 iterations

  - [x]* 3.3 Write property test for service error propagation (Property 11)
    - **Property 11: Service error propagation**
    - **Validates: Requirements 8.7**
    - For each service function, mock `apiClient` to reject and assert the returned promise rejects; run 100 iterations

- [x] 4. Implement ConfigurationForm
  - [x] 4.1 Create `src/features/configurations/ConfigurationForm.tsx`
    - Wrap a shadcn/ui `Dialog` around a `Form` (react-hook-form + zod schema from design)
    - Fields: `name` (Input, required), `type` (Select with COMMON/WORK options, required), `description` (Input/Textarea, optional)
    - Manage `submitting` and `submitError` state; disable submit button while submitting; show `FormMessage` for validation errors
    - In create mode all fields start empty; in edit mode pre-populate from `initialData`
    - Call `onSuccess` + `onClose` on successful submit; call `onClose` on cancel without any API call
    - _Requirements: 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 3.1, 3.2, 3.3, 3.4, 3.5, 5.1, 5.2, 5.3, 5.4, 5.5_

  - [x]* 4.2 Write property test for form submit routing (Property 2)
    - **Property 2: Form submit calls correct endpoint with form values**
    - **Validates: Requirements 2.3, 3.2**
    - Use `arbitraryConfigurationPayload()` to generate inputs; assert create mode calls `createConfiguration` with exact payload and edit mode calls `updateConfiguration` with id + payload; run 100 iterations

  - [x]* 4.3 Write property test for cancel makes no API call (Property 3)
    - **Property 3: Cancel makes no API call**
    - **Validates: Requirements 2.7, 4.6**
    - For any dialog state, clicking Cancel must not invoke any service function; run 100 iterations

  - [x]* 4.4 Write property test for edit mode pre-population (Property 4)
    - **Property 4: Edit mode pre-populates fields for any configuration**
    - **Validates: Requirements 3.1**
    - Use `arbitraryConfiguration()` to generate configs; assert rendered form fields match `name`, `description`, and `type`; run 100 iterations

  - [x]* 4.5 Write property test for required field validation (Property 7)
    - **Property 7: Required field validation rejects empty or whitespace inputs**
    - **Validates: Requirements 5.1, 5.2**
    - Use `arbitraryEmptyOrWhitespace()` for name and invalid type values; assert form rejects and no API call is made; run 100 iterations

  - [x]* 4.6 Write property test for optional description (Property 8)
    - **Property 8: Optional description allows submission**
    - **Validates: Requirements 5.3**
    - Use `arbitraryConfigurationPayload()` with description omitted or empty; assert form submits successfully; run 100 iterations

- [x] 5. Implement ConfigurationsPage
  - [x] 5.1 Create `src/features/configurations/ConfigurationsPage.tsx`
    - Fetch configurations on mount via `configurationService.getConfigurations()`; manage `configurations`, `loading`, `error`, `formMode`, `selectedConfig`, and `configToDelete` state
    - Render a shadcn/ui `Table` with columns: Name, Type, Description, Device Count; each row has a `DropdownMenu` with Edit and Delete actions
    - Show loading skeleton while fetching; show error banner with Retry button on failure; show "No configurations found" empty state
    - Render "New Configuration" button above the table that opens `ConfigurationForm` in create mode
    - Inline `DeleteDialog` (wrapping `AlertDialog`) that shows the config name, manages `deleting`/`deleteError` state, calls `configurationService.deleteConfiguration`, and refreshes on success
    - After any successful create, update, or delete, re-fetch the configuration list
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 2.1, 2.5, 3.4, 4.1, 4.4, 6.1, 6.2, 6.3_
    - **Note:** Delete UI lives in `ConfigurationDeleteDialog.tsx` (same AlertDialog pattern as devices).

  - [x]* 5.2 Write property test for table columns (Property 1)
    - **Property 1: Table columns rendered for any configuration list**
    - **Validates: Requirements 1.3, 6.1**
    - Use `fc.array(arbitraryConfiguration(), { minLength: 1 })` to generate lists; assert rendered table always has Name, Type, Description, and Device Count columns; run 100 iterations

  - [x]* 5.3 Write property test for device count display (Property 9)
    - **Property 9: Device count displayed correctly for any configuration**
    - **Validates: Requirements 6.2, 6.3**
    - Use `arbitraryConfiguration()` with `deviceCount` set to `fc.integer({ min: 0 })`; assert each row displays the exact count; run 100 iterations

  - [x]* 5.4 Write property test for delete dialog shows config name (Property 5)
    - **Property 5: Delete dialog shows configuration name for any configuration**
    - **Validates: Requirements 4.1**
    - Use `arbitraryConfiguration()` to generate configs; open delete dialog and assert config name appears in dialog body; run 100 iterations

  - [x]* 5.5 Write property test for delete confirm calls correct id (Property 6)
    - **Property 6: Delete confirm calls DELETE with correct id**
    - **Validates: Requirements 4.2**
    - Use `arbitraryConfiguration()` to generate configs; confirm deletion and assert `deleteConfiguration` called with exact `id`; run 100 iterations

  - [x]* 5.6 Write property test for null-safe rendering (Property 12)
    - **Property 12: Null-safe rendering for any configuration with null optional fields**
    - **Validates: Requirements 9.5**
    - Use `arbitraryConfigurationWithNulls()` (all optional fields null/[]); assert rendering table row and opening form does not throw; run 100 iterations

- [x] 6. Checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all unit and property tests pass; ask the user if questions arise.

- [x] 7. Wire routing and navigation
  - [x] 7.1 Update `src/app/App.tsx` — add `<Route path="/configurations" element={<ConfigurationsPage />} />`
    - Import `ConfigurationsPage` and register the route alongside the existing `/devices` route
    - _Requirements: 1.1, 7.2_

  - [x] 7.2 Update `src/features/layout/navItems.ts` — add Configurations entry after Devices
    - Add `{ label: 'Configurations', path: '/configurations', icon: Settings2 }` (or `LayoutList`) after the Devices entry
    - _Requirements: 7.1, 7.2, 7.3_
    - **Done:** `LayoutList` icon used.

- [x] 8. Final checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all tests pass end-to-end; ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Each task references specific requirements for traceability
- Property tests use fast-check with a minimum of 100 iterations each
- Each property test file must include a comment referencing the design property number and the requirements clause it validates
- Generators (`arbitraryConfiguration`, `arbitraryConfigurationPayload`, `arbitraryEmptyOrWhitespace`, `arbitraryConfigurationWithNulls`) should be defined once in a shared test helper file and imported by all property test files
