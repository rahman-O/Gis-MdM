# Design Document: HMDM Modern Architecture

## Overview

This design describes the architectural separation of the Headwind MDM project into a standalone Java/Maven backend and a new React + Vite + TypeScript frontend. The first feature delivered in the new frontend is user authentication (login/logout) with session persistence.

The backend is relocated to `/backend` and continues to serve its existing REST API unchanged. The frontend lives in `/frontend` and communicates with the backend exclusively through a centralized Axios-based API client. All UI is built with shadcn/ui components (Radix UI + Tailwind CSS).

### Goals

- Clean separation of concerns: backend is a pure REST API server; frontend is a pure SPA.
- Modern, maintainable frontend stack: React 18, Vite, TypeScript, React Router v6, shadcn/ui.
- Secure, token-based authentication with localStorage persistence.
- Feature-based directory structure that scales as new MDM features are added.

### Non-Goals

- Migrating any existing AngularJS/JSP UI logic to React (the legacy UI is replaced, not migrated).
- Modifying any backend Java source files.
- Implementing device management, user management, or any other MDM feature beyond authentication and the dashboard placeholder.

---

## Architecture

### Workspace Layout

```
/                          ← workspace root
├── backend/               ← relocated hmdm-server (Java/Maven, untouched)
│   ├── pom.xml
│   ├── common/
│   ├── server/
│   └── ...
├── frontend/              ← new React + Vite + TypeScript SPA
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   ├── components.json    ← shadcn/ui config
│   └── src/
│       ├── main.tsx
│       ├── app/
│       ├── features/
│       ├── shared/
│       └── services/
└── docs/                  ← project documentation
```

### Backend / Frontend Separation

```
┌─────────────────────────────────────────────────────────┐
│  Browser (SPA)                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │  React App (Vite dev server :5173 / static build) │  │
│  │  ┌────────────┐   ┌──────────────┐               │  │
│  │  │ Auth Pages │   │ Dashboard    │  (future pages)│  │
│  │  └─────┬──────┘   └──────┬───────┘               │  │
│  │        │                 │                        │  │
│  │  ┌─────▼─────────────────▼──────────────────┐    │  │
│  │  │         Axios API Client                  │    │  │
│  │  │  baseURL: http://localhost:8080/api        │    │  │
│  │  │  X-Auth-Token header (when present)        │    │  │
│  │  └──────────────────┬────────────────────────┘    │  │
│  └─────────────────────┼────────────────────────────-┘  │
└────────────────────────┼────────────────────────────────┘
                         │ HTTP/REST
┌────────────────────────▼────────────────────────────────┐
│  Java Backend (Tomcat :8080)                            │
│  POST /api/public/auth/login                            │
│  POST /api/public/auth/logout                           │
│  ... (all existing endpoints unchanged)                 │
└─────────────────────────────────────────────────────────┘
```

### Frontend Internal Architecture

The frontend follows a layered architecture:

```
┌──────────────────────────────────────────────────────┐
│  Routing Layer  (React Router v6)                    │
│  App.tsx → routes → Auth_Guard → Page components     │
├──────────────────────────────────────────────────────┤
│  Feature Layer  (src/features/)                      │
│  auth/  │  dashboard/  │  devices/ (future)          │
├──────────────────────────────────────────────────────┤
│  Shared Layer  (src/shared/)                         │
│  components/  │  hooks/  │  utils/  │  ui/           │
├──────────────────────────────────────────────────────┤
│  Service Layer  (src/services/)                      │
│  apiClient.ts  │  authService.ts                     │
├──────────────────────────────────────────────────────┤
│  State Layer  (React Context)                        │
│  AuthContext  (token, user, login, logout)           │
└──────────────────────────────────────────────────────┘
```

---

## Components and Interfaces

### Frontend Directory Structure

```
src/
├── main.tsx                        ← ReactDOM.createRoot, wraps App in providers
├── app/
│   ├── App.tsx                     ← Router setup, route definitions
│   └── providers.tsx               ← AuthProvider (and future providers)
├── features/
│   ├── auth/
│   │   ├── LoginPage.tsx           ← Login form page
│   │   ├── useLogin.ts             ← Login form logic hook
│   │   └── types.ts                ← LoginRequest, LoginResponse types
│   ├── dashboard/
│   │   └── DashboardPage.tsx       ← Protected dashboard placeholder
│   ├── devices/                    ← Reserved for future device management
│   └── users/                      ← Reserved for future user management
├── shared/
│   ├── components/
│   │   └── AuthGuard.tsx           ← Route protection component
│   ├── hooks/
│   │   └── useAuth.ts              ← Consumes AuthContext
│   ├── utils/
│   │   └── tokenStorage.ts         ← localStorage read/write helpers
│   └── ui/                         ← shadcn/ui generated components (Button, Input, etc.)
└── services/
    ├── apiClient.ts                ← Axios instance + interceptors
    └── authService.ts              ← login() and logout() API calls
```

