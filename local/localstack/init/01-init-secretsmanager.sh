#!/bin/bash

echo "Initializing Secrets Manager..."

# Escape JSON values for JSON compatibility using only standard POSIX tools
json_escape() {
  printf '%s' "$1" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))' 2>/dev/null

  if [ $? -ne 0 ]; then
    printf '%s' "$1" | sed 's/\\/\\\\/g' | sed 's/"/\\"/g' | awk '{printf "%s\\n", $0}' | sed 's/\\n$//'
  fi
}

# These keys are for publically testing LocalStack only. Refrain from modifying.
PAYLOAD=$(python3 -c "
import json, os
data = {
  'AWS_ELASTICACHE_PASSWORD': 'test-password',
  'AWS_ELASTICACHE_ENDPOINT': 'redis:6379',
  'AWS_ELASTICACHE_DB': '0',
  'JWT_PUBLIC_KEY': '-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAz+2qWNNzww859lLqoVHr\n/C90EKD/lwMaRtWkgys6to5eb+QIhP4IaFouYKwGO0OXm1x9iQ7E54ZocxRWrqDH\nNoH8nUDe/V2oK+PJ3wa2NWNDd3AbeIvWDwt5kipXgFs1diNPy/KRuZvQf0AVYs2A\nmgryzl3V5CiCjU14N1bjz89sYTqWP5jXfuOIy52BiJeSiXyCjgLjSC5no0QejegP\nJ4EIo2e2OVIVbFB3CUe47c5N1xikHlHBy/IeGuAAmiorLf1RKBQRRrjeo4T9S846\nZxm0gXUPpD9frSNDKKJzJIx6/EhPILZ5gRVq3YSj8Hp+S1t1rrvShRT6nvYcbbGl\njnkYbOxMhvGBgKaElaqWLY1yov9csJ9jywiGme/yXxAshq7lTn53Kl55mcjpeWdz\ni0VGWUc4mUiy0XOV1Hh1QnHBoLwrdD7Iud3433DdDbLoMlZQZOTRJTf/rzrtQQTw\ns6ppQKWJrBjmb7F8wpyBwGLbLYdW6lW8oMuj6GjtPQPYvxup3uVJYzbC0CD9lbZS\nAgxkEng3+lcC9gIDhiKiHmlgRKEwDA1yX6JWh7E3NVzg5YJ4x+Ch8OLp9rECIsv8\nZ7EGaT0l+6ArhP16S6nWfxwfBwu1Mu8HIZPofdJ85/6AigqhKin3Xuy2SuWtb1NV\njs05OvTCWC6YRNsdxKK3SO0CAwEAAQ==\n-----END PUBLIC KEY-----\n',
  'JWT_PRIVATE_KEY': '-----BEGIN PRIVATE KEY-----\n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\nWMFdOyMfOyZvaYSHiOr0QSILPQ7RZA==\n-----END PRIVATE KEY-----\n',
  'JWT_KID': 'test-kid',
  'JWT_ISSUER': 'test-issuer',
  'JWT_AUDIENCE': 'test-audience',
  'IP_WHITELIST': '',
  'IP_BLACKLIST': '',
  'RATE_LIMIT_RPS': '100',
  'RATE_LIMIT_BURST': '200',
  'IDEMPOTENCY_TTL_SEC': '300',
  'MAX_CONTENT_LENGTH': '4096',
  'BUCKET_TTL_SECOND': '300',
  'REGISTERED_CLIENTS': json.dumps({
    'test-client': {
      'name': 'Test Service',
      'secret': 'test-secret',
      'allowed_scopes': ['svc:user-api', 'svc:order-api', 'read', 'write'],
    },
    'komodo-user-api': {
      'name': 'User API',
      'secret': 'test-secret',
      'allowed_scopes': ['svc:user-api', 'svc:auth-api', 'read', 'write'],
    },
  }),
}
print(json.dumps(data))"
)

