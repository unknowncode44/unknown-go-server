package shipment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository defines persistence operations for Shipment and ShipmentItem
// entities. The service layer depends on this interface to stay
// persistence-agnostic.
type Repository interface {
	Create(s *entities.Shipment) (*entities.Shipment, error)
	FindAll(status *entities.ShipmentStatus) ([]entities.Shipment, error)
	FindByID(id uuid.UUID) (*entities.Shipment, error)
	Update(s *entities.Shipment) (*entities.Shipment, error)

	// Items
	AddItem(item *entities.ShipmentItem) (*entities.ShipmentItem, error)
	UpdateItem(item *entities.ShipmentItem) (*entities.ShipmentItem, error)
	FindItemByID(shipmentID, itemID uuid.UUID) (*entities.ShipmentItem, error)
	FindAssetItem(shipmentID, assetID uuid.UUID) (*entities.ShipmentItem, error)
	FindMaterialItem(shipmentID, materialID uuid.UUID) (*entities.ShipmentItem, error)
	RemoveItem(shipmentID, itemID uuid.UUID) error

	// Process atomically generates the stock movements for every item of a
	// shipment and marks it as PROCESSED. See the method comment for why the
	// writes live here instead of in the reused domain services.
	Process(shipmentID uuid.UUID, destinationType entities.LocationType, inventoryByMaterial map[uuid.UUID]uuid.UUID) (*entities.Shipment, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepo returns a new GORM-backed Repository.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(s *entities.Shipment) (*entities.Shipment, error) {
	if err := r.db.Create(s).Error; err != nil {
		return nil, err
	}
	return r.FindByID(s.ID)
}

// preloads applies the standard read-side eager loading for shipments.
func (r *repository) preloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Items", func(d *gorm.DB) *gorm.DB { return d.Order("created_at asc") }).
		Preload("Items.Asset").
		Preload("Items.Material").
		Preload("DestinationLocation").
		Preload("CreatedByUser")
}

func (r *repository) FindAll(status *entities.ShipmentStatus) ([]entities.Shipment, error) {
	var list []entities.Shipment
	q := r.preloads(r.db).Where("is_active = ?", true)
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	if err := q.Order("created_at desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) FindByID(id uuid.UUID) (*entities.Shipment, error) {
	var s entities.Shipment
	if err := r.preloads(r.db).First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// Update persists the mutable fields of a shipment (state transitions, ERP
// linkage, soft delete) using an updates map so falsey values persist.
func (r *repository) Update(s *entities.Shipment) (*entities.Shipment, error) {
	if err := r.db.Model(&entities.Shipment{}).
		Where("id = ?", s.ID).
		Updates(map[string]interface{}{
			"destination_location_id": s.DestinationLocationID,
			"status":                  s.Status,
			"requester_name":          s.RequesterName,
			"notes":                   s.Notes,
			"erp_movement_number":     s.ERPMovementNumber,
			"processed_at":            s.ProcessedAt,
			"reconciled_at":           s.ReconciledAt,
			"is_active":               s.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return r.FindByID(s.ID)
}

func (r *repository) AddItem(item *entities.ShipmentItem) (*entities.ShipmentItem, error) {
	if err := r.db.Create(item).Error; err != nil {
		return nil, err
	}
	return r.FindItemByID(item.ShipmentID, item.ID)
}

func (r *repository) UpdateItem(item *entities.ShipmentItem) (*entities.ShipmentItem, error) {
	if err := r.db.Model(&entities.ShipmentItem{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"quantity": item.Quantity,
		}).Error; err != nil {
		return nil, err
	}
	return r.FindItemByID(item.ShipmentID, item.ID)
}

func (r *repository) FindItemByID(shipmentID, itemID uuid.UUID) (*entities.ShipmentItem, error) {
	var item entities.ShipmentItem
	if err := r.db.Preload("Asset").Preload("Material").
		First(&item, "id = ? AND shipment_id = ?", itemID, shipmentID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) FindAssetItem(shipmentID, assetID uuid.UUID) (*entities.ShipmentItem, error) {
	var item entities.ShipmentItem
	if err := r.db.
		First(&item, "shipment_id = ? AND type = ? AND asset_id = ?", shipmentID, entities.ShipmentItemTypeAsset, assetID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) FindMaterialItem(shipmentID, materialID uuid.UUID) (*entities.ShipmentItem, error) {
	var item entities.ShipmentItem
	if err := r.db.
		First(&item, "shipment_id = ? AND type = ? AND material_id = ?", shipmentID, entities.ShipmentItemTypeMaterial, materialID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repository) RemoveItem(shipmentID, itemID uuid.UUID) error {
	return r.db.Delete(&entities.ShipmentItem{}, "id = ? AND shipment_id = ?", itemID, shipmentID).Error
}

// Process performs the shipment confirmation as a single database transaction.
//
// The reused domain services (asset_movement.Service, delivery_record.Service)
// each open their own transaction against the root *gorm.DB, so calling them
// here would NOT enrol in a shared rollback — a mid-way failure could leave
// orphan movements. To honor the atomicity requirement ("ni un movimiento ni un
// delivery record huérfano"), the writes are done directly on `tx`, following
// the same repository-level transaction pattern already used by
// asset.CreateWithSerial and need.CreateWithNumber. Because the reused service
// is bypassed, the same safeguards AssetMovement.CreateWithAssetUpdate enforces
// are replicated inline for each asset item: a row-level lock on the asset, a
// rejection of retired/inactive assets, the first-movement rule (an asset with
// no prior movement can't be transferred/installed, since POST /assets does not
// register an automatic INBOUND), and the status/current_location_id side
// effect for the INSTALL/TRANSFER cases.
//
// inventoryByMaterial maps each material item's MaterialID to its resolved
// MaterialInventory ID; the service resolves and validates these before calling
// Process, so a material without inventory aborts before any write happens.
func (r *repository) Process(shipmentID uuid.UUID, destinationType entities.LocationType, inventoryByMaterial map[uuid.UUID]uuid.UUID) (*entities.Shipment, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var s entities.Shipment
		if err := tx.Preload("Items").First(&s, "id = ?", shipmentID).Error; err != nil {
			return err
		}

		now := time.Now()
		note := "Generado automáticamente por el envío " + shipmentID.String()
		destinationID := s.DestinationLocationID

		// Asset movement type is decided by the single shipment destination.
		movementType := entities.AssetMovementTransfer
		newStatus := entities.AssetStatusInStock
		if destinationType == entities.LocationTypeClient {
			movementType = entities.AssetMovementInstall
			newStatus = entities.AssetStatusInstalled
		}

		for i := range s.Items {
			item := &s.Items[i]

			switch item.Type {
			case entities.ShipmentItemTypeAsset:
				if item.AssetID == nil {
					return gorm.ErrRecordNotFound
				}

				// Lock the asset row for the whole transaction so concurrent
				// movements on the same asset can't overwrite each other.
				var asset entities.Asset
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
					First(&asset, "id = ?", *item.AssetID).Error; err != nil {
					return err
				}

				// A retired/inactive asset must not move, same as the manual flow.
				if !asset.IsActive || asset.Status == entities.AssetStatusRetired {
					return fmt.Errorf("asset %s is retired or inactive, cannot process shipment", asset.SerialVisible)
				}

				// First-movement rule: POST /assets does not create an automatic
				// INBOUND, so an asset with zero movements can't be
				// transferred/installed yet — its first movement must be an
				// inbound registered through the asset-movement flow.
				var movementCount int64
				if err := tx.Model(&entities.AssetMovement{}).
					Where("asset_id = ?", *item.AssetID).
					Count(&movementCount).Error; err != nil {
					return err
				}
				if movementCount == 0 {
					return fmt.Errorf("asset %s has no inbound movement registered, cannot process shipment", asset.SerialVisible)
				}

				toID := destinationID
				am := &entities.AssetMovement{
					AssetID:        *item.AssetID,
					Type:           movementType,
					FromLocationID: asset.CurrentLocationID,
					ToLocationID:   &toID,
					MovementDate:   now,
					Notes:          &note,
				}
				if err := tx.Create(am).Error; err != nil {
					return err
				}

				// Replicate the Asset side effect (status + current location).
				if err := tx.Model(&entities.Asset{}).
					Where("id = ?", *item.AssetID).
					Updates(map[string]interface{}{
						"status":              newStatus,
						"current_location_id": destinationID,
					}).Error; err != nil {
					return err
				}

				if err := tx.Model(&entities.ShipmentItem{}).
					Where("id = ?", item.ID).
					Update("resulting_asset_movement_id", am.ID).Error; err != nil {
					return err
				}

			case entities.ShipmentItemTypeMaterial:
				if item.MaterialID == nil || item.Quantity == nil {
					return gorm.ErrRecordNotFound
				}
				invID, ok := inventoryByMaterial[*item.MaterialID]
				if !ok || invID == uuid.Nil {
					// Should not happen: the service validates inventory first.
					return gorm.ErrRecordNotFound
				}

				dr := &entities.DeliveryRecord{
					MaterialInventoryID: invID,
					Type:                entities.DeliveryOutbound,
					Quantity:            *item.Quantity,
					DeliveryDate:        now,
					Notes:               &note,
				}
				if err := tx.Create(dr).Error; err != nil {
					return err
				}

				if err := tx.Model(&entities.ShipmentItem{}).
					Where("id = ?", item.ID).
					Update("resulting_delivery_record_id", dr.ID).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Model(&entities.Shipment{}).
			Where("id = ?", shipmentID).
			Updates(map[string]interface{}{
				"status":       entities.ShipmentStatusProcessed,
				"processed_at": now,
			}).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.FindByID(shipmentID)
}
