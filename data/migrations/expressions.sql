CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS expressions
(
    id UUID PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    idempotency_key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    state VARCHAR(50) NOT NULL,
    result DOUBLE PRECISION,
    created_at timestamp NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, idempotency_key)
);
