# Payment Gateway gRPC Project

This project is a **simple payment system implementation** designed to mimic the workflow of a typical digital financial service. The system is built with a modular approach, where each service can operate independently but remains interconnected through a consistent database and API.

### Core Features

The key components that form the core of this project include:

- **🔐 User Authentication**
  Supports new account registration, login with credentials, and JWT token issuance and validation. The system also includes refresh token management to maintain user session security.

- **💳 Card Management**
  Each user can add cards to their account. Card details such as number, card type, expiration date, and CVV are stored securely and can be accessed for transactions.

- **🏬 Merchant Management**
  Merchants can be created and managed through the system. Each merchant has a unique identity in the form of a UUID merchant number and an API key used to process transactions.

- **💸 Transactions**
  Handles the payment process between a user's card and a merchant. Each transaction is recorded with a unique number, payment amount, and information about the receiving merchant.

- **🔄 Fund Transfers**
  Allows users to transfer balances between cards. The system records the source (from) and destination (to) of the transfer, ensuring sufficient balance before processing the transaction.

- **📈 Balance Top-Up**
  Users can add funds to their cards through a top-up process. Each top-up has a unique identity and automatically updates the respective card's balance.

- **🏧 Withdrawals**
  Users can withdraw funds from their cards. Like top-ups, this process is recorded, and the balance is updated according to the withdrawal amount.

- **💰 Balance Management**
  All cards are linked to a balance entity. The system is responsible for tracking, adding, and subtracting balances consistently after every financial operation (transaction, transfer, top-up, or withdrawal).

---

## 🚀 Project Features

- **REST API Client**: A RESTful client that interacts with the gRPC server.
- **gRPC Server**: The main server that handles all business logic.
- **Database Migration**: Uses `goose` to manage the database schema.
- **API Documentation**: Automatically generated Swagger documentation.
- **Docker Support**: Docker and Docker Compose configurations for an easy development environment.
- **Unit & Integration Tests**: Testing to ensure code reliability.
- **CI/CD**: GitHub Actions workflows for automated builds, formatting, and testing.

## 🧰 Technology Stack

- 🐹 **Go (Golang)** — Implementation language.
- 🌐 **Echo** — A minimalist web framework for building REST APIs.
- 🪵 **Zap Logger** — Structured logging for high-performance applications.
- 📦 **SQLC** — Generates type-safe Go code from SQL queries.
- 🚀 **gRPC** — High-performance RPC for internal service communication.
- 🧳 **Goose** — A migration tool for managing database schema changes.
- 🐳 **Docker** — A containerization platform for consistent development environments.
- 📄 **Swago** — Generates Swagger 2.0 documentation for Echo routes.
- 🔗 **Docker Compose** — Manages multi-container Docker applications.

---

## Architecture

This application is designed with a service-oriented monolithic architecture. The client-facing REST API acts as a gateway, translating HTTP requests into gRPC calls to the backend server. This server contains the core business logic and communicates with a PostgreSQL database.

```mermaid
graph TD
    subgraph "User Interaction"
        User -- "HTTP/REST (JSON)" --> Client[Client/API Gateway]
    end

    subgraph "Application Services"
        Client -- "gRPC (Protobuf)" --> Server[gRPC Server]
        Server -- "SQL" --> Database[(PostgreSQL)]
    end

    subgraph "Development & Operations"
        Migration[Migration Process] -- "SQL" --> Database
    end

    style Client fill:#d3869b,stroke:#3c3836,stroke-width:2px,color:#282828
    style Server fill:#83a598,stroke:#3c3836,stroke-width:2px,color:#282828
    style Database fill:#b8bb26,stroke:#3c3836,stroke-width:2px,color:#282828
    style Migration fill:#fe8019,stroke:#3c3836,stroke-width:2px,color:#282828
```

---

## Database Schema (ERD)

The following diagram illustrates the relationships between the tables in the database.

