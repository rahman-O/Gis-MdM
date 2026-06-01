-- 033: Enhance device_locations with additional tracking columns and composite index

ALTER TABLE device_locations
    ADD COLUMN IF NOT EXISTS altitude DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS battery_level DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS network_type VARCHAR(20) DEFAULT 'none',
    ADD COLUMN IF NOT EXISTS tracking_mode VARCHAR(20) DEFAULT 'normal',
    ADD COLUMN IF NOT EXISTS received_at TIMESTAMPTZ DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS month VARCHAR(7) NOT NULL DEFAULT to_char(NOW(), 'YYYY-MM');

CREATE INDEX IF NOT EXISTS idx_device_locations_device_timestamp
    ON device_locations (device_id, timestamp DESC);
