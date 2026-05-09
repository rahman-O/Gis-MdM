# Implementation Plan: Settings Management

## Overview

Implement the `/settings` page following the feature-slice pattern established by configurations-management and users-management. The route and navigation entry already exist; tasks proceed from scaffolding missing shadcn/ui components through types, service, page implementation, and final wiring.

## Tasks

- [ ] 1. Scaffold missing shadcn/ui components
  - Run `npx shadcn@latest add select` inside `frontend/` to generate `src/shared/ui/select.tsx`
  - Run `npx shadcn@latest add toast` inside `frontend/` to generate `src/shared/ui/toast.tsx`, `src/shared/ui/toaster.tsx`, and `src/shared/hooks/use-toast.ts`
  - Add `<Toaster />` to the root layout component (`src/features/layout/Layout.tsx` or `src/app/App.tsx`) so toasts render globally
  - Verify all generated files exist before proceeding
  - _Requirements: 3.3, 3.4, 5.1, 5.2, 5.3, 5.4_

- [ ] 2. Define TypeScript types
  - [ ] 2.1 Create `src/features/settings/types.ts`
    - Define `Settings`, `SettingsPayload`, and `ConfigurationOption` interfaces exactly as specified in the design
    - `SettingsPayload` omits `id` from `Settings` and marks all fields as required
    - `newDeviceConfigurationId` is `number | null` in both interfaces
    - _Requirements: 8.1, 8.2_

- [ ] 3. Implement the settings service
  - [ ] 3.1 Create `src/features/settings/settingsService.ts`
    - Implement `getSettings(): Promise<Settings>` calling `GET /rest/private/settings` via `apiClient` and unwrapping with `unwrapHmdmData`
    - Implement `updateSettings(data: SettingsPayload): Promise<Settings>` calling `PUT /rest/private/settings` via `apiClient` and unwrapping with `unwrapHmdmData`
    - All non-2xx responses or `status: "ERROR"` envelopes must propagate as thrown errors
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [ ]* 3.2 Write property test for service URL routing (Property 6)
    - **Property 6: Service routes to correct URLs**
    - **Validates: Requirements 7.1, 7.2**
    - Mock `apiClient`; assert `getSettings` calls `GET /rest/private/settings`; use `arbitrarySettingsPayload()` and assert `updateSettings` calls `PUT /rest/private/settings` with the exact payload; run 100 iterations

  - [ ]* 3.3 Write property test for service error propagation (Property 7)
    - **Property 7: Service error propagation**
    - **Validates: Requirements 7.4**
    - For each service function, mock `apiClient` to reject and assert the returned promise rejects; run 100 iterations

