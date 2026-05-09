# Requirements Document

## Introduction

This feature introduces a persistent Navigation/Layout system as the foundational shell for all protected pages in the Headwind MDM frontend. The layout replaces the ad-hoc header currently embedded in `DashboardPage` and provides a consistent, reusable structure — a collapsible sidebar with navigation links, a top header bar, and a main content area — that wraps every route behind the `AuthGuard`. All UI is built with shadcn/ui components on top of Tailwind CSS. The layout is the second deliverable in the HMDM Modern Architecture roadmap, enabling Dashboard improvements, Devices, and Users features to be built on top of it.

## Glossary

- **App_Layout**: The top-level shell component that renders the Sidebar, Header, and main content area, and wraps all protected routes.
- **Sidebar**: The persistent left-side navigation panel containing navigation links to all primary sections of the application.
- **Header**: The top bar rendered inside App_Layout, displaying the application title, the current user's login name, and a logout button.
- **Nav_Link**: A single navigation item inside the Sidebar that links to a named route and reflects the active state when its route is current.
- **Auth_Guard**: The existing route protection component that redirects unauthenticated users to the Login_Page (unchanged from the hmdm-modern-architecture spec).
- **Router**: The React Router v6 instance managing client-side navigation.
- **Auth_Context**: The existing React context providing the authenticated username and logout function.
- **Mobile_Breakpoint**: The viewport width threshold below which the Sidebar collapses into a hidden drawer. Defined as `768px` (Tailwind `md` breakpoint).

---

## Requirements

### Requirement 1: App Layout Shell

**User Story:** As a developer, I want a single layout component that wraps all protected pages, so that every page automatically gets the sidebar and header without duplicating markup.

#### Acceptance Criteria

1. THE App_Layout SHALL render a Sidebar, a Header, and a main content area as its three structural regions.
2. THE App_Layout SHALL accept a `children` prop and render it inside the main content area.
3. WHEN a protected route is accessed, THE Auth_Guard SHALL render the App_Layout as the outer wrapper, with the page component passed as children.
4. THE App_Layout SHALL use shadcn/ui primitives and Tailwind CSS for all styling; no custom CSS files SHALL be introduced.
5. THE App_Layout SHALL occupy the full viewport height and width.

---

### Requirement 2: Sidebar Navigation

**User Story:** As an MDM administrator, I want a persistent sidebar with navigation links, so that I can move between sections of the application without losing context.

#### Acceptance Criteria

1. THE Sidebar SHALL display navigation links for the following four sections, in order: Dashboard (`/dashboard`), Devices (`/devices`), Users (`/users`), Settings (`/settings`).
2. WHEN the user clicks a Nav_Link, THE Router SHALL navigate to the corresponding route.
3. WHEN the current route matches a Nav_Link's target path, THE Nav_Link SHALL render in an active visual state distinct from inactive links (e.g., highlighted background, bold text, or accent color).
4. THE Sidebar SHALL display the application name "Headwind MDM" as a branding element at the top of the panel.
5. THE Sidebar SHALL use shadcn/ui components where applicable and Tailwind CSS for layout and spacing.
6. WHEN the user navigates between routes, THE Sidebar SHALL remain visible and its position SHALL not change.

---

### Requirement 3: Top Header Bar

**User Story:** As an MDM administrator, I want a top header bar showing the app title and my account controls, so that I always know which application I am using and can log out easily.

#### Acceptance Criteria

1. THE Header SHALL display the current page title or application name in its left region.
2. THE Header SHALL display the authenticated user's login name, retrieved from Auth_Context, in its right region.
3. THE Header SHALL render a logout button in its right region that, when clicked, invokes the logout function from Auth_Context.
4. WHEN the logout button is clicked, THE Auth_Context SHALL execute the logout flow defined in the hmdm-modern-architecture spec (clear token, redirect to `/login`).
5. THE Header SHALL use shadcn/ui components and Tailwind CSS for all styling.
6. THE Header SHALL be visually separated from the main content area (e.g., via a bottom border).

