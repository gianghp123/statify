-- +goose Up
-- +goose StatementBegin
ALTER TABLE deployments ADD COLUMN is_spa BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE deployments DROP COLUMN is_spa;
-- +goose StatementEnd
