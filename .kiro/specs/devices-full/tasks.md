# Implementation Plan: Devices Full

## Overview

Implement the complete devices feature: paginated list, search, advanced filters, bulk actions, add/edit form, detail panel, status badges, delete, QR enrollment, configurable columns, and a Groups management page. Tasks proceed from scaffolding shared UI components through types, services, components, page wiring, routing, and navigation — each step building on the last.

## Tasks

- [ ] 1. Scaffold required shadcn/ui components
  - Run `npx shadcn@latest add command popover checkbox textarea` inside `frontend/` to generate the four components into `src/shared/ui/`
  - Verify `command.tsx`, `popover.tsx`, `checkbox.tsx`, and `textarea.tsx` exist in `src/shared/ui/`
  - _Requirements: 9.2, 9.4, 11.3, 15.1_

- [ ] 2. Define TypeScript types
  - [ ] 2.1 Create `src/features/devices/types.ts`
    - Define `LookupItem`, `ConfigurationView`, `ConfigurationOption`, `DeviceInfoView`, `DeviceView`, `DevicePayload`, `DeviceSearchRequest`, `DeviceListResponse`, and `DeviceFilters` interfaces exactly as specified in the design
    - All optional fields use `| null` not `| undefined` to match backend behavior
    - _Requirements: 21.1, 21.2, 21.3, 21.4, 21.5_

  - [ ]* 2.2 Write property test for DevicePayload serialization round-trip (Property 34)
    - **Property 34: DevicePayload serialization round-trip**
    - **Validates: Requirements 27.3**
    - Use `arbitraryDevicePayload()` to generate payloads; assert `JSON.parse(JSON.stringify(payload))` produces an equivalent object; run 100 iterations

  - [ ]* 2.3 Write property test for LookupItem serialization round-trip (Property 35)
    - **Property 35: LookupItem serialization round-trip**
    - **Validates: Requirements 27.4**
    - Use `arbitraryLookupItem()` to generate items; assert round-trip produces equivalent object; run 100 iterations

- [ ] 3. Implement device service
  - [ ] 3.1 Create `src/features/devices/deviceService.ts`
    - Implement all 8 functions: `getDevices`, `getDevice`, `createDevice`, `updateDevice`, `deleteDevice`, `getGroups`, `getConfigurations`, `getDeviceQr`
    - Use shared `apiClient` and `unwrapHmdmData` / `assertHmdmOk` from `src/services/hmdmEnvelope.ts`
    - All non-2xx responses or `status: "ERROR"` envelopes must propagate as thrown errors
    - _Requirements: 20.1–20.10_

  - [ ]* 3.2 Write property test for device service URL routing (Property 28)
    - **Property 28: Device service routes to correct URL for any operation**
    - **Validates: Requirements 20.1–20.8**
    - Use `arbitraryDeviceSearchRequest()`, `arbitraryDevicePayload()`, `fc.integer({ min: 1 })`, `fc.string({ minLength: 1 })` for inputs; mock `apiClient` and assert method + URL for each function; run 100 iterations

  - [ ]* 3.3 Write property test for device service error propagation (Property 29)
    - **Property 29: Device service error propagation**
    - **Validates: Requirements 20.10**
    - For each service function, mock `apiClient` to reject; assert the returned promise rejects; run 100 iterations

- [ ] 4. Implement group service
  - [ ] 4.1 Create `src/features/groups/groupService.ts`
    - Implement `getGroups`, `createGroup`, `updateGroup`, `deleteGroup` using shared `apiClient`
    - All non-2xx responses or `status: "ERROR"` envelopes must propagate as thrown errors
    - _Requirements: 26.1–26.6_

  - [ ]* 4.2 Write property test for group service URL routing (Property 32)
    - **Property 32: Group service routes to correct URL for any operation**
    - **Validates: Requirements 26.1–26.4**
    - Use `arbitraryLookupItem()` and `fc.integer({ min: 1 })` for inputs; mock `apiClient` and assert method + URL for each function; run 100 iterations

  - [ ]* 4.3 Write property test for group service error propagation (Property 33)
    - **Property 33: Group service error propagation**
    - **Validates: Requirements 26.6**
    - For each service function, mock `apiClient` to reject; assert the returned promise rejects; run 100 iterations

- [ ] 5. Implement StatusBadge component
  - [ ] 5.1 Create `src/features/devices/StatusBadge.tsx`
    - Map `statusCode` to Badge variant and label: `"green"` → default/green "Online"; `"red"` → destructive "Offline"; anything else → secondary "Unknown"
    - Accept `statusCode: string | null` prop; null treated as Unknown
    - _Requirements: 5.1, 5.4_

  - [ ]* 5.2 Write property test for status badge mapping (Property 1)
    - **Property 1: Status badge mapping is total and correct**
    - **Validates: Requirements 5.1, 5.4**
    - Use `fc.oneof(fc.constant("green"), fc.constant("red"), fc.string(), fc.constant(null))` for statusCode; assert correct label and variant for each; run 100 iterations

