package shipment

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/asset"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/location"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
	"github.com/unknowncode44/unknown-go-server/pkg/material_inventory"
)

// Service defines the business operations for the shipment workflow: a user
// builds a DRAFT shipment scanning assets/materials, processes it (which
// atomically generates the corresponding stock movements), and later links it
// to an external ERP movement number for accounting reconciliation.
type Service interface {
	Create(destinationLocationID uuid.UUID, createdByUserID uuid.UUID, requesterName string, notes string) (*entities.Shipment, error)
	FindAll(status *entities.ShipmentStatus) ([]entities.Shipment, error)
	FindByID(id uuid.UUID) (*entities.Shipment, error)
	AddAssetItem(shipmentID uuid.UUID, serial string) (*entities.ShipmentItem, error)
	AddMaterialItem(shipmentID uuid.UUID, code string, quantity float64) (*entities.ShipmentItem, error)
	RemoveItem(shipmentID, itemID uuid.UUID) error
	Process(shipmentID uuid.UUID) (*entities.Shipment, error)
	LinkERPMovement(shipmentID uuid.UUID, erpMovementNumber string) (*entities.Shipment, error)
	Cancel(shipmentID uuid.UUID) error
}

type service struct {
	repo Repository

	// Reused domain services for the read side of the workflow (resolving
	// serials/codes/locations/inventories). The transactional write side of
	// Process lives in the repository — see Repository.Process for why.
	assetSvc    asset.Service
	materialSvc material.Service
	invSvc      material_inventory.Service
	locationSvc location.Service
}

// NewService wires the shipment service with the reused domain services.
func NewService(
	repo Repository,
	assetSvc asset.Service,
	materialSvc material.Service,
	invSvc material_inventory.Service,
	locationSvc location.Service,
) Service {
	return &service{
		repo:        repo,
		assetSvc:    assetSvc,
		materialSvc: materialSvc,
		invSvc:      invSvc,
		locationSvc: locationSvc,
	}
}

func (s *service) Create(destinationLocationID uuid.UUID, createdByUserID uuid.UUID, requesterName string, notes string) (*entities.Shipment, error) {
	requesterName = strings.TrimSpace(requesterName)
	if requesterName == "" {
		return nil, errors.New("requester_name is required")
	}
	if destinationLocationID == uuid.Nil {
		return nil, errors.New("destination_location_id is required")
	}
	if createdByUserID == uuid.Nil {
		return nil, errors.New("created_by_user_id is required")
	}

	// validate destination location exists
	if _, err := s.locationSvc.FindByID(destinationLocationID); err != nil {
		return nil, errors.New("destination location not found")
	}

	shipment := &entities.Shipment{
		DestinationLocationID: destinationLocationID,
		CreatedByUserID:       createdByUserID,
		RequesterName:         requesterName,
		Notes:                 strings.TrimSpace(notes),
		Status:                entities.ShipmentStatusDraft,
		IsActive:              true,
	}
	return s.repo.Create(shipment)
}

func (s *service) FindAll(status *entities.ShipmentStatus) ([]entities.Shipment, error) {
	return s.repo.FindAll(status)
}

func (s *service) FindByID(id uuid.UUID) (*entities.Shipment, error) {
	return s.repo.FindByID(id)
}

// requireDraft loads a shipment and ensures it can still be edited.
func (s *service) requireDraft(shipmentID uuid.UUID) (*entities.Shipment, error) {
	shipment, err := s.repo.FindByID(shipmentID)
	if err != nil {
		return nil, errors.New("shipment not found")
	}
	if shipment.Status != entities.ShipmentStatusDraft {
		return nil, errors.New("shipment is not in DRAFT status")
	}
	return shipment, nil
}

func (s *service) AddAssetItem(shipmentID uuid.UUID, serial string) (*entities.ShipmentItem, error) {
	if _, err := s.requireDraft(shipmentID); err != nil {
		return nil, err
	}

	a, err := s.assetSvc.FindBySerial(serial)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	// reject the same asset twice in this shipment
	if existing, _ := s.repo.FindAssetItem(shipmentID, a.ID); existing != nil && existing.ID != uuid.Nil {
		return nil, errors.New("asset already added to this shipment")
	}

	assetID := a.ID
	item := &entities.ShipmentItem{
		ShipmentID: shipmentID,
		Type:       entities.ShipmentItemTypeAsset,
		AssetID:    &assetID,
	}
	return s.repo.AddItem(item)
}

