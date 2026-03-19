CREATE TABLE IF NOT EXISTS blog (
    id          SERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_deleted  BOOL NOT NULL DEFAULT FALSE,
    title       TEXT NOT NULL,
    content     TEXT
);

CREATE TABLE IF NOT EXISTS users (
    id       SERIAL PRIMARY KEY,
    username TEXT NOT NULL DEFAULT '',
    password TEXT NOT NULL DEFAULT '',
	session_token text not null default '',
	csrf_token text not null default ''
);
