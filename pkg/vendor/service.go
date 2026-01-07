package vendor

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Package vendor provides the domain service and repository contracts
// and default implementations for vendor (supplier) entities.

// Service defines the business API for vendor operations. Handlers use
// this interface to perform domain operations without depending on
// persistence details.
type Service interface {
	Create(vendor *entities.Vendor) (*entities.Vendor, error)
	FindAll() ([]entities.Vendor, error)
	FindByID(id uuid.UUID) (*entities.Vendor, error)
	Update(vendor *entities.Vendor) (*entities.Vendor, error)
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
func (s *service) Create(vendor *entities.Vendor) (*entities.Vendor, error) {
	if vendor == nil {
		return nil, errors.New("vendor is required")
	}

	if strings.TrimSpace(vendor.Name) == "" {
		return nil, errors.New("vendor name is required")
	}

	vendor.IsActive = true
	return s.repo.Create(vendor)
}

// FindAll returns all vendors.
func (s *service) FindAll() ([]entities.Vendor, error) {
	return s.repo.FindAll()
}

// FindByID returns a vendor by UUID.
func (s *service) FindByID(id uuid.UUID) (*entities.Vendor, error) {
	return s.repo.FindByID(id)
}

// Update validates the entity and delegates the update to the repository.
func (s *service) Update(vendor *entities.Vendor) (*entities.Vendor, error) {
	if vendor == nil || vendor.ID == uuid.Nil {
		return nil, errors.New("vendor id is required")
	}
	if strings.TrimSpace(vendor.Name) == "" {
		return nil, errors.New("vendor name is required")
	}

	return s.repo.Update(vendor)
}

// Deactivate performs a logical delete by setting IsActive to false.
// The operation is idempotent.
func (s *service) Deactivate(id uuid.UUID) error {
	v, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !v.IsActive {
		return nil
	}

	v.IsActive = false
	_, err = s.repo.Update(v)
	return err
}
