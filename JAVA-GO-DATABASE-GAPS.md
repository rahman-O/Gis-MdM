# تحليل فجوات قاعدة البيانات: Java (`backend`) → Go (`serverBackendGo`)

**التاريخ:** 2026-05-21  
**المراجع:**
- Java Liquibase: [`backend/server/src/main/resources/liquibase/db.changelog.xml`](backend/server/src/main/resources/liquibase/db.changelog.xml) + changelogs الإضافات في [`backend/plugins/`](backend/plugins/) و [`backend/notification/`](backend/notification/)
- Go migrations: [`serverBackendGo/db/migrations/`](serverBackendGo/db/migrations/) (`000001` … `000010` + `000008_devices_search_extras`)
- قائمة الجداول الكاملة (تثبيت Headwind): [`backend/hmdm_install.sh`](backend/hmdm_install.sh) (أمر `DROP TABLE`)

**مرتبط بـ:** [`JAVA-GO-BACKEND-GAPS.md`](JAVA-GO-BACKEND-GAPS.md) (سلوك REST)، [`JAVA-GO-MIGRATION-STATUS.md`](JAVA-GO-MIGRATION-STATUS.md)  
**مواصفة الإكمال:** [`specs/013-complete-database-gaps/spec.md`](specs/013-complete-database-gaps/spec.md)

---

## 1. الملخص التنفيذي

| المؤشر | Java (WAR كامل) | Go (`serverBackendGo`) |
|--------|-----------------|-------------------------|
| **جداول أساسية MDM** | ~28 | ~26 ✅ معظمها موجود |
| **جداول إضافات (plugins اختيارية)** | ~20 | ~8 جزئي |
| **جداول مفقودة بالكامل** | — | **~19** (انظر §3) |
| **أعمدة `configurations`** | ~60+ عمود منفصل | **مجمّعة في `settingsjson` JSONB** + جزء في أعمدة صريحة |
| **أعمدة `settings` (واجهة الأعمدة)** | ~25+ `columnDisplayed*` | **غير موجودة** (لا `userrolesettings`) |
| **دوال SQL مساعدة** | `mdm_device_launcher_version`, `mdm_config_app_upgrade`, … | **غير منقولة** |

**الاستنتاج:** قاعدة Go **كافية للمسار التشغيلي الحالي** (auth، أجهزة، تكوينات، ملفات، push، plugins أساسية). الفجوات الحرجة لـ **parity كامل مع Java + React** تتمحور حول: `devicestatuses`, `userrolesettings`, `configurationapplicationparameters`, `usagestats`, وجداول telemetry الفرعية لـ deviceinfo، وأعمدة `customers`/`settings` الإضافية.

---

## 2. منهجية المقارنة

1. استخراج أسماء الجداول من `hmdm_install.sh` (قائمة DROP) كمرجع **schema كامل** لنشر Java.
2. مقارنتها بجداول `CREATE TABLE` في migrations Go.
3. لكل جدول مشترك: مقارنة أعمدة Liquibase `ALTER TABLE … ADD` مقابل `ALTER`/`CREATE` في Go.
4. تمييز:
   - **⊘** جدول/عمود غير مطلوب إذا لم يُفعَّل plugin أو لم يُنفَّذ module في Go بعد.
   - **⚠️** موجود جزئياً أو بديل تصميمي (`settingsjson`).
   - **❌** ناقص ويؤثر على ميزة موثّقة في Java/React.

---

## 3. جداول Java — حالة التغطية في Go

### 3.1 جداول أساسية — موجودة في Go ✅