### Routing Structure

```
App.tsx
└── <BrowserRouter>
    └── <AuthProvider>
        └── <Routes>
            ├── <Route path="/login"     element={<LoginPage />} />
            ├── <Route path="/dashboard" element={<AuthGuard><DashboardPage /></AuthGuard>} />
            └── <Route path="/"          element={<Navigate to="/dashboard" />} />
```

Any route wrapped in `<AuthGuard>` checks for a valid token. If absent, it redirects to `/login`.

### Component Tree — Login Page

```
LoginPage
└── <div> (full-screen centering wrapper, Tailwind)
    └── <Card>                        ← shadcn/ui Card
        ├── <CardHeader>
        │   ├── <CardTitle>           ← "Sign In"
        │   └── <CardDescription>    ← subtitle text
        └── <CardContent>
            └── <Form>               ← react-hook-form + shadcn/ui Form
                ├── <FormField name="login">
                │   ├── <FormLabel>  ← shadcn/ui Label
                │   └── <FormControl>
                │       └── <Input type="text" />   ← shadcn/ui Input
                ├── <FormField name="password">
                │   ├── <FormLabel>
                │   └── <FormControl>
                │       └── <Input type="password" />
                ├── <FormMessage />  ← error display (shadcn/ui)
                └── <Button type="submit" disabled={isLoading}>
                        {isLoading ? <Loader2 className="animate-spin" /> : "Sign In"}
                    </Button>
```

### Component Tree — Dashboard Page

```
DashboardPage
└── <div> (layout wrapper)
    ├── <header>
    │   ├── <h1> "Headwind MDM Dashboard"
    │   ├── <span> "Welcome, {username}"
    │   └── <Button onClick={logout}>  ← shadcn/ui Button
    │           "Sign Out"
    │       </Button>
    └── <main>
        └── placeholder content
```

### AuthGuard Component

```tsx
// src/shared/components/AuthGuard.tsx
function AuthGuard({ children }: { children: ReactNode }) {
  const { token } = useAuth();
  if (!token) return <Navigate to="/login" replace />;
  return <>{children}</>;
}
```

### API Client Interface

```typescript
// src/services/apiClient.ts
const apiClient = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: { 'Content-Type': 'application/json' },
});

// Request interceptor: attach token
apiClient.interceptors.request.use((config) => {
  const token = getToken(); // reads localStorage
  if (token) config.headers['X-Auth-Token'] = token;
  return config;
});

// Response interceptor: handle 401
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      clearToken();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### Auth Service Interface

```typescript
// src/services/authService.ts
interface LoginRequest  { login: string; password: string; }
interface LoginResponse { authToken: string; }

async function login(credentials: LoginRequest): Promise<LoginResponse>
async function logout(): Promise<void>
```

### Auth Context Interface

```typescript
// src/app/providers.tsx
interface AuthContextValue {
  token: string | null;
  username: string | null;
  login: (credentials: LoginRequest) => Promise<void>;
  logout: () => Promise<void>;
}
```

---

## Data Models

### Token Storage

The token is stored in `localStorage` under the key `hmdm_auth_token`. Helper functions in `src/shared/utils/tokenStorage.ts` encapsulate all reads and writes.

```typescript
const TOKEN_KEY = 'hmdm_auth_token';

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}
```

### Login API Contract

**Request**
```
POST /api/public/auth/login
Content-Type: application/json

{ "login": "admin", "password": "secret" }
```

**Success Response** (HTTP 200)
```json
{ "authToken": "<opaque-token-string>" }
```

**Error Response** (HTTP 4xx / 5xx)
```json
{ "message": "Invalid credentials" }
```

### Logout API Contract

**Request**
```
POST /api/public/auth/logout
X-Auth-Token: <token>
```

**Response**: HTTP 200 (body ignored by frontend)

### Auth Context State Shape

```typescript
interface AuthState {
  token: string | null;   // null = unauthenticated
  username: string | null; // derived from token or login response
}
```

On app initialization, `AuthProvider` reads `localStorage` to hydrate `token`. If a token exists, the user is considered authenticated without a network round-trip.

### Login Flow (Data Flow Diagram)

```
User                LoginPage          useLogin hook       authService        apiClient         Backend
 │                      │                    │                  │                  │               │
 │── submit form ──────►│                    │                  │                  │               │
 │                      │── call login() ───►│                  │                  │               │
 │                      │                    │── login(creds) ─►│                  │               │
 │                      │                    │                  │── POST /auth/login ─────────────►│
 │                      │                    │                  │                  │               │
 │                      │                    │                  │◄── { authToken } ────────────────│
 │                      │                    │◄── authToken ────│                  │               │
 │                      │                    │── setToken() ────────────────────────────────────── │
 │                      │                    │── setAuth(token) │                  │               │
 │                      │◄── success ────────│                  │                  │               │
 │                      │── navigate(/dashboard)                │                  │               │
 │◄── Dashboard Page ───│                    │                  │                  │               │
