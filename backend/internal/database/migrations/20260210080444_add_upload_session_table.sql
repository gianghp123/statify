-- +goose Up
-- +goose StatementBegin
CREATE TYPE upload_status AS ENUM ('UPLOADED', 'WAITING', 'EXPIRED');
CREATE TABLE deployment_upload_session (
	id BIGSERIAL PRIMARY KEY,
	project_id BIGINT NOT NULL,
	upload_key VARCHAR(255) NOT NULL,
	output_prefix VARCHAR(255) NOT NULL,
	presigned_url VARCHAR(255) NOT NULL,
	status upload_status NOT NULL,
	expired_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	CONSTRAINT fk_project_id FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE upload_status;
DROP TABLE deployment_upload_session;
-- +goose StatementEnd