| الجدول (PostgreSQL lowercase) | Java | Go migration | ملاحظة |
|------------------------------|------|--------------|--------|
| `customers` | ✅ | `000001`, `000005`, `000008_files` | جزئي أعمدة (§4.1) |
| `users` | ✅ | `000001`, `000002` | ✅ |
| `permissions` | ✅ | `000001`, seeds في `000006+` | ✅ |
| `userroles` | ✅ | `000001` | ✅ |
| `userrolepermissions` | ✅ | `000001` | ✅ |
| `groups` | ✅ | `000001`, `000006` | `customerid` في `000006` |
| `configurations` | ✅ | `000001`, `000006`, `000007` | سياسات MDM في `settingsjson` |
| `devices` | ✅ | `000006`, `000008_devices_search_extras` | `info`, `infojson`, `imeiupdatets` |
| `devicegroups` | ✅ | `000006` | ✅ |
| `deviceapplicationsettings` | ✅ | `000006` | ✅ |
| `applications` | ✅ | `000007` | ✅ |
| `applicationversions` | ✅ | `000007`, `000016` | ✅ `apkhash` |
| `configurationapplications` | ✅ | `000007`, `000016` | ✅ `remove`, `longtap` |
| `configurationfiles` | ✅ | `000007`, `000008_files` | `fileid` في `000008` |
| `configurationapplicationsettings` | ✅ | `000007` | ✅ |
| `settings` | ✅ | `000001`, `000003`, `000015` | ✅ tenant columns؛ عرض الأعمدة في `userrolesettings` |
| `userdevicegroupsaccess` | ✅ | `000001` | ✅ |
| `userconfigurationaccess` | ✅ | `000001` | ✅ |
| `userhints` / `userhinttypes` | ✅ | `000004` | ✅ |
| `pendingsignup` | ✅ | `000002` | ✅ |
| `uploadedfiles` | ✅ | `000008_files` | ✅ |
| `icons` | ✅ | `000008_files` | ✅ |
| `pushmessages` / `pendingpushes` | ✅ | `000009` | ✅ |
| `plugins` / `pluginsdisabled` | ✅ | `000010` | ✅ |
| `plugin_audit_log` | ✅ | `000010` | ✅ |
| `plugin_messaging_messages` | ✅ | `000010` | ✅ |
| `plugin_push_messages` / `plugin_push_schedule` | ✅ | `000009`, `000010` | ✅ |
| `plugin_deviceinfo_settings` | ✅ | `000010` | ✅ |
| `plugin_deviceinfo_deviceparams` | ✅ | `000010` | جزئي (§3.2) |
| `plugin_deviceinfo_deviceparams_device` | ✅ | `000010` | ✅ |
| `plugin_devicelog_*` | ✅ | `000010` | ✅ |

### 3.2 جداول أساسية — مفقودة في Go ❌

| الجدول | Java / الاستخدام | الأولوية | تأثير |
|--------|------------------|----------|--------|
| **`devicestatuses`** | حالة تثبيت التطبيقات/الملفات لكل جهاز | **P1** | ✅ `000011` — devices search + summary |
| **`userrolesettings`** | أعمدة جدول الأجهزة (`columnDisplayed*`) | **P1** | ✅ `000012` — settings API |
| **`configurationapplicationparameters`** | `skipVersionCheck` لكل (config, app) | **P2** | ✅ `000013` — config save upsert |
| **`usagestats`** | إحصائيات نسخة الخادم | **P2** | ✅ `000014` — جاهز لـ 012 `stats` |
| **`trialkey`** | مفاتيح تجريبية | **P3** | تسجيل/ترخيص تجريبي |
| **`applicationfilestocopytemp`** | نسخ ملفات APK مؤقت | **P3** | رفع/نسخ تطبيقات |
| **`applicationversionstemp`** | جداول مؤقتة أثناء الرفع | **P3** | إدارة إصدارات |

### 3.3 جداول plugins — موجودة في Java فقط ⊘ / ❌ (deferred في 013 v1)

هذه الجداول تظهر في `hmdm_install.sh` عند تثبيت **حزمة plugins كاملة**؛ **خارج نطاق 013** — spec per-plugin فقط (لا `000018+`):

| الجدول | Plugin | في Go |
|--------|--------|-------|
| `plugin_deviceinfo_deviceparams_wifi` | deviceinfo | ❌ |
| `plugin_deviceinfo_deviceparams_gps` | deviceinfo | ❌ |
| `plugin_deviceinfo_deviceparams_mobile` | deviceinfo | ❌ |
| `plugin_deviceinfo_deviceparams_mobile2` | deviceinfo | ❌ |
| `plugin_devicelocations_history` | devicelocations | ❌ |
| `plugin_devicelocations_latest` | devicelocations | ❌ |
| `plugin_devicelocations_settings` | devicelocations | ❌ |
| `plugin_devicereset_status` | devicereset | ❌ |
| `plugin_apuppet_data` / `plugin_apuppet_settings` | apuppet | ❌ |
| `plugin_knox_rules` | knox | ❌ |
| `plugin_openvpn_defaults` | openvpn | ❌ |
| `plugin_photo_*` (4 جداول) | photo | ❌ |
| `plugin_urlfilter_lists` | urlfilter | ❌ |

**ملاحظة:** Go يوفّر جدولاً رئيسياً `plugin_deviceinfo_deviceparams` + `…_device` فقط؛ Java يقسّم GPS/WiFi/mobile في جداول فرعية.

---

## 4. فجوات الأعمدة (جداول مشتركة)

