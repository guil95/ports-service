# Set variables
MIGRATE_IMAGE = migrate/migrate
MIGRATION_DIR = ./migrations
DATABASE_URL = postgres://root:root@localhost:5432/ports_db?sslmode=disable
IMAGE_NAME = ports-service
CONTAINER_NAME = ports-service-container
CONTAINER_SERVER_NAME = ports-service-container-server
NETWORK = ports_service
APP_DIR=$(shell pwd)
DOCKER_COMPOSE := $(shell command -v docker-compose 2> /dev/null || echo "docker compose")

# Default target
.PHONY: help
help:
	@echo "Makefile for running database migrations using docker"
	@echo "Usage:"
	@echo "  make run-import file=path/to/file.json - Run imports and save/update database"
	@echo "  make run-server - Run the server"
	@echo "  make migrate-up - Apply all up migrations"
	@echo "  make migrate-down - Apply all down migrations"
	@echo "  make migrate-up-one - Apply the next up migration"
	@echo "  make migrate-down-one - Apply the next down migration"
	@echo "  make create-migration - Create a migration"
	@echo "  make up-dependencies - To start the database"
	@echo "  make down-dependencies - To finish the database"
	@echo "  make integration-tests - To run integration tests"
	@echo "  make unit-tests - To run unit tests"
	@echo "  make tests - To run all tests"
	@echo "  make lint - Check lint"
	@echo "  make test-lint - Check lint and tests"

.PHONY: create-network
create-network:
	@docker network inspect $(NETWORK) >/dev/null 2>&1 || docker network create $(NETWORK)

# Build image
.PHONY: build
build:
	@docker build -t $(IMAGE_NAME) .

# Run server using docker
.PHONY: run-server
run-server: create-network
	$(DOCKER_COMPOSE) -f build/docker-compose-app.yml up

# Run server locally
.PHONY: run-server-locally
run-server-locally:
	@go run cmd/main.go

# Run import ports using docker
.PHONY: run-import
run-import: create-network
	@if [ -z "$(FILE)" ]; then \
		echo "invalid file"; \
		exit 1; \
	fi

	@docker run --rm -v $(APP_DIR):/app -w /app --network  $(NETWORK) -e DB_HOST=postgres golang:latest sh -c "go run cmd/main.go import -f /app/$(FILE)"

# Run imports locally
.PHONY: run-import-locally
run-import-locally:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: FILE variable is not set. Please provide a file path."; \
		exit 1; \
	fi
	@go run cmd/main.go import -f $(FILE)


# Apply all up migrations
.PHONY: migrate-up
migrate-up:
	@docker run -v $(PWD)/$(MIGRATION_DIR):/migrations --network host $(MIGRATE_IMAGE) \
		-path=/migrations/ -database $(DATABASE_URL) up

# Apply all down migrations
.PHONY: migrate-down
migrate-down:
	@docker run -v $(PWD)/$(MIGRATION_DIR):/migrations --network host $(MIGRATE_IMAGE) \
		-path=/migrations/ -database $(DATABASE_URL) down

# Apply the next up migration
.PHONY: migrate-up-one
migrate-up-one:
	@docker run -v $(PWD)/$(MIGRATION_DIR):/migrations --network host $(MIGRATE_IMAGE) \
		-path=/migrations/ -database $(DATABASE_URL) up 1

# Apply the next down migration
.PHONY: migrate-down-one
migrate-down-one:
	@docker run -v $(PWD)/$(MIGRATION_DIR):/migrations --network host $(MIGRATE_IMAGE) \
		-path=/migrations/ -database $(DATABASE_URL) down 1

# Create a new migration file
.PHONY: create-migration
create-migration:
	@if [ -z "$(name)" ]; then echo "Error: Please provide a migration name with 'name=<migration_name>'"; exit 1; fi
	@docker run --rm -v $(PWD)/$(MIGRATION_DIR):/migrations $(MIGRATE_IMAGE) create -ext sql -dir ./migrations -seq $(name)

# Start dependencies (database)
.PHONY: up-dependencies
up-dependencies: create-network
	@$(DOCKER_COMPOSE) -f build/docker-compose-db.yml up -d

# Stop dependencies (database)
.PHONY: down-dependencies
down-dependencies:
	@$(DOCKER_COMPOSE) -f build/docker-compose-db.yml down -v

## Generate mocks
.PHONY: mock-generate
mock-generate:
	@go mod tidy
	@docker run --rm -v "$(PWD):/app" -w /app -t vektra/mockery --all --dir ./internal --case underscore

.PHONY: integration-tests
integration-tests:
	@echo "Running integration tests..."
	@go test -count=1 -tags=integration -v ./... | grep -v "\[no test files\]"

.PHONY: unit-tests
unit-tests:
	@echo "Running unit tests..."
	@go test -count=1 -tags=unit -v ./... | grep -v "\[no test files\]"

.PHONY: tests
tests: integration-tests unit-tests

.PHONY: install-lint
install-lint:
	@test -f ./bin/golangci-lint || curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_LINT}

.PHONY: lint
lint: install-lint
	@echo "Running golangci-lint..."
	@bin/golangci-lint run

.PHONY: test-lint
test-lint: tests lint

