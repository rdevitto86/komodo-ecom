#!/bin/bash

# deploy-infra.sh — Deploy or update the komodo-infra-<env> CloudFormation stack.
# Usage: ./infra/deploy/scripts/deploy-infra.sh <env>
# Example: ./infra/deploy/scripts/deploy-infra.sh dev

set -euo pipefail

ENV=${1:-}
if [[ -z "$ENV" || ! "$ENV" =~ ^(dev|stg|prod)$ ]]; then
  echo "Usage: $0 <dev|stg|prod>"
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CFN_DIR="$REPO_ROOT/infra/deploy/cfn"
STACK_NAME="komodo-infra-$ENV"

echo "==> Deploying infra stack: $STACK_NAME"

aws cloudformation deploy \
  --stack-name "$STACK_NAME" \
  --template-file "$CFN_DIR/infra.yaml" \
  --parameter-overrides "file://$CFN_DIR/parameters/$ENV.json" \
  --capabilities CAPABILITY_NAMED_IAM \
  --no-fail-on-empty-changeset \
  --tags \
    Project=komodo \
    Environment="$ENV" \
    ManagedBy=cloudformation

echo "==> $STACK_NAME deployed successfully."
echo ""
echo "Stack outputs:"
aws cloudformation describe-stacks \
  --stack-name "$STACK_NAME" \
  --query "Stacks[0].Outputs" \
  --output table
