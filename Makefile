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

up-infra:
	docker compose --profile infra up -d

up-auth:
	docker compose --profile auth up -d --build

up-backend:
	docker compose --profile backend up -d --build

up-ui:
	docker compose --profile ui-backend up -d --build

up-full:
	docker compose --profile full up -d --build

down:
	docker compose --profile full down --remove-orphans

down-v:
	docker compose --profile full down --remove-orphans --volumes

logs:
	docker compose --profile full logs -f

ps:
	docker compose --profile full ps
