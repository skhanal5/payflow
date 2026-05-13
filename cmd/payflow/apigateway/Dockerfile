FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /service ./cmd/payflow/apigateway/

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /service /service
ENTRYPOINT ["/service"]
