package delivery_record

import (
	"errors"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material_inventory"
)

// Service defines business operations for DeliveryRecord.
type Service interface {
	Create(dr *entities.DeliveryRecord) (*entities.DeliveryRecord, error)
	FindByInventoryID(inventoryID uuid.UUID) ([]entities.DeliveryRecord, error)
	FindAll() ([]entities.DeliveryRecord, error)
	GetTotalQuantity(inventoryID uuid.UUID) (total float64, err error)
}

type service struct {
	repo    Repository
	invRepo material_inventory.Repository
}

// NewService returns a new Service using the provided repositories.
func NewService(repo Repository, invRepo material_inventory.Repository) Service {
	return &service{repo: repo, invRepo: invRepo}
}

func (s *service) Create(dr *entities.DeliveryRecord) (*entities.DeliveryRecord, error) {
	if dr == nil {
		return nil, errors.New("delivery record is required")
	}
	if dr.MaterialInventoryID == uuid.Nil {
		return nil, errors.New("material_inventory_id is required")
	}
	if dr.Type != entities.DeliveryInbound && dr.Type != entities.DeliveryOutbound {
		return nil, errors.New("type must be INBOUND or OUTBOUND")
	}
	if dr.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	if dr.DeliveryDate.IsZero() {
		return nil, errors.New("delivery_date is required")
	}

	// validate inventory exists and is active
	inv, err := s.invRepo.FindByID(dr.MaterialInventoryID)
	if err != nil {
		return nil, errors.New("material inventory not found")
	}
	if !inv.IsActive {
		return nil, errors.New("material inventory is not active")
	}

	return s.repo.Create(dr)
}

func (s *service) FindByInventoryID(inventoryID uuid.UUID) ([]entities.DeliveryRecord, error) {
	return s.repo.FindByInventoryID(inventoryID)
}

func (s *service) FindAll() ([]entities.DeliveryRecord, error) {
	return s.repo.FindAll()
}

func (s *service) GetTotalQuantity(inventoryID uuid.UUID) (float64, error) {
	inv, err := s.invRepo.FindByID(inventoryID)
	if err != nil {
		return 0, errors.New("material inventory not found")
	}

	inbound, outbound, err := s.repo.SumByInventoryID(inventoryID)
	if err != nil {
		return 0, err
	}

	total := inv.InitialQuantity + inbound - outbound
	return total, nil
}
