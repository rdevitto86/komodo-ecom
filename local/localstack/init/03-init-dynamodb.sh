#!/bin/bash

# LocalStack initialization script for DynamoDB
# Creates DynamoDB tables for user and auth data

echo "Initializing DynamoDB in LocalStack..."

sleep 1

# Users table
echo "Creating Users table..."
awslocal dynamodb create-table \
  --table-name komodo-users-dev \
  --attribute-definitions \
    AttributeName=user_id,AttributeType=S \
    AttributeName=email,AttributeType=S \
  --key-schema \
    AttributeName=user_id,KeyType=HASH \
  --global-secondary-indexes \
    "IndexName=email-index,KeySchema=[{AttributeName=email,KeyType=HASH}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  2>/dev/null || echo "Users table already exists"

# User profiles table
echo "Creating UserProfiles table..."
awslocal dynamodb create-table \
  --table-name komodo-user-profiles-dev \
  --attribute-definitions \
    AttributeName=user_id,AttributeType=S \
  --key-schema \
    AttributeName=user_id,KeyType=HASH \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  2>/dev/null || echo "UserProfiles table already exists"

# Sessions table
echo "Creating Sessions table..."
awslocal dynamodb create-table \
  --table-name komodo-sessions-dev \
  --attribute-definitions \
    AttributeName=session_id,AttributeType=S \
    AttributeName=user_id,AttributeType=S \
  --key-schema \
    AttributeName=session_id,KeyType=HASH \
  --global-secondary-indexes \
    "IndexName=user-id-index,KeySchema=[{AttributeName=user_id,KeyType=HASH}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --stream-specification \
    StreamEnabled=true,StreamViewType=NEW_AND_OLD_IMAGES \
  2>/dev/null || echo "Sessions table already exists"

# OAuth tokens table
echo "Creating OAuthTokens table..."
awslocal dynamodb create-table \
  --table-name komodo-oauth-tokens-dev \
  --attribute-definitions \
    AttributeName=token_id,AttributeType=S \
    AttributeName=user_id,AttributeType=S \
  --key-schema \
    AttributeName=token_id,KeyType=HASH \
  --global-secondary-indexes \
    "IndexName=user-id-index,KeySchema=[{AttributeName=user_id,KeyType=HASH}],Projection={ProjectionType=ALL},ProvisionedThroughput={ReadCapacityUnits=5,WriteCapacityUnits=5}" \
  --provisioned-throughput \
    ReadCapacityUnits=5,WriteCapacityUnits=5 \
  2>/dev/null || echo "OAuthTokens table already exists"

# Single-table design for user-api (PK=USER#<id>, SK=PROFILE|ADDR#<id>|PMT#<id>|PREFS).
# This is what DYNAMODB_TABLE=komodo-users points at in all environments.
echo "Creating komodo-users single table..."
awslocal dynamodb create-table \
  --table-name komodo-users \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --billing-mode PAY_PER_REQUEST \
  2>/dev/null || echo "komodo-users table already exists"

echo "DynamoDB initialized successfully"
