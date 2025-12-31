COMPOSE_FILE ?= .docker/docker-compose.local.yaml
ENV_FILE ?= .env.local
ATLAS_CONFIG ?= file://internal/database/atlas/atlas.hcl

up:
	docker compose -f $(COMPOSE_FILE) --env-file .env.local up 

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

atlas-inspect:
	atlas schema inspect --env gorm --url "env://src" -c "$(ATLAS_CONFIG)"

atlas-migrate-diff:
	atlas migrate diff --env gorm -c "$(ATLAS_CONFIG)"

atlas-migrate-apply:
	atlas migrate apply --env local -c "$(ATLAS_CONFIG)"


