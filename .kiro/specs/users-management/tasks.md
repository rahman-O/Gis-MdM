# Implementation Plan: Users Management

## Overview

Implement the `/users` page following the feature-slice pattern established by configurations-management. The route and navigation entry already exist; tasks proceed from scaffolding missing shadcn/ui components through types, service, form, page, and final wiring.

## Tasks

- [ ] 1. Scaffold missing shadcn/ui components
  - Run `npx shadcn@latest add dialog select` inside `frontend/` to generate the two missing components into `src/shared/ui/`
  - Verify `dialog.tsx` and `select.tsx` exist in `src/shared/ui/`
  - _Requirements: 2.2, 5.1_

- [ ] 2. Define TypeScript types
  - [ ] 2.1 Create `src/features/users/types.ts`
    - Define `Role`, `User`, and `UserPayload` interfaces exactly as specified in the design
    - `User.role` is `Role | null`; `allDevicesAvailable` and `customerStaff` are `boolean` (null/absent from backend treated as `false`)
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5_

- [ ] 3. Implement the user service
  - [ ] 3.1 Create `src/features/users/userService.ts`
    - Implement `getUsers()`, `createUser(data)`, `updateUser(id, data)`, `deleteUser(id)`, `getRoles()` using the shared `apiClient` and `unwrapHmdmData` helper
    - All non-2xx responses or `status: "ERROR"` envelopes must propagate as thrown errors
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7_

  - [ ]* 3.2 Write property test for service URL routing (Property 9)
    - **Property 9: Service routes to correct URL for any operation**
    - **Validates: Requirements 8.1, 8.2, 8.3, 8.4, 8.5**
    - Use `fc.integer({ min: 1 })` for ids and `arbitraryUserPayload()` for payloads; mock `apiClient` and assert method + URL for each function; run 100 iterations

  - [ ]* 3.3 Write property test for service error propagation (Property 10)
    - **Property 10: Service error propagation**
    - **Validates: Requirements 8.7**
    - For each service function, mock `apiClient` to reject and assert the returned promise rejects; run 100 iterations

- [ ] 4. Implement UserForm
  - [ ] 4.1 Create `src/features/users/UserForm.tsx`
    - Wrap a shadcn/ui `Dialog` around a `Form` (react-hook-form + zod schemas from design)
    - Fields: `login` (Input, required), `name` (Input, required), `email` (Input, required, email pattern), `password` (Input type=password, required in create mode only), `roleId` (Select/Role_Select, required), `allDevicesAvailable` (Checkbox/Switch, default false), `customerStaff` (Checkbox/Switch, default false)
    - Call `userService.getRoles()` on mount; disable Role_Select and show loading state while roles are loading; show error if roles fetch fails
    - Manage `submitting`, `submitError`, `roles`, and `rolesLoading` state
    - Disable submit button while submitting; show `FormMessage` for validation errors adjacent to each field
    - In create mode all fields start empty; in edit mode pre-populate from `initialData` (including pre-selecting the matching role)
    - Call `onSuccess` + `onClose` on successful submit; call `onClose` on cancel without any API call
    - _Requirements: 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8, 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 5.1, 5.2, 5.3, 5.4, 5.5, 6.1, 6.2, 6.3, 6.4, 6.5_

  - [ ]* 4.2 Write property test for form submit routing (Property 2)
    - **Property 2: Form submit calls correct endpoint with form values**
    - **Validates: Requirements 2.4, 3.3**
    - Use `arbitraryUserPayload()` to generate inputs; assert create mode calls `createUser` with exact payload and edit mode calls `updateUser` with id + payload; run 100 iterations

  - [ ]* 4.3 Write property test for cancel makes no API call (Property 3)
    - **Property 3: Cancel makes no API call**
    - **Validates: Requirements 2.8, 4.6**
    - For any dialog state, clicking Cancel must not invoke any service function; run 100 iterations

  - [ ]* 4.4 Write property test for edit mode pre-population (Property 4)
    - **Property 4: Edit mode pre-populates fields for any user**
    - **Validates: Requirements 3.1**
    - Use `arbitraryUser()` to generate users; assert rendered form fields match `login`, `name`, `email`, `allDevicesAvailable`, `customerStaff`, and pre-selected `role.id`; run 100 iterations

  - [ ]* 4.5 Write property test for role select population and pre-selection (Property 7)
    - **Property 7: Role select populates and pre-selects correctly**
    - **Validates: Requirements 5.3, 5.4**
    - Use `fc.array(arbitraryRole(), { minLength: 1 })` for roles and `arbitraryUser()` for edit mode; assert options match roles and pre-selection matches user's role.id; run 100 iterations

  - [ ]* 4.6 Write property test for required field validation (Property 8)
    - **Property 8: Required field validation rejects invalid inputs**
    - **Validates: Requirements 5.5, 6.1, 6.2, 6.3, 6.4**
    - Use `arbitraryEmptyOrWhitespace()` for login/name/password, invalid email strings, and absent roleId; assert form rejects and no API call is made; run 100 iterations

