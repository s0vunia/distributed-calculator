CREATE TABLE IF NOT EXISTS agents
(
    id UUID PRIMARY KEY,
    heartbeat timestamp NOT NULL DEFAULT NOW()
);
