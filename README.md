# TapX Lite

A backend system that simulates offline payment reconciliation. Built to demonstrate practical skills in Go, PostgreSQL, Redis, and Docker.

---

## What It Does

Devices queue transactions offline. When connectivity is restored, the device submits a batch to this API. The backend validates each transaction and updates the ledger.

---

## Tech Stack

- **Go** REST API
- **PostgreSQL** primary database
- **Redis** included for future caching layer
- **Docker + Docker Compose** containerized local environment

---

## Project Structure

```
Tapx-lite/
├── cmd/server/main.go          # entry point, server lifecycle
├── internal/
│   ├── handler/                # HTTP layer
│   ├── service/                # business logic
│   ├── repository/             # database queries
│   └── model/                  # data structs
├── migrations/001_init.sql     # schema + seed data
├── Dockerfile                  # multi-stage build
└── docker-compose.yml          # app + postgres + redis
```

---

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | server liveness check |
| POST | `/sync` | submit batch of offline transactions |
| GET | `/balance?user_id={id}` | get account balance |

### POST /sync

```json
{
  "transactions": [
    {
      "id": "tx001",
      "user_id": 1,
      "amount": 500,
      "nonce": "unique-nonce-001"
    }
  ]
}
```

Response:
```json
{
  "accepted": ["tx001"],
  "rejected": null
}
```

---

## Key Concepts Demonstrated

- Layered architecture: handler → service → repository → model
- Atomic database transactions using `BEGIN` / `COMMIT` / `ROLLBACK`
- Nonce-based idempotency to prevent duplicate payments
- Database-level balance guard: `WHERE balance >= amount`
- Multi-stage Docker build
- Graceful shutdown on `SIGINT` / `SIGTERM`
- Environment-driven configuration with startup validation
- Connection pool tuning with `database/sql`

---

## Run Locally

```bash
git clone git@github.com:Shihab369/Tapx-lite.git
cd Tapx-lite

# start containers
docker compose up -d

# run migration
docker exec -i tapx-postgres psql -U postgres -d tapx < migrations/001_init.sql

# test
curl http://localhost:8080/health
```

---

## Environment Variables

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tapx
DB_SSLMODE=disable
SERVER_PORT=8080
```

---


## Documentation

- [Database Schema](docs/DATABASE.md)
- [Environment Setup](docs/dev-environment-setup.md)


## Author

Shihab — [github.com/Shihab369](https://github.com/Shihab369)