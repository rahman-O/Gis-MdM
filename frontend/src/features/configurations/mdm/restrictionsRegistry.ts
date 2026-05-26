/**
 * @file restrictionsRegistry.ts
 * @description Registry of all 42 Android UserRestrictions supported by the MDM agent.
 * Each restriction maps to an Android UserManager constant (e.g. DISALLOW_CAMERA → "no_camera").
 * The registry provides bilingual labels (EN/AR), category grouping, and minimum Android version.
 */

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

/** Restriction category used for grouping in the UI. */
export type RestrictionCategory =
  | 'security'
  | 'apps'
  | 'network'
  | 'accounts'
  | 'files'
  | 'system'
  | 'media';

/** Definition of a single Android UserRestriction. */
export interface RestrictionDefinition {
  /** Android restriction key (e.g. "no_camera"). */
  key: string;
  /** UI grouping category. */
  category: RestrictionCategory;
  /** English display label. */
  labelEn: string;
  /** Arabic display label. */
  labelAr: string;
  /** English description of what the restriction does. */
  descriptionEn: string;
  /** Arabic description of what the restriction does. */
  descriptionAr: string;
  /** Minimum Android version required (e.g. "9.0", "8.0", "5.0"). */
  minAndroid: string;
}

// ---------------------------------------------------------------------------
// Registry (42 restrictions)
// ---------------------------------------------------------------------------

