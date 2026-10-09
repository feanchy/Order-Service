CREATE TABLE processed_events (
    event_id TEXT PRIMARY KEY,
    processed_at TIMESTAMP NOT NULL DEFAULT NOW()
);