PG_USER = transaction_svc
PG_PASS = dev_pass # in production - we'd use secrets management for this

################################################################################
# Database

DB_IMAGE_NAME := txn_db
NETWORK_NAME := my_network

clean-db:
	docker stop ${DB_IMAGE_NAME} || true
	docker rm ${DB_IMAGE_NAME} || true

launch-db: clean-db
	docker network create ${NETWORK_NAME} || true
	docker run --name ${DB_IMAGE_NAME} --network ${NETWORK_NAME} -e POSTGRES_USER=${PG_USER} -e POSTGRES_PASSWORD=${PG_PASS} -d -p 5432:5432 postgres

bootstrap-db: launch-db
	sleep 1
	PGPASSWORD=${PG_PASS} psql -h localhost -U ${PG_USER} -f ./db/schema_postgres.sql

db-shell:
	PGPASSWORD=${PG_PASS} psql -h localhost -d transactions_db -p 5432 -U ${PG_USER}

################################################################################
# Service

# Define the name of the Docker image
SVC_IMAGE_NAME := transaction_service

clean-svc:
	docker stop $(SVC_IMAGE_NAME)_container || true
	docker rm $(SVC_IMAGE_NAME)_container || true

# Target to build the Docker image
build-svc:
	docker build -t $(SVC_IMAGE_NAME) .

# Target to run the Docker container in development mode
run-dev-svc: 
	docker run --network ${NETWORK_NAME} -p 8080:8080 --name $(SVC_IMAGE_NAME)_container $(SVC_IMAGE_NAME)


################################################################################

bootstrap-all: build-svc clean-db bootstrap-db clean-svc run-dev-svc
