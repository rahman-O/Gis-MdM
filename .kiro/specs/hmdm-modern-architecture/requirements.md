# Requirements Document

## Introduction

This feature introduces a clean modern architecture for the Headwind MDM (hmdm-server) project by separating the existing Java backend from a brand-new React/TypeScript frontend. The legacy JSP/AngularJS frontend is replaced entirely. The backend is relocated into a `/backend` directory and treated as a standalone REST API server. A new `/frontend` directory contains a React + Vite + TypeScript application that communicates with the backend exclusively via REST API calls. The first feature implemented in the new frontend is user authentication (login).

## Glossary

- **Workspace**: The root directory containing `/backend`, `/frontend`, and `/docs` subdirectories.
- **Backend**: The existing Java/Maven hmdm-server application relocated to `/backend`, serving REST API endpoints.
- **Frontend**: The new React + Vite + TypeScript application located in `/frontend`.
- **Auth_Service**: The frontend service module responsible for communicating with the backend authentication API.
- **Login_Page**: The React component that renders the login form and handles user interaction.
- **Token_Store**: The browser-side storage mechanism (localStorage) used to persist the authentication token.
- **API_Client**: The Axios-based HTTP client configured with the backend base URL.
- **Router**: The React Router instance managing client-side navigation between pages.
- **Dashboard_Page**: The protected React page that authenticated users are redirected to after login.
- **Auth_Guard**: The frontend route protection mechanism that redirects unauthenticated users to the Login_Page.
- **shadcn/ui**: The component library used for all UI elements, built on Radix UI primitives with Tailwind CSS styling.

---

## Requirements

### Requirement 1: Workspace Structure Separation

**User Story:** As a developer, I want the backend and frontend to live in separate top-level directories, so that each can be developed, built, and deployed independently.

#### Acceptance Criteria

1. THE Workspace SHALL contain a `/backend` directory holding the complete, unmodified hmdm-server Java/Maven project.
2. THE Workspace SHALL contain a `/frontend` directory holding the new React + Vite + TypeScript application.
3. WHERE documentation is provided, THE Workspace SHALL contain a `/docs` directory for project documentation.
4. THE Backend SHALL be runnable as a standalone server from within the `/backend` directory without any dependency on the `/frontend` directory.
5. THE Frontend SHALL be runnable as a standalone development server from within the `/frontend` directory without any dependency on the `/backend` source code.

---

### Requirement 2: Backend Preservation

**User Story:** As a developer, I want the backend code to remain untouched after relocation, so that existing API behavior is preserved.

#### Acceptance Criteria

1. THE Backend SHALL expose all existing REST API endpoints unchanged after relocation to `/backend`.
2. THE Backend SHALL retain its original Maven project structure (`pom.xml`, `src/`, module directories) inside `/backend`.
3. IF any backend source file is modified during relocation, THEN THE Backend SHALL fail its build, signaling an unintended change.
4. THE Backend SHALL serve its REST API on port 8080 by default, consistent with its original configuration.

---

### Requirement 3: Frontend Project Initialization

**User Story:** As a developer, I want a modern frontend scaffold, so that I can build new UI features with current tooling.

#### Acceptance Criteria

1. THE Frontend SHALL be initialized using Vite with the React + TypeScript template.
2. THE Frontend SHALL use TypeScript as the sole language for all application source files under `src/`.
3. THE Frontend SHALL include ESLint configured for React and TypeScript linting rules.
4. THE Frontend SHALL include Prettier configured for consistent code formatting.
5. THE Frontend SHALL include Axios as the HTTP client library for all API communication.
6. THE Frontend SHALL NOT contain any JSP files, jQuery code, or AngularJS code from the legacy application.
7. THE Frontend SHALL use shadcn/ui as the primary component library, built on top of Radix UI primitives and styled with Tailwind CSS.
8. THE Frontend SHALL include Tailwind CSS configured as the styling solution, required by shadcn/ui.

---

### Requirement 4: Frontend Directory Architecture

**User Story:** As a developer, I want a well-defined feature-based directory structure, so that the codebase scales cleanly as new features are added.

#### Acceptance Criteria

