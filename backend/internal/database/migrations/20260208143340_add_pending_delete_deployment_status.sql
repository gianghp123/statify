-- +goose Up
-- +goose StatementBegin
ALTER TYPE deployment_status ADD VALUE 'PENDING_DELETE';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TYPE deployment_status DELETE VALUE 'PENDING_DELETE';
-- +goose StatementEnd
