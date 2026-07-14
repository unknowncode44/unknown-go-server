# Unknown Go Server — Project Overview

> Backend API for **Quinar**: a materials, assets, costing and inventory
> management system. This document is a self-contained summary of the
> project — its purpose, architecture, technology stack, and data model
> (tables and relationships). It is meant to be read standalone (e.g. in a
> Claude for Desktop project) without browsing the source tree.

---

## 1. What this project is

`unknown-go-server` is a REST API written in **Go**. It manages the
operational data of a service/installation company:

- **Materials** — the internal catalog of items the company buys and installs.
- **Vendors** — suppliers, and the per-vendor codes/prices for each material.
- **Costs & Currencies** — append-only historical pricing of materials per
  vendor, in multiple currencies.
- **Assets** — physical, individually-tracked instances of materials
  (serialized, with a QR-friendly visible serial), their lifecycle status and
  current location.
- **Asset Movements** — the audit trail of everything that happens to an asset
  (inbound, transfer, install, repair, decommission, etc.).
- **Locations** — hierarchical places (internal warehouses or client sites)
  used as movement origins/destinations.
- **Material Inventory & Delivery Records** — bulk (non-serialized) stock
  tracking for materials measured by quantity, with inbound/outbound deliveries.
- **Purchase Order Sync** — ingestion endpoint that receives purchase-order
  rows (e.g. from an Excel/VBA macro) and links them to materials and vendors.
- **Users & Auth** — system accounts with email/password login (bcrypt),
  JWT-based authentication and 3 roles (`ADMIN`, `USER`, `OPERATIVE_USER`).

The system is the backend half of an internal tool; the configured database is
named `quinar`.

---

## 2. Architecture

The codebase follows a **clean, layered architecture** with strict separation
of concerns and dependency inversion (upper layers depend on interfaces, not
concrete implementations). Each business domain is a self-contained vertical
slice.

```
HTTP request
   │
   ▼
┌──────────────────────────────────────────────────────────────┐
│  api/routes        → registers Fiber route groups per domain   │
│  api/handlers      → parse request, call service, build reply  │
│  api/presenter     → request/response DTOs (decouple from DB)  │
└──────────────────────────────────────────────────────────────┘
   │  (depends on service interfaces)
   ▼
┌──────────────────────────────────────────────────────────────┐
│  pkg/<domain>/service.go     → business rules / use cases      │
└──────────────────────────────────────────────────────────────┘
   │  (depends on repository interfaces)
   ▼
┌──────────────────────────────────────────────────────────────┐
│  pkg/<domain>/repository.go  → persistence (GORM); no business │
│                                 rules, no HTTP knowledge       │
└──────────────────────────────────────────────────────────────┘
   │
   ▼
┌──────────────────────────────────────────────────────────────┐
│  pkg/entities/*.go   → domain models / GORM table definitions  │
│  db/                 → DB interface + Postgres implementation  │
└──────────────────────────────────────────────────────────────┘
```

### Layer responsibilities

| Layer | Location | Responsibility |
| ----- | -------- | -------------- |
| **Entry point** | `cmd/unknown-api/main.go` | Loads config, builds the DB, starts the Fiber server. |
| **Config** | `config/config.go`, `config.yaml` | Loads YAML + env via Viper into a singleton `Config`. |
| **Server** | `server/` | `Server` interface + `fiberServer` that wires every repo → service → handler and registers routes. |
| **Routes** | `api/routes/` | One file per domain; maps URL paths to handler methods. |
| **Handlers** | `api/handlers/` | HTTP I/O: parse params/body, call the service, map errors to HTTP status codes. |
| **Presenters** | `api/presenter/` | Request and response structs (DTOs) — the API contract, kept separate from entities. |
| **Services** | `pkg/<domain>/service.go` | Business rules, validation, orchestration. No Fiber, no GORM. |
| **Repositories** | `pkg/<domain>/repository.go` | Persistence via GORM, behind an interface. |
| **Entities** | `pkg/entities/` | Domain models with GORM tags = the database schema. |
| **Database** | `db/` | `Database` interface + Postgres/GORM implementation and `AutoMigrate`. |

### Key design decisions

- **Dependency inversion / decoupling** — `db.Database` and `server.Server` are
  interfaces, so the concrete engine (Postgres, Fiber) can be swapped with
  minimal impact. Repositories hide GORM; services never import Fiber or GORM.
- **Wiring is explicit** — `server/fiberServer.go` is the composition root: it
  constructs every `repo → service → handler` chain and injects dependencies
  (e.g. `assetService` receives both the asset repo and the material repo;
  `materialCostService` receives vendor-material, vendor and currency repos).
