# تحليل تكامل الفرونت (React) مع الباكند (Go)

**تاريخ التحليل:** 2026-05-21 (محدّث بعد 013 + اختبار API)  
**الفرونت:** [`frontend/`](frontend/) — Vite + React + Axios  
**الباكند:** [`serverBackendGo/`](serverBackendGo/) — Gin، مسارات `/rest/*`  
**مراجع:** [`JAVA-GO-MIGRATION-STATUS.md`](JAVA-GO-MIGRATION-STATUS.md) · [`JAVA-GO-BACKEND-GAPS.md`](JAVA-GO-BACKEND-GAPS.md) · [`JAVA-GO-DATABASE-GAPS.md`](JAVA-GO-DATABASE-GAPS.md)

---

## 1. الخلاصة التنفيذية

| البُعد | الحكم | ملاحظة |
|--------|--------|--------|
| **توافق المسارات (URL)** | **~96%** | كل استدعاء REST في React له مسار في Go؛ الاستثناءات: WebSocket، `stats`، `videos` |
| **توافق الحقول (payload/response)** | **~82–88%** | تحسّن بعد 013 (`userrolesettings`, `devicestatuses`)؛ تكوينات MDM policy في `settingsjson` |
| **توافق السلوك** | **~85%** | يعتمد على `.env`؛ بعض السلوكيات Java (bootstrap عميل، MQTT) غير منقولة |
| **جاهزية تشغيل UI الرئيسي** | **نعم** | Login، أجهزة، تكوينات، تطبيقات، إعدادات، لوحة تحكم — مع الفجوات أدناه |

**قاعدة التشغيل:** الفرونت يستخدم `baseURL: '/rest'`؛ Vite يوجّه إلى `http://localhost:8080` ([`frontend/vite.config.ts`](frontend/vite.config.ts)).

**قاعدة البيانات:** migrations حتى **`000017`** (013) مطلوبة لأعمدة الجدول، فلاتر `installationStatus`، و`userRole` columns — راجع [`serverBackendGo/docs/MIGRATION.md`](serverBackendGo/docs/MIGRATION.md).

---

## 2. البنية المشتركة

### 2.1 عميل HTTP

| العنصر | الفرونت | Go |
|--------|---------|-----|
| الملف | [`apiClient.ts`](frontend/src/services/apiClient.ts) | Router: `/rest/public`, `/rest/private`, `/rest/plugin/*`, `/rest/plugins/*`, `/rest/notifications` |
| Cookies | `withCredentials: true` | Session على `POST /public/auth/login` |
| JWT | `Authorization: Bearer` من `localStorage` | Middleware على `/rest/private/*` |
| Envelope | [`hmdmEnvelope.ts`](frontend/src/services/hmdmEnvelope.ts) | `status` OK/ERROR + `data` |

### 2.2 المصادقة

| الخطوة | مسار الفرونت | Go | الحالة |
|--------|--------------|-----|--------|
| خيارات | `GET /public/auth/options` | `auth` | ✅ |
| دخول | `POST /public/auth/login` | `auth` | ✅ (MD5 أو نص حسب `TRANSMIT_PASSWORD`) |
| خروج | `POST /public/auth/logout` | `auth` | ✅ |
| JWT | ⊘ (React لا يستدعيه) | `POST /public/jwt/login` | ⊘ — مفيد للاختبارات اليدوية |
| مستخدم حالي | `GET /private/users/current` | `users` | ✅ |
| 2FA | `/private/twofactor/*` | `twofactor` | ✅ |

### 2.3 WebSocket — غير مدعوم

| الفرونت | Go |
|---------|-----|
| [`websocket.ts`](frontend/src/services/websocket.ts) — غير مستورد في الصفحات | ❌ لا `/rest/ws/connect` |

### 2.4 وحدات Go المسجّلة (مرجع سريع)