/** Complete registry of all supported Android UserRestrictions. */
export const RESTRICTIONS_REGISTRY: RestrictionDefinition[] = [
  // ─── Security (5) ────────────────────────────────────────────────────────────
  {
    key: 'no_camera',
    category: 'security',
    labelEn: 'Disable Camera',
    labelAr: 'تعطيل الكاميرا',
    descriptionEn: 'Prevents use of all device cameras.',
    descriptionAr: 'يمنع استخدام جميع كاميرات الجهاز.',
    minAndroid: '5.0',
  },
  {
    key: 'no_factory_reset',
    category: 'security',
    labelEn: 'Disable Factory Reset',
    labelAr: 'تعطيل إعادة ضبط المصنع',
    descriptionEn: 'Prevents the user from performing a factory reset.',
    descriptionAr: 'يمنع المستخدم من إجراء إعادة ضبط المصنع.',
    minAndroid: '5.0',
  },
  {
    key: 'no_safe_boot',
    category: 'security',
    labelEn: 'Disable Safe Boot',
    labelAr: 'تعطيل التشغيل الآمن',
    descriptionEn: 'Prevents the user from booting into safe mode.',
    descriptionAr: 'يمنع المستخدم من التشغيل في الوضع الآمن.',
    minAndroid: '5.0',
  },
  {
    key: 'no_debugging_features',
    category: 'security',
    labelEn: 'Disable Debugging',
    labelAr: 'تعطيل التصحيح',
    descriptionEn: 'Prevents enabling developer options and USB debugging.',
    descriptionAr: 'يمنع تفعيل خيارات المطور وتصحيح USB.',
    minAndroid: '5.0',
  },
  {
    key: 'no_shutdown',
    category: 'security',
    labelEn: 'Disable Shutdown',
    labelAr: 'تعطيل إيقاف التشغيل',
    descriptionEn: 'Prevents the user from shutting down the device.',
    descriptionAr: 'يمنع المستخدم من إيقاف تشغيل الجهاز.',
    minAndroid: '9.0',
  },

  // ─── Apps (4) ────────────────────────────────────────────────────────────────
  {
    key: 'no_install_apps',
    category: 'apps',
    labelEn: 'Disable App Installation',
    labelAr: 'تعطيل تثبيت التطبيقات',
    descriptionEn: 'Prevents the user from installing new applications.',
    descriptionAr: 'يمنع المستخدم من تثبيت تطبيقات جديدة.',
    minAndroid: '5.0',
  },
  {
    key: 'no_uninstall_apps',
    category: 'apps',
    labelEn: 'Disable App Uninstallation',
    labelAr: 'تعطيل إزالة التطبيقات',
    descriptionEn: 'Prevents the user from uninstalling applications.',
    descriptionAr: 'يمنع المستخدم من إزالة التطبيقات.',
    minAndroid: '5.0',
  },
  {
    key: 'no_install_unknown_sources',
    category: 'apps',
    labelEn: 'Disable Unknown Sources',
    labelAr: 'تعطيل المصادر غير المعروفة',
    descriptionEn: 'Prevents installing apps from unknown sources.',
    descriptionAr: 'يمنع تثبيت التطبيقات من مصادر غير معروفة.',
    minAndroid: '5.0',
  },
  {
    key: 'no_install_unknown_sources_globally',
    category: 'apps',
    labelEn: 'Disable Unknown Sources Globally',
    labelAr: 'تعطيل المصادر غير المعروفة عالمياً',
    descriptionEn: 'Prevents all users from installing apps from unknown sources.',
    descriptionAr: 'يمنع جميع المستخدمين من تثبيت التطبيقات من مصادر غير معروفة.',
    minAndroid: '8.0',
  },

  // ─── Network (12) ────────────────────────────────────────────────────────────
  {
    key: 'no_config_wifi',
    category: 'network',
    labelEn: 'Disable Wi-Fi Configuration',
    labelAr: 'تعطيل إعدادات الواي فاي',
    descriptionEn: 'Prevents the user from changing Wi-Fi settings.',
    descriptionAr: 'يمنع المستخدم من تغيير إعدادات الواي فاي.',
    minAndroid: '5.0',
  },
  {
    key: 'no_config_vpn',
    category: 'network',
    labelEn: 'Disable VPN Configuration',
    labelAr: 'تعطيل إعدادات VPN',
    descriptionEn: 'Prevents the user from configuring VPN connections.',
    descriptionAr: 'يمنع المستخدم من إعداد اتصالات VPN.',
    minAndroid: '5.0',
  },
  {
    key: 'no_config_tethering',
    category: 'network',
    labelEn: 'Disable Tethering Configuration',
    labelAr: 'تعطيل إعدادات الربط',
    descriptionEn: 'Prevents the user from configuring tethering and hotspots.',
    descriptionAr: 'يمنع المستخدم من إعداد الربط ونقاط الاتصال.',
    minAndroid: '5.0',
  },
  {
    key: 'no_bluetooth_sharing',
    category: 'network',
    labelEn: 'Disable Bluetooth Sharing',
    labelAr: 'تعطيل مشاركة البلوتوث',
    descriptionEn: 'Prevents sharing files via Bluetooth.',
    descriptionAr: 'يمنع مشاركة الملفات عبر البلوتوث.',
    minAndroid: '5.0',
  },
  {
    key: 'no_bluetooth',
    category: 'network',
    labelEn: 'Disable Bluetooth',
    labelAr: 'تعطيل البلوتوث',
    descriptionEn: 'Prevents the user from enabling Bluetooth.',
    descriptionAr: 'يمنع المستخدم من تفعيل البلوتوث.',
    minAndroid: '5.0',
  },
  {
    key: 'no_sms',
    category: 'network',
    labelEn: 'Disable SMS',
    labelAr: 'تعطيل الرسائل القصيرة',
    descriptionEn: 'Prevents the user from sending or receiving SMS messages.',
    descriptionAr: 'يمنع المستخدم من إرسال أو استقبال الرسائل القصيرة.',
    minAndroid: '5.0',
  },
  {
    key: 'no_outgoing_calls',
    category: 'network',
    labelEn: 'Disable Outgoing Calls',
    labelAr: 'تعطيل المكالمات الصادرة',
    descriptionEn: 'Prevents the user from making outgoing phone calls.',
    descriptionAr: 'يمنع المستخدم من إجراء مكالمات هاتفية صادرة.',
    minAndroid: '5.0',
  },
  {
    key: 'no_config_mobile_networks',
    category: 'network',
    labelEn: 'Disable Mobile Network Configuration',
    labelAr: 'تعطيل إعدادات شبكة الجوال',
    descriptionEn: 'Prevents the user from changing mobile network settings.',
    descriptionAr: 'يمنع المستخدم من تغيير إعدادات شبكة الجوال.',
    minAndroid: '5.0',
  },
  {
    key: 'no_outgoing_beam',
    category: 'network',
    labelEn: 'Disable NFC Beam',
    labelAr: 'تعطيل شعاع NFC',
    descriptionEn: 'Prevents the user from using NFC beam to share data.',
    descriptionAr: 'يمنع المستخدم من استخدام شعاع NFC لمشاركة البيانات.',
    minAndroid: '5.0',
  },
  {
    key: 'no_airplane_mode',
    category: 'network',
    labelEn: 'Disable Airplane Mode',
    labelAr: 'تعطيل وضع الطيران',
    descriptionEn: 'Prevents the user from toggling airplane mode.',
    descriptionAr: 'يمنع المستخدم من تبديل وضع الطيران.',
    minAndroid: '9.0',
  },
  {
    key: 'no_config_cell_broadcasts',
    category: 'network',
    labelEn: 'Disable Cell Broadcast Configuration',
    labelAr: 'تعطيل إعدادات البث الخلوي',
    descriptionEn: 'Prevents the user from configuring cell broadcast channels.',
    descriptionAr: 'يمنع المستخدم من إعداد قنوات البث الخلوي.',
    minAndroid: '5.0',
  },
  {
    key: 'no_data_roaming',
    category: 'network',
    labelEn: 'Disable Data Roaming',
    labelAr: 'تعطيل تجوال البيانات',
    descriptionEn: 'Prevents the user from enabling data roaming.',
    descriptionAr: 'يمنع المستخدم من تفعيل تجوال البيانات.',
    minAndroid: '5.0',
  },

  // ─── Accounts (4) ────────────────────────────────────────────────────────────
  {
    key: 'no_modify_accounts',
    category: 'accounts',
    labelEn: 'Disable Account Modification',
    labelAr: 'تعطيل تعديل الحسابات',
    descriptionEn: 'Prevents the user from adding or removing accounts.',
    descriptionAr: 'يمنع المستخدم من إضافة أو إزالة الحسابات.',
    minAndroid: '5.0',
  },
  {
    key: 'no_add_user',
    category: 'accounts',
    labelEn: 'Disable Adding Users',
    labelAr: 'تعطيل إضافة المستخدمين',
    descriptionEn: 'Prevents the user from adding new device users.',
    descriptionAr: 'يمنع المستخدم من إضافة مستخدمين جدد للجهاز.',
    minAndroid: '5.0',
  },
  {
    key: 'no_remove_user',
    category: 'accounts',
    labelEn: 'Disable Removing Users',
    labelAr: 'تعطيل إزالة المستخدمين',
    descriptionEn: 'Prevents the user from removing device users.',
    descriptionAr: 'يمنع المستخدم من إزالة مستخدمي الجهاز.',
    minAndroid: '5.0',
  },
  {
    key: 'no_switch_user',
    category: 'accounts',
    labelEn: 'Disable User Switching',
    labelAr: 'تعطيل تبديل المستخدمين',
    descriptionEn: 'Prevents the user from switching between device users.',
    descriptionAr: 'يمنع المستخدم من التبديل بين مستخدمي الجهاز.',
    minAndroid: '9.0',
  },

  // ─── Files (4) ───────────────────────────────────────────────────────────────
  {
    key: 'no_usb_file_transfer',
    category: 'files',
    labelEn: 'Disable USB File Transfer',
    labelAr: 'تعطيل نقل الملفات عبر USB',
    descriptionEn: 'Prevents file transfer over USB connections.',
    descriptionAr: 'يمنع نقل الملفات عبر اتصالات USB.',
    minAndroid: '5.0',
  },
  {
    key: 'no_physical_media',
    category: 'files',
    labelEn: 'Disable Physical Media',
    labelAr: 'تعطيل الوسائط المادية',
    descriptionEn: 'Prevents mounting of external physical media (SD cards).',
    descriptionAr: 'يمنع تركيب الوسائط المادية الخارجية (بطاقات SD).',
    minAndroid: '5.0',
  },
  {
    key: 'no_share_location',
    category: 'files',
    labelEn: 'Disable Location Sharing',
    labelAr: 'تعطيل مشاركة الموقع',
    descriptionEn: 'Prevents the user from sharing their location.',
    descriptionAr: 'يمنع المستخدم من مشاركة موقعه.',
    minAndroid: '5.0',
  },
  {
    key: 'no_cross_profile_copy_paste',
    category: 'files',
    labelEn: 'Disable Cross-Profile Copy/Paste',
    labelAr: 'تعطيل النسخ واللصق بين الملفات الشخصية',
    descriptionEn: 'Prevents copy-paste between work and personal profiles.',
    descriptionAr: 'يمنع النسخ واللصق بين ملفات العمل والملفات الشخصية.',
    minAndroid: '5.0',
  },

  // ─── System (6) ──────────────────────────────────────────────────────────────
  {
    key: 'no_config_date_time',
    category: 'system',
    labelEn: 'Disable Date/Time Configuration',
    labelAr: 'تعطيل إعدادات التاريخ والوقت',
    descriptionEn: 'Prevents the user from changing date and time settings.',
    descriptionAr: 'يمنع المستخدم من تغيير إعدادات التاريخ والوقت.',
    minAndroid: '9.0',
  },
  {
    key: 'no_set_wallpaper',
    category: 'system',
    labelEn: 'Disable Wallpaper Change',
    labelAr: 'تعطيل تغيير الخلفية',
    descriptionEn: 'Prevents the user from changing the device wallpaper.',
    descriptionAr: 'يمنع المستخدم من تغيير خلفية الجهاز.',
    minAndroid: '5.0',
  },
  {
    key: 'no_config_locale',
    category: 'system',
    labelEn: 'Disable Locale Configuration',
    labelAr: 'تعطيل إعدادات اللغة',
    descriptionEn: 'Prevents the user from changing the device language.',
    descriptionAr: 'يمنع المستخدم من تغيير لغة الجهاز.',
    minAndroid: '9.0',
  },
  {
    key: 'no_config_brightness',
    category: 'system',
    labelEn: 'Disable Brightness Configuration',
    labelAr: 'تعطيل إعدادات السطوع',
    descriptionEn: 'Prevents the user from changing screen brightness.',
    descriptionAr: 'يمنع المستخدم من تغيير سطوع الشاشة.',
    minAndroid: '9.0',
  },
  {
    key: 'no_config_screen_timeout',
    category: 'system',
    labelEn: 'Disable Screen Timeout Configuration',
    labelAr: 'تعطيل إعدادات مهلة الشاشة',
    descriptionEn: 'Prevents the user from changing the screen timeout duration.',
    descriptionAr: 'يمنع المستخدم من تغيير مدة مهلة الشاشة.',
    minAndroid: '9.0',
  },
  {
    key: 'no_ambient_display',
    category: 'system',
    labelEn: 'Disable Ambient Display',
    labelAr: 'تعطيل العرض المحيط',
    descriptionEn: 'Prevents the ambient display (always-on display) feature.',
    descriptionAr: 'يمنع ميزة العرض المحيط (الشاشة الدائمة).',
    minAndroid: '9.0',
  },

  // ─── Media (7) ───────────────────────────────────────────────────────────────
  {
    key: 'no_fun',
    category: 'media',
    labelEn: 'Disable Fun',
    labelAr: 'تعطيل الترفيه',
    descriptionEn: 'Disables easter eggs and fun features on the device.',
    descriptionAr: 'يعطل ميزات الترفيه والمفاجآت في الجهاز.',
    minAndroid: '5.0',
  },
  {
    key: 'no_printing',
    category: 'media',
    labelEn: 'Disable Printing',
    labelAr: 'تعطيل الطباعة',
    descriptionEn: 'Prevents the user from printing documents.',
    descriptionAr: 'يمنع المستخدم من طباعة المستندات.',
    minAndroid: '9.0',
  },
  {
    key: 'no_run_in_background',
    category: 'media',
    labelEn: 'Disable Background Processes',
    labelAr: 'تعطيل العمليات في الخلفية',
    descriptionEn: 'Prevents apps from running in the background.',
    descriptionAr: 'يمنع التطبيقات من العمل في الخلفية.',
    minAndroid: '5.0',
  },
  {
    key: 'no_adjust_volume',
    category: 'media',
    labelEn: 'Disable Volume Adjustment',
    labelAr: 'تعطيل ضبط مستوى الصوت',
    descriptionEn: 'Prevents the user from adjusting the device volume.',
    descriptionAr: 'يمنع المستخدم من ضبط مستوى صوت الجهاز.',
    minAndroid: '5.0',
  },
  {
    key: 'no_unmute_microphone',
    category: 'media',
    labelEn: 'Disable Microphone Unmute',
    labelAr: 'تعطيل إلغاء كتم الميكروفون',
    descriptionEn: 'Prevents the user from unmuting the microphone.',
    descriptionAr: 'يمنع المستخدم من إلغاء كتم صوت الميكروفون.',
    minAndroid: '5.0',
  },
  {
    key: 'no_content_capture',
    category: 'media',
    labelEn: 'Disable Content Capture',
    labelAr: 'تعطيل التقاط المحتوى',
    descriptionEn: 'Prevents content capture services from analyzing screen content.',
    descriptionAr: 'يمنع خدمات التقاط المحتوى من تحليل محتوى الشاشة.',
    minAndroid: '10.0',
  },
  {
    key: 'no_content_suggestions',
    category: 'media',
    labelEn: 'Disable Content Suggestions',
    labelAr: 'تعطيل اقتراحات المحتوى',
    descriptionEn: 'Prevents the system from showing content suggestions.',
    descriptionAr: 'يمنع النظام من عرض اقتراحات المحتوى.',
    minAndroid: '10.0',
  },
];

