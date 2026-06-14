# Unknown Go Server — API Reference

Complete reference for every HTTP endpoint exposed by `unknown-go-server`
(the Quinar backend). It documents methods, paths, path/query params, request
payloads, and response shapes for all routes.

For the domain model and architecture, see
[PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md).

---

## Conventions

### Base paths
- **Authenticated/standard API:** all routes below are mounted under
  **`/api/v1`** (e.g. `POST /api/v1/materials`).
- **Public group:** mounted under **`/public`** (QR asset lookup, Excel sync).
- **Server port:** `4045` by default (from `config.yaml`). Examples use
  `http://localhost:4045`.
- **CORS:** enabled for all origins.

### Response envelopes
The API uses two different envelope shapes — note the **key inconsistency**
between success and error responses:

**Success** (most endpoints) — wraps the payload under `data`, with an `ok`
boolean flag:
```json
{ "ok": true, "data": { /* resource or array */ } }
```

**Error** (all `respondError` responses) — uses `success` (not `ok`):
```json
{ "success": false, "error": "human-readable message" }
```

**Sync endpoints** (`/public/sync/...`) use their own envelope:
```json
{ "ok": true, "data": {}, "message": "...", "error": "..." }
```

### Common status codes
| Code | Meaning |
|------|---------|
| `200 OK` | Successful read/update |
| `201 Created` | Resource created |
| `204 No Content` | Successful deactivate/delete (empty body) |
| `400 Bad Request` | Invalid body, invalid UUID, missing/invalid field |
| `404 Not Found` | Resource not found |
| `500 Internal Server Error` | Unexpected/persistence error |

### Data formats
- **IDs:** UUID strings (except Purchase Order Sync, which uses an integer id).
- **Dates in requests:** accept either **RFC3339** (`2026-04-20T15:00:00Z`) or
  **`YYYY-MM-DD`** (`2026-04-20`), unless noted.
- **Dates in responses:** RFC3339 timestamps.
- **Money/quantity:** JSON numbers (stored as `numeric(15,4)`).

### Endpoint index
| Resource | Base | Methods |
|----------|------|---------|
| Health | `/api/v1/health` | GET |
| Materials | `/api/v1/materials` | POST, GET, PUT, DELETE |
| Vendors | `/api/v1/vendors` | POST, GET, PUT, DELETE |
| Vendor-Materials | `/api/v1/vendor-materials` | POST, GET, PUT, DELETE |
| Currencies | `/api/v1/currencies` | POST, GET, PATCH, POST(deactivate) |
| Material Costs | `/api/v1/material-costs` + `/materials/:id/...` | POST, GET |
| Locations | `/api/v1/locations` | POST, GET, PUT, DELETE |
| Assets | `/api/v1/assets` | POST, GET, PUT, DELETE |
| Asset Movements | `/api/v1/asset-movements` + `/assets/:id/movements` | POST, GET |
| Material Inventory | `/api/v1/material-inventories` + `/materials/:id/inventory` | POST, GET, PUT, DELETE |
| Delivery Records | `/api/v1/delivery-records` + `/material-inventories/:id/...` | POST, GET |
| Purchase Order Sync (public) | `/public/sync/purchase-orders` | POST, GET, PATCH |
| Public Asset (QR) | `/public/asset/:serial` | GET |

---

## Health

### `GET /api/v1/health`
Liveness check.

- **Response:** `200 OK`, plain text body `Working Cool!` (not JSON).

---

## Materials

Base: `/api/v1/materials`

### `POST /api/v1/materials`
Create a material.

**Request body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | **yes** | rejected if empty |
| `sector` | string | no | business sector |
| `unitOfMeasure` | string | no | meter, unit, kg, … |
| `erpCode` | string | no | external ERP code |
| `internalCode` | string | no | required later to create Assets |
| `code` | string | no | unique business code |

```json
{
  "name": "Steel Bolt M8",
  "sector": "Hardware",
  "unitOfMeasure": "pcs",
  "erpCode": "ERP-12345",
  "internalCode": "INT-0001",
  "code": "BOLT-M8-001"
}
```

