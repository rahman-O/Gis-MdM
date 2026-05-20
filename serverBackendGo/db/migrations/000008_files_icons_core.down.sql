DROP TABLE IF EXISTS icons;
DROP TABLE IF EXISTS uploadedfiles;

ALTER TABLE configurationfiles DROP COLUMN IF EXISTS fileid;
ALTER TABLE customers DROP COLUMN IF EXISTS sizelimit;

DELETE FROM userrolepermissions
WHERE permissionid IN (SELECT id FROM permissions WHERE lower(name) IN ('files', 'edit_files'));

DELETE FROM permissions WHERE lower(name) IN ('files', 'edit_files');