// ---------------------------------------------------------------------------
// Category labels
// ---------------------------------------------------------------------------

/** Bilingual labels for each restriction category. */
export const CATEGORY_LABELS: Record<RestrictionCategory, { en: string; ar: string }> = {
  security: { en: 'Security', ar: 'الأمان' },
  apps: { en: 'Applications', ar: 'التطبيقات' },
  network: { en: 'Network & Communication', ar: 'الشبكة والاتصالات' },
  accounts: { en: 'Accounts & Users', ar: 'الحسابات والمستخدمون' },
  files: { en: 'Files & Sharing', ar: 'الملفات والمشاركة' },
  system: { en: 'System', ar: 'النظام' },
  media: { en: 'Media & Content', ar: 'الوسائط والمحتوى' },
};

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

/**
 * Parses a comma-separated restrictions string into a Set of restriction keys.
 * Returns an empty set for null/undefined/empty input.
 *
 * @example
 * restrictionsToSet("no_camera,no_sms") // Set {"no_camera", "no_sms"}
 * restrictionsToSet(null)               // Set {}
 */
export function restrictionsToSet(restrictions: string | null | undefined): Set<string> {
  if (!restrictions || restrictions.trim() === '') {
    return new Set<string>();
  }
  return new Set(
    restrictions
      .split(',')
      .map((r) => r.trim())
      .filter((r) => r.length > 0),
  );
}

/**
 * Serializes a Set of restriction keys back into a comma-separated string.
 * Returns an empty string for an empty set.
 *
 * @example
 * setToRestrictions(new Set(["no_camera", "no_sms"])) // "no_camera,no_sms"
 * setToRestrictions(new Set())                         // ""
 */
export function setToRestrictions(set: Set<string>): string {
  return Array.from(set).join(',');
}

/**
 * Returns all restriction definitions belonging to a given category.
 *
 * @example
 * getByCategory('security') // [{ key: 'no_camera', ... }, ...]
 */
export function getByCategory(category: RestrictionCategory): RestrictionDefinition[] {
  return RESTRICTIONS_REGISTRY.filter((r) => r.category === category);
}
