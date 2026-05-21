# تحليل الفجوات: Java (`backend/`) → Go (`serverBackendGo/`)

**تاريخ التحليل:** 2026-05-21  
**الغرض:** مرجع واحد لمعرفة ما نُقل بنجاح، ما نُقل جزئياً، وما لم يُنفَّذ بعد — لضمان اكتمال الهجرة دون كسر `frontend/` أو وكلاء Android.

**مصادر:**

| المصدر | المسار |
|--------|--------|
| Java (مرجع الحقيقة) | [`backend/`](backend/) — JAX-RS + Guice + MyBatis |
| Go (الهدف) | [`serverBackendGo/`](serverBackendGo/) — Gin + Postgres |
| خارطة المراحل | [`serverBackendGo/docs/MIGRATION.md`](serverBackendGo/docs/MIGRATION.md) |
| جداول التطابق لكل وحدة | [`serverBackendGo/docs/parity/`](serverBackendGo/docs/parity/) |
| تجربة Go قديمة (لا تُنسخ أعمى) | [`backend-go/`](backend-go/) |

---

## 1. الملخص التنفيذي

| المؤشر | القيمة التقريبية |
|--------|------------------|
| مراحل الهجرة المعلّمة **done** في `MIGRATION.md` | **8 / 8** (Phases 1–8) |
| موارد REST Java الرئيسية (`*Resource.java`) | **~35** ملفاً |
| وحدات Go مسجّلة (`internal/modules/*/module.go`) | **29** وحدة |
| حالة **منقول بالكامل** (مسارات + سلوك أساسي) | **~22** مجالاً |
| **جزئي** (مسار موجود، سلوك/تكامل ناقص) | **~15** مجالاً |
| **غير منقول** (لا يوجد مكافئ Go) | **4** REST + مكونات خلفية |
| إضافة **xtra** (Angular فقط، بدون REST) | خارج نطاق Go الحالي |

> **تنبيه:** وسم المرحلة **done** في `MIGRATION.md` يعني «المسارات الأساسية موجودة وتعمل في smoke»، وليس «تطابق 100% مع Java». هذا الملف يوضح الفجوات التي تبقى بعد إغلاق المراحل.

### نسبة التغطية التقريبية (حسب الوظيفة، وليس عدد الأسطر)

```
REST عام (لوحة التحكم + MDM core)     ████████████████████░░  ~85%
وكلاء (sync / notifications / QR)      █████████████████░░░░░  ~80%
إضافات (plugins)                       ███████████████░░░░░░░  ~75%
تكاملات خلفية (FCM, cron, audit auto)  ██████░░░░░░░░░░░░░░░░  ~30%
```

---

## 2. منهجية التحليل

1. جرد كل `@Path` في `backend/**/*Resource.java` ووحدة `backend/jwt` و`backend/notification`.
2. مقارنة كل مجال مع وحدة Go المقابلة وملف `docs/parity/<name>.md`.
3. فحص `internal/modules/*/module.go` للتأكد من التسجيل الفعلي.
4. تمييز: **Done** | **Partial** | **Missing** | **Out of scope** (مُعلَن في parity أو غير مستخدم من React).

**رموز الحالة:**

| الرمز | المعنى |
|-------|--------|
| ✅ | منقول — السلوك الأساسي متوافق مع React/الوكلاء |
| ⚠️ | جزئي — مسار موجود لكن سلوك/تكامل ناقص |
| ❌ | غير منقول — لا handler في Go |
| ⊘ | خارج النطاق — مُهمل في Java أو واجهة Angular فقط |

---

## 3. خريطة المراحل (Phases 1–8)

