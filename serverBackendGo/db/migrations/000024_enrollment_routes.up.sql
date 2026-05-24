-- Enrollment routes (017 US4)

CREATE TABLE IF NOT EXISTS enrollment_routes (
    id                      SERIAL PRIMARY KEY,
    customerid              INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    name                    VARCHAR(200) NOT NULL,
    description             TEXT,
    qrcodekey               VARCHAR(200),
    mainappid               INT,
    profile_version_id      INT REFERENCES profile_versions (id),
    default_tree_node_id    INT REFERENCES device_tree_nodes (id),
    default_device_id_mode  VARCHAR(20) NOT NULL DEFAULT 'imei',
    type                    INT NOT NULL DEFAULT 0,
    legacy_configuration_id INT UNIQUE REFERENCES configurations (id) ON DELETE SET NULL,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS enrollment_routes_name_customer_uidx
    ON enrollment_routes (customerid, lower(name));

CREATE UNIQUE INDEX IF NOT EXISTS enrollment_routes_qrcodekey_uidx
    ON enrollment_routes (lower(qrcodekey))
    WHERE qrcodekey IS NOT NULL AND trim(qrcodekey) <> '';

INSERT INTO enrollment_routes (
    customerid, name, description, qrcodekey, mainappid,
    profile_version_id, default_tree_node_id, default_device_id_mode, type,
    legacy_configuration_id
)
SELECT c.customerid, c.name, COALESCE(c.description, ''), c.qrcodekey, c.mainappid,
       p.published_version_id, c.default_tree_node_id,
       COALESCE(NULLIF(TRIM(c.default_device_id_mode), ''), 'imei'), c.type, c.id
FROM configurations c
LEFT JOIN profiles p ON p.legacy_configuration_id = c.id
WHERE NOT EXISTS (
    SELECT 1 FROM enrollment_routes er WHERE er.legacy_configuration_id = c.id
);

ALTER TABLE devices DROP CONSTRAINT IF EXISTS devices_enrollment_route_id_fkey;

UPDATE devices d
SET enrollment_route_id = er.id
FROM enrollment_routes er
WHERE er.legacy_configuration_id = d.configurationid;

ALTER TABLE devices
    ADD CONSTRAINT devices_enrollment_route_id_fkey
    FOREIGN KEY (enrollment_route_id) REFERENCES enrollment_routes (id);
