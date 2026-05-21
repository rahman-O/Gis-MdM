# تحليل تكامل الفرونت (React) مع الباكند (Go)

**تاريخ التحليل:** 2026-05-21  
**الفرونت:** [`frontend/`](frontend/) — Vite + React + Axios  
**الباكند:** [`serverBackendGo/`](serverBackendGo/) — Gin، مسارات `/rest/*`  
**مرجع هجرة Java:** [`JAVA-GO-MIGRATION-STATUS.md`](JAVA-GO-MIGRATION-STATUS.md)

---

## 1. الخلاصة التنفيذية

| البُعد | الحكم |
|--------|--------|
| **توافق المسارات (URL)** | **~95%** — كل استدعاء REST في React له مسار مسجّل في Go (ما عدا WebSocket) |
| **توافق الحقول (payload/response)** | **~75–85%** — غالباً يعمل؛ فجوات في بحث الأجهزة، تفاصيل الجهاز، وفلاتر متقدمة |
| **توافق السلوك** | **~80%** — يعتمد على تفعيل الوحدات في `.env` واكتمال منطق Go |
| **جاهزية تشغيل متكامل** | **نعم للمسارات الرئيسية** مع ملاحظات أدناه |

**القاعدة:** الفرونت يضيف `baseURL: '/rest'`؛ البروكسي في [`frontend/vite.config.ts`](frontend/vite.config.ts) يوجّه `/rest` إلى `http://localhost:8080` (أو `VITE_BACKEND_ORIGIN`).

---

## 2. البنية المشتركة

### 2.1 عميل HTTP

| العنصر | الفرونت | الباكند Go |
|--------|---------|------------|
| الملف | [`frontend/src/services/apiClient.ts`](frontend/src/services/apiClient.ts) | — |
| Base URL | `/rest` | Router groups: `/rest/public`, `/rest/private`, `/rest/plugin/main`, `/rest/plugins`, `/rest/notifications` |
| Cookies | `withCredentials: true` | Session cookie على login |
| JWT | `Authorization: Bearer` من `localStorage` | Middleware على `/rest/private/*` |
| Content-Type | `application/json` (ما عدا FormData) | `ShouldBindJSON` / multipart |

### 2.2 غلاف الاستجابة (Envelope)

| الحقل | الفرونت | Go |
|-------|---------|-----|
| `status` | `'OK' \| 'ERROR'` في [`hmdmEnvelope.ts`](frontend/src/services/hmdmEnvelope.ts) | `response.OK` / `response.ErrorEnvelope` |
| `message` | مفتاح i18n أو نص | نفس أسلوب Java |
| `data` | `unwrapHmdmData()` يتطلب `data` عند OK | `response.OK(c, payload)` |

### 2.3 المصادقة

| الخطوة | مسار الفرونت | Go module | ملاحظة |
|--------|--------------|-----------|--------|
| خيارات الدخول | `GET /public/auth/options` | `auth` | 404 → fallback في الفرونت |
| تسجيل الدخول | `POST /public/auth/login` `{ login, password }` | `auth` | كلمة المرور MD5 أو RSA حسب `TRANSMIT_PASSWORD` |
| JWT (اختياري) | — (الفرونت **لا يستدعي** `/public/jwt/login`) | `auth` | متوفر في Go؛ React يستخدم session + Bearer من login |
| الخروج | `POST /public/auth/logout` | `auth` | ✅ |
| المستخدم الحالي | `GET /private/users/current` | `users` | يحدّث permissions بعد login |
| 2FA | `/private/twofactor/*` | `twofactor` | يحجب private حتى التحقق |

### 2.4 بروكسي التطوير

```text
Browser → http://localhost:5173/rest/... 
       → Vite proxy → http://localhost:8080/rest/... (serverBackendGo)
```

متغيرات: `VITE_BACKEND_ORIGIN`, `VITE_BACKEND_CONTEXT`, `TOMCAT_PORT` (للـ Java القديم).

### 2.5 WebSocket — **غير متوافق**

| العنصر | الفرونت | Go |
|--------|---------|-----|
| الملف | [`frontend/src/services/websocket.ts`](frontend/src/services/websocket.ts) | ❌ لا يوجد `/rest/ws/connect` |
| الاستخدام | مثال/تجريبي؛ **لا يُستورد** في صفحات التطبيق الرئيسية | — |

