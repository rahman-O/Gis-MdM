-- Phase 4: devices, devicegroups, tenant-scoped groups/configurations, permissions seed.

ALTER TABLE groups ADD COLUMN IF NOT EXISTS customerid INT REFERENCES customers (id) ON DELETE CASCADE;
UPDATE groups SET customerid = 1 WHERE customerid IS NULL;
ALTER TABLE groups ALTER COLUMN customerid SET NOT NULL;

ALTER TABLE configurations ADD COLUMN IF NOT EXISTS customerid INT REFERENCES customers (id) ON DELETE CASCADE;
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS permissive BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS mainappid INT;
UPDATE configurations SET customerid = 1 WHERE customerid IS NULL;
ALTER TABLE configurations ALTER COLUMN customerid SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS groups_name_customer_uidx ON groups (customerid, lower(name));

CREATE TABLE IF NOT EXISTS devices (
    id              SERIAL PRIMARY KEY,
    number          VARCHAR(100) NOT NULL,
    description     TEXT,
    lastupdate      BIGINT NOT NULL DEFAULT 0,
    configurationid INT NOT NULL REFERENCES configurations (id),
    customerid      INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    info            TEXT,
    infojson        JSONB,
    imei            VARCHAR(50),
    phone           VARCHAR(50),
    enrolltime      BIGINT NOT NULL DEFAULT 0,
    publicip        VARCHAR(50),
    custom1         VARCHAR(200),
    custom2         VARCHAR(200),
    custom3         VARCHAR(200),
    oldnumber       VARCHAR(100),
    fastsearch      VARCHAR(200)
);

CREATE UNIQUE INDEX IF NOT EXISTS devices_number_customer_uidx ON devices (customerid, lower(number));

CREATE TABLE IF NOT EXISTS devicegroups (
    deviceid INT NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    groupid  INT NOT NULL REFERENCES groups (id) ON DELETE CASCADE,
    PRIMARY KEY (deviceid, groupid)
);

CREATE TABLE IF NOT EXISTS deviceapplicationsettings (
    id              SERIAL PRIMARY KEY,
    deviceid        INT NOT NULL REFERENCES devices (id) ON DELETE CASCADE,
    applicationpkg  VARCHAR(200) NOT NULL,
    name            VARCHAR(200),
    type            VARCHAR(50),
    value           TEXT
);

-- 000001 seeds permission id=1 without bumping the sequence; align before new rows.
SELECT setval('permissions_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM permissions), 1));

INSERT INTO permissions (name, description, superadmin)
SELECT 'edit_devices', 'Create and edit devices', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'edit_devices');

INSERT INTO permissions (name, description, superadmin)
SELECT 'edit_device_desc', 'Edit device descriptions', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'edit_device_desc');

INSERT INTO userrolepermissions (roleid, permissionid)
SELECT 2, p.id
FROM permissions p
WHERE lower(p.name) IN ('edit_devices', 'edit_device_desc', 'settings')
  AND NOT EXISTS (
      SELECT 1 FROM userrolepermissions urp
      WHERE urp.roleid = 2 AND urp.permissionid = p.id
  );

INSERT INTO groups (name, customerid)
SELECT 'General', 1
WHERE NOT EXISTS (SELECT 1 FROM groups WHERE customerid = 1 AND lower(name) = 'general');

INSERT INTO configurations (name, description, customerid, permissive)
SELECT 'Default', 'Default device configuration', 1, FALSE
WHERE NOT EXISTS (SELECT 1 FROM configurations WHERE customerid = 1 AND lower(name) = 'default');

INSERT INTO devices (number, description, lastupdate, configurationid, customerid, enrolltime)
SELECT 'hmdm-001', 'Sample device 1', EXTRACT(EPOCH FROM NOW())::BIGINT * 1000,
       (SELECT id FROM configurations WHERE customerid = 1 ORDER BY id LIMIT 1),
       1, EXTRACT(EPOCH FROM NOW())::BIGINT * 1000
WHERE NOT EXISTS (SELECT 1 FROM devices WHERE customerid = 1 AND number = 'hmdm-001');

INSERT INTO devices (number, description, lastupdate, configurationid, customerid, enrolltime)
SELECT 'hmdm-002', 'Sample device 2', EXTRACT(EPOCH FROM NOW())::BIGINT * 1000 - 7200000,
       (SELECT id FROM configurations WHERE customerid = 1 ORDER BY id LIMIT 1),
       1, EXTRACT(EPOCH FROM NOW())::BIGINT * 1000
WHERE NOT EXISTS (SELECT 1 FROM devices WHERE customerid = 1 AND number = 'hmdm-002');

INSERT INTO devicegroups (deviceid, groupid)
SELECT d.id, g.id
FROM devices d
CROSS JOIN groups g
WHERE d.customerid = 1 AND d.number IN ('hmdm-001', 'hmdm-002')
  AND g.customerid = 1 AND lower(g.name) = 'general'
ON CONFLICT DO NOTHING;