من [`internal/app/modules.go`](serverBackendGo/internal/app/modules.go): `auth`, `signup`, `passwordreset`, `users`, `twofactor`, `roles`, `customers`, `settings`, `hints`, `summary`, `devices`, `groups`, `applications`, `configurations`, `configfiles`, `files`, `icons`, `publicapi`, `sync`, `push`, `notifications`, `updates`, `qrcode`, `plugins/platform`, `plugins/audit`, `plugins/push`, `plugins/messaging`, `plugins/deviceinfo`, `plugins/devicelog`.

**غير مسجّلة:** `stats`, `videos` (`MODULE_STATS_ENABLED=false`, `MODULE_VIDEOS_ENABLED=false` في [`.env.example`](serverBackendGo/.env.example)).

---

## 3. خريطة ملفات الفرونت (REST)

| الملف | المجال |
|-------|--------|
| [`authService.ts`](frontend/src/services/authService.ts) | auth |
| [`deviceService.ts`](frontend/src/features/devices/deviceService.ts) | devices + groups/config list |
| [`configurationService.ts`](frontend/src/features/configurations/configurationService.ts) | configurations |
| [`applicationService.ts`](frontend/src/features/applications/services/applicationService.ts) | applications |
| [`webUiFilesService.ts`](frontend/src/features/applications/services/webUiFilesService.ts) | ملفات واجهة |
| [`settingsService.ts`](frontend/src/features/settings/settingsService.ts) | settings + user role columns |
| [`dashboardService.ts`](frontend/src/features/dashboard/dashboardService.ts) | summary + counts |
| [`pluginService.ts`](frontend/src/features/plugins/pluginService.ts) | plugin platform فقط |
| باقي `*Service.ts` | users, roles, customers, hints, push, updates, icons, files, signup, password reset, 2FA |

**صفحات بدون REST:** [`MapsInfoPage`](frontend/src/features/maps/MapsInfoPage.tsx) (معلومات فقط).

---

## 4. جدول المسارات — الفرونت ↔ Go

**رموز:** ✅ متطابق | ⚠️ مسار موجود — حقول/سلوك ناقص | ❌ الفرونت يحتاج وGo لا يوفّر | ⊘ Go يوفّر والفرونت لا يستخدم

### 4.1 عام / مصادقة

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET | `/public/auth/options` | ✅ | `auth` | ✅ |
| POST | `/public/auth/login` | ✅ | `auth` | ✅ |
| POST | `/public/auth/logout` | ✅ | `auth` | ✅ |
| POST | `/public/jwt/login` | ⊘ | `auth` | ⊘ |
| GET/POST | `/public/signup/*` | ✅ | `signup` | ⚠️ `MODULE_SIGNUP_ENABLED=false` افتراضياً |
| GET/POST | `/public/passwordReset/*` | ✅ | `passwordreset` | ⚠️ `MODULE_PASSWORDRESET_ENABLED` |
| GET | `/public/qr/:id`, `/public/qr/json/:id` | ✅ | `qrcode` | ✅ |
| GET | `/public/name`, `/public/logo` | ⊘ | `publicapi` | ⊘ |
| PUT | `/public/stats` | ⊘ | — | ❌ لا module (جدول `usagestats` ✅ بعد 013) |
| GET/POST | `/public/videos/*` | ⊘ | — | ❌ `MODULE_VIDEOS_ENABLED=false` |

### 4.2 مستخدمون، أدوار، ملف شخصي

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET | `/private/users/current` | ✅ | `users` | ✅ |
| PUT | `/private/users/details` | ✅ | `users` | ✅ |
| PUT | `/private/users/current` | ✅ | `users` | ✅ |
| GET/PUT/DELETE | `/private/users/*` | ✅ | `users` | ✅ |
| GET | `/private/roles/permissions`, `/all` | ✅ | `roles` | ✅ |
| PUT | `/private/roles` | ✅ | `roles` | ✅ |
| GET | `/private/twofactor/*` | ✅ | `twofactor` | ✅ |