### 4.1 `customers`

| العمود (Java) | Go | الحالة |
|---------------|-----|--------|
| `id`, `name`, `description`, `master`, `filesdir`, `prefix`, `lastlogintime` | ✅ `000001` | ✅ |
| `email`, `accounttype`, `customerstatus`, `registrationtime`, `expirytime`, `devicelimit`, `deviceconfigurationid` | ✅ `000005` | ✅ |
| `sizelimit` | ✅ `000008_files` | ✅ |
| `firstname`, `lastname`, `language` | ❌ | ⊘ تسويق/حساب متقدم |
| `inactivestate`, `pausestate`, `abandonstate` | ❌ | ⊘ حالات اشتراك |
| `signupstatus`, `signuptoken` | ❌ | ⊘ تسجيل ذاتي |

### 4.2 `devices`

| العمود (Java) | Go | الحالة |
|---------------|-----|--------|
| `number`, `description`, `lastupdate`, `configurationid`, `customerid` | ✅ | ✅ |
| `info`, `infojson`, `imei`, `phone`, `enrolltime`, `publicip` | ✅ | ✅ |
| `custom1`–`custom3`, `oldnumber`, `fastsearch` | ✅ | ✅ |
| `imeiupdatets` | ✅ `000008_devices_search_extras` | ✅ (012) |
| `groupid` (قديم — استبدل بـ `devicegroups`) | — | ✅ نموذج many-to-many |

### 4.3 `configurations`

Java يخزّن عشرات سياسات MDM كأعمدة منفصلة. Go يخزّن جزءاً صريحاً + **`settingsjson` JSONB** (`000007`).

| فئة الأعمدة (Java) | أمثلة | Go |
|--------------------|--------|-----|
| هوية/تصميم | `password`, `backgroundcolor`, `textcolor`, `qrcodekey`, `baseurl`, `defaultfilepath` | أعمدة صريحة ✅ |
| تطبيق رئيسي | `mainappid`, `contentappid`, `permissive` | ✅ (`mainappid` في `000006`) |
| شبكة/جهاز | `gps`, `bluetooth`, `wifi`, `mobiledata`, `usbstorage`, `wifissid`, … | ⚠️ داخل `settingsjson` إن كتبها الـ API |
| kiosk / launcher | `kioskmode`, `kioskhome`, `rundefaultlauncher`, … | ⚠️ `settingsjson` |
| تحديثات نظام/تطبيق | `systemupdatetype`, `scheduleappupdate`, `requestupdates`, … | ⚠️ `settingsjson` |
| سياسات متقدمة | `restrictions`, `allowedclasses`, `adminextras`, `encryptdevice`, … | ⚠️ `settingsjson` |

**فجوة عملية:** استيراد dump Java قديم → أعمدة `configurations` الفارغة في Go بينما Java يتوقع أعمدة NOT NULL؛ الحل عند الدمج: migration ترحيل عمود→JSON أو أعمدة صريحة لكل مفتاح يقرأه React.

**أعمدة Java غير ممثلة حتى في `settingsjson` (إن لم يُملأ من التطبيق):**

- `usedefaultdesignsettings`, `iconsize`, `desktopheader` (على مستوى configuration)
- `eventreceivingcomponent`, `pushoptions`, `desktopheadertemplate`
- `type` (قد يُستخدم في Go كـ `type INT` ✅)

### 4.4 `settings`

| العمود (Java) | Go (`000003` + `000001`) | الحالة |
|---------------|-------------------------|--------|
| `customerid`, `twofactor`, `idlelogout` | ✅ | ✅ |
| `language`, `usedefaultlanguage`, `createnewdevices`, `newdeviceconfigurationid` | ✅ | ✅ |
| `passwordlength`, `passwordstrength`, ألوان/تصميم عام | ✅ | ✅ |
| **`columndisplayeddevice*`** (12+ عمود) | ❌ | ❌ → يُفترض `userrolesettings` |
| `newdevicegroupid` | ❌ | ❌ |
| `phonenumberformat` | ❌ | P2 |
| `custompropertyname1`–`3`, `custommultiline*`, `customsend*` | ❌ | P2 |
| `desktopheadertemplate`, `senddescription`, `passwordreset` | ❌ | جزئي (`passwordreset` على `users`) |

### 4.5 `applications` / `applicationversions` / `configurationapplications`

