# backend-boiler

Baseline Golang API service structure using Fiber v3, GORM (PostgreSQL), Redis, RabbitMQ, slog, migration, and Swagger.

## Quick start

1. Copy `.env.example` to `.env`
2. Install dependencies:
   - `go mod tidy`
3. Run API:
   - `make run-api`
4. Run worker:
   - `make run-worker`

## Common commands

- `make test`
- `make build-api`
- `make build-worker`
- `make swagger`
- `make migrate-up`
