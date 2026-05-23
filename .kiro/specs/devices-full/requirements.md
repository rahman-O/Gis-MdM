# Requirements Document: Devices Full

## Introduction

This document is the single authoritative specification for all device-related functionality in the HMDM Modern Architecture frontend. It supersedes and replaces both `devices-management` and `devices-add-edit` specs. It covers the complete `/devices` page: paginated device list, search, advanced filtering with sorting, bulk actions (using dedicated bulk endpoints), add/edit form, device detail panel, status indicators, delete, QR code enrollment, configurable table columns, application settings per device, auto-refresh, and a Groups management section. All data is fetched from the HMDM backend REST API (base `/rest`) using the shared `apiClient` with automatic `X-Auth-Token` injection. The UI is built with React + Vite + TypeScript + shadcn/ui + Tailwind CSS.

---

## Glossary

- **Devices_Page**: The React page component rendered at the `/devices` route.
- **Device_Service**: The frontend service module (`src/features/devices/deviceService.ts`) wrapping all device-related `apiClient` calls.
- **Device_Table**: The shadcn/ui `Table` component rendering the paginated list of devices.
- **Device_Form**: The shadcn/ui `Dialog`-based form component used for both creating and editing devices.
- **Detail_Panel**: The shadcn/ui `Sheet` (side panel) displaying full device information for a selected device.
- **Delete_Dialog**: The shadcn/ui `AlertDialog` confirmation modal shown before a device is permanently deleted.
- **Status_Badge**: A colored `Badge` component communicating a device's connectivity state derived from `statusCode`.
- **Search_Input**: The `Input` component used to filter devices by device number or description.
- **Pagination_Controls**: The shadcn/ui `Pagination` component for navigating between pages — rendered both above and below the Device_Table.
- **Filter_Panel**: The collapsible section containing all advanced search/filter/sort options.
- **Bulk_Action_Bar**: The toolbar that appears when one or more device rows are selected, exposing bulk operations.
- **Group_Selector**: The multi-checkbox `Popover` inside `Device_Form` for assigning groups to a device.
- **Configuration_Selector**: The single `Select` inside `Device_Form` for assigning a configuration to a device.
- **Column_Visibility_Menu**: The dropdown that lets users toggle optional table columns on or off.
- **QR_Dialog**: The `Dialog` displaying a QR code for device enrollment.
- **App_Settings_Dialog**: The `Dialog` for managing per-device application settings.
- **Groups_Page**: The React page component rendered at the `/groups` route for managing device groups.
- **Group_Service**: The frontend service module (`src/features/groups/groupService.ts`) wrapping all group-related `apiClient` calls.
- **DevicePayload**: The request body sent to `POST /rest/private/devices` (create) or `PUT /rest/private/devices` (update).
- **DeviceSearchRequest**: The POST body sent to `POST /rest/private/devices/search` — includes pageNum, pageSize, value, groupId, configurationId, sortBy, sortDir, dateFrom, dateTo, onlineEarlierMillis, onlineLaterMillis, enrollmentDateFrom, enrollmentDateTo, mdmMode, kioskMode, androidVersion, launcherVersion, installationStatus, imeiChanged, fastSearch.
- **LookupItem**: A `{ id: number, name: string | null }` pair used for groups.
- **AppSetting**: A per-device application setting item with fields: `id`, `applicationId`, `name`, `value`, `comment`, `readonly`, `lastUpdate`.
- **Auto_Refresh**: The mechanism that re-fetches the device list every 60 seconds in the background without showing loading skeletons.- **Custom_Field**: One of `custom1`, `custom2`, or `custom3` — optional string device properties.


---

## Requirements

### Requirement 1: Devices List Page

**User Story:** As an MDM administrator, I want to see all enrolled devices in a paginated table, so that I can quickly survey the fleet.

#### Acceptance Criteria

1. WHEN a user navigates to `/devices`, THE Devices_Page SHALL render a full-page layout consistent with the existing navigation shell.
2. WHEN the Devices_Page mounts, THE Device_Service SHALL call `POST /rest/private/devices/search` with `{ pageNum: 1, pageSize: 20 }`.
3. WHEN the API response is received, THE Device_Table SHALL display the default visible columns: Status, Last Seen, Number, Configuration, Groups, Actions.
4. WHILE the API request is in flight, THE Devices_Page SHALL display a loading skeleton in place of the table rows.
5. IF the API request fails, THEN THE Devices_Page SHALL display an error message with a "Retry" action.
6. WHEN the device list is empty, THE Devices_Page SHALL display an empty-state message indicating no devices are enrolled.

---

### Requirement 2: Search and Advanced Filters

**User Story:** As an MDM administrator, I want to filter and sort the device list by multiple criteria, so that I can quickly locate specific devices in a large fleet.

#### Acceptance Criteria