| العنصر | Java | Go | الحالة |
|--------|------|-----|--------|
| `apkhash` على `applicationversions` | ✅ | ❌ | P2 تحقق APK |
| `remove` على `configurationapplications` | ✅ | ❌ | P2 إزالة تطبيق من config |
| `longtap` على `configurationapplications` | ✅ | ❌ | P3 UI launcher |
| `configurationapplications.remove` vs ملفات | `configurationfiles.remove` | ✅ في Go للملفات فقط | ⚠️ |

### 4.6 `users`

| العمود | Java | Go | ملاحظة |
|--------|------|-----|--------|
| `allconfigavailable` | ✅ | ❌ | قد يُستبدل بمنطق التطبيق |
| `userrole` (نص قديم) | ✅ | — | استُبدل بـ `userroleid` |
| حقول 2FA | ✅ | ✅ | ✅ |

### 4.7 `plugin_audit_log`

| العمود | Java (audit plugin) | Go `000010` | ملاحظة |
|--------|---------------------|-------------|--------|
| `payload` / `details` | تفاصيل الطلب | `payload` | تسمية مختلفة، نفس الغرض |
| حقول إضافية في Java changelog | راجع `audit.changelog.xml` | — | قارن عند تفعيل AuditFilter |

---

## 5. دوال وإجراءات SQL (Java فقط)

| الاسم | الغرض | في Go |
|------|--------|-------|
| `mdm_device_launcher_version(pkg, info)` | فلتر/ترتيب `launcherVersion` | ❌ — Go يستخدم `infojson->>'launcherVersion'` |
| `mdm_resolve_device_property(col, json)` | دمج IMEI/phone من عمود + infojson | ❌ — منطق في Go SQL بسيط |
| `mdm_device_permissions_index(info)` | ترتيب عمود Permissions | ❌ |
| `mdm_config_app_upgrade(configId, appId)` | ترقية إصدار التطبيق في التكوين | ❌ — منطق في `applications` service |

---

## 6. أدوات الهجرة والبيانات الوصفية

| الموضوع | Java | Go |
|---------|------|-----|
| أداة الهجرة | Liquibase (`databasechangelog`) | golang-migrate (`schema_migrations`) |
| استيراد dump Java → Go | — | قد تحتاج `000011_legacy_import.sql` لتطبيع أسماء الأعمدة |
| Seed | بيانات في Liquibase + install script | seeds في `000001`–`000010` |

---

## 7. خريطة الأولويات → migrations في Go (013)

| ID | Migration | الحالة | يرتبط بـ |
|----|-----------|--------|----------|
| `000011` | `devicestatuses_core` | ✅ 2026-05-21 | devices search، summary |
| `000012` | `userrolesettings_core` | ✅ | settings `userRole` API |
| `000013` | `configuration_application_parameters` | ✅ | config editor `skipVersionCheck` |
| `000014` | `usagestats_core` | ✅ | 012 `stats` INSERT |
| `000015` | `settings_columns_extend` | ✅ | tenant settings |
| `000016` | `applications_columns_extend` | ✅ | `apkhash`, `remove`, `longtap` |
| `000017` | `configurations_legacy_import` | ✅ no-op on greenfield | Java dump → `settingsjson` |
| `000018+` | plugins اختيارية | **⊘ deferred** | per-plugin specs |

---

## 8. توافق مع dump قاعدة Java موجود

إذا ربطت الفرونت/React بقاعدة **منسوخة من Java**:

1. **جداول ناقصة في Go** (`devicestatuses`, …) → استعلامات Java/MyBatis القديمة **تفشل** حتى لو REST في Go لا يستدعيها بعد.
2. **`configurations`**: dump يحتوي أعمدة؛ Go يتوقع `settingsjson` — يلزم **سكربت ترحيل بيانات** عند الانتقال.
3. **`devices.infojson`**: موجود في Go ✅؛ تأكد من `UPDATE devices SET infojson = info::jsonb` كما في Liquibase Java.

---

## 9. ملخص سريع للمطور

```
Java tables (install script)     ≈ 55
Go CREATE TABLE migrations       ≈ 36
Missing critical for core MDM    ≈  6  (devicestatuses, userrolesettings, cap, usagestats, trialkey, temps)
Missing optional plugins         ≈ 13
Configurations column parity     → mostly settingsjson (design shift)
Settings UI columns              → need userrolesettings table
```

**الخطوة التالية:** 012 REST (`stats` INSERT على `usagestats`، إغلاق gaps المتبقية في [`JAVA-GO-BACKEND-GAPS.md`](JAVA-GO-BACKEND-GAPS.md)).

---

*آخر تحديث: 2026-05-21 — فرع `013-complete-database-gaps`، migrations `000011`–`000017`.*
