# Implementation Plan: Devices Add/Edit

## Overview

Extend the existing devices feature with a `DeviceForm` dialog for creating and editing devices. Tasks proceed from scaffolding missing shadcn/ui components through service extension, form implementation, and page wiring — each step building on the last.

## Tasks

- [ ] 1. Scaffold missing shadcn/ui components
  - Run `npx shadcn@latest add command popover checkbox` inside `frontend/` to generate the three components into `src/shared/ui/`
  - Verify `command.tsx`, `popover.tsx`, and `checkbox.tsx` exist in `src/shared/ui/`
  - _Requirements: 8.3_

- [ ] 2. Extend device types
  - [ ] 2.1 Add `DevicePayload` and `ConfigurationOption` interfaces to `src/features/devices/types.ts`
    - `DevicePayload`: `id?`, `number`, `description?`, `configurationId?`, `groups?`, `imei?`, `phone?`, `custom1?`, `custom2?`, `custom3?`
    - `ConfigurationOption`: `id: number`, `name: string`
    - _Requirements: 10.1, 10.2_

  - [ ]* 2.2 Write property test for null-safe form rendering (Property 12)
    - **Property 12: Null-safe form rendering for any device with null optional fields**
    - **Validates: Requirements 3.1–3.7**
    - Use `arbitraryDeviceViewWithNulls()` with all optional fields null; assert opening `DeviceForm` in edit mode does not throw

- [ ] 3. Extend device service
  - [ ] 3.1 Add `createDevice(payload: DevicePayload): Promise<void>` to `src/features/devices/deviceService.ts`
    - Call `POST /rest/private/devices` with payload; use `assertHmdmOk` on response
    - _Requirements: 10.1, 10.3, 10.4_

  - [ ] 3.2 Add `getGroups(): Promise<LookupItem[]>` to `src/features/devices/deviceService.ts`
    - Call `GET /rest/private/groups`; use `unwrapHmdmData` to extract the array
    - _Requirements: 8.1_

  - [ ]* 3.3 Write unit tests for new service functions
    - Mock `apiClient`; assert `createDevice` calls `POST /rest/private/devices` with the payload
    - Assert `getGroups` calls `GET /rest/private/groups`
    - Assert both functions reject when `apiClient` rejects
    - _Requirements: 10.1, 10.3_

  - [ ]* 3.4 Write property test for service createDevice URL routing (Property 9)
    - **Property 9: Service createDevice routes to correct URL for any payload**
    - **Validates: Requirements 10.1, 10.3**
    - Use `arbitraryDevicePayload()` to generate payloads; mock `apiClient` and assert `POST /rest/private/devices` is called with exact payload; run 100 iterations

  - [ ]* 3.5 Write property test for service updateDevice URL routing (Property 10)
    - **Property 10: Service updateDevice routes to correct URL for any payload**
    - **Validates: Requirements 10.2, 10.3**
    - Use `arbitraryDevicePayload()` with `id` set; mock `apiClient` and assert `PUT /rest/private/devices` is called with exact payload; run 100 iterations

  - [ ]* 3.6 Write property test for service error propagation (Property 11)
    - **Property 11: Service error propagation**
    - **Validates: Requirements 10.4**
    - For `createDevice` and `updateDevice`, mock `apiClient` to reject; assert the returned promise rejects; run 100 iterations

