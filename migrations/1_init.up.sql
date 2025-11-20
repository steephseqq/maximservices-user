CREATE TABLE IF NOT EXISTS users(
    uuid TEXT NOT NULL,
    email TEXT NOT NULL,
    username TEXT,
    name TEXT,
    bio TEXT,
    pass_hash BYTEA NOT NULL,
    avatar_url TEXT NOT NULL,
    last_seen TIMESTAMP NOT NULL
)
