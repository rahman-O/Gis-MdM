# Requirements Document

## Introduction

The Devices Add/Edit feature extends the existing `/devices` page with the ability to create new devices and edit existing ones. An "Add Device" button on `DevicesPage` opens a `DeviceForm` dialog in create mode. An "Edit" action in each row's dropdown opens the same dialog in edit mode, pre-populated with the device's current values. The form validates required fields, fetches group and configuration options from the backend, and calls the appropriate REST endpoint on submit. On success the dialog closes and the device list refreshes. All UI is built from shadcn/ui components following the patterns established by configurations-management and users-management.

## Glossary

- **Device_Form**: The shadcn/ui `Dialog`-based form component (`src/features/devices/DeviceForm.tsx`) used for both creating and editing devices.
- **Device_Service**: The existing frontend service module (`src/features/devices/deviceService.ts`) extended with `createDevice` and updated `updateDevice` functions.
- **Group_Selector**: The multi-select control inside `Device_Form` that lets the user assign a device to one or more groups, populated from `GET /rest/private/groups`.
- **Configuration_Selector**: The single-select control inside `Device_Form` that assigns a configuration to the device, populated from `GET /rest/private/configurations`.
- **DevicePayload**: The request body sent to `POST /rest/private/devices` (create) or `PUT /rest/private/devices` (update).
- **Devices_Page**: The existing React page component at `/devices` that owns the "Add Device" button and the row action menu.
- **Custom_Field**: One of `custom1`, `custom2`, or `custom3` — optional string properties whose display labels come from the application settings.

---

## Requirements

### Requirement 1: Add Device Entry Point

**User Story:** As an MDM administrator, I want an "Add Device" button on the devices page, so that I can enroll a new device without leaving the page.

#### Acceptance Criteria

1. THE Devices_Page SHALL render an "Add Device" button in the page header area, alongside the existing search input.
2. WHEN a user clicks the "Add Device" button, THE Devices_Page SHALL open the Device_Form dialog in create mode with all fields empty.
3. WHEN the Device_Form is open in create mode, THE Device_Form SHALL display the title "Add Device".

### Requirement 2: Edit Device Entry Point

**User Story:** As an MDM administrator, I want to edit an existing device's details from the device list, so that I can update its configuration, groups, or metadata without re-enrolling it.

#### Acceptance Criteria

1. WHEN a user opens the row action dropdown for a device, THE Devices_Page SHALL include an "Edit" menu item.
2. WHEN a user clicks the "Edit" menu item for a device, THE Devices_Page SHALL open the Device_Form dialog in edit mode pre-populated with that device's current field values.
3. WHEN the Device_Form is open in edit mode, THE Device_Form SHALL display the title "Edit Device".
4. WHEN the Device_Form is open in edit mode, THE Device_Form SHALL render the `number` field as read-only (disabled).

### Requirement 3: Device Form Fields

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

### Requirement 4: Form Validation

**User Story:** As an MDM administrator, I want the form to catch invalid input before submission, so that I don't accidentally create a device with missing or malformed data.

#### Acceptance Criteria

1. WHEN a user submits the Device_Form with an empty or whitespace-only `number` field, THE Device_Form SHALL display a validation error adjacent to the `number` field and SHALL NOT call any API endpoint.
2. WHEN a user submits the Device_Form with an `imei` value that is not exactly 15 digits, THE Device_Form SHALL display a validation error adjacent to the `imei` field and SHALL NOT call any API endpoint.
3. WHEN all required fields are valid, THE Device_Form SHALL enable the submit button and allow submission.
4. WHEN a user submits the Device_Form with an empty `imei` field, THE Device_Form SHALL treat the field as absent and SHALL NOT apply the 15-digit validation rule.

### Requirement 5: Create Device

**User Story:** As an MDM administrator, I want to save a new device to the MDM system, so that it can be enrolled and managed.

#### Acceptance Criteria

