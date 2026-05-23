# Requirements Document

## Introduction

The Settings Management feature provides a dedicated page at `/settings` for viewing and updating global HMDM instance settings. The page renders a single form containing all configurable fields returned by `GET /rest/private/settings` and persists changes via `PUT /rest/private/settings`. On successful save the user receives a success toast; on failure an error toast is shown. Form validation prevents submission of incomplete required fields. The UI is built exclusively with shadcn/ui components (Form, Input, Switch, Select, Button) inside the existing React + Vite + TypeScript frontend.

## Glossary

- **Settings**: The global HMDM configuration object identified by a numeric `id`, containing fields that control device enrollment, security policy, and UI behaviour.
- **Settings_Page**: The React page component rendered at the `/settings` route.
- **Settings_Form**: The shadcn/ui `Form` rendered inside the Settings_Page, bound to all settings fields.
- **Settings_Service**: The frontend service module (`src/features/settings/settingsService.ts`) that wraps all `apiClient` calls for settings-related endpoints.
- **Password_Strength**: An enumerated value controlling the complexity requirement for device passwords; valid values are `0` (any), `1` (numeric), `2` (alphabetic), `3` (alphanumeric).
- **Toast**: A transient shadcn/ui notification displayed in response to a save operation outcome.

---

## Requirements

### Requirement 1: Settings Page Layout

**User Story:** As an MDM administrator, I want a dedicated settings page, so that I can view and manage all global HMDM configuration in one place.

#### Acceptance Criteria

1. WHEN a user navigates to `/settings`, THE Settings_Page SHALL render a full-page layout consistent with the existing navigation shell.
2. WHEN the Settings_Page mounts, THE Settings_Service SHALL call `GET /rest/private/settings`.
3. WHILE the GET request is in flight, THE Settings_Page SHALL display a loading skeleton or spinner in place of the form.
4. WHEN the GET response is received, THE Settings_Page SHALL populate the Settings_Form with the returned settings values.
5. IF the GET request fails, THEN THE Settings_Page SHALL display an error message with a retry action.

### Requirement 2: Settings Form Fields

**User Story:** As an MDM administrator, I want to see and edit all settings fields in a single form, so that I can configure the HMDM instance without navigating between multiple pages.

#### Acceptance Criteria

1. THE Settings_Form SHALL render an `Input` field bound to `customerName`.
2. THE Settings_Form SHALL render a `Switch` component bound to `createNewDevices`.
3. THE Settings_Form SHALL render a `Select` component bound to `newDeviceConfigurationId` populated with available configurations fetched from `GET /rest/private/configurations`.
4. THE Settings_Form SHALL render a `Select` component bound to `language` with options for supported UI languages.
5. THE Settings_Form SHALL render an `Input` field of type `number` bound to `passwordLength`.
6. THE Settings_Form SHALL render a `Select` component bound to `passwordStrength` with options: `0` (Any), `1` (Numeric), `2` (Alphabetic), `3` (Alphanumeric).
7. THE Settings_Form SHALL render an `Input` field of type `number` bound to `sendDeviceInfoExpiryDays`.
8. THE Settings_Form SHALL render a `Switch` component bound to `unsecureEnrollment`.
9. THE Settings_Form SHALL render a `Switch` component bound to `deviceFastSearch`.
10. THE Settings_Form SHALL render a "Save Settings" `Button` of type `submit`.

### Requirement 3: Save Settings

**User Story:** As an MDM administrator, I want to save my changes, so that the updated settings are persisted to the backend.

#### Acceptance Criteria

1. WHEN a user clicks the "Save Settings" button and the form is valid, THE Settings_Service SHALL call `PUT /rest/private/settings` with the current form values.
2. WHILE the PUT request is in flight, THE Settings_Form SHALL disable the "Save Settings" button and display a loading indicator on it.
3. WHEN the PUT request succeeds, THE Settings_Page SHALL display a success Toast with a confirmation message.
4. IF the PUT request fails, THEN THE Settings_Page SHALL display an error Toast with a descriptive message and keep the form populated with the submitted values.
5. WHEN the PUT request completes (success or failure), THE Settings_Form SHALL re-enable the "Save Settings" button.