- [ ] 6. Implement DeviceDetailPanel component
  - [ ] 6.1 Create `src/features/devices/DeviceDetailPanel.tsx`
    - Render inside a shadcn/ui `Sheet` (side="right")
    - Props: `deviceNumber: string | null`, `onClose: () => void`
    - Call `deviceService.getDevice(number)` when `deviceNumber` changes; manage `detail`, `detailLoading`, `detailError` state
    - Display: device number, StatusBadge, configuration name, group names, battery level (or "N/A"), last update timestamp (formatted as human-readable date/time), GPS coordinates or "Location unavailable"
    - Show loading skeleton while fetching; show error message on failure
    - _Requirements: 4.1–4.5_

  - [ ]* 6.2 Write property test for detail panel fetches by device number (Property 9)
    - **Property 9: Detail panel fetches by device number**
    - **Validates: Requirements 4.2**
    - Use `arbitraryDeviceView()` to generate devices; render panel with device.number; assert `getDevice` called with that number; run 100 iterations

  - [ ]* 6.3 Write property test for detail panel renders all required fields (Property 10)
    - **Property 10: Detail panel renders all required fields for any device**
    - **Validates: Requirements 4.4**
    - Use `arbitraryDeviceView()` to generate devices; assert rendered panel contains number, status badge, config name, group names, battery, timestamp, and location; run 100 iterations

- [ ] 7. Implement DeleteDialog component
  - [ ] 7.1 Create `src/features/devices/DeleteDialog.tsx`
    - Wrap shadcn/ui `AlertDialog`; props: `device: DeviceView | null`, `onConfirm: () => Promise<void>`, `onCancel: () => void`
    - Show device number in dialog body; manage `deleting` and `deleteError` state
    - Disable confirm button while deleting; show error message on failure; close on success
    - _Requirements: 6.1–6.6_

  - [ ]* 7.2 Write property test for delete dialog shows device number (Property 11)
    - **Property 11: Delete dialog shows device number for any device**
    - **Validates: Requirements 6.1**
    - Use `arbitraryDeviceView()` to generate devices; open delete dialog and assert device.number appears in dialog body; run 100 iterations

  - [ ]* 7.3 Write property test for delete confirmation calls correct id (Property 12)
    - **Property 12: Delete confirmation calls service with correct id**
    - **Validates: Requirements 6.2**
    - Use `arbitraryDeviceView()` to generate devices; confirm deletion and assert `deleteDevice` called with exact `id`; run 100 iterations

  - [ ]* 7.4 Write property test for cancel delete makes no API call (Property 13)
    - **Property 13: Cancel delete makes no API call**
    - **Validates: Requirements 6.6**
    - For any device, clicking Cancel must not invoke `deviceService.deleteDevice`; run 100 iterations


