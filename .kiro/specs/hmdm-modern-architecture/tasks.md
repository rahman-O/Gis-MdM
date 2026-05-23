# Implementation Plan: HMDM Modern Architecture

## Overview

Restructure the workspace by relocating the existing Java backend to `/backend`, scaffold a new React + Vite + TypeScript frontend in `/frontend`, and implement token-based authentication (login/logout) with session persistence using shadcn/ui, Axios, React Router v6, and React Context.

## Tasks

- [x] 1. Restructure workspace — move hmdm-server to backend/
  - Move the contents of `hmdm-server/` into a new `backend/` directory at the workspace root
  - Verify `backend/pom.xml` exists and the Maven module structure is intact
  - Create a `docs/` directory at the workspace root
  - _Requirements: 1.1, 1.3, 2.1, 2.2_

- [x] 2. Scaffold the frontend project
  - [x] 2.1 Initialize Vite + React + TypeScript project in `frontend/`
    - Run `npm create vite@latest frontend -- --template react-ts`
    - Install dependencies: `axios@^1.7.0`, `react-router-dom@^6.26.0`, `react-hook-form@^7.53.0`, `@hookform/resolvers@^3.9.0`, `zod@^3.23.0`
    - Install dev dependencies: `vitest@^2.1.0`, `@vitest/coverage-v8@^2.1.0`, `@testing-library/react@^16.0.0`, `@testing-library/jest-dom@^6.5.0`, `@testing-library/user-event@^14.5.0`, `fast-check@^3.22.0`, `msw@^2.4.0`, `eslint@^9.0.0`, `prettier@^3.3.0`
    - _Requirements: 3.1, 3.3, 3.4, 3.5_

  - [x] 2.2 Configure Tailwind CSS and shadcn/ui
    - Install and configure Tailwind CSS: `tailwindcss@^3.4.0`, `postcss`, `autoprefixer`
    - Create `tailwind.config.ts` with content paths covering `./src/**/*.{ts,tsx}`
    - Run `npx shadcn@latest init` to generate `components.json` and configure the `src/shared/ui/` output directory
    - Add shadcn/ui components: `button`, `card`, `form`, `input`, `label`
    - _Requirements: 3.7, 3.8_

  - [x] 2.3 Configure Vite, TypeScript, ESLint, and Prettier
    - Update `vite.config.ts` to add Vitest config block (`globals: true`, `environment: 'jsdom'`, `setupFiles: ['./src/test/setup.ts']`)
    - Create `src/test/setup.ts` importing `@testing-library/jest-dom`
    - Create `.eslintrc.cjs` with `plugin:react/recommended`, `plugin:@typescript-eslint/recommended` rules
    - Create `.prettierrc` with standard formatting options
    - _Requirements: 3.3, 3.4_

- [x] 3. Create the frontend directory structure and core types
  - Create empty placeholder files to establish the directory layout:
    - `src/app/App.tsx`, `src/app/providers.tsx`
    - `src/features/auth/types.ts`, `src/features/auth/LoginPage.tsx`, `src/features/auth/useLogin.ts`
    - `src/features/dashboard/DashboardPage.tsx`
    - `src/features/devices/.gitkeep`, `src/features/users/.gitkeep`
    - `src/shared/components/AuthGuard.tsx`
    - `src/shared/hooks/useAuth.ts`
    - `src/shared/utils/tokenStorage.ts`
    - `src/services/apiClient.ts`, `src/services/authService.ts`
  - _Requirements: 4.1–4.11_

- [x] 4. Implement token storage utilities
  - [x] 4.1 Implement `src/shared/utils/tokenStorage.ts`
    - Export `getToken()`, `setToken(token: string)`, `clearToken()` using `localStorage` key `hmdm_auth_token`
    - _Requirements: 6.3, 9.1_

  - [ ]* 4.2 Write property test for tokenStorage (Property 4)
    - **Property 4: Login stores the token (round trip)**
    - **Validates: Requirements 6.3**
    - Use `fc.string({ minLength: 1 })` to generate tokens; assert `getToken()` equals the stored value after `setToken()`
    - Assert `getToken()` returns `null` after `clearToken()`