1. THE Devices_Page SHALL render a Search_Input above the Device_Table.
2. WHEN a user types in the Search_Input, THE Devices_Page SHALL debounce the input by 300 ms before issuing a new API request.
3. WHEN a debounced search term is applied, THE Device_Service SHALL call `POST /rest/private/devices/search` with the `value` field set to the search term and `pageNum` reset to 1.
4. WHEN the search term is cleared, THE Device_Service SHALL call `POST /rest/private/devices/search` with no `value` field and `pageNum` reset to 1.
5. IF a search returns no results, THEN THE Devices_Page SHALL display a "no devices found" message that includes the search term.
6. THE Devices_Page SHALL render a Filter_Panel that is collapsed by default and expandable via a "More filters" toggle.
7. WHEN the Filter_Panel is expanded, THE Devices_Page SHALL display the following filter controls:
   - Group filter dropdown (populated from `GET /rest/private/groups`)
   - Configuration filter dropdown (populated from `GET /rest/private/configurations`)
   - Status filter dropdown (All / Online / Offline)
   - Android version text input (debounced 300ms)
   - Launcher version text input (debounced 300ms)
   - Enrollment date range (date-from / date-to)
   - Online since filter (onlineEarlierMillis / onlineLaterMillis)
   - MDM mode filter (All / Yes / No)
   - Kiosk mode filter (All / Yes / No)
   - Installation status filter dropdown
   - IMEI changed checkbox
   - Fast search checkbox
8. WHEN any filter value changes, THE Device_Service SHALL call `POST /rest/private/devices/search` with the updated filter fields and `pageNum` reset to 1.
9. WHEN a user clears all filters, THE Device_Service SHALL call `POST /rest/private/devices/search` with no filter fields and `pageNum` reset to 1.
10. THE Devices_Page SHALL support column sorting: WHEN a user clicks a sortable column header, THE Device_Service SHALL call `POST /rest/private/devices/search` with the `sortBy` field set to the column identifier and `sortDir` toggled between `asc` and `desc`.
11. WHEN any filter is active, THE Devices_Page SHALL display a visual indicator (badge count or highlighted toggle) on the Filter_Panel toggle button.
12. THE Devices_Page SHALL render an explicit "Search" button alongside the Search_Input that triggers an immediate search without waiting for the debounce period.

---

### Requirement 3: Pagination

**User Story:** As an MDM administrator, I want to navigate between pages of devices, so that I can browse large fleets without loading all records at once.

#### Acceptance Criteria

1. WHEN the API response `totalItemsCount` is greater than `pageSize`, THE Devices_Page SHALL render Pagination_Controls both above and below the Device_Table.
2. WHEN a user clicks a page number in the Pagination_Controls, THE Device_Service SHALL call `POST /rest/private/devices/search` with the updated `pageNum` (1-based).
3. WHEN a new page is requested, THE Devices_Page SHALL display a loading indicator while the request is in flight.
4. WHEN the search term or any active filter changes, THE Pagination_Controls SHALL reset to page 1.
5. THE Devices_Page SHALL display the current page range (e.g., "1–20 / 150") alongside the Pagination_Controls.


---

### Requirement 4: Device Detail Panel

**User Story:** As an MDM administrator, I want to inspect the full details of a device, so that I can diagnose issues and review device health.

#### Acceptance Criteria

1. WHEN a user clicks a device row in the Device_Table, THE Devices_Page SHALL open the Detail_Panel for that device.
2. WHEN the Detail_Panel opens, THE Device_Service SHALL call `GET /rest/private/devices/number/{number}` to fetch the full device record.
3. WHILE the detail request is in flight, THE Detail_Panel SHALL display a loading skeleton.
4. WHEN the detail response is received, THE Detail_Panel SHALL display: device number, Status_Badge with rich tooltip (humanized "last seen ago" + exact timestamp), configuration name (as a clickable link when the user has configuration permission), group names, battery level (percentage or "N/A"), model, last update timestamp (formatted as human-readable date/time), GPS coordinates or "Location unavailable" if absent, IMEI (with mismatch indicator if server value differs from device-reported value), phone number (with mismatch indicator), permissions status indicator, app installation status indicator, files status indicator, MDM mode, kiosk mode, Android version, launcher version (highlighted if it differs from the required configuration version), serial number, and enrollment date.
5. IF the detail request fails, THEN THE Detail_Panel SHALL display an error message inside the panel.
6. WHEN a user closes the Detail_Panel, THE Devices_Page SHALL return focus to the previously selected table row.

---

### Requirement 5: Status Indicators

**User Story:** As an MDM administrator, I want to see at a glance whether each device is online, offline, or in an unknown state, so that I can prioritize attention to unreachable devices.

#### Acceptance Criteria

1. THE Status_Badge SHALL map `statusCode` string values as follows: `"green"` SHALL render a green badge labeled "Online"; `"red"` SHALL render a red badge labeled "Offline"; any other value (including null) SHALL render a gray badge labeled "Unknown".
2. THE Status_Badge SHALL be rendered in the Status column of every device row in the Device_Table.
3. THE Status_Badge SHALL be rendered in the Detail_Panel header for the selected device.
4. WHEN the status is "Online", THE Status_Badge SHALL use the `default` or success variant; WHEN "Offline", the `destructive` variant; WHEN "Unknown", the `secondary` variant.

---

### Requirement 6: Delete Device

**User Story:** As an MDM administrator, I want to remove a device from the MDM system, so that I can keep the device inventory accurate.

#### Acceptance Criteria

1. WHEN a user activates the delete action for a device via the row action menu, THE Devices_Page SHALL open the Delete_Dialog identifying the device by number.
2. WHEN a user confirms deletion in the Delete_Dialog, THE Device_Service SHALL call `DELETE /rest/private/devices/{id}`.
3. WHILE the delete request is in flight, THE Delete_Dialog SHALL disable the confirm button and show a loading indicator.
4. WHEN the delete request succeeds, THE Devices_Page SHALL close the Delete_Dialog and refresh the device list.
5. IF the delete request fails, THEN THE Delete_Dialog SHALL display an error message and keep the dialog open.
6. WHEN a user cancels the Delete_Dialog, THE Devices_Page SHALL close the dialog without making any API call.


