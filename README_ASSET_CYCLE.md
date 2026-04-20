# Asset Movement Cycle (unknown-go-server)

This document explains the end-to-end cycle for turning a Material into an Asset and performing asset movements (inbound, transfer, install, etc.). It lists required resources, endpoints, and example payloads.

## Overview
- Materials are master records (`/api/v1/materials`).
- Assets are instances created from Materials (`/api/v1/assets`).
- Asset movements are recorded by creating `AssetMovement` entries (`/api/v1/asset-movements`) which atomically update asset status and current location.
- Locations are used as movement origins/destinations (`/api/v1/locations`).

## Preconditions
- Material must exist and be active.
- Material must have a non-empty `internalCode` to allow creating Assets (checked in `pkg/asset/service.go`).
- Locations referenced in movements must exist.
- IDs must be valid UUID strings.

## Endpoints (full cycle)
- Create Material: POST `/api/v1/materials`
- Create Location: POST `/api/v1/locations`
- Create Asset: POST `/api/v1/assets`
- Record Asset Movement: POST `/api/v1/asset-movements`
- List Asset Movements: GET `/api/v1/assets/:asset_id/movements`
- Get Last Movement: GET `/api/v1/assets/:asset_id/movements/last`
- Get Asset: GET `/api/v1/assets/:id` or `/api/v1/assets/serial/:serial`

## Required fields and validations
- Material creation requires: `name`, `sector`, `unitOfMeasure`, `code`.
- Location creation requires: `name`, `type`.
- Asset creation requires: `material_id` (UUID) and the material must be active with a non-empty `internalCode`.
- AssetMovement creation requires: `asset_id` (UUID), `type` (one of `INBOUND`, `TRANSFER`, `INSTALL`, `UNINSTALL`, `REPAIR`, `RETURN`, `DECOMMISSION`, `SCRAP`), and `movement_date` (RFC3339 or `YYYY-MM-DD`).
- First movement for an asset must be `INBOUND`.

## Example payloads

Create Material (required fields)

```json
{
  "name": "Steel Bolt M8",
  "sector": "Hardware",
  "unitOfMeasure": "pcs",
  "code": "BOLT-M8-001",
  "erpCode": "ERP-12345",
  "internalCode": "INT-0001"
}
```

Create Location

```json
{
  "name": "Main Warehouse",
  "type": "WAREHOUSE"
}
```

Create Asset (instance of a material)

```json
{
  "material_id": "11111111-2222-3333-4444-555555555555",
  "manufacturer_serial": "MFG-12345",
  "current_location_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}
```

Create Asset Movement — INBOUND (first movement)

```json
{
  "asset_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "type": "INBOUND",
  "to_location_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "movement_date": "2026-04-20T15:00:00Z",
  "notes": "Received from supplier"
}
```

Create Asset Movement — TRANSFER

```json
{
  "asset_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "type": "TRANSFER",
  "from_location_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "to_location_id": "99999999-8888-7777-6666-555555555555",
  "movement_date": "2026-04-21",
  "notes": "Moved to site B"
}
```

Decommission / Scrap (finalize asset)

```json
{
  "asset_id": "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
  "type": "DECOMMISSION",
  "movement_date": "2026-05-01",
  "notes": "End of life"
}
```

## Behavior notes
- Asset creation generates a visible serial (QR) using the Material `internalCode` (see `pkg/asset/service.go` and repository creation). The asset initial `status` defaults to `IN_STOCK` unless provided.
- AssetMovement creation updates `Asset.Status` and `Asset.CurrentLocationID` according to movement type (see `pkg/asset_movement/repository.go`).
- Movement creation is transactional: movement and asset update are committed atomically.
- If you need per-material stock counts (inbound/outbound totals), add a `material_stock` or `material_transaction` model and aggregate AssetMovement/transactions.

## Quick curl examples

Create material

```bash
curl -X POST http://localhost:8080/api/v1/materials \
  -H "Content-Type: application/json" \
  -d '{"name":"Steel Bolt M8","sector":"Hardware","unitOfMeasure":"pcs","code":"BOLT-M8-001","internalCode":"INT-0001"}'
```

Create asset

```bash
curl -X POST http://localhost:8080/api/v1/assets \
  -H "Content-Type: application/json" \
  -d '{"material_id":"11111111-2222-3333-4444-555555555555","manufacturer_serial":"MFG-12345"}'
```

Record inbound movement

```bash
curl -X POST http://localhost:8080/api/v1/asset-movements \
  -H "Content-Type: application/json" \
  -d '{"asset_id":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee","type":"INBOUND","movement_date":"2026-04-20T15:00:00Z","to_location_id":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"}'
```

---

File: README_ASSET_CYCLE.md
