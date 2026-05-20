-- Phase 6: uploaded files, icons, configuration file linkage, file permissions.

CREATE TABLE IF NOT EXISTS uploadedfiles (
    id               SERIAL PRIMARY KEY,
    customerid       INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    filepath         TEXT,
    description      TEXT,
    uploadtime       BIGINT NOT NULL DEFAULT (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT,
    devicepath       TEXT,
    external         BOOLEAN NOT NULL DEFAULT FALSE,
    externalurl      TEXT,
    replacevariables BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS uploadedfiles_customerid_idx ON uploadedfiles (customerid);

CREATE TABLE IF NOT EXISTS icons (
    id          SERIAL PRIMARY KEY,
    customerid  INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    name        VARCHAR(64) NOT NULL,
    fileid      INT NOT NULL REFERENCES uploadedfiles (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS icons_customerid_idx ON icons (customerid);

ALTER TABLE configurationfiles ADD COLUMN IF NOT EXISTS fileid INT REFERENCES uploadedfiles (id) ON DELETE CASCADE;

ALTER TABLE customers ADD COLUMN IF NOT EXISTS sizelimit INT NOT NULL DEFAULT 0;

INSERT INTO permissions (name, description, superadmin)
SELECT 'files', 'Browse uploaded files', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'files');

INSERT INTO permissions (name, description, superadmin)
SELECT 'edit_files', 'Upload and edit files', FALSE
WHERE NOT EXISTS (SELECT 1 FROM permissions WHERE lower(name) = 'edit_files');

INSERT INTO userrolepermissions (roleid, permissionid)
SELECT 2, p.id
FROM permissions p
WHERE lower(p.name) IN ('files', 'edit_files')
  AND NOT EXISTS (
      SELECT 1 FROM userrolepermissions urp
      WHERE urp.roleid = 2 AND urp.permissionid = p.id
  );

INSERT INTO uploadedfiles (customerid, filepath, description, uploadtime, devicepath, external, replacevariables)
SELECT 1, 'readme.txt', 'Sample file', (EXTRACT(EPOCH FROM NOW()) * 1000)::BIGINT, '/readme.txt', FALSE, FALSE
WHERE NOT EXISTS (SELECT 1 FROM uploadedfiles WHERE customerid = 1 AND filepath = 'readme.txt');

INSERT INTO icons (customerid, name, fileid)
SELECT 1, 'Default', uf.id
FROM uploadedfiles uf
WHERE uf.customerid = 1 AND uf.filepath = 'readme.txt'
  AND NOT EXISTS (SELECT 1 FROM icons WHERE customerid = 1 AND name = 'Default');
