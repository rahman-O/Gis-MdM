# Implementation Plan: Devices Management

## Overview

Implement the `/devices` page with a paginated, searchable device table, slide-in detail panel, status badges, and delete confirmation dialog. All API calls go through a dedicated `deviceService` using the shared `apiClient`. UI is built exclusively from shadcn/ui components.

## Tasks

- [x] 1. Scaffold shadcn/ui components and shared hook
  - Run `npx shadcn@latest add badge table skeleton pagination alert-dialog dropdown-menu` from the `frontend/` directory to generate components into `src/shared/ui/`
  - Create `src/shared/hooks/useDebounce.ts` — export `useDebounce<T>(value: T, delay: number): T` using `useState` + `useEffect`
  - _Requirements: 2.2_
  - **Done:** Same primitives installed manually (`badge`, `table`, `skeleton`, `pagination`, `alert-dialog`, `dropdown-menu`) under `src/shared/ui/` + Radix deps; `useDebounce` implemented.

- [x] 2. Define device types
  - [x] 2.1 Create `src/features/devices/types.ts` with all interfaces
    - `DeviceView`, `DeviceInfoView`, `LookupItem`, `ConfigurationView`, `DeviceSearchRequest` (pageNum 1-based, pageSize, value?), `DeviceListResponse` (devices.items, devices.totalItemsCount, configurations map)
    - `statusCode` is `string | null` — values are `"green"` / `"red"` / other
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5_

  - [ ]* 2.2 Write property test for null-safe field access (Property 12)
    - **Property 12: Null-safe field access**
    - **Validates: Requirements 8.5**
    - Use `fc.record` with nullable optional fields; assert no runtime error when rendering a row with all nulls

- [x] 3. Implement device service
  - [x] 3.1 Create `src/features/devices/deviceService.ts`
    - `getDevices(params: DeviceSearchRequest): Promise<DeviceListResponse>` — POST `/rest/private/devices/search`
    - `getDevice(number: string): Promise<DeviceView>` — GET `/rest/private/devices/number/{number}`
    - `updateDevice(device: DeviceView): Promise<DeviceView>` — PUT `/rest/private/devices`
    - `deleteDevice(id: number): Promise<void>` — DELETE `/rest/private/devices/{id}`
    - Unwrap `response.data.data` from the backend envelope; re-throw on non-2xx
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_

  - [x]* 3.2 Write unit tests for deviceService
    - Mock `apiClient`; assert correct method, URL, and body for each function
    - Assert `getDevices` sends POST with `{ pageNum, pageSize, value }` body
    - Assert `getDevice` calls GET with device number in path
    - Assert `deleteDevice` calls DELETE with numeric id in path
    - _Requirements: 7.1, 7.2, 7.4_

  - [ ]* 3.3 Write property test for service error propagation (Property 11)
    - **Property 11: Service error propagation**
    - **Validates: Requirements 7.6**
    - For any non-2xx status code, assert each service function rejects with the error

- [x] 4. Implement StatusBadge component
  - [x] 4.1 Create `src/features/devices/StatusBadge.tsx`
    - Accept `statusCode: string | null`; map `"green"` → `default` variant / "Online", `"red"` → `destructive` / "Offline", anything else → `secondary` / "Unknown"
    - _Requirements: 5.1, 5.4_

  - [x]* 4.2 Write property test for status badge mapping (Property 1)
    - **Property 1: Status badge mapping is total and correct**
    - **Validates: Requirements 5.1, 5.4**
    - Use `fc.oneof(fc.constant("green"), fc.constant("red"), fc.string())` — assert every input renders exactly one labeled badge

- [x] 5. Implement DeleteDialog component
  - [x] 5.1 Create `src/features/devices/DeleteDialog.tsx`
    - Wrap `AlertDialog`; props: `device: DeviceView | null`, `onConfirm(): Promise<void>`, `onCancel(): void`
    - Manage `deleting` and `deleteError` state internally
    - Disable confirm button and show spinner while `deleting`; display error and keep dialog open on failure; close on success
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

  - [x]* 5.2 Write unit tests for DeleteDialog
    - Assert confirm button is disabled while deleting
    - Assert error message shown on failure, dialog stays open
    - Assert `onConfirm` called with correct device id on confirm
    - Assert no API call on cancel
    - _Requirements: 6.3, 6.5, 6.6_

  - [ ]* 5.3 Write property test for cancel makes no API call (Property 10)
    - **Property 10: Cancel delete makes no API call**
    - **Validates: Requirements 6.6**
    - For any arbitrary `DeviceView`, clicking Cancel must not invoke `deleteDevice`

  - [ ]* 5.4 Write property test for delete calls service with correct ID (Property 9)
    - **Property 9: Delete confirmation calls service with correct ID**
    - **Validates: Requirements 6.2**
    - For any arbitrary `DeviceView`, confirming deletion must call `deleteDevice` with `device.id`

