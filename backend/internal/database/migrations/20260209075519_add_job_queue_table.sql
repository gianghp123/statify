-- +goose Up
-- +goose StatementBegin
CREATE TYPE job_queue_type AS ENUM ('deployment_delete', 'deployment_process');
CREATE TYPE job_queue_status AS ENUM ('pending', 'processing', 'success', 'failed');

CREATE TABLE job_queues (
    id BIGSERIAL PRIMARY KEY,
    type job_queue_type NOT NULL,
    deployment_id BIGINT NOT NULL,
    payload TEXT,
    status job_queue_status NOT NULL DEFAULT 'pending',
    retry_count INTEGER DEFAULT 0,
    error TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT fk_job_queue_deployment FOREIGN KEY (deployment_id) REFERENCES deployments(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE job_queues;
DROP TYPE job_queue_type;
DROP TYPE job_queue_status;
-- +goose StatementEnd
