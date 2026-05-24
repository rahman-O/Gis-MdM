# Contract: Frontend Control Plane UX

**Feature**: `017-device-control-plane` | **Audience**: React admin (`frontend/src`)

## Navigation (`navItems.ts`)

| Before | After |
|--------|-------|
| Configurations | **Profiles** (`/profiles`) |
| — | **Enrollment routes** (`/enrollment-routes`) |
| Devices (flat) | Devices with **tree sidebar** |

**FR-010a**: No parallel «Configurations» menu item.

## Routes

| Path | Page |
|------|------|
| `/devices` | Tree + device table |
| `/profiles` | Profile list |
| `/profiles/:id` | Profile editor (tabs: Restrictions, MDM, Apps, Design, Files) |
| `/enrollment-routes` | Route list |
| `/enrollment-routes/:id` | Route editor (binding only) + QR panel |
| `/onboarding` | Checklist / wizard (P2) |

## Devices page

- Left: `DeviceTreeSidebar` — expand/collapse, select node
- Right: filtered `DevicesTable`
- Actions: Move to folder (dialog), bulk move (later)

## Profile editor

Reuse `frontend/src/features/configurations/*` components renamed/moved to `features/profiles/*`:

- Same tab semantics as 016
- Header: **Usage panel** — device count, enrollment route count
- Publish button → impact preview → confirm if ≥50 devices

## Enrollment route editor

- Fields: name, profile version picker (published only), tree folder picker, device id mode, main app
- QR: `EnrollmentQrExperience` with `createOnDemand={true}` always
- Copy: Arabic/English helper — restrictions live in Profile

## Delete tree folder dialog

1. User clicks delete on node with devices
2. Modal: pick target folder (tree picker)
3. Confirm → API `POST .../delete` with `targetNodeId`
4. Refresh tree + table

## Onboarding (P2)

`OnboardingChecklist` on dashboard when:

- No published profile, OR
- No enrollment route, OR
- No tree beyond root

Wizard steps: Tree folder → Profile → Publish → Route → QR test.

## i18n keys (suggested)

- `nav.profiles`, `nav.enrollmentRoutes`
- `enrollmentRoute.help.profileOnly`
- `profile.publish.confirm.title` (≥50 devices)
