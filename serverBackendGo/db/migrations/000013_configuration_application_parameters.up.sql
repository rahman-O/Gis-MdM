CREATE TABLE IF NOT EXISTS configurationapplicationparameters (
    id SERIAL PRIMARY KEY,
    configurationid INT NOT NULL REFERENCES configurations (id) ON DELETE CASCADE,
    applicationid INT NOT NULL REFERENCES applications (id) ON DELETE CASCADE,
    skipversioncheck BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS cap_config_app_uidx
    ON configurationapplicationparameters (configurationid, applicationid);
