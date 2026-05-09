# Requirements Document

## Introduction

The Configurations Management feature provides a dedicated page at `/configurations` for viewing, creating, editing, and deleting MDM device configurations. A configuration is a named policy profile (type COMMON or WORK) that can be assigned to devices. The page renders a table of all configurations with key metadata, supports creating and editing configurations via a dialog form, and allows deletion with a confirmation prompt. A device count column shows how many devices use each configuration. All data is fetched from the existing HMDM backend REST API using the shared `apiClient` (base `/rest`, `X-Auth-Token` header injected automatically). The UI is built exclusively with shadcn/ui components inside the existing React + Vite + TypeScript frontend. The `/configurations` route is added to the navigation sidebar after the Devices entry.

## Glossary

- **Configuration**: An MDM policy profile identified by a unique numeric `id`, a human-readable `name`, an optional `description`, a `type` (either `COMMON` or `WORK`), an optional `applicationId`, and associated `files` and `applications` lists.
- **Configuration_List**: The collection of `Configuration` records returned by `GET /rest/private/configurations`.
- **Configurations_Page**: The React page component rendered at the `/configurations` route.
- **Configuration_Service**: The frontend service module (`src/features/configurations/configurationService.ts`) that wraps all `apiClient` calls for configuration-related endpoints.
- **Configuration_Form**: The shadcn/ui `Dialog` containing a `Form` used to create or edit a `Configuration`.
- **Delete_Dialog**: The shadcn/ui `AlertDialog` confirmation modal shown before a configuration is permanently deleted.
- **Type_Selector**: The shadcn/ui `Select` component used to choose between `COMMON` and `WORK` configuration types.
- **Device_Count**: The number of devices currently assigned to a given configuration, derived from the `applications` list or a dedicated count field returned by the API.

---

## Requirements

### Requirement 1: Configurations List Page

**User Story:** As an MDM administrator, I want to see all configurations in a table, so that I can survey available policy profiles and manage them.

#### Acceptance Criteria

1. WHEN a user navigates to `/configurations`, THE Configurations_Page SHALL render a full-page layout consistent with the existing navigation shell.
2. WHEN the Configurations_Page mounts, THE Configuration_Service SHALL call `GET /rest/private/configurations`.
3. WHEN the API response is received, THE Configurations_Page SHALL display a shadcn/ui `Table` with columns: Name, Type, Description, Device Count.
4. WHILE the API request is in flight, THE Configurations_Page SHALL display a loading skeleton or spinner in place of the table.
5. IF the API request fails, THEN THE Configurations_Page SHALL display an error message with a retry action.
6. WHEN the configuration list is empty, THE Configurations_Page SHALL display an empty-state message indicating no configurations exist.

### Requirement 2: Create Configuration

**User Story:** As an MDM administrator, I want to create a new configuration, so that I can define a new policy profile for devices.

#### Acceptance Criteria

1. THE Configurations_Page SHALL render a "New Configuration" `Button` above the table.
2. WHEN a user clicks the "New Configuration" button, THE Configurations_Page SHALL open the Configuration_Form dialog in create mode with all fields empty.
3. WHEN a user submits the Configuration_Form with valid data, THE Configuration_Service SHALL call `POST /rest/private/configurations` with the form values.
4. WHILE the create request is in flight, THE Configuration_Form SHALL disable the submit button and show a loading indicator.
5. WHEN the create request succeeds, THE Configurations_Page SHALL close the Configuration_Form and refresh the configuration list.
6. IF the create request fails, THEN THE Configuration_Form SHALL display an error message and keep the dialog open.
7. WHEN a user cancels the Configuration_Form, THE Configurations_Page SHALL close the dialog without making any API call.

### Requirement 3: Edit Configuration

**User Story:** As an MDM administrator, I want to edit an existing configuration, so that I can update its name, description, or type.

#### Acceptance Criteria

1. WHEN a user activates the edit action for a configuration (via a row action menu or button), THE Configurations_Page SHALL open the Configuration_Form dialog in edit mode pre-populated with the configuration's current `name`, `description`, and `type`.
2. WHEN the Configuration_Form is in edit mode and the user submits valid data, THE Configuration_Service SHALL call `PUT /rest/private/configurations/{id}` with the updated values.
3. WHILE the update request is in flight, THE Configuration_Form SHALL disable the submit button and show a loading indicator.
4. WHEN the update request succeeds, THE Configurations_Page SHALL close the Configuration_Form and refresh the configuration list.
5. IF the update request fails, THEN THE Configuration_Form SHALL display an error message and keep the dialog open.

### Requirement 4: Delete Configuration

**User Story:** As an MDM administrator, I want to delete a configuration, so that I can remove obsolete policy profiles.

#### Acceptance Criteria