| Phase | الوحدات Go | Java (إرشادي) | حالة الخارطة | فجوات جوهرية متبقية |
|-------|------------|---------------|--------------|---------------------|
| **1** | `auth` | `AuthResource`, `JWTAuthResource` | done | — |
| **1b** | `signup`, `passwordreset`, `twofactor` | `SignupResource`, `PasswordResetResource`, 2FA | done | Mailchimp، حدود الأجهزة عند التسجيل |
| **2** | `users`, `roles` | `UserResource`, `UserRoleResource` | done | impersonate / superadmin |
| **3** | `customers`, `settings`, `hints`, `summary` | `CustomerResource`, `SettingsResource`, `HintResource`, `SummaryResource` | done | إنشاء عميل بدون نسخ افتراضي؛ رسوم بيانية مبسّطة |
| **4** | `devices`, `groups`, `configurations` (list) | `DeviceResource`, `GroupResource`, `ConfigurationResource` | done | بحث متقدم، إثراء نتائج، FCM |
| **5** | `applications`, `configurations`, `configfiles` | `ApplicationResource`, `ConfigurationResource`, `ConfigurationFileResource` | done | push عند الحفظ؛ quota ملفات |
| **6** | `files`, `icons`, `publicapi` | `FilesResource`, `IconResource`, `PublicResource` | done | تحميل ثابت للوكلاء؛ **IconFileResource** ❌ |
| **7** | `sync`, `push`, `notifications`, `updates`, `qrcode` | `SyncResource`, `PushApiResource`, `NotificationResource`, `UpdateResource`, `QRCodeResource` | done | hooks، FCM حقيقي، APK بعيد، stats |
| **8** | `plugins/*` | `backend/plugins/*` | done | audit filter، cron جدولة push، deviceinfo/devicelog export |

---

## 4. جدول الوحدات — Java REST ↔ Go

### 4.1 Core server (`backend/server/.../resource/`)

| Java Resource | Prefix REST | Go module | الحالة | مرجع parity |
|---------------|-------------|-----------|--------|-------------|
| `AuthResource` | `/rest/public/auth` | `auth` | ✅ | [auth.md](serverBackendGo/docs/parity/auth.md) |
| `JWTAuthResource` | `/rest/public/jwt` | `auth` | ✅ | [auth.md](serverBackendGo/docs/parity/auth.md) |
| `SignupResource` | `/rest/public/signup` | `signup` | ⚠️ | Mailchimp، copy-settings |
| `PasswordResetResource` | `/rest/public/passwordReset` | `passwordreset` | ✅ | [auth.md](serverBackendGo/docs/parity/auth.md) |
| `UserResource` | `/rest/private/users` | `users` | ⚠️ | [users.md](serverBackendGo/docs/parity/users.md) — impersonate/superadmin ⊘ |
| `UserRoleResource` | `/rest/private/roles` | `roles` | ✅ | [roles.md](serverBackendGo/docs/parity/roles.md) |
| `CustomerResource` | `/rest/private/customers` | `customers` | ⚠️ | [customers.md](serverBackendGo/docs/parity/customers.md) |
| `SettingsResource` | `/rest/private/settings` | `settings` | ✅ | [settings.md](serverBackendGo/docs/parity/settings.md) |
| `HintResource` | `/rest/private/hints` | `hints` | ✅ | [hints.md](serverBackendGo/docs/parity/hints.md) |
| `SummaryResource` | `/rest/private/summary` | `summary` | ⚠️ | [summary.md](serverBackendGo/docs/parity/summary.md) — charts مبسّطة |
| `DeviceResource` | `/rest/private/devices` | `devices` | ⚠️ | Push notify **Done**; search enrichment still partial |
| `GroupResource` | `/rest/private/groups` | `groups` | ✅ | [groups.md](serverBackendGo/docs/parity/groups.md) |
| `ConfigurationResource` | `/rest/private/configurations` | `configurations` | ✅ | Push notify on save (Phase 9); extended fields in settingsjson |
| `ConfigurationFileResource` | `/rest/private/config-files` | `configfiles` | ⚠️ | [configfiles.md](serverBackendGo/docs/parity/configfiles.md) |
| `ApplicationResource` | `/rest/private/applications` | `applications` | ⚠️ | [applications.md](serverBackendGo/docs/parity/applications.md) — plugin hooks |
| `FilesResource` | `/rest/private/web-ui-files` | `files` | ⚠️ | [files.md](serverBackendGo/docs/parity/files.md) |
| `IconResource` | `/rest/private/icons` | `icons` | ✅ | [icons.md](serverBackendGo/docs/parity/icons.md) |
| **`IconFileResource`** | `/rest/private/icon-files` | `icons` | ✅ | Phase 9 — `POST` multipart + 144px PNG |
| `PublicResource` | `/rest/public` | `publicapi` | ✅ | [publicapi.md](serverBackendGo/docs/parity/publicapi.md) |
| **`PublicFilesResource`** | `/rest/public/files` | — | ⊘ | مُهمل في Java؛ parity يعلن out of scope |
| **`StatsResource`** | `/rest/public/stats` | — | ❌ | `PUT` — إحصائيات استخدام (`UsageStatsDAO`) |
| **`VideosResource`** | `/rest/public/videos` | — | ❌ | رفع/تحميل فيديوهات تدريبية |
| `SyncResource` | `/rest/public/sync` | `sync` | ⚠️ | [sync.md](serverBackendGo/docs/parity/sync.md) |
| `PushApiResource` | `/rest/private/push` | `push` | ⚠️ | [push.md](serverBackendGo/docs/parity/push.md) — FCM |
| `UpdateResource` | `/rest/private/update` | `updates` | ⚠️ | [updates.md](serverBackendGo/docs/parity/updates.md) |
| `QRCodeResource` | `/rest/public/qr` | `qrcode` | ✅ | [qrcode.md](serverBackendGo/docs/parity/qrcode.md) |

