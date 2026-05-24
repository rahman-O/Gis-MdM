-- 021 enrollment routes UX: bootstrap intent, stable channel, container ack

ALTER TABLE applicationversions
    ADD COLUMN IF NOT EXISTS is_recommended BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS applicationversions_one_recommended_uidx
    ON applicationversions (applicationid)
    WHERE is_recommended = TRUE;

ALTER TABLE enrollment_routes
    ADD COLUMN IF NOT EXISTS bootstrap_intent VARCHAR(20) NOT NULL DEFAULT 'stable',
    ADD COLUMN IF NOT EXISTS bootstrap_application_id INT REFERENCES applications (id),
    ADD COLUMN IF NOT EXISTS bootstrap_version_id INT REFERENCES applicationversions (id),
    ADD COLUMN IF NOT EXISTS container_placement_ack_at TIMESTAMPTZ;

ALTER TABLE devices DROP CONSTRAINT IF EXISTS devices_enrollment_route_id_fkey;

ALTER TABLE devices
    ADD CONSTRAINT devices_enrollment_route_id_fkey
    FOREIGN KEY (enrollment_route_id) REFERENCES enrollment_routes (id) ON DELETE SET NULL;

-- Backfill bootstrap fields from existing mainappid
UPDATE enrollment_routes er
SET bootstrap_application_id = av.applicationid,
    bootstrap_version_id = er.mainappid,
    bootstrap_intent = 'specific'
FROM applicationversions av
WHERE er.mainappid IS NOT NULL
  AND av.id = er.mainappid
  AND er.bootstrap_application_id IS NULL;

CREATE INDEX IF NOT EXISTS domain_events_enrollment_qr_viewed_idx
    ON domain_events (event_type, aggregate_id, created_at DESC)
    WHERE event_type = 'enrollment_route.qr_viewed';
