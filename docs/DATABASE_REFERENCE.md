# Unknown Go Server — Database Reference

> **Audience:** an AI agent (or developer) that needs a precise, complete model
> of this project's database. This document is the source of truth for tables,
> columns, types, constraints, indexes, relationships, enumerations, and the
> conventions used. It is derived directly from the GORM entity definitions in
> `pkg/entities/` and the migration in `db/postgres.go`.
>
> When reasoning about queries or schema, prefer the facts in this file. Where a
> behavior is enforced in application code rather than the database, it is
> explicitly marked **(app-level)**.

---

## 1. Engine & conventions

- **DBMS:** PostgreSQL.
- **ORM:** GORM v2 (`gorm.io/gorm`, `gorm.io/driver/postgres`).
- **Schema management:** GORM `AutoMigrate`, run at startup in
  `db/postgres.go` → `NewPostgresDatabase`. There is **no** standalone
  migration tool or versioned migration files. The schema is whatever
  `AutoMigrate` produces from the entity structs.
- **Extension:** `uuid-ossp` is enabled at startup
  (`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`) so the DB can generate UUIDs
  via `uuid_generate_v4()`.
- **Primary keys:** `uuid` for every table **except** `purchase_order_syncs`,
  which uses an auto-increment integer (`bigserial`).
- **Table naming:** GORM's default — the struct name is pluralized and
  snake_cased (e.g. `MaterialCost` → `material_costs`,
  `MaterialInventory` → `material_inventories`).
- **Column naming:** snake_case of the Go field name unless an explicit
  `column:` tag overrides it.
- **Timestamps:** `created_at` / `updated_at` are managed automatically by GORM
  (`CreatedAt` / `UpdatedAt` fields). Tables that only have `CreatedAt` are
  effectively append-only or event-style.
- **Money / quantities:** `numeric(15,4)` (or `decimal(19,4)` for the sync
  table). Always treat as exact decimals, not floats, when it matters.

### 1.1 Two kinds of relationship — IMPORTANT

This codebase models foreign keys in **two different ways**, and the
distinction matters for whether a real DB constraint exists:

1. **Association relationships** — the struct has both a `...ID uuid` column
   **and** a GORM association field (e.g. `Material Material
   ` + `gorm:"foreignKey:MaterialID"`). For these, `AutoMigrate` creates a
   **real foreign-key constraint** in PostgreSQL, and the relation can be
   eager-loaded with `Preload`.

2. **Bare-ID relationships** — the struct has only a `...ID uuid` column with
   no association field (e.g. `MaterialCost.VendorID`). These are
   **application-level references only**: there is *no* DB foreign-key
   constraint, just an indexed UUID column. Referential integrity is the
   responsibility of the service layer **(app-level)**.

Each relationship in §4 is tagged as **[FK constraint]** or
**[app-level ref]** accordingly.

---

## 2. Table catalog (quick map)

| # | Table | Struct | PK | Purpose |
|---|-------|--------|----|---------|
| 1 | `materials` | `Material` | uuid | Catalog item (master record) |
| 2 | `vendors` | `Vendor` | uuid | Supplier |
| 3 | `vendor_materials` | `VendorMaterial` | uuid | M:N link vendor ↔ material + vendor's code |
| 4 | `currencies` | `Currency` | uuid | Reference/master data (ISO 4217) |
| 5 | `material_costs` | `MaterialCost` | uuid | Append-only price history |
| 6 | `locations` | `Location` | uuid | Hierarchical place (warehouse / client site) |
| 7 | `assets` | `Asset` | uuid | Serialized physical instance of a material |
| 8 | `asset_serial_counters` | `AssetSerialCounter` | uuid | Per-material correlative for serial generation |
| 9 | `asset_movements` | `AssetMovement` | uuid | Audit-trail event for an asset |
| 10 | `material_inventories` | `MaterialInventory` | uuid | Bulk (non-serialized) stock header per material |
| 11 | `delivery_records` | `DeliveryRecord` | uuid | Inbound/outbound quantity vs. an inventory |
| 12 | `purchase_order_syncs` | `PurchaseOrderSync` | **bigint** | Imported purchase-order rows (Excel/VBA) |
| 13 | `users` | `User` | uuid | System account (auth + roles) |

> All 13 entities are registered in the `AutoMigrate(...)` call
> in `db/postgres.go` and are auto-created at startup.
> (`needs` / `need_items` / `need_counters` are also migrated; see the Need
> domain entities in `pkg/entities/`.)

---

## 3. Tables (detailed)

