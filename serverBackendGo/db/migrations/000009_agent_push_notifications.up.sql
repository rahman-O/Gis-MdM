-- Phase 7: push queue, plugin push tables, permissions.

CREATE TABLE IF NOT EXISTS pushmessages (
    id          SERIAL PRIMARY KEY,
    messagetype VARCHAR(50) NOT NULL,
    deviceid    INT NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    payload     TEXT
);

CREATE TABLE IF NOT EXISTS pendingpushes (
    id          SERIAL PRIMARY KEY,
    messageid   INT NOT NULL UNIQUE REFERENCES pushmessages (id) ON DELETE CASCADE,
    status      INT NOT NULL DEFAULT 0,
    createtime  BIGINT NOT NULL,
    sendtime    BIGINT
);

CREATE INDEX IF NOT EXISTS pushmessages_deviceid_idx ON pushmessages (deviceid);
CREATE INDEX IF NOT EXISTS pendingpushes_status_idx ON pendingpushes (status);

CREATE TABLE IF NOT EXISTS plugin_push_messages (
    id          SERIAL PRIMARY KEY,
    customerid  INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    deviceid    INT NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    ts          BIGINT NOT NULL,
    messagetype VARCHAR(255),
    payload     TEXT
);

CREATE INDEX IF NOT EXISTS plugin_push_messages_customer_idx ON plugin_push_messages (customerid);

CREATE TABLE IF NOT EXISTS plugin_push_schedule (
    id                SERIAL PRIMARY KEY,
    customerid        INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    deviceid          INT NOT NULL DEFAULT 0,
    groupid           INT NOT NULL DEFAULT 0,
    configurationid   INT NOT NULL DEFAULT 0,
    scope             VARCHAR(255),
    messagetype       VARCHAR(255),
    payload           TEXT,
    comment           TEXT,
    min               VARCHAR(1024),
    hour              VARCHAR(1024),
    day               VARCHAR(1024),
    weekday           VARCHAR(1024),
    month             VARCHAR(1024)
);

INSERT INTO permissions (name, description, superadmin)
SELECT 'push_api', 'Send Push messages to devices via REST API', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'push_api');

INSERT INTO permissions (name, description, superadmin)
SELECT 'plugin_push_send', 'Can send push messages to devices', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'plugin_push_send');

INSERT INTO permissions (name, description, superadmin)
SELECT 'plugin_push_delete', 'Can delete and update Push message history', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'plugin_push_delete');

INSERT INTO userrolepermissions (roleid, permissionid)
SELECT 2, p.id
FROM permissions p
WHERE lower(p.name) IN ('push_api', 'plugin_push_send', 'plugin_push_delete')
  AND NOT EXISTS (
      SELECT 1 FROM userrolepermissions urp
      WHERE urp.roleid = 2 AND urp.permissionid = p.id
  );

-- Smoke: queue message for hmdm-001 when device exists
INSERT INTO pushmessages (messagetype, deviceid, payload)
SELECT 'configUpdated', d.id, ''
FROM devices d
WHERE d.number = 'hmdm-001'
  AND NOT EXISTS (
      SELECT 1 FROM pushmessages pm
      JOIN pendingpushes pp ON pp.messageid = pm.id
      WHERE pm.deviceid = d.id AND pp.status = 0
  );

UPDATE configurations SET qrcodekey = 'default-qr'
WHERE customerid = 1 AND lower(name) = 'default' AND (qrcodekey IS NULL OR qrcodekey = '');

INSERT INTO pendingpushes (messageid, status, createtime)
SELECT pm.id, 0, EXTRACT(EPOCH FROM NOW())::BIGINT * 1000
FROM pushmessages pm
JOIN devices d ON d.id = pm.deviceid AND d.number = 'hmdm-001'
WHERE NOT EXISTS (
    SELECT 1 FROM pendingpushes pp WHERE pp.messageid = pm.id
);
