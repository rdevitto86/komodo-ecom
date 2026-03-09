#!/bin/bash
# local-run.sh — thin wrapper around the root docker-compose profiles.
# Prefer using `make` directly. See Makefile for all available targets.
#
# Usage: ./local-run.sh [start|stop|restart] [profile]
#
# Profiles:
#   infra      → localstack + redis
#   auth       → infra + auth-api
#   backend    → infra + auth + user + shop-items  (default)
#   ui         → backend + ui
#   full       → everything
#
# Examples:
#   ./local-run.sh start
#   ./local-run.sh start ui
#   ./local-run.sh stop

set -euo pipefail

PROFILE="${2:-backend}"

# Map friendly names to compose profile names
case "$PROFILE" in
  ui)      COMPOSE_PROFILE="ui-backend" ;;
  *)       COMPOSE_PROFILE="$PROFILE" ;;
esac

function start_stack() {
  echo "==> Starting profile: $COMPOSE_PROFILE"
  docker compose --profile "$COMPOSE_PROFILE" up -d --build

  echo ""
  echo "Stack is up. Active ports depend on profile:"
  echo "  LocalStack:     http://localhost:4566"
  echo "  Redis:          localhost:6379"
  echo "  Auth API:       http://localhost:7011  (internal: 7012)"
  echo "  User API:       http://localhost:7051  (internal: 7052)"
  echo "  Shop Items API: http://localhost:7041"
  echo "  UI:             http://localhost:7001"
}

function stop_stack() {
  echo "==> Stopping all services..."
  docker compose --profile full down --remove-orphans
}

function usage() {
  echo "Usage: $0 [start|stop|restart] [infra|auth|backend|ui|full]"
  exit 1
}

[[ $# -lt 1 ]] && usage

case "$1" in
  start)   start_stack ;;
  stop)    stop_stack ;;
  restart) stop_stack && start_stack ;;
  *)       usage ;;
esac
