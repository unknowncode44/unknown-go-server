package material_inventory

import (
	"errors"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
	"gorm.io/gorm"
)

// Service defines business operations for MaterialInventory.
type Service interface {
	Create(inv *entities.MaterialInventory) (*entities.MaterialInventory, error)
	FindAll() ([]entities.MaterialInventory, error)
	FindByID(id uuid.UUID) (*entities.MaterialInventory, error)
	FindByMaterialID(materialID uuid.UUID) (*entities.MaterialInventory, error)
	Update(inv *entities.MaterialInventory) (*entities.MaterialInventory, error)
	Deactivate(id uuid.UUID) error
}

type service struct {
	repo         Repository
	materialRepo material.Repository
}

// NewService returns a new Service using the provided repositories.
func NewService(repo Repository, materialRepo material.Repository) Service {
	return &service{repo: repo, materialRepo: materialRepo}
}

func (s *service) Create(inv *entities.MaterialInventory) (*entities.MaterialInventory, error) {
	if inv == nil {
		return nil, errors.New("material inventory is required")
	}
	if inv.MaterialID == uuid.Nil {
		return nil, errors.New("material_id is required")
	}
	if inv.StartDate.IsZero() {
		return nil, errors.New("start_date is required")
	}
	if inv.InitialQuantity < 0 {
		return nil, errors.New("initial_quantity must be >= 0")
	}

	// validate material exists and is active
	mat, err := s.materialRepo.FindByID(inv.MaterialID)
	if err != nil {
		return nil, errors.New("material not found")
	}
	if !mat.IsActive {
		return nil, errors.New("material is not active")
	}

	// check no duplicate active inventory for the same material
	existing, err := s.repo.FindByMaterialID(inv.MaterialID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil && existing.ID != uuid.Nil {
		return nil, errors.New("inventory already exists for this material")
	}

	inv.IsActive = true
	return s.repo.Create(inv)
}

func (s *service) FindAll() ([]entities.MaterialInventory, error) {
	return s.repo.FindAll()
}

func (s *service) FindByID(id uuid.UUID) (*entities.MaterialInventory, error) {
	return s.repo.FindByID(id)
}

func (s *service) FindByMaterialID(materialID uuid.UUID) (*entities.MaterialInventory, error) {
	return s.repo.FindByMaterialID(materialID)
}

func (s *service) Update(inv *entities.MaterialInventory) (*entities.MaterialInventory, error) {
	if inv == nil || inv.ID == uuid.Nil {
		return nil, errors.New("material inventory id is required")
	}
	return s.repo.Update(inv)
}

func (s *service) Deactivate(id uuid.UUID) error {
	inv, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !inv.IsActive {
		return nil
	}
	inv.IsActive = false
	_, err = s.repo.Update(inv)
	return err
}
