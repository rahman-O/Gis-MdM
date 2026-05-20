# Settings API parity

| Method | Go path | Java | Status |
|--------|---------|------|--------|
| GET | `/rest/private/settings` | `SettingsResource.getSettings` | Done |
| POST | `/rest/private/settings/misc` | `updateMiscSettings` | Done |
| POST | `/rest/private/settings/lang` | `updateLanguageSettings` | Done |
| POST | `/rest/private/settings/design` | `updateDefaultDesignSettings` | Done |
| GET | `/rest/private/settings/userRole/:roleId` | `getUserRoleSettings` | Done (defaults) |
| POST | `/rest/private/settings/userRoles/common` | `updateUserRoleCommonSettings` | Done (no-op OK) |

**Java:** `backend/server/src/main/java/com/hmdm/rest/resource/SettingsResource.java`

Migration `000003_settings_extend.up.sql` adds columns used by misc/lang/design saves.
