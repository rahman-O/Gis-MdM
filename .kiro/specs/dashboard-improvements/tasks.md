# Implementation Plan: Dashboard Improvements

## Overview

Replace the placeholder `DashboardPage` with a functional overview page showing five summary stat cards, a recent devices table (last 5 by `lastUpdate`), loading skeletons, and 60-second auto-refresh. All API calls go through a dedicated `dashboardService` using the shared `apiClient` and `HmdmEnvelope` unwrapping pattern.

## Tasks

- [ ] 1. Define dashboard types
  - Create `src/features/dashboard/types.ts`
  - `SummaryStats`: `{ deviceTotal, deviceOnline, deviceEnrolled, configurationCount, applicationCount }` (all `number`)
  - `RecentDeviceRow`: `{ id: number, number: string, description: string | null, statusCode: string | null, lastUpdate: number | null }`
  - _Requirements: 1.2, 2.3, 5.2_

- [ ] 2. Implement dashboardService
  - Create `src/features/dashboard/dashboardService.ts`
  - `getSummaryDevices(): Promise<DeviceSummaryPayload>` — GET `/rest/private/summary/devices` (not `/private/summary`; device stats live on the devices-specific endpoint), unwrap via `unwrapHmdmData`
  - `getRecentDevices(): Promise<DeviceView[]>` — reuse `deviceService.getDevices()` with `{ pageNum: 1, pageSize: 5, sortBy: LAST_UPDATE, sortDir: desc }`
  - Re-throw on non-OK envelope status or network error
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5_

- [ ]* 2.1 Write unit tests for dashboardService
  - Mock `apiClient`; assert `getSummaryDevices` sends GET to `/private/summary/devices` (via `dashboardService.test.ts`)
  - Assert `getRecentDevices` sends POST to `/rest/private/devices/search` with `pageSize: 5`
  - Assert both functions throw when envelope `status !== "OK"`
  - _Requirements: 5.1, 5.2, 5.4_

- [ ] 3. Implement DashboardPage component
  - [ ] 3.1 Create `src/features/dashboard/DashboardPage.tsx` with state and data fetching
    - State: `summary: SummaryStats | null`, `recentDevices: RecentDeviceRow[]`, `loading: boolean`, `error: string | null`
    - On mount: set `loading = true`, call `Promise.all([getSummary(), getRecentDevices()])`, set state, set `loading = false`
    - Start `setInterval` (60 000 ms) that re-fetches both endpoints without touching `loading` state (background refresh)
    - Return cleanup function that calls `clearInterval`
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 5.1_

  - [ ] 3.2 Add five StatCard components to DashboardPage
    - Render five cards: Total Devices (`deviceTotal`), Online Devices (`deviceOnline`), Enrolled Devices (`deviceEnrolled`), Configurations (`configurationCount`), Applications (`applicationCount`)
    - Each card uses shadcn/ui `Card` with a Lucide icon, numeric value, and label
    - Show `Skeleton` in place of the numeric value while `loading === true`
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 3.1, 3.3_

  - [ ] 3.3 Add RecentDevicesTable to DashboardPage
    - Render shadcn/ui `Table` with columns: Device Name, Last Seen, Status
    - Show 5 `Skeleton` rows while `loading === true`
    - Map `statusCode` to a `Badge`: `"green"` → "Online" (default variant), `"red"` → "Offline" (destructive), other → "Unknown" (secondary)
    - Format `lastUpdate` (Unix ms) as a human-readable date/time string; show "—" if null
    - Show "No recent devices" row when `recentDevices` is empty
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 3.2, 3.3_

  - [ ] 3.4 Add error handling to DashboardPage
    - When initial fetch fails, set `error` and show an error banner with a "Retry" button that re-triggers the initial fetch
    - When background refresh fails, silently ignore (do not overwrite existing data or show error)
    - _Requirements: 3.4, 5.4, 5.5_

- [ ]* 4. Write unit tests for DashboardPage
  - Assert 5 stat card skeletons rendered while loading
  - Assert 5 skeleton table rows rendered while loading
  - Assert error banner shown when initial fetch rejects
  - Assert loaded data displayed after successful fetch (no skeletons)
  - Assert `clearInterval` called on unmount
  - Assert background refresh does not show skeletons
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 4.3, 4.4_

- [ ]* 5. Write property-based tests (fast-check)

  - [ ]* 5.1 Write property test for stat card values match API response (Property 1)
    - **Property 1: Stat card values match API response**
    - **Validates: Requirements 1.2, 5.2**
    - For any `fc.record` of non-negative integers, assert each of the 5 cards displays the exact corresponding value

  - [ ]* 5.2 Write property test for recent devices sorted by lastUpdate descending (Property 2)
    - **Property 2: Recent devices sorted by lastUpdate descending**
    - **Validates: Requirements 2.2**
    - For any shuffled array of `RecentDeviceRow` with distinct `lastUpdate` values, assert rendered rows appear in descending order

  - [ ]* 5.3 Write property test for device row renders required fields (Property 3)
    - **Property 3: Device row renders required fields**
    - **Validates: Requirements 2.3**
    - For any `RecentDeviceRow` (including nullable fields), assert the rendered row contains device name, last seen, and a status badge

  - [ ]* 5.4 Write property test for badge correctly reflects online status (Property 4)
    - **Property 4: Badge correctly reflects online status**
    - **Validates: Requirements 2.4, 2.5**
    - For any `statusCode` value, assert the badge label is "Online" for `"green"`, "Offline" for `"red"`, "Unknown" for all others

  - [ ]* 5.5 Write property test for auto-refresh re-fetches after 60 seconds (Property 5)
    - **Property 5: Auto-refresh re-fetches after 60 seconds**
    - **Validates: Requirements 4.2**
    - Using `vi.useFakeTimers`, for any N elapsed 60-second intervals, assert each service function was called N+1 times total (1 initial + N refreshes)

  - [ ]* 5.6 Write property test for previously loaded data shown during background refresh (Property 6)
    - **Property 6: Previously loaded data shown during background refresh**
    - **Validates: Requirements 4.4**
    - For any previously loaded `SummaryStats`, assert no `Skeleton` components are rendered while a background refresh is in flight

  - [ ]* 5.7 Write property test for fetch failure shows error message (Property 7)
    - **Property 7: Fetch failure shows error message**
    - **Validates: Requirements 5.4, 5.5**
    - For any HTTP error status (400–599) or network rejection, assert an error message is visible in the rendered output

- [ ] 6. Checkpoint — Ensure all tests pass
  - Ensure all tests pass; ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- `statusCode` is a string from the backend (`"green"`, `"red"`, etc.) — same convention as the devices feature
- Background refresh must NOT set `loading = true` to avoid skeleton flash on every 60-second tick
- `getRecentDevices` reuses the existing `POST /rest/private/devices/search` endpoint — no new backend endpoint needed
- Property tests use fast-check; each test is tagged with its property number and requirement clause
- shadcn/ui components `Card`, `Skeleton`, `Badge`, and `Table` are already scaffolded — no `npx shadcn` commands needed