- [ ] 8. Implement DeviceForm component
  - [ ] 8.1 Create `src/features/devices/DeviceForm.tsx`
    - Wrap shadcn/ui `Dialog` around a `Form` (react-hook-form + zod schema from design)
    - Fields: `number` (Input, required, disabled in edit mode), `description` (Textarea, optional), `configurationId` (Configuration_Selector), `groups` (Group_Selector), `imei` (Input, 15-digit validation when non-empty), `phone`, `custom1`, `custom2`, `custom3` (Input, optional)
    - Fetch groups and configurations in parallel on mount via `Promise.all([getGroups(), getConfigurations()])`
    - In create mode all fields start empty; in edit mode pre-populate from `initialData`
    - On submit: call `createDevice` (create mode) or `updateDevice` with `{ ...formValues, id: initialData.id }` (edit mode)
    - Call `onSuccess()` then `onClose()` on successful submit; call `onClose()` on Cancel without any API call
    - _Requirements: 7.2–7.7, 8.2–8.8, 9.1–9.10, 10.1–10.4, 11.1–11.5, 12.1–12.5_

  - [ ] 8.2 Implement Group_Selector inside DeviceForm
    - Render a `Popover` trigger button showing selected group names (comma-separated, or "Select groups…" when empty)
    - Inside the popover render a scrollable list of `Checkbox` items — one per group
    - Disable the control while `groupsLoading` is true; show error text if `groupsError` is set
    - _Requirements: 11.2, 11.3, 11.5_

  - [ ] 8.3 Implement Configuration_Selector inside DeviceForm
    - Render a shadcn/ui `Select` with a "None" option (value `""`) and one `SelectItem` per configuration
    - Disable the control while `configurationsLoading` is true; show error text if `configurationsError` is set
    - Map selected value back to `configurationId: number | null` (empty string → null)
    - _Requirements: 12.2, 12.3, 12.5_

  - [ ]* 8.4 Write unit tests for DeviceForm
    - Assert create mode: title is "Add Device", `number` field is editable, all fields empty
    - Assert edit mode: title is "Edit Device", `number` field is disabled, fields pre-populated from `initialData`
    - Assert `getGroups` and `getConfigurations` called on mount
    - Assert submit button disabled while submitting; re-enabled on failure
    - Assert error message shown on submit failure; dialog stays open
    - Assert `onSuccess` and `onClose` called on successful submit
    - Assert `onClose` called on Cancel; no API call made
    - _Requirements: 7.3, 7.5, 7.6, 7.7, 8.3, 8.6, 8.7, 8.8, 9.9, 9.10_

  - [ ]* 8.5 Write property test for create mode submit (Property 14)
    - **Property 14: Create mode submit calls POST with form values and no id**
    - **Validates: Requirements 7.4**
    - Use `arbitraryDevicePayload()` to generate inputs; assert `createDevice` called with exact payload (no `id`); run 100 iterations

  - [ ]* 8.6 Write property test for edit mode submit (Property 15)
    - **Property 15: Edit mode submit calls PUT with form values and device id**
    - **Validates: Requirements 8.5**
    - Use `arbitraryDeviceView()` and updated form values; assert `updateDevice` called with payload including original `id`; run 100 iterations

  - [ ]* 8.7 Write property test for required field validation (Property 16)
    - **Property 16: Required field validation rejects empty or whitespace number**
    - **Validates: Requirements 10.1**
    - Use `arbitraryEmptyOrWhitespace()` for `number`; assert form rejects and no API call is made; run 100 iterations

  - [ ]* 8.8 Write property test for IMEI validation (Property 17)
    - **Property 17: IMEI validation accepts 15-digit strings and rejects all others**
    - **Validates: Requirements 10.2, 10.3**
    - Use `arbitraryValidImei()` and `arbitraryInvalidImei()`; assert valid IMEI passes and invalid IMEI fails; run 100 iterations each

  - [ ]* 8.9 Write property test for group selector pre-selection (Property 18)
    - **Property 18: Group selector pre-selects assigned groups for any device**
    - **Validates: Requirements 11.4**
    - Use `arbitraryDeviceView()` with non-empty `groups`; assert each group is checked in the Group_Selector; run 100 iterations

  - [ ]* 8.10 Write property test for configuration selector pre-selection (Property 19)
    - **Property 19: Configuration selector pre-selects assigned configuration for any device**
    - **Validates: Requirements 12.4**
    - Use `arbitraryDeviceView()` with non-null `configurationId`; assert matching configuration is selected; run 100 iterations

  - [ ]* 8.11 Write property test for null-safe form rendering (Property 30)
    - **Property 30: Null-safe rendering for any device with null optional fields**
    - **Validates: Requirements 21.6**
    - Use `arbitraryDeviceViewWithNulls()` with all optional fields null; assert opening `DeviceForm` in edit mode does not throw; run 100 iterations

- [ ] 9. Implement FilterPanel component
  - [ ] 9.1 Create `src/features/devices/FilterPanel.tsx`
    - Collapsible section collapsed by default; toggle button shows "More filters" / "Fewer filters"
    - Contains: group filter `Select`, configuration filter `Select`, status filter `Select` (All/Online/Offline), Android version `Input` (debounced 300ms)
    - Props: `filters: DeviceFilters`, `groups: LookupItem[]`, `configurations: ConfigurationOption[]`, `onChange: (filters: DeviceFilters) => void`
    - Show a badge count on the toggle button when any filter is active
    - _Requirements: 13.1–13.9_

  - [ ]* 9.2 Write property test for filter selection calls service with correct filter (Property 7)
    - **Property 7: Filter selection calls service with correct filter and resets page**
    - **Validates: Requirements 13.3, 13.4, 13.5, 13.6**
    - Use `fc.option(fc.integer({ min: 1 }))` for groupId/configurationId, `fc.option(fc.string())` for status/androidVersion; assert service called with correct filter and pageNum=1; run 100 iterations

  - [ ]* 9.3 Write property test for active filter indicator (Property 8)
    - **Property 8: Active filter indicator shown for any active filter**
    - **Validates: Requirements 13.9**
    - Use `fc.record` with at least one non-empty filter field; assert toggle button shows indicator; run 100 iterations

- [ ] 10. Implement BulkActionBar component
  - [ ] 10.1 Create `src/features/devices/BulkActionBar.tsx`
    - Props: `selectedCount: number`, `onDeleteSelected: () => void`, `onSetConfiguration: () => void`, `onSetGroup: () => void`
    - Display count of selected devices and three action buttons: "Delete Selected", "Set Configuration", "Set Group"
    - _Requirements: 15.5, 16.1, 17.1, 18.1_

  - [ ]* 10.2 Write property test for bulk delete calls DELETE for each selected device (Property 24)
    - **Property 24: Bulk delete calls DELETE for each selected device**
    - **Validates: Requirements 16.3**
    - Use `fc.set(fc.integer({ min: 1 }), { minLength: 1 })` for selected ids; confirm bulk delete and assert `deleteDevice` called once per id; run 100 iterations

  - [ ]* 10.3 Write property test for bulk set configuration calls PUT for each selected device (Property 25)
    - **Property 25: Bulk set configuration calls PUT for each selected device**
    - **Validates: Requirements 17.3**
    - Use `fc.set(fc.integer({ min: 1 }), { minLength: 1 })` for selected ids and `fc.integer({ min: 1 })` for configurationId; assert `updateDevice` called once per id with updated configurationId; run 100 iterations

  - [ ]* 10.4 Write property test for bulk set group calls PUT for each selected device (Property 26)
    - **Property 26: Bulk set group calls PUT for each selected device**
    - **Validates: Requirements 18.3**
    - Use `fc.set(fc.integer({ min: 1 }), { minLength: 1 })` for selected ids and `arbitraryLookupItem()` for group; assert `updateDevice` called once per id with updated groups; run 100 iterations