**التوصية:** لا تعتمد على WebSocket مع Go حالياً؛ استخدم polling الإشعارات للوكلاء فقط.

---

## 3. خريطة ملفات الفرونت (172 ملف `src/**`)

### 3.1 طبقة الخدمات (REST) — نقطة الربط مع الباكند

| الملف | المجال | عدد الاستدعاءات تقريباً |
|-------|--------|-------------------------|
| [`services/authService.ts`](frontend/src/services/authService.ts) | auth | 4 |
| [`services/apiClient.ts`](frontend/src/services/apiClient.ts) | HTTP | — |
| [`services/hmdmEnvelope.ts`](frontend/src/services/hmdmEnvelope.ts) | parsing | — |
| [`features/auth/twoFactorAuthService.ts`](frontend/src/features/auth/twoFactorAuthService.ts) | 2FA | 4 |
| [`features/auth/signupPublicService.ts`](frontend/src/features/auth/signupPublicService.ts) | signup | 4 |
| [`features/auth/passwordResetPublicService.ts`](frontend/src/features/auth/passwordResetPublicService.ts) | password reset | 3 |
| [`features/users/userService.ts`](frontend/src/features/users/userService.ts) | users | 5 |
| [`features/roles/roleService.ts`](frontend/src/features/roles/roleService.ts) | roles | 4 |
| [`features/profile/profileService.ts`](frontend/src/features/profile/profileService.ts) | profile | 3 |
| [`features/customers/customersService.ts`](frontend/src/features/customers/customersService.ts) | customers | 2 |
| [`features/devices/deviceService.ts`](frontend/src/features/devices/deviceService.ts) | devices + groups list | 12 |
| [`features/groups/groupService.ts`](frontend/src/features/groups/groupService.ts) | groups | 4 |
| [`features/configurations/configurationService.ts`](frontend/src/features/configurations/configurationService.ts) | configurations | 14 |
| [`features/applications/services/applicationService.ts`](frontend/src/features/applications/services/applicationService.ts) | applications | 17 |
| [`features/applications/services/webUiFilesService.ts`](frontend/src/features/applications/services/webUiFilesService.ts) | files upload | 6 |
| [`features/files/filesService.ts`](frontend/src/features/files/filesService.ts) | files list | 2 |
| [`features/icons/iconsService.ts`](frontend/src/features/icons/iconsService.ts) | icons metadata | 3 |
| [`features/settings/settingsService.ts`](frontend/src/features/settings/settingsService.ts) | settings | 8 |
| [`features/hints/hintsService.ts`](frontend/src/features/hints/hintsService.ts) | hints | 3 |
| [`features/dashboard/dashboardService.ts`](frontend/src/features/dashboard/dashboardService.ts) | summary + counts | 3 |
| [`features/dashboard/summaryService.ts`](frontend/src/features/dashboard/summaryService.ts) | (re-export) | — |
| [`features/push/pushService.ts`](frontend/src/features/push/pushService.ts) | push API | 1 |
| [`features/updates/updatesService.ts`](frontend/src/features/updates/updatesService.ts) | updates | 1 |
| [`features/plugins/pluginService.ts`](frontend/src/features/plugins/pluginService.ts) | plugin platform | 3 |
| [`features/devices/qrImage.ts`](frontend/src/features/devices/qrImage.ts) | QR PNG | blob GET |
| [`features/devices/enrollmentQrQuery.ts`](frontend/src/features/devices/enrollmentQrQuery.ts) | مسارات QR | بناء URL |

### 3.2 صفحات UI (تستهلك الخدمات أعلاه)

