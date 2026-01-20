package currency

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Service defines business operations for Currency entities.
type Service interface {
	Create(currency *entities.Currency) (*entities.Currency, error)
	FindAll() ([]entities.Currency, error)
	FindByID(id uuid.UUID) (*entities.Currency, error)
	FindByCode(code string) (*entities.Currency, error)
	Update(currency *entities.Currency) (*entities.Currency, error)
	Deactivate(id uuid.UUID) error
}

type service struct {
	repo Repository
}

// NewService returns a new Service using the provided Repository.
func NewService(r Repository) Service {
	return &service{repo: r}
}

// Create validates and creates a new currency. It normalizes code to uppercase
// and enforces uniqueness.
func (s *service) Create(currency *entities.Currency) (*entities.Currency, error) {
	if currency == nil {
		return nil, errors.New("currency is required")
	}

	currency.Code = strings.ToUpper(strings.TrimSpace(currency.Code))
	currency.Name = strings.TrimSpace(currency.Name)

	if currency.Code == "" {
		return nil, errors.New("currency code is required")
	}
	if len(currency.Code) != 3 {
		return nil, errors.New("currency code must be 3 characters")
	}
	if currency.Name == "" {
		return nil, errors.New("currency name is required")
	}

	// check duplicates
	if existing, err := s.repo.FindByCode(currency.Code); err == nil && existing != nil {
		return nil, errors.New("currency code already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	currency.IsActive = true
	return s.repo.Create(currency)
}

// FindAll returns all currencies.
func (s *service) FindAll() ([]entities.Currency, error) {
	return s.repo.FindAll()
}

// FindByID returns a currency by UUID.
func (s *service) FindByID(id uuid.UUID) (*entities.Currency, error) {
	return s.repo.FindByID(id)
}

// FindByCode returns a currency by code.
func (s *service) FindByCode(code string) (*entities.Currency, error) {
	return s.repo.FindByCode(strings.ToUpper(strings.TrimSpace(code)))
}

// Update allows updating only the name of a currency; code cannot be changed.
func (s *service) Update(currency *entities.Currency) (*entities.Currency, error) {
	if currency == nil || currency.ID == uuid.Nil {
		return nil, errors.New("currency id is required")
	}

	if strings.TrimSpace(currency.Name) == "" {
		return nil, errors.New("currency name is required")
	}

	// load existing to ensure code isn't changed
	existing, err := s.repo.FindByID(currency.ID)
	if err != nil {
		return nil, err
	}
	existing.Name = strings.TrimSpace(currency.Name)
	return s.repo.Update(existing)
}

// Deactivate performs a logical delete by setting IsActive to false.
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
