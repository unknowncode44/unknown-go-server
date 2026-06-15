package need

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository defines persistence operations for Need and NeedItem entities.
// The service layer depends on this interface to remain persistence-agnostic.
type Repository interface {
	Create(n *entities.Need) (*entities.Need, error)
	CreateWithNumber(n *entities.Need) (*entities.Need, error)
	FindAll() ([]entities.Need, error)
	FindByID(id uuid.UUID) (*entities.Need, error)
	Update(n *entities.Need) (*entities.Need, error)
	UpdateStatus(id uuid.UUID, status entities.NeedStatus) (*entities.Need, error)

	// Items
	AddItem(item *entities.NeedItem) (*entities.NeedItem, error)
	ListItems(needID uuid.UUID) ([]entities.NeedItem, error)
	FindItemByID(itemID uuid.UUID) (*entities.NeedItem, error)
	UpdateItem(item *entities.NeedItem) (*entities.NeedItem, error)
	RemoveItem(itemID uuid.UUID) error
}

// repository is the GORM-backed Repository implementation.
type repository struct {
	db *gorm.DB
}

// NewRepo returns a new Repository using the provided GORM DB.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new need record into the database.
func (r *repository) Create(n *entities.Need) (*entities.Need, error) {
	if err := r.db.Create(n).Error; err != nil {
		return nil, err
	}
	return n, nil
}

// CreateWithNumber generates a per-year correlative number for the need using a
// transactional counter (NEED-YYYY-NNNN) and creates the need atomically.
func (r *repository) CreateWithNumber(n *entities.Need) (*entities.Need, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()

	year := time.Now().Year()

	// lock the counter row for update
	var counter entities.NeedCounter
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&counter, "year = ?", year).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// create counter for this year
			counter = entities.NeedCounter{Year: year, Last: 1}
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
		if err := tx.Model(&entities.NeedCounter{}).Where("id = ?", counter.ID).Update("last", counter.Last).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// generate number NEED-YYYY-NNNN with 4-digit zero-padding
	n.Number = fmt.Sprintf("NEED-%d-%04d", year, counter.Last)

	if err := tx.Create(n).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return n, nil
}

// FindAll retrieves all needs with their items and the items' materials.
func (r *repository) FindAll() ([]entities.Need, error) {
	var needs []entities.Need
	if err := r.db.Preload("Items.Material").Find(&needs).Error; err != nil {
		return nil, err
	}
	return needs, nil
}

// FindByID returns a need by its UUID, with its items and their materials.
func (r *repository) FindByID(id uuid.UUID) (*entities.Need, error) {
	var n entities.Need
	if err := r.db.Preload("Items.Material").First(&n, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

// Update applies editable fields to the need and returns the updated entity.
// Uses an updates map so falsey values persist correctly.
func (r *repository) Update(n *entities.Need) (*entities.Need, error) {
	if err := r.db.Model(&entities.Need{}).
		Where("id = ?", n.ID).
		Updates(map[string]interface{}{
			"requester_name": n.RequesterName,
			"buyer_name":     n.BuyerName,
			"cost_center":    n.CostCenter,
			"justification":  n.Justification,
			"required_date":  n.RequiredDate,
			"status":         n.Status,
			"promoted_at":    n.PromotedAt,
			"converted_at":   n.ConvertedAt,
			"is_active":      n.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return r.FindByID(n.ID)
}

// UpdateStatus updates only the status of a need and returns the refreshed entity.
func (r *repository) UpdateStatus(id uuid.UUID, status entities.NeedStatus) (*entities.Need, error) {
	if err := r.db.Model(&entities.Need{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

// AddItem inserts a new item line for a need.
func (r *repository) AddItem(item *entities.NeedItem) (*entities.NeedItem, error) {
	if err := r.db.Create(item).Error; err != nil {
		return nil, err
	}
	return r.FindItemByID(item.ID)
}

// ListItems returns all items of a need, with their materials preloaded.
func (r *repository) ListItems(needID uuid.UUID) ([]entities.NeedItem, error) {
	var items []entities.NeedItem
	if err := r.db.Preload("Material").Where("need_id = ?", needID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindItemByID returns a single item by its UUID, with its material preloaded.
func (r *repository) FindItemByID(itemID uuid.UUID) (*entities.NeedItem, error) {
	var item entities.NeedItem
	if err := r.db.Preload("Material").First(&item, "id = ?", itemID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateItem applies editable fields to an item and returns the updated entity.
// Optional references (material_id, selected_cost_id) are only overwritten when
// provided, so omitting them in a partial update preserves the current values.
func (r *repository) UpdateItem(item *entities.NeedItem) (*entities.NeedItem, error) {
	existing, err := r.FindItemByID(item.ID)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"quantity": item.Quantity,
		"unit":     item.Unit,
		"notes":    item.Notes,
	}
	if item.MaterialID != uuid.Nil {
		updates["material_id"] = item.MaterialID
	}
	if item.SelectedCostID != nil {
		updates["selected_cost_id"] = item.SelectedCostID
	}

	if err := r.db.Model(&entities.NeedItem{}).
		Where("id = ?", existing.ID).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.FindItemByID(item.ID)
}

// RemoveItem physically deletes an item line from a need.
func (r *repository) RemoveItem(itemID uuid.UUID) error {
	return r.db.Delete(&entities.NeedItem{}, "id = ?", itemID).Error
}
