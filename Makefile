################################################################################
# Run the Service and DB containers

# This will:
# 1. build the container,
# 2. run the db container and create tables,
# 3. run the service container
bootstrap-all: build-svc clean-db bootstrap-db clean-svc run-dev-svc


################################################################################
# Running Tests

# Requires the service containers to be running locally
test-api: 
	hurl --very-verbose ./tests/api_test.hurl


test-unit: 
	go test -v ./...


################################################################################
# Database

DB_IMAGE_NAME := txn_db
NETWORK_NAME := my_network

PG_USER = transaction_svc
PG_PASS = dev_pass # in production - we'd use secrets management for this

# remove any containers for the postgres db
clean-db:
	docker stop ${DB_IMAGE_NAME} || true
	docker rm ${DB_IMAGE_NAME} || true

# run the postgres container
launch-db: clean-db
	docker network create ${NETWORK_NAME} || true
	docker run --name ${DB_IMAGE_NAME} --network ${NETWORK_NAME} -e POSTGRES_USER=${PG_USER} -e POSTGRES_PASSWORD=${PG_PASS} -d -p 5432:5432 postgres

# runs the postgres container, and then creates the relevant tables for the service
bootstrap-db: launch-db
	sleep 1
	PGPASSWORD=${PG_PASS} psql -h localhost -U ${PG_USER} -f ./db/schema_postgres.sql

# gives you a shell into the db for running manual SQL queries
db-shell:
	PGPASSWORD=${PG_PASS} psql -h localhost -d transactions_db -p 5432 -U ${PG_USER}

################################################################################
# Service

SVC_IMAGE_NAME := transaction_service

# removes any containers for the service
clean-svc:
	docker stop $(SVC_IMAGE_NAME)_container || true
	docker rm $(SVC_IMAGE_NAME)_container || true

# build the Docker image
build-svc:
	docker build -t $(SVC_IMAGE_NAME) .

# run the docker container for the service
run-dev-svc: 
	docker run --network ${NETWORK_NAME} -p 8080:8080 --name $(SVC_IMAGE_NAME)_container $(SVC_IMAGE_NAME)