---

### Requirement 7: Add Device

**User Story:** As an MDM administrator, I want an "Add Device" button on the devices page, so that I can enroll a new device without leaving the page.

#### Acceptance Criteria

1. THE Devices_Page SHALL render an "Add Device" button in the page header area alongside the Search_Input.
2. WHEN a user clicks the "Add Device" button, THE Devices_Page SHALL open the Device_Form dialog in create mode with all fields empty.
3. WHEN the Device_Form is open in create mode, THE Device_Form SHALL display the title "Add Device".
4. WHEN a user submits the Device_Form in create mode with valid data, THE Device_Service SHALL call `POST /rest/private/devices` with a DevicePayload containing the form values and no `id` field.
5. WHILE the create request is in flight, THE Device_Form SHALL disable the submit button and display a loading indicator.
6. WHEN the create request succeeds, THE Device_Form SHALL close the dialog and THE Devices_Page SHALL refresh the device list.
7. IF the create request fails, THEN THE Device_Form SHALL display an error message inside the dialog and keep the dialog open with the submit button re-enabled.

---

### Requirement 8: Edit Device

**User Story:** As an MDM administrator, I want to edit an existing device's details from the device list, so that I can update its configuration, groups, or metadata.

#### Acceptance Criteria

1. WHEN a user opens the row action dropdown for a device, THE Devices_Page SHALL include an "Edit" menu item.
2. WHEN a user clicks the "Edit" menu item, THE Devices_Page SHALL open the Device_Form dialog in edit mode pre-populated with that device's current field values.
3. WHEN the Device_Form is open in edit mode, THE Device_Form SHALL display the title "Edit Device".
4. WHEN the Device_Form is open in edit mode, THE Device_Form SHALL render the `number` field as read-only (disabled).
5. WHEN a user submits the Device_Form in edit mode with valid data, THE Device_Service SHALL call `PUT /rest/private/devices` with a DevicePayload containing the updated form values and the device's existing `id`.
6. WHILE the update request is in flight, THE Device_Form SHALL disable the submit button and display a loading indicator.
7. WHEN the update request succeeds, THE Device_Form SHALL close the dialog and THE Devices_Page SHALL refresh the device list.
8. IF the update request fails, THEN THE Device_Form SHALL display an error message inside the dialog and keep the dialog open with the submit button re-enabled.

---

### Requirement 9: Device Form Fields

**User Story:** As an MDM administrator, I want to fill in all relevant device properties in a single form, so that I can fully configure a device at enrollment or update time.

#### Acceptance Criteria

1. THE Device_Form SHALL render a `number` field (labeled "Device ID") as a required text input.
2. THE Device_Form SHALL render a `description` field as an optional textarea.
3. THE Device_Form SHALL render a `configurationId` field as an optional Configuration_Selector populated from `GET /rest/private/configurations`.
4. THE Device_Form SHALL render a `groups` field as an optional Group_Selector populated from `GET /rest/private/groups`.
5. THE Device_Form SHALL render an `imei` field as an optional text input.
6. THE Device_Form SHALL render a `phone` field as an optional text input.
7. THE Device_Form SHALL render `custom1`, `custom2`, and `custom3` fields as optional text inputs.
8. WHEN the Device_Form opens, THE Device_Form SHALL fetch groups and configurations in parallel so that both selectors are populated before the user interacts with them.
9. THE Device_Form SHALL render a "Cancel" button that closes the dialog without making any API call.
10. WHEN a user closes the Device_Form (via Cancel or the dialog close control), THE Devices_Page SHALL NOT refresh the device list unless a successful save occurred.


---

### Requirement 10: Form Validation

**User Story:** As an MDM administrator, I want the form to catch invalid input before submission, so that I don't accidentally create a device with missing or malformed data.

#### Acceptance Criteria

1. WHEN a user submits the Device_Form with an empty or whitespace-only `number` field, THE Device_Form SHALL display a validation error adjacent to the `number` field and SHALL NOT call any API endpoint.
2. WHEN a user submits the Device_Form with an `imei` value that is not exactly 15 digits, THE Device_Form SHALL display a validation error adjacent to the `imei` field and SHALL NOT call any API endpoint.
3. WHEN a user submits the Device_Form with an empty `imei` field, THE Device_Form SHALL treat the field as absent and SHALL NOT apply the 15-digit validation rule.
4. WHEN all required fields are valid, THE Device_Form SHALL enable the submit button and allow submission.

---

### Requirement 11: Group Selector

**User Story:** As an MDM administrator, I want to assign a device to one or more groups, so that group-based policies apply to it.

#### Acceptance Criteria

1. WHEN the Device_Form opens, THE Group_Selector SHALL fetch available groups from `GET /rest/private/groups`.
2. WHILE groups are loading, THE Group_Selector SHALL display a loading state and SHALL be disabled.
3. THE Group_Selector SHALL allow the user to select zero or more groups from the fetched list using a multi-checkbox popover.
4. WHEN the Device_Form is in edit mode, THE Group_Selector SHALL pre-select the groups already assigned to the device.
5. IF the groups fetch fails, THEN THE Group_Selector SHALL display an error message and SHALL remain disabled.

---

### Requirement 12: Configuration Selector

**User Story:** As an MDM administrator, I want to assign a configuration policy to a device, so that the correct MDM policy is applied.

