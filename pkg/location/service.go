package location

import (
	"errors"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Service defines business operations for Location.
type Service interface {
	Create(l *entities.Location) (*entities.Location, error)
	FindAll() ([]entities.Location, error)
	FindByID(id uuid.UUID) (*entities.Location, error)
	Update(l *entities.Location) (*entities.Location, error)
	Deactivate(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(l *entities.Location) (*entities.Location, error) {
	if l == nil {
		return nil, errors.New("location is required")
	}
	if l.Name == "" {
		return nil, errors.New("name is required")
	}
	if l.Type == "" {
		return nil, errors.New("type is required")
	}
	l.IsActive = true
	return s.repo.Create(l)
}

func (s *service) FindAll() ([]entities.Location, error) {
	return s.repo.FindAll()
}

func (s *service) FindByID(id uuid.UUID) (*entities.Location, error) {
	return s.repo.FindByID(id)
}

func (s *service) Update(l *entities.Location) (*entities.Location, error) {
	if l == nil || l.ID == uuid.Nil {
		return nil, errors.New("location id is required")
	}
	if l.Name == "" {
		return nil, errors.New("name is required")
	}
	return s.repo.Update(l)
}

func (s *service) Deactivate(id uuid.UUID) error {
	return s.repo.Deactivate(id)
}