- [ ] 4. Implement SettingsPage
  - [ ] 4.1 Replace `src/features/settings/SettingsPage.tsx` with full implementation
    - Fetch settings and configurations in parallel on mount via `Promise.all([settingsService.getSettings(), apiClient.get('/rest/private/configurations')])`
    - Manage `settings`, `configurations`, `loading`, `error`, and `submitting` state
    - While loading render `Skeleton` placeholders in place of the form
    - On GET error render an error banner with a "Retry" button; hide the form
    - Once data loads, render the form pre-populated using `form.reset()` with null-coerced values (booleans default `false`, numerics default `0`)
    - Form fields (all using shadcn/ui Form + react-hook-form + zod resolver):
      - `customerName`: `Input` (required)
      - `createNewDevices`: `Switch`
      - `newDeviceConfigurationId`: `Select` populated from configurations
      - `language`: `Select` with fixed language options (en, ru, de, fr, es, pt, zh)
      - `passwordLength`: `Input` type=number (required, min 1)
      - `passwordStrength`: `Select` with options 0–3
      - `sendDeviceInfoExpiryDays`: `Input` type=number (required, min 1)
      - `unsecureEnrollment`: `Switch`
      - `deviceFastSearch`: `Switch`
    - "Save Settings" `Button` type=submit; disabled and shows spinner while `submitting` is true
    - On valid submit call `settingsService.updateSettings(payload)`; on success show toast "Settings saved" and reset form to returned values; on failure show toast "Failed to save settings" with destructive style and retain form values
    - Validation errors displayed adjacent to each field via `FormMessage`
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8, 2.9, 2.10, 3.1, 3.2, 3.3, 3.4, 3.5, 4.1, 4.2, 4.3, 4.4, 4.5, 5.1, 5.2, 8.3, 8.4_

  - [ ]* 4.2 Write property test for form populated with any settings (Property 1)
    - **Property 1: Form populated with any settings object**
    - **Validates: Requirements 1.4, 8.3, 8.4**
    - Use `arbitrarySettings()` to generate settings objects; mock `settingsService.getSettings` to resolve with each; assert every form field reflects the corresponding value after load; run 100 iterations

  - [ ]* 4.3 Write property test for configurations populate select (Property 2)
    - **Property 2: Configurations populate select options**
    - **Validates: Requirements 2.3**
    - Use `fc.array(arbitraryConfigurationOption(), { minLength: 1 })` to generate configuration lists; mock the configurations endpoint; assert the `newDeviceConfigurationId` select renders exactly one option per configuration with correct label and value; run 100 iterations

  - [ ]* 4.4 Write property test for submit calls service with form values (Property 3)
    - **Property 3: Submit calls service with form values for any valid payload**
    - **Validates: Requirements 3.1, 7.2**
    - Use `arbitrarySettingsPayload()` to generate valid payloads; fill the form and submit; assert `settingsService.updateSettings` is called with exactly those values; run 100 iterations

  - [ ]* 4.5 Write property test for save button re-enabled after any PUT outcome (Property 4)
    - **Property 4: Save button re-enabled after any PUT outcome**
    - **Validates: Requirements 3.5**
    - For both success and failure outcomes of `settingsService.updateSettings`, assert the submit button is not disabled after the request completes; run 100 iterations

  - [ ]* 4.6 Write property test for validation rejects invalid required fields (Property 5)
    - **Property 5: Validation rejects invalid required fields**
    - **Validates: Requirements 4.1, 4.2, 4.3, 4.5**
    - Use `arbitraryInvalidCustomerName()` (whitespace-only strings) and `arbitraryInvalidPositiveInt()` (values ≤ 0) for `passwordLength` and `sendDeviceInfoExpiryDays`; assert form rejects submission, shows `FormMessage` errors, and makes no call to `settingsService.updateSettings`; run 100 iterations

  - [ ]* 4.7 Write property test for null-safe rendering (Property 8)
    - **Property 8: Null-safe rendering for any settings with null fields**
    - **Validates: Requirements 8.3, 8.4**
    - Use `arbitrarySettingsWithNulls()` (boolean and numeric fields forced to null); mock `settingsService.getSettings` to resolve with each; assert rendering does not throw and null booleans render as `false`, null numerics render as `0`; run 100 iterations

- [ ] 5. Checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all unit and property tests pass; ask the user if questions arise.

- [ ] 6. Verify routing and navigation
  - [ ] 6.1 Confirm `/settings` route is registered in `src/app/App.tsx` (already present — verify only, no change needed)
    - _Requirements: 1.1, 6.2_

  - [ ] 6.2 Confirm "Settings" entry is present in `src/features/layout/navItems.ts` (already present — verify only, no change needed)
    - _Requirements: 6.1, 6.2, 6.3_

- [ ] 7. Final checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all tests pass end-to-end; ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Each task references specific requirements for traceability
- Property tests use fast-check with a minimum of 100 iterations each
- Each property test file must include a comment referencing the design property number and the requirements clause it validates
- Generators (`arbitrarySettings`, `arbitrarySettingsPayload`, `arbitrarySettingsWithNulls`, `arbitraryConfigurationOption`, `arbitraryInvalidCustomerName`, `arbitraryInvalidPositiveInt`) should be defined once in a shared test helper file and imported by all property test files
- The `/settings` route and "Settings" nav entry already exist — tasks 6.1 and 6.2 are verification-only
- `select.tsx` and the toast components are the only shadcn/ui additions needed
- The `SettingsPage` is a single-page form (not a dialog) — no separate form component file is needed
- Configurations are fetched via the existing `apiClient` directly (no need for a separate `configurationService` import unless it already exists)