- [ ] 11. Implement QrDialog component
  - [ ] 11.1 Create `src/features/devices/QrDialog.tsx`
    - Render inside a shadcn/ui `Dialog`; props: `deviceId: number | null`, `onClose: () => void`
    - Call `deviceService.getDeviceQr(id)` on open; manage `qrData`, `loading`, `error` state
    - Show loading skeleton while fetching; render QR as `<img>` on success; show error message on failure
    - Provide a "Download" button that saves the QR code image to the user's device
    - _Requirements: 19.1–19.8_

  - [ ]* 11.2 Write property test for QR dialog fetches by device id (Property 27)
    - **Property 27: QR dialog fetches by device id**
    - **Validates: Requirements 19.3**
    - Use `fc.integer({ min: 1 })` for device ids; render QrDialog and assert `getDeviceQr` called with that id; run 100 iterations

- [ ] 12. Implement DevicesPage
  - [ ] 12.1 Create `src/features/devices/DevicesPage.tsx`
    - Fetch devices on mount and on state changes via `deviceService.getDevices()`; manage all state as specified in the design
    - Render shadcn/ui `Table` with default columns: Status (StatusBadge), Last Seen, Number, Configuration, Groups, Actions
    - Render optional columns (IMEI, Phone, Model, Battery Level, Android Version, Serial, Description) based on `columnVisibility` state
    - Show loading skeleton while fetching; show error banner with Retry button on failure; show empty-state message on empty list
    - Render "Add Device" button and Column_Visibility_Menu `DropdownMenu` in the page header
    - Render Search_Input above the table with 300ms debounce
    - Render FilterPanel below the search input; pass groups and configurations fetched in parallel on mount
    - Render Pagination_Controls below the table when `total > pageSize`; show current page and total count text
    - Render checkbox in leftmost column of every row and Select All checkbox in header
    - Show BulkActionBar above the table when `selectedIds.size > 0`
    - Each row has a `DropdownMenu` with: "View Details", "Edit", "Delete", and QR code icon button
    - _Requirements: 1.1–1.6, 2.1–2.5, 3.1–3.5, 4.1, 5.2, 6.1, 7.1, 8.1, 13.1, 14.1–14.5, 15.1–15.6_

  - [ ] 12.2 Wire DeviceDetailPanel into DevicesPage
    - Render `<DeviceDetailPanel deviceNumber={selectedDevice?.number ?? null} onClose={() => setSelectedDevice(null)} />`
    - On row click: set `selectedDevice = device`
    - _Requirements: 4.1, 4.6_

  - [ ] 12.3 Wire DeviceForm into DevicesPage
    - Render `<DeviceForm mode={formMode} initialData={deviceToEdit} onSuccess={refresh} onClose={() => setFormMode(null)} />` when `formMode` is non-null
    - "Add Device" button: set `formMode = 'create'`, `deviceToEdit = null`
    - "Edit" row action: set `deviceToEdit = device`, `formMode = 'edit'`
    - _Requirements: 7.1–7.7, 8.1–8.8_

  - [ ] 12.4 Wire DeleteDialog into DevicesPage
    - Render `<DeleteDialog device={deviceToDelete} onConfirm={handleDelete} onCancel={() => setDeviceToDelete(null)} />`
    - "Delete" row action: set `deviceToDelete = device`
    - On successful delete: close dialog, clear selection, refresh list
    - _Requirements: 6.1–6.6_

  - [ ] 12.5 Wire QrDialog into DevicesPage
    - Render `<QrDialog deviceId={deviceForQr?.id ?? null} onClose={() => setDeviceForQr(null)} />`
    - QR icon button in Actions column: set `deviceForQr = device`
    - _Requirements: 19.1–19.8_

  - [ ] 12.6 Implement bulk action dialogs in DevicesPage
    - "Delete Selected": open confirmation dialog showing count; on confirm call `deleteDevice` for each id in `selectedIds`; on success clear selection and refresh
    - "Set Configuration": open dialog with Configuration_Selector; on confirm call `updateDevice` for each selected device with new `configurationId`; on success clear selection and refresh
    - "Set Group": open dialog with Group_Selector; on confirm call `updateDevice` for each selected device with new `groups`; on success clear selection and refresh
    - Show partial failure errors listing devices that could not be updated/deleted
    - _Requirements: 16.1–16.6, 17.1–17.6, 18.1–18.6_

  - [ ]* 12.7 Write property test for every device row contains a status badge (Property 2)
    - **Property 2: Every device row contains a status badge**
    - **Validates: Requirements 5.2**
    - Use `fc.array(arbitraryDeviceView(), { minLength: 1 })` to generate lists; assert every rendered row contains a StatusBadge; run 100 iterations

  - [ ]* 12.8 Write property test for search resets page to 1 (Property 3)
    - **Property 3: Search resets page to 1 and includes search term**
    - **Validates: Requirements 2.3, 2.4, 3.4**
    - Use `fc.tuple(fc.integer({ min: 2 }), fc.string())` (page, searchTerm); assert applying search resets pageNum to 1 and includes value; run 100 iterations

  - [ ]* 12.9 Write property test for pagination controls appear iff total exceeds page size (Property 5)
    - **Property 5: Pagination controls appear iff total exceeds page size**
    - **Validates: Requirements 3.1, 3.5**
    - Use `fc.record({ totalItemsCount: fc.integer({ min: 0 }), pageSize: fc.integer({ min: 1 }) })`; assert pagination shown iff totalItemsCount > pageSize; run 100 iterations

  - [ ]* 12.10 Write property test for page navigation calls service with correct page number (Property 6)
    - **Property 6: Page navigation calls service with correct page number**
    - **Validates: Requirements 3.2**
    - Use `fc.integer({ min: 1 })` for page numbers; click page N and assert `getDevices` called with `pageNum: N`; run 100 iterations

  - [ ]* 12.11 Write property test for column toggle (Property 20)
    - **Property 20: Column toggle immediately shows or hides column**
    - **Validates: Requirements 14.4**
    - Use `fc.constantFrom('imei', 'phone', 'model', 'battery', 'android', 'serial', 'description')` for column names; toggle and assert column visibility changes; run 100 iterations

  - [ ]* 12.12 Write property test for Actions column always visible (Property 21)
    - **Property 21: Actions column is always visible**
    - **Validates: Requirements 14.5**
    - Use `fc.record` with any combination of column visibility booleans; assert Actions column always present; run 100 iterations

  - [ ]* 12.13 Write property test for Select All selects all rows (Property 22)
    - **Property 22: Select All selects all rows on current page**
    - **Validates: Requirements 15.3, 15.5**
    - Use `fc.array(arbitraryDeviceView(), { minLength: 1 })` for device lists; check Select All and assert all ids in selectedIds and BulkActionBar shown; run 100 iterations

  - [ ]* 12.14 Write property test for page or filter change clears selection (Property 23)
    - **Property 23: Page or filter change clears selection**
    - **Validates: Requirements 15.6**
    - With non-empty selection, change page or filter and assert selectedIds is empty; run 100 iterations

  - [ ]* 12.15 Write property test for null-safe row rendering (Property 30)
    - **Property 30: Null-safe rendering for any device with null optional fields**
    - **Validates: Requirements 21.6**
    - Use `arbitraryDeviceViewWithNulls()` to generate devices; assert rendering table row does not throw; run 100 iterations

