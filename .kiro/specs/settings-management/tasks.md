# Implementation Plan: Settings Management

## Overview

Implement the `/settings` page following the feature-slice pattern established by configurations-management and users-management. The route and navigation entry already exist; tasks proceed from scaffolding missing shadcn/ui components through types, service, page implementation, and final wiring.

## Tasks

- [x] 1. Scaffold missing shadcn/ui components
  - `select.tsx` was already present; added Radix-based `switch.tsx`, `toast.tsx`, `toaster.tsx`, and `src/shared/hooks/use-toast.ts` (equivalent to `npx shadcn add toast switch`)
  - Added `<Toaster />` to `src/features/layout/AppLayout.tsx` so toasts render globally
  - _Requirements: 3.3, 3.4, 5.1, 5.2, 5.3, 5.4_

- [x] 2. Define TypeScript types
  - [x] 2.1 Create `src/features/settings/types.ts`
    - `Settings`, `SettingsPayload`, and `ConfigurationOption` per design
    - _Requirements: 8.1, 8.2_

- [x] 3. Implement the settings service
  - [x] 3.1 Create `src/features/settings/settingsService.ts`
    - `GET /private/settings` (axios base URL `/rest`) — `unwrapHmdmData`
    - **Persist:** legacy backend uses `POST /private/settings/misc` then `POST /private/settings/lang` (not `PUT /private/settings`); service merges payload, then refetches with `GET`
    - Errors propagate via axios / envelope helpers
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [x]* 3.2 Write property test for service URL routing (Property 6)
  - [x]* 3.3 Write property test for service error propagation (Property 7)

- [x] 4. Implement SettingsPage
  - [x] 4.1 Replace `src/features/settings/SettingsPage.tsx` with full implementation
    - Parallel load: `Promise.allSettled([settingsService.getSettings(), apiClient.get('/private/configurations/search')])` (configurations use search endpoint used elsewhere in the app)
    - Loading skeleton, error banner + Retry, form with zod + RHF
    - Toasts on save success / failure
  - [x]* 4.2 Property 1 — `SettingsPage.property.test.tsx`
  - [x]* 4.3 Property 2
  - [x]* 4.4 Property 3
  - [x]* 4.5 Property 4
  - [x]* 4.6 Property 5
  - [x]* 4.7 Property 8

- [x] 5. Checkpoint — Ensure all tests pass
  - `npm run test` (vitest `--run`) in `frontend/` passes

- [x] 6. Verify routing and navigation
  - [x] 6.1 `/settings` in `src/app/App.tsx` — verified
  - [x] 6.2 `Settings` in `src/features/layout/navItems.ts` — verified

- [x] 7. Final checkpoint — Ensure all tests pass

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- **Backend save path:** Headwind HMDM saves miscellaneous + language settings via two POST endpoints; the original plan’s single `PUT` is not exposed on the Java resource.
- Configurations list: `GET /private/configurations/search` (permission-gated on the backend).
- Generators live in `frontend/src/features/settings/test/arbitraries.ts`.
