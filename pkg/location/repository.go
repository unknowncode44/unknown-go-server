package location

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository defines persistence operations for Location entities.
type Repository interface {
	Create(l *entities.Location) (*entities.Location, error)
	FindAll() ([]entities.Location, error)
	FindByID(id uuid.UUID) (*entities.Location, error)
	Update(l *entities.Location) (*entities.Location, error)
	Deactivate(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(l *entities.Location) (*entities.Location, error) {
	if err := r.db.Create(l).Error; err != nil {
		return nil, err
	}
	return l, nil
}

func (r *repository) FindAll() ([]entities.Location, error) {
	var list []entities.Location
	if err := r.db.Preload("ChildLocations").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) FindByID(id uuid.UUID) (*entities.Location, error) {
	var l entities.Location
	if err := r.db.Preload("ChildLocations").First(&l, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *repository) Update(l *entities.Location) (*entities.Location, error) {
	if err := r.db.Model(&entities.Location{}).
		Where("id = ?", l.ID).
		Updates(map[string]interface{}{
			"name":               l.Name,
			"type":               l.Type,
			"parent_location_id": l.ParentLocationID,
			"is_active":          l.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return l, nil
}

func (r *repository) Deactivate(id uuid.UUID) error {
	return r.db.Model(&entities.Location{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
