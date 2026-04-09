# Payment Gateway gRPC

[![Go Version](https://img.shields.io/badge/Go-1.25.0-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![gRPC](https://img.shields.io/badge/gRPC-v1.77.0-4285F4?style=for-the-badge&logo=grpc)](https://grpc.io/)
[![Echo](https://img.shields.io/badge/Echo-v4.13.2-00BFFF?style=for-the-badge&logo=echo)](https://echo.labstack.com/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)

A high-performance, Monolith payment system implementation. This project demonstrates a production-ready architecture using gRPC for internal service communication and Echo for a RESTful API gateway, all wrapped in a robust observability stack.

---

## Table of Contents
- [Core Features](#core-features)
- [Technology Stack](#technology-stack)
- [Architecture](#architecture)
- [Database Schema](#database-schema)
- [Getting Started](#getting-started)
- [Project Commands (Justfile)](#project-commands-justfile)
- [Observability](#observability)
- [Testing & Performance](#testing--performance)
- [Preview](#preview)

---

## Core Features

- **Robust Authentication**: JWT-based auth with Refresh Token management.
- **Card Management**: Secure storage and retrieval of user payment cards.
- **Merchant Ecosystem**: Merchant onboarding and API key-based transaction processing.
- **Transaction Engine**: Handles complex payment flows between users and merchants.
- **Internal Transfers**: Peer-to-peer balance transfers between cards.
- **Wallet Operations**: Seamless Top-up and Withdrawal workflows.
- **Atomic Balance Management**: Consistent balance updates across all financial operations.

---

## Technology Stack

| Category | Tools |
| :--- | :--- |
| **Core** | ![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white) ![gRPC](https://img.shields.io/badge/gRPC-4285F4?style=flat-square&logo=grpc&logoColor=white) ![Echo](https://img.shields.io/badge/Echo-00BFFF?style=flat-square&logo=echo&logoColor=white) |
| **Database** | ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat-square&logo=postgresql&logoColor=white) ![Redis](https://img.shields.io/badge/Redis-DC382D?style=flat-square&logo=redis&logoColor=white) ![SQLC](https://img.shields.io/badge/SQLC-Type--Safe-blue?style=flat-square) |
| **Observability** | ![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-000000?style=flat-square&logo=opentelemetry&logoColor=white) ![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=flat-square&logo=prometheus&logoColor=white) ![Grafana](https://img.shields.io/badge/Grafana-F46800?style=flat-square&logo=grafana&logoColor=white) ![Pyroscope](https://img.shields.io/badge/Pyroscope-Flame-red?style=flat-square) |
| **Testing** | ![k6](https://img.shields.io/badge/k6-7D64FF?style=flat-square&logo=k6&logoColor=white) ![Hurl](https://img.shields.io/badge/Hurl-Testing-orange?style=flat-square) ![Testcontainers](https://img.shields.io/badge/Testcontainers-Docker-blue?style=flat-square) |
| **DevOps** | ![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white) ![Just](https://img.shields.io/badge/Just-Task--Runner-black?style=flat-square) ![Goose](https://img.shields.io/badge/Goose-Migrations-brown?style=flat-square) |

---

## Architecture

The project follows a Monolith pattern. The REST API (Echo) acts as a gateway that communicates via gRPC with internal modules. Each module is logically separated but shares a common infrastructure.

```mermaid
graph TD
    subgraph "External World"
        User((User))
    end

    subgraph "API Gateway (Echo)"
        Gateway[REST API Gateway]
    end

    subgraph "Monolith (gRPC Server)"
        direction TB
        AuthService[Auth Module]
        CardService[Card Module]
        MerchantService[Merchant Module]
        TransactionService[Transaction Module]
    end

    subgraph "Persistence"
        Postgres[(PostgreSQL)]
        Redis[(Redis Cache)]
    end

    subgraph "Observability Stack"
        OTel[OTel Collector] --> Prometheus[Prometheus]
        OTel --> Jaeger[Jaeger/Tracing]
        OTel --> Pyroscope[Pyroscope]
        Prometheus --> Grafana[Grafana]
    end

    User -- "HTTP/JSON" --> Gateway
    Gateway -- "gRPC/Protobuf" --> AuthService
    Gateway -- "gRPC/Protobuf" --> CardService
    Gateway -- "gRPC/Protobuf" --> MerchantService
    Gateway -- "gRPC/Protobuf" --> TransactionService
    
    AuthService & CardService & MerchantService & TransactionService -- "SQL" --> Postgres
    AuthService & TransactionService -- "KV" --> Redis
    
    AuthService & CardService & MerchantService & TransactionService -- "Telemetry" --> OTel
```

---

## Database Schema

```mermaid
erDiagram
    users ||--o{ cards : "has"
    users ||--o{ merchants : "has"
    users ||--o{ refresh_tokens : "has"
    users }|--|{ user_roles : "has"
    roles }|--|{ user_roles : "has"
    cards ||--o{ saldos : "owns"
    cards ||--o{ transactions : "owns"
    cards ||--o{ topups : "owns"
    cards ||--o{ withdraws : "owns"
    cards ||--o{ transfers : "from"
    cards ||--o{ transfers : "to"
    merchants ||--o{ transactions : "receives"
```

> [!NOTE]
> This schema is managed using Goose migrations and SQLC for type-safe query generation.

---

## Getting Started

### Prerequisites
- Go 1.25.0+
- Docker & Docker Compose
- Just (Alternative to Make)

### Docker Setup (Fastest)
1.  **Initialize Environment**:
    ```bash
    cp docker.env .env
    ```
2.  **Start Services**:
    ```bash
    just docker-up
    ```
    API Gateway will be available at `http://localhost:5000`.

### Local Development
1.  **Spin up Dependencies**:
    ```bash
    docker-compose up -d postgres redis
    ```
2.  **Run Migrations**:
    ```bash
    just migrate
    ```
3.  **Start Components**:
    ```bash
    # Terminal 1: gRPC Server
    just run-server
    
    # Terminal 2: API Gateway
    just run-client
    ```

---

## Project Commands (Justfile)

We use `just` for task automation. Here are some common commands:

| Command | Description |
| :--- | :--- |
| `just migrate` | Apply database migrations |
| `just generate-proto` | Generate Go code from Protobuf definitions |
| `just generate-swagger` | Generate Swagger documentation |
| `just test` | Run tests with race detection |
| `just test-all` | Run all tests including integration tests |
| `just k6 <module> <type>` | Run performance tests (e.g., `just k6 card stress`) |
| `just hurl` | Run API integration tests using Hurl |
| `just fmt` | Format Go source code |

---

## Observability

This project features a comprehensive observability stack integrated via OpenTelemetry:

- **Metrics**: Exported to Prometheus and visualized in Grafana.
- **Tracing**: Distributed tracing via OTel gRPC/HTTP instrumentation.
- **Profiling**: Continuous profiling with Pyroscope.
- **Logging**: Structured logging using Uber Zap and OTel integration.

> Visit `http://localhost:3000` for pre-configured Grafana dashboards.

---

## Testing & Performance

### Unit & Integration Testing
We use Testcontainers to spin up ephemeral PostgreSQL and Redis instances for integration tests, ensuring tests are hermetic and reliable.
```bash
just test-all
```

### Performance Benchmarks (k6)
Comprehensive performance testing across modules:

| Module | VUs | Throughput | p95 Latency | Status |
| :--- | :--- | :--- | :--- | :--- |
| **User** | 1500 | ~3500 req/s | ~400ms | Stable |
| **Role** | 1500 | ~4000 req/s | ~350ms | Stable |
| **Card** | 1500 | ~3200 req/s | ~600ms | ⚠️ Error Spike (12.5%) |

---

## Preview

### Swagger Documentation
![Swagger API Documentation](./images/swagger_3.png)

### Frontend Previews
| Web | Desktop |
| :---: | :---: |
| ![Web Frontend](./images/preview_payment.png) | ![Desktop Frontend](./images/preview_payment_desktop.png) |