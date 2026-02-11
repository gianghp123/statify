-- +goose Up
-- +goose StatementBegin
ALTER TYPE job_queue_type ADD VALUE 'project_delete';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TYPE job_queue_type DELETE VALUE 'project_delete';
-- +goose StatementEnd
