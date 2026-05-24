DROP INDEX IF EXISTS devices_tree_node_id_idx;
ALTER TABLE devices DROP COLUMN IF EXISTS tree_node_id;
