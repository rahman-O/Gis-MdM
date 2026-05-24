-- Assign devices missing tree_node_id to customer root folder (017 polish T088)

UPDATE devices d
SET tree_node_id = root.id
FROM device_tree_nodes root
WHERE d.tree_node_id IS NULL
  AND root.customerid = d.customerid
  AND root.parent_id IS NULL;
