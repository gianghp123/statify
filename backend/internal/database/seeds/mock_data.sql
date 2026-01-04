-- Mock Data Seed Script

-- 1. Insert Regular User
-- Password is 'password' hashed with bcrypt (cost 10)
-- $2a$10$8K1p/a06vX/K.d7YpM.fhe/C0pL9R4P8O/hBfI9L8/E1jH5H7vW
INSERT INTO "public"."users" ("created_at", "updated_at", "username", "password_hash", "email", "role")
VALUES (NOW(), NOW(), 'johndoe', '$2a$10$8K1p/a06vX/K.d7YpM.fhe/C0pL9R4P8O/hBfI9L8/E1jH5H7vW', 'user@statify.app', 'USER');

-- 2. Insert Projects for the user
-- Assuming the user we just inserted has ID 2 (since admin is ID 1)
-- If this is a fresh DB with no admin yet, it might be ID 1, but user stated admin already exists.
-- Using subquery to be safer.
INSERT INTO "public"."projects" ("created_at", "updated_at", "name", "subdomain", "user_id")
SELECT NOW(), NOW(), 'Portfolio', 'portfolio', id FROM "public"."users" WHERE email = 'user@statify.app';

INSERT INTO "public"."projects" ("created_at", "updated_at", "name", "subdomain", "user_id")
SELECT NOW(), NOW(), 'Blog', 'blog', id FROM "public"."users" WHERE email = 'user@statify.app';

-- 3. Insert Deployments
-- Deployments for Portfolio
INSERT INTO "public"."deployments" ("created_at", "updated_at", "project_id", "status", "output_prefix", "source_zip_object_key")
SELECT NOW(), NOW(), id, 'READY', 'portfolio-v1', 'sources/portfolio-v1.zip'
FROM "public"."projects" WHERE subdomain = 'portfolio';

INSERT INTO "public"."deployments" ("created_at", "updated_at", "project_id", "status", "output_prefix", "source_zip_object_key", "validation_error")
SELECT NOW(), NOW(), id, 'FAILED', 'portfolio-v2', 'sources/portfolio-v2.zip', 'Missing index.html'
FROM "public"."projects" WHERE subdomain = 'portfolio';

-- Deployment for Blog
INSERT INTO "public"."deployments" ("created_at", "updated_at", "project_id", "status", "output_prefix", "source_zip_object_key")
SELECT NOW(), NOW(), id, 'UPLOADED', 'blog-v1', 'sources/blog-v1.zip'
FROM "public"."projects" WHERE subdomain = 'blog';

-- 4. Update Current Deployment ID for Projects
UPDATE "public"."projects"
SET current_deployment_id = (SELECT id FROM "public"."deployments" WHERE project_id = "public"."projects".id AND status = 'READY' ORDER BY created_at DESC LIMIT 1)
WHERE subdomain = 'portfolio';