### 4.3 عملاء (لوحة التحكم)

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| POST | `/private/customers/search` | ✅ | `customers` | ✅ |
| GET | `/private/customers/impersonate/:id` | ✅ | `customers` | ✅ |
| PUT/GET/DELETE | `/private/customers/*` | ⊘ | `customers` | ⊘ React لا يوفّر CRUD عميل |

### 4.4 أجهزة ومجموعات

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| POST | `/private/devices/search` | ✅ | `devices` | ✅ فلاتر React (012+013) |
| GET | `/private/devices/number/:number` | ✅ | `devices` | ✅ + `info` |
| PUT/DELETE/POST bulk… | ✅ | `devices` | ✅ |
| GET | `/private/groups/search` | ✅ | `groups` | ✅ |
| PUT/DELETE | `/private/groups` | ✅ | `groups` | ✅ |
| POST | `/private/groups/autocomplete` | ⊘ | `groups` | ⊘ |
| GET | `/private/groups/search/:value` | ⊘ | `groups` | ⊘ |

**فلاتر `POST /private/devices/search` — مطبّقة في Go:**

| فلتر React | Go | ملاحظة |
|------------|-----|--------|
| `status`, `androidVersion`, `mdmMode`, `kioskMode` | ✅ | |
| `dateFrom`/`dateTo`, `enrollmentDateFrom`/`To` | ✅ | |
| `onlineEarlierMillis`/`onlineLaterMillis` | ✅ | |
| `launcherVersion` | ✅ | عبر `infojson` |
| **`installationStatus`** | ✅ | **013:** `devicestatuses.applicationsstatus` (ليس `infojson` فقط) |
| `imeiChanged` | ✅ | عمود `imeiupdatets` |
| `sortBy` / `sortDir` | ✅ | يشمل `INSTALLATIONS`, `FILES` عبر `devicestatuses` |

**حقول قائمة الأجهزة:**

| حقل UI | في `POST /search` | في `GET /number/:n` |
|--------|-------------------|---------------------|
| `id`, `number`, `statusCode`, `groups`, `configurationId` | ✅ | ✅ |
| `model`, `batteryLevel`, `androidVersion`, … | ❌ مسطّحة | ✅ داخل `info` |
| `applications` / `files` على الجهاز | ❌ | ✅ داخل `info` |

**سلوك غير منقول (لا يؤثر على مسار REST):** تحديث `devicestatuses` تلقائياً من sync/agent — الجدول يُملأ افتراضياً عند الهجرة؛ القيم الحقيقية تحتاج منطق أعمال لاحق.

### 4.5 تكوينات

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET/PUT/DELETE | `/private/configurations/*` | ✅ | `configurations` | ✅ |
| POST | `/private/configurations/autocomplete` | ✅ | `configurations` | ✅ |
| PUT | `/private/configurations/application/upgrade` | ✅ | `configurations` | ✅ |
| POST | `/private/config-files` | ⊘ | `configfiles` | ⊘ UI يحفظ `files[]` داخل `PUT /configurations` |

**حقول محرر التكوين (Configuration):**

| نوع الحقل | Go | ملاحظة للفرونت |
|-----------|-----|----------------|
| `name`, `description`, `type`, `password`, ألوان، `qrCodeKey`, `baseUrl`, `mainAppId`, `contentAppId` | أعمدة SQL | ✅ |
| سياسات MDM (`gps`, `kioskMode`, `wifi`, …) | **`settingsjson` JSONB** | ✅ عند `GET/:id` يُدمج JSON في الجسم؛ عند `PUT` يُحفظ ما يُرسل في الحقول المدمجة |
| `skipVersionCheck` لكل تطبيق | **`configurationapplicationparameters`** (013) | ✅ عند الحفظ إن وُجد الحقل في عنصر `applications[]` |
| `remove`, `longTap` على ربط التطبيق | أعمدة `configurationapplications` (013) | ✅ في INSERT؛ قد لا تُعاد كلها في `GET applications/:id` — تحقق عند UI |
| ملفات التكوين | `configurationfiles` | ✅ |

