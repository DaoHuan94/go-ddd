## Database migrations with `github.com/golang-migrate/migrate`

This project uses PostgreSQL migrations via [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate).

### 1) Start PostgreSQL (via `compose.yaml`)

If you haven't already:

```bash
docker compose up -d
```

Your `compose.yaml` uses:

- database: `app`
- user: `app`
- password: `app`
- host/port: `localhost:5432`

So your DSN for the `migrate` CLI is typically:

```text
postgres://app:app@localhost:5432/app?sslmode=disable
```

### 2) Create a migrations folder

```bash
mkdir -p migrations
```

### 3) Install the `migrate` CLI (recommended)

The `migrate` repository provides a CLI tool you can install with Go:

```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

If this is the first time, make sure your `GOPATH/bin` (or equivalent) is in `PATH` so you can run `migrate`.

### 4) Create a new migration

Examples:

```bash
migrate create -ext sql -dir migrations create_users_table
```

This will create two files in `./migrations/` (names include a timestamp/version):

- `*.up.sql` (what to apply)
- `*.down.sql` (how to rollback)

### 5) Write SQL in the `up` and `down` files

Edit the generated files:

```sql
-- migrations/XXXXXXXXXX_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
  id   BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL
);
```

```sql
-- migrations/XXXXXXXXXX_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

### 6) Apply migrations

Run all pending migrations:

```bash
migrate -path migrations -database "postgres://app:app@localhost:5432/app?sslmode=disable" up
```

### 7) Check migration status / current version

```bash
migrate -path migrations -database "postgres://app:app@localhost:5432/app?sslmode=disable" version
```

### 8) Roll back migrations

Rollback by 1 step:

```bash
migrate -path migrations -database "postgres://app:app@localhost:5432/app?sslmode=disable" down 1
```

Rollback to a specific version:

```bash
migrate -path migrations -database "postgres://app:app@localhost:5432/app?sslmode=disable" down 202603190001
```

### 9) (Optional) Run migrations from Go code

If you prefer running migrations programmatically instead of using the CLI:

1. Add the dependency:

```bash
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
```

2. Use code like this (example):

```go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dsn := "postgres://app:app@localhost:5432/app?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// file:// points to a directory holding migrations.
	m, err := migrate.New("file://migrations", "postgres", db)
	if err != nil {
		panic(fmt.Errorf("migrate.New: %w", err))
	}

	if err := m.Up(); err != nil {
		// It's common for m.Up() to return migrate.ErrNoChange when everything is already applied.
		// You can ignore that error if you want.
		panic(err)
	}
}
```

Notes:

- The Go example uses the `file://migrations` source and the `postgres` database driver.
- You still need valid `*.up.sql`/`*.down.sql` files in `./migrations`.

