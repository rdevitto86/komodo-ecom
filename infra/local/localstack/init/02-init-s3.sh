#!/bin/bash

# LocalStack initialization script for S3
# Creates S3 buckets for file storage and assets

echo "Initializing S3 in LocalStack..."

sleep 1

# Create buckets
echo "Creating S3 buckets..."

# User uploads bucket
awslocal s3 mb s3://komodo-user-uploads-dev 2>/dev/null || echo "Bucket komodo-user-uploads-dev already exists"
awslocal s3api put-bucket-cors \
  --bucket komodo-user-uploads-dev \
  --cors-configuration '{
    "CORSRules": [{
      "AllowedOrigins": ["*"],
      "AllowedMethods": ["GET", "PUT", "POST", "DELETE"],
      "AllowedHeaders": ["*"],
      "MaxAgeSeconds": 3000
    }]
  }' 2>/dev/null

# Static assets bucket
awslocal s3 mb s3://komodo-static-assets-dev 2>/dev/null || echo "Bucket komodo-static-assets-dev already exists"
awslocal s3api put-bucket-policy \
  --bucket komodo-static-assets-dev \
  --policy '{
    "Version": "2012-10-17",
    "Statement": [{
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::komodo-static-assets-dev/*"
    }]
  }' 2>/dev/null

# Backups bucket
awslocal s3 mb s3://komodo-backups-dev 2>/dev/null || echo "Bucket komodo-backups-dev already exists"

echo "S3 initialized successfully"