- **Soft deletes over hard deletes** — most domains expose `Deactivate`
  (flip `IsActive = false`) instead of physical deletion, to preserve
  historical/referential integrity.
- **Append-only history** — `MaterialCost` records are never updated; a new row
  is inserted for each price, so cost history is fully auditable.
- **UUID primary keys** generated in the database via the `uuid-ossp`
  extension (`uuid_generate_v4()`), except the sync table which uses an
  auto-increment integer.
- **Singletons** — both config and the DB connection use `sync.Once` to ensure
  a single instance per process.
- **Authentication middleware** — `pkg/auth` provides `RequireAuth()` (valid
  HS256 JWT in `Authorization: Bearer <token>`, claims injected into
  `fiber.Locals`) and `RequireRole(roles...)`. In `server/fiberServer.go` the
  auth routes (`/auth/login` public, `/auth/me`) are registered first, then
  `api.Use(auth.RequireAuth())` protects every `/api/v1/*` route registered
  after it. `/locations` and `/users` additionally require the `ADMIN` role;
  `/public/*` routes and the health check stay unauthenticated. The role
  travels inside the token, so a role change only applies on the next login.
  On first boot, `SeedInitialAdmin` creates the first `ADMIN` account from the
  `SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` env vars (idempotent no-op once
  an active admin exists).

---

## 3. Technology stack

| Concern | Technology |
| ------- | ---------- |
| Language | **Go 1.22.1** |
| HTTP framework | **Fiber v2** (`github.com/gofiber/fiber/v2`) with CORS middleware |
| ORM | **GORM v1.25** (`gorm.io/gorm`) |
| Database | **PostgreSQL** (`gorm.io/driver/postgres`, pgx driver) |
| Config | **Viper** (`github.com/spf13/viper`) — YAML file + env overrides |
| UUIDs | `github.com/google/uuid` + Postgres `uuid-ossp` extension |
| Auth | **JWT HS256** (`github.com/golang-jwt/jwt/v5`) + **bcrypt** (`golang.org/x/crypto/bcrypt`) |
| Errors | `github.com/pkg/errors` |
| Module path | `github.com/unknowncode44/unknown-go-server` |

- **API base path:** `/api/v1` (plus a `/public` group and a health check at
  `GET /api/v1/health`).
- **Default port:** `4045` (configurable in `config.yaml`).
- Schema is created/updated automatically at startup with GORM `AutoMigrate`
  (no separate migration tool).

---

## 4. Data model — tables & relationships

All primary keys are UUIDs unless noted. `IsActive`, `CreatedAt`, `UpdatedAt`
are common bookkeeping fields. Monetary/quantity values use `numeric(15,4)`.

### 4.1 Tables

**Material** — the catalog item (master record).
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| Name | varchar(255) | not null |
| Sector | varchar(100) | business sector (electrical, RF, …) |
| UnitOfMeasure | varchar(50) | meter, unit, kg, … |
| ERPCode | string | indexed — external ERP code |
| InternalCode | string | unique — used to generate asset serials |
| Code | varchar(50) | unique, not null |
| IsActive | bool | soft-delete flag |

**Vendor** — a supplier.
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| Name | varchar(255) | not null |
| Code | varchar(100) | optional |
| TaxID | varchar(50) | optional — used to link purchase-order syncs |
| IsActive | bool | |

**VendorMaterial** — join table: which vendors supply which materials, with the vendor's own code/SKU.
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| VendorID | uuid | → Vendor, not null |
| MaterialID | uuid | → Material, not null |
| VendorCode | string | vendor's SKU/part number |
| IsActive | bool | |

**Currency** — reference/master data (ISO 4217).
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| Code | varchar(3) | unique — ARS, USD, EUR |
| Name | varchar(100) | not null |
| IsActive | bool | |

**MaterialCost** — append-only historical cost of a material from a vendor.
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| MaterialID | uuid | → Material, indexed |
| VendorID | uuid | → Vendor, indexed |
| CurrencyID | uuid | → Currency, indexed |
| Cost | numeric(15,4) | not null |
| CostDate | timestamp | date the cost becomes valid, indexed |

**Location** — hierarchical place (internal warehouse or client site).
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| Name | varchar(128) | not null |
| Type | varchar(32) | `INTERNAL` or `CLIENT` |
| ParentLocationID | uuid? | self-reference (tree) |
| IsActive | bool | |

