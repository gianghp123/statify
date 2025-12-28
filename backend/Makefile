COMPOSE_FILE ?= .docker/docker-compose.local.yaml
ENV_FILE ?= .env.local

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
