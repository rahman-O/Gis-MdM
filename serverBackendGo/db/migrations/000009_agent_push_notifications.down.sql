DROP TABLE IF EXISTS plugin_push_schedule;
DROP TABLE IF EXISTS plugin_push_messages;
DROP TABLE IF EXISTS pendingpushes;
DROP TABLE IF EXISTS pushmessages;

DELETE FROM userrolepermissions urp
USING permissions p
WHERE urp.permissionid = p.id
  AND lower(p.name) IN ('push_api', 'plugin_push_send', 'plugin_push_delete');

DELETE FROM permissions WHERE lower(name) IN ('push_api', 'plugin_push_send', 'plugin_push_delete');
