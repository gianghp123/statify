-- +goose Up
-- +goose StatementBegin
CREATE TYPE "public"."deployment_status" AS ENUM ('FAILED', 'UPLOADED', 'READY', 'DELETED', 'PROCESSING');
CREATE TYPE "public"."role" AS ENUM ('ADMIN', 'USER');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE "public"."deployment_status";
DROP TYPE "public"."role";
-- +goose StatementEnd
