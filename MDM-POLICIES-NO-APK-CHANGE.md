# سياسات وقيود MDM بدون تعديل APK

## المقدمة

هذا المستند يوثّق جميع السياسات والقيود التي يمكن تطبيقها على أجهزة Android المُدارة **بدون أي تعديل على كود تطبيق hmdm-android (APK)**. التطبيق الحالي يدعم هذه السياسات عبر آليتين:

1. **حقول `settingsjson`** — يقرأها التطبيق من `SyncResponse` ويطبّقها مباشرة
2. **حقل `restrictions`** — قائمة Android UserRestrictions يطبّقها عبر `DevicePolicyManager.addUserRestriction()`

**المتطلب الوحيد:** الجهاز مسجّل كـ Device Owner (عبر QR enrollment).

---

## القسم الأول: سياسات settingsjson (51 سياسة)

هذه سياسات يقرأها التطبيق مباشرة من استجابة المزامنة ويطبّقها.

### 1.1 التحكم بالأجهزة (Hardware Controls)

| # | السياسة | الحقل | القيم | الوصف |
|---|---------|-------|-------|-------|
| 1 | فرض GPS | `gps` | `true/false` | تشغيل أو إيقاف GPS إجبارياً |
| 2 | فرض Bluetooth | `bluetooth` | `true/false` | تشغيل أو إيقاف Bluetooth |
| 3 | فرض WiFi | `wifi` | `true/false` | تشغيل أو إيقاف WiFi |
| 4 | فرض بيانات الجوال | `mobileData` | `true/false` | تشغيل أو إيقاف Mobile Data |
| 5 | فرض USB Storage | `usbStorage` | `true/false` | تشغيل أو إيقاف USB |
| 6 | منع لقطات الشاشة | `disableScreenshots` | `true` | تعطيل Screenshots بالكامل |
| 7 | منع طلب صلاحية الموقع | `disableLocation` | `true` | منع التطبيقات من طلب GPS |

### 1.2 الشاشة والعرض (Display & UI)

| # | السياسة | الحقل | القيم | الوصف |
|---|---------|-------|-------|-------|
| 8 | قفل Status Bar | `lockStatusBar` | `true` | منع سحب شريط الإشعارات |
| 9 | سطوع تلقائي | `autoBrightness` | `true/false` | تفعيل/تعطيل السطوع التلقائي |
| 10 | تحديد السطوع | `brightness` | `0-255` | ضبط السطوع على قيمة ثابتة |
| 11 | إدارة Timeout | `manageTimeout` | `true` | تفعيل التحكم بوقت إطفاء الشاشة |
| 12 | تحديد Timeout | `timeout` | ثواني | وقت إطفاء الشاشة (مثلاً 60) |
| 13 | قفل الاتجاه | `orientation` | `0/1/2` | 0=حر, 1=عمودي, 2=أفقي |
| 14 | حجم الأيقونات | `iconSize` | `"SMALL"/"MEDIUM"/"LARGE"` أو `100/120/140` | حجم أيقونات اللانشر |
| 15 | لون الخلفية | `backgroundColor` | `"#RRGGBB"` | لون خلفية شاشة اللانشر |
| 16 | لون النص | `textColor` | `"#RRGGBB"` | لون نصوص اللانشر |
| 17 | صورة الخلفية | `backgroundImageUrl` | URL | صورة خلفية مخصصة |

### 1.3 الصوت (Audio)

| # | السياسة | الحقل | القيم | الوصف |
|---|---------|-------|-------|-------|
| 18 | قفل الصوت | `lockVolume` | `true` | منع المستخدم من تغيير الصوت |
| 19 | إدارة الصوت | `manageVolume` | `true` | تفعيل التحكم بمستوى الصوت |
| 20 | تحديد مستوى الصوت | `volume` | `0-100` | ضبط الصوت على نسبة ثابتة |