**Response** `201 Created`
```json
{
  "ok": true,
  "data": {
    "id": "uuid",
    "name": "Steel Bolt M8",
    "sector": "Hardware",
    "erpCode": "ERP-12345",
    "internalCode": "INT-0001",
    "unitOfMeasure": "pcs",
    "code": "BOLT-M8-001",
    "isActive": true,
    "createdAt": "2026-06-14T12:00:00Z",
    "updatedAt": "2026-06-14T12:00:00Z"
  }
}
```

### `GET /api/v1/materials`
List all materials. **Response** `200` → `{ "ok": true, "data": [ MaterialResponse, ... ] }`.

### `GET /api/v1/materials/:id`
Get a material by UUID. `400` on invalid UUID, `404` if not found.

> ⚠️ **Routing caveat:** the route `GET /materials/:erp_code` (intended to
> fetch by ERP code) is registered **after** `GET /materials/:id`. In Fiber the
> first matching pattern wins, so `:id` always matches a single path segment
> first and the `:erp_code` handler is effectively unreachable. Treat
> "get by ERP code" as not currently usable via this path.

### `PUT /api/v1/materials/:id`
Update a material. At least one of `name`, `sector`, `unitOfMeasure`,
`erpCode` must be provided, otherwise `400 "No fields provided for update"`.
(Only those four fields are applied by the handler.)

```json
{ "name": "Steel Bolt M8 (zinc)", "sector": "Hardware", "unitOfMeasure": "pcs", "erpCode": "ERP-99999" }
```
**Response** `200` → updated `MaterialResponse`.

### `DELETE /api/v1/materials/:id`
Soft-delete (deactivate). **Response** `204 No Content`. `404` if not found.

---

## Vendors

Base: `/api/v1/vendors`

### `POST /api/v1/vendors`
**Request body**
| Field | Type | Required |
|-------|------|----------|
| `name` | string | yes (service-level) |
| `code` | string | no |
| `tax_id` | string | no |

```json
{ "name": "ACME Supplies", "code": "ACME", "tax_id": "TAX12345678" }
```
**Response** `201` → `VendorResponse`:
```json
{ "ok": true, "data": { "id": "uuid", "name": "ACME Supplies", "code": "ACME", "tax_id": "TAX12345678", "is_active": true } }
```

### `GET /api/v1/vendors`
List all vendors → `{ "ok": true, "data": [ VendorResponse, ... ] }`.

### `GET /api/v1/vendors/:id`
Get vendor by UUID. `400` invalid UUID, `404` not found.

### `PUT /api/v1/vendors/:id`
Update vendor. `name` is required in the body (`400 "name is required"` if
empty). `code` and `tax_id` optional. **Response** `200` → `VendorResponse`.

### `DELETE /api/v1/vendors/:id`
Deactivate vendor. `204 No Content`.

---

## Vendor-Materials

Links a vendor to a material with the vendor's own code. Base:
`/api/v1/vendor-materials`

### `POST /api/v1/vendor-materials`
**Request body** — all three required:
| Field | Type | Required |
|-------|------|----------|
| `vendor_id` | UUID string | yes |
| `material_id` | UUID string | yes |
| `vendor_code` | string | yes |

```json
{ "vendor_id": "uuid", "material_id": "uuid", "vendor_code": "ACME-SKU-001" }
```
**Response** `201` → `VendorMaterialResponse`:
```json
{
  "ok": true,
  "data": {
    "id": "uuid",
    "vendor_id": "uuid",
    "material_id": "uuid",
    "vendor_code": "ACME-SKU-001",
    "is_active": true,
    "created_at": "...",
    "updated_at": "..."
  }
}
```

### `GET /api/v1/vendor-materials`
List all → array of `VendorMaterialResponse`.

### `GET /api/v1/vendor-materials/:id`
Get by UUID. `404` if not found.

### `PUT /api/v1/vendor-materials/:id`
Update. Only `vendor_code` is updatable; required (`400 "No fields provided for
update"` if empty).
```json
{ "vendor_code": "ACME-SKU-002" }
```

### `DELETE /api/v1/vendor-materials/:id`
Deactivate. `204 No Content`.

---

## Currencies

Base: `/api/v1/currencies`. Note this resource uses **PATCH** for update and a
dedicated **POST `/deactivate`** sub-route instead of DELETE.

