.DEFAULT_GOAL := bootstrap

APP_NAME := komodo-auth-api
ENV ?= local
ENV := $(strip $(ENV))
TAG ?= $(ENV)

DOCKERFILE := Dockerfile
COMPOSE_FILE := docker-compose.yaml
AWS_REGION ?= us-east-1
AWS_SECRET_PREFIX := komodo-auth-api/$(ENV)/
AWS_SECRET_BATCH := all-secrets
MEM_LIMIT ?= 512M
LOG_LEVEL ?= error
RESTART_POLICY ?= unless-stopped

ifeq ($(ENV),prod)
  AWS_ENDPOINT :=
  RESTART_POLICY := always
  MEM_LIMIT := 1g
  DISTROLESS_TAG := latest
else ifeq ($(ENV),staging)
  AWS_ENDPOINT :=
  MEM_LIMIT := 1g
  DISTROLESS_TAG := latest
else ifeq ($(ENV),dev)
  AWS_ENDPOINT := http://host.docker.internal:4566
  LOG_LEVEL := info
  DISTROLESS_TAG := debug
else
  AWS_ENDPOINT := http://host.docker.internal:4566
  RESTART_POLICY := no
  LOG_LEVEL := debug
  DISTROLESS_TAG := debug
endif

.PHONY: build run bootstrap stop restart clean test test_unit test_e2e lint

help:
	@echo "Targets:"
	@echo "  build             Build Docker images for ENV ($(ENV))"
	@echo "  run               Start containers via Docker Compose"
	@echo "  bootstrap         Build + run in one command"
	@echo "  stop              Stop and remove container stack"
	@echo "  restart           Restart container stack"
	@echo "  clean             Prune Docker and remove local build artifacts"
	@echo "  test              Run all tests"
	@echo "  test_unit         Run unit tests"
	@echo "  test_e2e          Run end-to-end tests"
	@echo "  lint              Run golangci-lint"
	@echo ""
	@echo "Supported Environments: local | dev | staging | prod"

build:
	@echo "Building $(APP_NAME):$(TAG) for ENV=$(ENV)"
	@docker build \
		-f $(DOCKERFILE) \
		-t $(APP_NAME)-public:$(TAG) \
		--build-arg BUILD_TARGET=public \
		--build-arg DISTROLESS_TAG=$(DISTROLESS_TAG) \
		..
	@docker build \
		-f $(DOCKERFILE) \
		-t $(APP_NAME)-private:$(TAG) \
		--build-arg BUILD_TARGET=private \
		--build-arg DISTROLESS_TAG=$(DISTROLESS_TAG) \
		..

run:
	@echo "Starting $(APP_NAME) for ENV=$(ENV)"
	@docker network create komodo-network 2>/dev/null || true
	@ENV=$(ENV) \
	APP_NAME=$(APP_NAME) \
	VERSION=$(TAG) \
	LOG_LEVEL=$(LOG_LEVEL) \
	AWS_REGION=$(AWS_REGION) \
	AWS_ENDPOINT=$(AWS_ENDPOINT) \
	AWS_SECRET_PREFIX=$(AWS_SECRET_PREFIX) \
	AWS_SECRET_BATCH=$(AWS_SECRET_BATCH) \
	RESTART_POLICY=$(RESTART_POLICY) \
	MEM_LIMIT=$(MEM_LIMIT) \
	DISTROLESS_TAG=$(DISTROLESS_TAG) \
	docker compose -p $(APP_NAME)-$(ENV) -f $(COMPOSE_FILE) up -d

stop:
	@echo "Stopping $(APP_NAME) for ENV=$(ENV)"
	@docker compose -p $(APP_NAME)-$(ENV) -f $(COMPOSE_FILE) down --remove-orphans

bootstrap: build run

restart: stop run

clean:
	docker container prune -f
	docker image prune -f
	docker network prune -f
	docker volume prune -f
	rm -rf bin

test:
	go test ./... -race

test_unit:
	go test -short ./...

test_e2e:
	go test -tags=e2e ./...

lint:
	@golangci-lint run ./...
