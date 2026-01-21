-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public.deployment_metrics_minute (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL,
    deployment_id BIGINT NOT NULL,
    minute_ts TIMESTAMP NOT NULL,

    request_count BIGINT NOT NULL DEFAULT 0,
    bytes_served BIGINT NOT NULL DEFAULT 0,
    total_duration BIGINT NOT NULL DEFAULT 0,

    status_2xx BIGINT NOT NULL DEFAULT 0,
    status_3xx BIGINT NOT NULL DEFAULT 0,
    status_4xx BIGINT NOT NULL DEFAULT 0,
    status_5xx BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_dmm_project
        FOREIGN KEY (project_id)
        REFERENCES public.projects (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_dmm_deployment
        FOREIGN KEY (deployment_id)
        REFERENCES public.deployments (id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS public.deployment_metrics_minute;

-- +goose StatementEnd
