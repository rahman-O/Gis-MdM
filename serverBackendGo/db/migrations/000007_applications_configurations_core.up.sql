-- Phase 5: applications, application versions, configuration junction tables, extended configurations.

ALTER TABLE configurations ADD COLUMN IF NOT EXISTS type INT NOT NULL DEFAULT 0;
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS password TEXT;
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS backgroundcolor VARCHAR(50);
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS textcolor VARCHAR(50);
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS backgroundimageurl TEXT;
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS qrcodekey VARCHAR(200);
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS baseurl TEXT;
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS defaultfilepath VARCHAR(500);
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS contentappid INT;
ALTER TABLE configurations ADD COLUMN IF NOT EXISTS settingsjson JSONB NOT NULL DEFAULT '{}';

CREATE UNIQUE INDEX IF NOT EXISTS configurations_name_customer_uidx
    ON configurations (customerid, lower(name));

CREATE TABLE IF NOT EXISTS applications (
    id              SERIAL PRIMARY KEY,
    pkg             VARCHAR(200) NOT NULL,
    name            VARCHAR(200) NOT NULL,
    customerid      INT REFERENCES customers (id) ON DELETE CASCADE,
    type            VARCHAR(20) NOT NULL DEFAULT 'app',
    common          BOOLEAN NOT NULL DEFAULT FALSE,
    showicon        BOOLEAN NOT NULL DEFAULT TRUE,
    system          BOOLEAN NOT NULL DEFAULT FALSE,
    url             TEXT,
    intent          TEXT,
    iconid          INT,
    runafterinstall BOOLEAN NOT NULL DEFAULT FALSE,
    runatboot       BOOLEAN NOT NULL DEFAULT FALSE,
    usekiosk        BOOLEAN NOT NULL DEFAULT FALSE,
    skipversion     BOOLEAN NOT NULL DEFAULT FALSE,
    icontext        VARCHAR(200),
    latestversion   INT
);

CREATE UNIQUE INDEX IF NOT EXISTS applications_pkg_customer_uidx
    ON applications (COALESCE(customerid, 0), lower(pkg));

CREATE TABLE IF NOT EXISTS applicationversions (
    id              SERIAL PRIMARY KEY,
    applicationid   INT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    version         VARCHAR(100),
    versioncode     INT NOT NULL DEFAULT 0,
    url             TEXT,
    urlarmeabi      TEXT,
    urlarm64        TEXT,
    filepath        TEXT,
    split           BOOLEAN NOT NULL DEFAULT FALSE,
    arch            VARCHAR(50),
    action          INT,
    showicon        BOOLEAN,
    screenorder     INT,
    keycode         INT,
    bottom          BOOLEAN,
    autoupdate      BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS configurationapplications (
    id                      SERIAL PRIMARY KEY,
    configurationid         INT NOT NULL REFERENCES configurations (id) ON DELETE CASCADE,
    applicationid           INT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    applicationversionid    INT REFERENCES applicationversions (id) ON DELETE SET NULL,
    action                  INT NOT NULL DEFAULT 1,
    showicon                BOOLEAN NOT NULL DEFAULT TRUE,
    screenorder             INT,
    keycode                 INT,
    bottom                  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS configurationapplications_uidx
    ON configurationapplications (configurationid, applicationid);

CREATE TABLE IF NOT EXISTS configurationfiles (
    id              SERIAL PRIMARY KEY,
    configurationid INT NOT NULL REFERENCES configurations (id) ON DELETE CASCADE,
    path            VARCHAR(500),
    externalurl     TEXT,
    url             TEXT,
    remove          BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS configurationapplicationsettings (
    id              SERIAL PRIMARY KEY,
    configurationid INT NOT NULL REFERENCES configurations (id) ON DELETE CASCADE,
    applicationid   INT REFERENCES applications (id) ON DELETE SET NULL,
    name            VARCHAR(200),
    type            VARCHAR(50),
    value           TEXT
);

INSERT INTO permissions (name, description, superadmin)
SELECT 'applications', 'Manage applications', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'applications');

INSERT INTO permissions (name, description, superadmin)
SELECT 'configurations', 'Manage device configurations', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'configurations');

INSERT INTO userrolepermissions (roleid, permissionid)
SELECT 2, p.id
FROM permissions p
WHERE lower(p.name) IN ('applications', 'configurations')
  AND NOT EXISTS (
      SELECT 1 FROM userrolepermissions urp
      WHERE urp.roleid = 2 AND urp.permissionid = p.id
  );

INSERT INTO applications (pkg, name, customerid, type, common, showicon, url)
SELECT 'com.headwind.mdm.demo', 'MDM Demo App', 1, 'app', FALSE, TRUE, ''
WHERE NOT EXISTS (
    SELECT 1 FROM applications WHERE customerid = 1 AND lower(pkg) = 'com.headwind.mdm.demo'
);

INSERT INTO applicationversions (applicationid, version, versioncode, url)
SELECT a.id, '1.0.0', 1, ''
FROM applications a
WHERE a.customerid = 1 AND lower(a.pkg) = 'com.headwind.mdm.demo'
  AND NOT EXISTS (
      SELECT 1 FROM applicationversions av
      WHERE av.applicationid = a.id AND av.version = '1.0.0'
  );

UPDATE configurations
SET description = COALESCE(description, 'Default device configuration'),
    type = 0,
    defaultfilepath = '/',
    settingsjson = COALESCE(settingsjson, '{}'::jsonb)
WHERE customerid = 1 AND lower(name) = 'default';

INSERT INTO configurationapplications (configurationid, applicationid, applicationversionid, action, showicon)
SELECT c.id, a.id, av.id, 1, TRUE
FROM configurations c
JOIN applications a ON a.customerid = c.customerid AND lower(a.pkg) = 'com.headwind.mdm.demo'
JOIN applicationversions av ON av.applicationid = a.id AND av.version = '1.0.0'
WHERE c.customerid = 1 AND lower(c.name) = 'default'
ON CONFLICT (configurationid, applicationid) DO NOTHING;
