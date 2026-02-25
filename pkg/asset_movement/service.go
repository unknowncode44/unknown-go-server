package asset_movement

import (
	"errors"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Service defines business operations for AssetMovement.
type Service interface {
	Create(am *entities.AssetMovement) (*entities.AssetMovement, error)
	FindByAsset(assetID uuid.UUID) ([]entities.AssetMovement, error)
	FindLastByAsset(assetID uuid.UUID) (*entities.AssetMovement, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(am *entities.AssetMovement) (*entities.AssetMovement, error) {
	if am == nil {
		return nil, errors.New("asset movement is required")
	}
	if am.AssetID == uuid.Nil {
		return nil, errors.New("asset_id is required")
	}
	if am.Type == "" {
		return nil, errors.New("movement type is required")
	}
	if am.MovementDate.IsZero() {
		return nil, errors.New("movement_date is required")
	}

	created, err := s.repo.CreateWithAssetUpdate(am)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *service) FindByAsset(assetID uuid.UUID) ([]entities.AssetMovement, error) {
	return s.repo.FindByAsset(assetID)
}

func (s *service) FindLastByAsset(assetID uuid.UUID) (*entities.AssetMovement, error) {
	return s.repo.FindLastByAsset(assetID)
}
