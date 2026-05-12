# Payflow — Agent Guide

Go monorepo for a transaction processing system: 2 gRPC services + HTTP API gateway, Kafka async messaging, PostgreSQL.

## Entrypoints

| Service | File | Status |
|---|---|---|
| API gateway (HTTP) | `cmd/payflow/apigateway/main.go` | Active — grpc-gateway proxying REST→gRPC, JWT auth middleware, `/health` endpoint |
| Product (gRPC) | `cmd/payflow/product/main.go` | Active — product catalog gRPC server, authz interceptor |
| Order (gRPC) | `cmd/payflow/order/main.go` | Deleted — entrypoint was commented out; `internal/order/{config,repository}` still compiled as library |

## Infrastructure (Docker Compose)

- **Kafka**: Apache Kafka 3.9.1 single-node (KRaft, port 9092)
- **order-db**: Postgres 15 on port 5432, init from `sql/init_order.sql`
- **inventory-db**: Postgres 15 on port 5433, init from `sql/init_inventory.sql`
- **pgAdmin** on port 8888, **Dozzle** (logs) on port 8080

No Flyway — schema is applied via `docker-entrypoint-initdb.d` scripts baked into the Postgres images.

## Key refactors (relative to `main`)

| `main` branch | This branch |
|---|---|
| `internal/utility/` | `internal/shared/{auth,db,env,logger}/` |
| `internal/inventory/` | `internal/product/` (richer model + gRPC server + authz interceptor) |
| `db/` (Flyway) | `sql/` (native init scripts) |
| Per-service proto files | `proto/` → `gen/go/` (shared protos, grpc-gateway stubs checked in) |

## Commands

| Command | What it does |
|---|---|
| `make fmt` | `go fmt ./...` |
| `make vet` | `go vet ./...` |
| `make lint` | `golangci-lint run ./...` (requires `golangci-lint` installed) |
| `make build` | `go build ./...` |
| `make rund` | `docker compose down -v` + prune + `up --build -d` (destroys volumes) |
| `make restart` | `docker compose down` + `up -d` (preserves volumes) |
| `make restart-db` | `docker restart order-db inventory-db` |
| `make inventory-service` | Builds product service binary with env vars from `.env` |

No test targets, no CI, no proto codegen target (generated files are checked in).

After cloning, run `git config core.hooksPath githooks` to enable the pre-push hook (runs `fmt` → `vet` → `lint` before each push).

## Required env vars (all panic if missing)

**Shared**: `KAFKA_BROKER`, `KAFKA_GROUPID`, `ENVIRONMENT`

**Order service**: `ORDER_REQUESTED_TOPIC`, `INVENTORY_CHECKED_TOPIC`, `ORDER_DB_HOST`, `ORDER_DB_USER`, `ORDER_DB_PASSWORD`, `ORDER_DB_PORT`

**Product service**: `ORDER_REQUESTED_TOPIC`, `INVENTORY_CHECKED_TOPIC`, `INVENTORY_DB_HOST`, `INVENTORY_DB_USERNAME`, `INVENTORY_DB_PASSWORD`, `INVENTORY_DB_PORT`, `PRODUCT_GRPC_PORT`

**API gateway**: `APIGATEWAY_PORT`, `ORDER_SERVICE`, `PRODUCT_SERVICE`, `JWT_SECRET_KEY`

## Known bugs

None currently known.
