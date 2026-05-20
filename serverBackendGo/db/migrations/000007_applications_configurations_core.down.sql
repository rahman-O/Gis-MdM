DROP TABLE IF EXISTS configurationapplicationsettings;
DROP TABLE IF EXISTS configurationfiles;
DROP TABLE IF EXISTS configurationapplications;
DROP TABLE IF EXISTS applicationversions;
DROP TABLE IF EXISTS applications;

DROP INDEX IF EXISTS configurations_name_customer_uidx;

ALTER TABLE configurations DROP COLUMN IF EXISTS settingsjson;
ALTER TABLE configurations DROP COLUMN IF EXISTS contentappid;
ALTER TABLE configurations DROP COLUMN IF EXISTS defaultfilepath;
ALTER TABLE configurations DROP COLUMN IF EXISTS baseurl;
ALTER TABLE configurations DROP COLUMN IF EXISTS qrcodekey;
ALTER TABLE configurations DROP COLUMN IF EXISTS backgroundimageurl;
ALTER TABLE configurations DROP COLUMN IF EXISTS textcolor;
ALTER TABLE configurations DROP COLUMN IF EXISTS backgroundcolor;
ALTER TABLE configurations DROP COLUMN IF EXISTS password;
ALTER TABLE configurations DROP COLUMN IF EXISTS type;

DELETE FROM userrolepermissions
WHERE permissionid IN (SELECT id FROM permissions WHERE lower(name) IN ('applications', 'configurations'));

DELETE FROM permissions WHERE lower(name) IN ('applications', 'configurations');
