# Requirements Document

## Introduction

The Devices Management feature provides a dedicated page at `/devices` for viewing, searching, and managing enrolled MDM devices. It replaces the placeholder currently at that route with a fully functional interface: a paginated, searchable table of devices, a detail panel for inspecting individual device information, color-coded status indicators, and a delete action with confirmation. All data is fetched from the existing HMDM backend REST API using the shared `apiClient` (base `/rest`, `X-Auth-Token` header injected automatically). The UI is built exclusively with shadcn/ui components inside the existing React + Vite + TypeScript frontend.

## Glossary

- **Device**: A mobile or desktop endpoint enrolled in Headwind MDM, identified by a unique numeric `id` and a human-readable `name` and `number` (device ID string).
- **Device_List**: The paginated collection of `Device` records returned by `GET /rest/private/devices`.
- **Device_Detail**: The full record for a single device returned by `GET /rest/private/devices/{id}`, including battery level, GPS location, and last update timestamp.
- **Status_Badge**: A colored `Badge` component that communicates a device's connectivity state (online / offline / unknown) derived from `statusCode`.
- **Devices_Page**: The React page component rendered at the `/devices` route.
- **Device_Service**: The frontend service module (`src/features/devices/deviceService.ts`) that wraps all `apiClient` calls for device-related endpoints.
- **Delete_Dialog**: The shadcn/ui `AlertDialog` confirmation modal shown before a device is permanently deleted.
- **Search_Input**: The shadcn/ui `Input` component used to filter devices by name or device number.
- **Pagination_Controls**: The shadcn/ui `Pagination` component that navigates between pages of the device list.
- **Detail_Panel**: The shadcn/ui `Sheet` (side panel) or `Dialog` that displays full `Device_Detail` information for a selected device.

---

## Requirements

### Requirement 1: Devices List Page

**User Story:** As an MDM administrator, I want to see all enrolled devices in a paginated table, so that I can quickly survey the fleet and navigate to individual devices.

#### Acceptance Criteria

1. WHEN a user navigates to `/devices`, THE Devices_Page SHALL render a full-page layout consistent with the existing navigation shell.
2. WHEN the Devices_Page mounts, THE Device_Service SHALL call `GET /rest/private/devices` with default pagination parameters (page 0, pageSize 20).
3. WHEN the API response is received, THE Devices_Page SHALL display a shadcn/ui `Table` with columns: Device Name, Device Number, Status, Configuration, Group, Last Seen.
4. WHILE the API request is in flight, THE Devices_Page SHALL display a loading skeleton or spinner in place of the table.
5. IF the API request fails, THEN THE Devices_Page SHALL display an error message with a retry action.
6. WHEN the device list is empty, THE Devices_Page SHALL display an empty-state message indicating no devices are enrolled.

### Requirement 2: Search

**User Story:** As an MDM administrator, I want to filter devices by name or device number, so that I can quickly locate a specific device in a large fleet.

#### Acceptance Criteria

1. THE Devices_Page SHALL render a Search_Input above the device table.
2. WHEN a user types in the Search_Input, THE Devices_Page SHALL debounce the input by 300 ms before issuing a new API request.
3. WHEN a debounced search term is applied, THE Device_Service SHALL call `GET /rest/private/devices` with the search term as a query parameter and reset the page to 0.
4. WHEN the search term is cleared, THE Device_Service SHALL call `GET /rest/private/devices` with no search parameter and reset the page to 0.
5. IF a search returns no results, THEN THE Devices_Page SHALL display a "no devices found" message that includes the search term.

### Requirement 3: Pagination

**User Story:** As an MDM administrator, I want to navigate between pages of devices, so that I can browse large fleets without loading all records at once.

#### Acceptance Criteria

1. WHEN the API response includes a total count greater than the page size, THE Devices_Page SHALL render Pagination_Controls below the table.
2. WHEN a user clicks a page number in the Pagination_Controls, THE Device_Service SHALL call `GET /rest/private/devices` with the updated page index.
3. WHEN a new page is requested, THE Devices_Page SHALL display a loading indicator while the request is in flight.
4. WHEN the search term changes, THE Pagination_Controls SHALL reset to page 1.
5. THE Devices_Page SHALL display the current page number and total device count as text alongside the Pagination_Controls.

### Requirement 4: Device Detail View

**User Story:** As an MDM administrator, I want to inspect the full details of a device, so that I can diagnose issues and review device health.

#### Acceptance Criteria

