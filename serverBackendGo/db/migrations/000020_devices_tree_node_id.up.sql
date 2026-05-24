-- Link devices to tree folders (017-device-control-plane US1)

ALTER TABLE devices ADD COLUMN IF NOT EXISTS tree_node_id INT REFERENCES device_tree_nodes (id);

CREATE INDEX IF NOT EXISTS devices_tree_node_id_idx ON devices (tree_node_id);

-- Backfill: assign existing devices to per-customer root folder when present
UPDATE devices d
SET tree_node_id = root.id
FROM device_tree_nodes root
WHERE d.tree_node_id IS NULL
  AND root.customerid = d.customerid
  AND root.parent_id IS NULL
  AND root.path = '/' || root.id::TEXT || '/';
