# Configurations Legacy Details (Source: `backend`)

This document captures how **Configurations** work in the legacy Headwind app (AngularJS + Java backend), to serve as the migration source of truth for the new frontend/backend.

---

## 1) Core Backend API (`/rest/private/configurations`)

Primary resource:
- `backend/server/src/main/java/com/hmdm/rest/resource/ConfigurationResource.java`

### Endpoints
- `GET /search`
  - Returns all configurations (requires `configurations` permission).
  - Backend injects `baseUrl` into each configuration.
- `GET /list`
  - Returns lightweight `LookupItem[]` (`id`, `name`).
  - Used for selectors and allowed even with broad permissions.
- `GET /search/{value}`
  - Filtered search by text.
- `POST /autocomplete`
  - Returns lookup items for autocomplete.
- `PUT /`
  - Create/update configuration.
  - Create path (`id == null`) requires `add_config`.
  - Update path requires access check by configuration (`hasConfigurationAccess`).
  - On update, triggers `pushService.notifyDevicesOnUpdate(configurationId)`.
- `PUT /copy`
  - Copy configuration (requires `copy_config` and access to source config).
- `DELETE /{id}`
  - Delete configuration (requires `copy_config` + access).
  - Can fail with device-reference constraint (`CONFIGURATION_DEVICE_REFERENCE_EXISTS`).
- `PUT /application/upgrade`
  - Upgrade specific app used in configuration to latest.
- `GET /applications`
  - Returns all applications.
- `GET /applications/{id}`
  - Returns app list in context of a configuration.
- `GET /{id}`
  - Full configuration payload used by editor screen.

### Permissions/Gates
- `configurations` (view/edit base access)
- `add_config` (create)
- `copy_config` (copy/delete)
- `hasConfigurationAccess(id)` for per-config authorization

---

## 2) Legacy Frontend Service Contract (Angular)

Service file:
- `backend/server/src/main/webapp/app/components/main/service/main.service.js`

`configurationService` methods used by UI:
- `getAllConfigurations` -> `GET rest/private/configurations/search/:value`
- `getAllConfigNames` -> `GET rest/private/configurations/list`
- `getById` -> `GET rest/private/configurations/:id`
- `getAllTypicalConfigurations` -> `GET rest/private/configurations/typical/search/:value`
- `updateConfiguration` -> `PUT rest/private/configurations`
- `copyConfiguration` -> `PUT rest/private/configurations/copy`
- `upgradeConfigurationApplication` -> `PUT rest/private/configurations/:id/application/:appId/upgrade` (legacy route shape)
- `removeConfiguration` -> `DELETE rest/private/configurations/:id`
- `getApplications` -> `GET rest/private/configurations/applications/:id`

---

## 3) Legacy Configurations List Behavior

Controller:
- `backend/server/src/main/webapp/app/components/main/controller/configurations.controller.js`
  - `ConfigurationsTabController`

Key behavior:
- Supports normal and typical tabs (`isTypical`).
- Search by `searchValue`.
- Checks permissions with `authService.hasPermission`.
- QR availability rule in list:
  - `qrCodeKey` exists
  - `mainAppId > 0`
  - `eventReceivingComponent` non-empty
- QR action opens:
  - `configuration.baseUrl + "/#/qr/" + configuration.qrCodeKey + "/"`
- Add config shows warning confirmation first.
- Edit navigates to `configEditor`.
- Copy opens modal (name required, duplicate-name protection).
- Delete confirmation varies if list has only one configuration.

---

## 4) Legacy Configuration Editor (Very Important)

Main files:
- Controller: `configurations.controller.js` -> `ConfigurationEditorController`
- View: `backend/server/src/main/webapp/app/components/main/view/configuration.html`

Editor is multi-tab and large; this is where most migration gaps exist.

### 4.1 Common Settings tab
- Name (required)
- Description
- Admin password (required)
- Request updates mode (`DONOTTRACK`, `GPS`, `WIFI`)
- App permissions (`GRANTALL`, `ASKLOCATION`, `DENYLOCATION`, `ASKALL`)
- Push options (`mqttAlarm`, `polling`; legacy `mqttWorker` hint logic exists)
- MQTT keepalive (for specific push option)
- Device radios/toggles:
  - GPS, Bluetooth, Wi-Fi, Mobile Data (tri-state patterns)
  - USB storage
  - Brightness management + value slider
  - Timeout management + value
  - Volume lock/manage + value
  - Timezone mode + manual timezone string
  - System update strategy + scheduled time window
  - App update schedule + time window
  - Download updates policy
  - Password mode
  - Show Wi-Fi
  - Use default launcher
  - Disable screenshots
  - Autostart foreground