- [ ] 13. Implement GroupsPage
  - [ ] 13.1 Create `src/features/groups/GroupsPage.tsx`
    - Fetch groups on mount via `groupService.getGroups()`; manage `groups`, `loading`, `error`, `formMode`, `groupToEdit`, `groupToDelete` state
    - Render shadcn/ui `Table` with columns: Group Name, Actions
    - Show loading skeleton while fetching; show error banner with Retry button on failure; show "No groups found" empty state
    - Render "Add Group" button in the page header
    - Inline group form dialog (shadcn/ui `Dialog`) with a required "Group Name" `Input` field; validate non-empty on submit
    - In create mode: call `groupService.createGroup(name)` on submit
    - In edit mode: pre-populate with group name; call `groupService.updateGroup(group)` on submit
    - Inline delete confirmation dialog (`AlertDialog`) showing group name; call `groupService.deleteGroup(id)` on confirm
    - After any successful create, update, or delete, re-fetch the groups list
    - _Requirements: 22.1–22.6, 23.1–23.6, 24.1–24.4, 25.1–25.5_

  - [ ]* 13.2 Write property test for groups table columns (Property 31)
    - **Property 31: Groups table columns rendered for any group list**
    - **Validates: Requirements 22.3**
    - Use `fc.array(arbitraryLookupItem(), { minLength: 1 })` to generate lists; assert rendered table has Group Name and Actions columns; run 100 iterations

  - [ ]* 13.3 Write property test for group name validation (Property 36)
    - **Property 36: Group name validation rejects empty or whitespace names**
    - **Validates: Requirements 23.6**
    - Use `arbitraryEmptyOrWhitespace()` for group name; assert form rejects and no API call is made; run 100 iterations

  - [ ]* 13.4 Write property test for group delete confirmation calls correct id (Property 37)
    - **Property 37: Group delete confirmation calls service with correct id**
    - **Validates: Requirements 25.2**
    - Use `arbitraryLookupItem()` to generate groups; confirm deletion and assert `deleteGroup` called with exact `id`; run 100 iterations

  - [ ]* 13.5 Write property test for cancel group delete makes no API call (Property 38)
    - **Property 38: Cancel group delete makes no API call**
    - **Validates: Requirements 25.5**
    - For any group, clicking Cancel must not invoke `groupService.deleteGroup`; run 100 iterations

