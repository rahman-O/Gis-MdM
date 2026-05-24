-- Profile version junction tables (mirror configuration junctions)

CREATE TABLE IF NOT EXISTS profile_version_applications (
    id                      SERIAL PRIMARY KEY,
    profile_version_id      INT NOT NULL REFERENCES profile_versions (id) ON DELETE CASCADE,
    applicationid           INT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    applicationversionid    INT REFERENCES applicationversions (id) ON DELETE SET NULL,
    action                  INT NOT NULL DEFAULT 1,
    showicon                BOOLEAN NOT NULL DEFAULT TRUE,
    screenorder             INT,
    keycode                 INT,
    bottom                  BOOLEAN NOT NULL DEFAULT FALSE,
    remove                  BOOLEAN NOT NULL DEFAULT FALSE,
    longtap                 BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS profile_version_applications_uidx
    ON profile_version_applications (profile_version_id, applicationid);

CREATE TABLE IF NOT EXISTS profile_version_files (
    id                  SERIAL PRIMARY KEY,
    profile_version_id  INT NOT NULL REFERENCES profile_versions (id) ON DELETE CASCADE,
    path                VARCHAR(500),
    externalurl         TEXT,
    url                 TEXT,
    remove              BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS profile_version_application_settings (
    id                  SERIAL PRIMARY KEY,
    profile_version_id  INT NOT NULL REFERENCES profile_versions (id) ON DELETE CASCADE,
    applicationid       INT REFERENCES applications (id) ON DELETE SET NULL,
    name                VARCHAR(200),
    type                VARCHAR(50),
    value               TEXT,
    readonly            BOOLEAN NOT NULL DEFAULT FALSE,
    comment             TEXT,
    variable            VARCHAR(200)
);

CREATE TABLE IF NOT EXISTS profile_version_application_parameters (
    profile_version_id  INT NOT NULL REFERENCES profile_versions (id) ON DELETE CASCADE,
    applicationid       INT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    value               TEXT,
    PRIMARY KEY (profile_version_id, applicationid, name)
);

-- Backfill from existing configurations (1:1 profile + published v1)
INSERT INTO profiles (customerid, name, description, legacy_configuration_id)
SELECT c.customerid, c.name, COALESCE(c.description, ''), c.id
FROM configurations c
WHERE NOT EXISTS (SELECT 1 FROM profiles p WHERE p.legacy_configuration_id = c.id);

INSERT INTO profile_versions (
    profile_id, version_number, status, type, password, backgroundcolor, textcolor,
    backgroundimageurl, qrcodekey, baseurl, defaultfilepath, mainappid, contentappid,
    permissive, settingsjson, published_at
)
SELECT p.id, 1, 'published', c.type, c.password, c.backgroundcolor, c.textcolor,
       c.backgroundimageurl, c.qrcodekey, c.baseurl, c.defaultfilepath, c.mainappid, c.contentappid,
       c.permissive, c.settingsjson, NOW()
FROM configurations c
JOIN profiles p ON p.legacy_configuration_id = c.id
WHERE NOT EXISTS (
    SELECT 1 FROM profile_versions pv WHERE pv.profile_id = p.id AND pv.version_number = 1
);

UPDATE profiles p
SET published_version_id = pv.id
FROM profile_versions pv
WHERE pv.profile_id = p.id AND pv.status = 'published' AND pv.version_number = 1
  AND p.published_version_id IS NULL;

INSERT INTO profile_version_applications (
    profile_version_id, applicationid, applicationversionid, action, showicon, screenorder, keycode, bottom, remove, longtap
)
SELECT pv.id, ca.applicationid, ca.applicationversionid, ca.action, ca.showicon, ca.screenorder, ca.keycode, ca.bottom,
       COALESCE(ca.remove, FALSE), COALESCE(ca.longtap, FALSE)
FROM configurationapplications ca
JOIN profiles p ON p.legacy_configuration_id = ca.configurationid
JOIN profile_versions pv ON pv.profile_id = p.id AND pv.version_number = 1
ON CONFLICT DO NOTHING;

INSERT INTO profile_version_files (profile_version_id, path, externalurl, url, remove)
SELECT pv.id, cf.path, cf.externalurl, cf.url, cf.remove
FROM configurationfiles cf
JOIN profiles p ON p.legacy_configuration_id = cf.configurationid
JOIN profile_versions pv ON pv.profile_id = p.id AND pv.version_number = 1
ON CONFLICT DO NOTHING;

INSERT INTO profile_version_application_settings (
    profile_version_id, applicationid, name, type, value
)
SELECT pv.id, cas.applicationid, cas.name, cas.type, cas.value
FROM configurationapplicationsettings cas
JOIN profiles p ON p.legacy_configuration_id = cas.configurationid
JOIN profile_versions pv ON pv.profile_id = p.id AND pv.version_number = 1;

INSERT INTO profile_version_application_parameters (profile_version_id, applicationid, name, value)
SELECT pv.id, cap.applicationid, 'skipVersionCheck',
       CASE WHEN cap.skipversioncheck THEN 'true' ELSE 'false' END
FROM configurationapplicationparameters cap
JOIN profiles p ON p.legacy_configuration_id = cap.configurationid
JOIN profile_versions pv ON pv.profile_id = p.id AND pv.version_number = 1
ON CONFLICT DO NOTHING;
