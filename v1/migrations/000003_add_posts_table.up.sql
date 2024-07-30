CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    user_id bigserial NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title citext NOT NULL,
    body citext NOT NULL
);