| المسار React | الصفحة | الخدمات |
|--------------|--------|---------|
| `/login` | LoginPage | authService |
| `/signup`, `/signup-complete/:token` | Signup | signupPublicService |
| `/password-recovery`, `/password-reset/:token` | Password | passwordResetPublicService |
| `/twofactor` | TwofactorPage | twoFactorAuthService |
| `/dashboard` | DashboardPage | dashboardService |
| `/devices` | DevicesPage | deviceService, configurationService |
| `/qr/:qrCodeKey` | EnrollmentQrPage | qrImage, public QR |
| `/applications` | ApplicationsPage | applicationService, webUiFilesService |
| `/applications/admin` | AdminApplicationsPage | applicationService (admin) |
| `/application/:id/versions` | ApplicationVersionsPage | applicationService |
| `/configurations`, `.../edit` | Configurations | configurationService |
| `/files` | FilesPage | filesService |
| `/icons` | IconsPage | iconsService |
| `/groups` | GroupsPage | groupService |
| `/users` | UsersPage | userService |
| `/roles` | RolesPage | roleService |
| `/settings` | SettingsPage | settingsService, configurationService |
| `/hints` | HintsPage | hintsService |
| `/updates` | UpdatesPage | updatesService |
| `/push` | PushPage | pushService |
| `/control-panel` | ControlPanelPage | customersService |
| `/plugin-settings` | PluginSettingsPage | pluginService |
| `/profile` | ProfilePage | profileService |
| `/maps` | MapsInfoPage | **لا API** (معلومات فقط) |

### 3.3 ملفات بدون REST مباشر

| نوع | أمثلة |
|-----|--------|
| UI / مكوّنات | `shared/ui/*`, `DeviceForm`, `ConfigurationForm` |
| صلاحيات محلية | [`features/auth/permissions.ts`](frontend/src/features/auth/permissions.ts), `session.ts` |
| ترميز كلمات المرور | `loginPasswordEncode.ts`, `userPasswordEncode.ts` |
| i18n | `i18n/locales/en.json` |
| اختبارات | `*.test.ts`, `*.property.test.tsx` |

---

## 4. جدول المسارات — الفرونت ↔ Go

**رموز:** ✅ متطابق | ⚠️ مسار موجود، حقول/سلوك ناقص | ❌ الفرونت يستدعي وGo لا يوفّر | ⊘ Go يوفّر والفرونت لا يستخدم

### 4.1 عام / مصادقة

| Method | مسار (بعد `/rest`) | الفرونت | Go | الحالة |
|--------|-------------------|---------|-----|--------|
| GET | `/public/auth/options` | ✅ | `auth` | ✅ |
| POST | `/public/auth/login` | ✅ | `auth` | ✅ |
| POST | `/public/auth/logout` | ✅ | `auth` | ✅ |
| POST | `/public/jwt/login` | ⊘ | `auth` | ⊘ |
| GET | `/public/signup/verifyToken/:token` | ✅ | `signup` | ⚠️ يتطلب `MODULE_SIGNUP_ENABLED=true` |
| POST | `/public/signup/verifyEmail` | ✅ | `signup` | ⚠️ نفسه |
| POST | `/public/signup/complete` | ✅ | `signup` | ⚠️ نفسه |
| GET | `/public/signup/canSignup` | ⊘ | `signup` | ⊘ |
| GET | `/public/passwordReset/recover/:user` | ✅ | `passwordreset` | ⚠️ `MODULE_PASSWORDRESET_ENABLED` |
| GET | `/public/passwordReset/settings/:token` | ✅ | `passwordreset` | ✅ |
| POST | `/public/passwordReset/reset` | ✅ | `passwordreset` | ✅ |
| GET | `/public/qr/:id` | ✅ blob | `qrcode` | ✅ |
| GET | `/public/qr/json/:id` | ✅ | `qrcode` | ✅ |
| GET | `/public/name`, `/public/logo` | ⊘ (روابط ثابتة محتملة) | `publicapi` | ⊘ |

### 4.2 مستخدمون، أدوار، ملف شخصي

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET | `/private/users/current` | ✅ | `users` | ✅ |
| PUT | `/private/users/details` | ✅ | `users` | ✅ |
| PUT | `/private/users/current` | ✅ | `users` | ✅ (كلمة مرور) |
| GET | `/private/users/all` | ✅ | `users` | ✅ |
| PUT | `/private/users` أو `/private/users/` | ✅ | `users` | ✅ |
| DELETE | `/private/users/other/:id` | ✅ | `users` | ✅ |
| GET | `/private/users/roles` | ✅ | `users` | ✅ |
| GET | `/private/roles/permissions` | ✅ | `roles` | ✅ |
| GET | `/private/roles/all` | ✅ | `roles` | ✅ |
| PUT | `/private/roles` | ✅ | `roles` PUT `/` | ✅ |
| DELETE | `/private/roles/:id` | ✅ | `roles` | ✅ |
| GET | `/private/twofactor/qr/:userId` | ✅ image | `twofactor` | ✅ |
| GET | `/private/twofactor/verify/:userId/:code` | ✅ | `twofactor` | ✅ |
| GET | `/private/twofactor/set` | ✅ | `twofactor` | ✅ |
| GET | `/private/twofactor/reset` | ✅ | `twofactor` | ✅ |