1. THE Frontend SHALL organize application source code under `src/` with the following top-level directories: `app/`, `features/`, `shared/`, and `services/`.
2. THE Frontend SHALL contain a `src/features/auth/` directory for all authentication-related components, hooks, and types.
3. THE Frontend SHALL contain a `src/features/devices/` directory reserved for device management features.
4. THE Frontend SHALL contain a `src/features/users/` directory reserved for user management features.
5. THE Frontend SHALL contain a `src/features/dashboard/` directory for the main dashboard feature.
6. THE Frontend SHALL contain a `src/shared/components/` directory for reusable UI components shared across features.
7. THE Frontend SHALL contain a `src/shared/hooks/` directory for reusable React hooks.
8. THE Frontend SHALL contain a `src/shared/utils/` directory for utility functions.
9. THE Frontend SHALL contain a `src/shared/ui/` directory for base UI primitives.
10. THE Frontend SHALL contain a `src/services/` directory for API service modules.
11. THE Frontend SHALL have `src/main.tsx` as the application entry point.

---

### Requirement 5: API Client Configuration

**User Story:** As a developer, I want a centralized API client, so that all HTTP requests share a consistent base URL and configuration.

#### Acceptance Criteria

1. THE API_Client SHALL be configured with a base URL of `http://localhost:8080/api`.
2. WHEN a request is made, THE API_Client SHALL include the `Content-Type: application/json` header by default.
3. WHEN an authentication token is present in the Token_Store, THE API_Client SHALL attach it as an `X-Auth-Token` request header on every outgoing request.
4. IF the backend returns an HTTP 401 response, THEN THE API_Client SHALL clear the Token_Store and redirect the user to the Login_Page.
5. THE API_Client SHALL be the sole mechanism used by all service modules to communicate with the Backend.

---

### Requirement 6: Authentication — Login

**User Story:** As an MDM administrator, I want to log in with my credentials, so that I can access the management dashboard.

#### Acceptance Criteria

1. THE Login_Page SHALL render a form containing a username input field, a password input field, and a submit button.
2. WHEN the user submits the login form, THE Auth_Service SHALL send a POST request to `/api/public/auth/login` with a JSON body containing `login` and `password` fields.
3. WHEN the backend returns a successful response containing an `authToken`, THE Token_Store SHALL persist the token in localStorage under the key `hmdm_auth_token`.
4. WHEN the token is successfully stored, THE Router SHALL redirect the user to the Dashboard_Page at `/dashboard`.
5. IF the backend returns an error response to the login request, THEN THE Login_Page SHALL display a human-readable error message to the user without clearing the input fields.
6. WHILE the login request is in flight, THE Login_Page SHALL display a loading indicator and disable the submit button to prevent duplicate submissions.
7. THE Login_Page SHALL be accessible at the `/login` route.
8. WHEN an unauthenticated user navigates to a protected route, THE Auth_Guard SHALL redirect the user to the Login_Page.

---

### Requirement 7: Authentication — Logout

**User Story:** As an MDM administrator, I want to log out, so that my session is terminated securely.

#### Acceptance Criteria

1. WHEN the user triggers a logout action, THE Auth_Service SHALL send a POST request to `/api/public/auth/logout`.
2. WHEN the logout request completes (success or failure), THE Token_Store SHALL remove the `hmdm_auth_token` entry from localStorage.
3. WHEN the token is cleared, THE Router SHALL redirect the user to the Login_Page.

---

### Requirement 8: Dashboard Placeholder

**User Story:** As an MDM administrator, I want to see a dashboard after login, so that I have a landing page to navigate from.

#### Acceptance Criteria

1. THE Dashboard_Page SHALL be accessible at the `/dashboard` route.
2. WHEN an unauthenticated user navigates to `/dashboard`, THE Auth_Guard SHALL redirect the user to `/login`.
3. THE Dashboard_Page SHALL display the currently authenticated user's login name retrieved from the Token_Store or application state.
4. THE Dashboard_Page SHALL provide a logout control that triggers the logout flow defined in Requirement 7.

---

### Requirement 9: Token Persistence and Session Restoration

**User Story:** As an MDM administrator, I want my session to persist across browser refreshes, so that I do not have to log in repeatedly.

#### Acceptance Criteria

1. WHEN the Frontend application initializes, THE Auth_Service SHALL read the `hmdm_auth_token` key from localStorage.
2. IF a valid token is found on initialization, THEN THE Auth_Guard SHALL allow navigation to protected routes without requiring a new login.
3. IF no token is found on initialization, THEN THE Auth_Guard SHALL redirect the user to the Login_Page when a protected route is accessed.
