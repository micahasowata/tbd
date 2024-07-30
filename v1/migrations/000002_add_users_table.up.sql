CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY NOT NULL,
    name citext NOT NULL,
    email citext UNIQUE NOT NULL,
    password bytea NOT NULL
);