**Asset** — a physical, individually tracked instance of a material.
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| SerialVisible | varchar(64) | unique — public QR serial, generated from material InternalCode |
| MaterialID | uuid | → Material, not null |
| ManufacturerSerial | varchar(128)? | optional |
| Status | varchar(32) | `IN_STOCK` / `INSTALLED` / `IN_REPAIR` / `RETIRED` |
| ParentAssetID | uuid? | self-reference for kits/hierarchies |
| CurrentLocationID | uuid? | → Location (derived state, updated by movements) |
| IsActive | bool | |

**AssetSerialCounter** — last correlative number per material, used to generate `SerialVisible` safely under locking.
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| MaterialID | uuid | unique |
| Last | int | last used number |

**AssetMovement** — audit-trail event for an asset.
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| AssetID | uuid | → Asset, not null |
| Type | varchar(32) | `INBOUND`, `TRANSFER`, `INSTALL`, `UNINSTALL`, `REPAIR`, `RETURN`, `DECOMMISSION`, `SCRAP` |
| FromLocationID | uuid? | → Location |
| ToLocationID | uuid? | → Location |
| MovementDate | timestamp | not null |
| Notes | text? | optional |

**MaterialInventory** — bulk (non-serialized) stock header for a material.
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| MaterialID | uuid | → Material, unique (one inventory per material) |
| InitialQuantity | numeric(15,4) | opening balance |
| StartDate | timestamp | not null |
| IsActive | bool | |

**DeliveryRecord** — an inbound/outbound quantity movement against a MaterialInventory.
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| MaterialInventoryID | uuid | → MaterialInventory, indexed |
| Type | varchar(16) | `INBOUND` or `OUTBOUND` |
| Quantity | numeric(15,4) | not null |
| DeliveryDate | timestamp | indexed |
| Notes | text? | optional |

**User** — a system account with access to the API (auth + roles).
| Field | Type | Notes |
|-------|------|-------|
| ID | uuid | PK |
| Name | varchar(255) | not null |
| Email | varchar(255) | unique, not null — login identifier |
| PasswordHash | varchar(255) | bcrypt hash; never exposed by the API |
| Role | varchar(20) | `ADMIN` / `USER` / `OPERATIVE_USER` (default `USER`) |
| IsActive | bool | soft-delete flag; the last active `ADMIN` cannot be deactivated |

**PurchaseOrderSync** — raw purchase-order rows ingested from an external source (Excel/VBA), later linked to a material/vendor.
| Field | Type | Notes |
|-------|------|-------|
| ID | uint | **integer** PK (auto-increment) |
| CCO, NotaPedido, FechaNP, OCBejerman, Articulo, Descripcion, Proveedor, Moneda | string/text | raw imported fields |
| FechaOC | date | |
| Cantidad, ImporteUnitario | decimal(19,4) | |
| SincronizadoEn | timestamp | auto-set on insert |
| MaterialID | uuid? | → Material (nullable until linked), indexed |
| VendorID | uuid? | → Vendor (nullable until linked), indexed |

### 4.2 How the relationships work

Relationships are modeled via **foreign-key UUID columns**, and in some cases
GORM association structs are declared for eager loading. The conceptual graph:

```
                       ┌───────────┐
                       │ Currency  │
                       └─────┬─────┘
                             │ (CurrencyID)
        ┌──────────┐   ┌─────┴────────┐   ┌──────────┐
        │ Material │◄──┤ MaterialCost ├──►│  Vendor  │
        └────┬─────┘   └──────────────┘   └────┬─────┘
             │  ▲                               │
             │  │  (MaterialID / VendorID)      │
             │  └──────────┬────────────────────┘
             │      ┌───────┴────────┐
             │      │ VendorMaterial │   (which vendor supplies which material)
             │      └────────────────┘
             │
             ├──────────────► ┌────────────────────┐  (1 material → 1 inventory)
             │                │ MaterialInventory  │
             │                └─────────┬──────────┘
             │                          │ (MaterialInventoryID)
             │                ┌─────────┴──────────┐
             │                │  DeliveryRecord    │  INBOUND / OUTBOUND qty
             │                └────────────────────┘
             │
             ├──────────────► ┌────────────────────┐  (correlative generator)
             │                │ AssetSerialCounter │
             │                └────────────────────┘
             │
             ▼  (MaterialID)
        ┌──────────┐
        │  Asset   │──self──► ParentAssetID (kits / sub-assets)
        └────┬─────┘
             │ (CurrentLocationID, derived)
             │           ┌──────────┐
             ├──────────►│ Location │──self──► ParentLocationID (tree)
             │           └────┬─────┘
             │                │ (From/ToLocationID)
             ▼  (AssetID)     │
        ┌────────────────┐    │
        │ AssetMovement  ├────┘   audit trail; updates Asset.Status + CurrentLocation
        └────────────────┘

        ┌────────────────────┐
        │ PurchaseOrderSync  │── MaterialID? / VendorID?  (linked after import)
        └────────────────────┘
```

