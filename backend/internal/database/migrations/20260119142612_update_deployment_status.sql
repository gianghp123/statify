-- +goose Up
-- +goose NO TRANSACTION
DROP TYPE IF EXISTS deployment_status_new;

CREATE TYPE deployment_status_new AS ENUM (
  'QUEUED',
  'FAILED',
  'READY',
  'DELETED',
  'PROCESSING',
  'LIVE'
);

-- 2. Drop default FIRST
ALTER TABLE deployments
  ALTER COLUMN status DROP DEFAULT;

-- 3. Change column type
ALTER TABLE deployments
  ALTER COLUMN status TYPE deployment_status_new
  USING status::text::deployment_status_new;

-- 4. Set new default (IMPORTANT)
ALTER TABLE deployments
  ALTER COLUMN status SET DEFAULT 'QUEUED';

-- 5. Swap enum types
DROP TYPE deployment_status;
ALTER TYPE deployment_status_new RENAME TO deployment_status;



-- +goose Down
-- +goose NO TRANSACTION
DROP TYPE IF EXISTS deployment_status_old;

CREATE TYPE deployment_status_old AS ENUM (
  'UPLOADED',
  'FAILED',
  'READY',
  'DELETED',
  'PROCESSING',
  'LIVE'
);

ALTER TABLE deployments
  ALTER COLUMN status DROP DEFAULT;

ALTER TABLE deployments
  ALTER COLUMN status TYPE deployment_status_old
  USING status::text::deployment_status_old;

ALTER TABLE deployments
  ALTER COLUMN status SET DEFAULT 'UPLOADED';

DROP TYPE deployment_status;
ALTER TYPE deployment_status_old RENAME TO deployment_status;