#### Acceptance Criteria

1. WHEN the Device_Form opens, THE Configuration_Selector SHALL fetch available configurations from `GET /rest/private/configurations`.
2. WHILE configurations are loading, THE Configuration_Selector SHALL display a loading state and SHALL be disabled.
3. THE Configuration_Selector SHALL allow the user to select at most one configuration, or leave the field empty (None).
4. WHEN the Device_Form is in edit mode, THE Configuration_Selector SHALL pre-select the configuration matching the device's `configurationId`.
5. IF the configurations fetch fails, THEN THE Configuration_Selector SHALL display an error message and SHALL remain disabled.


---

### Requirement 13: Advanced Filters

**User Story:** As an MDM administrator, I want to filter the device list by group, configuration, status, and Android version, so that I can quickly narrow down devices matching specific criteria.

#### Acceptance Criteria

1. THE Devices_Page SHALL render a Filter_Panel that is collapsed by default and expandable via a "More filters" toggle.
2. WHEN the Filter_Panel is expanded, THE Devices_Page SHALL display: a group filter dropdown, a configuration filter dropdown, a status filter dropdown (All / Online / Offline), and an Android version text input.
3. WHEN a user selects a group in the group filter, THE Device_Service SHALL call `POST /rest/private/devices/search` with the selected `groupId` and `pageNum` reset to 1.
4. WHEN a user selects a configuration in the configuration filter, THE Device_Service SHALL call `POST /rest/private/devices/search` with the selected `configurationId` and `pageNum` reset to 1.
5. WHEN a user selects a status in the status filter, THE Device_Service SHALL call `POST /rest/private/devices/search` with the corresponding `status` value and `pageNum` reset to 1.
6. WHEN a user types in the Android version input, THE Devices_Page SHALL debounce the input by 300 ms before calling `POST /rest/private/devices/search` with the `androidVersion` value and `pageNum` reset to 1.
7. WHEN a user clears all filters, THE Device_Service SHALL call `POST /rest/private/devices/search` with no filter fields and `pageNum` reset to 1.
8. THE Devices_Page SHALL populate the group filter dropdown from `GET /rest/private/groups` and the configuration filter dropdown from `GET /rest/private/configurations`.
9. WHEN any filter is active, THE Devices_Page SHALL display a visual indicator (e.g., badge count or highlighted toggle) showing that filters are applied.

---

### Requirement 14: Configurable Table Columns

**User Story:** As an MDM administrator, I want to choose which device properties are visible in the table, so that I can focus on the information most relevant to my workflow.

#### Acceptance Criteria

1. THE Devices_Page SHALL render a Column_Visibility_Menu button that opens a dropdown listing all available columns.
2. THE Device_Table SHALL display the following columns by default: Status, Last Seen, Number, Configuration, Groups, Actions.
3. THE Device_Table SHALL support the following optional columns that are hidden by default: IMEI, Phone, Model, Battery Level, Android Version, Serial, Description.
4. WHEN a user toggles a column in the Column_Visibility_Menu, THE Device_Table SHALL immediately show or hide that column without a page reload.
5. THE Device_Table SHALL always display the Actions column and SHALL NOT allow it to be hidden.


---

### Requirement 15: Row Selection

**User Story:** As an MDM administrator, I want to select multiple devices at once, so that I can perform bulk operations efficiently.

#### Acceptance Criteria

1. THE Device_Table SHALL render a checkbox in the leftmost column of every device row.
2. THE Device_Table SHALL render a "Select All" checkbox in the table header that selects or deselects all currently visible rows.
3. WHEN a user checks the "Select All" checkbox, THE Device_Table SHALL mark all rows on the current page as selected.
4. WHEN a user unchecks the "Select All" checkbox, THE Device_Table SHALL deselect all rows.
5. WHEN at least one row is selected, THE Devices_Page SHALL display the Bulk_Action_Bar above the Device_Table showing the count of selected devices.
6. WHEN the page changes or a new search or filter is applied, THE Devices_Page SHALL clear all row selections.

---

### Requirement 16: Bulk Delete

**User Story:** As an MDM administrator, I want to delete multiple devices at once, so that I can clean up the device inventory efficiently.

#### Acceptance Criteria

1. WHEN at least one device row is selected, THE Bulk_Action_Bar SHALL display a "Delete Selected" button.
2. WHEN a user clicks "Delete Selected", THE Devices_Page SHALL open a confirmation dialog identifying the number of devices to be deleted.
3. WHEN a user confirms bulk deletion, THE Device_Service SHALL call `DELETE /rest/private/devices/{id}` for each selected device.
4. WHILE bulk deletion is in progress, THE Devices_Page SHALL disable the confirm button and show a progress indicator.
5. WHEN all delete requests succeed, THE Devices_Page SHALL close the confirmation dialog, clear the selection, and refresh the device list.
6. IF any delete request fails, THEN THE Devices_Page SHALL display an error message listing the devices that could not be deleted and keep the dialog open.

---

### Requirement 17: Bulk Set Configuration

**User Story:** As an MDM administrator, I want to assign a configuration to multiple devices at once, so that I can apply policies to a fleet segment efficiently.

#### Acceptance Criteria

