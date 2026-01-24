package asset_movement

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository defines persistence operations for AssetMovement entities.
type Repository interface {
	Create(am *entities.AssetMovement) (*entities.AssetMovement, error)
	CreateWithAssetUpdate(am *entities.AssetMovement) (*entities.AssetMovement, error)
	FindByAsset(assetID uuid.UUID) ([]entities.AssetMovement, error)
	FindLastByAsset(assetID uuid.UUID) (*entities.AssetMovement, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(am *entities.AssetMovement) (*entities.AssetMovement, error) {
	if err := r.db.Create(am).Error; err != nil {
		return nil, err
	}
	return am, nil
}

// CreateWithAssetUpdate performs insertion of AssetMovement and updates the related Asset
// atomically using a transaction. It enforces append-only behavior and first-movement rule.
func (r *repository) CreateWithAssetUpdate(am *entities.AssetMovement) (*entities.AssetMovement, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// check if asset exists and lock it for update
	var asset entities.Asset
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&asset, "id = ?", am.AssetID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// check if asset retired or inactive
	if !asset.IsActive || asset.Status == entities.AssetStatusRetired {
		tx.Rollback()
		return nil, gorm.ErrInvalidTransaction
	}

	// ensure first movement rule: if no movements exist, only INBOUND allowed
	var count int64
	if err := tx.Model(&entities.AssetMovement{}).Where("asset_id = ?", am.AssetID).Count(&count).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if count == 0 && am.Type != entities.AssetMovementInbound {
		tx.Rollback()
		return nil, gorm.ErrInvalidTransaction
	}

	// create movement
	if err := tx.Create(am).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// update asset status and current location depending on movement type
	switch am.Type {
	case entities.AssetMovementInbound:
		asset.Status = entities.AssetStatusInStock
		asset.CurrentLocationID = am.ToLocationID
	case entities.AssetMovementTransfer:
		asset.Status = entities.AssetStatusInStock
		asset.CurrentLocationID = am.ToLocationID
	case entities.AssetMovementInstall:
		asset.Status = entities.AssetStatusInstalled
		asset.CurrentLocationID = am.ToLocationID
	case entities.AssetMovementUninstall:
		asset.Status = entities.AssetStatusInStock
		asset.CurrentLocationID = am.ToLocationID
	case entities.AssetMovementRepair:
		asset.Status = entities.AssetStatusInRepair
		asset.CurrentLocationID = am.ToLocationID
	case entities.AssetMovementReturn:
		asset.Status = entities.AssetStatusInStock
		asset.CurrentLocationID = am.ToLocationID
	case entities.AssetMovementDecommission, entities.AssetMovementScrap:
		asset.Status = entities.AssetStatusRetired
		asset.IsActive = false
		asset.CurrentLocationID = am.ToLocationID
	}

	if err := tx.Model(&entities.Asset{}).Where("id = ?", asset.ID).Updates(map[string]interface{}{
		"status":              asset.Status,
		"current_location_id": asset.CurrentLocationID,
		"is_active":           asset.IsActive,
	}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return am, nil
}

func (r *repository) FindByAsset(assetID uuid.UUID) ([]entities.AssetMovement, error) {
	var list []entities.AssetMovement
	if err := r.db.Preload("FromLocation").Preload("ToLocation").Where("asset_id = ?", assetID).Order("movement_date desc, created_at desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) FindLastByAsset(assetID uuid.UUID) (*entities.AssetMovement, error) {
	var am entities.AssetMovement
	if err := r.db.Preload("FromLocation").Preload("ToLocation").Where("asset_id = ?", assetID).Order("movement_date desc, created_at desc").First(&am).Error; err != nil {
		return nil, err
	}
	return &am, nil
}
