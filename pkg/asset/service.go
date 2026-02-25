package asset

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
)

// Service defines business operations for Asset entities.
type Service interface {
	Create(a *entities.Asset) (*entities.Asset, error)
	FindAll() ([]entities.Asset, error)
	FindByID(id uuid.UUID) (*entities.Asset, error)
	FindBySerial(serial string) (*entities.Asset, error)
	Update(a *entities.Asset) (*entities.Asset, error)
	Deactivate(id uuid.UUID) error
}

type service struct {
	repo         Repository
	materialRepo material.Repository
}

// NewService returns a new Asset Service.
func NewService(r Repository, mr material.Repository) Service {
	return &service{repo: r, materialRepo: mr}
}

func (s *service) Create(a *entities.Asset) (*entities.Asset, error) {
	if a == nil {
		return nil, errors.New("asset is required")
	}
	if a.MaterialID == uuid.Nil {
		return nil, errors.New("material_id is required")
	}

	// load material to get code
	mat, err := s.materialRepo.FindByID(a.MaterialID)
	if err != nil {
		return nil, err
	}
	if !mat.IsActive {
		return nil, errors.New("material is inactive")
	}
	code := strings.TrimSpace(mat.InternalCode)
	if code == "" {
		return nil, errors.New("material internal code is empty")
	}

	// initial state
	a.IsActive = true
	if a.Status == "" {
		a.Status = entities.AssetStatusInStock
	}

	created, err := s.repo.CreateWithSerial(a.MaterialID, code, a)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *service) FindAll() ([]entities.Asset, error) {
	return s.repo.FindAll()
}

func (s *service) FindByID(id uuid.UUID) (*entities.Asset, error) {
	return s.repo.FindByID(id)
}

func (s *service) FindBySerial(serial string) (*entities.Asset, error) {
	serial = strings.TrimSpace(serial)
	if serial == "" {
		return nil, errors.New("serial is required")
	}
	return s.repo.FindBySerial(serial)
}

func (s *service) Update(a *entities.Asset) (*entities.Asset, error) {
	if a == nil || a.ID == uuid.Nil {
		return nil, errors.New("asset id is required")
	}
	return s.repo.Update(a)
}

func (s *service) Deactivate(id uuid.UUID) error {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !a.IsActive {
		return nil
	}
	a.IsActive = false
	_, err = s.repo.Update(a)
	return err
}
