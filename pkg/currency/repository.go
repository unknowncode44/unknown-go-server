package currency

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository defines persistence operations for Currency entities.
type Repository interface {
	Create(currency *entities.Currency) (*entities.Currency, error)
	FindAll() ([]entities.Currency, error)
	FindByID(id uuid.UUID) (*entities.Currency, error)
	FindByCode(code string) (*entities.Currency, error)
	Update(currency *entities.Currency) (*entities.Currency, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepo returns a new Repository using the provided GORM DB.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new currency record into the database.
func (r *repository) Create(currency *entities.Currency) (*entities.Currency, error) {
	if err := r.db.Create(currency).Error; err != nil {
		return nil, err
	}
	return currency, nil
}

// FindAll retrieves all currency records.
func (r *repository) FindAll() ([]entities.Currency, error) {
	var items []entities.Currency
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID returns a currency by its UUID.
func (r *repository) FindByID(id uuid.UUID) (*entities.Currency, error) {
	var cur entities.Currency
	if err := r.db.First(&cur, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cur, nil
}

// FindByCode returns a currency by its code (case-sensitive match expected by caller).
func (r *repository) FindByCode(code string) (*entities.Currency, error) {
	var cur entities.Currency
	if err := r.db.First(&cur, "code = ?", code).Error; err != nil {
		return nil, err
	}
	return &cur, nil
}

// Update applies allowed changes to the currency record and returns the updated entity.
// Note: repository does not change the `code` field; service layer enforces that rule.
func (r *repository) Update(currency *entities.Currency) (*entities.Currency, error) {
	if err := r.db.Model(&entities.Currency{}).
		Where("id = ?", currency.ID).
		Updates(map[string]interface{}{
			"name":      currency.Name,
			"is_active": currency.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return currency, nil
}
