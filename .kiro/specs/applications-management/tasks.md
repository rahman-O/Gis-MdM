# Implementation Plan: Applications Management

## Overview

Implement the `/applications` page following the feature-slice pattern established by devices-management and configurations-management. Tasks proceed from scaffolding the one missing shadcn/ui component through types, service, form, page, routing, and navigation — each step building on the last.

## Tasks

- [ ] 1. Scaffold required shadcn/ui components
  - Run `npx shadcn@latest add switch` inside `frontend/` to generate `switch.tsx` into `src/shared/ui/`
  - Verify `switch.tsx` exists in `src/shared/ui/`
  - All other required components (`table`, `dialog`, `alert-dialog`, `dropdown-menu`, `skeleton`, `badge`, `button`, `input`, `form`) are already present
  - _Requirements: 5.5_

- [ ] 2. Define TypeScript types
  - [ ] 2.1 Create `src/features/applications/types.ts`
    - Define `Application`, `ApplicationPayload`, and `ApkUploadResponse` interfaces exactly as specified in the design
    - `Application` boolean fields (`system`, `showIcon`, `runAfterInstall`, `runAtBoot`) must be `boolean` (not `boolean | null`) — callers coerce null to false at the service boundary
    - _Requirements: 9.1, 9.2, 9.3_

- [ ] 3. Implement the application service
  - [ ] 3.1 Create `src/features/applications/applicationService.ts`
    - Implement `getApplications()`, `createApplication(data)`, `updateApplication(id, data)`, `deleteApplication(id)`, and `uploadApk(file)` using the shared `apiClient` and `unwrapHmdmData` helper
    - `uploadApk` must send the file as `multipart/form-data` using a `FormData` object; coerce null boolean fields to `false` in `getApplications`
    - All non-2xx responses or `status: "ERROR"` envelopes must propagate as thrown errors
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7_

  - [ ]* 3.2 Write property test for service URL routing (Property 12)
    - **Property 12: Service routes to correct URL for any operation**
    - **Validates: Requirements 8.1, 8.2, 8.3, 8.4, 8.5**
    - Use `fc.integer({ min: 1 })` for ids and `arbitraryApplicationPayload()` for payloads; mock `apiClient` and assert method + URL for each function; run 100 iterations

  - [ ]* 3.3 Write property test for service error propagation (Property 13)
    - **Property 13: Service error propagation**
    - **Validates: Requirements 8.7**
    - For each service function, mock `apiClient` to reject and assert the returned promise rejects; run 100 iterations

- [ ] 4. Implement ApplicationForm
  - [ ] 4.1 Create `src/features/applications/ApplicationForm.tsx`
    - Wrap a shadcn/ui `Dialog` around a `Form` (react-hook-form + zod schema from design)
    - Fields: `name` (Input, required), `pkg` (Input, required), `version` (Input, required), `url` (Input, required), APK file input (`accept=".apk"`, optional), `system` / `showIcon` / `runAfterInstall` / `runAtBoot` (Switch, default false)
    - Manage `submitting`, `submitError`, `uploading`, and `uploadError` state
    - On file selection: validate `.apk` extension, call `applicationService.uploadApk`, populate `url` and auto-fill `name`/`pkg`/`version` if returned; show upload progress and disable file input while uploading
    - Disable submit button while submitting; show `FormMessage` for validation errors
    - In create mode all fields start empty; in edit mode pre-populate from `initialData`
    - Call `onSuccess` + `onClose` on successful submit; call `onClose` on cancel without any API call
    - _Requirements: 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.9, 2.10, 3.1, 3.2, 3.3, 3.5, 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

  - [ ]* 4.2 Write property test for form submit routing (Property 3)
    - **Property 3: Form submit calls correct endpoint based on mode**
    - **Validates: Requirements 2.6, 3.2**
    - Use `arbitraryApplicationPayload()` to generate inputs; assert create mode calls `createApplication` with exact payload and edit mode calls `updateApplication` with id + payload; run 100 iterations

  - [ ]* 4.3 Write property test for cancel makes no API call (Property 4)
    - **Property 4: Cancel makes no API call**
    - **Validates: Requirements 2.10, 4.6**
    - For any dialog state, clicking Cancel must not invoke any service function; run 100 iterations

  - [ ]* 4.4 Write property test for edit mode pre-population (Property 5)
    - **Property 5: Edit mode pre-populates all fields for any application**
    - **Validates: Requirements 3.1**
    - Use `arbitraryApplication()` to generate apps; assert rendered form fields match all eight fields; run 100 iterations

  - [ ]* 4.5 Write property test for required field validation (Property 8)
    - **Property 8: Required field validation rejects empty or whitespace inputs**
    - **Validates: Requirements 5.1, 5.2, 5.3, 5.4**
    - Use `arbitraryEmptyOrWhitespace()` for each required field (name, pkg, version, url); assert form rejects and no API call is made; run 100 iterations

  - [ ]* 4.6 Write property test for non-APK file rejection (Property 10)
    - **Property 10: Non-APK file rejected before upload starts**
    - **Validates: Requirements 6.2**
    - Use `arbitraryNonApkFile()` to generate file names; assert validation error shown and `uploadApk` not called; run 100 iterations

  - [ ]* 4.7 Write property test for APK upload sets url field (Property 11)
    - **Property 11: APK upload sets url field to returned URL**
    - **Validates: Requirements 2.5, 6.5**
    - Use `arbitraryApkUploadResponse()` to generate responses; mock `uploadApk` and assert url field is set to response url; run 100 iterations

