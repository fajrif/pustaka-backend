.PHONY: install run migrate seed

install:
	go mod download
	go mod tidy

run:
	go run main.go

migrate:
	psql -U postgres -d pustaka -f database/migrations.sql

seed:
	@echo "Seeding admin user..."
	@psql -U postgres -d pustaka -c "INSERT INTO users (email, password_hash, full_name, role) VALUES ('admin@pustaka.com', '\$$2a\$$10\$$YtcSqB5h7rKYOhB5YjW8/.fKj0Z4HQlGLj6kZ.rX6YqJXJqYQJhfy', 'Administrator', 'admin') ON CONFLICT (email) DO NOTHING;"

build:
	go build -o bin/pustaka main.go

clean:
	rm -rf bin/
