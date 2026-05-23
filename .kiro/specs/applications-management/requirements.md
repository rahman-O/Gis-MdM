# Requirements Document

## Introduction

The Applications Management feature provides a dedicated page at `/applications` for viewing, uploading, editing, and deleting Android applications managed by Headwind MDM. An application entry represents an APK that can be deployed to enrolled devices, identified by a package name and version. The page renders a table of all applications with key metadata, supports uploading new APKs or adding application entries via a dialog form, allows editing application metadata, and allows deletion with a confirmation prompt. All data is fetched from the existing HMDM backend REST API using the shared `apiClient` (base `/rest`, `X-Auth-Token` header injected automatically). The UI is built exclusively with shadcn/ui components inside the existing React + Vite + TypeScript frontend. The `/applications` route is added to the navigation sidebar.

## Glossary

- **Application**: An Android app entry managed by Headwind MDM, identified by a unique numeric `id`, a human-readable `name`, a `pkg` (Android package name), a `version` string, a `url` (APK download URL), a `system` boolean flag, a `showIcon` boolean, a `runAfterInstall` boolean, and a `runAtBoot` boolean.
- **Application_List**: The collection of `Application` records returned by `GET /rest/private/applications`.
- **Applications_Page**: The React page component rendered at the `/applications` route.
- **Application_Service**: The frontend service module (`src/features/applications/applicationService.ts`) that wraps all `apiClient` calls for application-related endpoints.
- **Application_Form**: The shadcn/ui `Dialog` containing a `Form` used to upload a new APK or add/edit an application entry.
- **Delete_Dialog**: The shadcn/ui `AlertDialog` confirmation modal shown before an application is permanently deleted.
- **APK_Upload**: The file input mechanism within the Application_Form that sends an APK binary to `POST /rest/private/web-ui-files` and returns a download URL.
- **System_Flag**: A boolean field on an `Application` indicating whether the app is a system application pre-installed on the device.

---

## Requirements

### Requirement 1: Applications List Page

**User Story:** As an MDM administrator, I want to see all managed applications in a table, so that I can survey the application catalog and manage individual entries.

#### Acceptance Criteria

1. WHEN a user navigates to `/applications`, THE Applications_Page SHALL render a full-page layout consistent with the existing navigation shell.
2. WHEN the Applications_Page mounts, THE Application_Service SHALL call `GET /rest/private/applications`.
3. WHEN the API response is received, THE Applications_Page SHALL display a shadcn/ui `Table` with columns: Name, Package, Version, System.
4. WHEN the `system` field of an application is `true`, THE Applications_Page SHALL render a visual indicator (such as a `Badge` or checkmark icon) in the System column; WHEN `system` is `false`, THE Applications_Page SHALL render an empty or negative indicator.
5. WHILE the API request is in flight, THE Applications_Page SHALL display a loading skeleton or spinner in place of the table.
6. IF the API request fails, THEN THE Applications_Page SHALL display an error message with a retry action.
7. WHEN the application list is empty, THE Applications_Page SHALL display an empty-state message indicating no applications are available.

### Requirement 2: Upload / Add Application

**User Story:** As an MDM administrator, I want to upload a new APK or add an application entry manually, so that I can make new apps available for deployment to devices.

#### Acceptance Criteria

1. THE Applications_Page SHALL render an "Add Application" `Button` above the table.
2. WHEN a user clicks the "Add Application" button, THE Applications_Page SHALL open the Application_Form dialog in create mode with all fields empty.
3. WHEN a user selects an APK file in the Application_Form, THE Application_Service SHALL call `POST /rest/private/web-ui-files` with the file as `multipart/form-data` to upload the APK.
4. WHILE the APK upload request is in flight, THE Application_Form SHALL display an upload progress indicator and disable the file input.
5. WHEN the APK upload succeeds, THE Application_Form SHALL populate the `url` field with the returned download URL and auto-fill the `name`, `pkg`, and `version` fields if the API response includes them.
6. WHEN a user submits the Application_Form with valid data, THE Application_Service SHALL call `POST /rest/private/applications` with the form values.
7. WHILE the create request is in flight, THE Application_Form SHALL disable the submit button and show a loading indicator.
8. WHEN the create request succeeds, THE Applications_Page SHALL close the Application_Form and refresh the application list.
9. IF the create request fails, THEN THE Application_Form SHALL display an error message and keep the dialog open.
10. WHEN a user cancels the Application_Form, THE Applications_Page SHALL close the dialog without making any API call.

### Requirement 3: Edit Application

**User Story:** As an MDM administrator, I want to edit an existing application's metadata, so that I can correct names, versions, or behavioral flags without re-uploading the APK.

#### Acceptance Criteria

1. WHEN a user activates the edit action for an application (via a row action menu or button), THE Applications_Page SHALL open the Application_Form dialog in edit mode pre-populated with the application's current `name`, `pkg`, `version`, `url`, `system`, `showIcon`, `runAfterInstall`, and `runAtBoot` values.
2. WHEN the Application_Form is in edit mode and the user submits valid data, THE Application_Service SHALL call `PUT /rest/private/applications/{id}` with the updated values.
3. WHILE the update request is in flight, THE Application_Form SHALL disable the submit button and show a loading indicator.
4. WHEN the update request succeeds, THE Applications_Page SHALL close the Application_Form and refresh the application list.
5. IF the update request fails, THEN THE Application_Form SHALL display an error message and keep the dialog open.

### Requirement 4: Delete Application

**User Story:** As an MDM administrator, I want to delete an application from the catalog, so that I can remove obsolete or unwanted apps.

