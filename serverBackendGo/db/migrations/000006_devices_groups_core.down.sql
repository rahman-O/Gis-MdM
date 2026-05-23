DROP TABLE IF EXISTS deviceapplicationsettings;
DROP TABLE IF EXISTS devicegroups;
DROP TABLE IF EXISTS devices;

DELETE FROM userrolepermissions urp
USING permissions p
WHERE urp.permissionid = p.id AND lower(p.name) IN ('edit_devices', 'edit_device_desc');

DELETE FROM permissions WHERE lower(name) IN ('edit_devices', 'edit_device_desc');

DROP INDEX IF EXISTS groups_name_customer_uidx;
DROP INDEX IF EXISTS devices_number_customer_uidx;

ALTER TABLE configurations DROP COLUMN IF EXISTS mainappid;
ALTER TABLE configurations DROP COLUMN IF EXISTS permissive;
ALTER TABLE configurations DROP COLUMN IF EXISTS customerid;

ALTER TABLE groups DROP COLUMN IF EXISTS customerid;
