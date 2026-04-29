#!/bin/bash

# DynamoDB initialization for LocalStack.
#
# Source of truth for table design (keys, GSIs, billing, streams):
#   infra/deploy/cfn/infra.yaml          — shared/platform tables (auth, user)
#   apis/<service>/docs/data-model.md    — per-service tables
#
# Table ownership — one service, one table:
#   komodo-auth-api       → komodo-sessions, komodo-oauth-tokens
#   komodo-user-api       → komodo-users  (single-table: profiles, addresses, preferences)
#   komodo-shop-items-api → komodo-shop-items    [TODO: add when data model is finalised]
#   komodo-cart-api       → komodo-cart          [TODO]
#   komodo-order-api      → komodo-orders, komodo-returns (merged)   [TODO]
#   komodo-inventory-api  → komodo-inventory     [TODO]
#   komodo-payments-api   → komodo-payments      [TODO]
#   komodo-address-api    → komodo-addresses     [TODO]
#   komodo-loyalty-api    → komodo-loyalty, komodo-reviews (merged)  [TODO]
#   komodo-event-bus-api  → komodo-events, komodo-event-subscriptions
#
# Naming: local tables have no environment suffix (komodo-users, not komodo-users-local).
# CFn appends -${Environment} for dev/stg/prod.

echo "Initializing DynamoDB in LocalStack..."

sleep 1

# ── komodo-auth-api ───────────────────────────────────────────────────────

echo "Creating Sessions table (komodo-auth-api)..."
awslocal dynamodb create-table \
  --table-name komodo-sessions \
  --attribute-definitions \
    AttributeName=session_id,AttributeType=S \
    AttributeName=user_id,AttributeType=S \
  --key-schema \
    AttributeName=session_id,KeyType=HASH \
  --global-secondary-indexes \
    "IndexName=user-id-index,KeySchema=[{AttributeName=user_id,KeyType=HASH}],Projection={ProjectionType=ALL}" \
  --billing-mode PAY_PER_REQUEST \
  --stream-specification StreamEnabled=true,StreamViewType=NEW_AND_OLD_IMAGES \
  2>/dev/null || echo "Sessions table already exists"

echo "Creating OAuthTokens table (komodo-auth-api)..."
awslocal dynamodb create-table \
  --table-name komodo-oauth-tokens \
  --attribute-definitions \
    AttributeName=token_id,AttributeType=S \
    AttributeName=user_id,AttributeType=S \
  --key-schema \
    AttributeName=token_id,KeyType=HASH \
  --global-secondary-indexes \
    "IndexName=user-id-index,KeySchema=[{AttributeName=user_id,KeyType=HASH}],Projection={ProjectionType=ALL}" \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "OAuthTokens table already exists"

# ── komodo-user-api ───────────────────────────────────────────────────────
# Single-table design. PK=USER#<id>, SK=PROFILE | ADDR#<id> | PREFS
# GSI on email supports login lookups. Stream feeds event-bus CDC Lambda.

echo "Creating Users table (komodo-user-api)..."
awslocal dynamodb create-table \
  --table-name komodo-users \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
    AttributeName=email,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --global-secondary-indexes \
    "IndexName=email-index,KeySchema=[{AttributeName=email,KeyType=HASH}],Projection={ProjectionType=ALL}" \
  --billing-mode PAY_PER_REQUEST \
  --stream-specification StreamEnabled=true,StreamViewType=NEW_AND_OLD_IMAGES \
  2>/dev/null || echo "Users table already exists"

# ── Per-service tables ────────────────────────────────────────────────────
# Add each table here as the service's data model is finalised.
# Template:
#
#   awslocal dynamodb create-table \
#     --table-name komodo-<service> \
#     --attribute-definitions ... \
#     --key-schema ... \
#     --billing-mode PAY_PER_REQUEST \
#     --stream-specification StreamEnabled=true,StreamViewType=NEW_AND_OLD_IMAGES \
#     2>/dev/null || echo "<Service> table already exists"

# ── komodo-event-bus-api ──────────────────────────────────────────────────

echo "Creating Events table (komodo-event-bus-api)..."
awslocal dynamodb create-table \
  --table-name komodo-events \
  --attribute-definitions \
    AttributeName=event_id,AttributeType=S \
    AttributeName=domain,AttributeType=S \
  --key-schema \
    AttributeName=event_id,KeyType=HASH \
    AttributeName=domain,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  --stream-specification StreamEnabled=true,StreamViewType=NEW_AND_OLD_IMAGES \
  2>/dev/null || echo "Events table already exists"

awslocal dynamodb update-time-to-live \
  --table-name komodo-events \
  --time-to-live-specification "Enabled=true,AttributeName=expires_at" \
  2>/dev/null || true

echo "Creating EventSubscriptions table (komodo-event-bus-api)..."
awslocal dynamodb create-table \
  --table-name komodo-event-subscriptions \
  --attribute-definitions \
    AttributeName=event_type,AttributeType=S \
    AttributeName=subscriber_url,AttributeType=S \
  --key-schema \
    AttributeName=event_type,KeyType=HASH \
    AttributeName=subscriber_url,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "EventSubscriptions table already exists"

echo "DynamoDB initialized successfully"