Column legend: **PK** = primary key, **FK** = participates in a relationship
(see §4), **UQ** = unique index, **IX** = non-unique index, **NN** = NOT NULL,
**?** = nullable (Go pointer).

### 3.1 `materials`
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `name` | Name | varchar(255) | NN | |
| `sector` | Sector | varchar(100) | NN | business sector (electrical, RF…) |
| `material_group` | Group | varchar(10) | NN, default `'BDC'` | business group: `BDC` (Bienes de Cambio) or `BDU` (Bienes de Uso); explicit column name because `group` is a reserved word in SQL |
| `unit_of_measure` | UnitOfMeasure | varchar(50) | NN | meter, unit, kg… |
| `erp_code` | ERPCode | text | IX | external ERP code |
| `internal_code` | InternalCode | text | UQ | used to generate asset serials |
| `code` | Code | varchar(50) | UQ, NN | business code |
| `is_active` | IsActive | boolean | NN, default `true` | soft-delete flag |
| `created_at` | CreatedAt | timestamptz | | |
| `updated_at` | UpdatedAt | timestamptz | | |

### 3.2 `vendors`
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `name` | Name | varchar(255) | NN | |
| `code` | Code | varchar(100) | | optional |
| `tax_id` | TaxID | varchar(50) | | optional; used to link purchase-order syncs |
| `is_active` | IsActive | boolean | NN, default `true` | |
| `created_at` | CreatedAt | timestamptz | | |
| `updated_at` | UpdatedAt | timestamptz | | |

### 3.3 `vendor_materials`
Join table: which vendor supplies which material, with the vendor's own code.
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `vendor_id` | VendorID | uuid | NN, FK → `vendors` **[app-level ref]** | |
| `material_id` | MaterialID | uuid | NN, FK → `materials` **[app-level ref]** | |
| `vendor_code` | VendorCode | text | NN | vendor SKU / part number |
| `is_active` | IsActive | boolean | default `true` | |
| `created_at` | CreatedAt | timestamptz | | |
| `updated_at` | UpdatedAt | timestamptz | | |

### 3.4 `currencies`
Master/reference data.
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `code` | Code | varchar(3) | NN, UQ | ISO 4217 (ARS, USD, EUR) |
| `name` | Name | varchar(100) | NN | |
| `is_active` | IsActive | boolean | default `true` | |
| `created_at` | CreatedAt | timestamptz | | |
| `updated_at` | UpdatedAt | timestamptz | | |

### 3.5 `material_costs`
**Append-only** price history (no `updated_at`). One row per (material, vendor,
currency, date) price point.
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `material_id` | MaterialID | uuid | NN, IX, FK → `materials` **[app-level ref]** | |
| `vendor_id` | VendorID | uuid | NN, IX, FK → `vendors` **[app-level ref]** | |
| `currency_id` | CurrencyID | uuid | NN, IX, FK → `currencies` **[app-level ref]** | |
| `cost` | Cost | numeric(15,4) | NN | |
| `cost_date` | CostDate | timestamptz | NN, IX | date from which this cost is valid |
| `created_at` | CreatedAt | timestamptz | | |

### 3.6 `locations`
Self-referencing tree of places.
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `name` | Name | varchar(128) | NN | |
| `type` | Type | varchar(32) | NN | enum `LocationType` (§5) |
| `parent_location_id` | ParentLocationID | uuid? | FK → `locations` (self) **[FK constraint]** | null = root |
| `is_active` | IsActive | boolean | default `true` | |
| `created_at` | CreatedAt | timestamptz | | |
| `updated_at` | UpdatedAt | timestamptz | | |

Association fields: `ParentLocation *Location`, `ChildLocations []Location`
(both keyed on `parent_location_id`).

### 3.7 `assets`
A serialized physical instance of a material.
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `serial_visible` | SerialVisible | varchar(64) | NN, UQ | public QR serial, generated from material `internal_code` |
| `material_id` | MaterialID | uuid | NN, FK → `materials` **[FK constraint]** | |
| `manufacturer_serial` | ManufacturerSerial | varchar(128)? | | optional |
| `status` | Status | varchar(32) | NN | enum `AssetStatus` (§5); defaults to `IN_STOCK` at creation **(app-level)** |
| `parent_asset_id` | ParentAssetID | uuid? | FK → `assets` (self) **[FK constraint]** | kits / hierarchy |
| `current_location_id` | CurrentLocationID | uuid? | FK → `locations` **[FK constraint]** | derived state, maintained by movements |
| `is_active` | IsActive | boolean | default `true` | |
| `created_at` | CreatedAt | timestamptz | | |
| `updated_at` | UpdatedAt | timestamptz | | |

