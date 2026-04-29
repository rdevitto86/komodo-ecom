#!/bin/bash

echo "Initializing S3 in LocalStack..."

sleep 1

echo "Creating S3 buckets..."

# User uploads
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

# Static assets (public read)
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

# Shop item images (public read, CORS for browser uploads)
awslocal s3 mb s3://komodo-shop-items-assets 2>/dev/null || echo "Bucket komodo-shop-items-assets already exists"
awslocal s3api put-bucket-cors \
  --bucket komodo-shop-items-assets \
  --cors-configuration '{
    "CORSRules": [{
      "AllowedOrigins": ["*"],
      "AllowedMethods": ["GET", "PUT"],
      "AllowedHeaders": ["*"],
      "MaxAgeSeconds": 3000
    }]
  }' 2>/dev/null
awslocal s3api put-bucket-policy \
  --bucket komodo-shop-items-assets \
  --policy '{
    "Version": "2012-10-17",
    "Statement": [{
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::komodo-shop-items-assets/*"
    }]
  }' 2>/dev/null

# Email templates (private, read-only for communications-api)
awslocal s3 mb s3://komodo-email-templates 2>/dev/null || echo "Bucket komodo-email-templates already exists"

# Backups
awslocal s3 mb s3://komodo-backups-dev 2>/dev/null || echo "Bucket komodo-backups-dev already exists"

echo "S3 initialized successfully"