1. WHEN at least one device row is selected, THE Bulk_Action_Bar SHALL display a "Set Configuration" button.
2. WHEN a user clicks "Set Configuration", THE Devices_Page SHALL open a dialog containing a Configuration_Selector populated from `GET /rest/private/configurations`.
3. WHEN a user selects a configuration and confirms, THE Device_Service SHALL call `PUT /rest/private/devices` for each selected device with the updated `configurationId`.
4. WHILE bulk update is in progress, THE Devices_Page SHALL disable the confirm button and show a progress indicator.
5. WHEN all update requests succeed, THE Devices_Page SHALL close the dialog, clear the selection, and refresh the device list.
6. IF any update request fails, THEN THE Devices_Page SHALL display an error message listing the devices that could not be updated.

---

### Requirement 18: Bulk Set Group

**User Story:** As an MDM administrator, I want to assign a group to multiple devices at once, so that I can organize fleet segments efficiently.

#### Acceptance Criteria

1. WHEN at least one device row is selected, THE Bulk_Action_Bar SHALL display a "Set Group" button.
2. WHEN a user clicks "Set Group", THE Devices_Page SHALL open a dialog containing a Group_Selector populated from `GET /rest/private/groups`.
3. WHEN a user selects a group and confirms, THE Device_Service SHALL call `PUT /rest/private/devices` for each selected device with the updated `groups` array.
4. WHILE bulk update is in progress, THE Devices_Page SHALL disable the confirm button and show a progress indicator.
5. WHEN all update requests succeed, THE Devices_Page SHALL close the dialog, clear the selection, and refresh the device list.
6. IF any update request fails, THEN THE Devices_Page SHALL display an error message listing the devices that could not be updated.


---

### Requirement 19: QR Code Enrollment

**User Story:** As an MDM administrator, I want to generate a QR code for a device, so that I can enroll it by scanning without manual configuration.

#### Acceptance Criteria

1. THE Device_Table SHALL render a QR code icon button in the Actions column of every device row.
2. WHEN a user clicks the QR code button for a device, THE Devices_Page SHALL open the QR_Dialog for that device.
3. WHEN the QR_Dialog opens, THE Device_Service SHALL call `GET /rest/private/devices/{id}/qr` to fetch the QR code data.
4. WHILE the QR code request is in flight, THE QR_Dialog SHALL display a loading skeleton.
5. WHEN the QR code response is received, THE QR_Dialog SHALL render the QR code as an image or SVG that can be scanned by a device camera.
6. IF the QR code request fails, THEN THE QR_Dialog SHALL display an error message.
7. THE QR_Dialog SHALL provide a "Download" button that saves the QR code image to the user's device.
8. WHEN a user closes the QR_Dialog, THE Devices_Page SHALL return to the normal table view without refreshing the device list.

---

### Requirement 20: Device Service Layer

**User Story:** As a frontend developer, I want a dedicated service module for all device API calls, so that all backend communication is centralized and testable.

#### Acceptance Criteria

1. THE Device_Service SHALL expose `getDevices(params: DeviceSearchRequest): Promise<DeviceListResponse>` calling `POST /rest/private/devices/search`.
2. THE Device_Service SHALL expose `getDevice(number: string): Promise<DeviceView>` calling `GET /rest/private/devices/number/{number}`.
3. THE Device_Service SHALL expose `createDevice(payload: DevicePayload): Promise<void>` calling `POST /rest/private/devices`.
4. THE Device_Service SHALL expose `updateDevice(payload: DevicePayload): Promise<void>` calling `PUT /rest/private/devices`.
5. THE Device_Service SHALL expose `deleteDevice(id: number): Promise<void>` calling `DELETE /rest/private/devices/{id}`.
6. THE Device_Service SHALL expose `getGroups(): Promise<LookupItem[]>` calling `GET /rest/private/groups`.
7. THE Device_Service SHALL expose `getConfigurations(): Promise<ConfigurationOption[]>` calling `GET /rest/private/configurations`.
8. THE Device_Service SHALL expose `getDeviceQr(id: number): Promise<string>` calling `GET /rest/private/devices/{id}/qr`.
9. THE Device_Service SHALL use the shared `apiClient` instance for all HTTP calls so that the `X-Auth-Token` header is attached automatically.
10. IF any Device_Service call receives a non-2xx response or a `status: "ERROR"` envelope, THEN THE Device_Service SHALL propagate the error to the caller without swallowing it.

---

### Requirement 21: Device Data Types

**User Story:** As a frontend developer, I want well-typed TypeScript interfaces for all device-related data, so that the compiler catches integration errors early.

#### Acceptance Criteria

1. THE frontend SHALL define `DeviceView` with fields: `id: number`, `number: string`, `description: string | null`, `configurationId: number | null`, `groups: LookupItem[]`, `statusCode: string | null`, `lastUpdate: number | null`, `info: DeviceInfoView | null`, `androidVersion: string | null`, `serial: string | null`, `imei: string | null`, `phone: string | null`, `custom1: string | null`, `custom2: string | null`, `custom3: string | null`.
2. THE frontend SHALL define `DeviceInfoView` with fields: `batteryLevel: number | null`, `model: string | null`, `imei: string | null`, `phone: string | null`, `androidVersion: string | null`, `latitude: number | null`, `longitude: number | null`.
3. THE frontend SHALL define `DevicePayload` with fields: `id?: number`, `number: string`, `description?: string | null`, `configurationId?: number | null`, `groups?: LookupItem[]`, `imei?: string | null`, `phone?: string | null`, `custom1?: string | null`, `custom2?: string | null`, `custom3?: string | null`.
4. THE frontend SHALL define `DeviceSearchRequest` with fields: `pageNum: number`, `pageSize: number`, `value?: string`, `groupId?: number`, `configurationId?: number`, `status?: string`, `androidVersion?: string`.
5. THE frontend SHALL define `DeviceListResponse` with fields: `devices: { items: DeviceView[], totalItemsCount: number }`, `configurations: Record<string, ConfigurationView>`.
6. WHEN the backend returns a field that is absent or null, THE frontend SHALL treat it as `null` rather than throwing a runtime error.