Association fields: `Material`, `ParentAsset *Asset`, `ChildAssets []Asset`,
`CurrentLocation *Location`.

### 3.8 `asset_serial_counters`
Generator support: last correlative used per material (locked during serial
generation to avoid collisions).
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `material_id` | MaterialID | uuid | NN, UQ, FK → `materials` **[app-level ref]** | one counter per material |
| `last` | Last | integer | NN, default `0` | last used number |
| `created_at` | CreatedAt | timestamptz | | |
| `updated_at` | UpdatedAt | timestamptz | | |

### 3.9 `asset_movements`
Append-only audit-trail event (no `updated_at`). Creating one transactionally
updates the parent asset's `status` and `current_location_id` **(app-level)**.
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `asset_id` | AssetID | uuid | NN, FK → `assets` **[FK constraint]** | |
| `type` | Type | varchar(32) | NN | enum `AssetMovementType` (§5) |
| `from_location_id` | FromLocationID | uuid? | FK → `locations` **[FK constraint]** | |
| `to_location_id` | ToLocationID | uuid? | FK → `locations` **[FK constraint]** | |
| `movement_date` | MovementDate | timestamptz | NN | effective date of the movement |
| `notes` | Notes | text? | | |
| `created_at` | CreatedAt | timestamptz | | |

Association fields: `Asset`, `FromLocation *Location`, `ToLocation *Location`.

### 3.10 `material_inventories`
Bulk stock header — one per material (unique `material_id`).
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `material_id` | MaterialID | uuid | NN, UQ, FK → `materials` **[FK constraint]** | one inventory per material |
| `initial_quantity` | InitialQuantity | numeric(15,4) | NN, default `0` | opening balance |
| `start_date` | StartDate | timestamptz | NN | |
| `is_active` | IsActive | boolean | NN, default `true` | |
| `created_at` | CreatedAt | timestamptz | | |
| `updated_at` | UpdatedAt | timestamptz | | |

Association field: `Material`.

### 3.11 `delivery_records`
Quantity in/out against a `material_inventories` row.
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `material_inventory_id` | MaterialInventoryID | uuid | NN, IX, FK → `material_inventories` **[FK constraint]** | |
| `type` | Type | varchar(16) | NN | enum `DeliveryType` (§5) |
| `quantity` | Quantity | numeric(15,4) | NN | |
| `delivery_date` | DeliveryDate | timestamptz | NN, IX | |
| `notes` | Notes | text? | | |
| `created_at` | CreatedAt | timestamptz | | |

Association field: `MaterialInventory`.

> **Derived total (not stored):** current stock for an inventory is computed in
> the application as
> `initial_quantity + Σ(INBOUND.quantity) − Σ(OUTBOUND.quantity)`
> over its delivery records (see the `/material-inventories/:id/total`
> endpoint). It is not persisted.

### 3.12 `purchase_order_syncs`
Raw purchase-order rows imported from an external Excel/VBA macro, later linked
to catalog data. **Integer PK.**
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | bigint | PK, auto-increment | |
| `cco` | CCO | varchar(50) | | |
| `nota_pedido` | NotaPedido | varchar(100) | | |
| `fecha_np` | FechaNP | varchar(20) | | stored as string |
| `fecha_oc` | FechaOC | date | | |
| `oc_bejerman` | OCBejerman | varchar(100) | | |
| `articulo` | Articulo | varchar(100) | | |
| `descripcion` | Descripcion | text | | |
| `cantidad` | Cantidad | decimal(19,4) | | |
| `proveedor` | Proveedor | varchar(255) | | free-text supplier name |
| `moneda` | Moneda | varchar(10) | | free-text currency |
| `importe_unitario` | ImporteUnitario | decimal(19,4) | | |
| `sincronizado_en` | SincronizadoEn | timestamptz | autoCreateTime | set on insert |
| `material_id` | MaterialID | uuid? | IX, FK → `materials` **[app-level ref]** | null until linked |
| `vendor_id` | VendorID | uuid? | IX, FK → `vendors` **[app-level ref]** | null until linked |

### 3.13 `users`
System accounts with API access (email/password login, JWT, roles). No
relationships to other tables.
| Column | Go field | Type | Constraints | Notes |
|--------|----------|------|-------------|-------|
| `id` | ID | uuid | PK, default `uuid_generate_v4()` | |
| `name` | Name | varchar(255) | NN | |
| `email` | Email | varchar(255) | UQ, NN | login identifier; not editable via API |
| `password_hash` | PasswordHash | varchar(255) | NN | bcrypt hash; never returned by the API |
| `role` | Role | varchar(20) | NN, default `'USER'` | `ADMIN` / `USER` / `OPERATIVE_USER` (app-level enum) |
| `is_active` | IsActive | boolean | NN, default `true` | soft-delete flag |
| `created_at` | CreatedAt | timestamptz | | |
| `updated_at` | UpdatedAt | timestamptz | | |

