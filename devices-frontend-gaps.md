# Devices Frontend Gaps (vs Backend)

This document lists missing items in the new `frontend` devices feature compared to backend capabilities in `DeviceResource`.

## 1) Missing Device Endpoints in Frontend

The backend exposes additional endpoints that are not implemented/used in the new frontend devices flow:

- `POST /private/devices/autocomplete`
- `POST /private/devices/deleteBulk`
- `POST /private/devices/groupBulk`
- `GET /private/devices/{id}/applicationSettings`
- `POST /private/devices/{id}/applicationSettings`
- `POST /private/devices/{id}/applicationSettings/notify`
- `POST /private/devices/{id}/description`

## 2) Search/Filter Capability Gap

Current frontend search request sends mainly:

- `pageNum`
- `pageSize`
- `value`

Backend `DeviceSearchRequest` supports many more filters/sort options not yet wired in frontend:

- `groupId`
- `configurationId`
- `sortBy`
- `sortDir`
- `dateFrom`, `dateTo`
- `onlineEarlierMillis`, `onlineLaterMillis`
- `enrollmentDateFrom`, `enrollmentDateTo`
- `mdmMode`
- `kioskMode`
- `androidVersion`
- `launcherVersion`
- `installationStatus`
- `imeiChanged`
- `fastSearch`

## 3) Update Device Flow Not Exposed in UI

- Frontend has `updateDevice(...)` in `deviceService`.
- No actual edit form/page action in `DevicesPage` or `DeviceDetailPanel` currently calls it.

## 4) Status Mapping Is Incomplete

Backend allows status codes:

- `green`
- `red`
- `yellow`
- `brown`
- `grey`

Current frontend `StatusBadge` maps only:

- `green` -> Online
- `red` -> Offline
- all others -> Unknown

## 5) Device Details Coverage Gap

Backend `DeviceInfoView` includes more data than what the current detail panel renders, such as:

- `model`
- `permissions`
- `applications`
- `files`
- `defaultLauncher`

Current `DeviceDetailPanel` displays a subset only (number, status, configuration, groups, battery, last update, location fallback).

## 6) Location Contract Mismatch Risk

- Frontend expects `device.info.location.lat/lon`.
- Backend `DeviceInfoView` does not explicitly expose a `location` field.
- Result: location may remain unavailable unless backend contract is extended or mapped differently.

## 7) Missing Device Table Columns and Advanced Statuses

Compared to the old devices page, the new frontend currently misses multiple table columns and advanced computed statuses.

### Missing status columns

- `Permission Status`
  - Old behavior: computes permission health from `device.info.permissions` (3 checks) and kiosk/permissive mode.
  - UI indicator mapping in old app: green/amber/red icon + tooltip details.
- `Installation Status`
  - Old behavior: compares configuration applications vs installed apps from `device.info.applications`.
  - Detects: not installed, version mismatch, up-to-date, removal-required.
- `Files status`
  - Old behavior: compares configuration files vs `device.info.files`.
  - Detects: missing file, timestamp mismatch, good state.

### Other missing optional columns that existed in old UI

- IMEI
- Phone
- Model
- Description
- Launcher version
- Battery level
- Default launcher / Background mode
- MDM mode
- Kiosk mode
- Android version
- Enrollment date
- Serial
- Public IP
- Custom1 / Custom2 / Custom3

## 8) Missing Actions in Devices Row

The old page had richer row actions than the current React page.

- `Edit`
  - Present in old UI (`button.change`), opens full edit modal.
  - Uses `PUT /private/devices`.
  - Current React page: not implemented as a UI action.
- `QR code`
  - Present in old UI (`button.qrcode`), opens enrollment QR flow (using config `qrCodeKey` + base URL).
  - Current React page: not implemented.
- `Delete`
  - Present and implemented in current React page.
- `More` dropdown in old UI:
  - `Application settings` (device-level app settings CRUD + notify).
  - Plugin-driven actions (device context actions), e.g. detailed info/logs/messaging/push.