- [x] 5. Implement the API client
  - [x] 5.1 Implement `src/services/apiClient.ts`
    - Create Axios instance with `baseURL: 'http://localhost:8080/api'` and `Content-Type: application/json` default header
    - Add request interceptor: read token via `getToken()`; if present, set `X-Auth-Token` header
    - Add response interceptor: on HTTP 401, call `clearToken()` then redirect via `window.location.href = '/login'`
    - _Requirements: 5.1, 5.2, 5.3, 5.4_

  - [ ]* 5.2 Write property test for request interceptor (Property 2)
    - **Property 2: Request headers are set correctly by the API client interceptor**
    - **Validates: Requirements 5.2, 5.3**
    - Use `fc.string({ minLength: 1 })` for token values; assert `X-Auth-Token` equals the token when present; assert `Content-Type` is always `application/json`

  - [ ]* 5.3 Write property test for 401 interceptor (Property 3)
    - **Property 3: 401 response clears the token**
    - **Validates: Requirements 5.4**
    - Use `msw` to return 401 for any request; assert `localStorage.getItem('hmdm_auth_token')` is `null` after the interceptor runs

- [x] 6. Implement the auth service
  - [x] 6.1 Implement `src/services/authService.ts`
    - Define `LoginRequest` and `LoginResponse` types (or import from `features/auth/types.ts`)
    - Implement `login(credentials: LoginRequest): Promise<LoginResponse>` — POST to `/public/auth/login` via `apiClient`
    - Implement `logout(): Promise<void>` — POST to `/public/auth/logout` via `apiClient`
    - _Requirements: 6.2, 7.1, 5.5_

  - [x] 6.2 Define auth types in `src/features/auth/types.ts`
    - Export `LoginRequest { login: string; password: string }`
    - Export `LoginResponse { authToken: string }`
    - _Requirements: 6.2_

  - [ ]* 6.3 Write unit tests for authService
    - Use `msw` to mock `POST /api/public/auth/login` and `POST /api/public/auth/logout`
    - Assert `login()` returns the `authToken` from the response
    - Assert `logout()` calls `POST /api/public/auth/logout`
    - _Requirements: 6.2, 7.1_

- [x] 7. Implement AuthContext and AuthProvider
  - [x] 7.1 Implement `src/app/providers.tsx`
    - Define `AuthContextValue { token, username, login, logout }`
    - `AuthProvider` reads `localStorage` on mount to hydrate `token` state
    - `login()` calls `authService.login()`, calls `setToken()`, updates state, navigates to `/dashboard`
    - `logout()` calls `authService.logout()` in a try/finally; `finally` calls `clearToken()`, resets state, navigates to `/login`
    - Export `AuthContext` and `AuthProvider`
    - _Requirements: 6.3, 6.4, 7.2, 7.3, 9.1, 9.2_

  - [x] 7.2 Implement `src/shared/hooks/useAuth.ts`
    - Export `useAuth()` hook that consumes `AuthContext` and throws if used outside provider
    - _Requirements: 6.4, 7.3_

  - [ ]* 7.3 Write property test for session restoration (Property 12)
    - **Property 12: Session restoration from localStorage**
    - **Validates: Requirements 9.1, 9.2**
    - Use `fc.string({ minLength: 1 })` to generate token strings; pre-populate `localStorage` before mounting `AuthProvider`; assert `useAuth().token` equals the pre-populated token

  - [ ]* 7.4 Write property test for logout token clearing (Property 9)
    - **Property 9: Logout clears the token regardless of response outcome**
    - **Validates: Requirements 7.2**
    - Use `msw` to simulate both success and error responses for logout; assert `localStorage.getItem('hmdm_auth_token')` is `null` after logout in both cases

- [x] 8. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 9. Implement the Login feature
  - [x] 9.1 Implement `src/features/auth/useLogin.ts`
    - Manage `isLoading` and `error` state
    - On submit: set `isLoading = true`, call `AuthContext.login()`, catch errors and set `error` string, always reset `isLoading = false` in finally
    - Return `{ register, handleSubmit, isLoading, error }` (react-hook-form integration)
    - _Requirements: 6.5, 6.6_

  - [x] 9.2 Implement `src/features/auth/LoginPage.tsx`
    - Use shadcn/ui `Card`, `CardHeader`, `CardTitle`, `CardDescription`, `CardContent`
    - Use shadcn/ui `Form`, `FormField`, `FormLabel`, `FormControl`, `FormMessage`
    - Use shadcn/ui `Input` for username and password fields
    - Use shadcn/ui `Button` with `disabled={isLoading}` and `Loader2` spinner when loading
    - Render error message in `<FormMessage>` without clearing input fields on error
    - _Requirements: 6.1, 6.5, 6.6, 6.7_

  - [ ]* 9.3 Write property test for login token round-trip (Property 4 — component level)
    - **Property 4: Login stores the token (round trip)**
    - **Validates: Requirements 6.3**
    - Use `fc.string({ minLength: 1 })` for `authToken`; mock login endpoint with `msw`; render full app; assert `localStorage.getItem('hmdm_auth_token')` equals the generated token after form submission

  - [ ]* 9.4 Write property test for login navigation (Property 5)
    - **Property 5: Successful login navigates to dashboard**
    - **Validates: Requirements 6.4**
    - Use `fc.record({ login: fc.string({ minLength: 1 }), password: fc.string({ minLength: 1 }) })` for credentials; mock successful login; assert router location is `/dashboard` after submission

  - [ ]* 9.5 Write property test for login error preserving inputs (Property 6)
    - **Property 6: Login error preserves input field values**
    - **Validates: Requirements 6.5**
    - Use `fc.string({ minLength: 1 })` for username and password; mock error response; assert input values are unchanged and error message element is present in DOM

  - [ ]* 9.6 Write property test for submit button disabled during loading (Property 7)
    - **Property 7: Submit button is disabled while login is in flight**
    - **Validates: Requirements 6.6**
    - Intercept the login request before it resolves; assert submit button has `disabled` attribute set