**حقول `PUT /private/users` (الفرونت → Go):**

| حقل الفرونت | Go / DB | ملاحظة |
|-------------|---------|--------|
| `login`, `name`, `email` | ✅ | |
| `userRole.id` | ✅ | |
| `newPassword` (MD5 hex) | ✅ | عبر `userPasswordEncode` |
| `allDevicesAvailable`, `allConfigAvailable` | ✅ | |
| `groups[]`, `configurations[]` | ✅ | null عند "الكل" |

### 4.3 عملاء (لوحة التحكم)

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| POST | `/private/customers/search` | ✅ | `customers` | ✅ |
| GET | `/private/customers/impersonate/:id` | ✅ | `customers` | ✅ |
| PUT | `/private/customers` | ⊘ | `customers` | ⊘ React لا يوفّر إنشاء/تعديل عميل |
| GET | `/private/customers/:id/edit` | ⊘ | `customers` | ⊘ |
| DELETE | `/private/customers/:id` | ⊘ | `customers` | ⊘ |

### 4.4 أجهزة ومجموعات

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| POST | `/private/devices/search` | ✅ | `devices` | ✅ فلاتر React (012 US1); `installationStatus` عبر infojson فقط |
| GET | `/private/devices/number/:number` | ✅ | `devices` | ✅ يتضمن `info` (012 US1) |
| PUT | `/private/devices` | ✅ | `devices` | ✅ |
| DELETE | `/private/devices/:id` | ✅ | `devices` | ✅ |
| POST | `/private/devices/deleteBulk` | ✅ | `devices` | ✅ |
| POST | `/private/devices/groupBulk` | ✅ | `devices` | ✅ |
| POST | `/private/devices/autocomplete` | ✅ | `devices` | ✅ |
| POST | `/private/devices/:id/description` | ✅ | `devices` | ✅ |
| GET/POST | `.../applicationSettings` | ✅ | `devices` | ✅ |
| POST | `.../applicationSettings/notify` | ✅ | `devices` | ✅ (push queue) |
| GET | `/private/groups/search` | ✅ | `groups` | ✅ |
| PUT | `/private/groups` | ✅ | `groups` | ✅ |
| DELETE | `/private/groups/:id` | ✅ | `groups` | ✅ |
| POST | `/private/groups/autocomplete` | ⊘ | `groups` | ⊘ |

**فلاتر `POST /private/devices/search` — يُطبَّق في Go (012 US1):**

`status`, `androidVersion`, `dateFrom`, `dateTo`, `onlineEarlierMillis`, `onlineLaterMillis`, `enrollmentDateFrom`, `enrollmentDateTo`, `mdmMode`, `kioskMode`, `launcherVersion`, `installationStatus` (infojson), `imeiChanged` (يتطلب عمود `imeiupdatets`), `sortBy`, `sortDir`, بالإضافة إلى `pageNum`, `pageSize`, `value`, `groupId`, `configurationId`, `fastSearch`

**حقول الاستجابة `DeviceView`:**

| حقل UI | Go `DeviceView` | ملاحظة |
|--------|-----------------|--------|
| `id`, `number`, `description`, `lastUpdate`, `imei`, `phone`, `statusCode`, `groups` | ✅ | |
| `configurationId` | ✅ | |
| `model`, `batteryLevel`, `androidVersion`, `serial` | ❌ في القائمة | موجودة داخل `info` في `GET /number/{n}` |
| `info` (كائن `DeviceInfoView`) | ✅ في التفاصيل | `GET /private/devices/number/:number` |
| `applications`, `files` على الجهاز | ✅ داخل `info` | ليس في صفوف `POST /search` |

