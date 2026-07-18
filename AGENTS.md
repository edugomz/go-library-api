# AGENTS.md

This file provides guidance to coding agents (including Claude Code) when working with code in this repository.

## Commit style

Use Conventional Commits: `<type>(<scope>): <description>`

Common types: `feat`, `fix`, `chore`, `refactor`, `test`, `docs`

Examples:
- `feat(auth): add JWT login and register endpoints`
- `fix(middleware): return 401 on missing Bearer prefix`
- `chore(deps): add golang-jwt/jwt/v5`

## Commands

```bash
# Start dev infrastructure (PostgreSQL on :5432, pgAdmin on :5050)
docker-compose up -d

# Run database migrations (one-off; server no longer runs them on boot)
go run ./cmd/api --migrate-only

# Run the application
go run ./cmd/api/main.go

# Run all tests
go test ./...

# Run a single test
go test ./internal/service/... -run TestGetUsers
go test ./internal/repository/... -run TestCreateUser

# Integration tests require the test DB first
docker-compose -f docker-compose.test.yml up -d
go test ./internal/repository/...
```

## Configuration

The app reads from `.env` or environment variables. Required vars:

| Variable | Default | Notes |
|---|---|---|
| `DB_HOST` | — | Required |
| `DB_USER` | — | Required |
| `DB_NAME` | — | Required |
| `DB_PASSWORD` | — | Optional |
| `DB_PORT` | `5432` | |
| `APP_PORT` | `8080` | |
| `LOG_LEVEL` | `info` | `debug\|info\|warn\|error` |
| `JWT_SECRET` | — | Required; used to sign/verify HS256 tokens |

Default docker-compose credentials: `postgres/postgres`, DB: `library`.

## Architecture

The dependency chain flows strictly one way: `handler → service → repository → DB`.

```
cmd/api/main.go          Entry point: bootstraps app, starts server (or runs migrations and exits if --migrate-only)
internal/app/
  app.go                 Application struct; wires config, DB, routes, and Gin engine
  handler.go             Constructs all handlers via explicit dependency injection
  routes.go              Registers route groups: /api/v1/* and / (web)
internal/config/         Loads and validates env-based config; builds DSN
internal/db/             Opens GORM connection; exposes package-level db.DB singleton
internal/models/         GORM model structs (User, Author, Book, Review, ReadingList)
internal/repository/     DB access layer — concrete GORM implementations
internal/service/        Business logic; defines repository interfaces (not the repo package)
internal/handlers/api/   Gin HTTP handlers for REST endpoints (auth, user, book, author)
internal/middleware/      Gin middleware; auth.go validates Bearer JWT and sets userID in context
internal/handlers/web/   Gin handler for HTML views (internal/views/*.html)
migrations/              Versioned SQL migrations (golang-migrate); run explicitly via `go run ./cmd/api --migrate-only`, not on server boot
```

**Key design decision:** repository interfaces are defined in the `service` package (not `repository`), so services depend only on abstractions. This is why unit tests for services use in-package mocks without importing the repository package.

**Three test strategies in use:**
- `internal/service/*_test.go` — unit tests with in-memory mock repositories (no DB needed)
- `internal/repository/*_test.go` — integration tests that connect to a real PostgreSQL instance on port `5433` (started via `docker-compose.test.yml`); share DB setup via `setupTestDB(t)` in `testdb_test.go` (the DB isn't truncated between runs, so tests use unique values via `uniqueSuffix()` rather than asserting exact row counts)
- `internal/handlers/api/*_test.go` and `internal/middleware/auth_test.go` — HTTP-layer tests using `httptest` + a real `gin.Engine`; handlers are exercised with in-package mocks of the `service` package's exported repository interfaces (e.g. `service.AuthorRepository`), not a mocked service layer

## CI

`.github/workflows/test.yml` runs on every PR and on push to `main`: spins up `docker-compose.test.yml`'s `db-test` service, then runs `go build ./...` and `make test`.

## Adding a new resource

Follow the pattern established by User/Book/Author:
1. Add model struct in `internal/models/`
2. Add repository in `internal/repository/` implementing an interface
3. Define the interface and service in `internal/service/`
4. Add handler in `internal/handlers/api/`
5. Wire handler in `internal/app/handler.go` → `NewHandlers()`
6. Register routes in `internal/app/routes.go`
7. Add a new numbered up/down SQL migration pair under `migrations/sql/` (e.g. `000002_add_x.up.sql` / `.down.sql`)

## Web frontend

Server-rendered `html/template` pages (Gin `LoadHTMLGlob`) + htmx (CDN `<script>` tag, no build step). See `ui-roadmap.md` for what's built and what's deferred.

- `internal/views/nav.html` defines a shared `{{define "nav"}}` partial included via `{{template "nav" .}}` at the top of every page's `<body>`; pass `"LoggedIn": h.isLoggedIn(c)` in the template data for it to render correctly.
- `static/` holds the one stylesheet, served via `r.Static("/static", "./static")` in `app.go`.
- Web routes (`internal/handlers/web`) call `service` packages directly, the same way the API handlers do, they don't go through `/api/v1/*` over HTTP and aren't behind `middleware.RequireAuth`.
- Login/register set/clear an `HttpOnly` JWT cookie (`token`) directly from `AuthService`; no web route currently enforces it, it only toggles nav display. See the auth tradeoff note in `ui-roadmap.md` before changing this.

## Maintaining this file

Keep this file for knowledge useful to almost every future agent session in this project.
Do not repeat what the codebase already shows; point to the authoritative file or command instead.
Prefer rewriting or pruning existing entries over appending new ones.
When updating this file, preserve this bar for all agents and keep entries concise.
