-- 018 profile rollout: tree assignments, device rollout status, profile enable flag

ALTER TABLE profiles
    ADD COLUMN IF NOT EXISTS enabled BOOLEAN NOT NULL DEFAULT TRUE;

CREATE TABLE IF NOT EXISTS profile_tree_assignments (
    id                  SERIAL PRIMARY KEY,
    customerid          INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    profile_id          INT NOT NULL REFERENCES profiles (id) ON DELETE CASCADE,
    profile_version_id  INT NOT NULL REFERENCES profile_versions (id) ON DELETE CASCADE,
    tree_node_id        INT NOT NULL REFERENCES device_tree_nodes (id) ON DELETE CASCADE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by          INT REFERENCES users (id) ON DELETE SET NULL,
    UNIQUE (customerid, tree_node_id)
);

CREATE INDEX IF NOT EXISTS profile_tree_assignments_profile_idx
    ON profile_tree_assignments (profile_id);

ALTER TABLE devices
    ADD COLUMN IF NOT EXISTS target_profile_version_id INT REFERENCES profile_versions (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS applied_profile_version_id INT REFERENCES profile_versions (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS profile_rollout_status VARCHAR(20),
    ADD COLUMN IF NOT EXISTS profile_rollout_reason TEXT,
    ADD COLUMN IF NOT EXISTS profile_rollout_updated_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS devices_profile_rollout_status_idx
    ON devices (customerid, profile_rollout_status);

CREATE INDEX IF NOT EXISTS devices_target_profile_version_idx
    ON devices (target_profile_version_id);