Relationship summary:

- **Material 1—N VendorMaterial N—1 Vendor** — many-to-many between materials
  and vendors, resolved through the `VendorMaterial` join table (which also
  stores the vendor's own code).
- **MaterialCost N—1 Material / Vendor / Currency** — each cost row points at
  one material, one vendor and one currency; append-only history.
- **Asset N—1 Material** — every asset is an instance of exactly one material.
- **Asset self-reference** — `ParentAssetID` builds kits/hierarchies of assets.
- **Asset N—1 Location** (`CurrentLocationID`) — derived "current location",
  kept in sync by asset movements.
- **AssetMovement N—1 Asset**, and N—1 to `Location` twice (from / to). The
  first movement of an asset must be `INBOUND`; creating a movement
  transactionally updates the asset's `Status` and `CurrentLocationID`.
- **Location self-reference** — `ParentLocationID` builds a location tree.
- **MaterialInventory 1—1 Material** (unique `MaterialID`); **DeliveryRecord
  N—1 MaterialInventory** records each quantity in/out.
- **PurchaseOrderSync** holds optional `MaterialID` / `VendorID` that start
  null and are filled in once a synced row is linked to catalog data.

---

## 5. Key business flows

### Asset lifecycle
1. Create a **Material** (must have a non-empty `InternalCode`).
2. Create an **Asset** for that material — the service loads the material,
   verifies it is active, and generates a unique `SerialVisible` via the
   `AssetSerialCounter` (under locking). Initial status defaults to `IN_STOCK`.
3. Record **AssetMovements** (`INBOUND` first, then `TRANSFER`, `INSTALL`, …).
   Each movement is committed atomically and updates the asset's status and
   current location.
4. `DECOMMISSION` / `SCRAP` finalize the asset.

See `README_ASSET_CYCLE.md` for full endpoint payloads.

### Material costing
A `MaterialCost` row ties a material + vendor + currency + cost + valid-from
date. New prices are appended (never updated), giving a complete cost history
for comparison and "last updated cost" queries.

### Bulk inventory
For materials tracked by quantity rather than serial, a single
`MaterialInventory` holds the opening balance and `DeliveryRecord` rows add/
remove quantity (`INBOUND`/`OUTBOUND`).

### Purchase order sync
An external Excel/VBA macro POSTs an array of purchase-order rows to a `/public`
endpoint; they are stored in `PurchaseOrderSync` and can later be linked to a
material (`MaterialID`) and vendor (`VendorID` / vendor TAX ID).

---

## 6. Configuration & running

`config.yaml` (overridable by environment variables via Viper):

```yaml
server:
  port: 4045
db:
  host: <postgres-host>
  port: <port>
  user: <user>
  password: <password>
  dbname: quinar
  timezone: america/buenos_aires
auth:
  jwt_secret: <long-random-secret>   # in production set via env var AUTH_JWT_SECRET
  jwt_expiry_hours: 12
```

Auth-related environment variables:

- `AUTH_JWT_SECRET` — overrides `auth.jwt_secret` (Viper maps `.` → `_`).
  **Never commit the real secret.**
- `SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` — export at least once before the
  first start on a fresh database so the initial `ADMIN` gets created (the
  seed is an idempotent no-op once an active admin exists).

Run from the entry point package:

```bash
go run ./cmd/unknown-api
```

On startup the server connects to Postgres, installs the `uuid-ossp`
extension, runs `AutoMigrate` for all entities, wires up every domain, and
listens on the configured port.

---

## 7. Repository layout

```
cmd/unknown-api/      → main entry point
config/               → Viper config loader (+ config.yaml at repo root)
db/                   → Database interface + Postgres/GORM implementation
server/               → Server interface + Fiber composition root
api/
  routes/             → route registration per domain
  handlers/           → HTTP handlers per domain
  presenter/          → request/response DTOs
pkg/
  entities/           → GORM models (the schema)
  <domain>/           → repository.go + service.go per domain
                        (material, vendor, vendor_material, currency,
                         material_cost, asset, asset_movement, location,
                         material_inventory, delivery_record,
                         sync_purchase_order, user)
  auth/               → JWT generation/validation + RequireAuth/RequireRole
                        middlewares (no repository/entity of its own)
docs/                 → this overview
README.MD             → Material module notes
README_ASSET_CYCLE.md → asset movement cycle guide
```