### `POST /api/v1/currencies`
**Request body** — both required:
| Field | Type | Notes |
|-------|------|-------|
| `code` | string | ISO 4217, uppercased & trimmed server-side (e.g. `USD`) |
| `name` | string | trimmed |

```json
{ "code": "usd", "name": "US Dollar" }
```
**Response** `201` → `CurrencyResponse`:
```json
{
  "ok": true,
  "data": {
    "id": "uuid",
    "code": "USD",
    "name": "US Dollar",
    "isActive": true,
    "createdAt": "...",
    "updatedAt": "..."
  }
}
```

### `GET /api/v1/currencies`
List all → array of `CurrencyResponse`.

### `GET /api/v1/currencies/:id`
Get by UUID. `400` invalid UUID, `404` not found.

### `PATCH /api/v1/currencies/:id`
Update. Only `name` is updatable; required (`400 "No fields provided for
update"` if empty).
```json
{ "name": "United States Dollar" }
```

### `POST /api/v1/currencies/:id/deactivate`
Logical delete. **Response** `204 No Content`. `404` if not found.

---

## Material Costs

Append-only cost history (material + vendor + currency + cost + valid-from
date). Base: `/api/v1/material-costs`, plus material-scoped query routes.

### `POST /api/v1/material-costs`
**Request body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `material_id` | UUID string | yes | |
| `vendor_id` | UUID string | yes | |
| `currency_id` | UUID string | yes | |
| `cost` | number | yes | |
| `cost_date` | string | yes | RFC3339 or `YYYY-MM-DD`; valid-from date |

```json
{
  "material_id": "uuid",
  "vendor_id": "uuid",
  "currency_id": "uuid",
  "cost": 12.55,
  "cost_date": "2026-04-01"
}
```
**Response** `201` → `MaterialCostResponse`:
```json
{
  "ok": true,
  "data": {
    "id": "uuid",
    "material_id": "uuid",
    "vendor_id": "uuid",
    "currency_id": "uuid",
    "cost": 12.55,
    "cost_date": "2026-04-01T00:00:00Z",
    "created_at": "..."
  }
}
```

### `GET /api/v1/material-costs`
List all cost records → array of `MaterialCostResponse`.

### `GET /api/v1/material-costs/:id`
Get a single cost record by UUID. `404` if not found.

### `GET /api/v1/materials/:id/costs`
All cost records for a given material → array of `MaterialCostResponse`.

### `GET /api/v1/materials/:materialId/vendors/:vendorId/current-cost`
The currently-valid cost for a material from a specific vendor.

- **Query params:** `at=YYYY-MM-DD` (optional) — evaluate the "current" cost as
  of this date instead of now.
- **Response** `200` → `CurrentCostData`:
```json
{
  "ok": true,
  "data": {
    "material_id": "uuid",
    "vendor_id": "uuid",
    "currency": "uuid-of-currency",
    "cost": 12.55,
    "cost_date": "2026-04-01T00:00:00Z"
  }
}
```
> Note: the `currency` field currently returns the **currency UUID string**,
> not the ISO code (a known handler limitation). `404` if no cost exists.

### `GET /api/v1/materials/:materialId/vendor-comparison`
Compare the current cost across all vendors for a material.

- **Query params:** `at=YYYY-MM-DD` (optional).
- **Response** `200` → `VendorComparisonResponse`:
```json
{
  "ok": true,
  "data": {
    "material_id": "uuid",
    "at": "2026-06-14",
    "vendors": [
      { "id": "vendor-uuid", "name": "ACME Supplies", "currency": "USD", "cost": 12.55, "cost_date": "2026-04-01T00:00:00Z" }
    ]
  }
}
```

---

## Locations

Hierarchical places (warehouses / client sites). Base: `/api/v1/locations`

### `POST /api/v1/locations`
**Request body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | yes | |
| `type` | string | yes | `INTERNAL` or `CLIENT` |
| `parent_location_id` | UUID string | no | parent in the location tree |

```json
{ "name": "Main Warehouse", "type": "INTERNAL", "parent_location_id": "uuid-optional" }
```
**Response** `201` → `LocationResponse`:
```json
{
  "ok": true,
  "data": {
    "id": "uuid",
    "name": "Main Warehouse",
    "type": "INTERNAL",
    "parent_location_id": "uuid-or-omitted",
    "is_active": true,
    "created_at": "...",
    "updated_at": "..."
  }
}
```