---

## 4. Relationships

Format: `child.column → parent.table` — cardinality — kind.

| Relationship | From (FK column) | To | Cardinality | Kind |
|--------------|------------------|----|-------------|------|
| Vendor supplies materials | `vendor_materials.vendor_id` | `vendors.id` | N:1 | app-level ref |
| Material offered by vendors | `vendor_materials.material_id` | `materials.id` | N:1 | app-level ref |
| Cost → material | `material_costs.material_id` | `materials.id` | N:1 | app-level ref |
| Cost → vendor | `material_costs.vendor_id` | `vendors.id` | N:1 | app-level ref |
| Cost → currency | `material_costs.currency_id` | `currencies.id` | N:1 | app-level ref |
| Location tree | `locations.parent_location_id` | `locations.id` | N:1 (self) | FK constraint |
| Asset → material | `assets.material_id` | `materials.id` | N:1 | FK constraint |
| Asset kit/hierarchy | `assets.parent_asset_id` | `assets.id` | N:1 (self) | FK constraint |
| Asset current location | `assets.current_location_id` | `locations.id` | N:1 | FK constraint |
| Serial counter → material | `asset_serial_counters.material_id` | `materials.id` | 1:1 (unique) | app-level ref |
| Movement → asset | `asset_movements.asset_id` | `assets.id` | N:1 | FK constraint |
| Movement origin | `asset_movements.from_location_id` | `locations.id` | N:1 | FK constraint |
| Movement destination | `asset_movements.to_location_id` | `locations.id` | N:1 | FK constraint |
| Inventory → material | `material_inventories.material_id` | `materials.id` | 1:1 (unique) | FK constraint |
| Delivery → inventory | `delivery_records.material_inventory_id` | `material_inventories.id` | N:1 | FK constraint |
| Sync row → material | `purchase_order_syncs.material_id` | `materials.id` | N:1 (nullable) | app-level ref |
| Sync row → vendor | `purchase_order_syncs.vendor_id` | `vendors.id` | N:1 (nullable) | app-level ref |

**Derived many-to-many:** `materials` ↔ `vendors` through `vendor_materials`
(plus `material_costs` also associates the two with a currency and a price).

### 4.1 Entity-relationship diagram (text)

```
                          ┌────────────┐
                          │ currencies │
                          └─────┬──────┘
                                │ currency_id
   ┌───────────┐   ┌────────────┴─────┐   ┌──────────┐
   │ materials │◄──┤  material_costs  ├──►│ vendors  │
   └─────┬─────┘   └──────────────────┘   └────┬─────┘
         │  ▲ material_id        vendor_id ▲    │
         │  │      ┌──────────────────┐    │    │
         │  └──────┤ vendor_materials ├────┘    │   (M:N materials↔vendors)
         │         └──────────────────┘         │
         │                                       │
         │ material_id (unique)                  │
         ├───────────────► ┌───────────────────────┐
         │                 │ material_inventories  │ (1 per material)
         │                 └──────────┬────────────┘
         │                            │ material_inventory_id
         │                 ┌──────────┴────────────┐
         │                 │   delivery_records    │ INBOUND / OUTBOUND
         │                 └───────────────────────┘
         │
         │ material_id (unique)
         ├───────────────► ┌────────────────────────┐
         │                 │ asset_serial_counters  │ (serial generator)
         │                 └────────────────────────┘
         │
         │ material_id
         ▼
   ┌──────────┐   parent_asset_id (self)
   │  assets  │──────────────► assets
   └────┬─────┘
        │ current_location_id          ┌────────────┐
        ├─────────────────────────────►│ locations  │── parent_location_id (self)
        │                              └─────┬──────┘
        │ asset_id                            │ from_location_id / to_location_id
        ▼                                     │
   ┌────────────────┐ ◄───────────────────────┘
   │ asset_movements│  (updates asset.status + current_location_id)
   └────────────────┘

   ┌───────────────────────┐ material_id?  vendor_id?
   │ purchase_order_syncs  │──────────────► materials / vendors (nullable, set on link)
   └───────────────────────┘
```

---

## 5. Enumerations

These are stored as strings (`varchar`) in the listed columns. There are **no**
DB-level CHECK constraints; valid values are enforced **(app-level)**.