- [x] 6. Implement DeviceDetailPanel component
  - [x] 6.1 Create `src/features/devices/DeviceDetailPanel.tsx`
    - Render inside `Sheet` (side="right"); props: `deviceNumber: string | null`, `onClose(): void`
    - Fetch via `deviceService.getDevice(deviceNumber)` when `deviceNumber` changes
    - Show `Skeleton` while loading; show error message on failure
    - Display: device number, status badge, configuration name, group names, battery level (or "N/A"), last update (formatted date/time), GPS coordinates or "Location unavailable"
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_

  - [ ]* 6.2 Write unit tests for DeviceDetailPanel
    - Assert skeleton shown while loading
    - Assert error message shown on fetch failure
    - Assert "Location unavailable" when `info.latitude` is null
    - Assert all required fields rendered on success
    - _Requirements: 4.3, 4.4, 4.5_

  - [ ]* 6.3 Write property test for detail panel renders all required fields (Property 8)
    - **Property 8: Detail panel renders all required fields**
    - **Validates: Requirements 4.4**
    - For any arbitrary `DeviceView` with non-null detail, assert all required fields are present in the rendered output

  - [ ]* 6.4 Write property test for detail panel fetches by device number (Property 7)
    - **Property 7: Detail panel fetches by device number**
    - **Validates: Requirements 4.2**
    - For any arbitrary `DeviceView`, assert `getDevice` is called with `device.number` (not `device.id`)

- [x] 7. Implement DevicesPage component
  - [x] 7.1 Create `src/features/devices/DevicesPage.tsx` with state and data fetching
    - State: `devices`, `configurations`, `total`, `page` (1-based), `pageSize` (20), `search`, `debouncedSearch`, `loading`, `error`, `selectedDevice`, `deviceToDelete`
    - Use `useDebounce(search, 300)` for `debouncedSearch`
    - `useEffect` on `[debouncedSearch, page]` → call `deviceService.getDevices`; reset `page` to 1 when `debouncedSearch` changes
    - _Requirements: 1.1, 1.2, 1.4, 1.5, 1.6, 2.1, 2.2, 2.3, 2.4, 3.3_

  - [x] 7.2 Add table, search input, and row actions to DevicesPage
    - Render shadcn/ui `Table` with columns: Device Number, Status, Configuration, Groups, Last Seen, Actions
    - Render `StatusBadge` in Status column for every row
    - Render `DropdownMenu` in Actions column with "View Details" and "Delete" items
    - Show `Skeleton` rows while `loading`; show error banner with "Retry" button on `error`; show empty-state message when list is empty or search returns no results
    - _Requirements: 1.3, 1.4, 1.5, 1.6, 2.5, 5.2_

  - [x] 7.3 Add pagination controls to DevicesPage
    - Render shadcn/ui `Pagination` below the table when `total > pageSize`
    - Display current page number and total device count as text
    - On page click call `setPage(n)` which triggers re-fetch
    - _Requirements: 3.1, 3.2, 3.3, 3.5_

  - [x] 7.4 Wire DeviceDetailPanel and DeleteDialog into DevicesPage
    - Pass `selectedDevice?.number` to `DeviceDetailPanel`; on close set `selectedDevice` to null and return focus to the previously selected row
    - Pass `deviceToDelete` to `DeleteDialog`; on confirm call `deviceService.deleteDevice` then refresh list; on cancel set `deviceToDelete` to null
    - _Requirements: 4.1, 4.6, 6.1, 6.4_

  - [x]* 7.5 Write unit tests for DevicesPage
    - Assert skeleton shown while loading
    - Assert error banner with retry shown on failure
    - Assert empty-state message when `devices` is empty
    - Assert "no devices found for '{term}'" when search returns empty
    - _Requirements: 1.4, 1.5, 1.6, 2.5_

  - [ ]* 7.6 Write property test for every row has a status badge (Property 2)
    - **Property 2: Every device row contains a status badge**
    - **Validates: Requirements 5.2**
    - For any non-empty `fc.array(arbitraryDeviceView(), { minLength: 1 })`, assert each rendered row contains a `StatusBadge`

  - [ ]* 7.7 Write property test for search resets page to 1 (Property 3)
    - **Property 3: Search resets page to 1**
    - **Validates: Requirements 2.3, 2.4, 3.4**
    - For any `(currentPage, newSearchTerm)`, assert `page` is reset to 1 before the API call is issued

  - [ ]* 7.8 Write property test for debounce collapses rapid inputs (Property 4)
    - **Property 4: Debounce collapses rapid inputs into one call**
    - **Validates: Requirements 2.2**
    - For any sequence of keystrokes within 300ms, assert `getDevices` is called at most once after the debounce period

  - [ ]* 7.9 Write property test for pagination controls visibility (Property 5)
    - **Property 5: Pagination controls appear iff total exceeds page size**
    - **Validates: Requirements 3.1, 3.5**
    - For any `(totalItemsCount, pageSize)`, assert pagination is rendered iff `totalItemsCount > pageSize`

  - [ ]* 7.10 Write property test for page navigation calls service with correct page (Property 6)
    - **Property 6: Page navigation calls service with correct page number**
    - **Validates: Requirements 3.2**
    - For any page number N clicked, assert `getDevices` is called with `pageNum: N`

- [x] 8. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- `statusCode` is a string from the backend (`"green"`, `"red"`, etc.) — not a number
- Pagination is 1-based (`pageNum`) to match the backend `POST /rest/private/devices/search` contract
- `getDevice` looks up by device `number` (string), not numeric `id`
- Property tests use fast-check; each test is tagged with its property number and requirement clause