```mermaid
erDiagram
    users {
        INT user_id PK
        VARCHAR firstname
        VARCHAR lastname
        VARCHAR email
        VARCHAR password
    }

    roles {
        INT role_id PK
        VARCHAR role_name
    }

    user_roles {
        INT user_role_id PK
        INT user_id FK
        INT role_id FK
    }

    cards {
        INT card_id PK
        INT user_id FK
        VARCHAR card_number
        VARCHAR card_type
        DATE expire_date
        VARCHAR cvv
    }

    merchants {
        INT merchant_id PK
        UUID merchant_no
        VARCHAR name
        VARCHAR api_key
        INT user_id FK
    }

    saldos {
        INT saldo_id PK
        VARCHAR card_number FK
        INT total_balance
    }

    transactions {
        INT transaction_id PK
        UUID transaction_no
        VARCHAR card_number FK
        INT amount
        INT merchant_id FK
    }

    transfers {
        INT transfer_id PK
        UUID transfer_no
        VARCHAR transfer_from FK
        VARCHAR transfer_to FK
        INT transfer_amount
    }

    topups {
        INT topup_id PK
        UUID topup_no
        VARCHAR card_number FK
        INT topup_amount
    }

    withdraws {
        INT withdraw_id PK
        UUID withdraw_no
        VARCHAR card_number FK
        INT withdraw_amount
    }

    refresh_tokens {
        INT refresh_token_id PK
        INT user_id FK
        VARCHAR token
    }

    %% Relationships
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

---

## Getting Started

You can run this project either locally with a Go environment or using Docker.

### Prerequisites

- Go (version 1.21 or newer)
- Docker and Docker Compose
- The `make` command-line tool
- An `.env` file with the necessary environment variables. You can copy from `docker.env` as a template.

### 1. Clone the Repository

```bash
git clone https://github.com/hoover/payment-gateway-grpc.git
cd payment-gateway-grpc
```

### 2. Running with Docker (Recommended)

This is the easiest way to run all services.

1.  **Create Environment File:**
    Copy `docker.env` to a new `.env` file.

    ```bash
    cp docker.env .env
    ```

2.  **Build and Run Services:**
    Use the `make` command to build the images and start the containers in detached mode.

    ```bash
    make docker-up
    ```

    This will start the `postgres`, `migrate`, `server`, and `client` services. The client will be available at `http://localhost:5000`.

3.  **Stopping Services:**
    To stop all running containers, use:
    ```bash
    make docker-down
    ```

### 3. Running Locally

1.  **Start the Database:**
    You can use the provided Docker Compose file to run only the PostgreSQL database.

    ```bash
    docker-compose up -d postgres
    ```

2.  **Set Up Environment:**
    Create an `.env` file in the root directory and fill in the database connection details and other required variables.

3.  **Run Database Migrations:**
    Apply the latest database schema.

    ```bash
    make migrate
    ```

4.  **Run Services:**
    Open two separate terminal windows.

    In the first terminal, run the gRPC server:

    ```bash
    make run-server
    ```

    In the second terminal, run the client:

    ```bash
    make run-client
    ```

    The client API will be accessible at `http://localhost:5000`.

---

## Available `make` Commands

- `migrate`: Applies database migrations.
- `migrate-down`: Reverts the last database migration.
- `run-server`: Starts the gRPC server locally.
- `run-client`: Starts the client API gateway locally.
- `docker-up`: Builds and starts all services with Docker Compose.
- `docker-down`: Stops and removes all services started with Docker Compose.
- `test`: Runs unit tests.
- `test-all`: Runs all tests (unit and integration).
- `fmt`: Formats the Go source code.
- `lint`: Lints the codebase.
- `generate-proto`: Generates Go code from Protobuf files.

---

## Preview

### Swagger API Documentation

<img src="./images/swagger_3.png" alt="swagger" />

### Frontend Preview

#### Web

<img src="./images/preview_payment.png" alt="Web Frontend Preview" />

#### Desktop

<img src="./images/preview_payment_desktop.png" alt="Desktop Frontend Preview" />

## Performance & Scability Summary (k6)

### User Module

