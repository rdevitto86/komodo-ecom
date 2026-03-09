#!/bin/bash
# filepath: /Users/rad/komodo-apis/komodo-address-api/deploy/deploy_docker_staging.sh

set -e

# Set environment
ENV=staging
COMPOSE_STAGING="build/docker-compose.staging.yaml"

echo "Deploying komodo-address-api to Staging Docker cluster..."

# Build the Docker image with Staging settings
docker build -f build/Dockerfile -t komodo-address-api:${ENV} --build-arg ENV=${ENV} .

# Start the Staging stack using Docker Compose overlays
docker compose -f ${COMPOSE_STAGING} up -d --remove-orphans

echo "Staging deployment complete."
docker compose -f ${COMPOSE_STAGING} ps
