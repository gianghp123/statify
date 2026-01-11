-- +goose Up
-- +goose StatementBegin

-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "username" text NOT NULL,
  "password_hash" text NOT NULL,
  "email" text NOT NULL,
  "role" "public"."role" NULL DEFAULT 'USER',
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email"),
  CONSTRAINT "uni_users_username" UNIQUE ("username")
);
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");

-- Create "projects" table
CREATE TABLE "public"."projects" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NOT NULL,
  "subdomain" text NOT NULL,
  "user_id" bigint NULL,
  "current_deployment_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_projects_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE INDEX "idx_projects_deleted_at" ON "public"."projects" ("deleted_at");

-- Create "deployments" table
CREATE TABLE "public"."deployments" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "project_id" bigint NULL,
  "status" "public"."deployment_status" NULL DEFAULT 'UPLOADED',
  "output_prefix" text NOT NULL,
  "source_zip_object_key" text NOT NULL,
  "validation_error" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_deployments_project" FOREIGN KEY ("project_id") REFERENCES "public"."projects" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE INDEX "idx_deployments_deleted_at" ON "public"."deployments" ("deleted_at");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop tables in reverse order to handle Foreign Key constraints
DROP TABLE IF EXISTS "public"."deployments";
DROP TABLE IF EXISTS "public"."projects";
DROP TABLE IF EXISTS "public"."users";

-- Drop custom enum types
DROP TYPE IF EXISTS "public"."deployment_status";
DROP TYPE IF EXISTS "public"."role";
-- +goose StatementEnd