---

### Requirement 4: Responsive Sidebar — Mobile Collapse

**User Story:** As an MDM administrator using a mobile or narrow-viewport device, I want the sidebar to collapse so that the main content area is not obscured.

#### Acceptance Criteria

1. WHILE the viewport width is below the Mobile_Breakpoint, THE Sidebar SHALL be hidden by default and SHALL NOT occupy horizontal space in the layout.
2. WHILE the viewport width is at or above the Mobile_Breakpoint, THE Sidebar SHALL be permanently visible and SHALL occupy a fixed horizontal width.
3. WHEN the viewport width is below the Mobile_Breakpoint, THE Header SHALL display a menu toggle button that, when clicked, opens the Sidebar as an overlay or drawer.
4. WHEN the Sidebar is open as a drawer and the user clicks a Nav_Link, THE Sidebar SHALL close automatically.
5. WHEN the Sidebar is open as a drawer and the user clicks outside the Sidebar area, THE Sidebar SHALL close.
6. THE responsive behavior SHALL be implemented using Tailwind CSS responsive prefixes (`md:`) and/or shadcn/ui Sheet component; no third-party responsive libraries SHALL be introduced.

---

### Requirement 5: Active Route Highlighting

**User Story:** As an MDM administrator, I want the current section highlighted in the sidebar, so that I always know where I am in the application.

#### Acceptance Criteria

1. WHEN the Router's current pathname matches a Nav_Link's target path exactly, THE Nav_Link SHALL apply an active style class distinct from the default style.
2. WHEN the Router's current pathname does not match a Nav_Link's target path, THE Nav_Link SHALL render in its default (inactive) style.
3. THE active state SHALL be determined using React Router v6's `useLocation` hook or `NavLink` component; no manual URL string comparisons against `window.location` SHALL be used.
4. WHEN the user navigates to a new route, THE previously active Nav_Link SHALL revert to its inactive style and the newly matching Nav_Link SHALL become active, without a full page reload.

---

### Requirement 6: Layout Integration with Protected Routes

**User Story:** As a developer, I want the layout to automatically apply to all protected routes, so that adding a new page only requires defining a route — not re-implementing the shell.

#### Acceptance Criteria

1. THE App_Layout SHALL be applied to all routes wrapped by Auth_Guard, including `/dashboard`, and any future routes such as `/devices`, `/users`, and `/settings`.
2. WHEN a new protected route is added to the Router, THE App_Layout SHALL render for that route without any modification to the App_Layout component itself.
3. THE existing `DashboardPage` component SHALL be refactored to remove its inline header markup, delegating header and navigation rendering entirely to App_Layout.
4. THE App_Layout SHALL be located at `src/features/layout/AppLayout.tsx` within the feature-based directory structure.
5. THE Sidebar component SHALL be located at `src/features/layout/Sidebar.tsx`.
6. THE Header component SHALL be located at `src/features/layout/Header.tsx`.

---

### Requirement 7: Accessibility

**User Story:** As an MDM administrator using assistive technology, I want the navigation to be keyboard-navigable and screen-reader-friendly, so that I can use the application without a mouse.

#### Acceptance Criteria

1. THE Sidebar SHALL render navigation links as native `<a>` elements or components that resolve to `<a>` elements, so that keyboard focus and screen reader announcement work correctly.
2. THE Nav_Link that is currently active SHALL have an `aria-current="page"` attribute set.
3. THE Sidebar SHALL be wrapped in a `<nav>` element with an `aria-label` of `"Main navigation"`.
4. THE Header's logout button SHALL have a descriptive accessible label (visible text or `aria-label`) that identifies its action.
5. WHEN the mobile drawer is open, THE Sidebar SHALL be reachable via keyboard Tab navigation before the main content area.
6. WHEN the mobile drawer is closed, THE Sidebar links SHALL NOT be reachable via keyboard Tab navigation.