### 4.2 Notification (`backend/notification/`)

| Java | Prefix | Go | الحالة |
|------|--------|-----|--------|
| `NotificationResource` | `/rest/notifications` | `notifications` | ✅ |
| `LongPollingServlet` | `/rest/notification/polling/{deviceNumber}` | `notifications` | ✅ (مسار مفرد `notification`) |

### 4.3 Plugins (`backend/plugins/`)

| Java | Go path | Go module | الحالة | ملاحظات |
|------|---------|-----------|--------|---------|
| `PluginResource` | `/rest/plugin/main` | `plugins/platform` | ✅ | [plugins-platform.md](serverBackendGo/docs/parity/plugins-platform.md) |
| `AuditResource` | `/rest/plugins/audit` | `plugins/audit` | ⚠️ | بحث فقط — بدون `AuditFilter` |
| `MessagingResource` | `/rest/plugins/messaging` | `plugins/messaging` | ✅ | [plugins-messaging.md](serverBackendGo/docs/parity/plugins-messaging.md) |
| `PushResource` (plugin) | `/rest/plugins/push` | `plugins/push` | ✅ | CRUD + **cron runner** (Phase 9) |
| `DeviceInfoResource` + settings | `/rest/plugins/deviceinfo/...` | `plugins/deviceinfo` | ⚠️ | [plugins-deviceinfo.md](serverBackendGo/docs/parity/plugins-deviceinfo.md) |
| `DeviceLogResource` + settings | `/rest/plugins/devicelog/...` | `plugins/devicelog` | ⚠️ | [plugins-devicelog.md](serverBackendGo/docs/parity/plugins-devicelog.md) |
| **xtra** (Angular + Liquibase) | — | — | ⊘ | لا REST؛ واجهة قديمة فقط |

### 4.4 Two-factor (ليست Resource منفصلة في server)

| المسار | Go | الحالة |
|--------|-----|--------|
| `/rest/private/twofactor/*` | `twofactor` | ✅ — [auth.md](serverBackendGo/docs/parity/auth.md) |

---

## 5. ما لم يُنقل — تفصيل (❌)

| # | المكوّن Java | المسار / السلوك | التأثير | أولوية مقترحة |
|---|-------------|-----------------|---------|---------------|
| 1 | `StatsResource` | `PUT /rest/public/stats` | تجميع إحصائيات استخدام من الوكلاء/الواجهة | متوسطة — إن كان العميل يرسل stats |
| 2 | `VideosResource` | `POST/GET /rest/public/videos/{fileName}` | رفع فيديوهات مساعدة | منخفضة — إن لم تُستخدم في React الحالي |
| 3 | `IconFileResource` | `POST /rest/private/icon-files` | رفع صورة أيقونة مع resize | **عالية** — قد يكمل `icons` (metadata فقط حالياً) |
| 4 | `PushScheduleTaskModule` | cron لجدولة `plugin_push_schedule` | رسائل مجدولة لا تُرسل تلقائياً | **عالية** للمستخدمين الذين يعتمدون الجدولة |
| 5 | `AuditFilter` (servlet) | تسجيل تلقائي لكل طلب خاص | سجل تدقيق غير مكتمل بدون كتابة تلقائية | متوسطة |
| 6 | `SyncResponseHook` | توسيع `SyncResponse` من plugins | إعدادات/حقول إضافية عند المزامنة | متوسطة — حسب plugins مفعّلة |
| 7 | FCM / `PushSenderMqtt` | إرسال push فعلي للأجهزة | `NoopPush` في `devices` و`configurations` | **عالية** لإشعارات فورية |
| 8 | `MailchimpService` | signup + customer create | تسويق/اشتراك | منخفضة — اختياري في Java أصلاً |

