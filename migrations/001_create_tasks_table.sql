CREATE TABLE IF NOT EXISTS tasks
(
    id         uuid PRIMARY KEY,
    execute_at TIMESTAMP    NOT NULL,
    method     VARCHAR(15),
    url        VARCHAR(100) NOT NULL,
    payload    BYTEA
);
