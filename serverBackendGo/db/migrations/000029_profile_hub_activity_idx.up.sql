-- 019 profile hub: faster activity timeline lookups

CREATE INDEX IF NOT EXISTS domain_events_aggregate_type_idx
    ON domain_events (aggregate_id, event_type, created_at DESC);
