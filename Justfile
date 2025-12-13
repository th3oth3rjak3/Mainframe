# Conditionally set the shell based on the operating system
set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]
set shell := ["sh", "-c"]

# Default task.
default: dev

# --- DATABASE MIGRATIONS ------------------------------------------------------
# Create a new goose SQL migration:
#   just migrate-create add_users_table
migrate-create name:
    goose create "{{name}}" sql -dir internal/data/migrations

# Run all migrations:
migrate-up:
    goose -dir internal/data/migrations sqlite3 internal/data/mainframe.db up

# Roll back the most recent migration:
migrate-down:
    goose -dir internal/data/migrations sqlite3 internal/data/mainframe.db down

# Reset the database entirely:
migrate-reset:
    goose -dir internal/data/migrations sqlite3 internal/data/mainframe.db reset

# --- SWAGGO / OPENAPI ---------------------------------------------------------
# Regenerate Swagger docs:
swag:
    swag init -g cmd/api/main.go -o internal/docs

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

build-id-win:
    go build -o bin/id_generator.exe ./cmd/id_generator

build-id:
    go build -o bin/id_generator ./cmd/id_generator

build-hmac:
    go build -o bin/hmac_key ./cmd/hmac_key

build-hmac-win:
    go build -o bin/hmac_key.exe ./cmd/hmac_key

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

run-id: build-id
    ./bin/id_generator

run-id-win: build-id-win
    ./bin/id_generator

run-hmac: build-hmac
    ./bin/hmac_key

run-hmac-win: build-hmac-win
    ./bin/hmac_key

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

