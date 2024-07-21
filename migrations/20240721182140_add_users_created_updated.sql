-- +goose Up
-- +goose StatementBegin
ALTER TABLE USERS ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT now();
ALTER TABLE USERS ADD COLUMN updated_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE USERS DROP COLUMN created_at;
ALTER TABLE USERS DROP COLUMN updated_at;
-- +goose StatementEnd