- [ ] 5. Implement ApplicationsPage
  - [ ] 5.1 Create `src/features/applications/ApplicationsPage.tsx`
    - Fetch applications on mount via `applicationService.getApplications()`; manage `applications`, `loading`, `error`, `formMode`, `selectedApp`, and `appToDelete` state
    - Render a shadcn/ui `Table` with columns: Name, Package, Version, System; each row has a `DropdownMenu` with Edit and Delete actions
    - System column: render a `Badge` labeled "System" when `system` is `true`; render a muted dash when `false`
    - Show loading skeleton while fetching; show error banner with Retry button on failure; show "No applications found" empty state
    - Render "Add Application" button above the table that opens `ApplicationForm` in create mode
    - Inline `DeleteDialog` (wrapping `AlertDialog`) that shows the app name and package, manages `deleting`/`deleteError` state, calls `applicationService.deleteApplication`, and refreshes on success
    - After any successful create, update, or delete, re-fetch the application list
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 2.1, 2.8, 3.4, 4.1, 4.3, 4.4_

  - [ ]* 5.2 Write property test for table columns (Property 1)
    - **Property 1: Table columns rendered for any application list**
    - **Validates: Requirements 1.3**
    - Use `fc.array(arbitraryApplication(), { minLength: 1 })` to generate lists; assert rendered table always has Name, Package, Version, and System columns; run 100 iterations

  - [ ]* 5.3 Write property test for system flag indicator (Property 2)
    - **Property 2: System flag indicator rendered correctly for any application**
    - **Validates: Requirements 1.4**
    - Use `arbitraryApplication()` with `fc.boolean()` for `system`; assert System column renders badge when true and empty/dash when false; run 100 iterations

  - [ ]* 5.4 Write property test for delete dialog shows name and package (Property 6)
    - **Property 6: Delete dialog shows application name and package for any application**
    - **Validates: Requirements 4.1**
    - Use `arbitraryApplication()` to generate apps; open delete dialog and assert both `name` and `pkg` appear in dialog body; run 100 iterations

  - [ ]* 5.5 Write property test for delete confirm calls correct id (Property 7)
    - **Property 7: Delete confirm calls DELETE with correct id**
    - **Validates: Requirements 4.2**
    - Use `arbitraryApplication()` to generate apps; confirm deletion and assert `deleteApplication` called with exact `id`; run 100 iterations

  - [ ]* 5.6 Write property test for null-safe boolean field handling (Property 14)
    - **Property 14: Null-safe boolean field handling**
    - **Validates: Requirements 9.4**
    - Use `arbitraryApplicationWithNullBooleans()` (all boolean fields null); assert rendering table row and opening form does not throw and booleans are treated as false; run 100 iterations

- [ ] 6. Checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all unit and property tests pass; ask the user if questions arise.

- [ ] 7. Wire routing and navigation
  - [ ] 7.1 Update `src/app/App.tsx` — add `<Route path="/applications" element={<ApplicationsPage />} />`
    - Import `ApplicationsPage` and register the route alongside the existing `/devices` and `/configurations` routes
    - _Requirements: 1.1, 7.2_

  - [ ] 7.2 Update `src/features/layout/navItems.ts` — add Applications entry
    - Add `{ label: 'Applications', path: '/applications', icon: Package }` (from `lucide-react`) after the Configurations entry
    - _Requirements: 7.1, 7.2, 7.3_

- [ ] 8. Final checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all tests pass end-to-end; ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Each task references specific requirements for traceability
- Property tests use fast-check with a minimum of 100 iterations each
- Each property test file must include a comment referencing the design property number and the requirements clause it validates
- Generators (`arbitraryApplication`, `arbitraryApplicationPayload`, `arbitraryEmptyOrWhitespace`, `arbitraryNonApkFile`, `arbitraryApkUploadResponse`, `arbitraryApplicationWithNullBooleans`) should be defined once in a shared test helper file and imported by all property test files
