# Requirements Document

## Introduction

The Users Management feature provides a dedicated page at `/users` for viewing, creating, editing, and deleting HMDM administrator accounts. Each user has a login, display name, email, assigned role, and access flags. The page renders a table of all users with key metadata, supports creating and editing users via a dialog form, and allows deletion with a confirmation prompt. Role options are fetched from the backend and presented via a Select component. All data is fetched from the existing HMDM backend REST API using the shared `apiClient` (base `/rest`, `X-Auth-Token` header injected automatically). The UI is built exclusively with shadcn/ui components inside the existing React + Vite + TypeScript frontend. The `/users` route is added to the navigation sidebar.

## Glossary

- **User**: An HMDM administrator account identified by a unique numeric `id`, a `login` string, a `name` string, an `email` string, a `role` object (`id`, `name`), an `allDevicesAvailable` boolean, and a `customerStaff` boolean.
- **Role**: A named permission level assigned to a user, identified by a numeric `id` and a `name` string, returned by `GET /rest/private/roles`.
- **User_List**: The collection of `User` records returned by `GET /rest/private/users`.
- **Users_Page**: The React page component rendered at the `/users` route.
- **User_Service**: The frontend service module (`src/features/users/userService.ts`) that wraps all `apiClient` calls for user-related endpoints.
- **User_Form**: The shadcn/ui `Dialog` containing a `Form` used to create or edit a user.
- **Delete_Dialog**: The shadcn/ui `AlertDialog` confirmation modal shown before a user is permanently deleted.
- **Role_Select**: The shadcn/ui `Select` component used to assign a role to a user within the User_Form.

---

## Requirements

### Requirement 1: Users List Page

**User Story:** As an MDM administrator, I want to see all user accounts in a table, so that I can survey and manage who has access to the system.

#### Acceptance Criteria

1. WHEN a user navigates to `/users`, THE Users_Page SHALL render a full-page layout consistent with the existing navigation shell.
2. WHEN the Users_Page mounts, THE User_Service SHALL call `GET /rest/private/users`.
3. WHEN the API response is received, THE Users_Page SHALL display a shadcn/ui `Table` with columns: Login, Name, Email, Role, Status.
4. WHILE the API request is in flight, THE Users_Page SHALL display a loading skeleton or spinner in place of the table.
5. IF the API request fails, THEN THE Users_Page SHALL display an error message with a retry action.
6. WHEN the user list is empty, THE Users_Page SHALL display an empty-state message indicating no users are available.

### Requirement 2: Create User

**User Story:** As an MDM administrator, I want to create a new user account, so that I can grant access to additional administrators.

#### Acceptance Criteria

1. THE Users_Page SHALL render an "Add User" `Button` above the table.
2. WHEN a user clicks the "Add User" button, THE Users_Page SHALL open the User_Form dialog in create mode with all fields empty.
3. WHEN the User_Form is in create mode, THE User_Form SHALL render a password field.
4. WHEN a user submits the User_Form with valid data in create mode, THE User_Service SHALL call `POST /rest/private/users` with the form values.
5. WHILE the create request is in flight, THE User_Form SHALL disable the submit button and show a loading indicator.
6. WHEN the create request succeeds, THE Users_Page SHALL close the User_Form and refresh the user list.
7. IF the create request fails, THEN THE User_Form SHALL display an error message and keep the dialog open.
8. WHEN a user cancels the User_Form, THE Users_Page SHALL close the dialog without making any API call.

### Requirement 3: Edit User

**User Story:** As an MDM administrator, I want to edit an existing user's details, so that I can update their name, email, or role without recreating the account.

#### Acceptance Criteria

1. WHEN a user activates the edit action for a user (via a row action menu or button), THE Users_Page SHALL open the User_Form dialog in edit mode pre-populated with the user's current `login`, `name`, `email`, `role`, `allDevicesAvailable`, and `customerStaff` values.
2. WHEN the User_Form is in edit mode, THE User_Form SHALL NOT render a password field.
3. WHEN the User_Form is in edit mode and the user submits valid data, THE User_Service SHALL call `PUT /rest/private/users/{id}` with the updated values.
4. WHILE the update request is in flight, THE User_Form SHALL disable the submit button and show a loading indicator.
5. WHEN the update request succeeds, THE Users_Page SHALL close the User_Form and refresh the user list.
6. IF the update request fails, THEN THE User_Form SHALL display an error message and keep the dialog open.

### Requirement 4: Delete User

**User Story:** As an MDM administrator, I want to delete a user account, so that I can revoke access for administrators who no longer need it.

#### Acceptance Criteria

