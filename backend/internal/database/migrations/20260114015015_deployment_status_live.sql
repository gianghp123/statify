-- +goose Up
-- +goose StatementBegin
ALTER TYPE deployment_status ADD VALUE 'LIVE';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TYPE deployment_status DROP VALUE 'LIVE';
-- +goose StatementEnd
