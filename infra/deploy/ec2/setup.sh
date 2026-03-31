#!/bin/bash
# setup.sh — Bootstrap a fresh Amazon Linux 2023 / Ubuntu EC2 instance.
# Run once after launch. After this, use docker compose to manage services.
#
# Usage: sudo bash setup.sh
#
# Prerequisites:
#   - EC2 instance with an IAM role that has:
#       ecr:GetAuthorizationToken, ecr:BatchGetImage, ecr:GetDownloadUrlForLayer
#   - Security group: inbound 80, 443 from 0.0.0.0/0; 22 from your IP

set -euo pipefail

AWS_REGION="${AWS_REGION:-us-east-2}"
REPO_DIR="/opt/komodo"

echo "==> Installing Docker"
if command -v apt-get &>/dev/null; then
  # Ubuntu/Debian
  apt-get update -q
  apt-get install -y docker.io docker-compose-plugin awscli certbot nginx
elif command -v dnf &>/dev/null; then
  # Amazon Linux 2023
  dnf install -y docker docker-compose-plugin awscli certbot python3-certbot-nginx
fi

systemctl enable --now docker
usermod -aG docker ec2-user 2>/dev/null || usermod -aG docker ubuntu 2>/dev/null || true

echo "==> Logging into ECR"
aws ecr get-login-password --region "$AWS_REGION" \
  | docker login --username AWS --password-stdin \
    "$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com"

echo "==> Setting up repo directory"
mkdir -p "$REPO_DIR"
cp -r . "$REPO_DIR/"

echo ""
echo "==> Setup complete."
echo ""
echo "Next steps:"
echo "  1. Copy your .env file to $REPO_DIR/deploy/ec2/.env"
echo "  2. Edit $REPO_DIR/deploy/ec2/nginx.conf — replace yourdomain.com placeholders"
echo "  3. Obtain TLS cert: certbot --nginx -d auth.yourdomain.com -d users.yourdomain.com -d items.yourdomain.com"
echo "  4. Pull images and start:"
echo "       cd $REPO_DIR"
echo "       docker compose -f deploy/ec2/docker-compose.yaml --env-file deploy/ec2/.env pull"
echo "       docker compose -f deploy/ec2/docker-compose.yaml --env-file deploy/ec2/.env up -d"
echo ""
echo "To deploy a new image tag:"
echo "       IMAGE_TAG=<sha> docker compose -f deploy/ec2/docker-compose.yaml --env-file deploy/ec2/.env up -d --no-deps <service>"
