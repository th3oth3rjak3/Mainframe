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
# Build the API binary into the /bin folder:
build-win:
    go build -o bin/mainframe.exe ./cmd/api

build: 
    go build -o bin/mainframe ./cmd/api

build-hasher-win:
    go build -o bin/hasher.exe ./cmd/hasher

build-hasher:
    go build -o bin/hasher ./cmd/hasher

# Run the API in dev mode:
dev: swag
    go run ./cmd/api/main.go

# Build and run the compiled API binary (avoids go run):
run: swag build
    ./bin/mainframe

# Build and run the compiled API binary (avoids go run):
run-win: swag build-win
    ./bin/mainframe

run-hasher pw: build-hasher
    ./bin/hasher {{pw}}

run-hasher-win pw: build-hasher-win
    ./bin/hasher {{pw}}

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