---

## 6. منقول جزئياً — تفصيل (⚠️)

### 6.1 مصادقة وتسجيل

| المجال | Java | Go | الفجوة |
|--------|------|-----|--------|
| Signup | `SignupResource` + Mailchimp | `signup` | بدون Mailchimp، بدون نسخ إعدادات افتراضية كاملة |
| Users | `impersonate`, `superadmin/*` | ⊘ | مُعلَن out of scope — React لا يستخدم `GET /{id}` |

### 6.2 عملاء وأجهزة ولوحة

| المجال | الفجوة |
|--------|--------|
| Customers `PUT /` | لا نسخ أجهزة/تكوينات افتراضية عند إنشاء tenant |
| Devices `POST /search` | فلاتر `mdmMode`, `launcherVersion`, `deviceStatuses`؛ إثراء apps/files |
| Devices notify | `NoopPush` — لا FCM |
| Summary | `devicestatuses` / charts تثبيت حسب config — مبسّط |

### 6.3 تكوينات وملفات وتطبيقات

| المجال | الفجوة |
|--------|--------|
| Configurations save | `NoopPushNotifier` — لا إشعار أجهزة بعد الحفظ |
| Config files | قرص فقط؛ `uploadedfiles` + quota `sizeLimit` |
| Applications | plugin hooks عند الحفظ |
| Files | push عند ربط configurations؛ تحميل `/files/*` للوكلاء |

### 6.4 وكلاء (Phase 7)

| المجال | الفجوة |
|--------|--------|
| Sync | `SyncResponseHook`، دمج settings كامل |
| Updates | تحميل APK بعيد؛ `sendStats` stub |
| Push API | إرسال فعلي للجهاز |

### 6.5 Plugins (Phase 8)

| Plugin | Endpoints Java غير موجودة في Go | الفجوة سلوكية |
|--------|----------------------------------|---------------|
| **deviceinfo** | `POST .../private/search/device`, `POST .../private/export`, `GET settings/device/{deviceNumber}` | GPS/WiFi multi-table، تصدير |
| **devicelog** | `POST .../private/search/export`, `GET .../rules/{deviceNumber}` | تصدير CSV؛ قواعد لكل جهاز |
| **audit** | — (يُكتب عبر Filter) | لا تسجيل تلقائي |
| **push** | — | لا `PushScheduleTaskModule` cron |

---

## 7. مكونات Java غير REST (خلفية)

| المكوّن | الموقع التقريبي | في Go |
|---------|-----------------|-------|
| `AuditFilter` | `plugins/audit` | ❌ |
| `LongPollingServlet` | `notification` | ✅ (HTTP handler) |
| `BackgroundTaskRunnerService` | `common` | جزئي — لا جدولة push |
| `SyncResponseHook` + Guice multibind | `SyncResource` | ❌ |
| `CustomerCreatedEventListener` | devicelog / deviceinfo plugins | ❌ |
| `InsertDeviceLogRecordsTask` | devicelog | ❌ (معالجة دفعات) |
| MyBatis DAOs | `backend/common/.../persistence/` | Postgres repos — غالباً نفس الجداول |
| Liquibase changelogs | `backend/**/liquibase/` | `serverBackendGo/db/migrations/` — مراجعة يدوية لكل جدول |

---

## 8. Endpoints Java — قائمة سريعة بالحالة

### منقول ✅ (مسارات رئيسية موجودة في Go)

- Auth: `options`, `login`, `logout`
- JWT: `login`
- Password reset + signup (core paths)
- Users: `current`, `details`, `all`, CRUD `other`, `roles`
- Roles: permissions, CRUD
- Customers: search, impersonate, edit, prefix, delete (create partial)
- Settings, hints, summary/devices
- Devices, groups — CRUD + bulk + app settings
- Configurations — full CRUD + copy + applications
- Applications — full catalog + admin common
- Config-files POST
- Files, icons (metadata CRUD)
- Public: name, logo, applications/upload
- Sync: configuration GET/POST, info, applicationSettings
- Notifications + polling
- Private push POST, QR, updates check
- Plugin platform, audit search, messaging, push plugin CRUD, deviceinfo core, devicelog core

### غير منقول ❌

```
PUT  /rest/public/stats
POST /rest/public/videos
GET  /rest/public/videos/{fileName}
POST /rest/private/icon-files
```

### جزئي / ناقص endpoint ⚠️

