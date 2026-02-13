-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects
ADD COLUMN effective_status deployment_status NULL,
ADD COLUMN latest_deployment_id BIGINT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects
DROP COLUMN effective_status,
DROP COLUMN latest_deployment_id;
-- +goose StatementEnd