1. WHEN a user submits the Device_Form in create mode with valid data, THE Device_Service SHALL call `POST /rest/private/devices` with a DevicePayload containing the form values.
2. WHILE the create request is in flight, THE Device_Form SHALL disable the submit button and display a loading indicator.
3. WHEN the create request succeeds, THE Device_Form SHALL close the dialog and THE Devices_Page SHALL refresh the device list.
4. IF the create request fails, THEN THE Device_Form SHALL display an error message inside the dialog and SHALL keep the dialog open with the submit button re-enabled.

### Requirement 6: Update Device

**User Story:** As an MDM administrator, I want to save changes to an existing device, so that its configuration, groups, and metadata stay current.

#### Acceptance Criteria

1. WHEN a user submits the Device_Form in edit mode with valid data, THE Device_Service SHALL call `PUT /rest/private/devices` with a DevicePayload containing the updated form values and the device's existing `id`.
2. WHILE the update request is in flight, THE Device_Form SHALL disable the submit button and display a loading indicator.
3. WHEN the update request succeeds, THE Device_Form SHALL close the dialog and THE Devices_Page SHALL refresh the device list.
4. IF the update request fails, THEN THE Device_Form SHALL display an error message inside the dialog and SHALL keep the dialog open with the submit button re-enabled.

### Requirement 7: Cancel and Close

**User Story:** As an MDM administrator, I want to dismiss the device form without saving, so that I can abort an accidental or incomplete edit.

#### Acceptance Criteria

1. THE Device_Form SHALL render a "Cancel" button that closes the dialog without making any API call.
2. WHEN a user closes the Device_Form dialog (via Cancel or the dialog's close control), THE Devices_Page SHALL NOT refresh the device list unless a successful save occurred.

### Requirement 8: Group Selector

**User Story:** As an MDM administrator, I want to assign a device to one or more groups, so that group-based policies apply to it.

#### Acceptance Criteria

1. WHEN the Device_Form opens, THE Group_Selector SHALL fetch available groups from `GET /rest/private/groups`.
2. WHILE groups are loading, THE Group_Selector SHALL display a loading state and SHALL be disabled.
3. THE Group_Selector SHALL allow the user to select zero or more groups from the fetched list.
4. WHEN the Device_Form is in edit mode, THE Group_Selector SHALL pre-select the groups already assigned to the device.
5. IF the groups fetch fails, THEN THE Group_Selector SHALL display an error message and SHALL remain disabled.

### Requirement 9: Configuration Selector

**User Story:** As an MDM administrator, I want to assign a configuration policy to a device, so that the correct MDM policy is applied.

#### Acceptance Criteria

1. WHEN the Device_Form opens, THE Configuration_Selector SHALL fetch available configurations from `GET /rest/private/configurations`.
2. WHILE configurations are loading, THE Configuration_Selector SHALL display a loading state and SHALL be disabled.
3. THE Configuration_Selector SHALL allow the user to select at most one configuration, or leave the field empty.
4. WHEN the Device_Form is in edit mode, THE Configuration_Selector SHALL pre-select the configuration matching the device's `configurationId`.
5. IF the configurations fetch fails, THEN THE Configuration_Selector SHALL display an error message and SHALL remain disabled.

### Requirement 10: Device Service Extension

**User Story:** As a frontend developer, I want the device service to expose create and update functions, so that the form has a clean API boundary for backend communication.

#### Acceptance Criteria

1. THE Device_Service SHALL expose a `createDevice(payload: DevicePayload): Promise<void>` function that calls `POST /rest/private/devices`.
2. THE Device_Service SHALL expose an `updateDevice(payload: DevicePayload): Promise<void>` function that calls `PUT /rest/private/devices`.
3. THE Device_Service SHALL use the shared `apiClient` instance for all HTTP calls so that the `X-Auth-Token` header is attached automatically.
4. IF any Device_Service call receives a non-2xx response or a `status: "ERROR"` envelope, THEN THE Device_Service SHALL propagate the error to the caller without swallowing it.
