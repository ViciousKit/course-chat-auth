-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (name) VALUES ('unknown'), ('user'), ('admin');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM roles WHERE name IN ('unknown', 'user', 'admin');
-- +goose StatementEnd

