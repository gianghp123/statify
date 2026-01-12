-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects ADD CONSTRAINT unique_subdomain UNIQUE (subdomain);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects DROP CONSTRAINT unique_subdomain;
-- +goose StatementEnd