### Requirement 4: Form Validation

**User Story:** As an MDM administrator, I want the settings form to validate my input, so that I cannot submit incomplete or invalid data.

#### Acceptance Criteria

1. THE Settings_Form SHALL require the `customerName` field; WHEN the `customerName` field is empty on submit, THE Settings_Form SHALL display a validation error on the `customerName` field and prevent submission.
2. THE Settings_Form SHALL require the `passwordLength` field; WHEN the `passwordLength` field is empty or contains a value less than `1` on submit, THE Settings_Form SHALL display a validation error on the `passwordLength` field and prevent submission.
3. THE Settings_Form SHALL require the `sendDeviceInfoExpiryDays` field; WHEN the `sendDeviceInfoExpiryDays` field is empty or contains a value less than `1` on submit, THE Settings_Form SHALL display a validation error on the `sendDeviceInfoExpiryDays` field and prevent submission.
4. WHEN a validation error is present, THE Settings_Form SHALL display the error message adjacent to the invalid field using the shadcn/ui `FormMessage` component.
5. WHEN a validation error is present, THE Settings_Form SHALL NOT call `PUT /rest/private/settings`.

### Requirement 5: Toast Notifications

**User Story:** As an MDM administrator, I want clear feedback after saving settings, so that I know whether my changes were applied successfully.

#### Acceptance Criteria

1. WHEN the PUT request succeeds, THE Settings_Page SHALL display a Toast with a title of "Settings saved" and a success visual style.
2. IF the PUT request fails with a network or server error, THEN THE Settings_Page SHALL display a Toast with a title of "Failed to save settings" and a destructive visual style.
3. THE Toast SHALL be dismissible by the user.
4. THE Toast SHALL auto-dismiss after a reasonable timeout without user interaction.

### Requirement 6: Navigation Integration

**User Story:** As an MDM administrator, I want a Settings link in the sidebar, so that I can navigate to the settings page from anywhere in the application.

#### Acceptance Criteria

1. THE navigation sidebar SHALL include a "Settings" entry.
2. WHEN a user clicks the "Settings" sidebar entry, THE application SHALL navigate to `/settings`.
3. WHEN the current route is `/settings`, THE navigation sidebar SHALL render the "Settings" entry in its active/selected visual state.

### Requirement 7: Settings Service Layer

**User Story:** As a frontend developer, I want a dedicated service module for settings API calls, so that all backend communication is centralized and testable.

#### Acceptance Criteria

1. THE Settings_Service SHALL expose a `getSettings(): Promise<Settings>` function that calls `GET /rest/private/settings`.
2. THE Settings_Service SHALL expose an `updateSettings(data: SettingsPayload): Promise<Settings>` function that calls `PUT /rest/private/settings`.
3. THE Settings_Service SHALL use the shared `apiClient` instance for all HTTP calls so that the `X-Auth-Token` header is attached automatically.
4. IF any Settings_Service call receives a non-2xx response, THEN THE Settings_Service SHALL propagate the error to the caller without swallowing it.

### Requirement 8: Data Types

**User Story:** As a frontend developer, I want well-typed TypeScript interfaces for all settings-related data, so that the compiler catches integration errors early.

#### Acceptance Criteria

1. THE frontend SHALL define a `Settings` interface with fields: `id: number`, `customerName: string`, `createNewDevices: boolean`, `newDeviceConfigurationId: number | null`, `language: string`, `passwordLength: number`, `passwordStrength: number`, `sendDeviceInfoExpiryDays: number`, `unsecureEnrollment: boolean`, `deviceFastSearch: boolean`.
2. THE frontend SHALL define a `SettingsPayload` interface that omits the `id` field from `Settings` and marks all fields as required.
3. WHEN the backend returns a boolean field that is absent or null, THE frontend SHALL treat it as `false` rather than throwing a runtime error.
4. WHEN the backend returns a numeric field that is absent or null, THE frontend SHALL treat it as `0` rather than throwing a runtime error.