### 4.2 Design Settings tab
- Use default design settings toggle
- Background color / text color
- Background image URL (+ upload helper)
- Icon size (`SMALL`, `MEDIUM`, `LARGE`)
- Desktop header mode:
  - `NO_HEADER`, `DEVICE_ID`, `DESCRIPTION`, `TEMPLATE`, `CUSTOM1/2/3`
- Desktop header template text
- Orientation
- Display status toggle

### 4.3 Applications tab
- App list with actions:
  - install / not install / remove / permit / prohibit
- Per-app version selection and upgrade flow
- Main/content app synchronization logic
- Filtering, sorting, pagination
- Show/hide system apps
- Add app modal (can create new app from modal)
- App details modal

### 4.4 MDM Settings tab (critical for QR)
- `mainAppId` selector (typeahead)
- `eventReceivingComponent` input
- `contentAppId` selector (kiosk mode)
- Kiosk behavior flags:
  - home, recents, notifications, system info, keyguard, lock buttons, exit, keep screen on
- Provisioning Wi-Fi settings:
  - SSID, password, security type
- Launcher URL override
- QR parameters (`qrParameters`)
- Admin extras (`adminExtras`)
- Mobile enrollment
- Encrypt device
- Permissive mode + restrictions + allowed classes + lock safe settings
- New server URL for migration
- QR URL display condition:
  - only if main app + event receiver set
  - otherwise hint text is shown

### 4.5 App Settings tab
- Configuration-level `applicationSettings` CRUD
- Filter by app and text
- Modal editor with validation (name + app required)

### 4.6 Files tab
- `defaultFilePath`
- Linked files list CRUD per configuration
- Upload/create/edit/remove file links
- Replace variables / remove flags

---

## 5) Save Validation Rules (Legacy)

From `ConfigurationEditorController.save()`:
- `pushOptions` required
- `name` required
- `password` required
- If `kioskMode` then content app must be selected
- Protects against lost/invalid `mainApp` state (`bConfigurationWasLost`)
- Normalizes outgoing request:
  - `request.type = isTypical ? 1 : 0`
  - app lists, update windows, timezone mode conversions
  - null normalization for optional fields

---

## 6) QR-Specific Legacy Rules

QR works only when configuration is complete for provisioning:
- `qrCodeKey` exists
- `mainAppId` is set
- `eventReceivingComponent` exists
- main app has valid downloadable URL/version

Server evidence:
- `backend/server/src/main/java/com/hmdm/rest/resource/QRCodeResource.java`
  - Logs and returns empty/non-image when prerequisites are missing.
  - Explicitly checks for main app and URL before generating QR PNG.

---

## 7) Configuration Data Model (Backend Domain)

Model file:
- `backend/common/src/main/java/com/hmdm/persistence/domain/Configuration.java`

Contains broad field groups:
- Common management/security/network/device behavior
- MDM enrollment/provisioning and kiosk controls
- Design/theming
- Applications + app usage parameters
- App settings
- Files/default path
- QR key/base URL metadata

This model is significantly richer than current MVP UI in new frontend.

---

## 8) Migration Gaps To Implement In New Project

High-priority:
- Full Configuration editor (tabbed or sectional) parity with legacy.
- Main App + Event Receiving Component selection/editing.
- QR eligibility feedback and actionable fixes in UI.
- Application management in config (action, version, order, details).
- App settings CRUD.
- Files tab + links + default file path.
- Type handling (`typical/common`) and permission gating.
- Save normalization logic matching legacy.

Medium-priority:
- Upgrade app-in-config flow.
- Design tab parity.
- Full validation parity and helper hints.
- Legacy onboarding/warnings/dirty-form leave guards.

---

## 9) Recommended Implementation Order (for new frontend)

1. Expand `Configuration` types to full backend contract.
2. Add missing service endpoints and payload normalization.
3. Build minimal MDM block first:
   - Main App
   - Event Receiving Component
   - Content App
   - QR-related fields
4. Add save validation parity.
5. Add Applications tab logic.
6. Add App Settings and Files tabs.
7. Add design/common advanced fields.
8. Add permission-based UI gating.

---

## 10) Source Files Referenced

- `backend/server/src/main/java/com/hmdm/rest/resource/ConfigurationResource.java`
- `backend/common/src/main/java/com/hmdm/persistence/domain/Configuration.java`
- `backend/server/src/main/java/com/hmdm/rest/resource/QRCodeResource.java`
- `backend/server/src/main/webapp/app/components/main/service/main.service.js`
- `backend/server/src/main/webapp/app/components/main/controller/configurations.controller.js`
- `backend/server/src/main/webapp/app/components/main/view/configuration.html`
- `backend/server/src/main/webapp/app/components/main/controller/qr.controller.js`

