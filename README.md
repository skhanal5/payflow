## Payflow

### About
Learning how to build a transaction processing system.

### Local Development

#### Environment Variables
First set up the following environment variables in a `.env` file

Example: 
```
DB_USER=root
DB_PASSWORD=changeme
PGADMIN_EMAIL=root@localhost.com
PGADMIN_PASSWORD=changeme
```

#### Docker
Use `make rund` to build the containers from scratch or `make restart` to restart the containers.

This will do the following:
- Dozzle will run on http://localhost:8080
- pgAdmin will run on http://localhost:8888
- PostgreSQL will be running on host: db and port: 5432
