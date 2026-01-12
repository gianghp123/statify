-- +goose Up
-- 1. Identify the current constraint name (usually 'fk_deployments_project')
-- 2. Drop it
ALTER TABLE deployments 
DROP CONSTRAINT IF EXISTS fk_deployments_project;

-- 3. Add it back with CASCADE
ALTER TABLE deployments
ADD CONSTRAINT fk_deployments_project
FOREIGN KEY (project_id) 
REFERENCES projects(id) 
ON DELETE CASCADE;

-- +goose Down
-- Revert back to a standard foreign key (RESTRICT/NO ACTION)
ALTER TABLE deployments 
DROP CONSTRAINT IF EXISTS fk_deployments_project;

ALTER TABLE deployments
ADD CONSTRAINT fk_deployments_project
FOREIGN KEY (project_id) 
REFERENCES projects(id);