### 1.4 وضع Kiosk (Kiosk Mode)

| # | السياسة | الحقل | القيم | الوصف |
|---|---------|-------|-------|-------|
| 21 | تفعيل Kiosk | `kioskMode` | `true` | قفل الجهاز على تطبيق واحد |
| 22 | زر Home في Kiosk | `kioskHome` | `true` | إظهار/إخفاء زر Home |
| 23 | قائمة Recent في Kiosk | `kioskRecents` | `true` | إظهار/إخفاء Recent Apps |
| 24 | الإشعارات في Kiosk | `kioskNotifications` | `true` | إظهار/إخفاء الإشعارات |
| 25 | معلومات النظام في Kiosk | `kioskSystemInfo` | `true` | إظهار/إخفاء System Info |
| 26 | Keyguard في Kiosk | `kioskKeyguard` | `true` | إظهار/إخفاء شاشة القفل |
| 27 | قفل الأزرار في Kiosk | `kioskLockButtons` | `true` | قفل أزرار الجهاز الفيزيائية |
| 28 | إبقاء الشاشة مضاءة | `kioskScreenOn` | `true` | منع إطفاء الشاشة في Kiosk |
| 29 | طريقة الخروج من Kiosk | `kioskExit` | `"password"/"back"/"none"` | كيف يخرج المشرف |

### 1.5 الأمان وكلمات المرور (Security)

| # | السياسة | الحقل | القيم | الوصف |
|---|---------|-------|-------|-------|
| 30 | سياسة كلمة المرور | `passwordMode` | `"strong"/"numeric"/"any"/"none"` | نوع كلمة المرور المطلوبة |
| 31 | كلمة مرور المشرف | `password` | نص | كلمة مرور لوحة المشرف على الجهاز |
| 32 | قفل الإعدادات الآمنة | `lockSafeSettings` | `true` | منع تغيير إعدادات محددة |

### 1.6 التطبيقات والتحديثات (Apps & Updates)

| # | السياسة | الحقل | القيم | الوصف |
|---|---------|-------|-------|-------|
| 33 | وضع Permissive | `permissive` | `true` | السماح بكل التطبيقات (لا قيود) |
| 34 | التطبيق الرئيسي | `mainApp` | `"com.pkg.name"` | تحديد التطبيق الرئيسي |
| 35 | تشغيل كـ Launcher | `runDefaultLauncher` | `true/false` | هل يعمل كلانشر افتراضي |
| 36 | Autostart في المقدمة | `autostartForeground` | `true` | إبقاء التطبيقات المشغّلة في المقدمة |
| 37 | جدولة تحديث التطبيقات | `scheduleAppUpdate` | `true` | تحديث تلقائي مجدول |
| 38 | صلاحيات التطبيقات | `appPermissions` | `"grant"/"deny"/"default"` | منح/رفض صلاحيات تلقائياً |
| 39 | الأنشطة المسموحة | `allowedClasses` | `"com.android.settings.Settings,..."` | whitelist activities |
| 40 | نوع تحديث النظام | `systemUpdateType` | `0/1/2` | 0=تلقائي, 1=مؤجل, 2=نافذة زمنية |
| 41 | بداية نافذة تحديث النظام | `systemUpdateFrom` | `"02:00"` | وقت بداية التحديث |
| 42 | نهاية نافذة تحديث النظام | `systemUpdateTo` | `"05:00"` | وقت نهاية التحديث |
| 43 | بداية نافذة تحديث التطبيقات | `appUpdateFrom` | `"01:00"` | وقت بداية تحديث apps |
| 44 | نهاية نافذة تحديث التطبيقات | `appUpdateTo` | `"04:00"` | وقت نهاية تحديث apps |
| 45 | طريقة تحميل التحديثات | `downloadUpdates` | `"wifi"/"any"` | WiFi فقط أو أي شبكة |

### 1.7 الشبكة والاتصال (Network & Communication)

