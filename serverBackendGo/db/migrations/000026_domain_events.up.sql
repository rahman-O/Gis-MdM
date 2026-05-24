-- Domain events outbox (017 US5 publish notifications)

CREATE TABLE IF NOT EXISTS domain_events (
    id              BIGSERIAL PRIMARY KEY,
    event_type      VARCHAR(64) NOT NULL,
    aggregate_id    VARCHAR(64) NOT NULL,
    payload         JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS domain_events_unprocessed_idx
    ON domain_events (created_at)
    WHERE processed_at IS NULL;
