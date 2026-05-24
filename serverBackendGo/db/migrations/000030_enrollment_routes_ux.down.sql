DROP INDEX IF EXISTS domain_events_enrollment_qr_viewed_idx;

ALTER TABLE devices DROP CONSTRAINT IF EXISTS devices_enrollment_route_id_fkey;

ALTER TABLE devices
    ADD CONSTRAINT devices_enrollment_route_id_fkey
    FOREIGN KEY (enrollment_route_id) REFERENCES enrollment_routes (id);

ALTER TABLE enrollment_routes
    DROP COLUMN IF EXISTS container_placement_ack_at,
    DROP COLUMN IF EXISTS bootstrap_version_id,
    DROP COLUMN IF EXISTS bootstrap_application_id,
    DROP COLUMN IF EXISTS bootstrap_intent;

DROP INDEX IF EXISTS applicationversions_one_recommended_uidx;

ALTER TABLE applicationversions DROP COLUMN IF EXISTS is_recommended;