```

### Logout Flow (Data Flow Diagram)

```
User             DashboardPage       AuthContext         authService        apiClient         Backend
 │                    │                   │                  │                  │               │
 │── click logout ───►│                   │                  │                  │               │
 │                    │── logout() ──────►│                  │                  │               │
 │                    │                   │── logout() ─────►│                  │               │
 │                    │                   │                  │── POST /auth/logout ────────────►│
 │                    │                   │                  │◄── 200 OK ───────────────────────│
 │                    │                   │── clearToken() ──│                  │               │
 │                    │                   │── setAuth(null)  │                  │               │
 │                    │◄── navigate(/login)│                  │                  │               │
 │◄── Login Page ─────│                   │                  │                  │               │
```

### Session Restoration Flow

```
App Init
   │
   ├── AuthProvider mounts
   │       │
   │       ├── reads localStorage['hmdm_auth_token']
   │       │
   │       ├── token found? ──YES──► set token in state → user navigates freely
   │       │
   │       └── token absent? ──────► token = null → AuthGuard redirects to /login
```

---

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: No legacy code in frontend source

*For any* file found under `frontend/src/`, the file extension must be `.ts` or `.tsx` (TypeScript only), and the file content must not contain JSP syntax (`<%`), jQuery patterns (`$(` or `jQuery(`), or AngularJS directives (`ng-`).

**Validates: Requirements 3.2, 3.6**

---

### Property 2: Request headers are set correctly by the API client interceptor

*For any* outgoing request made through the API client, the `Content-Type` header must equal `application/json`. Additionally, *for any* request made when a token is present in `localStorage['hmdm_auth_token']`, the `X-Auth-Token` header must equal that token value.

**Validates: Requirements 5.2, 5.3**

---

### Property 3: 401 response clears the token

*For any* API response with HTTP status 401, after the response interceptor runs, `localStorage.getItem('hmdm_auth_token')` must return `null`.

**Validates: Requirements 5.4**

---

### Property 4: Login stores the token (round trip)

*For any* `authToken` string returned in a successful login response, after the login flow completes, `localStorage.getItem('hmdm_auth_token')` must equal that same `authToken` string.

**Validates: Requirements 6.3**

---

### Property 5: Successful login navigates to dashboard

*For any* valid credentials that produce a successful login response, after the login flow completes, the router's current location must be `/dashboard`.

**Validates: Requirements 6.4**

---

### Property 6: Login error preserves input field values

*For any* error response returned by the login endpoint, after the error is handled, the username and password input fields must retain the values that were entered before submission, and an error message element must be present in the DOM.

**Validates: Requirements 6.5**

---

### Property 7: Submit button is disabled while login is in flight

*For any* login submission that is pending (request not yet resolved), the submit button must have the `disabled` attribute set to `true`.

**Validates: Requirements 6.6**

---

### Property 8: AuthGuard redirects unauthenticated users

*For any* protected route accessed when `localStorage['hmdm_auth_token']` is absent or null, the rendered output must be a redirect to `/login` rather than the protected page content.

**Validates: Requirements 6.8, 8.2, 9.3**

---

### Property 9: Logout clears the token regardless of response outcome

*For any* logout operation — whether the backend responds with success or an error — after the logout flow completes, `localStorage.getItem('hmdm_auth_token')` must return `null`.

**Validates: Requirements 7.2**

---

### Property 10: Logout navigates to login page

*For any* logout action triggered by the user, after the logout flow completes, the router's current location must be `/login`.

**Validates: Requirements 7.3**

---

### Property 11: Dashboard displays the authenticated username

*For any* authenticated session with a known username, the rendered `DashboardPage` must contain a text node that includes that username.

**Validates: Requirements 8.3**

---

### Property 12: Session restoration from localStorage

*For any* token string pre-populated in `localStorage['hmdm_auth_token']` before `AuthProvider` mounts, after the provider initializes, the auth context's `token` value must equal that pre-populated token, and `AuthGuard` must allow access to protected routes.

**Validates: Requirements 9.1, 9.2**

---

## Error Handling

### API Client Errors

| Scenario | Handling |
|---|---|
| HTTP 401 from any endpoint | Clear `hmdm_auth_token` from localStorage; redirect to `/login` via `window.location.href` |
| HTTP 4xx (non-401) | Reject the promise; caller handles the error |
| HTTP 5xx | Reject the promise; caller displays a generic error message |
| Network error (no response) | Reject the promise; caller displays a connectivity error |

### Login Page Errors

- On any rejected login promise, the `useLogin` hook sets an `error` string in local state.
- The `LoginPage` renders this error string inside a `<FormMessage>` (shadcn/ui) below the form fields.
- Input fields are **not** cleared on error so the user can correct their credentials.
- The loading state is reset to `false` on both success and failure.

### Logout Errors

- The logout API call is wrapped in a `try/finally` block.
- Regardless of whether the POST succeeds or throws, the `finally` block calls `clearToken()` and navigates to `/login`.
- This ensures the user is always logged out locally even if the backend is unreachable.

### AuthGuard

- If `token` is `null` (absent or cleared), `AuthGuard` renders `<Navigate to="/login" replace />`.
- The `replace` flag prevents the protected route from appearing in browser history, so the back button does not return the user to a page they cannot access.

### Token Expiry

- The frontend does not decode or validate the token's expiry client-side.
- If the token has expired, the backend will return a 401 on the next API call, which triggers the 401 interceptor (clear + redirect).

---

## Testing Strategy

### Dual Testing Approach

Both unit tests and property-based tests are required. They are complementary:

- **Unit tests** verify specific examples, integration points, and edge cases.
- **Property-based tests** verify universal properties across many generated inputs.

### Tooling

| Purpose | Library |
|---|---|
| Test runner | Vitest |
| Component testing | React Testing Library |
| Property-based testing | `fast-check` |
| HTTP mocking | `msw` (Mock Service Worker) |

### Unit Tests (specific examples and integration)

- `LoginPage` renders username input, password input, and submit button.
- `LoginPage` at `/login` route is accessible without authentication.
- `DashboardPage` at `/dashboard` route is accessible when authenticated.
- `authService.logout()` calls `POST /api/public/auth/logout`.
- `DashboardPage` renders a logout button.
- Workspace structure: `/backend/pom.xml` exists, `/frontend/package.json` exists.
- `apiClient` base URL is `http://localhost:8080/api`.
- ESLint and Prettier config files exist in `/frontend`.
- shadcn/ui `components.json` exists in `/frontend`.

### Property-Based Tests

Each property test must run a minimum of **100 iterations**. Each test must include a comment referencing the design property it validates.

**Tag format:** `Feature: hmdm-modern-architecture, Property {N}: {property_text}`

| Property | Test Description |
|---|---|
| Property 1 | Generate random file paths under `src/`; assert all have `.ts`/`.tsx` extension and no legacy patterns |
| Property 2 | Generate random request configs with/without token; assert correct headers are set by interceptor |
| Property 3 | Generate mock 401 responses; assert token is null after interceptor runs |
| Property 4 | Generate random `authToken` strings; assert localStorage contains the token after login |
| Property 5 | Generate valid credential pairs; assert router location is `/dashboard` after login |
| Property 6 | Generate random credentials + error responses; assert inputs retain values and error is shown |
| Property 7 | Simulate in-flight login; assert submit button is disabled |
| Property 8 | Generate random protected route paths with no token; assert redirect to `/login` |
| Property 9 | Generate random logout outcomes (success/error); assert token is null after logout |
| Property 10 | Generate logout triggers; assert router location is `/login` after logout |
| Property 11 | Generate random username strings; assert username appears in rendered dashboard |
| Property 12 | Generate random token strings pre-set in localStorage; assert auth context token matches and protected routes are accessible |

### Example Property Test (fast-check)

```typescript
// Feature: hmdm-modern-architecture, Property 4: Login stores the token (round trip)
it('stores any authToken returned by login in localStorage', async () => {
  await fc.assert(
    fc.asyncProperty(fc.string({ minLength: 1 }), async (authToken) => {
      localStorage.clear();
      server.use(
        rest.post('/api/public/auth/login', (_, res, ctx) =>
          res(ctx.json({ authToken }))
        )
      );
      await authService.login({ login: 'user', password: 'pass' });
      expect(localStorage.getItem('hmdm_auth_token')).toBe(authToken);
    }),
    { numRuns: 100 }
  );
});
```

### Coverage Targets

- Service layer (`apiClient.ts`, `authService.ts`): 100% branch coverage
- Auth feature (`LoginPage`, `useLogin`, `AuthGuard`): 90%+ line coverage
- Shared utilities (`tokenStorage.ts`): 100% line coverage
