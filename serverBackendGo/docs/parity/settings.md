# Settings API parity

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| GET | `/rest/private/settings` | `SettingsResource.getSettings` | Done |
| POST | `/rest/private/settings/misc` | `updateMiscSettings` | Done |
| POST | `/rest/private/settings/lang` | `updateLanguageSettings` | Done |
| POST | `/rest/private/settings/design` | `updateDefaultDesignSettings` | Done |
| GET | `/rest/private/settings/userRole/:roleId` | `getUserRoleSettings` | **Done** — `userrolesettings` table (`000012`) |
| POST | `/rest/private/settings/userRoles/common` | `updateUserRoleCommonSettings` | **Done** — upsert `userrolesettings` |

**Java:** `backend/server/src/main/java/com/hmdm/rest/resource/SettingsResource.java`

Migrations: `000003_settings_extend` (misc/lang/design), `000012_userrolesettings_core`, `000015_settings_columns_extend`.

**014:** GET/POST misc include tenant fields from `000015`: `newDeviceGroupId`, `phoneNumberFormat`, `customPropertyName1`–`3`, `customMultiline1`–`3`, `customSend1`–`3`, `desktopHeaderTemplate`, `sendDescription`. React `SettingsPage` edits and saves via existing misc/lang paths.