### 4.6 تطبيقات

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET/PUT/POST/DELETE | `/private/applications/*` | ✅ | `applications` | ✅ |
| `apkhash` على الإصدار | ⚠️ | `applications` | ⚠️ عمود DB موجود (013)；واجهة React قد لا ترسل `apkHash` بعد |

### 4.7 ملفات الواجهة (`web-ui-files`)

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET/POST | `/private/web-ui-files/*` | ✅ | `files` | ✅ |
| GET | `/private/web-ui-files/configurations/:id` | ⊘ | `files` | ⊘ |

### 4.8 أيقونات

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET/PUT/DELETE | `/private/icons/*` | ✅ | `icons` | ✅ |
| POST | `/private/icon-files` | ⊘ | `icons` | ⊘ [`IconsPage`](frontend/src/features/icons/IconsPage.tsx) يدخل **fileId يدوياً** |

### 4.9 إعدادات، تلميحات، ملخص

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| GET/POST | `/private/settings`, `misc`, `lang`, `design` | ✅ | `settings` | ✅ |
| GET | `/private/settings/userRole/:roleId` | ✅ | `settings` | ✅ **013:** كل `columnDisplayed*` من `userrolesettings` |
| POST | `/private/settings/userRoles/common` | ✅ مصفوفة | `settings` | ✅ **مصفوفة** `UserRoleSettings[]` (مطابق Java) |
| GET/POST | `/private/hints/*` | ✅ | `hints` | ✅ |
| POST | `/private/hints/history` (mark shown) | ⊘ | `hints` | ⊘ |
| GET | `/private/summary/devices` | ✅ | `summary` | ✅ **013:** `installSummary` + `app*ByConfig` من `devicestatuses` |

**حقول `GET /private/settings` — الفرونت لا يعرضها كلها:**

| حقل DB (013 `000015`) | في `normalizeSettings()` | تأثير UI |
|------------------------|--------------------------|----------|
| `newdevicegroupid`, `phonenumberformat` | ❌ غير مُعرَّف | ⚠️ لا تظهر في صفحة Settings |
| `custompropertyname1`–`3`, `custommultiline*`, `customsend*` | ❌ | ⚠️ تبويب أعمدة الأجهزة يقرأ الأسماء من `fetchRawSettings` جزئياً |
| `desktopheadertemplate`, `senddescription` | ❌ | ⚠️ |

### 4.10 Push، تحديثات، Plugins

| Method | مسار | الفرونت | Go | الحالة |
|--------|------|---------|-----|--------|
| POST | `/private/push` | ✅ | `push` | ✅ |
| GET | `/private/update/check` | ✅ | `updates` | ✅ |
| POST | `/private/update` | ⊘ | `updates` | ⊘ |
| GET/POST | `/plugin/main/private/*` | ✅ | `plugins/platform` | ✅ |
| `/rest/plugins/push|messaging|audit|deviceinfo|devicelog/*` | ⊘ | plugins | ⊘ لا صفحات React |

---

## 5. فجوات REST — الفرونت يحتاجها وGo لا يوفّرها بعد

| الأولوية | المسار / الميزة | تأثير على React |
|----------|-----------------|-----------------|
| **P1** | `PUT /rest/public/stats` | لا يؤثر على UI؛ أدوات/ترخيص خارج الواجهة |
| **P2** | `GET/POST /rest/public/videos/{fileName}` | روابط تدريب قديمة فقط |
| **P2** | `POST /private/icon-files` + ربط Icons UI | رفع أيقونة بدل fileId يدوي |
| **P3** | صفحات plugin push/messaging/audit/deviceinfo/devicelog | لوحة Angular القديمة |
| **P3** | WebSocket `/rest/ws/connect` | غير مستخدم في React |

