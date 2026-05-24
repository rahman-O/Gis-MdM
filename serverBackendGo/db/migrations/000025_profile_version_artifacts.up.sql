-- Profile version compiled artifacts (017 US5)

CREATE TABLE IF NOT EXISTS profile_version_artifacts (
    profile_version_id  INT PRIMARY KEY REFERENCES profile_versions (id) ON DELETE CASCADE,
    artifact_json       JSONB NOT NULL,
    artifact_hash       VARCHAR(64) NOT NULL,
    compiled_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS profile_version_artifacts_hash_idx
    ON profile_version_artifacts (artifact_hash);
