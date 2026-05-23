# Design Document: Navigation Layout

## Overview

This document describes the technical design for the Navigation/Layout feature — a persistent application shell that wraps all protected routes in the Headwind MDM frontend. The shell consists of three regions: a fixed sidebar with navigation links, a top header bar, and a main content area. It replaces the ad-hoc header currently embedded in `DashboardPage` and becomes the foundation for all future protected pages.

The implementation uses the existing stack: React 18, TypeScript, React Router v6, shadcn/ui, and Tailwind CSS. No new dependencies are introduced beyond the shadcn/ui Sheet component (already part of the shadcn/ui package, just not yet scaffolded into `src/shared/ui/`).

### Key Design Decisions

- **AppLayout as a layout route wrapper**: `App.tsx` is refactored to use a React Router v6 layout route pattern — a parent `<Route>` renders `AppLayout` with an `<Outlet />`, and child routes render their page components into that outlet. This means adding a new protected route requires only a single `<Route>` addition, with no changes to `AppLayout`.
- **shadcn/ui Sheet for mobile drawer**: The Sheet component provides accessible focus trapping, outside-click dismissal, and keyboard handling out of the box, satisfying the accessibility requirements without custom logic.
- **React Router NavLink for active state**: `NavLink` from `react-router-dom` provides the `isActive` callback and sets `aria-current="page"` automatically, eliminating manual `window.location` comparisons.
- **Lucide React icons**: Already available in the project via shadcn/ui's dependency chain; `LayoutDashboard`, `Smartphone`, `Users`, `Settings`, and `Menu` icons are used.

---

## Architecture

The layout integrates into the existing routing tree as a layout route. The component hierarchy for a protected page is:

```
BrowserRouter (main.tsx)
└── AuthProvider (providers.tsx)
    └── App (App.tsx)
        └── Routes
            ├── /login → LoginPage
            └── / (AuthGuard + AppLayout — layout route)
                ├── /dashboard → DashboardPage
                ├── /devices  → DevicesPage (future)
                ├── /users    → UsersPage (future)
                └── /settings → SettingsPage (future)
```

`AuthGuard` wraps the layout route element. `AppLayout` renders `<Outlet />` as its children slot, so React Router injects the matched child page component automatically.

### Desktop Layout (≥ 768px)

```
┌──────────────────────────────────────────────┐
│  Sidebar (240px fixed)  │  Header             │
│  ─────────────────────  │  ─────────────────  │
│  Headwind MDM           │  [page title] [user]│
│                         │  [logout]           │
│  ● Dashboard            ├─────────────────────┤
│    Devices              │                     │
│    Users                │   <Outlet />        │
│    Settings             │   (page content)    │
│                         │                     │
└──────────────────────────────────────────────┘
```

### Mobile Layout (< 768px)

```
┌──────────────────────────────┐
│  [☰] Header                  │
│  ──────────────────────────  │
│  [menu] [page title] [user]  │
│  [logout]                    │
├──────────────────────────────┤
│                              │
│   <Outlet />                 │
│   (page content)             │
│                              │
└──────────────────────────────┘

  [Sheet drawer slides in from left on menu click]
```

---

## Components and Interfaces

### NavItem type

```typescript
// src/features/layout/types.ts
export interface NavItem {
  label: string
  path: string
  icon: React.ComponentType<{ className?: string }>
}
```

### Navigation configuration

```typescript
// src/features/layout/navItems.ts
import { LayoutDashboard, Smartphone, Users, Settings } from 'lucide-react'
import type { NavItem } from './types'

export const NAV_ITEMS: NavItem[] = [
  { label: 'Dashboard', path: '/dashboard', icon: LayoutDashboard },
  { label: 'Devices',   path: '/devices',   icon: Smartphone },
  { label: 'Users',     path: '/users',     icon: Users },
  { label: 'Settings',  path: '/settings',  icon: Settings },
]
```

### AppLayout (`src/features/layout/AppLayout.tsx`)

```typescript
interface AppLayoutProps {
  // No props required — uses React Router <Outlet /> internally
}
```

Renders the full-viewport shell: `Sidebar` on the left (hidden on mobile), `Header` + `<Outlet />` on the right. Owns the `mobileOpen: boolean` state that controls the Sheet drawer, passing `open`/`onOpenChange` down to `Sidebar` and a `onMenuClick` handler down to `Header`.

### Sidebar (`src/features/layout/Sidebar.tsx`)

```typescript
interface SidebarProps {
  mobileOpen: boolean
  onMobileClose: () => void
}
```

Renders two representations:
1. **Desktop**: a `<nav aria-label="Main navigation">` inside a fixed-width `div` (visible only at `md:` and above via Tailwind).
2. **Mobile**: a shadcn/ui `<Sheet>` containing the same `<nav>`, controlled by `mobileOpen`/`onMobileClose`. Clicking any `NavLink` inside the Sheet calls `onMobileClose`.

The nav content is extracted into a shared `<SidebarNav onNavigate?>` sub-component to avoid duplication.

### Header (`src/features/layout/Header.tsx`)

```typescript
interface HeaderProps {
  onMenuClick: () => void
}
```

Renders a `<header>` with a bottom border. Left region: menu toggle button (visible only below `md:`) + application name. Right region: username from `useAuth()` + logout `<Button>`.

### Sheet (`src/shared/ui/sheet.tsx`)

Scaffolded from shadcn/ui. Provides `Sheet`, `SheetContent`, `SheetHeader`, `SheetTitle` primitives used by `Sidebar` for the mobile drawer.

---

## Data Models

