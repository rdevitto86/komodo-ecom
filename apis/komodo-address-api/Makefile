# Default to Docker-first workflow: `make` -> build image + up (dev)
.DEFAULT_GOAL := bootstrap

# === Service metadata ===
APP := komodo-address-api
ENV ?= dev
ENV := $(strip $(ENV))
TAG ?= $(ENV)

# Paths
BUILD_DIR := build

# Compose (base + overlay)
COMPOSE_FILE := $(BUILD_DIR)/docker-compose.$(ENV).yaml

.PHONY: local build run bootstrap stop deploy_staging deploy_prod restart clean test_all test_unit test_int test_e2e fmt

help:
	@echo "Targets:"
	@echo "  build             Build Docker image for ENV ($(ENV))"
	@echo "  run               Build image and start container via Docker Compose"
	@echo "  local             Run app locally (no Docker, loads .env)"
	@echo "  stop              Stop and remove container stack"
	@echo "  restart           Restart container stack"
	@echo "  clean             Prune Docker and remove local build artifacts"
	@echo "  test_all          Run all tests"
	@echo "  test_unit         Run unit tests"
	@echo "  test_int          Run integration tests"
	@echo "  test_e2e          Run end-to-end tests"
	@echo "  fmt               Run go fmt and go vet"
	@echo ""
	@echo "Examples:"
	@echo "  make build ENV=prod"
	@echo "  make run ENV=dev"
	@echo "  make local ENV=dev"
	@echo "  make stop ENV=prod"
	@echo "  make restart ENV=dev"
	@echo "  make clean"
	@echo "  make test_all"

# Local Docker Commands
build:
	docker build -f $(BUILD_DIR)/Dockerfile -t $(APP):$(TAG) --build-arg ENV=$(ENV) .

run:
	docker compose -f $(COMPOSE_FILE) up -d

stop:
	docker compose -f $(COMPOSE_FILE) down --remove-orphans

bootstrap:
	$(MAKE) build
	$(MAKE) run

restart:
	$(MAKE) stop ENV=$(ENV)
	$(MAKE) run ENV=$(ENV)

clean:
	docker container prune -f
	docker image prune -f
	docker network prune -f
	docker volume prune -f
	rm -rf bin

# ENV Deploy Commands
deploy_dev:
	$(MAKE) bootstrap ENV=dev

deploy_staging:
	bash deploy/staging/deploy_docker_staging.sh

deploy_prod:
	bash deploy/prod/deploy_docker_prod.sh

# Test Commands
test_all:
	go test ./... -race

test_unit:
	go test -short ./...

test_int:
	go test -tags=integration ./...

test_e2e:
	go test -tags=e2e ./...

fmt:
	go fmt ./... && go vet ./...