## 8.1) Missing Header Action: Add Device

Old devices page had an explicit **Add Device** button in the header.

- Old behavior:
  - Opens device modal in create mode.
  - Submits through `PUT /private/devices` with new device payload.
  - Refreshes list after successful create.
- Current React page:
  - No `Add Device` button and no create-device dialog flow.

## 9) Missing Application Settings Feature (Core Device Action)

Old frontend implemented full device application settings management from device row actions.

- Endpoints:
  - `GET /private/devices/{id}/applicationSettings`
  - `POST /private/devices/{id}/applicationSettings`
  - `POST /private/devices/{id}/applicationSettings/notify`
- Capabilities in old app:
  - List settings
  - Add / edit / delete setting items
  - Save settings
  - Send notification to device after update
- Current React page: not implemented.

## 10) Missing Bulk Operations from Old UI

Old devices page supported grouped actions for selected rows:

- Bulk set configuration
  - `PUT /private/devices` with `{ ids, configurationId }`.
- Bulk set/clear groups
  - `POST /private/devices/groupBulk` with `{ ids, action, groups }`.
- Bulk delete
  - `POST /private/devices/deleteBulk` with `{ ids }`.

Current React page: row selection + bulk action bar is not implemented.

## 11) Detailed Information / Logs / Messaging / Push Messages Are Plugin-Based

Items like:

- Detailed information
- Logs
- Messaging
- Push messages

were exposed in old UI through plugin device-actions, not hardcoded core device endpoints.

### How old app wired this

- It loaded available plugins from:
  - `GET rest/plugin/main/private/available`
- For each plugin with:
  - `enabledForDevice = true`
  - `functionsViewTemplate` present
  - permission granted (`deviceFunctionsPermission`)
  it rendered an action under the row `More` menu.
- On click, it emitted a plugin-specific event:
  - `plugin-{identifier}-device-selected`
  and navigated to plugin screen with selected device context.

### Relevant plugin modules and APIs

- `deviceinfo` plugin (Detailed information):
  - `GET rest/plugins/deviceinfo/deviceinfo/private/:deviceNumber`
  - dynamic data search/export APIs
- `devicelog` plugin (Logs):
  - `POST rest/plugins/devicelog/log/private/search`
  - export API
- `messaging` plugin (Messaging):
  - search/send/delete message APIs under `rest/plugins/messaging/...`
- `push` plugin (Push messages):
  - search/send/delete/schedule APIs under `rest/plugins/push/...`

Current React page: plugin-based device actions are not implemented.

## 12) Device Description-Only Edit Path Missing

Old app supported edit of description only (for users without full edit permission):

- `POST /private/devices/{id}/description`

Current React page does not expose this permission-based fallback flow.

## 13) Missing Advanced Device Form Behavior

Old `device.html` modal includes behavior that is not present in the new frontend:

- Validation/constraints:
  - disallow certain characters in device number (`/ ? &`)
  - enforce configuration selection
  - enforce at least one group in restricted user mode
- Conditional edit behavior:
  - lock device number for migration scenarios (`oldNumber` handling)
  - explicit migration confirmation when changing number of an enrolled device
- Field rendering from settings:
  - custom property labels from common settings
  - optional multiline rendering per custom field (`customMultiline1..3`)
  - phone mask and IMEI mask handling

## 14) Missing Row Selection UX

Old devices page supports full row selection workflow:

- Row checkbox per device.
- Header "Select all" checkbox.
- Selection state used to enable/disable group bulk actions.
- Selection reset after list refresh.

Current React page has no row-selection model.

## 15) Missing Configuration Link Action in Table

Old devices table lets user click configuration name and jump to configuration editor (when permission allows).

- Old behavior: `editConfiguration(device.configuration)` transitions to config editor state.
- Current React page: configuration name is plain text only.

## 16) Missing Periodic Auto-Refresh