### `GET /api/v1/locations`
List all → array of `LocationResponse`.

### `GET /api/v1/locations/:id`
Get by UUID. `404` if not found.

### `PUT /api/v1/locations/:id`
Update. All fields optional; provided fields overwrite. `is_active` is a
nullable boolean — include it to toggle activation.
```json
{ "name": "Central Warehouse", "type": "INTERNAL", "parent_location_id": "uuid", "is_active": false }
```

### `DELETE /api/v1/locations/:id`
Deactivate. `204 No Content`.

---

## Assets

Serialized instances of a material. Base: `/api/v1/assets`

### `POST /api/v1/assets`
Creates an asset and auto-generates `serial_visible` from the material's
`internalCode`. The referenced material must be **active** and have a
non-empty `internalCode`. Initial `status` defaults to `IN_STOCK`.

**Request body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `material_id` | UUID string | yes | |
| `manufacturer_serial` | string | no | |
| `parent_asset_id` | UUID string | no | for kits/hierarchies |
| `current_location_id` | UUID string | no | |

```json
{
  "material_id": "uuid",
  "manufacturer_serial": "MFG-12345",
  "parent_asset_id": "uuid-optional",
  "current_location_id": "uuid-optional"
}
```
**Response** `201` → `AssetResponse`:
```json
{
  "ok": true,
  "data": {
    "id": "uuid",
    "serial_visible": "INT-0001-000001",
    "material_id": "uuid",
    "manufacturer_serial": "MFG-12345",
    "status": "IN_STOCK",
    "parent_asset_id": "uuid-or-omitted",
    "current_location_id": "uuid-or-omitted",
    "is_active": true,
    "created_at": "...",
    "updated_at": "..."
  }
}
```
Errors (`400`): `invalid material_id`, `material is inactive`,
`material internal code is empty`.

### `GET /api/v1/assets`
List all assets. **Response** `200` → `{ "ok": true, "data": [ ...raw asset entities... ] }`.
> Note: this list endpoint returns the raw entities (including nested
> associations) rather than the trimmed `AssetResponse` used by the
> single-asset endpoints.

### `GET /api/v1/assets/:id`
Get asset by UUID → `AssetResponse`. `404` if not found.

### `GET /api/v1/assets/serial/:serial`
Get asset by its visible serial → `AssetResponse`. `404` if not found.

### `PUT /api/v1/assets/:id`
Update. All fields optional; provided fields overwrite.
| Field | Type | Notes |
|-------|------|-------|
| `manufacturer_serial` | string | |
| `parent_asset_id` | UUID string | |
| `current_location_id` | UUID string | |
| `status` | string | `IN_STOCK` / `INSTALLED` / `IN_REPAIR` / `RETIRED` |
| `is_active` | bool (nullable) | |

```json
{ "status": "INSTALLED", "current_location_id": "uuid", "is_active": true }
```
**Response** `200` → `AssetResponse`.

### `DELETE /api/v1/assets/:id`
Deactivate. `204 No Content`.

---

## Asset Movements

Audit-trail events that also update the asset's status/location. Base:
`/api/v1/asset-movements` plus asset-scoped read routes.

### `POST /api/v1/asset-movements`
Records a movement. Committed transactionally with the asset update. The
**first** movement of an asset must be `INBOUND`.

**Request body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `asset_id` | UUID string | yes | |
| `type` | string | yes | `INBOUND`, `TRANSFER`, `INSTALL`, `UNINSTALL`, `REPAIR`, `RETURN`, `DECOMMISSION`, `SCRAP` |
| `from_location_id` | UUID string | no | |
| `to_location_id` | UUID string | no | |
| `movement_date` | string | yes | RFC3339 or `YYYY-MM-DD` |
| `notes` | string | no | |