- [ ] 4. Implement DeviceForm component
  - [ ] 4.1 Create `src/features/devices/DeviceForm.tsx`
    - Wrap a shadcn/ui `Dialog` around a `Form` (react-hook-form + zod schema from design)
    - Fields: `number` (Input, required, disabled in edit mode), `description` (Textarea, optional), `configurationId` (Configuration_Selector, optional), `groups` (Group_Selector, optional), `imei` (Input, optional, 15-digit validation when non-empty), `phone` (Input, optional), `custom1`, `custom2`, `custom3` (Input, optional)
    - Fetch groups and configurations in parallel on mount via `Promise.all([getGroups(), getConfigurations()])`
    - Manage `submitting`, `submitError`, `groups`, `groupsLoading`, `groupsError`, `configurations`, `configurationsLoading`, `configurationsError` state
    - In create mode all fields start empty; in edit mode pre-populate from `initialData`
    - On submit: call `createDevice` (create mode) or `updateDevice` with `{ ...formValues, id: initialData.id }` (edit mode)
    - Call `onSuccess()` then `onClose()` on successful submit; call `onClose()` on Cancel without any API call
    - _Requirements: 2.3, 2.4, 3.1–3.8, 4.1–4.4, 5.1–5.4, 6.1–6.4, 7.1, 8.1–8.5, 9.1–9.5_

  - [ ] 4.2 Implement Group_Selector inside DeviceForm
    - Render a `Popover` trigger button showing selected group names (comma-separated, or "Select groups…" when empty)
    - Inside the popover render a scrollable list of `Checkbox` items — one per group from the fetched list
    - Disable the entire control while `groupsLoading` is true; show error text if `groupsError` is set
    - _Requirements: 8.2, 8.3, 8.5_

  - [ ] 4.3 Implement Configuration_Selector inside DeviceForm
    - Render a shadcn/ui `Select` with a "None" option (value `""`) and one `SelectItem` per configuration
    - Disable the control while `configurationsLoading` is true; show error text if `configurationsError` is set
    - Map selected value back to `configurationId: number | null` (empty string → null)
    - _Requirements: 9.2, 9.3, 9.5_

  - [ ]* 4.4 Write unit tests for DeviceForm
    - Assert create mode: title is "Add Device", `number` field is editable, all fields empty
    - Assert edit mode: title is "Edit Device", `number` field is disabled, fields pre-populated from `initialData`
    - Assert `getGroups` and `getConfigurations` called on mount
    - Assert submit button disabled while submitting; re-enabled on failure
    - Assert error message shown on submit failure; dialog stays open
    - Assert `onSuccess` and `onClose` called on successful submit
    - Assert `onClose` called on Cancel; no API call made
    - _Requirements: 2.3, 2.4, 3.8, 4.1, 5.2–5.4, 6.2–6.4, 7.1_

  - [ ]* 4.5 Write property test for edit mode pre-population (Property 1)
    - **Property 1: Edit mode pre-populates all fields for any device**
    - **Validates: Requirements 2.2, 3.1–3.7**
    - Use `arbitraryDeviceView()` to generate devices; assert rendered form fields match device values; run 100 iterations

  - [ ]* 4.6 Write property test for create mode submit (Property 2)
    - **Property 2: Create mode submit calls POST with form values**
    - **Validates: Requirements 5.1, 10.1**
    - Use `arbitraryDevicePayload()` to generate inputs; assert `createDevice` called with exact payload (no `id`); run 100 iterations

  - [ ]* 4.7 Write property test for edit mode submit (Property 3)
    - **Property 3: Edit mode submit calls PUT with form values and device id**
    - **Validates: Requirements 6.1, 10.2**
    - Use `arbitraryDeviceView()` and updated form values; assert `updateDevice` called with payload including original `id`; run 100 iterations

  - [ ]* 4.8 Write property test for required field validation (Property 4)
    - **Property 4: Required field validation rejects empty or whitespace number**
    - **Validates: Requirements 4.1**
    - Use `arbitraryEmptyOrWhitespace()` for `number`; assert form rejects and no API call is made; run 100 iterations

  - [ ]* 4.9 Write property test for IMEI validation (Property 5)
    - **Property 5: IMEI validation accepts 15-digit strings and rejects all others**
    - **Validates: Requirements 4.2, 4.4**
    - Use `arbitraryValidImei()` and `arbitraryInvalidImei()`; assert valid IMEI passes and invalid IMEI fails; run 100 iterations each

  - [ ]* 4.10 Write property test for cancel makes no API call (Property 6)
    - **Property 6: Cancel makes no API call**
    - **Validates: Requirements 7.1**
    - For any dialog state, clicking Cancel must not invoke `createDevice` or `updateDevice`; run 100 iterations

  - [ ]* 4.11 Write property test for group selector pre-selection (Property 7)
    - **Property 7: Group selector pre-selects assigned groups for any device**
    - **Validates: Requirements 8.4**
    - Use `arbitraryDeviceView()` with non-empty `groups`; assert each group is checked in the Group_Selector; run 100 iterations

  - [ ]* 4.12 Write property test for configuration selector pre-selection (Property 8)
    - **Property 8: Configuration selector pre-selects assigned configuration for any device**
    - **Validates: Requirements 9.4**
    - Use `arbitraryDeviceView()` with non-null `configurationId`; assert matching configuration is selected; run 100 iterations

- [ ] 5. Wire DeviceForm into DevicesPage
  - [ ] 5.1 Add `formMode` and `deviceToEdit` state to `src/features/devices/DevicesPage.tsx`
    - `formMode: 'create' | 'edit' | null` — controls dialog open state
    - `deviceToEdit: DeviceView | null` — device being edited; null in create mode
    - _Requirements: 1.1, 2.1_

  - [ ] 5.2 Add "Add Device" button to DevicesPage header
    - Render a `Button` labeled "Add Device" in the header area alongside the search input
    - On click: set `formMode = 'create'` and `deviceToEdit = null`
    - _Requirements: 1.1, 1.2_

  - [ ] 5.3 Add "Edit" item to RowActionsMenu in DevicesPage
    - Add an "Edit" `DropdownMenuItem` above the existing "Delete" item in the row actions dropdown
    - On click: set `deviceToEdit = device` and `formMode = 'edit'`
    - _Requirements: 2.1, 2.2_

  - [ ] 5.4 Render DeviceForm dialog in DevicesPage
    - Render `<DeviceForm mode={formMode} initialData={deviceToEdit} onSuccess={refresh} onClose={() => setFormMode(null)} />` when `formMode` is non-null
    - Pass `open={formMode !== null}` to the underlying Dialog
    - _Requirements: 1.2, 1.3, 2.2, 2.3, 5.3, 6.3_

  - [ ]* 5.5 Write unit tests for DevicesPage wiring
    - Assert "Add Device" button is rendered
    - Assert clicking "Add Device" opens `DeviceForm` in create mode with empty fields
    - Assert clicking "Edit" row action opens `DeviceForm` in edit mode with device data
    - Assert `DeviceForm` is not rendered when `formMode` is null
    - _Requirements: 1.1, 1.2, 2.1, 2.2_

- [ ] 6. Checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all unit and property tests pass; ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- The existing `updateDevice` in `deviceService.ts` accepts a full `DeviceView`; `DeviceForm` constructs the payload by spreading form values and adding `id: initialData.id`
- `groups` is sent as `LookupItem[]` — the same shape stored on `DeviceView.groups`
- `configurationId` maps to a `Select` value; empty string selection maps to `null` in the payload
- Custom field labels (`custom1`, `custom2`, `custom3`) ideally come from settings, but for MVP they can be labeled "Custom 1", "Custom 2", "Custom 3"
- Property tests use fast-check with a minimum of 100 iterations each
- Generators (`arbitraryDeviceView`, `arbitraryDevicePayload`, `arbitraryValidImei`, `arbitraryInvalidImei`, `arbitraryEmptyOrWhitespace`, `arbitraryDeviceViewWithNulls`) should be defined in a shared test helper and imported by all property test files