1. WHEN a user clicks a device row in the table, THE Devices_Page SHALL open the Detail_Panel for that device.
2. WHEN the Detail_Panel opens, THE Device_Service SHALL call `GET /rest/private/devices/{id}` to fetch the full device record.
3. WHILE the detail request is in flight, THE Detail_Panel SHALL display a loading skeleton.
4. WHEN the detail response is received, THE Detail_Panel SHALL display: device name, device number, status badge, configuration name, group name, battery level (percentage), last update timestamp (formatted as a human-readable date/time), and GPS coordinates or "Location unavailable" if absent.
5. IF the detail request fails, THEN THE Detail_Panel SHALL display an error message.
6. WHEN a user closes the Detail_Panel, THE Devices_Page SHALL return focus to the previously selected table row.

### Requirement 5: Status Indicators

**User Story:** As an MDM administrator, I want to see at a glance whether each device is online, offline, or in an unknown state, so that I can prioritize attention to unreachable devices.

#### Acceptance Criteria

1. THE Status_Badge SHALL map `statusCode` values to display states: a `statusCode` of `1` SHALL render a green badge labeled "Online"; a `statusCode` of `2` SHALL render a red badge labeled "Offline"; any other `statusCode` value SHALL render a gray badge labeled "Unknown".
2. THE Status_Badge SHALL be rendered in the Status column of every device row in the table.
3. THE Status_Badge SHALL be rendered in the Detail_Panel header for the selected device.
4. WHEN the status is "Online", THE Status_Badge SHALL use the `default` or `success` variant; WHEN the status is "Offline", THE Status_Badge SHALL use the `destructive` variant; WHEN the status is "Unknown", THE Status_Badge SHALL use the `secondary` variant.

### Requirement 6: Delete Device

**User Story:** As an MDM administrator, I want to remove a device from the MDM system, so that I can keep the device inventory accurate.

#### Acceptance Criteria

1. WHEN a user activates the delete action for a device (via a row action menu or button), THE Devices_Page SHALL open the Delete_Dialog identifying the device by name.
2. WHEN a user confirms deletion in the Delete_Dialog, THE Device_Service SHALL call `DELETE /rest/private/devices/{id}`.
3. WHILE the delete request is in flight, THE Delete_Dialog SHALL disable the confirm button and show a loading indicator.
4. WHEN the delete request succeeds, THE Devices_Page SHALL close the Delete_Dialog and refresh the device list.
5. IF the delete request fails, THEN THE Delete_Dialog SHALL display an error message and keep the dialog open.
6. WHEN a user cancels the Delete_Dialog, THE Devices_Page SHALL close the dialog without making any API call.

### Requirement 7: Device Service Layer

**User Story:** As a frontend developer, I want a dedicated service module for device API calls, so that all backend communication is centralized and testable.

#### Acceptance Criteria

1. THE Device_Service SHALL expose a `getDevices(params: DeviceListParams): Promise<DeviceListResponse>` function that calls `GET /rest/private/devices`.
2. THE Device_Service SHALL expose a `getDevice(id: number): Promise<Device>` function that calls `GET /rest/private/devices/{id}`.
3. THE Device_Service SHALL expose a `updateDevice(device: Device): Promise<Device>` function that calls `PUT /rest/private/devices`.
4. THE Device_Service SHALL expose a `deleteDevice(id: number): Promise<void>` function that calls `DELETE /rest/private/devices/{id}`.
5. THE Device_Service SHALL use the shared `apiClient` instance for all HTTP calls so that the `X-Auth-Token` header is attached automatically.
6. IF any Device_Service call receives a non-2xx response, THEN THE Device_Service SHALL propagate the error to the caller without swallowing it.

### Requirement 8: Data Types

**User Story:** As a frontend developer, I want well-typed TypeScript interfaces for all device-related data, so that the compiler catches integration errors early.

#### Acceptance Criteria

1. THE frontend SHALL define a `Device` interface with fields: `id: number`, `name: string`, `number: string`, `statusCode: number`, `configurationId: number | null`, `groupId: number | null`, `lastUpdate: number | null` (Unix timestamp ms), `info: DeviceInfo | null`.
2. THE frontend SHALL define a `DeviceInfo` interface with fields: `batteryLevel: number | null`, `latitude: number | null`, `longitude: number | null`.
3. THE frontend SHALL define a `DeviceListParams` interface with fields: `page: number`, `pageSize: number`, `search?: string`.
4. THE frontend SHALL define a `DeviceListResponse` interface with fields: `data: Device[]`, `total: number`, `page: number`, `pageSize: number`.
5. WHEN the backend returns a field that is absent or null, THE frontend SHALL treat it as `null` rather than throwing a runtime error.
