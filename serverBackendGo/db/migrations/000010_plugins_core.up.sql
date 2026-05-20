-- Phase 8: plugin platform catalog, extension tables, permissions.

CREATE TABLE IF NOT EXISTS plugins (
    id                      SERIAL PRIMARY KEY,
    identifier              VARCHAR(50) NOT NULL UNIQUE,
    name                    TEXT NOT NULL,
    description             TEXT,
    createtime              TIMESTAMP NOT NULL DEFAULT NOW(),
    disabled                BOOLEAN NOT NULL DEFAULT FALSE,
    javascriptmodulefile    VARCHAR(200),
    functionsviewtemplate   VARCHAR(200),
    settingsviewtemplate    VARCHAR(200),
    namelocalizationkey     VARCHAR(200) NOT NULL DEFAULT 'plugin.name.not.specified',
    settingspermission      VARCHAR(200),
    functionspermission     VARCHAR(200),
    devicefunctionspermission VARCHAR(200),
    enabledfordevice        BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS pluginsdisabled (
    pluginid    INT NOT NULL REFERENCES plugins (id) ON DELETE CASCADE,
    customerid  INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    PRIMARY KEY (pluginid, customerid)
);

CREATE TABLE IF NOT EXISTS plugin_audit_log (
    id          SERIAL PRIMARY KEY,
    createtime  BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    customerid  INT REFERENCES customers (id) ON DELETE CASCADE,
    userid      INT,
    login       VARCHAR(100),
    action      VARCHAR(100),
    payload     TEXT,
    ipaddress   VARCHAR(500),
    errorcode   INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS plugin_messaging_messages (
    id          SERIAL PRIMARY KEY,
    customerid  INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    deviceid    INT NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    ts          BIGINT NOT NULL,
    message     VARCHAR(5000),
    status      INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS plugin_deviceinfo_settings (
    id                  SERIAL PRIMARY KEY,
    customerid          INT NOT NULL UNIQUE REFERENCES customers (id) ON DELETE CASCADE,
    datapreserveperiod  INT NOT NULL DEFAULT 30,
    senddata            BOOLEAN NOT NULL DEFAULT FALSE,
    intervalmins        INT NOT NULL DEFAULT 15
);

CREATE TABLE IF NOT EXISTS plugin_deviceinfo_deviceparams (
    id          SERIAL PRIMARY KEY,
    deviceid    INT NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    customerid  INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    ts          BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS plugin_deviceinfo_deviceparams_device (
    id              SERIAL PRIMARY KEY,
    recordid        INT NOT NULL UNIQUE REFERENCES plugin_deviceinfo_deviceparams (id) ON DELETE CASCADE,
    batterylevel    INT,
    batterycharging VARCHAR(20),
    ip              VARCHAR(50),
    wifienabled     BOOLEAN,
    gpsenabled      BOOLEAN
);

CREATE TABLE IF NOT EXISTS plugin_devicelog_settings (
    id                  SERIAL PRIMARY KEY,
    customerid          INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    logspreserveperiod  INT NOT NULL DEFAULT 30
);

CREATE UNIQUE INDEX IF NOT EXISTS plugin_devicelog_settings_customer_uq
    ON plugin_devicelog_settings (customerid);

CREATE TABLE IF NOT EXISTS plugin_devicelog_settings_rules (
    id                SERIAL PRIMARY KEY,
    settingid         INT NOT NULL REFERENCES plugin_devicelog_settings (id) ON DELETE CASCADE,
    name              VARCHAR(120) NOT NULL,
    active            BOOLEAN NOT NULL DEFAULT TRUE,
    applicationid     INT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    severity          TEXT NOT NULL,
    filter            TEXT,
    groupid           INT REFERENCES groups (id) ON DELETE CASCADE,
    configurationid   INT REFERENCES configurations (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS plugin_devicelog_setting_rule_devices (
    ruleid    INT NOT NULL REFERENCES plugin_devicelog_settings_rules (id) ON DELETE CASCADE,
    deviceid  INT NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    PRIMARY KEY (ruleid, deviceid)
);

CREATE TABLE IF NOT EXISTS plugin_devicelog_log (
    id              SERIAL PRIMARY KEY,
    createtime      BIGINT,
    customerid      INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    deviceid        INT NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    applicationid   INT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    ipaddress       VARCHAR(512),
    severity        TEXT,
    severityorder   INT,
    message         TEXT
);

-- Catalog seeds (open-source plugins)
INSERT INTO plugins (identifier, name, description, namelocalizationkey, enabledfordevice)
SELECT 'audit', 'Audit', 'User action audit', 'plugin.audit.localization.key.name', FALSE
WHERE NOT EXISTS (SELECT 1 FROM plugins WHERE identifier = 'audit');

INSERT INTO plugins (identifier, name, description, namelocalizationkey, enabledfordevice)
SELECT 'push', 'Push', 'Push messaging plugin', 'plugin.push.localization.key.name', FALSE
WHERE NOT EXISTS (SELECT 1 FROM plugins WHERE identifier = 'push');

INSERT INTO plugins (identifier, name, description, namelocalizationkey, enabledfordevice)
SELECT 'messaging', 'Messaging', 'Device messaging', 'plugin.messaging.localization.key.name', TRUE
WHERE NOT EXISTS (SELECT 1 FROM plugins WHERE identifier = 'messaging');

INSERT INTO plugins (identifier, name, description, namelocalizationkey, enabledfordevice)
SELECT 'deviceinfo', 'Device Info', 'Device telemetry', 'plugin.deviceinfo.localization.key.name', TRUE
WHERE NOT EXISTS (SELECT 1 FROM plugins WHERE identifier = 'deviceinfo');

INSERT INTO plugins (identifier, name, description, namelocalizationkey, enabledfordevice)
SELECT 'devicelog', 'Device Log', 'Device logs', 'plugin.devicelog.localization.key.name', TRUE
WHERE NOT EXISTS (SELECT 1 FROM plugins WHERE identifier = 'devicelog');

INSERT INTO permissions (name, description, superadmin)
SELECT 'plugins_customer_access_management', 'Manage tenant plugin list', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'plugins_customer_access_management');

INSERT INTO permissions (name, description, superadmin)
SELECT 'plugin_audit_access', 'Access audit log', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'plugin_audit_access');

INSERT INTO permissions (name, description, superadmin)
SELECT 'plugin_messaging_send', 'Send messaging to devices', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'plugin_messaging_send');

INSERT INTO permissions (name, description, superadmin)
SELECT 'plugin_messaging_delete', 'Delete messaging history', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'plugin_messaging_delete');

INSERT INTO permissions (name, description, superadmin)
SELECT 'plugin_deviceinfo_access', 'Access device info plugin', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'plugin_deviceinfo_access');

INSERT INTO permissions (name, description, superadmin)
SELECT 'plugin_devicelog_access', 'Access device log plugin', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'plugin_devicelog_access');

