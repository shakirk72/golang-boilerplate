# go-mvc-app

A Go REST API built with Gin, GORM (PostgreSQL), Zap logging, JWT auth, and Swagger docs — following an MVC structure.

---

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.25+ | https://go.dev/dl |
| PostgreSQL | 16+ | https://www.postgresql.org/download or via Docker |
| `swag` CLI | latest | `make install-swag` |
| `air` CLI | latest | `make install-air` *(only for live-reload dev)* |

---

## 1 — Clone & install dependencies

```bash
git clone <repo-url>
cd <project-name>

# Download Go modules
make tidy
```

---

## 2 — Configure environment

```bash
cp .env.example .env
```

Open `.env` and fill in the required values:

```env
APP_ENV=dev          # dev | qa | prod
APP_PORT=8080

# PostgreSQL — required
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=gomvc_dev

# External API — required
EXTERNAL_API_BASE_URL=https://api.example.com
EXTERNAL_API_KEY=your-api-key

# JWT — required
JWT_SECRET=your-jwt-secret

# CORS — comma-separated origins, or * for dev
CORS_ORIGINS=*

# Log sinks: files | console | files,console
LOGS_TARGETS=files,console
```

> Optional overrides (`DB_PORT`, `DB_SSLMODE`, `DB_TIMEZONE`, `EXTERNAL_API_TIMEOUT`, `JWT_EXPIRY_HOURS`) are commented out in `.env.example` — uncomment to use.

---

## 3 — Start PostgreSQL

**Option A — Docker (recommended for local dev)**

```bash
docker compose up db -d
```

This starts Postgres 16 on port `5432` with DB `gomvc_dev` / user `postgres` / password `postgres`.
Update `.env` to match if you changed those values.

**Option B — Local Postgres**

Create the database manually:

```bash
psql -U postgres -c "CREATE DATABASE gomvc_dev;"
```

---

## 4 — Run the app

### Development (one-shot)

Regenerates Swagger docs then starts the server:

```bash
make dev
```

### Development with live reload

Watches `.go`, `.toml`, and `.env` files — rebuilds and restarts on every save:

```bash
make watch
```

### Other environments

```bash
make qa    # APP_ENV=qa
make prod  # APP_ENV=prod
```

### Build a binary

```bash
make build
./dist/server
```

---

## 5 — Verify it's running

```bash
curl http://localhost:8080/health
```

Swagger UI is available at:

```
http://localhost:8080/swagger/index.html
```

---

## 6 — Run with Docker Compose (app + db together)

```bash
docker compose up --build
```

The `app` service reads `.env` automatically and waits for Postgres to be healthy before starting.

---

## Project structure

```
src/
  cmd/            # Entry point (main.go)
  config/         # Env loading, common constants
  controllers/    # HTTP handlers
  services/       # Business logic
  repositories/   # DB layer (GORM)
  models/         # GORM models
  routes/         # Route registration
  middleware/     # Auth, logging, CORS
  helpers/        # Shared utilities
  utils/          # Low-level helpers
docs/             # Auto-generated Swagger files
logs/             # Log output (when LOGS_TARGETS=files)
```

---

## Makefile targets

| Target | Description |
|--------|-------------|
| `make tidy` | Download / tidy Go modules |
| `make dev` | Regen Swagger + run (dev) |
| `make watch` | Live-reload dev server |
| `make qa` | Run with `APP_ENV=qa` |
| `make prod` | Run with `APP_ENV=prod` |
| `make build` | Compile binary to `dist/server` |
| `make swagger` | Regenerate Swagger docs only |
| `make install-swag` | Install `swag` CLI (one-time) |
| `make install-air` | Install `air` CLI (one-time) |