No new server-side data models are introduced. The layout consumes existing client-side state:

| Source | Data | Consumer |
|---|---|---|
| `AuthContext` | `username: string \| null` | `Header` (display) |
| `AuthContext` | `logout: () => Promise<void>` | `Header` (button handler) |
| React Router | current pathname | `NavLink` (active state, automatic) |
| Local state in `AppLayout` | `mobileOpen: boolean` | `Sidebar` (Sheet open/close), `Header` (menu button) |

The `NAV_ITEMS` array is a static constant — no API calls or dynamic data are needed for the navigation structure.

---

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Children pass-through

*For any* React node passed as children (via `<Outlet />`) to `AppLayout`, that content should appear inside the main content region and not inside the sidebar or header regions.

**Validates: Requirements 1.2**

---

### Property 2: Nav link navigation

*For any* nav item in `NAV_ITEMS`, clicking its rendered link in the Sidebar should cause the router's current pathname to equal that nav item's `path`.

**Validates: Requirements 2.2**

---

### Property 3: Active route state and accessibility

*For any* nav item in `NAV_ITEMS` and any current route pathname, the nav item's active visual style should be applied if and only if its `path` equals the current pathname, and `aria-current="page"` should be present on the active item and absent on all others.

**Validates: Requirements 2.3, 5.1, 5.2, 5.4, 7.2**

---

### Property 4: Username display

*For any* non-null username value provided by `AuthContext`, the `Header` component should render that exact username string in its right region.

**Validates: Requirements 3.2**

---

### Property 5: Logout invocation

*For any* logout function provided by `AuthContext`, clicking the logout button in `Header` should invoke that function exactly once.

**Validates: Requirements 3.3**

---

### Property 6: Menu button opens drawer

*For any* initial state where the mobile drawer is closed, clicking the menu toggle button in `Header` should transition the drawer to the open state.

**Validates: Requirements 4.3**

---

### Property 7: Nav click closes drawer

*For any* nav item, clicking its link while the mobile drawer is open should transition the drawer to the closed state.

**Validates: Requirements 4.4**

---

### Property 8: Nav links render as anchor elements

*For any* nav item in `NAV_ITEMS`, the rendered DOM element for that nav link should be an `<a>` element (or resolve to one), ensuring keyboard focus and screen reader compatibility.

**Validates: Requirements 7.1**

---

## Error Handling

| Scenario | Handling |
|---|---|
| `useAuth()` returns `null` username | `Header` renders a fallback (e.g., empty string or "—"); no crash |
| `logout()` throws | `AuthContext.logout` already has a `try/finally` that clears state regardless; `Header` does not need additional error handling |
| Navigation to an unmatched route | React Router renders nothing in `<Outlet />`; the shell (sidebar + header) remains visible |
| Sheet component not yet installed | Build-time TypeScript error; resolved by scaffolding `src/shared/ui/sheet.tsx` before implementing `Sidebar` |

---

## Testing Strategy

### Dual Testing Approach

Both unit tests and property-based tests are required. Unit tests cover specific examples, integration points, and edge cases. Property-based tests verify universal behaviors across generated inputs.

### Unit Tests

Focus areas:
- Render `AppLayout` with mock children and assert the three structural regions are present (Req 1.1)
- Render `Sidebar` and assert the four nav items appear in order with correct labels and paths (Req 2.1)
- Render `Sidebar` and assert the branding text "Headwind MDM" is present (Req 2.4)
- Render `Header` with a mock `AuthContext` and assert the username and logout button are present (Req 3.1, 3.2)
- Render `Sidebar` and assert the `<nav aria-label="Main navigation">` element exists (Req 7.3)
- Render `Header` and assert the logout button has accessible text (Req 7.4)
- Render the full router with a protected route and assert `AppLayout` wraps the page (Req 6.1)
- Render `DashboardPage` and assert it no longer contains a `<header>` element (Req 6.3)

### Property-Based Tests

Library: **fast-check** (install as dev dependency: `npm install --save-dev fast-check`)

Configuration: minimum **100 iterations** per property test.

Each test is tagged with a comment in the format:
`// Feature: navigation-layout, Property N: <property text>`

| Property | Test Description |
|---|---|
| P1: Children pass-through | Generate arbitrary React string content; render `AppLayout` with it; assert content appears in `<main>` and not in `<aside>` or `<header>` |
| P2: Nav link navigation | For each item in `NAV_ITEMS`, render with a `MemoryRouter`, click the link, assert `location.pathname === item.path` |
| P3: Active route state | For each item in `NAV_ITEMS` × each item's path as current route, assert active class and `aria-current="page"` on the matching item only |
| P4: Username display | Generate arbitrary non-empty username strings; render `Header` with mocked `AuthContext`; assert the username appears in the rendered output |
| P5: Logout invocation | Generate a mock logout function; render `Header`; click logout; assert mock was called exactly once |
| P6: Menu button opens drawer | Render `AppLayout` in mobile viewport; assert drawer is closed; click menu button; assert drawer is open |
| P7: Nav click closes drawer | Render `Sidebar` with `mobileOpen=true`; click any nav link; assert `onMobileClose` was called |
| P8: Nav links render as anchors | For each item in `NAV_ITEMS`, render `Sidebar`; query the link element; assert `tagName === 'A'` |

Note: P2 and P3 iterate over the fixed `NAV_ITEMS` array rather than generating arbitrary nav items, since the nav structure is a static configuration. P4 and P5 use fast-check to generate arbitrary string/function inputs to verify the Header's behavior holds for all valid auth context values.
