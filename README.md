# Squirrel CLI

This repository now provides both a simple GraphQL API server backed by PostgreSQL and a command line client.

The CLI communicates with the server instead of storing data locally. Login stores a temporary session token under `data/`.

## Commands

```bash
squirrel create-account --name "Jane" --email jane@example.com --password secret --squirrel "Sprinkles"
squirrel login --email jane@example.com --password secret
squirrel add "Buy milk" --tags errands,personal
```

Items starting with `http://` or `https://` are treated as articles. Everything else is saved as a note.

Run the API server with:

```bash
go run ./cmd/server
```

The server listens on `localhost:8080` and expects a PostgreSQL instance configured via `DATABASE_URL` (defaults to `postgres://postgres:postgres@localhost:5432/squirrel?sslmode=disable`).
