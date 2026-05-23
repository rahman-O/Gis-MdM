# Implementation Plan: Navigation Layout

## Overview

Implement the persistent application shell using React Router v6 layout routes, shadcn/ui Sheet for mobile drawer, and NavLink for active state. Tasks build incrementally from types → static config → UI components → routing refactor → placeholder pages → tests.

## Tasks

- [x] 1. Scaffold Sheet component and define shared types
  - [x] 1.1 Scaffold shadcn/ui Sheet into `src/shared/ui/sheet.tsx`
    - Run `npx shadcn@latest add sheet` or manually scaffold the Sheet primitives (`Sheet`, `SheetContent`, `SheetHeader`, `SheetTitle`, `SheetDescription`) using the shadcn/ui pattern already established in `src/shared/ui/`
    - _Requirements: 4.3, 4.5, 4.6_
  - [x] 1.2 Create `src/features/layout/types.ts` with `NavItem` interface
    - Define `NavItem { label: string; path: string; icon: React.ComponentType<{ className?: string }> }`
    - _Requirements: 2.1_
  - [x] 1.3 Create `src/features/layout/navItems.ts` with `NAV_ITEMS` constant
    - Import `LayoutDashboard`, `Smartphone`, `Users`, `Settings` from `lucide-react`
    - Export `NAV_ITEMS: NavItem[]` with four entries: Dashboard `/dashboard`, Devices `/devices`, Users `/users`, Settings `/settings`
    - _Requirements: 2.1_

- [x] 2. Implement Header component
  - [x] 2.1 Create `src/features/layout/Header.tsx`
    - Accept `onMenuClick: () => void` prop
    - Render a `<header>` with bottom border using Tailwind
    - Left region: menu toggle `<Button>` with `Menu` icon (visible only below `md:`) that calls `onMenuClick`; application name text
    - Right region: username from `useAuth()` + logout `<Button>` that calls `logout` from `useAuth()`
    - Logout button must have descriptive accessible label
    - _Requirements: 3.1, 3.2, 3.3, 3.5, 3.6, 4.3, 7.4_
  - [ ]* 2.2 Write property test for username display (Property 4)
    - **Property 4: Username display**
    - **Validates: Requirements 3.2**
    - Use `fc.string({ minLength: 1 })` to generate arbitrary non-empty usernames; render `Header` with mocked `AuthContext`; assert the username string appears in the rendered output
    - Minimum 100 iterations; tag comment: `// Feature: navigation-layout, Property 4: Username display`
  - [ ]* 2.3 Write property test for logout invocation (Property 5)
    - **Property 5: Logout invocation**
    - **Validates: Requirements 3.3**
    - Generate a `vi.fn()` logout mock; render `Header`; click the logout button; assert mock was called exactly once
    - Minimum 100 iterations; tag comment: `// Feature: navigation-layout, Property 5: Logout invocation`
  - [ ]* 2.4 Write unit tests for Header
    - Assert username and logout button are present (Req 3.1, 3.2)
    - Assert logout button has accessible text (Req 7.4)
    - Assert menu toggle button is present (Req 4.3)

- [x] 3. Implement Sidebar component
  - [x] 3.1 Create `src/features/layout/Sidebar.tsx`
    - Accept `mobileOpen: boolean` and `onMobileClose: () => void` props
    - Extract shared `SidebarNav` sub-component that maps `NAV_ITEMS` to `NavLink` elements with icon + label
    - `NavLink` active callback applies highlighted style and sets `aria-current="page"` on the active item
    - Desktop: `<nav aria-label="Main navigation">` inside a fixed-width `div`, hidden below `md:` via Tailwind
    - Mobile: shadcn/ui `<Sheet>` controlled by `mobileOpen`/`onMobileClose`; each `NavLink` click calls `onMobileClose`
    - Branding text "Headwind MDM" at top of panel
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 4.1, 4.2, 4.4, 4.5, 4.6, 5.1, 5.2, 5.3, 5.4, 7.1, 7.2, 7.3_
  - [ ]* 3.2 Write property test for active route state (Property 3)
    - **Property 3: Active route state and accessibility**
    - **Validates: Requirements 2.3, 5.1, 5.2, 5.4, 7.2**
    - For each item in `NAV_ITEMS` × each item's path as current route via `MemoryRouter initialEntries`; assert active class and `aria-current="page"` on the matching item only, and absent on all others
    - Minimum 100 iterations; tag comment: `// Feature: navigation-layout, Property 3: Active route state and accessibility`
  - [ ]* 3.3 Write property test for nav click closes drawer (Property 7)
    - **Property 7: Nav click closes drawer**
    - **Validates: Requirements 4.4**
    - Render `Sidebar` with `mobileOpen=true` and a `vi.fn()` `onMobileClose`; click each nav link; assert `onMobileClose` was called
    - Minimum 100 iterations; tag comment: `// Feature: navigation-layout, Property 7: Nav click closes drawer`
  - [ ]* 3.4 Write property test for nav links render as anchors (Property 8)
    - **Property 8: Nav links render as anchor elements**
    - **Validates: Requirements 7.1**
    - For each item in `NAV_ITEMS`, render `Sidebar` inside `MemoryRouter`; query the link element; assert `tagName === 'A'`
    - Minimum 100 iterations; tag comment: `// Feature: navigation-layout, Property 8: Nav links render as anchor elements`
  - [ ]* 3.5 Write unit tests for Sidebar
    - Assert four nav items appear in order with correct labels and paths (Req 2.1)
    - Assert branding text "Headwind MDM" is present (Req 2.4)
    - Assert `<nav aria-label="Main navigation">` element exists (Req 7.3)