#### Acceptance Criteria

1. WHEN a user activates the delete action for an application (via a row action menu or button), THE Applications_Page SHALL open the Delete_Dialog identifying the application by name and package.
2. WHEN a user confirms deletion in the Delete_Dialog, THE Application_Service SHALL call `DELETE /rest/private/applications/{id}`.
3. WHILE the delete request is in flight, THE Delete_Dialog SHALL disable the confirm button and show a loading indicator.
4. WHEN the delete request succeeds, THE Applications_Page SHALL close the Delete_Dialog and refresh the application list.
5. IF the delete request fails, THEN THE Delete_Dialog SHALL display an error message and keep the dialog open.
6. WHEN a user cancels the Delete_Dialog, THE Applications_Page SHALL close the dialog without making any API call.

### Requirement 5: Application Form Validation

**User Story:** As an MDM administrator, I want the application form to validate my input, so that I cannot submit incomplete or invalid data.

#### Acceptance Criteria

1. THE Application_Form SHALL require the `name` field; WHEN the `name` field is empty on submit, THE Application_Form SHALL display a validation error on the `name` field and prevent submission.
2. THE Application_Form SHALL require the `pkg` field; WHEN the `pkg` field is empty on submit, THE Application_Form SHALL display a validation error on the `pkg` field and prevent submission.
3. THE Application_Form SHALL require the `version` field; WHEN the `version` field is empty on submit, THE Application_Form SHALL display a validation error on the `version` field and prevent submission.
4. THE Application_Form SHALL require the `url` field; WHEN the `url` field is empty on submit, THE Application_Form SHALL display a validation error on the `url` field and prevent submission.
5. THE Application_Form SHALL render the `system`, `showIcon`, `runAfterInstall`, and `runAtBoot` fields as shadcn/ui `Switch` components defaulting to `false`.
6. WHEN a validation error is present, THE Application_Form SHALL display the error message adjacent to the invalid field using the shadcn/ui `FormMessage` component.

### Requirement 6: APK File Upload

**User Story:** As an MDM administrator, I want to upload an APK file directly from the form, so that the system stores the binary and provides a download URL automatically.

#### Acceptance Criteria

1. THE Application_Form SHALL render a file input that accepts only files with the `.apk` extension.
2. WHEN a user selects a non-APK file, THE Application_Form SHALL display a validation error on the file input and prevent the upload from starting.
3. WHEN the APK upload request is in flight, THE Application_Form SHALL show upload progress and prevent the user from selecting a different file.
4. IF the APK upload request fails, THEN THE Application_Form SHALL display an error message on the file input and allow the user to retry by selecting the file again.
5. WHEN the APK upload succeeds, THE Application_Form SHALL set the `url` field value to the URL returned by `POST /rest/private/web-ui-files`.
6. THE Application_Form SHALL allow a user to provide a `url` value manually without uploading an APK file, so that externally hosted APKs can be referenced.

### Requirement 7: Navigation Integration

**User Story:** As an MDM administrator, I want an Applications link in the sidebar, so that I can navigate to the applications page from anywhere in the application.

#### Acceptance Criteria

1. THE navigation sidebar SHALL include an "Applications" entry.
2. WHEN a user clicks the "Applications" sidebar entry, THE application SHALL navigate to `/applications`.
3. WHEN the current route is `/applications`, THE navigation sidebar SHALL render the "Applications" entry in its active/selected visual state.

### Requirement 8: Application Service Layer

**User Story:** As a frontend developer, I want a dedicated service module for application API calls, so that all backend communication is centralized and testable.

#### Acceptance Criteria

1. THE Application_Service SHALL expose a `getApplications(): Promise<Application[]>` function that calls `GET /rest/private/applications`.
2. THE Application_Service SHALL expose a `createApplication(data: ApplicationPayload): Promise<Application>` function that calls `POST /rest/private/applications`.
3. THE Application_Service SHALL expose an `updateApplication(id: number, data: ApplicationPayload): Promise<Application>` function that calls `PUT /rest/private/applications/{id}`.
4. THE Application_Service SHALL expose a `deleteApplication(id: number): Promise<void>` function that calls `DELETE /rest/private/applications/{id}`.
5. THE Application_Service SHALL expose an `uploadApk(file: File): Promise<ApkUploadResponse>` function that calls `POST /rest/private/web-ui-files` with `Content-Type: multipart/form-data`.
6. THE Application_Service SHALL use the shared `apiClient` instance for all HTTP calls so that the `X-Auth-Token` header is attached automatically.
7. IF any Application_Service call receives a non-2xx response, THEN THE Application_Service SHALL propagate the error to the caller without swallowing it.

### Requirement 9: Data Types

**User Story:** As a frontend developer, I want well-typed TypeScript interfaces for all application-related data, so that the compiler catches integration errors early.

#### Acceptance Criteria

1. THE frontend SHALL define an `Application` interface with fields: `id: number`, `name: string`, `pkg: string`, `version: string`, `url: string`, `system: boolean`, `showIcon: boolean`, `runAfterInstall: boolean`, `runAtBoot: boolean`.
2. THE frontend SHALL define an `ApplicationPayload` interface with fields: `name: string`, `pkg: string`, `version: string`, `url: string`, `system: boolean`, `showIcon: boolean`, `runAfterInstall: boolean`, `runAtBoot: boolean`.
3. THE frontend SHALL define an `ApkUploadResponse` interface with at minimum the field `url: string` representing the download URL of the uploaded APK.
4. WHEN the backend returns a boolean field that is absent or null, THE frontend SHALL treat it as `false` rather than throwing a runtime error.
