# marketflow
Real-Time Cryptocurrency Market Data Processing System using Go, Redis, and PostgreSQL with support for concurrency patterns and hexagonal architecture.

---

## 🧭 Table of Contents

* [🚀 Launch the Project](#launch-the-project)
* [📁 Project Structure](#project-structure)
* [🌐 API Endpoints](#api-endpoints)

---

## 🚀 Launch the Project

Make sure you have Docker installed, then run:

```bash
docker compose up --build
```

🔧 Don’t forget to configure `.env` properly before launch.

---

## 📁 Project Structure

```bash
marketflow/
.
├── Dockerfile
├── README.md
├── Task.md
├── cmd
│   └── app.go
├── configs
│   └── config.json
├── docker-compose.yml
├── exchange_tars
│   ├── exchange1_amd64.tar
│   ├── exchange2_amd64.tar
│   └── exchange3_amd64.tar
├── go.mod
├── go.sum
├── internal
│   ├── adapters
│   │   ├── cache
│   │   │   └── redis.go
│   │   ├── exchange
│   │   │   ├── generator.go
│   │   │   ├── tcp_client.go
│   │   │   └── wire.go
│   │   ├── postgres
│   │   │   ├── connect.go
│   │   │   └── repository.go
│   │   └── web
│   │       ├── handlers_market_data.go
│   │       ├── handlers_system.go
│   │       ├── helpers.go
│   │       ├── router.go
│   │       └── server.go
│   ├── app
│   │   ├── market_data_service.go
│   │   ├── pipeline.go
│   │   ├── system_service.go
│   │   └── window_store.go
│   ├── config
│   │   ├── config.go
│   │   └── flags.go
│   ├── domain
│   │   ├── errors.go
│   │   ├── model.go
│   │   └── ports.go
│   └── infra
│       └── log.go
├── main.go
├── makefile
└── migrations
    └── init.sql

15 directories, 35 files
```

🧱 **Architecture:** Clean Hexagonal (Ports & Adapters)<br>
⚙️ **Patterns Used:** Fan-in, Fan-out, Worker Pool, Generator

---

## 🌐 API Endpoints

### 📊 Price Endpoints

| Method | Endpoint                                                |
| ------ | ------------------------------------------------------- |
| GET    | `/prices/latest/{symbol}`                               |
| GET    | `/prices/latest/{exchange}/{symbol}`                    |
| GET    | `/prices/highest/{symbol}`                              |
| GET    | `/prices/highest/{exchange}/{symbol}`                   |
| GET    | `/prices/highest/{symbol}?period={duration}`            |
| GET    | `/prices/highest/{exchange}/{symbol}?period={duration}` |
| GET    | `/prices/lowest/{symbol}`                               |
| GET    | `/prices/lowest/{exchange}/{symbol}`                    |
| GET    | `/prices/lowest/{symbol}?period={duration}`             |
| GET    | `/prices/lowest/{exchange}/{symbol}?period={duration}`  |
| GET    | `/prices/average/{symbol}`                              |
| GET    | `/prices/average/{exchange}/{symbol}`                   |
| GET    | `/prices/average/{symbol}?period={duration}`            |
| GET    | `/prices/average/{exchange}/{symbol}?period={duration}` |

### 🔁 Mode Switching

| Method | Endpoint                               |
| ------ | -------------------------------------- |
| POST   | `/mode/live` – Switch to **Live Mode** |
| POST   | `/mode/test` – Switch to **Test Mode** |

### 🩺 Health Check

| Method | Endpoint                          |
| ------ | --------------------------------- |
| GET    | `/health` – Returns system status |

---