- [x] 4. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 5. Implement AppLayout and wire routing
  - [x] 5.1 Create `src/features/layout/AppLayout.tsx`
    - Own `mobileOpen: boolean` state (default `false`)
    - Render full-viewport shell: `Sidebar` (left, desktop only) + right column containing `Header` + `<main>` with `<Outlet />`
    - Pass `mobileOpen`/`onMobileClose` to `Sidebar`; pass `onMenuClick` to `Header`
    - Use Tailwind for full-height flex layout; no custom CSS files
    - _Requirements: 1.1, 1.2, 1.4, 1.5, 6.4_
  - [ ]* 5.2 Write property test for children pass-through (Property 1)
    - **Property 1: Children pass-through**
    - **Validates: Requirements 1.2**
    - Use `fc.string({ minLength: 1 })` to generate arbitrary string content; render `AppLayout` with a `MemoryRouter` and the string as a child route; assert content appears inside `<main>` and not inside `<aside>` or `<header>`
    - Minimum 100 iterations; tag comment: `// Feature: navigation-layout, Property 1: Children pass-through`
  - [ ]* 5.3 Write property test for menu button opens drawer (Property 6)
    - **Property 6: Menu button opens drawer**
    - **Validates: Requirements 4.3**
    - Render `AppLayout` in a mobile-width container; assert Sheet is closed; click menu toggle button; assert Sheet transitions to open state
    - Minimum 100 iterations; tag comment: `// Feature: navigation-layout, Property 6: Menu button opens drawer`
  - [ ]* 5.4 Write unit tests for AppLayout
    - Assert three structural regions (sidebar, header, main) are present (Req 1.1)
    - Render with a mock child route and assert it appears in `<main>` (Req 1.2)
    - Render full router and assert `AppLayout` wraps the page (Req 6.1)

- [x] 6. Refactor App.tsx to layout route pattern
  - [x] 6.1 Refactor `src/app/App.tsx` to use React Router v6 layout route
    - Replace the per-route `<AuthGuard>` wrappers with a single parent `<Route>` whose element is `<AuthGuard><AppLayout /></AuthGuard>`
    - Add child `<Route>` entries for `/dashboard`, `/devices`, `/users`, `/settings` as children of the layout route
    - Add a root redirect `<Route path="/" element={<Navigate to="/dashboard" replace />} />`
    - _Requirements: 1.3, 6.1, 6.2_

- [x] 7. Refactor DashboardPage and add placeholder pages
  - [x] 7.1 Refactor `src/features/dashboard/DashboardPage.tsx`
    - Remove the inline `<header>` element and all auth-related imports (`useAuth`, `Button`)
    - Retain only the main content area markup
    - _Requirements: 6.3_
  - [x] 7.2 Create placeholder pages: `src/features/devices/DevicesPage.tsx`, `src/features/users/UsersPage.tsx`, `src/features/settings/SettingsPage.tsx`
    - Each page renders a minimal `<div>` with a heading identifying the section
    - _Requirements: 6.1, 6.2_
  - [ ]* 7.3 Write property test for nav link navigation (Property 2)
    - **Property 2: Nav link navigation**
    - **Validates: Requirements 2.2**
    - For each item in `NAV_ITEMS`, render the full router with `MemoryRouter`; click the nav link; assert `location.pathname === item.path`
    - Minimum 100 iterations; tag comment: `// Feature: navigation-layout, Property 2: Nav link navigation`
  - [ ]* 7.4 Write unit test for DashboardPage refactor
    - Assert `DashboardPage` no longer contains a `<header>` element (Req 6.3)

- [x] 8. Final checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Each task references specific requirements for traceability
- Property tests use `fast-check` with minimum 100 iterations per property
- Unit tests cover specific examples, edge cases, and integration points
- `fast-check` must be installed as a dev dependency: `npm install --save-dev fast-check`