**`MaterialGroup`** — `materials.material_group`
| Value | Meaning |
|-------|---------|
| `BDC` | Bienes de Cambio — materials sold to / installed at clients (bulk or serialized) |
| `BDU` | Bienes de Uso — Quinar's own infrastructure assets (always serialized) |

**`AssetStatus`** — `assets.status`
| Value | Meaning |
|-------|---------|
| `IN_STOCK` | in stock (default on creation) |
| `INSTALLED` | installed |
| `IN_REPAIR` | in repair |
| `RETIRED` | retired |

**`AssetMovementType`** — `asset_movements.type`
| Value | Meaning |
|-------|---------|
| `INBOUND` | entry into stock (**must be the first movement** of an asset, app-level) |
| `TRANSFER` | internal transfer |
| `INSTALL` | installation |
| `UNINSTALL` | removal |
| `REPAIR` | sent to repair |
| `RETURN` | returned from repair |
| `DECOMMISSION` | permanent retirement |
| `SCRAP` | permanent retirement due to breakage/obsolescence |
| `SOLD` | sale of a serialized BDC asset to an end client (same asset effect as `DECOMMISSION`/`SCRAP`: retired + inactive) |

**`LocationType`** — `locations.type`
| Value | Meaning |
|-------|---------|
| `INTERNAL` | internal location |
| `CLIENT` | client site |

**`DeliveryType`** — `delivery_records.type`
| Value | Meaning |
|-------|---------|
| `INBOUND` | quantity in |
| `OUTBOUND` | quantity out |

**`UserRole`** — `users.role`
| Value | Meaning |
|-------|---------|
| `ADMIN` | manages locations (warehouses/sectors/shelves) and users |
| `USER` | manages materials / Bienes de Cambio |
| `OPERATIVE_USER` | manages assets / Bienes de Uso (BDU) |

---

## 6. Behavioral / integrity rules (enforced in application code)

These rules are **not** in the database schema; an agent writing raw SQL must
uphold them manually, and an agent reading data should expect them:

1. **Soft deletes.** Most "delete" operations set `is_active = false` rather
   than removing rows (`materials`, `vendors`, `vendor_materials`,
   `currencies`, `locations`, `assets`, `material_inventories`). Filter on
   `is_active = true` for active sets. `material_costs`, `asset_movements`,
   `delivery_records` have no `is_active` (event/append-only).
2. **Append-only history.** `material_costs` rows are never updated; a new row
   is inserted per price change. "Current cost" = latest `cost_date` ≤ a given
   date for a (material, vendor).
3. **Asset serial generation.** On asset creation, `serial_visible` is built
   from the material's `internal_code` using `asset_serial_counters` under a
   lock/transaction. The referenced material must be active and have a
   non-empty `internal_code`.
4. **Movement side effects.** Inserting an `asset_movements` row updates the
   asset's `status` and `current_location_id` atomically (same transaction).
   The first movement for an asset must be `INBOUND`.
5. **One inventory per material / one counter per material.** Enforced by the
   unique index on `material_id` in `material_inventories` and
   `asset_serial_counters`.
6. **Referential integrity for app-level refs** (`vendor_materials`,
   `material_costs`, `asset_serial_counters`, `purchase_order_syncs`): there is
   no DB FK, so the service layer must ensure referenced rows exist. Orphans
   are possible if code bypasses the services.
7. **Users.** Passwords are stored only as bcrypt hashes; email is unique and
   not editable via the API. The last active `ADMIN` cannot be deactivated
   (app-level guard). On startup, `SeedInitialAdmin` inserts the first `ADMIN`
   from `SEED_ADMIN_EMAIL` / `SEED_ADMIN_PASSWORD` env vars if no active admin
   exists.

---

## 7. Notes for an agent modifying the schema

- To add a table: define the struct in `pkg/entities/`, add GORM tags, **and**
  register it in the `AutoMigrate(...)` slice in `db/postgres.go` — otherwise
  the table is never created.
- To make a bare-ID reference a real DB FK, add a GORM association field
  (e.g. `Vendor Vendor `gorm:"foreignKey:VendorID"``) or an explicit
  `constraint:` tag, then re-run `AutoMigrate`.
- `AutoMigrate` is additive: it creates missing tables/columns/indexes and
  adds missing FK constraints, but it does **not** drop columns or alter
  existing column types destructively. Renames are treated as add + orphan.
- UUID defaults rely on the `uuid-ossp` extension; keep the
  `CREATE EXTENSION` call before `AutoMigrate`.
```