func (s *service) AddMaterialItem(shipmentID uuid.UUID, code string, quantity float64) (*entities.ShipmentItem, error) {
	if _, err := s.requireDraft(shipmentID); err != nil {
		return nil, err
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	mat, err := s.materialSvc.FindByCode(strings.TrimSpace(code))
	if err != nil {
		return nil, errors.New("material not found")
	}

	// if the material is already on the shipment, add up the quantity instead
	// of creating a duplicate line.
	if existing, _ := s.repo.FindMaterialItem(shipmentID, mat.ID); existing != nil && existing.ID != uuid.Nil {
		current := 0.0
		if existing.Quantity != nil {
			current = *existing.Quantity
		}
		newQty := current + quantity
		existing.Quantity = &newQty
		return s.repo.UpdateItem(existing)
	}

	materialID := mat.ID
	qty := quantity
	item := &entities.ShipmentItem{
		ShipmentID: shipmentID,
		Type:       entities.ShipmentItemTypeMaterial,
		MaterialID: &materialID,
		Quantity:   &qty,
	}
	return s.repo.AddItem(item)
}

func (s *service) RemoveItem(shipmentID, itemID uuid.UUID) error {
	if _, err := s.requireDraft(shipmentID); err != nil {
		return err
	}
	if _, err := s.repo.FindItemByID(shipmentID, itemID); err != nil {
		return errors.New("item not found in this shipment")
	}
	return s.repo.RemoveItem(shipmentID, itemID)
}

func (s *service) Process(shipmentID uuid.UUID) (*entities.Shipment, error) {
	shipment, err := s.repo.FindByID(shipmentID)
	if err != nil {
		return nil, errors.New("shipment not found")
	}
	if shipment.Status != entities.ShipmentStatusDraft {
		return nil, errors.New("shipment is not in DRAFT status")
	}
	if len(shipment.Items) == 0 {
		return nil, errors.New("shipment has no items")
	}

	// destination type drives the asset movement type (INSTALL for CLIENT,
	// TRANSFER for INTERNAL).
	dest, err := s.locationSvc.FindByID(shipment.DestinationLocationID)
	if err != nil {
		return nil, errors.New("destination location not found")
	}

	// Resolve (and validate) every material's inventory BEFORE any write, so a
	// material without inventory aborts the whole operation with no orphan
	// records. The map is handed to the repository transaction.
	inventoryByMaterial := make(map[uuid.UUID]uuid.UUID)
	for i := range shipment.Items {
		item := shipment.Items[i]
		if item.Type != entities.ShipmentItemTypeMaterial || item.MaterialID == nil {
			continue
		}
		if _, done := inventoryByMaterial[*item.MaterialID]; done {
			continue
		}
		code := ""
		if item.Material != nil {
			code = item.Material.Code
		}
		inv, err := s.invSvc.FindByMaterialID(*item.MaterialID)
		if err != nil || inv == nil {
			return nil, fmt.Errorf("no inventory found for material %s, cannot process shipment", code)
		}
		// Parity with delivery_record.Service.Create, which rejects OUTBOUND
		// against an inactive inventory (FindByMaterialID does not filter it out).
		if !inv.IsActive {
			return nil, fmt.Errorf("inventory for material %s is not active, cannot process shipment", code)
		}
		inventoryByMaterial[*item.MaterialID] = inv.ID
	}

	return s.repo.Process(shipmentID, dest.Type, inventoryByMaterial)
}

func (s *service) LinkERPMovement(shipmentID uuid.UUID, erpMovementNumber string) (*entities.Shipment, error) {
	erpMovementNumber = strings.TrimSpace(erpMovementNumber)
	if erpMovementNumber == "" {
		return nil, errors.New("erp_movement_number is required")
	}

	shipment, err := s.repo.FindByID(shipmentID)
	if err != nil {
		return nil, errors.New("shipment not found")
	}
	if shipment.Status != entities.ShipmentStatusProcessed {
		return nil, errors.New("shipment is not in PROCESSED status")
	}

	now := time.Now()
	shipment.ERPMovementNumber = &erpMovementNumber
	shipment.Status = entities.ShipmentStatusReconciled
	shipment.ReconciledAt = &now
	return s.repo.Update(shipment)
}

func (s *service) Cancel(shipmentID uuid.UUID) error {
	shipment, err := s.repo.FindByID(shipmentID)
	if err != nil {
		return errors.New("shipment not found")
	}
	if shipment.Status != entities.ShipmentStatusDraft {
		return errors.New("only DRAFT shipments can be cancelled")
	}
	shipment.IsActive = false
	_, err = s.repo.Update(shipment)
	return err
}
