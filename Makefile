.PHONY: clean rund restart migrate ps help

clean:
	docker compose down -v
	docker network prune -f
	docker volume prune -f

rund: clean 
	docker compose up --build -d

restart:
	docker compose down
	docker compose up -d

migrate:
	docker compose stop flyway
	docker compose rm -f flyway
	docker compose up flyway

ps:
	docker-compose ps

help:
	@echo "Available targets:"
	@echo "  ps       - List the status of the containers"
	@echo "  rund     - Remove and run the containers"
	@echo "  restart  - Restart the containers"
	@echo "  migrate  - Run flyway migrations"
	@echo "  prune    - Remove containers, networks, and volumes"