# 🏦 E-Wallet Microservices

> Sistem E-Wallet berbasis Microservices dengan arsitektur Advanced

[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-24.0-2496ED?style=flat&logo=docker&logoColor=white)](https://docker.com)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.28-326CE5?style=flat&logo=kubernetes&logoColor=white)](https://kubernetes.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [Microservices](#-microservices)
- [Getting Started](#-getting-started)
- [API Documentation](#-api-documentation)
- [Deployment](#-deployment)
- [Contributing](#-contributing)
- [License](#-license)

---

## ✨ Features

### 🔐 Authentication
- User Registration & Login
- JWT Token Authentication
- Token Refresh

### 💰 Wallet Management
- Multi-Currency Support
- Real-time Balance
- Top Up via Stripe

### 💸 Transfer
- Internal Transfer (User to User)
- External Transfer (Gopay, OVO, DANA, etc.)
- Transfer Status Tracking

### 📊 Transaction History
- Complete Transaction Log
- Filter & Search
- Export to CSV/PDF

### 🔔 Notifications
- Email Notifications
- Push Notifications
- Real-time Updates

### 📈 Analytics
- Transaction Summary
- Revenue Analytics
- User Insights

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      CLIENT LAYER                           │
│                    (Web / Mobile App)                        │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    API GATEWAY (Traefik)                     │
│                  Rate Limiting, SSL, Routing                 │
└─────────────────────────┬───────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│  AUTH SERVICE │ │WALLET SERVICE │ │TRANSFER SVC   │
│  (Port: 8081) │ │ (Port: 8082)  │ │ (Port: 8083)  │
└───────┬───────┘ └───────┬───────┘ └───────┬───────┘
        │                 │                 │
        ▼                 ▼                 ▼
┌───────────────┐ ┌───────────────┐ ┌───────────────┐
│ NOTIFICATION  │ │ PAYMENT GW    │ │LEDGER SERVICE │
│ SERVICE       │ │ SERVICE       │ │ (Port: 8086)  │
│ (Port: 8084)  │ │ (Port: 8085)  │ └───────────────┘
└───────────────┘ └───────────────┘
                          │
                          ▼
                ┌───────────────┐
                │  ANALYTICS    │
                │  SERVICE      │
                │ (Port: 8087)  │
                └───────────────┘

┌─────────────────────────────────────────────────────────────┐
│                     DATA LAYER                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────┐  │
│  │ MySQL   │  │ MongoDB │  │  Redis  │  │  RabbitMQ   │  │
│  │ (Data)  │  │ (Logs)  │  │ (Cache) │  │ (Messages)  │  │
│  └─────────┘  └─────────┘  └─────────┘  └─────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 🛠️ Tech Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Backend** | Go 1.21 | Microservices implementation |
| **API Gateway** | Traefik v2.10 | Reverse proxy, load balancing |
| **Database (SQL)** | MySQL 8.0 | Users, Wallets, Transactions |
| **Database (NoSQL)** | MongoDB 6 | Logs, Notifications, Analytics |
| **Cache** | Redis 7 | Session, Rate Limiting, Cache |
| **Message Queue** | RabbitMQ 3 | Async communication |
| **Payment** | Stripe (Sandbox) | Payment processing |
| **Container** | Docker | Containerization |
| **Orchestration** | Docker Compose / K8s | Deployment |

---

## 📦 Microservices

| Service | Port | Description | Database |
|---------|------|-------------|----------|
| `auth-service` | 8081 | Authentication & Authorization | MySQL |
| `wallet-service` | 8082 | Balance & Top Up Management | MySQL + Redis |
| `transfer-service` | 8083 | Internal & External Transfers | MySQL |
| `payment-gateway` | 8084 | Stripe Integration | MySQL |
| `notification-service` | 8085 | Email & Push Notifications | MongoDB |
| `ledger-service` | 8086 | Audit Log & Transaction History | MongoDB |
| `analytics-service` | 8087 | Reporting & Analytics | MongoDB |

---

## 🚀 Getting Started

### Prerequisites

- [Docker](https://docker.com) & [Docker Compose](https://docs.docker.com/compose/)
- [Go 1.21+](https://golang.org/dl/) (for development)
- [Stripe Account](https://stripe.com) (for payment integration)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/ikhdamw/e-wallet.git
   cd e-wallet
   ```

2. **Setup environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start all services**
   ```bash
   docker-compose up -d
   ```

4. **Verify services**
   ```bash
   docker-compose ps
   ```

5. **Access services**
   - API Gateway: http://localhost:80
   - Traefik Dashboard: http://localhost:8080
   - RabbitMQ: http://localhost:15672

---

## 📡 API Documentation

### Authentication

```bash
# Register
POST /api/auth/register
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}

# Login
POST /api/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Wallet

```bash
# Get Balance
GET /api/wallet/balance
Authorization: Bearer <token>

# Top Up
POST /api/wallet/topup
Authorization: Bearer <token>
{
  "amount": 100000,
  "payment_method": "stripe"
}

# Transaction History
GET /api/wallet/history?page=1&limit=10
Authorization: Bearer <token>
```

### Transfer

```bash
# Internal Transfer
POST /api/transfer/internal
Authorization: Bearer <token>
{
  "recipient_email": "recipient@example.com",
  "amount": 50000,
  "description": "Payment for coffee"
}

# External Transfer
POST /api/transfer/external
Authorization: Bearer <token>
{
  "provider": "gopay",
  "recipient_account": "081234567890",
  "amount": 100000,
  "description": "Transfer to Gopay"
}

# Check Status
GET /api/transfer/status/:id
Authorization: Bearer <token>
```

### Payment Gateway

```bash
# Create Stripe Payment
POST /api/payment/stripe/create
Authorization: Bearer <token>
{
  "amount": 100000,
  "currency": "IDR"
}

# Stripe Webhook
POST /api/payment/stripe/webhook
```

---

## 🐳 Deployment

### Docker Compose (Development)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### Kubernetes (Production)

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check pods
kubectl get pods

# Check services
kubectl get services
```

---

## 🔧 Development

### Project Structure

```
e-wallet/
├── docker-compose.yml          # Docker Compose configuration
├── .env.example                # Environment variables template
├── README.md                   # This file
├── api-gateway/                # Traefik configuration
│   └── traefik.yml
├── k8s/                        # Kubernetes manifests
│   ├── auth-service/
│   ├── wallet-service/
│   ├── transfer-service/
│   ├── payment-gateway/
│   ├── notification-service/
│   ├── ledger-service/
│   └── analytics-service/
├── services/                   # Microservices source code
│   ├── auth-service/
│   ├── wallet-service/
│   ├── transfer-service/
│   ├── payment-gateway/
│   ├── notification-service/
│   ├── ledger-service/
│   └── analytics-service/
└── databases/                  # Database initialization
    ├── mysql/
    └── mongodb/
```

### Adding a New Service

1. Create service directory in `services/`
2. Implement service in Go
3. Create Dockerfile
4. Add service to `docker-compose.yml`
5. Add Traefik labels for routing

---

## 📊 Monitoring

### Health Check

```bash
# Check all services
curl http://localhost:8081/health  # Auth Service
curl http://localhost:8082/health  # Wallet Service
curl http://localhost:8083/health  # Transfer Service
```

### Logs

```bash
# View specific service logs
docker-compose logs -f auth-service
docker-compose logs -f wallet-service
```

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👨‍💻 Author

**Ikhda Muhammad Wildani**
- GitHub: [@ikhdamw](https://github.com/ikhdamw)
- Email: ikhda.wizards@gmail.com

---

## 🙏 Acknowledgments

- [Go](https://golang.org)
- [Docker](https://docker.com)
- [Kubernetes](https://kubernetes.io)
- [Traefik](https://traefik.io)
- [Stripe](https://stripe.com)

---

⭐ **Star this repository if you find it helpful!**
