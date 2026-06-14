# Change Log

Running log of features, bug fixes, and other changes that have been
**committed, pushed, and merged** into `quinar_main_branch`. Newest entries on
top. Each entry records what changed, why, and any verification performed.

> Policy: every time a change is committed, pushed, and merged, add a new entry
> here (see CLAUDE.md → "Change log policy").

---

## 2026-06-14 — Project documentation + purchase-order sync wiring

**Branch / PR:** `docs/project-overview` → merged into `quinar_main_branch`
via PR #21 (merge `ce284fb`, source commit `166bdaa`).

### Changes
- **docs:** Added project documentation set under `docs/`:
  - `PROJECT_OVERVIEW.md` — purpose, architecture, tech stack, data-model summary.
  - `API_REFERENCE.md` — full reference of every route, payload, and response.
  - `DATABASE_REFERENCE.md` — complete schema/tables/relations reference for agents.
  - `CLAUDE.md` (repo root) — guidance for future Claude Code sessions.
- **fix (wiring):** `server/fiberServer.go` now constructs the
  `SyncPurchaseOrderHandler` (repo → service → handler) and passes it to
  `PublicSyncPurchaseOrdersRoute`. Previously `spoHandler` was referenced but
  never created, so the `server` package did not compile.
- **fix (migration):** Registered `&entities.PurchaseOrderSync{}` in the
  `AutoMigrate(...)` call in `db/postgres.go` so the `purchase_order_syncs`
  table is created at startup. Without it the now-reachable sync endpoints
  would fail with `relation "purchase_order_syncs" does not exist`.

### Verification — "does it run?" (on `quinar_main_branch` after merge)
| Check | Result |
|-------|--------|
| `go build ./...` | ✅ exit 0 |
| `go vet ./...` | ✅ exit 0 |
| DB connection (remote `quinar`) | ✅ "Conexion a DB Exitosa!" |
| `AutoMigrate` incl. new `PurchaseOrderSync` | ✅ "Migracion exitosa" |
| Server boot | ✅ Fiber v2.52.5 on `0.0.0.0:4045`, 87 handlers |
| `GET /api/v1/health` | ✅ HTTP 200, body `Working Cool!` |
| `GET /public/sync/purchase-orders` (newly wired) | ✅ HTTP 200, returns records |

### Notes / follow-ups (not blocking)
- Slow-SQL warnings (~200–445 ms) during AutoMigrate are GORM introspecting the
  remote DB over the network at startup; not a runtime concern.
- `GET /public/sync/purchase-orders` serializes raw entity field names
  (`"ID"`, `"CCO"`, …) instead of snake_case because `GetPurchaseOrders`
  returns entities without a presenter — casing differs from the rest of the API.
- Minor open items in `pkg/sync_purchase_order/service.go`: `SyncOrders`
  swallows a bad date-parse error (stores zero time); `BatchCreate` on an empty
  array returns a GORM "empty slice" error (would surface as a 500).