### 4.5 تكوينات

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET | `/private/configurations/search` | ✅ | `configurations` | ✅ |
| GET | `/private/configurations/search/:value` | ✅ | `configurations` | ✅ |
| GET | `/private/configurations/list` | ✅ | `configurations` | ✅ |
| POST | `/private/configurations/autocomplete` | ✅ body: string | `configurations` | ✅ JSON string |
| GET | `/private/configurations/:id` | ✅ | `configurations` | ✅ |
| PUT | `/private/configurations` | ✅ | `configurations` | ⚠️ حقول تصميم كثيرة؛ Go يحفظ subset |
| PUT | `/private/configurations/copy` | ✅ | `configurations` | ✅ |
| DELETE | `/private/configurations/:id` | ✅ | `configurations` | ✅ |
| GET | `/private/configurations/applications` | ✅ + fallback | `configurations` | ✅ |
| GET | `/private/configurations/applications/:id` | ✅ | `configurations` | ✅ |
| PUT | `/private/configurations/application/upgrade` | ✅ | `configurations` | ✅ |
| POST | `/private/config-files` | ⊘ | `configfiles` | ⊘ الفرونت يعدّل `files[]` داخل PUT configuration فقط |
| GET | `/private/configurations/list` | ✅ (أجهزة) | `configurations` | ✅ |

**ملاحظة:** تبويب ملفات التكوين في UI يعدّل `configuration.files[]` محلياً ويحفظ مع `PUT /configurations` — لا يرفع ملفات عبر `config-files` منفصل.

### 4.6 تطبيقات

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET | `/private/applications/search` | ✅ | `applications` | ✅ |
| GET | `/private/applications/search/:value` | ✅ | `applications` | ✅ |
| POST | `/private/applications/autocomplete` | ✅ | `applications` | ✅ |
| GET | `/private/applications/admin/search` | ✅ | `applications` | ✅ (super admin) |
| GET | `/private/applications/admin/search/:value` | ✅ | `applications` | ✅ |
| GET | `/private/applications/admin/common/:id` | ✅ | `applications` | ✅ |
| GET | `/private/applications/:id` | ✅ | `applications` | ✅ |
| GET | `/private/applications/:id/versions` | ✅ | `applications` | ✅ |
| PUT | `/private/applications/android` | ✅ | `applications` | ✅ |
| PUT | `/private/applications/web` | ✅ | `applications` | ✅ |
| PUT | `/private/applications/versions` | ✅ | `applications` | ✅ |
| PUT | `/private/applications/validatePkg` | ✅ | `applications` | ✅ |
| DELETE | `/private/applications/:id` | ✅ | `applications` | ✅ |
| DELETE | `/private/applications/versions/:id` | ✅ | `applications` | ✅ |
| GET/POST | `/private/applications/configurations...` | ✅ | `applications` | ✅ |
| GET/POST | `/private/applications/version/.../configurations` | ✅ | `applications` | ✅ |

### 4.7 ملفات الواجهة (`web-ui-files`)

| Method | مسار | الفرونت | Go (`files` module) | الحالة |
|--------|------|---------|---------------------|--------|
| GET | `/private/web-ui-files/search` | ✅ | ✅ Group `/web-ui-files` | ✅ |
| POST | `/private/web-ui-files` | ✅ multipart | ✅ Upload | ✅ |
| POST | `/private/web-ui-files/raw` | ✅ | ✅ | ✅ |
| POST | `/private/web-ui-files/update` | ✅ | ✅ | ✅ |
| POST | `/private/web-ui-files/remove` | ✅ | ✅ | ✅ |
| GET | `/private/web-ui-files/limit` | ✅ | ✅ | ⚠️ quota قد تكون مبسّطة |
| GET | `/private/web-ui-files/apps/:url` | ✅ | ✅ | ✅ |
| GET | `/private/web-ui-files/configurations/:id` | ⊘ | ✅ | ⊘ |
| POST | `/private/web-ui-files/configurations` | ⊘ | ✅ | ⊘ |

### 4.8 أيقونات

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET | `/private/icons/search` | ✅ | `icons` | ✅ |
| PUT | `/private/icons` | ✅ `{ name, fileId }` | `icons` | ✅ |
| DELETE | `/private/icons/:id` | ✅ | `icons` | ✅ |
| POST | `/private/icon-files` | ⊘ | `icons` | ⊘ UI يطلب **fileId يدوياً** بدل رفع صورة |