INSERT INTO userrolepermissions (roleid, permissionid)
SELECT 2, p.id FROM permissions p
WHERE lower(p.name) IN (
    'plugins_customer_access_management', 'plugin_audit_access',
    'plugin_messaging_send', 'plugin_messaging_delete',
    'plugin_deviceinfo_access', 'plugin_devicelog_access'
)
AND NOT EXISTS (SELECT 1 FROM userrolepermissions urp WHERE urp.roleid = 2 AND urp.permissionid = p.id);

INSERT INTO plugin_devicelog_settings (customerid)
SELECT c.id FROM customers c
WHERE NOT EXISTS (SELECT 1 FROM plugin_devicelog_settings s WHERE s.customerid = c.id);

INSERT INTO plugin_deviceinfo_settings (customerid)
SELECT c.id FROM customers c
WHERE NOT EXISTS (SELECT 1 FROM plugin_deviceinfo_settings s WHERE s.customerid = c.id);

INSERT INTO plugin_audit_log (createtime, customerid, userid, login, action, payload, ipaddress)
SELECT EXTRACT(EPOCH FROM NOW())::BIGINT * 1000, 1, 1, 'admin', 'login', '{"source":"seed"}', '127.0.0.1'
WHERE EXISTS (SELECT 1 FROM customers WHERE id = 1)
  AND NOT EXISTS (SELECT 1 FROM plugin_audit_log WHERE customerid = 1 AND action = 'login' AND login = 'admin');