مرجع تفصيلي: [`JAVA-GO-BACKEND-GAPS.md`](JAVA-GO-BACKEND-GAPS.md) §4.

---

## 6. فجوات حقول/سلوك — مسار موجود لكن التجربة ناقصة

### 6.1 أجهزة

| الفجوة | التفاصيل | إجراء مقترح |
|--------|----------|-------------|
| أعمدة الجدول في البحث | `model`, `batteryLevel`… فقط في `info` عند التفاصيل | اختياري: إسقاط حقول من `infojson` في `DeviceView` للقائمة |
| `launcherVersion` في الصف | غير عمود في الاستجابة | فلتر يعمل؛ عمود UI قد يظهر فارغاً |
| تحديث `devicestatuses` | لا يُحدَّث من `sync/info` تلقائياً | منطق أعمال في `sync` أو خدمة devices |

### 6.2 تكوينات

| الفجوة | التفاصيل |
|--------|----------|
| سياسات MDM كثيرة | تُخزَّن في `settingsjson` — يجب أن يرسل الفرونت نفس مفاتيح Java (camelCase) عند الحفظ |
| حقول تصميم على مستوى configuration | `iconSize`, `desktopHeader` داخل JSON وليس أعمدة منفصلة |
| قراءة `skipVersionCheck` | قد لا تُحمَّل في `GET applications/:id` — تحقق عند تبويب التطبيقات |

### 6.3 إعدادات

| الفجوة | التفاصيل |
|--------|----------|
| أعمدة tenant الجديدة (013) | DB جاهز؛ `GET /settings` لا يملأ كل الحقول في [`normalizeSettings`](frontend/src/features/settings/settingsService.ts) |
| أسماء `custom1`–`custom3` في تبويب الأعمدة | يعتمد `fetchRawSettings()` — يعمل إن رجعت من DB |

### 6.4 عملاء

| الفجوة | التفاصيل |
|--------|----------|
| Bootstrap عند إنشاء عميل | Java ينسخ قوالب؛ Go `PUT /customers` بدون نسخ كامل |
| Mailchimp | غير منقول |

### 6.5 لوحة التحكم (Summary)

| العنصر | الحالة بعد 013 |
|--------|----------------|
| `statusSummary`, `devicesTotal`, `installSummary` | ✅ |
| `appSuccessByConfig`, `appFailureByConfig`, … | ✅ (حتى 10 تكوينات) |
| `devicesEnrolledMonthly` | ⚠️ سلسلة شهرية مبسّطة/أصفار |

### 6.6 عملاء Android (وكلاء)

| المسار | Go | الفرونت |
|--------|-----|---------|
| `/rest/public/sync/*` | ✅ | ⊘ |
| `/rest/notifications/*` | ✅ | ⊘ |

---

## 7. ما أُضيف في schema (013) وتأثيره على الفرونت

| Migration | الجدول/العمود | تأثير مباشر على React |
|-----------|---------------|------------------------|
| `000011` | `devicestatuses` | فلتر **App install status** + ترتيب INSTALLATIONS/FILES + charts لوحة التحكم |
| `000012` | `userrolesettings` | تبويب **Settings → Role columns** |
| `000013` | `configurationapplicationparameters` | `skipVersionCheck` في محرر التكوين (عند الإرسال) |
| `000014` | `usagestats` | لا UI؛ يحتاج module `stats` (012) |
| `000015` | أعمدة `settings` | حقول tenant إضافية — **واجهة React لم تُحدَّث بعد** |
| `000016` | `apkhash`, `remove`, `longtap` | تطبيقات/تكوين — جزئي في API |
| `000017` | استيراد Java → `settingsjson` | عند استعادة dump Java فقط |

---

## 8. مسارات Go غير مستخدمة من React (للتوسعة)

