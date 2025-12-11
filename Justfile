# Conditionally set the shell based on the operating system
set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]
set shell := ["sh", "-c"]

# Default task.
default: dev

# --- DATABASE MIGRATIONS ------------------------------------------------------

# Create a new goose SQL migration:
#   just migrate-create add_users_table
migrate-create name:
    goose create "{{name}}" sql -dir migrations

# Run all migrations:
migrate-up:
    goose -dir migrations sqlite3 data/mainframe.db up

# Roll back the most recent migration:
migrate-down:
    goose -dir migrations sqlite3 data/mainframe.db down

# Reset the database entirely:
migrate-reset:
    goose -dir migrations sqlite3 data/mainframe.db reset


# --- SWAGGO / OPENAPI ---------------------------------------------------------

# Regenerate Swagger docs:
swag:
    swag init -g cmd/api/main.go -o docs


# --- SERVER -------------------------------------------------------------------

# Run the API in dev mode:
dev: swag
    go run ./cmd/api/main.go

# Build the API binary:
build:
    go build -o bin/mainframe ./cmd/api


# --- TESTING / LINTING --------------------------------------------------------

fmt:
    go fmt ./...

test:
    go test ./...

vet:
    go vet ./...

staticcheck:
    staticcheck ./...

lint: vet staticcheck