- [x] 10. Implement AuthGuard and routing
  - [x] 10.1 Implement `src/shared/components/AuthGuard.tsx`
    - Consume `useAuth()`; if `token` is null, render `<Navigate to="/login" replace />`; otherwise render `children`
    - _Requirements: 6.8, 8.2, 9.3_

  - [x] 10.2 Implement `src/app/App.tsx`
    - Set up `<BrowserRouter>` with `<Routes>`:
      - `<Route path="/login" element={<LoginPage />} />`
      - `<Route path="/dashboard" element={<AuthGuard><DashboardPage /></AuthGuard>} />`
      - `<Route path="/" element={<Navigate to="/dashboard" />} />`
    - _Requirements: 6.7, 8.1_

  - [x] 10.3 Update `src/main.tsx`
    - Wrap `<App />` in `<AuthProvider>` from `providers.tsx`
    - _Requirements: 4.11_

  - [ ]* 10.4 Write property test for AuthGuard redirect (Property 8)
    - **Property 8: AuthGuard redirects unauthenticated users**
    - **Validates: Requirements 6.8, 8.2, 9.3**
    - Use `fc.constantFrom('/dashboard')` (and future protected paths); ensure `localStorage` is empty; assert rendered output is a redirect to `/login` and protected content is not present

- [x] 11. Implement the Dashboard placeholder
  - [x] 11.1 Implement `src/features/dashboard/DashboardPage.tsx`
    - Render a layout with a `<header>` containing the page title, `"Welcome, {username}"` text, and a shadcn/ui `Button` that calls `logout()` from `useAuth()`
    - Render a `<main>` section with placeholder content
    - _Requirements: 8.1, 8.3, 8.4_

  - [ ]* 11.2 Write property test for dashboard username display (Property 11)
    - **Property 11: Dashboard displays the authenticated username**
    - **Validates: Requirements 8.3**
    - Use `fc.string({ minLength: 1 })` for username; render `DashboardPage` with auth context providing that username; assert a text node containing the username is present in the DOM

  - [ ]* 11.3 Write property test for logout navigation (Property 10)
    - **Property 10: Logout navigates to login page**
    - **Validates: Requirements 7.3**
    - Render authenticated app; trigger logout button click; assert router location is `/login`

- [x] 12. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 13. Verify workspace structure and configuration files
  - [x] 13.1 Write unit tests for workspace structure
    - Assert `backend/pom.xml` exists (Requirement 1.1, 2.2)
    - Assert `frontend/package.json` exists (Requirement 1.2)
    - Assert `frontend/components.json` exists (Requirement 3.7)
    - Assert `frontend/.eslintrc.cjs` (or `eslint.config.js`) exists (Requirement 3.3)
    - Assert `frontend/.prettierrc` exists (Requirement 3.4)
    - Assert `frontend/src/main.tsx` exists (Requirement 4.11)
    - _Requirements: 1.1, 1.2, 2.2, 3.3, 3.4, 3.7, 4.11_

  - [ ]* 13.2 Write property test for no legacy code in frontend (Property 1)
    - **Property 1: No legacy code in frontend source**
    - **Validates: Requirements 3.2, 3.6**
    - Enumerate all files under `frontend/src/`; assert each has `.ts` or `.tsx` extension; assert file content contains none of `<%`, `$(`, `jQuery(`, `ng-`

- [x] 14. Final checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Each task references specific requirements for traceability
- Property tests use `fast-check` with `numRuns: 100` minimum; each test must include a comment with the tag `Feature: hmdm-modern-architecture, Property {N}: {property_text}`
- shadcn/ui components are generated into `src/shared/ui/` via the CLI; do not hand-write them
- The backend source in `hmdm-server/` must not be modified — only moved to `backend/`
- The `msw` service worker setup (`public/mockServiceWorker.js`) is required for browser-environment tests
