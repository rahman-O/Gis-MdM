ALTER TABLE devices DROP CONSTRAINT IF EXISTS devices_enrollment_route_id_fkey;

ALTER TABLE devices
    ADD CONSTRAINT devices_enrollment_route_id_fkey
    FOREIGN KEY (enrollment_route_id) REFERENCES configurations (id);

UPDATE devices d
SET enrollment_route_id = er.legacy_configuration_id
FROM enrollment_routes er
WHERE er.id = d.enrollment_route_id AND er.legacy_configuration_id IS NOT NULL;

DROP TABLE IF EXISTS enrollment_routes;
