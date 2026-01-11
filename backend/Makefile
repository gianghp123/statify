COMPOSE_FILE ?= .docker/docker-compose.local.yaml
ENV_FILE ?= .env.local

ifneq (,$(wildcard $(ENV_FILE)))
    include $(ENV_FILE)
    export
endif

# 2. Build the Goose DSN using your specific variables
# Note: Using $(POSTGRES_DB) ensures it hits 'statify' instead of the default 'postgres'
DB_URL="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"
MIG_DIR=internal/database/migrations

up:
	docker compose -f $(COMPOSE_FILE) --env-file .env.local up -d

up-build:
	docker compose -f $(COMPOSE_FILE) --env-file .env.local up --build

down:
	docker compose -f $(COMPOSE_FILE) down

down-v:
	docker compose -f $(COMPOSE_FILE) down -v

restart:
	docker compose -f $(COMPOSE_FILE) down
	docker compose -f $(COMPOSE_FILE) --env-file .env.local up

restart-build:
	docker compose -f $(COMPOSE_FILE) down -v
	docker compose -f $(COMPOSE_FILE) --env-file .env.local up --build

# Apply all migrations
goose-up:
	goose -dir $(MIG_DIR) postgres $(DB_URL) up

# Check what has been applied vs pending
goose-status:
	goose -dir $(MIG_DIR) postgres $(DB_URL) status

# Undo the last migration
goose-down:
	goose -dir $(MIG_DIR) postgres $(DB_URL) down

# Create a new migration file (Usage: make goose-create name=add_users)
goose-create:
	goose -dir $(MIG_DIR) create $(name) sql


