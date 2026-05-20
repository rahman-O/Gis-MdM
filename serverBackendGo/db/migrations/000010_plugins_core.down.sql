DELETE FROM plugin_audit_log;
DROP TABLE IF EXISTS plugin_devicelog_log;
DROP TABLE IF EXISTS plugin_devicelog_setting_rule_devices;
DROP TABLE IF EXISTS plugin_devicelog_settings_rules;
DROP TABLE IF EXISTS plugin_devicelog_settings;
DROP TABLE IF EXISTS plugin_deviceinfo_deviceparams_device;
DROP TABLE IF EXISTS plugin_deviceinfo_deviceparams;
DROP TABLE IF EXISTS plugin_deviceinfo_settings;
DROP TABLE IF EXISTS plugin_messaging_messages;
DROP TABLE IF EXISTS plugin_audit_log;
DROP TABLE IF EXISTS pluginsdisabled;

DELETE FROM userrolepermissions urp
USING permissions p
WHERE urp.permissionid = p.id
  AND lower(p.name) IN (
    'plugins_customer_access_management', 'plugin_audit_access',
    'plugin_messaging_send', 'plugin_messaging_delete',
    'plugin_deviceinfo_access', 'plugin_devicelog_access'
  );

DELETE FROM permissions WHERE lower(name) IN (
    'plugins_customer_access_management', 'plugin_audit_access',
    'plugin_messaging_send', 'plugin_messaging_delete',
    'plugin_deviceinfo_access', 'plugin_devicelog_access'
);

DELETE FROM plugins WHERE identifier IN ('audit', 'push', 'messaging', 'deviceinfo', 'devicelog');
DROP TABLE IF EXISTS plugins;