```json
{
  "asset_id": "uuid",
  "type": "INBOUND",
  "to_location_id": "uuid",
  "movement_date": "2026-04-20T15:00:00Z",
  "notes": "Received from supplier"
}
```
**Response** `201` → `AssetMovementResponse`:
```json
{
  "ok": true,
  "data": {
    "id": "uuid",
    "asset_id": "uuid",
    "type": "INBOUND",
    "from_location_id": "uuid-or-omitted",
    "to_location_id": "uuid-or-omitted",
    "movement_date": "2026-04-20T15:00:00Z",
    "notes": "Received from supplier",
    "created_at": "..."
  }
}
```

### `GET /api/v1/assets/:asset_id/movements`
All movements for an asset → array of `AssetMovementResponse`.

### `GET /api/v1/assets/:asset_id/movements/last`
The most recent movement for an asset → single `AssetMovementResponse`.
`404` if the asset has no movements.

---

## Material Inventory

Bulk (non-serialized) stock per material. Base: `/api/v1/material-inventories`,
plus a material-scoped read route. One inventory per material.

### `POST /api/v1/material-inventories`
**Request body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `material_id` | UUID string | yes | |
| `initial_quantity` | number | no | opening balance (default 0) |
| `start_date` | string | yes | RFC3339 or `YYYY-MM-DD` |

```json
{ "material_id": "uuid", "initial_quantity": 100, "start_date": "2026-01-01" }
```
**Response** `201` → `MaterialInventoryResponse`:
```json
{
  "ok": true,
  "data": {
    "id": "uuid",
    "material_id": "uuid",
    "material_name": "Steel Bolt M8",
    "material_code": "BOLT-M8-001",
    "initial_quantity": 100,
    "start_date": "2026-01-01T00:00:00Z",
    "is_active": true,
    "created_at": "...",
    "updated_at": "..."
  }
}
```

### `GET /api/v1/material-inventories`
List all → array of `MaterialInventoryResponse`.

### `GET /api/v1/material-inventories/:id`
Get inventory by its UUID. `404` if not found.

### `GET /api/v1/material-inventories/:id/total`
Computed stock totals for an inventory (initial + inbound − outbound from its
delivery records).

**Response** `200` → `InventoryTotalResponse`:
```json
{
  "ok": true,
  "data": {
    "material_inventory_id": "uuid",
    "initial_quantity": 100,
    "total_inbound": 50,
    "total_outbound": 30,
    "current_total": 120
  }
}
```

### `GET /api/v1/material-inventories/:id/deliveries`
All delivery records for an inventory → array of `DeliveryRecordResponse`
(see Delivery Records below).

### `GET /api/v1/materials/:id/inventory`
Get the inventory for a given **material** UUID → `MaterialInventoryResponse`.
`404` if none.

### `PUT /api/v1/material-inventories/:id`
Update. Fields optional.
| Field | Type | Notes |
|-------|------|-------|
| `initial_quantity` | number (nullable) | |
| `start_date` | string | RFC3339 or `YYYY-MM-DD` |

```json
{ "initial_quantity": 120, "start_date": "2026-02-01" }
```

### `DELETE /api/v1/material-inventories/:id`
Deactivate. `204 No Content`.

---

## Delivery Records

Inbound/outbound quantity movements against an inventory. Base:
`/api/v1/delivery-records`.

### `POST /api/v1/delivery-records`
**Request body**
| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `material_inventory_id` | UUID string | yes | |
| `type` | string | yes | `INBOUND` or `OUTBOUND` |
| `quantity` | number | yes | |
| `delivery_date` | string | yes | RFC3339 or `YYYY-MM-DD` |
| `notes` | string | no | |

```json
{
  "material_inventory_id": "uuid",
  "type": "OUTBOUND",
  "quantity": 30,
  "delivery_date": "2026-03-15",
  "notes": "Sent to site B"
}
```
**Response** `201` → `DeliveryRecordResponse`:
```json
{
  "ok": true,
  "data": {
    "id": "uuid",
    "material_inventory_id": "uuid",
    "type": "OUTBOUND",
    "quantity": 30,
    "delivery_date": "2026-03-15T00:00:00Z",
    "notes": "Sent to site B",
    "created_at": "..."
  }
}
```

### `GET /api/v1/delivery-records`
List all delivery records → array of `DeliveryRecordResponse`.

> Per-inventory delivery listing is available via
> `GET /api/v1/material-inventories/:id/deliveries` (documented above).