### 4.9 إعدادات، تلميحات، ملخص

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET | `/private/settings` | ✅ | `settings` | ✅ |
| POST | `/private/settings/misc` | ✅ | `settings` | ✅ |
| POST | `/private/settings/lang` | ✅ | `settings` | ✅ |
| POST | `/private/settings/design` | ✅ | `settings` | ✅ |
| GET | `/private/settings/userRole/:roleId` | ✅ | `settings` | ✅ |
| POST | `/private/settings/userRoles/common` | ✅ | `settings` | ✅ |
| GET | `/private/hints/history` | ✅ | `hints` | ✅ |
| POST | `/private/hints/enable` | ✅ | `hints` | ✅ |
| POST | `/private/hints/disable` | ✅ | `hints` | ✅ |
| POST | `/private/hints/history` | ⊘ | `hints` MarkShown | ⊘ |
| GET | `/private/summary/devices` | ✅ | `summary` | ⚠️ شكل OK؛ بيانات charts قد تكون فارغة/مبسّطة |

**`DeviceSummaryPayload`:** الفرونت يتوقع `statusSummary`, `installSummary`, `devicesTotal`, … — Go يُرجع نفس أسماء JSON؛ قيمة `stringAttr` للألوان (`red`/`yellow`/`green`) متوافقة مع [`dashboardService.ts`](frontend/src/features/dashboard/dashboardService.ts).

### 4.10 Push، تحديثات، Plugins

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| POST | `/private/push` | ✅ | `push` | ✅ |
| GET | `/private/update/check` | ✅ | `updates` | ✅ |
| POST | `/private/update` | ⊘ | `updates` Apply | ⊘ صفحة Updates تعرض فقط |
| GET | `/plugin/main/private/active` | ✅ | `plugins/platform` | ✅ |
| GET | `/plugin/main/private/available` | ✅ | `plugins/platform` | ✅ |
| POST | `/plugin/main/private/disabled` | ✅ | `plugins/platform` | ✅ |
| POST | `/rest/plugins/push/private/search` | ⊘ | `plugins/push` | ⊘ لا صفحة React |
| POST | `/rest/plugins/messaging/...` | ⊘ | `plugins/messaging` | ⊘ |
| POST | `/rest/plugins/audit/...` | ⊘ | `plugins/audit` | ⊘ |
| `/rest/plugins/deviceinfo/...` | ⊘ | `plugins/deviceinfo` | ⊘ |
| `/rest/plugins/devicelog/...` | ⊘ | `plugins/devicelog` | ⊘ |

**Push body (متطابق):**

```json
{
  "messageType": "string",
  "payload": "string",
  "deviceNumbers": ["..."],
  "groups": ["..."],
  "broadcast": false
}
```

---

## 5. فجوات الحقول حسب الميزة

### 5.1 تسجيل الدخول

| الحقل | الفرونت يتوقع | Go يرسل |
|-------|---------------|---------|
| `authToken` | مطلوب | ✅ من login |
| `userRole.permissions[]` | لـ `permissions.ts` | ✅ إن وُجدت في DB |
| `superAdmin`, `singleCustomer` | ✅ | ✅ |
| `twoFactor`, `twoFactorAccepted` | توجيه `/twofactor` | ✅ |
| `passwordReset`, `passwordResetToken` | توجيه reset | ✅ |

### 5.2 مستخدمون

| حقل | ملاحظة |
|-----|--------|
| `newPassword` | MD5 uppercase hex — يجب أن يطابق ما يتوقعه Go auth |
| `groups` / `configurations` | مصفوفة `{ id }` — متوافق |

### 5.3 تكوين (Configuration)

الفرونت يرسل كائناً كبيراً (`Configuration` في [`types.ts`](frontend/src/features/configurations/types.ts)). Go [`domain/configuration.go`](serverBackendGo/internal/modules/configurations/domain/configuration.go) يدعم جزءاً من الحقول؛ الباقي قد يُهمل أو يُخزَّن في `Extra` إن وُجد منطق merge — راجع parity عند مشاكل حفظ.

حقول حرجة للـ QR:

| حقل | استخدام UI |
|-----|------------|
| `qrCodeKey` | Enrollment QR |
| `baseUrl` | QR URL |
| `mainAppId`, `contentAppId` | تطبيقات التكوين |

### 5.4 ملفات (رفع)

