# Change Log

Running log of features, bug fixes, and other changes that have been
**committed, pushed, and merged** into `quinar_main_branch`. Newest entries on
top. Each entry records what changed, why, and any verification performed.

> Policy: every time a change is committed, pushed, and merged, add a new entry
> here (see CLAUDE.md → "Change log policy").

---

## 2026-07-14 — Auth module: JWT login, users domain, role-based access

**Branch / PR:** `feature/auth-roles` → merged into `quinar_main_branch` via
PR #26 (merge `bfafff5`; source commit `b4b5cdd`).

### Changes
- **feat (entities/db):** new `User` entity (`pkg/entities/user.go`) with
  unique email, bcrypt `password_hash`, `role`
  (`ADMIN`/`USER`/`OPERATIVE_USER`) and `is_active`; registered in
  `AutoMigrate` (`db/postgres.go`) → creates the `users` table.
- **feat (pkg/user):** repository + service mirroring the `vendor` pattern,
  plus `FindByEmail`, `CountActiveByRole`, `SetPassword`, `Authenticate`
  (generic "credenciales inválidas" error to avoid user enumeration) and a
  guard that refuses to deactivate the **last active ADMIN**
  (`ErrLastActiveAdmin` → HTTP 409). `SeedInitialAdmin` creates the first
  admin from `SEED_ADMIN_EMAIL`/`SEED_ADMIN_PASSWORD` env vars (idempotent),
  called from `cmd/unknown-api/main.go`.
- **feat (pkg/auth):** JWT HS256 generation/validation
  (`github.com/golang-jwt/jwt/v5`, new dependency) + `RequireAuth()` and
  `RequireRole()` Fiber middlewares. Role travels inside the token (a role
  change applies on next login).
- **feat (api):** presenters/handlers/routes for `POST /api/v1/auth/login`
  (public), `GET /api/v1/auth/me`, and admin-only `/api/v1/users` CRUD +
  `PUT /users/:id/password`. Password hash never exposed.
- **feat (server):** in `fiberServer.go`, auth routes register first, then
  `api.Use(auth.RequireAuth())` protects every `/api/v1/*` route registered
  after it. `/locations` additionally requires `ADMIN`
  (`api/routes/location_routes.go`). `/public/*` and `/api/v1/health` remain
  public.
- **feat (config):** new `auth` section (`jwt_secret`, `jwt_expiry_hours`) in
  `config.yaml`/`config.go`; in production the secret is set via
  `AUTH_JWT_SECRET`.
- **docs:** PROJECT_OVERVIEW (architecture note, `users` table, env vars),
  DATABASE_REFERENCE (§3.13 `users`, `UserRole` enum, integrity rule),
  API_REFERENCE (auth conventions, 401/403/409, Auth & Users sections).

### Verification
| Check | Result |
|-------|--------|
| `go build ./...` / `go vet ./...` | ✅ exit 0 |
| Server boot → `users` table migrated + initial admin seeded | ✅ |
| `POST /auth/login` (seed admin) | ✅ 200 + token (local & VPS) |
| `GET /materials` sin token | ✅ 401 (antes 200 público) |
| `GET /materials` con token | ✅ 200 |
| `GET /auth/me` | ✅ claims del token |
| `POST /users` con admin (crear `OPERATIVE_USER`/`USER`) | ✅ 201 |
| `POST /locations` / `GET /users` con token rol `USER` | ✅ 403 |
| `DELETE /users/:id` del único admin activo | ✅ 409 con mensaje del guard |
| `/public/sync/purchase-orders` y `/health` sin token | ✅ 200 |
| Post-deploy en VPS (`https://api.quinar.hvdevs.com`) | ✅ misma matriz: 401 sin token, login OK, 200 con token |

### Notes / follow-ups (not blocking)
- Fine-grained `OPERATIVE_USER` permissions over `/assets` /
  `/asset-movements` are **pending** (today those routes only require login) —
  separate design task.
- `SEED_ADMIN_EMAIL`/`SEED_ADMIN_PASSWORD` must be exported at least once on a
  fresh DB; afterwards the seed is a no-op. `AUTH_JWT_SECRET` was set on the
  VPS (the value in `config.yaml` is a placeholder).
- No refresh tokens, email-based password recovery, rate limiting or public
  signup — out of scope by design.

---

## 2026-06-15 — Persist `erp_code` on material update

**Branch / PR:** `fix_material_erp_update` → to be merged into
`quinar_main_branch` (merge hash TBD).

### Changes
- **fix (repository):** `material.Update` omitted `erp_code` from its `Updates`
  map, so `PUT /api/v1/materials/:id` accepted and echoed a new `erpCode` but
  never persisted it (silent no-op in the DB). Added `"erp_code"` to the map in
  `pkg/material/repository.go`.

### Context
- `erp_code` is **optional at creation** (only `name`, `sector`,
  `unit_of_measure` and `code` are required; `code` must be unique) and is now
  **editable** afterwards via `PUT`.