---

## Purchase Order Sync (public)

Ingestion of purchase-order rows from an external Excel/VBA macro, plus linking
them to catalog data. Mounted under the **`/public`** group.

> ⚠️ **Status caveat:** these endpoints are fully implemented in the handler,
> service and routes, but on the current branch the handler is **not wired up**
> in `server/fiberServer.go` (`spoHandler` is referenced but never
> constructed), so the package does not compile as-is. The routes below reflect
> the intended contract once wiring is restored.

### `POST /public/sync/purchase-orders`
Accepts a **JSON array** of purchase-order rows; upserts them and returns the
count synced.

**Request body** — array of:
| Field | Type | Notes |
|-------|------|-------|
| `cco` | string | |
| `nota_pedido` | string | |
| `fecha_oc` | string | `YYYY-MM-DD` |
| `oc_bejerman` | string | |
| `articulo` | string | |
| `descripcion` | string | |
| `cantidad` | number | |
| `importe_unitario` | number | |
| `moneda` | string | |
| `proveedor` | string | |

```json
[
  {
    "cco": "CCO-1",
    "nota_pedido": "NP-100",
    "fecha_oc": "2026-05-01",
    "oc_bejerman": "OC-200",
    "articulo": "ART-1",
    "descripcion": "Cable 2.5mm",
    "cantidad": 500,
    "importe_unitario": 1.25,
    "moneda": "ARS",
    "proveedor": "ACME"
  }
]
```
**Response** `201`
```json
{ "ok": true, "message": "Sync completed", "data": { "records_synced": 1 } }
```
Errors: `400` (bad JSON), `500` (sync failure) → `{ "ok": false, "error": "..." }`.

### `GET /public/sync/purchase-orders`
Returns all synced purchase-order records.

**Response** `200`
```json
{ "ok": true, "data": [ /* PurchaseOrderSync records */ ] }
```

### `PATCH /public/sync/purchase-orders/:id/:material_id/:vendor_tax_id`
Links a synced record (integer `:id`) to a material UUID and a vendor (by TAX
id). All three are **path params**, no body.

- `:id` — integer record id
- `:material_id` — material UUID
- `:vendor_tax_id` — vendor TAX id string

Example:
```
PATCH /public/sync/purchase-orders/125/550e8400-e29b-41d4-a716-446655440000/TAX12345678
```
**Response** `200`
```json
{ "ok": true, "data": { /* updated record */ } }
```
Errors: `400` (invalid material id / missing vendor tax id), `500` (parse/update failure).

---

## Public Asset (QR)

### `GET /public/asset/:serial`
Unauthenticated, limited view of an asset for QR-code scans (looked up by its
visible serial). Best-effort includes current location and last movement.

**Response** `200` → `PublicAssetResponse` (returned directly, **not** wrapped
in the `ok/data` envelope):
```json
{
  "serial_visible": "INT-0001-000001",
  "material": { "name": "Steel Bolt M8", "code": "BOLT-M8-001" },
  "status": "INSTALLED",
  "current_location": { "id": "uuid", "name": "Site B" },
  "last_movement": { "type": "INSTALL", "date": "2026-04-21T10:00:00Z" }
}
```
`current_location` and `last_movement` are omitted if unavailable. `404` if the
serial is not found.

---

## Quick curl examples

```bash
# Create a material
curl -X POST http://localhost:4045/api/v1/materials \
  -H "Content-Type: application/json" \
  -d '{"name":"Steel Bolt M8","sector":"Hardware","unitOfMeasure":"pcs","code":"BOLT-M8-001","internalCode":"INT-0001"}'

# Create an asset from that material
curl -X POST http://localhost:4045/api/v1/assets \
  -H "Content-Type: application/json" \
  -d '{"material_id":"<material-uuid>","manufacturer_serial":"MFG-123"}'

# Record the first (INBOUND) movement
curl -X POST http://localhost:4045/api/v1/asset-movements \
  -H "Content-Type: application/json" \
  -d '{"asset_id":"<asset-uuid>","type":"INBOUND","movement_date":"2026-04-20","to_location_id":"<loc-uuid>"}'

# Public QR lookup
curl http://localhost:4045/public/asset/INT-0001-000001
```
