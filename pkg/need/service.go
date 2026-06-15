package need

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
)

// Service defines business operations for Need entities.
type Service interface {
	Create(n *entities.Need) (*entities.Need, error)
	FindAll() ([]entities.Need, error)
	FindByID(id uuid.UUID) (*entities.Need, error)
	Update(n *entities.Need) (*entities.Need, error)

	// Controlled state transitions
	Promote(id uuid.UUID) (*entities.Need, error)
	MarkConverted(id uuid.UUID) (*entities.Need, error)
	Discard(id uuid.UUID) (*entities.Need, error)

	// Items
	AddItem(needID uuid.UUID, item *entities.NeedItem) (*entities.NeedItem, error)
	ListItems(needID uuid.UUID) ([]entities.NeedItem, error)
	UpdateItem(item *entities.NeedItem) (*entities.NeedItem, error)
	RemoveItem(itemID uuid.UUID) error
}

type service struct {
	repo         Repository
	materialRepo material.Repository
}

// NewService returns a new Need Service. It depends on the material repository
// to validate referenced materials.
func NewService(r Repository, mr material.Repository) Service {
	return &service{repo: r, materialRepo: mr}
}

// Create validates the required fields, normalizes strings and persists a new
// need with an auto-generated correlative number.
func (s *service) Create(n *entities.Need) (*entities.Need, error) {
	if n == nil {
		return nil, errors.New("need is required")
	}

	n.RequesterName = strings.TrimSpace(n.RequesterName)
	n.BuyerName = strings.TrimSpace(n.BuyerName)
	n.CostCenter = strings.TrimSpace(n.CostCenter)
	n.Justification = strings.TrimSpace(n.Justification)

	if n.RequesterName == "" {
		return nil, errors.New("requester_name is required")
	}
	if n.BuyerName == "" {
		return nil, errors.New("buyer_name is required")
	}

	n.Status = entities.NeedStatusInProgress
	n.IsActive = true

	return s.repo.CreateWithNumber(n)
}

func (s *service) FindAll() ([]entities.Need, error) {
	return s.repo.FindAll()
}

func (s *service) FindByID(id uuid.UUID) (*entities.Need, error) {
	if id == uuid.Nil {
		return nil, errors.New("need id is required")
	}
	return s.repo.FindByID(id)
}

// Update applies the editable fields of an existing need.
func (s *service) Update(n *entities.Need) (*entities.Need, error) {
	if n == nil || n.ID == uuid.Nil {
		return nil, errors.New("need id is required")
	}

	n.RequesterName = strings.TrimSpace(n.RequesterName)
	n.BuyerName = strings.TrimSpace(n.BuyerName)
	n.CostCenter = strings.TrimSpace(n.CostCenter)
	n.Justification = strings.TrimSpace(n.Justification)

	if n.RequesterName == "" {
		return nil, errors.New("requester_name is required")
	}
	if n.BuyerName == "" {
		return nil, errors.New("buyer_name is required")
	}

	return s.repo.Update(n)
}

// Promote moves a need from IN_PROGRESS to READY. It requires at least one
// item and stamps the promotion time.
func (s *service) Promote(id uuid.UUID) (*entities.Need, error) {
	n, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if n.Status != entities.NeedStatusInProgress {
		return nil, errors.New("only needs in IN_PROGRESS can be promoted")
	}

	items, err := s.repo.ListItems(id)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("a need must have at least one item to be promoted")
	}

	now := time.Now()
	n.PromotedAt = &now
	n.Status = entities.NeedStatusReady

	return s.repo.Update(n)
}

// MarkConverted moves a need from READY to CONVERTED once the external system
// confirms the Purchase Requisition was created.
func (s *service) MarkConverted(id uuid.UUID) (*entities.Need, error) {
	n, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if n.Status != entities.NeedStatusReady {
		return nil, errors.New("only needs in READY can be marked as converted")
	}

	now := time.Now()
	n.ConvertedAt = &now
	n.Status = entities.NeedStatusConverted

	return s.repo.Update(n)
}

// Discard deactivates a need from any state except CONVERTED.
func (s *service) Discard(id uuid.UUID) (*entities.Need, error) {
	n, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if n.Status == entities.NeedStatusConverted {
		return nil, errors.New("a converted need cannot be discarded")
	}

	n.IsActive = false
	n.Status = entities.NeedStatusDiscarded

	return s.repo.Update(n)
}

// AddItem validates the item, the referenced material and the parent need's
// state before appending the line.
func (s *service) AddItem(needID uuid.UUID, item *entities.NeedItem) (*entities.NeedItem, error) {
	if item == nil {
		return nil, errors.New("item is required")
	}
	if item.MaterialID == uuid.Nil {
		return nil, errors.New("material_id is required")
	}
	if item.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	item.Unit = strings.TrimSpace(item.Unit)
	if item.Unit == "" {
		return nil, errors.New("unit is required")
	}

	// the parent need must exist and still be editable
	n, err := s.repo.FindByID(needID)
	if err != nil {
		return nil, err
	}
	if n.Status != entities.NeedStatusInProgress {
		return nil, errors.New("items can only be added to a need in IN_PROGRESS")
	}

	// the referenced material must exist and be active
	mat, err := s.materialRepo.FindByID(item.MaterialID)
	if err != nil {
		return nil, errors.New("material not found")
	}
	if !mat.IsActive {
		return nil, errors.New("material is inactive")
	}

	item.NeedID = needID
	return s.repo.AddItem(item)
}

func (s *service) ListItems(needID uuid.UUID) ([]entities.NeedItem, error) {
	if needID == uuid.Nil {
		return nil, errors.New("need id is required")
	}
	return s.repo.ListItems(needID)
}

// UpdateItem validates and persists changes to an existing item line.
func (s *service) UpdateItem(item *entities.NeedItem) (*entities.NeedItem, error) {
	if item == nil || item.ID == uuid.Nil {
		return nil, errors.New("item id is required")
	}
	if item.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}
	item.Unit = strings.TrimSpace(item.Unit)
	if item.Unit == "" {
		return nil, errors.New("unit is required")
	}

	if item.MaterialID != uuid.Nil {
		mat, err := s.materialRepo.FindByID(item.MaterialID)
		if err != nil {
			return nil, errors.New("material not found")
		}
		if !mat.IsActive {
			return nil, errors.New("material is inactive")
		}
	}

	return s.repo.UpdateItem(item)
}

func (s *service) RemoveItem(itemID uuid.UUID) error {
	if itemID == uuid.Nil {
		return errors.New("item id is required")
	}
	return s.repo.RemoveItem(itemID)
}
