# TapX Lite — Offline Payment Sync Backend

A production-grade backend system that simulates offline financial transactions and reconciles them with a central ledger. Built with Go, PostgreSQL, Redis, and Docker. Designed to reflect real-world engineering challenges including idempotency, atomicity, and offline-first transaction handling.

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [System Architecture](#2-system-architecture)
3. [Technology Stack](#3-technology-stack)
4. [Project Structure](#4-project-structure)
5. [Data Models](#5-data-models)
6. [API Reference](#6-api-reference)
7. [Core Engineering Decisions](#7-core-engineering-decisions)
8. [Database Schema](#8-database-schema)
9. [Configuration](#9-configuration)
10. [Local Development](#10-local-development)
11. [Docker Deployment](#11-docker-deployment)
12. [Code Walkthrough](#12-code-walkthrough)
13. [Testing](#13-testing)
14. [Known Limitations and Future Work](#14-known-limitations-and-future-work)

---

## 1. Project Overview

### Problem Statement

Mobile payment systems in Bangladesh (bKash, Nagad, Rocket) are entirely dependent on live internet connectivity. When connectivity drops — in rural markets, crowded urban areas, or building basements — transactions fail. There is no offline-first payment infrastructure available for institutional closed-loop environments such as university canteens, corporate cafeterias, or cooperative retail.

### What TapX Lite Solves

TapX Lite is the backend engine for an offline-first payment system. Devices queue transactions locally when offline. When connectivity is restored, the device submits the batch to this backend for reconciliation. The backend validates each transaction atomically — checking for duplicate submissions and sufficient balance — and updates the ledger accordingly.

### Scope

This repository covers the backend sync service only. It is intentionally scoped to demonstrate:

- Offline transaction reconciliation via a batch sync API
- Idempotency using nonce-based duplicate detection
- Atomic balance updates using database-level transactions
- Production-ready Go service with proper configuration, middleware, and graceful shutdown
- Containerized deployment using Docker and Docker Compose

---

## 2. System Architecture

### High-Level Overview

```
Offline Device (Android / Simulated Client)
            │
            │  POST /sync  (batch of queued transactions)
            ▼
  ┌─────────────────────┐
  │   Go REST API        │  ← TapX Lite Backend (this repo)
  │   port 8080          │
  └────────┬────────────┘
           │
     ┌─────┴──────┐
     │            │
     ▼            ▼
PostgreSQL      Redis
  port 5432    port 6379
(ledger +      (future: nonce
 accounts)      cache layer)
```

### Layered Architecture

The codebase follows a strict layered architecture where each layer has a single responsibility and communicates only with the layer directly below it.

```
HTTP Request
     │
     ▼
┌──────────────────────────────────────┐
│  Handler Layer  (internal/handler)   │  Receives HTTP requests, validates input,
│                                      │  returns HTTP responses. No business logic.
└──────────────────┬───────────────────┘
                   │
                   ▼
┌──────────────────────────────────────┐
│  Service Layer  (internal/service)   │  Owns all business logic. Orchestrates
│                                      │  repository calls. Makes accept/reject decisions.
└──────────────────┬───────────────────┘
                   │
                   ▼
┌──────────────────────────────────────┐
│  Repository Layer (internal/repo)    │  Database access only. No business logic.
│                                      │  All SQL queries live here.
└──────────────────┬───────────────────┘
                   │
                   ▼
┌──────────────────────────────────────┐
│  Model Layer    (internal/model)     │  Pure data structures shared across all layers.
│                                      │  No methods, no logic.
└──────────────────────────────────────┘
```

### Transaction Processing Flow

```
POST /sync received
        │
        ▼
Decode and validate request body
        │
        ▼
For each transaction in batch:
        │
        ├──▶ BEGIN database transaction
        │
        ├──▶ INSERT transaction record
        │         │
        │         ├── UNIQUE nonce constraint violated?
        │         │         └──▶ ROLLBACK → mark rejected
        │         │
        │         └── Insert succeeds → continue
        │
        ├──▶ UPDATE accounts SET balance = balance - amount
        │         WHERE id = user_id AND balance >= amount
        │         │
        │         ├── 0 rows affected (insufficient balance)?
        │         │         └──▶ ROLLBACK → mark rejected
        │         │
        │         └── 1 row affected → continue
        │
        ├──▶ COMMIT
        │         └──▶ mark accepted
        │
        └── Next transaction
```

---

## 3. Technology Stack

| Layer | Technology | Version | Reason |
|-------|-----------|---------|--------|
| Language | Go | 1.26 | Performance, strong stdlib, ideal for DevOps and fintech backends |
| Database | PostgreSQL | 15 | ACID compliance, strong constraint support, production standard |
| Cache | Redis | 7 | Future nonce caching layer; included for extensibility |
| Containerization | Docker | 29.3.1 | Environment consistency across dev and production |
| Orchestration | Docker Compose | v5.1.1 | Local multi-container management |
| CI/CD | GitHub Actions | — | Automated testing and build on push |
| Cloud | AWS EC2 | — | Target deployment environment (Phase 2) |

---

## 4. Project Structure

```
Tapx-lite/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point. Wires all layers together.
│                                # Owns server lifecycle: config, DB, middleware, shutdown.
│
├── internal/
│   ├── handler/
│   │   └── transaction.go       # HTTP handlers for /health, /sync, /balance.
│   │                            # Validates input. Delegates to service. Returns responses.
│   │
│   ├── service/
│   │   └── transaction.go       # Business logic for sync and balance operations.
│   │                            # Owns accept/reject decisions. Manages DB transactions.
│   │
│   ├── repository/
│   │   └── transaction.go       # All PostgreSQL queries. No logic beyond data access.
│   │                            # Exposes BeginTx for atomic multi-step operations.
│   │
│   └── model/
│       └── transaction.go       # Shared data structures: Transaction, SyncRequest,
│                                # SyncResponse. No methods or business logic.
│
├── migrations/
│   └── 001_init.sql             # Initial schema: accounts and transactions tables.
│                                # Includes seed data for local development.
│
├── docs/
│   ├── dev-environment-setup.md # Full environment setup guide for Ubuntu 24.04.
│   └── day1-work-log.md         # Day 1 progress log.
│
├── Dockerfile                   # Multi-stage build: golang:1.26-alpine → alpine:latest.
├── docker-compose.yml           # Local stack: app + postgres + redis with healthchecks.
├── .env                         # Environment variables. Never committed to Git.
├── .gitignore                   # Excludes .env and build artifacts.
├── go.mod                       # Go module definition.
├── go.sum                       # Dependency checksums.
└── README.md                    # This file.
```

---

## 5. Data Models

### Transaction

Represents a single offline payment record submitted for reconciliation.

```go
type Transaction struct {
    ID        string    `json:"id"`
    UserID    int       `json:"user_id"`
    Amount    int       `json:"amount"`    // stored in smallest currency unit (paisa)
    Nonce     string    `json:"nonce"`     // unique one-time value for idempotency
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}
```

**Design note:** `Amount` is stored as an integer in the smallest currency unit (paisa) to avoid floating-point precision errors. 100 taka = 10000 paisa.

### SyncRequest

```go
type SyncRequest struct {
    Transactions []Transaction `json:"transactions"`
}
```

### SyncResponse

```go
type SyncResponse struct {
    Accepted []string `json:"accepted"`
    Rejected []string `json:"rejected"`
}
```

---

## 6. API Reference

### GET /health

Returns server liveness status. Used by load balancers and container orchestrators.

**Response 200**
```json
{
  "status": "ok",
  "service": "tapx-lite"
}
```

---

### POST /sync

Accepts a batch of offline transactions and reconciles them against the ledger. Each transaction is evaluated independently using an atomic database transaction.

**Request Body**
```json
{
  "transactions": [
    {
      "id": "tx001",
      "user_id": 1,
      "amount": 500,
      "nonce": "device-uuid-timestamp-001"
    }
  ]
}
```

**Response 200**
```json
{
  "accepted": ["tx001"],
  "rejected": null
}
```

**Rejection reasons (not exposed in response, logged server-side):**
- Duplicate nonce: transaction was already processed
- Insufficient balance: account balance is less than transaction amount
- Database error: internal failure during atomic operation

---

### GET /balance?user_id={id}

Returns the current account balance for the specified user.

**Response 200**
```json
{
  "user_id": 1,
  "balance": 9500
}
```

**Response 404**
```json
account not found for user 99
```

---

## 7. Core Engineering Decisions

### 7.1 Nonce-Based Idempotency

**Problem:** An offline device may queue the same transaction twice (e.g., double tap, retry logic). The backend must detect and reject duplicates without relying on application-level checks alone.

**Solution:** A `UNIQUE` constraint on the `nonce` column in the `transactions` table. When `InsertTransactionTx` is called with a duplicate nonce, PostgreSQL raises a constraint violation error. The service layer catches this, rolls back the database transaction, and marks the payment as rejected.

This gives us two independent layers of duplicate detection: application-level nonce check (future Redis cache) and database-level constraint enforcement.

### 7.2 Atomic Balance Updates

**Problem:** A naive implementation would read the balance, check it in application code, then update it in a separate query. This creates a race condition — two concurrent requests for the same user could both pass the balance check before either deducts.

**Solution:** The balance deduction is performed in a single atomic SQL statement:

```sql
UPDATE accounts
SET balance = balance - $1
WHERE id = $2 AND balance >= $1
```

The `WHERE balance >= $1` guard means the update only executes if the balance is sufficient. If zero rows are affected, the service layer interprets this as insufficient funds and rolls back. No separate SELECT is needed. No race condition is possible.

### 7.3 Database Transactions Per Payment

**Problem:** Inserting a transaction record and updating the account balance are two separate database operations. If the process crashes between them, the ledger becomes inconsistent.

**Solution:** Each payment is wrapped in a PostgreSQL database transaction (`BEGIN` / `COMMIT` / `ROLLBACK`). Both operations succeed together or neither does. The `defer dbTx.Rollback()` pattern in `processSingle` ensures rollback always occurs on any failure path, even panics.

### 7.4 processSingle Extraction (defer-in-loop fix)

**Problem:** Using `defer` inside a `for` loop means deferred calls execute at the end of the enclosing function, not at the end of each loop iteration. In a batch of 100 transactions, all 100 database connections would remain open until the entire function returned.

**Solution:** The per-transaction logic was extracted into a private `processSingle` method. `defer dbTx.Rollback()` now executes at the end of each `processSingle` call, releasing the database connection immediately after each transaction is processed.

### 7.5 Multi-Stage Docker Build

**Problem:** A Go application built inside a `golang` base image produces a final image of ~300MB, most of which is the Go toolchain — unnecessary at runtime.

**Solution:** A two-stage Dockerfile. Stage 1 uses `golang:1.26-alpine` to compile a statically linked binary (`CGO_ENABLED=0`). Stage 2 copies only that binary into a clean `alpine:latest` image (~5MB). The Go compiler, source code, and build tools are discarded.

### 7.6 Environment Validation at Startup

**Problem:** A missing environment variable (e.g., `DB_PASSWORD`) causes a runtime failure deep inside the application, often with a cryptic error.

**Solution:** `loadConfig()` reads all environment variables at startup and explicitly checks each required field. If any are missing, the application logs a clear error listing all missing variables and exits immediately before attempting any connections. Fast failure with a clear message is preferable to a silent misconfiguration.

### 7.7 Graceful Shutdown

**Problem:** Killing a server process immediately (`SIGKILL`) drops in-flight requests. For a payment system, an interrupted request mid-commit could leave the ledger in an inconsistent state.

**Solution:** The server listens for `SIGINT` and `SIGTERM` signals. On receipt, `server.Shutdown(ctx)` is called with a 10-second context. This stops accepting new connections and waits for all in-flight requests to complete before the process exits.

---

## 8. Database Schema

```sql
-- accounts: stores user balances for the closed-loop payment system.
CREATE TABLE IF NOT EXISTS accounts (
    id         SERIAL PRIMARY KEY,
    balance    INT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- transactions: immutable ledger of all processed payments.
CREATE TABLE IF NOT EXISTS transactions (
    id         TEXT PRIMARY KEY,
    user_id    INT NOT NULL REFERENCES accounts(id),
    amount     INT NOT NULL CHECK (amount > 0),
    nonce      TEXT NOT NULL UNIQUE,
    status     TEXT NOT NULL DEFAULT 'accepted',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

**Key constraints:**

`CHECK (balance >= 0)` — third layer of protection against negative balances, enforced at the database level regardless of application logic.

`CHECK (amount > 0)` — prevents zero or negative amount transactions from being stored.

`UNIQUE (nonce)` — enforces idempotency at the database level. Duplicate transaction submissions fail with a constraint violation before any balance is touched.

`REFERENCES accounts(id)` — foreign key integrity. A transaction cannot reference a non-existent account.

---

## 9. Configuration

All configuration is driven by environment variables. No values are hardcoded.

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_HOST` | Yes | — | PostgreSQL host |
| `DB_PORT` | Yes | — | PostgreSQL port |
| `DB_USER` | Yes | — | PostgreSQL username |
| `DB_PASSWORD` | Yes | — | PostgreSQL password |
| `DB_NAME` | Yes | — | PostgreSQL database name |
| `DB_SSLMODE` | No | `disable` | SSL mode. Use `require` in production |
| `SERVER_PORT` | No | `8080` | HTTP server port |

**Local development `.env`:**
```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tapx
DB_SSLMODE=disable
SERVER_PORT=8080
```

**Security note:** `.env` is listed in `.gitignore` and must never be committed to version control.

---

## 10. Local Development

### Prerequisites

- Ubuntu 24.04 LTS
- Go 1.26+
- Docker 28+
- Docker Compose v2+
- Git

### Setup

```bash
# Clone the repository
git clone git@github.com:Shihab369/Tapx-lite.git
cd Tapx-lite

# Copy environment file
cp .env.example .env
# Edit .env with your values

# Install Go dependencies
go mod download

# Start infrastructure (PostgreSQL + Redis)
docker compose up -d postgres redis

# Run database migrations
docker exec -i tapx-postgres psql -U postgres -d tapx < migrations/001_init.sql

# Run the application locally
go run ./cmd/server
```

### Verify

```bash
curl http://localhost:8080/health
# {"service":"tapx-lite","status":"ok"}
```

---

## 11. Docker Deployment

### Build Image

```bash
docker build -t tapx-lite:latest .
```

### Run Full Stack

```bash
docker compose up -d
```

This starts three containers:

| Container | Image | Port | Role |
|-----------|-------|------|------|
| `tapx-postgres` | postgres:15-alpine | 5432 | Primary database |
| `tapx-redis` | redis:7 | 6379 | Cache layer |
| `tapx-app` | tapx-lite:latest | 8080 | Go API server |

The `app` service uses `depends_on` with `condition: service_healthy` to ensure PostgreSQL and Redis pass their healthchecks before the application starts.

### Check Status

```bash
docker compose ps
docker logs tapx-app
```

Expected startup logs:
```
configuration loaded
database connected
tapx-lite listening on port 8080
```

### Teardown

```bash
docker compose down
```

---

## 12. Code Walkthrough

### cmd/server/main.go

The application entry point. Responsible for the full server lifecycle.

```
loadConfig()
```
Reads all environment variables into a `Config` struct. Validates that required fields are present. Fails fast with a descriptive error if any are missing. This prevents silent misconfigurations from reaching production.

```
connectDB(cfg)
```
Opens a PostgreSQL connection using the `lib/pq` driver. Configures the connection pool with production-appropriate settings: 25 maximum open connections, 10 idle connections, 5-minute connection lifetime. Verifies connectivity with a context-bound ping that times out after 5 seconds.

```
loggingMiddleware
```
Wraps every HTTP handler. Records method, path, response duration, and client IP for every request. Essential for production observability.

```
jsonMiddleware
```
Sets `Content-Type: application/json` on every response. Eliminates repetition across individual handlers.

```
chainMiddleware
```
Applies middleware in declaration order using a reverse-iteration pattern. This ensures the first middleware listed is the outermost wrapper, which is the intuitive behavior.

```
http.Server with timeouts
```
`ReadTimeout: 10s` — maximum time to read the full request including body.
`WriteTimeout: 15s` — maximum time to write the full response.
`IdleTimeout: 60s` — maximum time to keep an idle keep-alive connection open.
These prevent slow-client attacks (Slowloris) from exhausting server resources.

```
Graceful shutdown
```
The server runs in a goroutine. The main goroutine blocks on a channel waiting for `SIGINT` or `SIGTERM`. On signal receipt, `server.Shutdown()` is called with a 10-second deadline, allowing in-flight requests to complete before the process exits.

---

### internal/model/transaction.go

Pure data structures with no methods or logic. Used by all layers. Amount is stored as an integer in paisa to avoid floating-point arithmetic errors.

---

### internal/repository/transaction.go

All PostgreSQL interactions. Three design principles apply here:

Every method accepts a `context.Context` as its first argument, enabling request-scoped cancellation and timeout propagation.

`BeginTx` exposes database transaction management to the service layer, which is the correct owner of transaction boundaries.

`UpdateBalanceTx` performs the balance deduction and the sufficiency check in a single atomic SQL statement, eliminating the read-check-write race condition.

---

### internal/service/transaction.go

Owns all business logic. `ProcessSync` iterates over the batch and delegates each transaction to `processSingle`.

`processSingle` is a private method extracted from the loop specifically to enable safe use of `defer`. Each call to `processSingle` has its own stack frame, so `defer dbTx.Rollback()` executes immediately when that call returns — not at the end of the entire batch loop. This ensures database connections are released promptly.

The `defer dbTx.Rollback()` after a successful `Commit()` is safe because PostgreSQL ignores a rollback on an already-committed transaction.

---

### internal/handler/transaction.go

HTTP boundary layer. Each handler follows the same pattern: validate HTTP method, decode and validate request body, call service, encode and return response. Handlers have no knowledge of database operations or business rules.

---

## 13. Testing

### Manual API Tests

**Health check:**
```bash
curl http://localhost:8080/health
```
Expected: `{"service":"tapx-lite","status":"ok"}`

**Balance query:**
```bash
curl "http://localhost:8080/balance?user_id=1"
```
Expected: `{"balance":10000,"user_id":1}`

**Valid transaction:**
```bash
curl -X POST http://localhost:8080/sync \
  -H "Content-Type: application/json" \
  -d '{"transactions":[{"id":"tx001","user_id":1,"amount":500,"nonce":"nonce-001"}]}'
```
Expected: `{"accepted":["tx001"],"rejected":null}`

**Duplicate nonce rejection:**
```bash
curl -X POST http://localhost:8080/sync \
  -H "Content-Type: application/json" \
  -d '{"transactions":[{"id":"tx002","user_id":1,"amount":500,"nonce":"nonce-001"}]}'
```
Expected: `{"accepted":null,"rejected":["tx002"]}`

**Insufficient balance:**
```bash
curl -X POST http://localhost:8080/sync \
  -H "Content-Type: application/json" \
  -d '{"transactions":[{"id":"tx003","user_id":2,"amount":999999,"nonce":"nonce-999"}]}'
```
Expected: `{"accepted":null,"rejected":["tx003"]}`

---

## 14. Known Limitations and Future Work

### Current Limitations

**No authentication.** The API is open. Any client can submit transactions for any user. JWT-based authentication is required before any public deployment.

**Redis is unused.** The Redis container is running but not integrated. The intended use is as a fast nonce cache to reduce database load on high-volume sync operations.

**No rate limiting.** A malicious client can flood the `/sync` endpoint. Rate limiting middleware is required.

**sslmode=disable.** Acceptable for local development. Production deployment must set `DB_SSLMODE=require` and provision SSL certificates.

**No unit tests.** Service and repository logic requires test coverage before this can be considered production-ready.

### Planned Enhancements

- JWT authentication middleware
- Redis nonce cache integration
- Rate limiting per client IP
- Unit tests for service layer with mock repository
- GitHub Actions CI pipeline (automated test and build on push)
- AWS EC2 deployment with RDS and ElastiCache
- Structured JSON logging (replacing `log.Printf`)
- NFC transaction token signing and verification

---

## Author

Shihab — [@Shihab369](https://github.com/Shihab369)

Built as part of the TapX fintech infrastructure project targeting the Bangladesh institutional payment market.
