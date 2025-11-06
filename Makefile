.PHONY: help run build test clean install-deps migrate frontend-dev frontend-build

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-deps: ## Install Go dependencies
	go mod download
	go mod tidy

run: ## Run the backend server (uses config.local.yaml)
	cp config.local.yaml config.yaml
	go run cmd/server/main.go

run-docker: ## Run with Docker config
	cp config.docker.yaml config.yaml
	go run cmd/server/main.go

build: ## Build the backend binary
	go build -o bin/rotki-demo cmd/server/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf frontend/dist/

migrate: ## Run database migrations (requires running server once)
	@echo "Migrations run automatically on server start"

migrate-rpc: ## Run RPC nodes table migration
	@echo "Running RPC nodes migration..."
	mysql -u root -p rotki_demo < migrations/002_add_rpc_nodes_table.sql
	@echo "Migration complete!"

test-rpc: ## Test RPC nodes API endpoints
	@echo "Testing RPC nodes API..."
	@bash scripts/test_rpc_nodes.sh

frontend-install: ## Install frontend dependencies
	cd frontend && npm install

frontend-dev: ## Run frontend development server
	cd frontend && npm run dev

frontend-build: ## Build frontend for production
	cd frontend && npm run build

docker-build: ## Build Docker image
	docker-compose build

docker-up: ## Start all Docker services
	docker-compose up -d

docker-down: ## Stop all Docker services
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

docker-db: ## Start only database services (MySQL + Redis)
	docker-compose up -d mysql redis

docker-restart: ## Restart all Docker services
	docker-compose restart

docker-clean: ## Stop and remove all containers, networks, and volumes
	docker-compose down -v

docker-ps: ## Show running containers
	docker-compose ps

dev: ## Run both backend and frontend in development mode
	@echo "Starting backend and frontend..."
	@trap 'kill 0' EXIT; \
	go run cmd/server/main.go & \
	cd frontend && npm run dev

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

fmt-frontend: ## Format frontend code
	cd frontend && npm run format

all: clean install-deps build frontend-build ## Clean, install dependencies, and build everything
