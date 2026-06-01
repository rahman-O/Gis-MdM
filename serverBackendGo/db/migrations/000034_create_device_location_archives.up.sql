-- 034: Create device_location_archives table for hourly summaries

CREATE TABLE IF NOT EXISTS device_location_archives (
    id BIGSERIAL PRIMARY KEY,
    device_id INTEGER NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    hour_start BIGINT NOT NULL,
    start_latitude DOUBLE PRECISION NOT NULL,
    start_longitude DOUBLE PRECISION NOT NULL,
    end_latitude DOUBLE PRECISION NOT NULL,
    end_longitude DOUBLE PRECISION NOT NULL,
    distance_traveled DOUBLE PRECISION NOT NULL DEFAULT 0,
    point_count INTEGER NOT NULL DEFAULT 0,
    archived_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (device_id, hour_start)
);

CREATE INDEX IF NOT EXISTS idx_device_location_archives_device_hour
    ON device_location_archives (device_id, hour_start);
