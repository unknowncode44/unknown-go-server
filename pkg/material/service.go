package material

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Package material provides domain interfaces and implementations for
// material entities used by the application.

// Service defines the business API for material operations.
type Service interface {
	Create(material *entities.Material) (*entities.Material, error)
	FindAll() ([]entities.Material, error)
	FindByID(id uuid.UUID) (*entities.Material, error)
	Update(material *entities.Material) (*entities.Material, error)
	Deactivate(id uuid.UUID) error
}

// service is the default Service implementation.
type service struct {
	repo Repository
}

// NewService returns a new Service using the provided Repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Create validates business rules and delegates persistence to the repository.
// Required fields are trimmed and must be non-empty.
func (s *service) Create(material *entities.Material) (*entities.Material, error) {
	if material == nil {
		return nil, errors.New("material is required")
	}

	if strings.TrimSpace(material.Name) == "" {
		return nil, errors.New("material name is required")
	}
	if strings.TrimSpace(material.Sector) == "" {
		return nil, errors.New("material sector is required")
	}
	if strings.TrimSpace(material.UnitOfMeasure) == "" {
		return nil, errors.New("unit of measure is required")
	}
	material.Code = strings.TrimSpace(material.Code)
	if material.Code == "" {
		return nil, errors.New("material code is required")
	}

	// check duplicates
	if existing, err := s.repo.FindByCode(material.Code); err == nil && existing != nil {
		return nil, errors.New("material code already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	material.IsActive = true
	return s.repo.Create(material)
}

// FindAll returns all materials.
func (s *service) FindAll() ([]entities.Material, error) {
	return s.repo.FindAll()
}

// FindByID returns a material by UUID.
func (s *service) FindByID(id uuid.UUID) (*entities.Material, error) {
	return s.repo.FindByID(id)
}

// Update validates the entity and delegates the update to the repository.
func (s *service) Update(material *entities.Material) (*entities.Material, error) {
	if material == nil || material.ID == uuid.Nil {
		return nil, errors.New("material id is required")
	}
	if strings.TrimSpace(material.Name) == "" {
		return nil, errors.New("material name is required")
	}
	return s.repo.Update(material)
}

// Deactivate performs a logical delete by setting IsActive to false.
// The operation is idempotent.
func (s *service) Deactivate(id uuid.UUID) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !m.IsActive {
		return nil
	}
	m.IsActive = false
	_, err = s.repo.Update(m)
	return err
}
