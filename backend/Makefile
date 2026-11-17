.PHONY: run build test test-integration lint migrate-up migrate-down seed docker-up docker-down tidy

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
BINARY_NAME=api
WORKER_BINARY=worker
SEED_BINARY=seed
MIGRATE_BINARY=migrate
MAIN_PATH=./cmd/api
WORKER_PATH=./cmd/worker
SEED_PATH=./cmd/seed
MIGRATE_PATH=./cmd/migrate

# Database
MIGRATIONS_PATH=./migrations
DB_URL?=postgres://complylite:complylite@localhost:5432/complylite?sslmode=disable

# Docker
COMPOSE_FILE=../docker-compose.yml

# Build
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) $(MAIN_PATH)
	$(GOBUILD) -o bin/$(WORKER_BINARY) $(WORKER_PATH)
	$(GOBUILD) -o bin/$(SEED_BINARY) $(SEED_PATH)
	$(GOBUILD) -o bin/$(MIGRATE_BINARY) $(MIGRATE_PATH)

run: build
	./bin/$(BINARY_NAME)

worker: build
	./bin/$(WORKER_BINARY)

seed: build
	./bin/$(SEED_BINARY)

# Test
test:
	$(GOTEST) -race ./...

test-integration:
	TEST_DATABASE_URL=$(DB_URL) $(GOTEST) -race -tags=integration ./...

# Lint
lint:
	golangci-lint run ./...

# Database migrations
migrate-up:
	./bin/$(MIGRATE_BINARY) -path=$(MIGRATIONS_PATH) -database="$(DB_URL)" up

migrate-down:
	./bin/$(MIGRATE_BINARY) -path=$(MIGRATIONS_PATH) -database="$(DB_URL)" down

migrate-create:
	@read -p "Migration name: " name; \
	$(MIGRATE_BINARY) create -ext sql -dir $(MIGRATIONS_PATH) -seq $$name

# Dependencies
tidy:
	$(GOMOD) tidy

# Docker
docker-up:
	docker compose -f $(COMPOSE_FILE) up -d

docker-down:
	docker compose -f $(COMPOSE_FILE) down

docker-logs:
	docker compose -f $(COMPOSE_FILE) logs -f

# Generate mocks (if needed)
mocks:
	go generate ./...

# Clean
clean:
	$(GOCLEAN)
	rm -rf bin/