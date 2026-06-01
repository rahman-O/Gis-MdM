-- Rollback 033: Remove enhanced columns and index from device_locations

DROP INDEX IF EXISTS idx_device_locations_device_timestamp;

ALTER TABLE device_locations
    DROP COLUMN IF EXISTS altitude,
    DROP COLUMN IF EXISTS battery_level,
    DROP COLUMN IF EXISTS network_type,
    DROP COLUMN IF EXISTS tracking_mode,
    DROP COLUMN IF EXISTS received_at,
    DROP COLUMN IF EXISTS month;