- [ ] 14. Checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all unit and property tests pass; ask the user if questions arise.

- [ ] 15. Wire routing and navigation
  - [ ] 15.1 Update `src/app/App.tsx` — add routes for `/devices` and `/groups`
    - Add `<Route path="/devices" element={<DevicesPage />} />` (replacing any existing placeholder)
    - Add `<Route path="/groups" element={<GroupsPage />} />`
    - Import `DevicesPage` and `GroupsPage`
    - _Requirements: 1.1, 22.1_

  - [ ] 15.2 Update `src/features/layout/navItems.ts` — add or update Devices and Groups entries
    - Ensure Devices entry exists at `/devices` with appropriate icon
    - Add Groups entry at `/groups` with appropriate icon (e.g., `Layers` or `FolderOpen`) after Devices
    - _Requirements: 1.1, 22.1_

- [ ] 16. Final checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all tests pass end-to-end; ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Each task references specific requirements for traceability
- Property tests use fast-check with a minimum of 100 iterations each
- Each property test file must include a comment referencing the design property number and the requirements clause it validates
- Generators (`arbitraryLookupItem`, `arbitraryDeviceView`, `arbitraryDeviceViewWithNulls`, `arbitraryDevicePayload`, `arbitraryDeviceSearchRequest`, `arbitraryValidImei`, `arbitraryInvalidImei`, `arbitraryEmptyOrWhitespace`, `arbitraryConfigurationOption`) should be defined once in `src/features/devices/__tests__/generators.ts` and imported by all property test files
- `groups` is sent as `LookupItem[]` — the same shape stored on `DeviceView.groups`
- `configurationId` maps to a `Select` value; empty string selection maps to `null` in the payload
- Column visibility defaults: all optional columns hidden (`{ imei: false, phone: false, model: false, battery: false, android: false, serial: false, description: false }`)
- Bulk actions use sequential `Promise.allSettled` to collect partial failures rather than `Promise.all` which would abort on first failure
- Custom field labels (`custom1`, `custom2`, `custom3`) can be labeled "Custom 1", "Custom 2", "Custom 3" for MVP
- The `deviceService.getGroups` and `groupService.getGroups` both call `GET /rest/private/groups` — they can share the same underlying call or `deviceService.getGroups` can delegate to `groupService.getGroups`

---

## Additional Tasks (from gaps analysis)

- [ ] 17. Scaffold Tooltip component
  - Run `npx shadcn@latest add tooltip` inside `frontend/` to generate `src/shared/ui/tooltip.tsx`
  - _Requirements: 30.2_

- [ ] 18. Extend TypeScript types with all missing fields
  - [ ] 18.1 Extend `DeviceSearchRequest` in `src/features/devices/types.ts` with all sort/filter fields: `sortBy`, `sortDir`, `dateFrom`, `dateTo`, `onlineEarlierMillis`, `onlineLaterMillis`, `enrollmentDateFrom`, `enrollmentDateTo`, `mdmMode`, `kioskMode`, `launcherVersion`, `installationStatus`, `imeiChanged`, `fastSearch`
    - _Requirements: 37.1_
  - [ ] 18.2 Extend `DeviceInfoView` with: `permissions`, `applications`, `files`, `defaultLauncher`, `mdmMode`, `kioskMode`, `enrollTime`, `publicIp`, `launcherVersion`
    - _Requirements: 37.5_
  - [ ] 18.3 Add `BulkDeletePayload`, `GroupBulkPayload`, `AppSetting`, `DevicePermissions`, `DeviceApplication`, `DeviceFile` interfaces
    - _Requirements: 37.2, 37.3, 37.4_

