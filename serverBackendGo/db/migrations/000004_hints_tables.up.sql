-- Hint tracking tables (Liquibase 19.09.19 parity).

CREATE TABLE IF NOT EXISTS userhinttypes (
    hintkey VARCHAR(100) NOT NULL PRIMARY KEY
);

INSERT INTO userhinttypes (hintkey) VALUES
    ('hint.step.1'),
    ('hint.step.2'),
    ('hint.step.3'),
    ('hint.step.4')
ON CONFLICT (hintkey) DO NOTHING;

CREATE TABLE IF NOT EXISTS userhints (
    id         SERIAL PRIMARY KEY,
    userid     INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    hintkey    VARCHAR(100) NOT NULL,
    created    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT userhints_userid_hintkey_unique UNIQUE (userid, hintkey)
);

CREATE INDEX IF NOT EXISTS userhints_userid_idx ON userhints (userid);