# Auth API Secrets
awslocal secretsmanager create-secret \
  --name "komodo-auth-api/${ENV:-local}/all-secrets" \
  --description "All secrets for komodo auth api" \
  --secret-string "$PAYLOAD" || echo "Secret already exists or failed to create."

# User API Secrets
# User API Secrets — shares the same RSA key pair as auth-api so user-api can
# validate tokens. JWT_PRIVATE_KEY is included because InitializeKeys() requires
# both; the private key is not used for signing in user-api.
echo "Creating user-api secrets..."
USER_PAYLOAD=$(python3 -c "
import json, os
data = {
  'DYNAMODB_ENDPOINT': 'http://host.docker.internal:4566',
  'DYNAMODB_ACCESS_KEY': 'test',
  'DYNAMODB_SECRET_KEY': 'test',
  'DYNAMODB_TABLE': 'komodo-users',
  'USER_API_CLIENT_ID': 'test-client',
  'USER_API_CLIENT_SECRET': 'test-secret',
  'JWT_PUBLIC_KEY': '-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAz+2qWNNzww859lLqoVHr\n/C90EKD/lwMaRtWkgys6to5eb+QIhP4IaFouYKwGO0OXm1x9iQ7E54ZocxRWrqDH\nNoH8nUDe/V2oK+PJ3wa2NWNDd3AbeIvWDwt5kipXgFs1diNPy/KRuZvQf0AVYs2A\nmgryzl3V5CiCjU14N1bjz89sYTqWP5jXfuOIy52BiJeSiXyCjgLjSC5no0QejegP\nJ4EIo2e2OVIVbFB3CUe47c5N1xikHlHBy/IeGuAAmiorLf1RKBQRRrjeo4T9S846\nZxm0gXUPpD9frSNDKKJzJIx6/EhPILZ5gRVq3YSj8Hp+S1t1rrvShRT6nvYcbbGl\njnkYbOxMhvGBgKaElaqWLY1yov9csJ9jywiGme/yXxAshq7lTn53Kl55mcjpeWdz\ni0VGWUc4mUiy0XOV1Hh1QnHBoLwrdD7Iud3433DdDbLoMlZQZOTRJTf/rzrtQQTw\ns6ppQKWJrBjmb7F8wpyBwGLbLYdW6lW8oMuj6GjtPQPYvxup3uVJYzbC0CD9lbZS\nAgxkEng3+lcC9gIDhiKiHmlgRKEwDA1yX6JWh7E3NVzg5YJ4x+Ch8OLp9rECIsv8\nZ7EGaT0l+6ArhP16S6nWfxwfBwu1Mu8HIZPofdJ85/6AigqhKin3Xuy2SuWtb1NV\njs05OvTCWC6YRNsdxKK3SO0CAwEAAQ==\n-----END PUBLIC KEY-----\n',
  'JWT_PRIVATE_KEY': '-----BEGIN PRIVATE KEY-----\n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\nWMFdOyMfOyZvaYSHiOr0QSILPQ7RZA==\n-----END PRIVATE KEY-----\n',
  'JWT_KID': 'test-kid',
  'JWT_ISSUER': 'test-issuer',
  'JWT_AUDIENCE': 'test-audience',
  'IP_WHITELIST': '',
  'IP_BLACKLIST': '',
  'RATE_LIMIT_RPS': '100',
  'RATE_LIMIT_BURST': '200',
  'IDEMPOTENCY_TTL_SEC': '300',
  'MAX_CONTENT_LENGTH': '4096',
  'BUCKET_TTL_SECOND': '300',
}
print(json.dumps(data))"
)
awslocal secretsmanager create-secret \
  --name "komodo-user-api/${ENV:-local}/all-secrets" \
  --description "All secrets for komodo user api" \
  --secret-string "$USER_PAYLOAD" || echo "User API secret already exists or failed to create."

echo "Listing created secrets:"
awslocal secretsmanager list-secrets --query 'SecretList[*].Name' --output table
