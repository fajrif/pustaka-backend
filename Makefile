.PHONY: help install run build test clean dev
.PHONY: db-create db-drop db-migrate db-setup db-reset db-rollback db-status db-version db-console db-info
.PHONY: migrate migrate-status migrate-rollback migrate-create build-migrate

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Default target
help:
	@echo "=========================================="
	@echo "  Merk Buku API - Makefile Commands"
	@echo "=========================================="
	@echo ""
	@echo "📦 Application Commands:"
	@echo "  make install          - Install Go dependencies"
	@echo "  make run              - Run the application"
	@echo "  make build            - Build the application binary"
	@echo "  make dev              - Run with hot reload (requires air)"
	@echo "  make test             - Run tests"
	@echo "  make clean            - Clean build artifacts"
	@echo ""
	@echo "🗄️  Database Commands (Rails-style with dash):"
	@echo "  make db-create        - Create database"
	@echo "  make db-drop          - Drop database (⚠ dangerous!)"
	@echo "  make db-migrate       - Run pending migrations"
	@echo "  make db-setup         - Create + migrate (first time)"
	@echo "  make db-reset         - Drop + create + migrate (⚠ dangerous!)"
	@echo "  make db-rollback      - Rollback last migration"
	@echo "  make db-status        - Show migration status"
	@echo "  make db-version       - Show current migration version"
	@echo "  make db-console       - Open PostgreSQL console"
	@echo "  make db-info          - Show database information"
	@echo ""
	@echo "🔄 Alternative Commands:"
	@echo "  make migrate          - Same as db-migrate"
	@echo "  make migrate-status   - Same as db-status"
	@echo "  make migrate-rollback - Same as db-rollback"
	@echo "  make migrate-create   - Create new migration file"
	@echo ""
	@echo "💡 Examples:"
	@echo "  make db-setup                      # First time setup"
	@echo "  make db-migrate                    # Run migrations"
	@echo "  make db-drop db-create db-migrate  # Full reset"
	@echo "  make db-reset                      # Full reset (one command)"
	@echo ""
	@echo "Note: Commands use dash (-) instead of colon (:) for Make compatibility"
	@echo "      Shell scripts (./db.sh) can still use original names"
	@echo ""

# ============================================================================
# Application Commands
# ============================================================================

install:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies installed"

run:
	@echo "🚀 Starting application..."
	go run main.go

build:
	@echo "🔨 Building application..."
	mkdir -p bin
	go build -o bin/app main.go
	@echo "✓ Binary built: bin/app"

test:
	@echo "🧪 Running tests..."
	go test -v ./...

clean:
	@echo "🧹 Cleaning build artifacts..."
	go clean
	rm -rf bin/
	@echo "✓ Clean completed"

dev:
	@echo "🔥 Starting development server with hot reload..."
	@if ! command -v air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

# ============================================================================
# Database Commands (using dash instead of colon for Make compatibility)
# ============================================================================

db-create:
	@chmod +x db.sh
	@./db.sh create

db-drop:
	@chmod +x db.sh
	@./db.sh drop

db-migrate:
	@chmod +x db.sh
	@./db.sh migrate

db-setup:
	@chmod +x db.sh
	@./db.sh setup

db-reset:
	@chmod +x db.sh
	@./db.sh reset

db-rollback:
	@chmod +x db.sh
	@./db.sh rollback

db-status:
	@chmod +x db.sh
	@./db.sh status

db-version:
	@chmod +x db.sh
	@./db.sh version

db-console:
	@chmod +x db.sh
	@./db.sh console

db-info:
	@chmod +x db.sh
	@./db.sh info

# ============================================================================
# Alternative Migration Commands (without db- prefix)
# ============================================================================

migrate: db-migrate

migrate-status: db-status

migrate-rollback: db-rollback

migrate-create:
	@chmod +x migrate.sh
	@./migrate.sh create

# ============================================================================
# Build Tools
# ============================================================================

build-migrate:
	@echo "🔨 Building migration tool..."
	mkdir -p bin
	go build -o bin/migrate migrate.go
	@echo "✓ Migration tool built: bin/migrate"