- Regression-safe: the update handler loads the full material before applying
  changes, so partial updates (and `Deactivate`, which reuses `Update`) preserve
  the existing `erp_code` rather than blanking it.

### Verification — "does it run?" (on `fix_material_erp_update`)
| Check | Result |
|-------|--------|
| `go build ./...` | ✅ exit 0 |
| `go vet ./...` | ✅ exit 0 |
| Server boot + DB connection (remote `quinar`) | ✅ |
| `PUT /materials/:id` with `erpCode` then re-`GET` | ✅ new value persisted |
| `PUT /materials/:id` with only `name` | ✅ `erp_code` preserved |

### Notes / follow-ups (not blocking)
- `internal_code` is similarly absent from the update map, but it is intentionally
  left out: it seeds asset serial generation, so editing it after assets exist is
  not a safe one-line change. Out of scope here.

---

## 2026-06-14 — Material ERP-code bugfixes + new `need` purchase domain

**Branch / PR:** `needs_features` → merged into `quinar_main_branch` via PR #22
(merge `09c9ece`; source commits `d5a47f0`, `b5a0129`, `2cf6088`, `a0be250`).

### Changes

**Bugfixes (material ERP lookup):**
- **fix (routing):** `GET /materials/:erp_code` collided with
  `GET /materials/:id` (Fiber matches the first single-segment pattern), so the
  ERP lookup was unreachable. Moved it to `GET /materials/by-erp/:erp_code`
  (`api/routes/material_routes.go`).
- **fix (repository):** `material.FindByERPCode` used `db.First` on a slice and
  returned a single row. Switched to `Where().Find()` so all materials sharing a
  Bejerman "bag" ERP code (e.g. `0 MAT GOP21`) are returned
  (`pkg/material/repository.go`).
- **fix (encoding):** `GetMaterialsByERP` now `url.PathUnescape`s the param
  before querying, since Fiber's default config does not unescape path params —
  spaced ERP codes must be percent-encoded by the client
  (`api/handlers/material_handler.go`).

**New domain — `need` (pre-Purchase-Requisition stage):**
- **entities:** `Need`, `NeedItem`, `NeedCounter` (per-year correlative). `Need`
  carries a server-generated `Number` (`NEED-YYYY-NNNN`) and a status lifecycle
  `IN_PROGRESS → READY → CONVERTED` (plus `DISCARDED`).
- **repository:** `CreateWithNumber` generates the number under a transactional
  row lock (same pattern as `AssetSerialCounter`); reads preload
  `Items.Material`; item partial-update preserves omitted optional references.
- **service:** required-field validation, material existence/active checks, and
  controlled transitions `Promote` / `MarkConverted` / `Discard`.
- **API layer:** presenter, handler and routes following the `material` pattern;
  item responses embed denormalized `material_name` / `material_erp_code`.
  Routes under `/api/v1/needs` (+ `/:id/items`, `/items/:itemId`, and the
  `promote` / `mark-converted` / `discard` actions).
- **wiring/migration:** `repo → service → handler` chain added to
  `server/fiberServer.go`; `Need`, `NeedItem`, `NeedCounter` registered in the
  `AutoMigrate(...)` call in `db/postgres.go`.
- **docs:** added Needs + Need Items sections (and index row) to
  `API_REFERENCE.md`; documented the `by-erp` route change.

> Out of scope this iteration (planned for later): `NeedAttachment`, the outbound
> webhook fired on promote, and authentication/authorization.

### Verification — "does it run?" (on `needs_features`)
| Check | Result |
|-------|--------|
| `go build ./...` | ✅ exit 0 |
| `go vet ./...` | ✅ exit 0 |
| DB connection (remote `quinar`) | ✅ "Conexion a DB Exitosa!" |
| `AutoMigrate` creates `needs`, `need_items`, `need_counters` | ✅ "Migracion exitosa" |
| Server boot | ✅ Fiber v2.52.5 on `0.0.0.0:4045` |
| `GET /api/v1/health` | ✅ HTTP 200, body `Working Cool!` |
| `GET /materials/by-erp/<code>` returns all matches (incl. `%20`-encoded) | ✅ |
| Need flow: create → add item → promote → mark-converted | ✅ statuses/timestamps correct |
| Need correlative across multiple creates | ✅ `NEED-2026-0001/0002/0003` |
| Need guard errors (promote-empty, item on promoted, invalid material, etc.) | ✅ 400 with clear messages |

### Notes / follow-ups (not blocking)
- The `Need` entity gained an `Items []NeedItem` has-many association (not in the
  original spec) because §2.2 requires `FindAll`/`FindByID` to preload items;
  table names are unaffected.
- `selected_cost_id` on a need item cannot currently be *cleared* via update
  (omitting it preserves the existing value); revisit if "unset price" is needed.

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
