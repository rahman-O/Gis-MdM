-- Profiles + versions (017 US3)

CREATE TABLE IF NOT EXISTS profiles (
    id              SERIAL PRIMARY KEY,
    customerid      INT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    draft_version_id INT,
    published_version_id INT,
    legacy_configuration_id INT REFERENCES configurations (id) ON DELETE SET NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS profiles_name_customer_uidx
    ON profiles (customerid, lower(name));

CREATE TABLE IF NOT EXISTS profile_versions (
    id              SERIAL PRIMARY KEY,
    profile_id      INT NOT NULL REFERENCES profiles (id) ON DELETE CASCADE,
    version_number  INT NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'draft',
    type            INT NOT NULL DEFAULT 0,
    password        TEXT,
    backgroundcolor VARCHAR(50),
    textcolor       VARCHAR(50),
    backgroundimageurl TEXT,
    qrcodekey       VARCHAR(200),
    baseurl         TEXT,
    defaultfilepath VARCHAR(500),
    mainappid       INT,
    contentappid    INT,
    permissive      BOOLEAN NOT NULL DEFAULT FALSE,
    settingsjson    JSONB NOT NULL DEFAULT '{}',
    published_at    TIMESTAMPTZ,
    published_by    INT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (profile_id, version_number)
);

CREATE INDEX IF NOT EXISTS profile_versions_profile_status_idx
    ON profile_versions (profile_id, status);
