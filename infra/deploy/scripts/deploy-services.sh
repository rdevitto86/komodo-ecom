#!/bin/bash

# deploy-services.sh — Deploy or update the komodo-services-<env> CloudFormation stack.
# Requires the komodo-infra-<env> stack to already exist.
#
# Usage:   ./infra/deploy/scripts/deploy-services.sh <env> [image-tag]
# Example: ./infra/deploy/scripts/deploy-services.sh dev abc1234

set -euo pipefail

ENV=${1:-}
IMAGE_TAG=${2:-latest}

if [[ -z "$ENV" || ! "$ENV" =~ ^(dev|stg|prod)$ ]]; then
  echo "Usage: $0 <dev|stg|prod> [image-tag]"
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CFN_DIR="$REPO_ROOT/infra/deploy/cfn"
STACK_NAME="komodo-services-$ENV"

AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
AWS_REGION=${AWS_DEFAULT_REGION:-us-east-1}

echo "==> Deploying services stack: $STACK_NAME"
echo "    Image tag: $IMAGE_TAG"
echo "    Account:   $AWS_ACCOUNT_ID"
echo "    Region:    $AWS_REGION"

aws cloudformation deploy \
  --stack-name "$STACK_NAME" \
  --template-file "$CFN_DIR/services.yaml" \
  --parameter-overrides \
    "file://$CFN_DIR/parameters/$ENV.json" \
    "ImageTag=$IMAGE_TAG" \
    "AWSAccountId=$AWS_ACCOUNT_ID" \
    "AWSRegion=$AWS_REGION" \
  --capabilities CAPABILITY_NAMED_IAM \
  --no-fail-on-empty-changeset \
  --tags \
    Project=komodo \
    Environment="$ENV" \
    ManagedBy=cloudformation \
    ImageTag="$IMAGE_TAG"

echo "==> $STACK_NAME deployed successfully."
