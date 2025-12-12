-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id TEXT PRIMARY KEY NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    last_login DATETIME,
    failed_login_attempts INTEGER DEFAULT 0,
    last_failed_login_attempt DATETIME,
    is_disabled INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users_username;
DROP TABLE users;
-- +goose StatementEnd
