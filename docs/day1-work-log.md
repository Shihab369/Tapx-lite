# TapX Lite - Day 1 Work Log
### Date: April 05, 2026

---

## Summary

Aaj prothom din. Environment setup theke shuru kore actual Go code lekha porjonto pouchhechi. Nicher sob kaj completed.

---

## What I Completed Today

### 1. Environment Setup and Fixes

| Task | Status |
|------|--------|
| Docker official repository theke install | Done |
| Docker service start + enable | Done |
| User docker group e add | Done |
| Git identity set (name + email) | Done |
| SSH key generate (ed25519) | Done |
| GitHub SSH connection verify | Done |
| Zsh Powerlevel10k typo fix | Done |

### 2. GitHub Repository

- Repository name: `Tapx-lite`
- URL: `https://github.com/Shihab369/Tapx-lite`
- Visibility: Public
- License: MIT
- .gitignore: Go

### 3. Project Structure Created

```
Tapx-lite/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   │   └── transaction.go
│   ├── service/
│   │   └── transaction.go
│   ├── repository/
│   │   └── transaction.go
│   └── model/
│       └── transaction.go
├── docs/
│   └── env_report_20260405_003757.md
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── .env
├── .gitignore
├── LICENSE
└── README.md
```

### 4. Go Files Written

#### internal/model/transaction.go
Data structures define kora hoyeche:
- `Transaction` struct - ekta transaction er shape
- `SyncRequest` struct - batch transaction er request format
- `SyncResponse` struct - accepted r rejected er response format

#### internal/repository/transaction.go
Database layer lekha hoyeche:
- `NonceExists` - duplicate transaction check
- `InsertTransaction` - notun transaction save
- `GetBalance` - user er balance fetch
- `UpdateBalance` - transaction er pore balance update

#### internal/service/transaction.go
Business logic lekha hoyeche:
- `ProcessSync` - core logic, 4 step e transaction process kore
  - Step 1: Nonce check (duplicate detect)
  - Step 2: Balance check (sufficient funds)
  - Step 3: Transaction save
  - Step 4: Balance update
- `GetBalance` - balance fetch with error handling

#### internal/handler/transaction.go
HTTP layer lekha hoyeche:
- `HealthCheck` - GET /health, server alive kina check
- `Sync` - POST /sync, batch transaction receive
- `GetBalance` - GET /balance?user_id=1, balance query

---

## Architecture Flow Learned Today

```
HTTP Request
     ↓
Handler    (request receive, validate)
     ↓
Service    (business logic, decisions)
     ↓
Repository (database queries only)
     ↓
Model      (data shape, used by all layers)
```

---

## Key Concepts Learned Today

**Layered Architecture:** Protita layer er ekta kaj. Handler database janena, repository business logic janena. Ei separation er karone code maintainable thake.

**Struct Tags:** `json:"id"` diye Go ke bola hoy JSON e field ta ki name e thakbe.

**Constructor Pattern:** Go te class nei, `NewXxx()` function diye object banano hoy.

**Context:** Protita function e `context.Context` pass kora hoy. Timeout r cancellation handle kore.

**Nonce:** One-time unique value. Duplicate transaction detect korte use hoy. Same nonce duibar ashle second ta reject.

**Layered Error Handling:** Service layer e protita step fail hole `continue` diye porer transaction e jay. Kono ekta transaction er failure puro batch k block kore na.

---

## What Is Remaining (Next Session)

- [ ] `main.go` lekha - sob layer ek sathe jora lagabe
- [ ] `docker-compose.yml` lekha - PostgreSQL + Redis local e chalabe
- [ ] `.env` file complete kora
- [ ] Database migration file lekha
- [ ] `Dockerfile` lekha
- [ ] Local e run kore test kora
- [ ] GitHub Actions CI/CD pipeline setup
- [ ] AWS deployment

---

## Commands Used Today

```bash
# Docker install
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker shihab
newgrp docker

# Git setup
git config --global user.name "Shihab"
git config --global user.email "shihabud121@outlook.com"

# SSH key
ssh-keygen -t ed25519 -C "shihabud121@outlook.com"
ssh -T git@github.com

# Project setup
mkdir -p ~/projects
cd ~/projects
git clone git@github.com:Shihab369/Tapx-lite.git
cd Tapx-lite
go mod init github.com/Shihab369/Tapx-lite

# Git workflow
git add .
git commit -m "chore: initial project structure setup"
git push
```

---

## Mentor Notes

Aaj onek kaj hoyeche. Ekjon beginner er jonno ei amount of work first day e excellent.

Porer session e main.go lekhar sathe puro system locally chalabe. Seটাই prothom real milestone - ekta working API docker e chole, database connected, curl diye test hobe.

Rest nao. Fresh mind e code valo hoy.

---

> Day 1 complete. Environment ready. Architecture understood. Code started.