1. WHEN a user activates the delete action for a configuration (via a row action menu or button), THE Configurations_Page SHALL open the Delete_Dialog identifying the configuration by name.
2. WHEN a user confirms deletion in the Delete_Dialog, THE Configuration_Service SHALL call `DELETE /rest/private/configurations/{id}`.
3. WHILE the delete request is in flight, THE Delete_Dialog SHALL disable the confirm button and show a loading indicator.
4. WHEN the delete request succeeds, THE Configurations_Page SHALL close the Delete_Dialog and refresh the configuration list.
5. IF the delete request fails, THEN THE Delete_Dialog SHALL display an error message and keep the dialog open.
6. WHEN a user cancels the Delete_Dialog, THE Configurations_Page SHALL close the dialog without making any API call.

### Requirement 5: Configuration Form Validation

**User Story:** As an MDM administrator, I want the configuration form to validate my input, so that I cannot submit incomplete or invalid data.

#### Acceptance Criteria

1. THE Configuration_Form SHALL require the `name` field; WHEN the `name` field is empty on submit, THE Configuration_Form SHALL display a validation error on the `name` field and prevent submission.
2. THE Configuration_Form SHALL require the `type` field; WHEN the `type` field is unselected on submit, THE Configuration_Form SHALL display a validation error on the `type` field and prevent submission.
3. THE Configuration_Form SHALL render the `description` field as optional; WHEN the `description` field is empty, THE Configuration_Form SHALL allow submission.
4. THE Type_Selector SHALL present exactly two options: `COMMON` and `WORK`.
5. WHEN a validation error is present, THE Configuration_Form SHALL display the error message adjacent to the invalid field using the shadcn/ui `FormMessage` component.

### Requirement 6: Device Count Display

**User Story:** As an MDM administrator, I want to see how many devices use each configuration, so that I can assess the impact of changes before editing or deleting.

#### Acceptance Criteria

1. THE Configurations_Page SHALL display a Device Count column in the configuration table.
2. WHEN the API response includes device assignment data for a configuration, THE Configurations_Page SHALL display the count of assigned devices in the Device Count column.
3. WHEN a configuration has zero assigned devices, THE Configurations_Page SHALL display `0` in the Device Count column.

### Requirement 7: Navigation Integration

**User Story:** As an MDM administrator, I want a Configurations link in the sidebar, so that I can navigate to the configurations page from anywhere in the application.

#### Acceptance Criteria

1. THE navigation sidebar SHALL include a "Configurations" entry positioned after the "Devices" entry.
2. WHEN a user clicks the "Configurations" sidebar entry, THE application SHALL navigate to `/configurations`.
3. WHEN the current route is `/configurations`, THE navigation sidebar SHALL render the "Configurations" entry in its active/selected visual state.

### Requirement 8: Configuration Service Layer

**User Story:** As a frontend developer, I want a dedicated service module for configuration API calls, so that all backend communication is centralized and testable.

#### Acceptance Criteria

1. THE Configuration_Service SHALL expose a `getConfigurations(): Promise<Configuration[]>` function that calls `GET /rest/private/configurations`.
2. THE Configuration_Service SHALL expose a `getConfiguration(id: number): Promise<Configuration>` function that calls `GET /rest/private/configurations/{id}`.
3. THE Configuration_Service SHALL expose a `createConfiguration(data: ConfigurationPayload): Promise<Configuration>` function that calls `POST /rest/private/configurations`.
4. THE Configuration_Service SHALL expose an `updateConfiguration(id: number, data: ConfigurationPayload): Promise<Configuration>` function that calls `PUT /rest/private/configurations/{id}`.
5. THE Configuration_Service SHALL expose a `deleteConfiguration(id: number): Promise<void>` function that calls `DELETE /rest/private/configurations/{id}`.
6. THE Configuration_Service SHALL use the shared `apiClient` instance for all HTTP calls so that the `X-Auth-Token` header is attached automatically.
7. IF any Configuration_Service call receives a non-2xx response, THEN THE Configuration_Service SHALL propagate the error to the caller without swallowing it.

### Requirement 9: Data Types

**User Story:** As a frontend developer, I want well-typed TypeScript interfaces for all configuration-related data, so that the compiler catches integration errors early.

#### Acceptance Criteria

1. THE frontend SHALL define a `Configuration` interface with fields: `id: number`, `name: string`, `description: string | null`, `type: 'COMMON' | 'WORK'`, `applicationId: number | null`, `files: ConfigurationFile[]`, `applications: ConfigurationApplication[]`.
2. THE frontend SHALL define a `ConfigurationFile` interface with at minimum the field `id: number`.
3. THE frontend SHALL define a `ConfigurationApplication` interface with at minimum the field `id: number`.
4. THE frontend SHALL define a `ConfigurationPayload` interface with fields: `name: string`, `description?: string`, `type: 'COMMON' | 'WORK'`.
5. WHEN the backend returns a field that is absent or null, THE frontend SHALL treat it as `null` or an empty array rather than throwing a runtime error.
