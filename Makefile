.PHONY: help up-infra up-auth up-backend up-ui up-full down down-v logs ps

# Default target
help:
	@echo ""
	@echo "  Komodo Local Dev"
	@echo ""
	@echo "  make up-infra      localstack + redis only"
	@echo "  make up-auth       infra + auth-api"
	@echo "  make up-backend    infra + auth + user + shop-items"
	@echo "  make up-ui         backend + ui  (full E2E local)"
	@echo "  make up-full       everything"
	@echo ""
	@echo "  make down          stop all running services"
	@echo "  make down-v        stop all + remove volumes"
	@echo "  make logs          follow logs for all running services"
	@echo "  make ps            show running containers"
	@echo ""

COMPOSE := docker compose -f infra/local/docker-compose.yml

up-infra:
	$(COMPOSE) --profile infra up -d

up-auth:
	$(COMPOSE) --profile auth up -d --build

up-backend:
	$(COMPOSE) --profile backend up -d --build

up-ui:
	$(COMPOSE) --profile ui-backend up -d --build

up-full:
	$(COMPOSE) --profile full up -d --build

down:
	$(COMPOSE) --profile full down --remove-orphans

down-v:
	$(COMPOSE) --profile full down --remove-orphans --volumes

logs:
	$(COMPOSE) --profile full logs -f

ps:
	$(COMPOSE) --profile full ps