| Test Type  | VUs  | Throughput (req/s) | Error Rate | p95 Latency | Notes                             |
| ---------- | ---- | ------------------ | ---------- | ----------- | --------------------------------- |
| Smoke      | 1    | –                  | 0%         | <10 ms      | Baseline validated                |
| Capability | 150  | ~960               | ~0%        | ~6–7 ms     | Highly efficient, CPU-light       |
| Load       | 1000 | ~3800              | 0%         | ~386 ms     | Tail latency increases but stable |
| Stress     | 1500 | ~3560              | ~0%        | ~408 ms     | Graceful degradation              |
| Spike      | 1000 | ~3235              | 0%         | ~336 ms     | Clean and fast recovery           |

### Role Module

| Test Type  | VUs  | Throughput (req/s) | Error Rate | p95 Latency | Notes                             |
| ---------- | ---- | ------------------ | ---------- | ----------- | --------------------------------- |
| Smoke      | 1    | –                  | 0%         | <10 ms      | Baseline validated                |
| Capability | 154  | ~1200              | 0%         | ~7.5 ms     | Similar to User, slightly heavier |
| Load       | 1000 | ~4200              | ~0%        | ~333 ms     | p95 exceeds soft threshold        |
| Stress     | 1500 | ~3980              | ~0%        | ~349 ms     | Stable under sustained pressure   |
| Spike      | 1000 | ~3850              | 0%         | ~271 ms     | Fast recovery after spike         |

### Card Module

| Test Type  | VUs  | Throughput (req/s) | Error Rate | p95 Latency | Notes                                 |
| ---------- | ---- | ------------------ | ---------- | ----------- | ------------------------------------- |
| Smoke      | 1    | –                  | 0%         | <100 ms     | Acceptable under low load             |
| Capability | 900  | ~2240              | **12.5%**  | ~797 ms     | Structural failure observed           |
| Load       | 1000 | ~2600              | **12.5%**  | ~806 ms     | Latency and error thresholds breached |
| Stress     | 1500 | ~3230              | **12.5%**  | ~587 ms     | Consistent error pattern              |
| Spike      | 1000 | ~3088              | **12.5%**  | ~360 ms     | Fast failure, consistent behavior     |

---

## Performance Test Visualizations

### Card Module

| Test Type       | Visualization                                                                                                              |
| --------------- | -------------------------------------------------------------------------------------------------------------------------- |
| Capability Test | <img src="./images/load_test/card/card_capability.png" alt="Card module capability test results" width="500" />            |
| Load Test       | <img src="./images/load_test/card/card_load_test.png" alt="Card module load test results near capacity" width="500" />     |
| Stress Test     | <img src="./images/load_test/card/card_rampling.png" alt="Card module stress test behavior beyond capacity" width="500" /> |
| Spike Test      | <img src="./images/load_test/card/card_spike.png" alt="Card module spike test results" width="500" />                      |

### Role Module

| Test Type       | Visualization                                                                                                              |
| --------------- | -------------------------------------------------------------------------------------------------------------------------- |
| Capability Test | <img src="./images/load_test/role/role_capability.png" alt="Role module capability test results" width="500" />            |
| Load Test       | <img src="./images/load_test/role/role_load_test.png" alt="Role module load test results near capacity" width="500" />     |
| Stress Test     | <img src="./images/load_test/role/role_rampling.png" alt="Role module stress test behavior beyond capacity" width="500" /> |
| Spike Test      | <img src="./images/load_test/role/role_spike.png" alt="Role module spike test results" width="500" />                      |

### User Module

| Test Type       | Visualization                                                                                                              |
| --------------- | -------------------------------------------------------------------------------------------------------------------------- |
| Capability Test | <img src="./images/load_test/user/user_capability.png" alt="User module capability test results" width="500" />            |
| Load Test       | <img src="./images/load_test/user/user_load_test.png" alt="User module load test results near capacity" width="500" />     |
| Stress Test     | <img src="./images/load_test/user/user_rampling.png" alt="User module stress test behavior beyond capacity" width="500" /> |
| Spike Test      | <img src="./images/load_test/user/user_spike.png" alt="User module spike test results" width="500" />                      |