```
POST /rest/plugins/deviceinfo/deviceinfo/private/search/device
POST /rest/plugins/deviceinfo/deviceinfo/private/export
GET  /rest/plugins/deviceinfo/deviceinfo-plugin-settings/device/{deviceNumber}
POST /rest/plugins/devicelog/log/private/search/export
GET  /rest/plugins/devicelog/log/rules/{deviceNumber}
```

(+ سلوك: FCM, configuration push notify, push schedule cron, audit auto-write)

---

## 9. وحدات Go بدون مكافئ Java مباشر

| Go module | ملاحظة |
|-----------|--------|
| `twofactor` | منطق webapp قديم — مجمّع في Go |
| `shared/targets`, `shared/status` | مساعدة لـ plugins push/messaging |

---

## 10. قاعدة البيانات والمخطط

- Go يستخدم **نفس مخطط Postgres legacy** عبر migrations في `serverBackendGo/db/migrations/`.
- جداول plugins: `000010_plugins_core` وما قبلها.
- **تحقق يدوي مطلوب** لأي جدول يُستخدم في Java DAO ولم يُذكر في migration Go (خصوصاً deviceinfo GPS tables، `usagestats`, `videos`).

---

## 11. تأثير الواجهة الأمامية (`frontend/`)

| الخدمة | المسار المتوقع | مخاطر إن بقيت فجوة |
|--------|----------------|---------------------|
| `customersService` | customers search/impersonate | إنشاء tenant ناقص |
| `pushService` | private + plugin push | جدولة لا تُنفَّذ بدون cron |
| تطبيقات/تكوينات | configurations + files | عدم وصول push للأجهزة |
| Plugins UI (Angular assets في Java) | `/plugins/*/webapp` | **خارج نطاق Go** — يحتاج React أو proxy static |

---

## 12. خطة عمل مقترحة (لإغلاق الهجرة)

مرتبة حسب الأثر على التشغيل اليومي:

| P | المهمة | Java مرجع | Go هدف |
|---|--------|-----------|--------|
| **P0** | FCM / push notifier حقيقي | `PushSender*`, notify في Device/Configuration | `devices`, `configurations`, `files` |
| **P0** | Push schedule cron | `PushScheduleTaskModule` | `plugins/push` worker |
| **P1** | `IconFileResource` | `IconFileResource.java` | `icons` أو `files` |
| **P1** | deviceinfo export + search/device | `DeviceInfoResource` | `plugins/deviceinfo` |
| **P1** | devicelog export + rules/{device} | `DeviceLogResource` | `plugins/devicelog` |
| **P2** | AuditFilter → middleware | `AuditFilter` | `plugins/audit` + platform middleware |
| **P2** | SyncResponseHook registry | `SyncResource` | `sync` + plugin registry |
| **P2** | Customers create defaults | `CustomerResource` | `customers` |
| **P3** | Stats, Videos | `StatsResource`, `VideosResource` | وحدات جديدة أو دمج `publicapi` |
| **P3** | Updates APK download + stats POST | `UpdateResource` | `updates` |
| **P3** | Devices search enrichment | `DeviceResource` | `devices` |

---

## 13. كيف تُحدَّث هذا الوثيقة

1. بعد كل phase: حدّث الصف في §4 وملف `serverBackendGo/docs/parity/<module>.md`.
2. شغّل smoke من `serverBackendGo/scripts/dev.sh` وسجّل النتيجة في §12.
3. عند إغلاق فجوة: غيّر الحالة من ⚠️/❌ إلى ✅ وأزل السطر من §5 أو §6.

---

## 14. مراجع سريعة

- Spec Kit feature (تحليل الفجوات): [`specs/010-java-go-gap-analysis/spec.md`](specs/010-java-go-gap-analysis/spec.md)
- Spec Kit feature (إكمال الفجوات): [`specs/011-complete-migration-gaps/spec.md`](specs/011-complete-migration-gaps/spec.md)
- الخطوات التالية: [`serverBackendGo/docs/NEXT_STEPS.md`](serverBackendGo/docs/NEXT_STEPS.md)
- دستور المشروع: [`.specify/memory/constitution.md`](.specify/memory/constitution.md)

---

*تم إنشاء هذا الملف آلياً من جرد الكود ووثائق parity — راجع الكود قبل اعتبار أي بند ❌ نهائياً إذا أُضيفت مسارات بعد تاريخ التحليل.*
