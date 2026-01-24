package asset

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository defines persistence operations for Asset entities.
type Repository interface {
	Create(a *entities.Asset) (*entities.Asset, error)
	CreateWithSerial(materialID uuid.UUID, materialCode string, a *entities.Asset) (*entities.Asset, error)
	FindAll() ([]entities.Asset, error)
	FindByID(id uuid.UUID) (*entities.Asset, error)
	FindBySerial(serial string) (*entities.Asset, error)
	Update(a *entities.Asset) (*entities.Asset, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepo returns a new GORM-backed Repository.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(a *entities.Asset) (*entities.Asset, error) {
	if err := r.db.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

// CreateWithSerial generates a serial for the given material using a
// transactional counter and creates the asset atomically.
func (r *repository) CreateWithSerial(materialID uuid.UUID, materialCode string, a *entities.Asset) (*entities.Asset, error) {
	// start transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// lock the counter row for update
	var counter entities.AssetSerialCounter
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&counter, "material_id = ?", materialID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// create counter
			counter = entities.AssetSerialCounter{MaterialID: materialID, Last: 1}
			if err := tx.Create(&counter).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			tx.Rollback()
			return nil, err
		}
	} else {
		counter.Last = counter.Last + 1
		if err := tx.Model(&entities.AssetSerialCounter{}).Where("id = ?", counter.ID).Update("last", counter.Last).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// generate serial
	serial := fmt.Sprintf("%s-%04d", materialCode, counter.Last)
	a.SerialVisible = serial
	a.CreatedAt = time.Now()

	if err := tx.Create(a).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return a, nil
}

func (r *repository) FindAll() ([]entities.Asset, error) {
	var list []entities.Asset
	if err := r.db.Preload("Material").Preload("CurrentLocation").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) FindByID(id uuid.UUID) (*entities.Asset, error) {
	var a entities.Asset
	if err := r.db.Preload("Material").Preload("CurrentLocation").First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) FindBySerial(serial string) (*entities.Asset, error) {
	var a entities.Asset
	if err := r.db.Preload("Material").Preload("CurrentLocation").First(&a, "serial_visible = ?", serial).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) Update(a *entities.Asset) (*entities.Asset, error) {
	if err := r.db.Model(&entities.Asset{}).
		Where("id = ?", a.ID).
		Updates(map[string]interface{}{
			"serial_visible":      a.SerialVisible,
			"manufacturer_serial": a.ManufacturerSerial,
			"status":              a.Status,
			"parent_asset_id":     a.ParentAssetID,
			"current_location_id": a.CurrentLocationID,
			"is_active":           a.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return a, nil
}
