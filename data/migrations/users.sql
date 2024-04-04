CREATE TABLE IF NOT EXISTS users
(
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE,
    Pass_hash Varchar
);
