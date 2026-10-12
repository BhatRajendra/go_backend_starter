# Go Backend Starter Pack — Setup Guide

This guide sets up a reusable Go backend with Chi, PostgreSQL (`pgxpool`), `slog` logging, Air, Goose, and SQLC.

## 1. Navigate to the repository

```bash
cd ~/Yuddhaa/go_backend_starter
```

If you are creating the Go module for the first time:

```bash
go mod init github.com/Yuddhaa/go_backend_starter
```

> Run `go mod init` only if `go.mod` does not already exist. If you use a different GitHub repository path, use that path as the module name.

Check the installed Go version:

```bash
go version
```

## 2. Install the Go dependencies

```bash
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/joho/godotenv
```

These dependencies provide:

- `chi/v5` — HTTP router and middleware.
- `pgxpool` — PostgreSQL connection pooling.
- `godotenv` — Loading environment variables from `.env`.

The following are part of Go's standard library, so they do not require installation:

- `log/slog` — Structured logging.
- `net/http` — HTTP server.
- `context` — Request and database contexts.
- `runtime` — Caller information for logging.

## 3. Install development tools

### Air — hot reloading

```bash
go install github.com/air-verse/air@latest
```

Verify:

```bash
air -v
```

Generate the initial configuration if needed:

```bash
air init
```

This creates `.air.toml`. Configure it to build your application's `main.go`, run the resulting binary, watch Go source files, and exclude temporary build output and `.git`.

Start the development server with:

```bash
air
```

### Goose — database migrations

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Verify:

```bash
goose -version
```

### SQLC — type-safe database code generation

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Verify:

```bash
sqlc version
```

If any installed command is not found, make sure your Go binary directory is in `PATH`:

```bash
go env GOPATH
```

For a typical macOS setup using zsh:

```bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## 4. Create the directory structure

```bash
mkdir -p internal/app
mkdir -p internal/logging
mkdir -p internal/db

mkdir -p sql/migrations
mkdir -p sql/queries
mkdir -p sql/schema
```

The resulting structure should be:

```text
go_backend_starter/
├── main.go
├── go.mod
├── go.sum
├── .env
├── .envEg
├── .gitignore
├── .air.toml
├── sqlc.yaml
├── internal/
│   ├── app/
│   │   ├── helper.go
│   │   ├── middleware.go
│   │   └── server.go
│   ├── logging/
│   │   └── log.go
│   └── db/                 # SQLC-generated code
└── sql/
    ├── migrations/         # Goose migration files
    ├── queries/            # SQL queries for SQLC
    └── schema/
        └── schema.sql     # Schema source for SQLC
```

If your current `sqlc.yaml` generates code into a root-level `db/` directory, either retain that layout or change the configuration to `internal/db`. The output directory and Go import paths must agree.

## 5. Configure environment variables

Keep `.envEg` as a template and create your local environment file:

```bash
cp .envEg .env
```

Edit `.env` and configure the variables your application requires, for example:

```env
APP_ENV=DEV
PORT=4000
DB=postgresql://USER:PASSWORD@HOST/DATABASE?sslmode=require
```

Replace the database URL with your actual PostgreSQL connection string.

Ensure `.env` is excluded from Git. Keep `.envEg` free of real passwords, tokens, and connection credentials.

Your `main.go` should load the environment variables before using them, validate required values, initialize the connection pool, and ping the database.

## 6. Configure SQLC

Keep `sqlc.yaml` in the repository root. A configuration for the `internal/db` output directory can look like this:

```yaml
version: "2"

sql:
  - engine: "postgresql"
    schema: "sql/schema/schema.sql"
    queries: "sql/queries"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
```

If you already have a working `sqlc.yaml`, retain it and verify that the paths match your actual directories.

**Important:** Keep `sql/schema/schema.sql` as the schema source for SQLC, even when removing the Sandhi-specific tables. You can leave a comment in it until the next project defines its schema.

Generate Go code after adding the schema and SQL queries:

```bash
sqlc generate
```

You can also explicitly specify the configuration file:

```bash
sqlc generate -f sqlc.yaml
```

SQLC reads the schema and query files, then generates type-safe Go code. It does not apply migrations to the database.

## 7. Set up Goose migrations

Remove the Sandhi-specific migration from the starter pack:

```bash
rm -f sql/migrations/001_initial_schema.sql
```

Keep the `sql/migrations/` directory available for future projects.

When you start a new project and are ready to define its schema, create a migration:

```bash
goose -dir sql/migrations create init_schema sql
```

Goose creates a timestamped SQL migration file. Add the project's schema changes to that file.

To apply migrations to your configured database, load the environment variables into your shell:

```bash
set -a
source .env
set +a
```

Then run:

```bash
goose -dir sql/migrations postgres "$DB" up
```

Useful Goose commands:

```bash
# Show migration status
goose -dir sql/migrations postgres "$DB" status

# Roll back the latest migration
goose -dir sql/migrations postgres "$DB" down

# Apply all pending migrations
goose -dir sql/migrations postgres "$DB" up
```

Use `down` carefully, especially against production databases.

## 8. Keep SQLC schema and Goose migrations consistent

These two mechanisms serve different purposes:

- `sql/migrations/` contains the incremental database changes applied by Goose.
- `sql/schema/schema.sql` describes the database schema SQLC uses to generate code.
- `sql/queries/` contains the SQL statements SQLC converts into Go methods.
- `internal/db/` contains the generated Go code.

When you add or change a table, update the migration and the SQLC schema source consistently. Then add queries and regenerate the code:

```bash
sqlc generate
```

Do not expect `sqlc generate` to modify your actual database; Goose handles that.

## 9. Format and tidy the project

Run these commands after setting up the source files:

```bash
gofmt -w main.go internal/
go mod tidy
```

Check that the project builds:

```bash
go build ./...
```

Run the tests:

```bash
go test ./...
```

Run the backend in development mode:

```bash
air
```

## 10. Initial Git cleanup

Check that project-specific files are no longer included:

```bash
git status
```

If the original Sandhi migration was already tracked, removing it will show up as a deletion. Commit that change when the starter pack is ready.

Ensure these are not committed:

- `.env` with real credentials.
- Compiled binaries.
- Air's temporary build directory, commonly `tmp/`.
- Project-specific database credentials or secrets.

Keep `.envEg`, `.air.toml`, `sqlc.yaml`, the generic application scaffolding, and the logging implementation.

## Quick setup checklist

- [ ] Go module initialized with the starter repository's module path.
- [ ] Chi, pgxpool, and godotenv installed.
- [ ] Air, Goose, and SQLC installed.
- [ ] `.env` created locally and ignored by Git.
- [ ] `.air.toml` configured.
- [ ] `sqlc.yaml` retained and paths verified.
- [ ] Sandhi-specific migration removed.
- [ ] `sql/schema/schema.sql` retained as the SQLC schema source.
- [ ] `go mod tidy`, `go build ./...`, and `go test ./...` run successfully.
