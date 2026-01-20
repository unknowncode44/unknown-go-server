package material_cost

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Service defines business operations for MaterialCost.
type Service interface {
	Create(mc *entities.MaterialCost) (*entities.MaterialCost, error)
	FindAll() ([]entities.MaterialCost, error)
	FindByID(id uuid.UUID) (*entities.MaterialCost, error)
	FindByMaterial(materialID uuid.UUID) ([]entities.MaterialCost, error)
	FindLatestByMaterialAndVendor(materialID uuid.UUID, vendorID uuid.UUID) (*entities.MaterialCost, error)
}

type service struct {
	repo Repository
}

// NewService returns a new Service using the provided Repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(mc *entities.MaterialCost) (*entities.MaterialCost, error) {
	if mc == nil {
		return nil, errors.New("material cost is required")
	}

	if mc.MaterialID == uuid.Nil {
		return nil, errors.New("material_id is required")
	}
	if mc.VendorID == uuid.Nil {
		return nil, errors.New("vendor_id is required")
	}
	if mc.CurrencyID == uuid.Nil {
		return nil, errors.New("currency_id is required")
	}
	if mc.Cost <= 0 {
		return nil, errors.New("cost must be greater than zero")
	}
	if mc.CostDate.IsZero() || mc.CostDate.Equal(time.Time{}) {
		return nil, errors.New("cost_date is required")
	}

	return s.repo.Create(mc)
}

func (s *service) FindAll() ([]entities.MaterialCost, error) {
	return s.repo.FindAll()
}

func (s *service) FindByID(id uuid.UUID) (*entities.MaterialCost, error) {
	return s.repo.FindByID(id)
}

func (s *service) FindByMaterial(materialID uuid.UUID) ([]entities.MaterialCost, error) {
	return s.repo.FindByMaterial(materialID)
}

func (s *service) FindLatestByMaterialAndVendor(materialID uuid.UUID, vendorID uuid.UUID) (*entities.MaterialCost, error) {
	return s.repo.FindLatestByMaterialAndVendor(materialID, vendorID)
}
