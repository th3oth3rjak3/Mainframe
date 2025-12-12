-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL
);

INSERT INTO roles (id, name)
VALUES 
    ('05c9b67e-5cfa-4f01-974d-a77632637e23', 'Administrator'),
    ('1bd84e72-796c-453c-8083-91d42465830f', 'Basic User'),
    ('ce7ed876-d7fc-41f8-a3d3-d245e6d725c8', 'Recipe User');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE roles;
-- +goose StatementEnd