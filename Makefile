.PHONY: fmt clean rund restart ps help

fmt: 
	go fmt ./...
	
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

ps:
	docker-compose ps

help:
	@echo "Available targets:"
	@echo "  ps       - List the status of the containers"
	@echo "  rund     - Remove and run the containers"
	@echo "  restart  - Restart all containers"
	@echo "  prune    - Remove containers, networks, and volumes"