1. WHEN a user activates the delete action for a user (via a row action menu or button), THE Users_Page SHALL open the Delete_Dialog identifying the user by login and name.
2. WHEN a user confirms deletion in the Delete_Dialog, THE User_Service SHALL call `DELETE /rest/private/users/{id}`.
3. WHILE the delete request is in flight, THE Delete_Dialog SHALL disable the confirm button and show a loading indicator.
4. WHEN the delete request succeeds, THE Users_Page SHALL close the Delete_Dialog and refresh the user list.
5. IF the delete request fails, THEN THE Delete_Dialog SHALL display an error message and keep the dialog open.
6. WHEN a user cancels the Delete_Dialog, THE Users_Page SHALL close the dialog without making any API call.

### Requirement 5: Role Assignment

**User Story:** As an MDM administrator, I want to assign a role to a user from a list of available roles, so that I can control what each administrator is permitted to do.

#### Acceptance Criteria

1. WHEN the User_Form opens (create or edit mode), THE User_Service SHALL call `GET /rest/private/roles` to fetch available roles.
2. WHILE the roles request is in flight, THE Role_Select SHALL display a loading state and be disabled.
3. WHEN the roles response is received, THE Role_Select SHALL populate its options with the `name` of each role and use the role `id` as the option value.
4. WHEN the User_Form is in edit mode, THE Role_Select SHALL pre-select the role matching the user's current `role.id`.
5. THE User_Form SHALL require a role selection; WHEN no role is selected on submit, THE User_Form SHALL display a validation error on the Role_Select and prevent submission.

### Requirement 6: User Form Validation

**User Story:** As an MDM administrator, I want the user form to validate my input, so that I cannot submit incomplete or invalid data.

#### Acceptance Criteria

1. THE User_Form SHALL require the `login` field; WHEN the `login` field is empty on submit, THE User_Form SHALL display a validation error on the `login` field and prevent submission.
2. THE User_Form SHALL require the `name` field; WHEN the `name` field is empty on submit, THE User_Form SHALL display a validation error on the `name` field and prevent submission.
3. THE User_Form SHALL require the `email` field; WHEN the `email` field contains a value that does not match the pattern `[^@]+@[^.]+\..+` on submit, THE User_Form SHALL display a validation error on the `email` field and prevent submission.
4. WHEN the User_Form is in create mode, THE User_Form SHALL require the `password` field; WHEN the `password` field is empty on submit, THE User_Form SHALL display a validation error on the `password` field and prevent submission.
5. WHEN a validation error is present, THE User_Form SHALL display the error message adjacent to the invalid field using the shadcn/ui `FormMessage` component.

### Requirement 7: Navigation Integration

**User Story:** As an MDM administrator, I want a Users link in the sidebar, so that I can navigate to the users page from anywhere in the application.

#### Acceptance Criteria

1. THE navigation sidebar SHALL include a "Users" entry.
2. WHEN a user clicks the "Users" sidebar entry, THE application SHALL navigate to `/users`.
3. WHEN the current route is `/users`, THE navigation sidebar SHALL render the "Users" entry in its active/selected visual state.

### Requirement 8: User Service Layer

**User Story:** As a frontend developer, I want a dedicated service module for user API calls, so that all backend communication is centralized and testable.

#### Acceptance Criteria

1. THE User_Service SHALL expose a `getUsers(): Promise<User[]>` function that calls `GET /rest/private/users`.
2. THE User_Service SHALL expose a `createUser(data: UserPayload): Promise<User>` function that calls `POST /rest/private/users`.
3. THE User_Service SHALL expose an `updateUser(id: number, data: UserPayload): Promise<User>` function that calls `PUT /rest/private/users/{id}`.
4. THE User_Service SHALL expose a `deleteUser(id: number): Promise<void>` function that calls `DELETE /rest/private/users/{id}`.
5. THE User_Service SHALL expose a `getRoles(): Promise<Role[]>` function that calls `GET /rest/private/roles`.
6. THE User_Service SHALL use the shared `apiClient` instance for all HTTP calls so that the `X-Auth-Token` header is attached automatically.
7. IF any User_Service call receives a non-2xx response, THEN THE User_Service SHALL propagate the error to the caller without swallowing it.

### Requirement 9: Data Types

**User Story:** As a frontend developer, I want well-typed TypeScript interfaces for all user-related data, so that the compiler catches integration errors early.

#### Acceptance Criteria

1. THE frontend SHALL define a `Role` interface with fields: `id: number`, `name: string`.
2. THE frontend SHALL define a `User` interface with fields: `id: number`, `login: string`, `name: string`, `email: string`, `role: Role`, `allDevicesAvailable: boolean`, `customerStaff: boolean`.
3. THE frontend SHALL define a `UserPayload` interface with fields: `login: string`, `name: string`, `email: string`, `password?: string`, `roleId: number`, `allDevicesAvailable: boolean`, `customerStaff: boolean`.
4. WHEN the backend returns a boolean field that is absent or null, THE frontend SHALL treat it as `false` rather than throwing a runtime error.
5. WHEN the backend returns the `role` field as absent or null, THE frontend SHALL treat it as `null` and render an empty cell in the Role column rather than throwing a runtime error.
