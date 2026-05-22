CREATE TABLE IF NOT EXISTS devicestatuses (
    deviceid INT NOT NULL PRIMARY KEY REFERENCES devices (id) ON DELETE CASCADE,
    configfilesstatus VARCHAR(100),
    applicationsstatus VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS devicestatuses_apps_status_idx ON devicestatuses (applicationsstatus);

INSERT INTO devicestatuses (deviceid, configfilesstatus, applicationsstatus)
SELECT d.id, 'OTHER', 'FAILURE'
FROM devices d
WHERE NOT EXISTS (SELECT 1 FROM devicestatuses ds WHERE ds.deviceid = d.id);