- [ ] 19. Extend device service with all missing endpoints
  - [ ] 19.1 Add `deleteBulk(ids: number[]): Promise<void>` — calls `POST /rest/private/devices/deleteBulk`
    - _Requirements: 28.1, 36.1_
  - [ ] 19.2 Add `groupBulk(payload: GroupBulkPayload): Promise<void>` — calls `POST /rest/private/devices/groupBulk`
    - _Requirements: 29.1, 36.2_
  - [ ] 19.3 Add `getAppSettings(deviceId: number): Promise<AppSetting[]>` — calls `GET /rest/private/devices/{id}/applicationSettings`
    - _Requirements: 32.3, 36.3_
  - [ ] 19.4 Add `saveAppSettings(deviceId: number, settings: AppSetting[]): Promise<void>` — calls `POST /rest/private/devices/{id}/applicationSettings`
    - _Requirements: 32.7, 36.4_
  - [ ] 19.5 Add `notifyAppSettings(deviceId: number): Promise<void>` — calls `POST /rest/private/devices/{id}/applicationSettings/notify`
    - _Requirements: 32.8, 36.5_
  - [ ] 19.6 Add `updateDescription(deviceId: number, description: string): Promise<void>` — calls `POST /rest/private/devices/{id}/description`
    - _Requirements: 34.3, 36.6_
  - [ ] 19.7 Add `autocomplete(value: string): Promise<string[]>` — calls `POST /rest/private/devices/autocomplete`
    - _Requirements: 36.7_

  - [ ]* 19.8 Write property test for bulk delete uses dedicated endpoint (Property 39)
    - **Property 39: Bulk delete uses dedicated endpoint**
    - **Validates: Requirements 28.1**
    - Use `fc.set(fc.integer({ min: 1 }), { minLength: 1 })` for ids; confirm bulk delete and assert `deleteBulk` called with `{ ids }` — NOT individual `deleteDevice` calls; run 100 iterations

  - [ ]* 19.9 Write property test for bulk group uses dedicated endpoint (Property 40)
    - **Property 40: Bulk group uses dedicated endpoint**
    - **Validates: Requirements 29.1**
    - Use `fc.set(fc.integer({ min: 1 }), { minLength: 1 })` for ids and `arbitraryLookupItem()` for group; assert `groupBulk` called with correct `GroupBulkPayload`; run 100 iterations

  - [ ]* 19.10 Write property test for app settings save calls notify (Property 46)
    - **Property 46: App settings save calls notify after save**
    - **Validates: Requirements 32.8**
    - For any deviceId and settings array, assert `notifyAppSettings` is called with the same deviceId after `saveAppSettings` succeeds; run 100 iterations

  - [ ]* 19.11 Write property test for description-only edit endpoint (Property 47)
    - **Property 47: Description-only edit calls correct endpoint**
    - **Validates: Requirements 34.3, 36.6**
    - Use `fc.integer({ min: 1 })` for deviceId and `fc.string()` for description; assert `POST /rest/private/devices/{id}/description` is called; run 100 iterations

- [ ] 20. Update StatusBadge with full status code mapping and rich tooltip
  - [ ] 20.1 Update `src/features/devices/StatusBadge.tsx` to map all 5 status codes: `"green"` → Online (green), `"red"` → Offline (destructive), `"yellow"` → Warning (amber), `"brown"` → Inactive (orange), `"grey"` / null / other → Unknown (secondary)
    - _Requirements: 30.1_
  - [ ] 20.2 Wrap the badge in a shadcn/ui `Tooltip` that shows humanized elapsed time (e.g., "3 hours ago") on line 1 and exact formatted timestamp on line 2, derived from `lastUpdate` prop
    - _Requirements: 30.2_

  - [ ]* 20.3 Write property test for full status code mapping (Property 41)
    - **Property 41: Status badge maps all 5 status codes correctly**
    - **Validates: Requirements 30.1**
    - Use `fc.oneof(fc.constant("green"), fc.constant("red"), fc.constant("yellow"), fc.constant("brown"), fc.constant("grey"), fc.string(), fc.constant(null))` for statusCode; assert correct label and variant; run 100 iterations

  - [ ]* 20.4 Write property test for rich tooltip (Property 42)
    - **Property 42: Status badge tooltip shows humanized elapsed time**
    - **Validates: Requirements 30.2**
    - Use `fc.integer({ min: 0 })` for lastUpdate timestamps; assert tooltip contains humanized string and exact timestamp; run 100 iterations

- [ ] 21. Extend FilterPanel with full filter/sort set
  - [ ] 21.1 Update `src/features/devices/FilterPanel.tsx` to add all missing filter controls: launcher version input (debounced), enrollment date range pickers, online since range selectors, MDM mode select, kiosk mode select, installation status dropdown, IMEI changed checkbox, fast search checkbox, and sort controls (sortBy + sortDir toggle)
    - _Requirements: 2.7, 2.10_

  - [ ]* 21.2 Write property test for sort direction toggle (Property 50)
    - **Property 50: Sort direction toggles on repeated column header click**
    - **Validates: Requirements 2.10**
    - Click same column header twice; assert sortDir toggles asc→desc; click different column; assert sortDir resets to asc; run 100 iterations

- [ ] 22. Add explicit Search button
  - Update `src/features/devices/DevicesPage.tsx` to render a "Search" button alongside the Search_Input that triggers an immediate API call without waiting for the debounce
  - _Requirements: 2.12_

  - [ ]* 22.1 Write property test for explicit search button (Property 49)
    - **Property 49: Explicit search button triggers immediate fetch**
    - **Validates: Requirements 2.12**
    - For any search input value, clicking "Search" must trigger API call immediately; assert no debounce delay; run 100 iterations

- [ ] 23. Add dual pagination placement
  - Update `src/features/devices/DevicesPage.tsx` to render Pagination_Controls both above and below the Device_Table when `total > pageSize`
  - _Requirements: 3.1_

  - [ ]* 23.1 Write property test for dual pagination (Property 48)
    - **Property 48: Dual pagination rendered when total exceeds page size**
    - **Validates: Requirements 3.1**
    - Use `fc.record({ totalItemsCount: fc.integer({ min: 1 }), pageSize: fc.integer({ min: 1 }) })` where totalItemsCount > pageSize; assert Pagination_Controls rendered twice; run 100 iterations