- [ ] 5. Implement UsersPage
  - [ ] 5.1 Replace `src/features/users/UsersPage.tsx` with full implementation
    - Fetch users on mount via `userService.getUsers()`; manage `users`, `loading`, `error`, `formMode`, `selectedUser`, and `userToDelete` state
    - Render a shadcn/ui `Table` with columns: Login, Name, Email, Role, Status; each row has a `DropdownMenu` with Edit and Delete actions
    - Show loading skeleton while fetching; show error banner with Retry button on failure; show "No users found" empty state
    - Render "Add User" button above the table that opens `UserForm` in create mode
    - Inline `DeleteDialog` (wrapping `AlertDialog`) that shows the user's login and name, manages `deleting`/`deleteError` state, calls `userService.deleteUser`, and refreshes on success
    - After any successful create, update, or delete, re-fetch the user list
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 2.1, 2.5, 3.4, 4.1, 4.3, 4.4_

  - [ ]* 5.2 Write property test for table columns (Property 1)
    - **Property 1: Table columns rendered for any user list**
    - **Validates: Requirements 1.3**
    - Use `fc.array(arbitraryUser(), { minLength: 1 })` to generate lists; assert rendered table always has Login, Name, Email, Role, and Status columns; run 100 iterations

  - [ ]* 5.3 Write property test for delete dialog shows user identity (Property 5)
    - **Property 5: Delete dialog shows user login and name for any user**
    - **Validates: Requirements 4.1**
    - Use `arbitraryUser()` to generate users; open delete dialog and assert both login and name appear in dialog body; run 100 iterations

  - [ ]* 5.4 Write property test for delete confirm calls correct id (Property 6)
    - **Property 6: Delete confirm calls DELETE with correct id**
    - **Validates: Requirements 4.2**
    - Use `arbitraryUser()` to generate users; confirm deletion and assert `deleteUser` called with exact `id`; run 100 iterations

  - [ ]* 5.5 Write property test for null-safe rendering (Property 11)
    - **Property 11: Null-safe rendering for any user with null optional fields**
    - **Validates: Requirements 9.4, 9.5**
    - Use `arbitraryUserWithNulls()` (role null, booleans null); assert rendering table row and opening form does not throw; run 100 iterations

- [ ] 6. Checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all unit and property tests pass; ask the user if questions arise.

- [ ] 7. Verify routing and navigation
  - [ ] 7.1 Confirm `/users` route is registered in `src/app/App.tsx` (already present — verify only, no change needed)
    - _Requirements: 1.1, 7.2_

  - [ ] 7.2 Confirm "Users" entry is present in `src/features/layout/navItems.ts` (already present — verify only, no change needed)
    - _Requirements: 7.1, 7.2, 7.3_

- [ ] 8. Final checkpoint — Ensure all tests pass
  - Run `vitest --run` inside `frontend/`; ensure all tests pass end-to-end; ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for a faster MVP
- Each task references specific requirements for traceability
- Property tests use fast-check with a minimum of 100 iterations each
- Each property test file must include a comment referencing the design property number and the requirements clause it validates
- Generators (`arbitraryRole`, `arbitraryUser`, `arbitraryUserPayload`, `arbitraryEmptyOrWhitespace`, `arbitraryUserWithNulls`) should be defined once in a shared test helper file and imported by all property test files
- The `/users` route and "Users" nav entry already exist — tasks 7.1 and 7.2 are verification-only
- `dialog.tsx` and `select.tsx` are the only shadcn/ui components that need to be added