---

### Requirement 22: Groups Management Page

**User Story:** As an MDM administrator, I want a dedicated page to create, view, edit, and delete device groups, so that I can organize my fleet into logical segments.

#### Acceptance Criteria

1. WHEN a user navigates to `/groups`, THE Groups_Page SHALL render a full-page layout consistent with the existing navigation shell.
2. WHEN the Groups_Page mounts, THE Group_Service SHALL call `GET /rest/private/groups` to fetch all groups.
3. WHEN the API response is received, THE Groups_Page SHALL display a table with columns: Group Name, Actions.
4. WHILE the API request is in flight, THE Groups_Page SHALL display a loading skeleton.
5. IF the API request fails, THEN THE Groups_Page SHALL display an error message with a "Retry" action.
6. WHEN the groups list is empty, THE Groups_Page SHALL display an empty-state message.

---

### Requirement 23: Create Group

**User Story:** As an MDM administrator, I want to create a new device group, so that I can organize devices into logical segments.

#### Acceptance Criteria

1. THE Groups_Page SHALL render an "Add Group" button in the page header.
2. WHEN a user clicks "Add Group", THE Groups_Page SHALL open a dialog with a required "Group Name" text input.
3. WHEN a user submits the dialog with a valid group name, THE Group_Service SHALL call `POST /rest/private/groups` with the group name.
4. WHEN the create request succeeds, THE Groups_Page SHALL close the dialog and refresh the groups list.
5. IF the create request fails, THEN THE Groups_Page SHALL display an error message inside the dialog and keep it open.
6. WHEN a user submits the dialog with an empty group name, THE Groups_Page SHALL display a validation error and SHALL NOT call any API endpoint.

---

### Requirement 24: Edit Group

**User Story:** As an MDM administrator, I want to rename an existing group, so that I can keep group names accurate as my organization evolves.

#### Acceptance Criteria

1. WHEN a user activates the edit action for a group, THE Groups_Page SHALL open a dialog pre-populated with the group's current name.
2. WHEN a user submits the edit dialog with a valid name, THE Group_Service SHALL call `PUT /rest/private/groups` with the updated group data.
3. WHEN the update request succeeds, THE Groups_Page SHALL close the dialog and refresh the groups list.
4. IF the update request fails, THEN THE Groups_Page SHALL display an error message inside the dialog and keep it open.

---

### Requirement 25: Delete Group

**User Story:** As an MDM administrator, I want to delete a group, so that I can remove obsolete organizational segments.

#### Acceptance Criteria

1. WHEN a user activates the delete action for a group, THE Groups_Page SHALL open a confirmation dialog identifying the group by name.
2. WHEN a user confirms deletion, THE Group_Service SHALL call `DELETE /rest/private/groups/{id}`.
3. WHEN the delete request succeeds, THE Groups_Page SHALL close the dialog and refresh the groups list.
4. IF the delete request fails, THEN THE Groups_Page SHALL display an error message and keep the dialog open.
5. WHEN a user cancels the delete dialog, THE Groups_Page SHALL close the dialog without making any API call.

---

### Requirement 26: Group Service Layer

**User Story:** As a frontend developer, I want a dedicated service module for group API calls, so that all group-related backend communication is centralized and testable.

#### Acceptance Criteria

1. THE Group_Service SHALL expose `getGroups(): Promise<LookupItem[]>` calling `GET /rest/private/groups`.
2. THE Group_Service SHALL expose `createGroup(name: string): Promise<LookupItem>` calling `POST /rest/private/groups`.
3. THE Group_Service SHALL expose `updateGroup(group: LookupItem): Promise<void>` calling `PUT /rest/private/groups`.
4. THE Group_Service SHALL expose `deleteGroup(id: number): Promise<void>` calling `DELETE /rest/private/groups/{id}`.
5. THE Group_Service SHALL use the shared `apiClient` instance for all HTTP calls so that the `X-Auth-Token` header is attached automatically.
6. IF any Group_Service call receives a non-2xx response or a `status: "ERROR"` envelope, THEN THE Group_Service SHALL propagate the error to the caller without swallowing it.


---

### Requirement 27: Parser and Serializer Round-Trip

**User Story:** As a frontend developer, I want all data transformations between API responses and UI state to be lossless, so that device data is never corrupted by serialization bugs.

#### Acceptance Criteria

1. THE Device_Service SHALL parse `DeviceListResponse` from the backend envelope such that all fields present in the API response are accessible on the resulting `DeviceListResponse` object.
2. THE Device_Service SHALL serialize `DevicePayload` to the request body such that all non-undefined fields in the payload are present in the outgoing JSON body.
3. FOR ALL valid `DevicePayload` objects, serializing to JSON and deserializing back SHALL produce an equivalent object (round-trip property).
4. FOR ALL valid `LookupItem` objects, serializing to JSON and deserializing back SHALL produce an equivalent object (round-trip property).
5. WHEN the backend returns a `DeviceView` with null optional fields, THE Device_Service SHALL return a `DeviceView` object where those fields are `null` rather than `undefined` or missing.