| العنصر | الفرونت | Go |
|--------|---------|-----|
| حقل multipart | `file` (typical) | تحقق من handler `files` |
| استجابة الرفع | `fileId`, `url` | يجب أن تطابق `FileUploadResult` في webUiFilesService |

---

## 6. مسارات Go غير مستخدمة من React (لكن موجودة)

مفيدة للوكلاء أو توسعة UI لاحقاً:

| المسار | الوحدة | الغرض |
|--------|--------|--------|
| `/rest/public/sync/*` | sync | وكلاء Android |
| `/rest/notifications/*`, `/rest/notification/polling/*` | notifications | وكلاء |
| `/rest/public/update` (agent) | — | مختلف عن `/private/update` |
| `/rest/private/config-files` POST | configfiles | رفع ملف تكوين |
| `/rest/private/icon-files` POST | icons | رفع أيقونة |
| `/rest/plugins/*` (push, messaging, audit, deviceinfo, devicelog) | plugins | لوحة Angular قديمة / مستقبل |
| `/rest/public/applications/upload` | publicapi | رفع APK عام |

---

## 7. أعلام التفعيل في الباكند (`.env`)

إذا كان المسار 404 أو الوحدة لا تُسجَّل:

| المتغير | يؤثر على |
|---------|----------|
| `MODULE_AUTH_ENABLED` | auth |
| `MODULE_SIGNUP_ENABLED` | signup — **افتراضي false** → صفحة التسجيل تفشل |
| `MODULE_PASSWORDRESET_ENABLED` | password reset |
| `MODULE_*_ENABLED` لكل phase | انظر [`serverBackendGo/.env.example`](serverBackendGo/.env.example) |

**محلي:** طابق [`serverBackendGo/.env`](serverBackendGo/.env) مع `.env.example`؛ `FILES_DIRECTORY=./data/files`.

---

## 8. قائمة تحقق تكامل (UAT)

```text
[ ] Vite proxy → :8080 و make dev يعمل
[ ] Login → dashboard (session + Bearer)
[ ] Devices: بحث نصي + pagination + group/config filters
[ ] Device detail: قبول أن info/apps قد تكون N/A
[ ] Configurations: CRUD + QR enrollment
[ ] Applications + versions + file upload
[ ] Files page: search / remove
[ ] Icons: metadata (fileId يدوي أو ربط برفع لاحق)
[ ] Users / Roles / Settings / Hints
[ ] Push message
[ ] Control panel impersonate (super admin)
[ ] Plugin settings toggles
[ ] 2FA flow إن مُفعّل للمستخدم
[ ] Signup / password reset إن MODULE_*_ENABLED=true
```

---

## 9. توصيات لضمان التزامن

| الأولوية | الإجراء |
|----------|---------|
| **P0** | توسيع `devices.SearchRequest` + SQL في Go ليطابق فلاتر [`FilterPanel`](frontend/src/features/devices/FilterPanel.tsx) |
| **P0** | إثراء `GetByNumber` بـ `info` أو حقول مسطحة يتوقعها `DeviceDetailPanel` |
| **P1** | صفحة Icons: استدعاء `POST /private/icon-files` بدل إدخال fileId يدوياً |
| **P1** | تفعيل `MODULE_SIGNUP_ENABLED` في dev إن اختبرت التسجيل |
| **P2** | صفحات React لـ plugin push/messaging/audit (اختياري) |
| **P2** | إزالة أو تعطيل `websocket.ts` أو توثيق أنه غير مدعوم مع Go |
| **P3** | `POST /private/update` في UpdatesPage إن أردت تطبيق التحديثات من UI |

---

## 10. مراجع

| الوثيقة | الغرض |
|---------|--------|
| [`JAVA-GO-MIGRATION-STATUS.md`](JAVA-GO-MIGRATION-STATUS.md) | حالة نقل Java → Go |
| [`serverBackendGo/docs/parity/`](serverBackendGo/docs/parity/) | جداول endpoint لكل وحدة |
| [`serverBackendGo/docs/MIGRATION.md`](serverBackendGo/docs/MIGRATION.md) | مراحل الباكند |
| [`frontend/vite.config.ts`](frontend/vite.config.ts) | بروكسي API |

---

*عند تغيير أي `*Service.ts` في الفرونت أو handlers في Go، حدّث قسم §4 و§5 في هذا الملف.*