- [ ] 24. Extend DeviceDetailPanel with all missing fields
  - Update `src/features/devices/DeviceDetailPanel.tsx` to display all fields from the extended `DeviceInfoView`: model, IMEI (with mismatch indicator), phone (with mismatch indicator), permissions status indicator, app installation status indicator, files status indicator, MDM mode, kiosk mode, launcher version (highlighted if mismatch), serial, enrollment date, public IP
  - _Requirements: 4.4_

- [ ] 25. Add computed status columns to DevicesPage
  - [ ] 25.1 Add "Permission Status" optional column: compute from `device.info.permissions` — green/amber/red icon with tooltip
    - _Requirements: 31.1_
  - [ ] 25.2 Add "Installation Status" optional column: compare `device.info.applications` vs configuration apps — green/amber/red icon with tooltip
    - _Requirements: 31.2_
  - [ ] 25.3 Add "Files Status" optional column: compare `device.info.files` vs configuration files — green/amber/red icon with tooltip
    - _Requirements: 31.3_
  - [ ] 25.4 Add all three computed columns to Column_Visibility_Menu (hidden by default)
    - _Requirements: 31.5_

- [ ] 26. Implement AppSettingsDialog
  - Create `src/features/devices/AppSettingsDialog.tsx`
  - On open: call `deviceService.getAppSettings(deviceId)`; show loading skeleton while fetching
  - Render list of `AppSetting` items with inline add/edit/delete capability
  - "Save" button: call `saveAppSettings` then `notifyAppSettings`; show error on failure; close on success
  - Add "App Settings" item to row action dropdown in DevicesPage
  - _Requirements: 32.1–32.10_

- [ ] 27. Implement auto-refresh
  - Update `src/features/devices/DevicesPage.tsx` to start a `setInterval` (60 000ms) on mount that re-fetches using current state without setting `loading = true`
  - Clear the interval on unmount
  - _Requirements: 33.1–33.5_

  - [ ]* 27.1 Write property test for auto-refresh no skeletons (Property 43)
    - **Property 43: Auto-refresh does not show skeletons**
    - **Validates: Requirements 33.4**
    - After initial load, advance fake timer by 60s; assert no Skeleton components rendered during background refresh; run 100 iterations

  - [ ]* 27.2 Write property test for auto-refresh interval cleared on unmount (Property 44)
    - **Property 44: Auto-refresh interval is cleared on unmount**
    - **Validates: Requirements 33.3**
    - Mount and unmount DevicesPage; advance fake timer; assert no additional API calls after unmount; run 100 iterations

- [ ] 28. Implement description-only edit path
  - Add a minimal description-edit `Dialog` in `DevicesPage` with a single `Textarea`
  - Show "Edit Description" in row dropdown when user has `edit_device_desc` permission but not `edit_devices`
  - On submit: call `deviceService.updateDescription(id, description)`
  - _Requirements: 34.1–34.5_

- [ ] 29. Add advanced form validation to DeviceForm
  - [ ] 29.1 Reject device number values containing `/`, `?`, or `&` with a validation error
    - _Requirements: 35.1_
  - [ ] 29.2 Show migration confirmation dialog when changing `number` of an already-enrolled device in edit mode
    - _Requirements: 35.2_
  - [ ] 29.3 Render custom field labels from application settings when available; fall back to "Custom 1/2/3"
    - _Requirements: 35.3_
  - [ ] 29.4 Render custom fields as `Textarea` when `customMultiline1/2/3` is true in settings
    - _Requirements: 35.4_

  - [ ]* 29.5 Write property test for device number rejects forbidden characters (Property 45)
    - **Property 45: Device number rejects forbidden characters**
    - **Validates: Requirements 35.1**
    - Use `fc.string().filter(s => /[/?&]/.test(s))` for device numbers; assert form rejects and no API call is made; run 100 iterations

- [ ] 30. Update bulk actions to use dedicated endpoints
  - [ ] 30.1 Update bulk delete in `DevicesPage` to call `deviceService.deleteBulk(ids)` instead of individual `deleteDevice` calls
    - _Requirements: 28.1_
  - [ ] 30.2 Update bulk set group in `DevicesPage` to call `deviceService.groupBulk({ ids, action: 'set', groups })` instead of individual `updateDevice` calls
    - _Requirements: 29.1_

- [ ] 31. Final checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all tests pass end-to-end; ask the user if questions arise.

## Notes (updated)

- `deleteBulk` and `groupBulk` use dedicated backend endpoints — do NOT use individual per-device calls for bulk operations
- `statusCode` has 5 distinct values: `"green"`, `"red"`, `"yellow"`, `"brown"`, `"grey"` — all must be mapped distinctly
- Auto-refresh must NOT set `loading = true` to avoid skeleton flash on every 60-second tick
- The `Tooltip` component from shadcn/ui is required for the rich status tooltip
- Computed status columns (Permission, Installation, Files) are computed client-side from `DeviceInfoView` — no additional API calls needed
- Description-only edit uses `POST /rest/private/devices/{id}/description` — a separate endpoint from the full `PUT /rest/private/devices`
- Device number validation must reject `/`, `?`, `&` characters in addition to the existing empty/whitespace check
