-- Device tree folders (017-device-control-plane US1)

CREATE TABLE IF NOT EXISTS device_tree_nodes (
    id          SERIAL PRIMARY KEY,
    customerid  INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    parent_id   INT REFERENCES device_tree_nodes (id) ON DELETE CASCADE,
    name        VARCHAR(200) NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    path        TEXT NOT NULL,
    depth       INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS device_tree_nodes_parent_idx ON device_tree_nodes (parent_id);
CREATE INDEX IF NOT EXISTS device_tree_nodes_path_idx ON device_tree_nodes (customerid, path);
CREATE UNIQUE INDEX IF NOT EXISTS device_tree_nodes_name_parent_uidx
    ON device_tree_nodes (customerid, parent_id, lower(name));
