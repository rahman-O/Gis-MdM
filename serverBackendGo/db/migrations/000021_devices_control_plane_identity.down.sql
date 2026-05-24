ALTER TABLE configurations
    DROP COLUMN IF EXISTS default_device_id_mode,
    DROP COLUMN IF EXISTS default_tree_node_id;

DROP INDEX IF EXISTS devices_agent_id_uidx;

ALTER TABLE devices
    DROP COLUMN IF EXISTS enrollment_route_id,
    DROP COLUMN IF EXISTS enrollment_state,
    DROP COLUMN IF EXISTS agent_id;