---

### Requirement 28: Bulk Delete via Dedicated Endpoint

**User Story:** As an MDM administrator, I want bulk delete to use the dedicated backend endpoint, so that the operation is atomic and efficient.

#### Acceptance Criteria

1. WHEN a user confirms bulk deletion of selected devices, THE Device_Service SHALL call `POST /rest/private/devices/deleteBulk` with `{ ids: number[] }` instead of issuing individual DELETE requests.
2. WHILE the bulk delete request is in flight, THE Devices_Page SHALL disable the confirm button and show a progress indicator.
3. WHEN the bulk delete request succeeds, THE Devices_Page SHALL close the confirmation dialog, clear the selection, and refresh the device list.
4. IF the bulk delete request fails, THEN THE Devices_Page SHALL display an error message and keep the dialog open.

---

### Requirement 29: Bulk Set Group via Dedicated Endpoint

**User Story:** As an MDM administrator, I want bulk group assignment to use the dedicated backend endpoint, so that the operation is atomic and efficient.

#### Acceptance Criteria

1. WHEN a user confirms bulk group assignment, THE Device_Service SHALL call `POST /rest/private/devices/groupBulk` with `{ ids: number[], action: 'set' | 'clear', groups: LookupItem[] }`.
2. WHILE the bulk group request is in flight, THE Devices_Page SHALL disable the confirm button and show a progress indicator.
3. WHEN the bulk group request succeeds, THE Devices_Page SHALL close the dialog, clear the selection, and refresh the device list.
4. IF the bulk group request fails, THEN THE Devices_Page SHALL display an error message and keep the dialog open.

---

### Requirement 30: Status Badge — Full Status Code Mapping

**User Story:** As an MDM administrator, I want all backend status codes to be visually distinct, so that I can differentiate between warning states and fully offline devices.

#### Acceptance Criteria

1. THE Status_Badge SHALL map all backend `statusCode` values as follows:
   - `"green"` → green badge labeled "Online"
   - `"red"` → red badge labeled "Offline"
   - `"yellow"` → yellow/amber badge labeled "Warning"
   - `"brown"` → orange badge labeled "Inactive"
   - `"grey"` or any other value → gray badge labeled "Unknown"
2. THE Status_Badge SHALL display a rich tooltip on hover showing: the humanized elapsed time since last contact (e.g., "3 hours ago") and the exact formatted timestamp on a second line.
3. THE Status_Badge SHALL be rendered in the Status column of every device row and in the Detail_Panel header.

---

### Requirement 31: Computed Status Columns

**User Story:** As an MDM administrator, I want to see permission, installation, and file compliance status at a glance in the device table, so that I can identify non-compliant devices quickly.

#### Acceptance Criteria

1. THE Device_Table SHALL support an optional "Permission Status" column that computes a health indicator from `device.info.permissions`: green when all required permissions are granted, amber when some are missing, red when critical permissions are denied.
2. THE Device_Table SHALL support an optional "Installation Status" column that compares the configuration's required applications against `device.info.applications`: green when all apps are installed at the correct version, amber when versions mismatch, red when apps are missing.
3. THE Device_Table SHALL support an optional "Files Status" column that compares the configuration's required files against `device.info.files`: green when all files are present and current, amber when timestamps mismatch, red when files are missing.
4. WHEN a user hovers over a computed status indicator, THE Device_Table SHALL display a tooltip with the specific details of the compliance check.
5. THE computed status columns SHALL be hidden by default and toggleable via the Column_Visibility_Menu.

---

### Requirement 32: Application Settings per Device

**User Story:** As an MDM administrator, I want to manage per-device application settings, so that I can customize app behavior for individual devices.

#### Acceptance Criteria

1. WHEN a user opens the row action dropdown for a device, THE Devices_Page SHALL include an "App Settings" menu item.
2. WHEN a user clicks "App Settings", THE Devices_Page SHALL open the App_Settings_Dialog for that device.
3. WHEN the App_Settings_Dialog opens, THE Device_Service SHALL call `GET /rest/private/devices/{id}/applicationSettings` to fetch the current settings list.
4. WHILE the fetch is in flight, THE App_Settings_Dialog SHALL display a loading skeleton.
5. WHEN the settings response is received, THE App_Settings_Dialog SHALL display a list of settings with fields: application name, setting name, value, comment, and a read-only indicator.
6. THE App_Settings_Dialog SHALL allow the user to add, edit, and delete setting items.
7. WHEN a user saves changes, THE Device_Service SHALL call `POST /rest/private/devices/{id}/applicationSettings` with the updated settings array.
8. WHEN the save request succeeds, THE Device_Service SHALL call `POST /rest/private/devices/{id}/applicationSettings/notify` to push the update to the device.
9. IF the save or notify request fails, THEN THE App_Settings_Dialog SHALL display an error message and keep the dialog open.
10. WHEN a user closes the App_Settings_Dialog without saving, THE Devices_Page SHALL NOT make any API call.

---

### Requirement 33: Auto-Refresh

**User Story:** As an MDM administrator, I want the device list to refresh automatically, so that I always see up-to-date device status without manually reloading.

#### Acceptance Criteria

