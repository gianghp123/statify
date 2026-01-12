-- +goose Up
-- +goose StatementBegin
ALTER TABLE deployments DROP COLUMN source_zip_object_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE deployments ADD COLUMN source_zip_object_key VARCHAR(255);
-- +goose StatementEnd
