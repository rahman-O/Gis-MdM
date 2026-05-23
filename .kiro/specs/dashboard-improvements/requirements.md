# Requirements Document

## Introduction

The Dashboard Improvements feature replaces the current placeholder dashboard with a functional, data-driven overview page for the Headwind MDM admin interface. It displays summary statistics fetched from the backend, a recent devices table, loading skeletons during data fetches, and auto-refreshes every 60 seconds. All UI is built with shadcn/ui components (Card, Skeleton, Badge) on top of the existing React + Vite + TypeScript + Tailwind CSS stack.

## Glossary

- **Dashboard**: The main overview page rendered at `/dashboard` after authentication.
- **Summary_API**: The backend endpoint `GET /rest/private/summary` that returns aggregate statistics.
- **Devices_API**: The backend endpoint used to retrieve device records for the recent devices table.
- **Stat_Card**: A shadcn/ui `Card` component displaying an icon, a numeric value, and a label for a single summary statistic.
- **Summary_Stats**: The five aggregate values returned by `Summary_API`: `deviceTotal`, `deviceOnline`, `deviceEnrolled`, `configurationCount`, `applicationCount`.
- **Recent_Devices_Table**: A table showing the 5 most recently updated devices, sorted descending by `lastUpdate`.
- **Skeleton**: A shadcn/ui `Skeleton` component used as a placeholder while data is loading.
- **Auto_Refresh**: The mechanism that re-fetches dashboard data every 60 seconds without user interaction.
- **Badge**: A shadcn/ui `Badge` component used to indicate device online/offline status in the Recent_Devices_Table.
- **Dashboard_Page**: The React component that composes all dashboard UI elements.

---

## Requirements

### Requirement 1: Summary Statistics Cards

**User Story:** As an MDM administrator, I want to see key fleet statistics at a glance on the dashboard, so that I can quickly assess the state of my device fleet without navigating to individual sections.

#### Acceptance Criteria

1. WHEN the Dashboard_Page renders, THE Dashboard_Page SHALL display exactly five Stat_Cards: Total Devices, Online Devices, Enrolled Devices, Configurations, and Applications.
2. WHEN Summary_API returns a successful response, THE Dashboard_Page SHALL populate each Stat_Card with the corresponding value from `Summary_Stats` (`deviceTotal`, `deviceOnline`, `deviceEnrolled`, `configurationCount`, `applicationCount`).
3. THE Dashboard_Page SHALL render each Stat_Card using the shadcn/ui `Card` component containing an icon, a numeric value, and a text label.
4. WHEN Summary_API returns a successful response, THE Dashboard_Page SHALL display each numeric value as a non-negative integer.

---

### Requirement 2: Recent Devices Table

**User Story:** As an MDM administrator, I want to see the most recently active devices on the dashboard, so that I can quickly identify recent activity without going to the full devices list.

#### Acceptance Criteria

1. WHEN the Dashboard_Page renders, THE Dashboard_Page SHALL display a Recent_Devices_Table showing up to 5 device records.
2. WHEN device data is available, THE Recent_Devices_Table SHALL display records sorted in descending order by `lastUpdate` timestamp.
3. WHEN device data is available, THE Recent_Devices_Table SHALL display at minimum the device name, `lastUpdate` value, and online status for each record.
4. WHEN a device record has an online status of `true`, THE Recent_Devices_Table SHALL render a `Badge` component indicating the device is online.
5. WHEN a device record has an online status of `false`, THE Recent_Devices_Table SHALL render a `Badge` component indicating the device is offline.

---

### Requirement 3: Loading Skeletons

**User Story:** As an MDM administrator, I want to see loading indicators while dashboard data is being fetched, so that I understand the page is loading and not broken.

#### Acceptance Criteria

1. WHILE a data fetch is in progress, THE Dashboard_Page SHALL render `Skeleton` components in place of each Stat_Card's numeric value.
2. WHILE a data fetch is in progress, THE Dashboard_Page SHALL render `Skeleton` components in place of the Recent_Devices_Table rows.
3. WHEN a data fetch completes successfully, THE Dashboard_Page SHALL replace all `Skeleton` components with the actual data.
4. WHEN a data fetch completes with an error, THE Dashboard_Page SHALL replace all `Skeleton` components with an error message.

---

### Requirement 4: Auto-Refresh

**User Story:** As an MDM administrator, I want the dashboard to refresh automatically, so that I always see up-to-date fleet statistics without manually reloading the page.

#### Acceptance Criteria

1. WHEN the Dashboard_Page mounts, THE Dashboard_Page SHALL initiate an automatic refresh interval of 60 seconds.
2. WHEN 60 seconds elapse since the last successful data fetch, THE Dashboard_Page SHALL re-fetch data from Summary_API without user interaction.
3. WHEN the Dashboard_Page unmounts, THE Dashboard_Page SHALL cancel the auto-refresh interval to prevent memory leaks.
4. WHILE an auto-refresh fetch is in progress, THE Dashboard_Page SHALL continue displaying the previously loaded data rather than showing Skeleton components.

---

### Requirement 5: Data Fetching and API Integration

**User Story:** As an MDM administrator, I want the dashboard to retrieve live data from the backend, so that the statistics and device list reflect the current state of the system.

#### Acceptance Criteria

1. WHEN the Dashboard_Page mounts, THE Dashboard_Page SHALL fetch summary statistics by sending a `GET` request to `/rest/private/summary`.
2. WHEN Summary_API returns HTTP 200, THE Dashboard_Page SHALL parse the response body and extract `deviceTotal`, `deviceOnline`, `deviceEnrolled`, `configurationCount`, and `applicationCount`.
3. THE Dashboard_Page SHALL use the shared `apiClient` instance for all API requests so that the `X-Auth-Token` header is attached automatically.
4. IF Summary_API returns a non-200 HTTP status, THEN THE Dashboard_Page SHALL display an error message indicating that statistics could not be loaded.
5. IF a network error occurs during any data fetch, THEN THE Dashboard_Page SHALL display an error message indicating that the server could not be reached.

---

### Requirement 6: Component and Style Standards

**User Story:** As a frontend developer, I want the dashboard to follow the project's established component and styling conventions, so that the codebase remains consistent and maintainable.

#### Acceptance Criteria

1. THE Dashboard_Page SHALL use shadcn/ui `Card`, `Skeleton`, and `Badge` components exclusively for their respective UI roles.
2. THE Dashboard_Page SHALL use Tailwind CSS utility classes for all layout and spacing, with no inline styles.
3. THE Dashboard_Page SHALL be implemented as a TypeScript React component with explicit type annotations for all props and API response shapes.
4. THE Dashboard_Page SHALL be located at `frontend/src/features/dashboard/DashboardPage.tsx`.