| # | السياسة | الحقل | القيم | الوصف |
|---|---------|-------|-------|-------|
| 46 | خيارات Push | `pushOptions` | `"mqttAlarm"/"polling"` | طريقة استقبال الأوامر |
| 47 | تتبع الموقع | `requestUpdates` | `"1"/"5"/"15"` | فترة إرسال الموقع (دقائق) |
| 48 | Keepalive time | `keepaliveTime` | ثواني | فترة keepalive لـ MQTT |
| 49 | إظهار WiFi عند خطأ | `showWifi` | `true` | إظهار إعدادات WiFi تلقائياً عند فقد الاتصال |

### 1.8 إدارة الجهاز (Device Management)

| # | السياسة | الحقل | القيم | الوصف |
|---|---------|-------|-------|-------|
| 50 | المنطقة الزمنية | `timeZone` | `"Asia/Baghdad"/"auto"` | فرض timezone محدد |
| 51 | رابط سيرفر جديد | `newServerUrl` | `"https://new.server"` | ترحيل الجهاز لسيرفر آخر |

---

## القسم الثاني: قيود restrictions (Android UserRestrictions) — 42 قيد

هذه قيود يطبّقها التطبيق عبر `DevicePolicyManager.addUserRestriction()`. تُضاف كقائمة مفصولة بفواصل في حقل `restrictions`.

### 2.1 أمان وحماية

| # | القيد | القيمة | الوصف | Android Min |
|---|-------|--------|-------|-------------|
| 52 | منع الكاميرا | `no_camera` | تعطيل الكاميرا بالكامل | 5.0 |
| 53 | منع Factory Reset | `no_factory_reset` | منع إعادة ضبط المصنع | 5.0 |
| 54 | منع Safe Boot | `no_safe_boot` | منع الدخول لـ Safe Mode | 5.0 |
| 55 | منع USB Debugging | `no_debugging_features` | تعطيل تصحيح USB | 5.0 |
| 56 | منع إيقاف التشغيل | `no_shutdown` | منع إطفاء الجهاز | 9.0 |

### 2.2 تطبيقات

| # | القيد | القيمة | الوصف | Android Min |
|---|-------|--------|-------|-------------|
| 57 | منع تثبيت تطبيقات | `no_install_apps` | منع تثبيت أي تطبيق يدوياً | 5.0 |
| 58 | منع حذف تطبيقات | `no_uninstall_apps` | منع حذف أي تطبيق | 5.0 |
| 59 | منع مصادر مجهولة | `no_install_unknown_sources` | حظر sideloading | 5.0 |
| 60 | منع مصادر مجهولة (عام) | `no_install_unknown_sources_globally` | حظر كامل لكل المستخدمين | 8.0 |

### 2.3 شبكة واتصالات

| # | القيد | القيمة | الوصف | Android Min |
|---|-------|--------|-------|-------------|
| 61 | منع تغيير WiFi | `no_config_wifi` | لا يمكن إضافة/حذف شبكات | 5.0 |
| 62 | منع تغيير VPN | `no_config_vpn` | لا يمكن إعداد VPN | 5.0 |
| 63 | منع Tethering | `no_config_tethering` | منع مشاركة الإنترنت (hotspot) | 5.0 |
| 64 | منع مشاركة Bluetooth | `no_bluetooth_sharing` | منع إرسال ملفات عبر BT | 5.0 |
| 65 | منع Bluetooth بالكامل | `no_bluetooth` | تعطيل Bluetooth كلياً | 5.0 |
| 66 | منع SMS | `no_sms` | منع إرسال/استقبال رسائل | 5.0 |
| 67 | منع المكالمات الصادرة | `no_outgoing_calls` | منع الاتصال | 5.0 |
| 68 | منع تغيير شبكات الجوال | `no_config_mobile_networks` | منع تغيير إعدادات APN | 5.0 |
| 69 | منع NFC Beam | `no_outgoing_beam` | منع إرسال بيانات عبر NFC | 5.0 |
| 70 | منع وضع الطيران | `no_airplane_mode` | منع تفعيل Airplane Mode | 9.0 |
| 71 | منع تغيير Cell Broadcasts | `no_config_cell_broadcasts` | منع تغيير إعدادات CBS | 5.0 |
| 72 | منع Data Roaming | `no_data_roaming` | منع التجوال | 5.0 |

