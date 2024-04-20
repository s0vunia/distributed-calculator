CREATE TABLE IF NOT EXISTS apps
(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50),
    secret Varchar
);

INSERT INTO apps(id, name, secret) VALUES (1, 'orchestrator', 'une-3r0yj*1+le22$x2y8=q%nag2q1(8brlbmmr(6ixh_$qa-#')
