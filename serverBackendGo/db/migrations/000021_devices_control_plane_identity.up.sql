-- Device control plane identity + enrollment binding (017 US2)

ALTER TABLE configurations
    ADD COLUMN IF NOT EXISTS default_tree_node_id INT REFERENCES device_tree_nodes (id),
    ADD COLUMN IF NOT EXISTS default_device_id_mode VARCHAR(20) NOT NULL DEFAULT 'imei';

ALTER TABLE devices
    ADD COLUMN IF NOT EXISTS agent_id UUID NOT NULL DEFAULT gen_random_uuid(),
    ADD COLUMN IF NOT EXISTS enrollment_state VARCHAR(20) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS enrollment_route_id INT REFERENCES configurations (id);

CREATE UNIQUE INDEX IF NOT EXISTS devices_agent_id_uidx ON devices (agent_id);

UPDATE devices
SET enrollment_route_id = configurationid
WHERE enrollment_route_id IS NULL;

UPDATE devices
SET enrollment_state = CASE
    WHEN lastupdate > 0 THEN 'active'
    WHEN enrolltime > 0 THEN 'enrolled'
    ELSE 'pending'
END;

UPDATE configurations c
SET default_tree_node_id = root.id
FROM device_tree_nodes root
WHERE c.default_tree_node_id IS NULL
  AND root.customerid = c.customerid
  AND root.parent_id IS NULL;

UPDATE devices d
SET tree_node_id = c.default_tree_node_id
FROM configurations c
WHERE d.tree_node_id IS NULL
  AND d.configurationid = c.id
  AND c.default_tree_node_id IS NOT NULL;