### 2.4 حسابات ومستخدمين

| # | القيد | القيمة | الوصف | Android Min |
|---|-------|--------|-------|-------------|
| 73 | منع تعديل الحسابات | `no_modify_accounts` | لا يمكن إضافة/حذف حسابات | 5.0 |
| 74 | منع إضافة مستخدمين | `no_add_user` | منع إنشاء مستخدمين جدد | 5.0 |
| 75 | منع حذف مستخدمين | `no_remove_user` | منع حذف مستخدمين | 5.0 |
| 76 | منع تبديل المستخدمين | `no_switch_user` | منع التبديل بين المستخدمين | 9.0 |

### 2.5 ملفات ونقل بيانات

| # | القيد | القيمة | الوصف | Android Min |
|---|-------|--------|-------|-------------|
| 77 | منع USB file transfer | `no_usb_file_transfer` | منع نقل ملفات عبر USB | 5.0 |
| 78 | منع SD Card | `no_physical_media` | منع mount وسائط خارجية | 5.0 |
| 79 | منع مشاركة الموقع | `no_share_location` | منع مشاركة GPS مع تطبيقات | 5.0 |
| 80 | منع Cross-profile copy | `no_cross_profile_copy_paste` | منع النسخ بين profiles | 5.0 |

### 2.6 إعدادات النظام

| # | القيد | القيمة | الوصف | Android Min |
|---|-------|--------|-------|-------------|
| 81 | منع تغيير التاريخ/الوقت | `no_config_date_time` | قفل الوقت والتاريخ | 9.0 |
| 82 | منع تغيير الخلفية | `no_set_wallpaper` | قفل خلفية الشاشة | 5.0 |
| 83 | منع تغيير اللغة | `no_config_locale` | قفل لغة النظام | 9.0 |
| 84 | منع تغيير السطوع | `no_config_brightness` | قفل إعدادات السطوع | 9.0 |
| 85 | منع تغيير Screen Timeout | `no_config_screen_timeout` | قفل وقت إطفاء الشاشة | 9.0 |
| 86 | منع تغيير نغمة الرنين | `no_ambient_display` | منع تغيير Ambient Display | 9.0 |

### 2.7 محتوى ووسائط

| # | القيد | القيمة | الوصف | Android Min |
|---|-------|--------|-------|-------------|
| 87 | منع Easter Eggs | `no_fun` | إخفاء Easter eggs من النظام | 5.0 |
| 88 | منع طباعة | `no_printing` | تعطيل الطباعة | 9.0 |
| 89 | منع تشغيل في الخلفية | `no_run_in_background` | منع تطبيقات الخلفية | 5.0 |
| 90 | منع تعديل الصوت | `no_adjust_volume` | منع تغيير مستوى الصوت | 5.0 |
| 91 | منع Microphone | `no_unmute_microphone` | كتم الميكروفون دائماً | 5.0 |
| 92 | منع Content Capture | `no_content_capture` | منع التقاط المحتوى | 10.0 |
| 93 | منع Content Suggestions | `no_content_suggestions` | منع اقتراحات المحتوى | 10.0 |

---

## القسم الثالث: سياسات مركّبة (أمثلة عملية)

### مثال 1: جهاز مبيعات (Sales Kiosk)

```json
{
  "kioskMode": true,
  "kioskExit": "password",
  "kioskScreenOn": true,
  "kioskLockButtons": true,
  "orientation": 2,
  "lockVolume": true,
  "volume": 70,
  "restrictions": "no_install_apps,no_factory_reset,no_usb_file_transfer,no_camera,no_safe_boot"
}
```

