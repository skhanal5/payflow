## Payflow

Learning how to build a transaction processing system — Go monorepo with 2 gRPC services, an HTTP API gateway, Kafka async messaging, and PostgreSQL.

### Services

| Service | Entrypoint | Port |
|---|---|---|
| API gateway (HTTP) | `cmd/payflow/apigateway/main.go` | Configurable via `APIGATEWAY_PORT` |
| Product (gRPC + Kafka consumer) | `cmd/payflow/product/main.go` | Configurable via `PRODUCT_GRPC_PORT` |
| Order (gRPC + Kafka producer/consumer) | `cmd/payflow/order/main.go` | Configurable via `ORDER_GRPC_PORT` |

### Infrastructure (Docker Compose)

- **order-db**: Postgres 15 on port 5432
- **inventory-db**: Postgres 15 on port 5433
- **Kafka**: Apache Kafka 3.9.1 (KRaft, port 9092)
- **pgAdmin**: Port 8888
- **Dozzle**: Port 8080
- **product**: gRPC server on port 50052, Kafka consumer
- **order**: gRPC server on port 50051, Kafka producer/consumer
- **apigateway**: HTTP server on port 8080, grpc-gateway proxying to product + order

### Environment Variables

Create a `.env` file with these:

```
# Docker
DB_USER=root
DB_PASSWORD=changeme
PGADMIN_EMAIL=root@localhost.com
PGADMIN_PASSWORD=changeme

# Shared
KAFKA_BROKER=localhost:9092
KAFKA_GROUPID=payflow-group
ENVIRONMENT=development

# Order service
ORDER_REQUESTED_TOPIC=order.requested
INVENTORY_CHECKED_TOPIC=inventory.checked
ORDER_DB_HOST=localhost
ORDER_DB_USER=root
ORDER_DB_PASSWORD=changeme
ORDER_DB_PORT=5432
ORDER_GRPC_PORT=:50051

# Product service
INVENTORY_DB_HOST=localhost
INVENTORY_DB_USERNAME=root
INVENTORY_DB_PASSWORD=changeme
INVENTORY_DB_PORT=5433
PRODUCT_GRPC_PORT=:50052

# API gateway
APIGATEWAY_PORT=:8080
ORDER_SERVICE=localhost:50051
PRODUCT_SERVICE=localhost:50052
JWT_SECRET_KEY=change-me
```

### Local Development

```sh
# Build and start everything (infrastructure + Go services)
make rund

# Or run a service natively for faster iteration
go run ./cmd/payflow/product/main.go
```

### Git Hooks

After cloning:

```sh
git config core.hooksPath githooks
```

This enables the pre-push hook that runs `go fmt` → `go vet` → `golangci-lint` before each push.