Old devices page refreshes list automatically every 60 seconds:

- Calls search in background (`spinnerHidden`) and keeps list up to date.

Current React page performs fetch on user/search/page actions only.

## 17) Missing Core Plugin Action Host in Devices Page

Beyond listing plugin names, old devices page acts as a host for device-scoped plugin actions:

- Loads plugin metadata from `plugin/main/private/available`.
- Filters by:
  - `enabledForDevice`
  - presence of `functionsViewTemplate`
  - permission gate `deviceFunctionsPermission`
- Emits selected device context to plugin channels:
  - `plugin-deviceinfo-device-selected`
  - `plugin-devicelog-device-selected`
  - `plugin-messaging-device-selected`
  - `plugin-push-device-selected`

Current React devices page has no equivalent plugin action host architecture.

## 18) Missing Persistent Search State (Cookies Restore)

Old devices page persists and restores search context using cookies:

- saved state includes:
  - search params
  - paging
  - additional filters
  - group/config selections
- on reopening devices tab, previous context is restored automatically.

Current React page has no equivalent persistence/restore mechanism.

## 19) Missing License-Limit Guard on Add Device

Old devices header disables Add button when tenant license limit is reached:

- condition: `deviceLimit > 0 && deviceCount >= deviceLimit`
- prevents opening add flow when capacity is exhausted.

Current React page has no equivalent guard (and currently no add flow).

## 20) Missing Account-Expired Device List Handling

Old devices logic applies account-expiry behavior in list rendering:

- warns user about expiry
- marks account-expired mode
- progressively hides/fades rows after a threshold.

Current React devices page has no account-expiry rendering behavior.

## 21) Missing Fine-Grained Permission Matrix in UI

Old page uses granular permissions for action visibility, not just generic edit:

- `edit_devices` -> edit/delete/bulk actions
- `edit_device_desc` -> description-only edit path
- `enroll_devices` -> QR action visibility
- `edit_device_app_settings` -> Application settings action
- `configurations` -> clickable configuration link access
- plugin action permissions via `deviceFunctionsPermission`

Current React page does not replicate this full permission matrix.

## 22) Missing User-Scoped Configuration Availability Filtering

Old page restricts configuration visibility per user scope:

- `availableConfigs` derived from user grants
- `configAvailable(config)` gate for configuration navigation/access.

Current React page does not implement equivalent user-scoped config gating.

## 23) Missing IMEI/Phone Consistency Diagnostics

Old page compares server fields against device-reported info and shows mismatch diagnostics:

- `displayedIMEI`, `imeiTooltip`, `imeiTooltipClass`
- `displayedPhone`, `phoneTooltip`, `phoneTooltipClass`
- highlights conformance issues and missing data conditions.

Current React table shows plain values without mismatch diagnostics.

## 24) Missing Rich Status Tooltip ("Last Seen Ago")

Old status indicator includes humanized elapsed tooltip:

- minutes/hours/days/weeks/months/years ago
- exact formatted timestamp on a second line.

Current React status badge does not provide equivalent rich status tooltip.

## 25) Missing Launcher Version Mismatch Highlight

Old devices table colors launcher version when installed launcher version differs from required configuration version.

Current React page does not provide launcher version mismatch highlighting.

## 26) Missing Dual Pagination Placement

Old devices page renders pagination block both above and below the table.

Current React page renders pagination once (bottom area only).

## 27) Missing Explicit Manual Search Trigger UX

Old page includes explicit Search button in addition to keyboard submit and fast-search checkbox workflow.

Current React page relies on debounced input only and does not expose equivalent explicit search trigger UX.

## Suggested Priority

1. Implement advanced search/filter + sorting.
2. Add bulk actions (`deleteBulk`, `groupBulk`) in devices table.
3. Add edit device UI that uses `updateDevice`.
4. Extend detail panel with key backend fields.
5. Normalize status mapping for all backend status codes.