### مثال 2: جهاز موظف (Corporate Device)

```json
{
  "permissive": false,
  "passwordMode": "strong",
  "disableScreenshots": true,
  "lockSafeSettings": true,
  "systemUpdateType": 2,
  "systemUpdateFrom": "02:00",
  "systemUpdateTo": "05:00",
  "restrictions": "no_install_unknown_sources,no_factory_reset,no_usb_file_transfer,no_modify_accounts,no_config_vpn,no_debugging_features"
}
```

### مثال 3: جهاز طفل (Parental Control)

```json
{
  "permissive": false,
  "manageTimeout": true,
  "timeout": 1800,
  "lockVolume": true,
  "volume": 50,
  "restrictions": "no_install_apps,no_uninstall_apps,no_camera,no_sms,no_outgoing_calls,no_modify_accounts,no_config_wifi,no_install_unknown_sources"
}
```

### مثال 4: جهاز مستودع (Warehouse Scanner)

```json
{
  "kioskMode": true,
  "kioskExit": "password",
  "gps": true,
  "requestUpdates": "5",
  "orientation": 1,
  "restrictions": "no_install_apps,no_factory_reset,no_config_wifi,no_config_mobile_networks,no_modify_accounts,no_safe_boot,no_usb_file_transfer"
}
```

### مثال 5: جهاز عرض (Digital Signage)

```json
{
  "kioskMode": true,
  "kioskScreenOn": true,
  "kioskLockButtons": true,
  "kioskExit": "none",
  "orientation": 2,
  "autoBrightness": false,
  "brightness": 255,
  "lockVolume": true,
  "volume": 0,
  "restrictions": "no_install_apps,no_factory_reset,no_safe_boot,no_shutdown,no_camera,no_usb_file_transfer,no_modify_accounts,no_config_wifi,no_adjust_volume"
}
```

---

## القسم الرابع: كيفية التطبيق

### الخطوة 1: من لوحة الإدارة (Admin Panel)

1. اذهب إلى **Profiles** → اختر البروفايل → **MDM tab**
2. أضف القيود في حقل `restrictions` (مفصولة بفواصل)
3. عدّل باقي الحقول (kiosk, volume, brightness, etc.)
4. **Publish** البروفايل

### الخطوة 2: المزامنة التلقائية

- الجهاز يستقبل الإعدادات الجديدة في الـ sync التالي (خلال دقيقة عادةً)
- أو أرسل Push notification لفرض sync فوري

### الخطوة 3: التحقق

- افتح صفحة الجهاز في لوحة الإدارة
- تحقق من Device Info → MDM Mode = true
- تحقق أن القيود مطبّقة على الجهاز

---

## القسم الخامس: ملاحظات مهمة

| الملاحظة | التفاصيل |
|----------|---------|
| **Device Owner مطلوب** | معظم القيود تعمل فقط إذا التطبيق مثبّت كـ Device Owner |
| **Android version** | بعض القيود تحتاج Android 9+ (مذكور في العمود) |
| **قابلة للإلغاء** | إزالة القيد من `restrictions` يلغيه في الـ sync التالي |
| **لا تحتاج إعادة تشغيل** | معظم القيود تُطبّق فوراً بدون reboot |
| **تراكمية** | يمكن دمج أي عدد من القيود معاً |
| **لا تعارض** | settingsjson و restrictions يعملان معاً بدون تعارض |

---

## الإحصائيات

| الفئة | العدد |
|-------|-------|
| سياسات settingsjson | 51 |
| قيود restrictions | 42 |
| **المجموع الكلي** | **93 سياسة/قيد** |

جميعها تعمل **بدون أي تعديل على كود APK** — فقط تغيير الإعدادات من لوحة الإدارة.
