# CLAUDE.md

Guidance for Claude Code (claude.com/claude-code) when working in this repo.

## Project

`unknown-go-server` — a Go REST API (module
`github.com/unknowncode44/unknown-go-server`) backing **Quinar**, an internal
system for managing materials, vendors, costs, serialized assets, asset
movements, locations, bulk inventory and purchase-order syncs. The configured
database is `quinar` (PostgreSQL).

A full, standalone description of the domain, architecture and data model lives
in [docs/PROJECT_OVERVIEW.md](docs/PROJECT_OVERVIEW.md) — read it first for
context. Asset-flow specifics are in
[README_ASSET_CYCLE.md](README_ASSET_CYCLE.md).

## Tech stack

- **Go 1.22.1**
- **Fiber v2** — HTTP framework (CORS enabled). API base path `/api/v1`,
  plus a `/public` group; health check at `GET /api/v1/health`.
- **GORM v1.25** over **PostgreSQL** (pgx driver). Schema is created/updated
  with `AutoMigrate` at startup — there is no separate migration tool.
- **Viper** for config (`config.yaml` + env overrides).
- **google/uuid** + Postgres `uuid-ossp` extension for UUID primary keys.

## Architecture (clean / layered)

Each domain is a vertical slice. Dependencies point inward via interfaces.

```
api/routes → api/handlers (+ api/presenter DTOs)
          → pkg/<domain>/service.go      (business rules; no Fiber, no GORM)
          → pkg/<domain>/repository.go   (GORM persistence behind an interface)
          → pkg/entities/*.go            (GORM models = the schema)
db/ + config/ + server/ wire everything together
```

- **Composition root:** `server/fiberServer.go` builds every
  `repo → service → handler` chain and registers routes. Add new domains here.
- **Entry point:** `cmd/unknown-api/main.go` → config → DB → Fiber server.
- **Migrations:** register new entities in the `AutoMigrate(...)` call in
  `db/postgres.go`, otherwise their tables won't be created.
- Services depend on repository **interfaces**; some services take multiple
  repos (e.g. `asset` needs material repo; `material_cost` needs
  vendor-material, vendor, currency repos).

## Conventions

- **UUID primary keys** via `uuid_generate_v4()` (except `PurchaseOrderSync`,
  which uses an auto-increment `uint`).
- **Soft deletes:** prefer a `Deactivate` use case (`IsActive = false`) over
  physical deletion to preserve history/referential integrity.
- **Append-only history:** `MaterialCost` is never updated — insert a new row.
- **Layer purity:** repositories must not contain business rules; services must
  not import Fiber or GORM; handlers only do HTTP I/O and call services.
- **Presenters** (`api/presenter`) are the API request/response contract —
  keep them separate from `pkg/entities`.
- Each domain folder under `pkg/` contains `repository.go` + `service.go`;
  matching `handlers`, `routes`, `presenter` files live under `api/`.
- Some source comments are in Spanish; match the surrounding style when editing.

## Adding a new domain (checklist)

1. Entity in `pkg/entities/<name>.go` (GORM tags).
2. Add it to `AutoMigrate(...)` in `db/postgres.go`.
3. `pkg/<name>/repository.go` (interface + GORM impl) and `service.go`
   (interface + business logic).
4. `api/handlers/<name>_handler.go`, `api/presenter/<name>_presenter.go`,
   `api/routes/<name>_routes.go`.
5. Wire repo → service → handler and register routes in
   `server/fiberServer.go`.

## Common commands

```bash
# Run the server (loads config.yaml, migrates, listens on config port 4045)
go run ./cmd/unknown-api

# Build
go build ./...

# Tidy / verify modules
go mod tidy

# Format & vet
gofmt -w .
go vet ./...
```

> No automated test suite or CI is present in the repo at this time.

## Configuration

`config.yaml` at the repo root (env vars override via Viper). Holds `server.port`
and the `db.*` connection settings (host, port, user, password, dbname,
timezone). Do not commit real credentials.

## Git / workflow notes

- The active development branch is `quinar_main_branch`; the `main` branch is
  stale (only legacy `Company`/`Branch` entities) and does **not** reflect the
  real project — branch from `quinar_main_branch`.
- GitHub: https://github.com/unknowncode44/unknown-go-server
- Create a feature branch for changes; only commit/push when the user asks.

## Change log policy

Any time a change (feature, bug fix, or any other edit) is **committed, pushed,
and merged** into `quinar_main_branch`, add a new entry to
[docs/CHANGELOG.md](docs/CHANGELOG.md). Newest entry on top. Each entry should
record:

- date, branch/PR (and commit/merge hashes if known);
- what changed and why;
- any verification performed (build/vet/run, endpoints hit, etc.);
- notes or non-blocking follow-ups.

Only log changes that actually landed on `quinar_main_branch` (post-merge), not
work-in-progress.
