-- +goose Up
-- +goose StatementBegin
ALTER TYPE deployment_status ADD VALUE 'WAITING_UPLOAD';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TYPE deployment_status DROP VALUE 'WAITING_UPLOAD';
-- +goose StatementEnd