| المسار | الغرض |
|--------|--------|
| `/rest/public/sync/*` | وكلاء |
| `/rest/notifications/*` | polling وكلاء |
| `/rest/private/config-files` POST | رفع ملف تكوين منفصل |
| `/rest/private/icon-files` POST | رفع أيقونة |
| `/rest/public/applications/upload` | رفع APK عام |
| `/rest/plugins/*` | إدارة plugins متقدمة |
| `/rest/private/web-ui-files/configurations` | ربط ملف↔تكوينات |

---

## 9. أعلام `.env` المؤثرة على الفرونت

| المتغير | التأثير |
|---------|---------|
| `MODULE_SIGNUP_ENABLED` | صفحة `/signup` — **false** افتراضياً |
| `MODULE_PASSWORDRESET_ENABLED` | استعادة كلمة المرور |
| `MODULE_*` لكل phase | إيقاف تسجيل routes → 404 |
| `FILES_DIRECTORY` | رفع ملفات/QR — استخدم `./data/files` محلياً |
| `MODULE_STATS_ENABLED` | **true** في dev لاختبار `PUT /public/stats` (014) |
| `MODULE_VIDEOS_ENABLED` | **false** — لا REST فيديو |

---

## 10. قائمة تحقق UAT (محدّثة)

```text
[x] Vite proxy → :8080 و make dev
[x] Login → dashboard
[x] Devices: فلاتر متقدمة + installationStatus (devicestatuses)
[x] Device detail: GET /number/{n} + info
[x] Settings: userRole columns GET + POST array common
[x] Summary: installSummary + per-config app status
[x] Configurations: settingsjson + skipVersionCheck/remove/longTap (014 — اختبار يدوي لكل تبويب موصى به)
[x] Icons: رفع عبر icon-files (014)
[x] Settings: حقول tenant 000015 (014)
[x] Stats: PUT /public/stats عند MODULE_STATS_ENABLED=true
[x] Sync: devicestatuses من sync/info (014 مبسّط)
[x] Updates: POST /private/update من UpdatesPage
[x] Hints: POST /private/hints/history (mark shown)
[ ] Signup / password reset مع MODULE_*=true
[ ] Control panel impersonate (super admin)
[ ] Plugin sub-pages (push/messaging…) — غير مطلوب لـ React الحالي
```

---

## 11. توصيات أولويات (مزامنة الفرونت ↔ Go)

| الأولوية | الإجراء | الملفات |
|----------|---------|---------|
| **P0** | ~~`installationStatus` + userRole columns~~ | **منجز 013** |
| ~~**P1**~~ | ~~tenant settings + config MDM + icons~~ | **منجز 014** |
| ~~**P2**~~ | ~~stats + sync devicestatuses~~ | **منجز 014** |
| **P2** | أعمدة `model`/`batteryLevel` من `infojson` (اختياري) | `device_repo.go` |
| **P3** | صفحات React للـ plugins أو الإبقاء على ⊘ | — |
| ~~**P3**~~ | ~~Updates + hints history~~ | **منجز 014** |

---

## 12. مراجع وصيانة الملف

| الوثيقة | الغرض |
|---------|--------|
| [`serverBackendGo/docs/parity/`](serverBackendGo/docs/parity/) | endpoint لكل module |
| [`JAVA-GO-DATABASE-GAPS.md`](JAVA-GO-DATABASE-GAPS.md) | schema |
| [`JAVA-GO-BACKEND-GAPS.md`](JAVA-GO-BACKEND-GAPS.md) | REST/سلوك ناقص |
| [`specs/013-complete-database-gaps/quickstart.md`](specs/013-complete-database-gaps/quickstart.md) | smoke tests |

**قواعد التحديث:** عند تغيير `*Service.ts` أو handler Go — حدّث §4 و§6. عند migration جديد — أضف صفاً في §7.

---

*آخر تحديث: 2026-05-21 — بعد migrations `000011`–`000017` واختبار API (installationStatus، userRole، summary).*
