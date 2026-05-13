.PHONY: fmt vet lint build clean rund restart restart-db inventory-service frontend-install frontend-dev frontend-build ps help

fmt: 
	go fmt ./...
	
vet:
	go vet ./...
	
lint:
	PATH="$$(go env GOPATH)/bin:$$PATH" golangci-lint run ./...

build:
	go build ./...

clean:
	docker compose down -v
	docker network prune -f
	docker volume prune -f

rund: clean 
	docker compose up --build -d

restart:
	docker compose down
	docker compose up -d

restart-db:
	docker restart order-db
	docker restart inventory-db

inventory-service:
	mkdir -p bin
	export $$(cat .env | xargs) && go build -o bin/inventory-service ./cmd/payflow/product/main.go

ps:
	docker-compose ps

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

help:
	@echo "Available targets:"
	@echo "  ps       - List the status of the containers"
	@echo "  rund     - Remove and run the containers"
	@echo "  restart  - Restart all containers"
	@echo "  prune    - Remove containers, networks, and volumes"