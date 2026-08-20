# useradmin example

A self-contained example server that boots `useradmin` on an in-memory
SQLite database. No external services required.

## Run

```bash
go run ./example
```

Then open http://localhost:8080/ in your browser.

## Hot reload (recommended for development)

The repo includes a [`taskfile.yml`](../taskfile.yml) and [`.air.toml`](../.air.toml)
for hot-reload development with [Air](https://github.com/air-verse/air)
and [Task](https://taskfile.dev).

Install the tools once:

```bash
# Install task:  https://taskfile.dev/installation/
# Install air:
task air:install
```

Then start the example with hot reload:

```bash
task dev
```

Or use Air directly:

```bash
air
```

## What it does

- Creates an in-memory SQLite database (reset on every restart)
- Initializes `userstore`, `geostore` (with auto-seed), and
  `sessionstore` with auto-migration
- Seeds 40 sample users (mixed statuses) plus 1 administrator
- Mounts the user admin panel at `/admin/users`
- Serves a landing page at `/` that links into the admin

## Authentication

The example is open (no authentication) so you can click around
immediately. In a real integration, provide `AuthUserID` and `AuthUser`
callbacks that read your session/JWT from the request context, and gate
the `/admin/users` route behind your auth middleware.

## Persistence

To use a file-based database instead of in-memory, change `dbFile`:

```go
const dbFile = "useradmin_example.db"
```

Data will persist across restarts.
