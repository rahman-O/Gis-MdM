DROP INDEX IF EXISTS devices_target_profile_version_idx;
DROP INDEX IF EXISTS devices_profile_rollout_status_idx;

ALTER TABLE devices
    DROP COLUMN IF EXISTS profile_rollout_updated_at,
    DROP COLUMN IF EXISTS profile_rollout_reason,
    DROP COLUMN IF EXISTS profile_rollout_status,
    DROP COLUMN IF EXISTS applied_profile_version_id,
    DROP COLUMN IF EXISTS target_profile_version_id;

DROP TABLE IF EXISTS profile_tree_assignments;

ALTER TABLE profiles DROP COLUMN IF EXISTS enabled;
