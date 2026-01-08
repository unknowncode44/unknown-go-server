package vendor_material

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Service defines business operations for VendorMaterial.
type Service interface {
	Create(vm *entities.VendorMaterial) (*entities.VendorMaterial, error)
	FindAll() ([]entities.VendorMaterial, error)
	FindByID(id uuid.UUID) (*entities.VendorMaterial, error)
	Update(vm *entities.VendorMaterial) (*entities.VendorMaterial, error)
	Deactivate(id uuid.UUID) error
}

type service struct {
	repo Repository
}

// NewService returns a new Service using the provided Repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(vm *entities.VendorMaterial) (*entities.VendorMaterial, error) {
	if vm == nil {
		return nil, errors.New("vendor_material is required")
	}
	if vm.VendorID == uuid.Nil {
		return nil, errors.New("vendor_id is required")
	}
	if vm.MaterialID == uuid.Nil {
		return nil, errors.New("material_id is required")
	}
	if strings.TrimSpace(vm.VendorCode) == "" {
		return nil, errors.New("vendor_code is required")
	}

	vm.IsActive = true
	return s.repo.Create(vm)
}

func (s *service) FindAll() ([]entities.VendorMaterial, error) {
	return s.repo.FindAll()
}

func (s *service) FindByID(id uuid.UUID) (*entities.VendorMaterial, error) {
	return s.repo.FindByID(id)
}

func (s *service) Update(vm *entities.VendorMaterial) (*entities.VendorMaterial, error) {
	if vm == nil || vm.ID == uuid.Nil {
		return nil, errors.New("vendor_material id is required")
	}

	// Only vendor_code is updatable via API for now; keep validations minimal.
	if strings.TrimSpace(vm.VendorCode) == "" {
		return nil, errors.New("vendor_code is required")
	}

	return s.repo.Update(vm)
}

func (s *service) Deactivate(id uuid.UUID) error {
	vm, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !vm.IsActive {
		return nil
	}
	vm.IsActive = false
	_, err = s.repo.Update(vm)
	return err
}
