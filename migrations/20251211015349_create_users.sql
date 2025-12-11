-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id TEXT PRIMARY KEY,
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

INSERT INTO users (id, username, email, first_name, last_name, password_hash) 
VALUES ('88d178e9-a2b8-48e3-b991-c2ebef654092', 'admin', 'admin@localhost', 'Local', 'Administrator', 'fake pw hash for now');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_users_username;
DROP TABLE users;
-- +goose StatementEnd
