CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY NOT NULL,
    username TEXT UNIQUE NOT NULL CHECK (username <> ''),
    password BYTEA NOT NULL
)
