.PHONY: help build run dev test clean docker-up docker-down migrate swagger lint

# Default target
help:
	@echo "PaperBeginner - AI-powered Academic Research Guide"
	@echo ""
	@echo "Usage:"
	@echo "  make dev          - Start development environment"
	@echo "  make build        - Build all services"
	@echo "  make run          - Run the API server locally"
	@echo "  make test         - Run tests"
	@echo "  make lint         - Run linters"
	@echo "  make swagger      - Generate Swagger documentation"
	@echo "  make migrate      - Run database migrations"
	@echo "  make docker-up    - Start all Docker services"
	@echo "  make docker-down  - Stop all Docker services"
	@echo "  make clean        - Clean build artifacts"

# Development
dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up --build

dev-infra:
	docker compose up postgres redis minio -d

# Build
build:
	cd backend && go build -o ../bin/api ./cmd/api
	cd backend && go build -o ../bin/worker ./cmd/worker

build-docker:
	docker compose build

# Run locally (requires dev-infra)
run:
	cd backend && go mod tidy && go run ./cmd/api

run-worker:
	cd backend && go run ./cmd/worker

# Testing
test:
	cd backend && go test -v -race ./...

test-coverage:
	cd backend && go test -v -race -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html

# Linting
lint:
	cd backend && golangci-lint run ./...

# Swagger documentation
swagger:
	cd backend && swag init -g cmd/api/main.go -o api/docs

# Database migrations
migrate:
	cd backend && go run ./cmd/migrate up

migrate-down:
	cd backend && go run ./cmd/migrate down

migrate-create:
	@read -p "Migration name: " name; \
	cd backend && go run ./cmd/migrate create $$name

# Docker commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-clean:
	docker compose down -v --rmi local

# Frontend
frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

# Clean
clean:
	rm -rf bin/
	rm -rf backend/api/docs/
	cd frontend && rm -rf dist/ node_modules/

# Initialize project (first time setup)
init: dev-infra
	@echo "Waiting for services to be ready..."
	@sleep 5
	$(MAKE) migrate
	@echo "Project initialized successfully!"

