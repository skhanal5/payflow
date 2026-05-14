# Payflow — Agent Guide

Go monorepo: 2 gRPC services (order, product) + HTTP API gateway, Kafka async messaging, PostgreSQL, React frontend.

## Entrypoints

| Service | File | gRPC Port |
|---|---|---|
| API gateway (HTTP) | `cmd/payflow/apigateway/main.go` | N/A (HTTP on `APIGATEWAY_PORT`) |
| Order | `cmd/payflow/order/main.go` | `ORDER_GRPC_PORT` |
| Product | `cmd/payflow/product/main.go` | `PRODUCT_GRPC_PORT` |

Order emits `order.requested` → Kafka, consumes `inventory.checked` from Kafka.
Product consumes `order.requested` → Kafka, emits `inventory.checked` to Kafka.

## Infrastructure (Docker Compose)

- **broker**: Kafka 3.9.1 single-node (KRaft, port 9092)
- **order-db**: Postgres 15 on port 5432, init from `sql/init_order.sql`
- **inventory-db**: Postgres 15 on port 5433, init from `sql/init_inventory.sql`
- **apigateway**: HTTP on host port 8080, proxies REST→gRPC to order + product
- **frontend**: nginx serving SPA on host port 3000
- **pgAdmin** on port 8888, **Dozzle** (logs) on port 8080 (note: port conflict with apigateway)

Schema applied via Postgres `docker-entrypoint-initdb.d` scripts — no Flyway.

## Commands

| Command | What it does |
|---|---|
| `make fmt && make vet && make lint` | Pre-push gate (`golangci-lint` required locally) |
| `make build` | `go build ./...` |
| `make rund` | `docker compose down -v` + prune + `up --build -d` (destroys volumes) |
| `make restart` | `docker compose down` + `up -d` (preserves volumes) |
| `make inventory-service` | Builds **product** service binary to `bin/` with `.env` vars |
| `go run ./cmd/payflow/product/` | Run a single service natively (substitute order/apigateway) |
| `make frontend-dev` | `cd frontend && npm run dev` (Vite dev server) |
| `make frontend-build` | `cd frontend && npm run build` |
| `make frontend-install` | `cd frontend && npm install` |

Requires `.env` file (see `README.md` for template). Missing env vars panic at startup via `internal/shared/env/env.go`.

## Proto codegen

Generated stubs checked in under `gen/go/`. Proto sources in `proto/`. Codegen tools declared in `go.mod` `tool` directive:
- `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-grpc-gateway`, `protoc-gen-openapiv2`

No make target — run manually if protos change.

## Frontend

React 19 + Vite + Tailwind CSS 4 + shadcn/ui + Radix UI. TypeScript (.tsx/.ts). ESLint (`typescript-eslint`) + Prettier for linting/formatting.

## Git hooks

```sh
git config core.hooksPath githooks
```

Enables pre-push hook: `go fmt` → `go vet` → `golangci-lint`.

## Required env vars

**Shared**: `KAFKA_BROKER`, `KAFKA_GROUPID`, `ENVIRONMENT`

**Order**: `ORDER_REQUESTED_TOPIC`, `INVENTORY_CHECKED_TOPIC`, `ORDER_DB_HOST`, `ORDER_DB_USER`, `ORDER_DB_PASSWORD`, `ORDER_DB_PORT`, `ORDER_GRPC_PORT`

**Product**: `ORDER_REQUESTED_TOPIC`, `INVENTORY_CHECKED_TOPIC`, `INVENTORY_DB_HOST`, `INVENTORY_DB_USERNAME`, `INVENTORY_DB_PASSWORD`, `INVENTORY_DB_PORT`, `PRODUCT_GRPC_PORT`

**API gateway**: `APIGATEWAY_PORT`, `ORDER_SERVICE`, `PRODUCT_SERVICE`, `JWT_SECRET_KEY`
