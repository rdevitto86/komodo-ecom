# Komodo E-Commerce Platform

Welcome to the Komodo E-Commerce Platform! This project is a modern e-commerce web application built with Svelte and TypeScript, designed to provide a seamless shopping experience.

## Features

- User authentication and authorization
- Product browsing and search
- Shopping cart and checkout
- Order history and tracking
- Admin dashboard for managing products and orders

## Getting Started

### Prerequisites

- Node.js 18 or higher
- npm or yarn
- Docker and Docker Compose (optional, for containerized development)

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/komodo-ecom.git
   cd komodo-ecom
   ```

2. Install dependencies:
   ```
   npm install
   ```

3. Start the development server:
   ```
   npm start
   ```

4. Open your browser and navigate to `http://localhost:3000`.

---

## 🚀 Deployment Guide
This guide covers deploying your Svelte TypeScript application across three environments: Local Development, DEV, STG (Staging), and PROD.

📋 Prerequisites
- Docker and Docker Compose installed
- Node.js 18+ and npm
- AWS CLI configured (for production)
- Make utility (optional, for easier commands)

🛠️ Local Development
Quick Start
bash# Start development server
make dev
# or
docker-compose --profile dev up

# Access at: http://localhost:3000
Development Commands
bashmake dev          # Start development with hot reload
make dev-clean    # Clean restart development
npm run dev       # Run without Docker (faster)
🏠 On-Premises Deployment (Home Server)
Initial Setup
bash# Run setup script on your home server
chmod +x scripts/setup-onprem.sh
./scripts/setup-onprem.sh
Deploy to Home Server
bash# Deploy with reverse proxy
make onprem

# Or simple deployment
docker-compose --profile onprem up -d --build

Access Your App
- App: http://YOUR_SERVER_IP:8080
- Traefik Dashboard: http://YOUR_SERVER_IP:8090