1. WHEN the Devices_Page mounts, THE Devices_Page SHALL initiate an automatic refresh interval of 60 seconds.
2. WHEN 60 seconds elapse since the last successful fetch, THE Devices_Page SHALL re-fetch the device list using the current search, filter, and page state — without showing loading skeletons (background refresh).
3. WHEN the Devices_Page unmounts, THE Devices_Page SHALL cancel the auto-refresh interval to prevent memory leaks.
4. WHILE an auto-refresh fetch is in progress, THE Devices_Page SHALL continue displaying the previously loaded data rather than showing skeleton components.
5. IF an auto-refresh fetch fails, THE Devices_Page SHALL silently ignore the error and retry on the next interval.

---

### Requirement 34: Description-Only Edit

**User Story:** As an MDM administrator with limited permissions, I want to edit only the device description, so that I can annotate devices without full edit access.

#### Acceptance Criteria

1. WHEN a user has the `edit_device_desc` permission but not `edit_devices`, THE Devices_Page SHALL render an "Edit Description" action in the row dropdown instead of the full "Edit" action.
2. WHEN a user clicks "Edit Description", THE Devices_Page SHALL open a minimal dialog containing only the `description` textarea field.
3. WHEN a user submits the description dialog, THE Device_Service SHALL call `POST /rest/private/devices/{id}/description` with the updated description.
4. WHEN the request succeeds, THE Devices_Page SHALL close the dialog and refresh the device list.
5. IF the request fails, THEN THE Devices_Page SHALL display an error message inside the dialog and keep it open.

---

### Requirement 35: Advanced Device Form Validation

**User Story:** As an MDM administrator, I want the device form to enforce all business rules, so that invalid device records cannot be created.

#### Acceptance Criteria

1. THE Device_Form SHALL reject a `number` field value that contains any of the characters `/`, `?`, or `&`, and SHALL display a validation error.
2. WHEN a device is already enrolled and the user changes its `number` field in edit mode, THE Device_Form SHALL display a migration confirmation dialog warning that changing the device ID will affect the enrolled device.
3. THE Device_Form SHALL render custom field labels (`custom1`, `custom2`, `custom3`) using the label names configured in the application settings when available, falling back to "Custom 1", "Custom 2", "Custom 3".
4. WHEN the application settings specify `customMultiline1`, `customMultiline2`, or `customMultiline3` as true, THE Device_Form SHALL render the corresponding custom field as a `Textarea` instead of a single-line `Input`.

---

### Requirement 36: Updated Device Service Layer (Extended)

**User Story:** As a frontend developer, I want the device service to expose all backend endpoints, so that every feature has a clean API boundary.

#### Acceptance Criteria

1. THE Device_Service SHALL expose `deleteBulk(ids: number[]): Promise<void>` calling `POST /rest/private/devices/deleteBulk`.
2. THE Device_Service SHALL expose `groupBulk(payload: GroupBulkPayload): Promise<void>` calling `POST /rest/private/devices/groupBulk`.
3. THE Device_Service SHALL expose `getAppSettings(deviceId: number): Promise<AppSetting[]>` calling `GET /rest/private/devices/{id}/applicationSettings`.
4. THE Device_Service SHALL expose `saveAppSettings(deviceId: number, settings: AppSetting[]): Promise<void>` calling `POST /rest/private/devices/{id}/applicationSettings`.
5. THE Device_Service SHALL expose `notifyAppSettings(deviceId: number): Promise<void>` calling `POST /rest/private/devices/{id}/applicationSettings/notify`.
6. THE Device_Service SHALL expose `updateDescription(deviceId: number, description: string): Promise<void>` calling `POST /rest/private/devices/{id}/description`.
7. THE Device_Service SHALL expose `autocomplete(value: string): Promise<string[]>` calling `POST /rest/private/devices/autocomplete`.
8. IF any Device_Service call receives a non-2xx response or a `status: "ERROR"` envelope, THEN THE Device_Service SHALL propagate the error to the caller without swallowing it.

---

### Requirement 37: Extended Device Data Types

**User Story:** As a frontend developer, I want complete TypeScript interfaces for all device-related data including the extended search request and bulk payloads.

#### Acceptance Criteria

1. THE frontend SHALL extend `DeviceSearchRequest` to include all backend-supported fields: `sortBy?: string`, `sortDir?: 'asc' | 'desc'`, `dateFrom?: number`, `dateTo?: number`, `onlineEarlierMillis?: number`, `onlineLaterMillis?: number`, `enrollmentDateFrom?: number`, `enrollmentDateTo?: number`, `mdmMode?: boolean`, `kioskMode?: boolean`, `launcherVersion?: string`, `installationStatus?: string`, `imeiChanged?: boolean`, `fastSearch?: boolean`.
2. THE frontend SHALL define `GroupBulkPayload` with fields: `ids: number[]`, `action: 'set' | 'clear'`, `groups: LookupItem[]`.
3. THE frontend SHALL define `AppSetting` with fields: `id?: number`, `applicationId: number`, `name: string`, `value: string`, `comment?: string`, `readonly?: boolean`, `lastUpdate?: number`.
4. THE frontend SHALL define `BulkDeletePayload` with field: `ids: number[]`.
5. THE frontend SHALL extend `DeviceInfoView` to include: `permissions: DevicePermissions | null`, `applications: DeviceApplication[] | null`, `files: DeviceFile[] | null`, `defaultLauncher: boolean | null`, `mdmMode: boolean | null`, `kioskMode: boolean | null`, `enrollTime: number | null`, `publicIp: string | null`, `launcherVersion: string | null`.
