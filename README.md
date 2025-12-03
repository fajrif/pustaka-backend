## BACKEND API PUSTAKA DIGITAL
This is the backend API for the Digital Library application, built using
Go and PostgreSQL.

## 🚀 Installation
Before using this backend API, follow these steps to set
up your environment:

```bash
# 1. Make scripts executable
chmod +x db.sh
chmod +x migrate.sh

# 2. Install dependencies (if any)
make install

# 3. Install Air for live reloading during development
go install github.com/air-verse/air@latest

# 4. Set environment variables
cp .env.example .env
vim .env  # Edit with your DB credentials

# 5. (Optional) Create database for the first time
./db.sh create

```

---

## 🎯 Complete Database Commands Reference

This section provides a comprehensive reference for all database-related
commands available in this project. You can use either the Makefile
commands (with `db-` prefix) or the shell script commands (without
prefix) to manage your database.

### Using Makefile (with dash):

```bash
# Database lifecycle
make db-create        # Create database
make db-drop          # Drop database
make db-migrate       # Run migrations
make db-setup         # Create + migrate
make db-reset         # Drop + create + migrate

# Information
make db-status        # Show migration status
make db-version       # Show current version
make db-info          # Show database info

# Utilities
make db-console       # Open psql
make db-rollback      # Rollback last migration

# Alternative (without db- prefix)
make migrate          # Same as db-migrate
make migrate-status   # Same as db-status
make migrate-rollback # Same as db-rollback
```

### Using Shell Script (original names):

```bash
# Database lifecycle
./db.sh create        # Create database
./db.sh drop          # Drop database
./db.sh migrate       # Run migrations
./db.sh setup         # Create + migrate
./db.sh reset         # Drop + create + migrate

# Information
./db.sh status        # Show migration status
./db.sh version       # Show current version
./db.sh info          # Show database info

# Utilities
./db.sh console       # Open psql
./db.sh rollback      # Rollback last migration
```

---

