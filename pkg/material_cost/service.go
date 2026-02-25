package material_cost

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/currency"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/vendor"
	"github.com/unknowncode44/unknown-go-server/pkg/vendor_material"
	"gorm.io/gorm"
)

// VendorCostComparison represents the comparison view for a vendor offering.
type VendorCostComparison struct {
	VendorID   uuid.UUID
	VendorName string
	Currency   string
	Cost       float64
	CostDate   time.Time
}

// Service defines business operations for MaterialCost.
type Service interface {
	Create(mc *entities.MaterialCost) (*entities.MaterialCost, error)
	FindAll() ([]entities.MaterialCost, error)
	FindByID(id uuid.UUID) (*entities.MaterialCost, error)
	FindByMaterial(materialID uuid.UUID) ([]entities.MaterialCost, error)
	FindLatestByMaterialAndVendor(materialID uuid.UUID, vendorID uuid.UUID) (*entities.MaterialCost, error)
	GetCurrentCost(materialID uuid.UUID, vendorID uuid.UUID, at *time.Time) (*entities.MaterialCost, error)
	CompareVendorsByMaterial(materialID uuid.UUID, at *time.Time) ([]VendorCostComparison, error)
}

type service struct {
	repo         Repository
	vmRepo       vendor_material.Repository
	vendorRepo   vendor.Repository
	currencyRepo currency.Repository
}

// NewService returns a new Service using the provided Repository and collaborators.
func NewService(repo Repository, vmRepo vendor_material.Repository, vendorRepo vendor.Repository, currencyRepo currency.Repository) Service {
	return &service{repo: repo, vmRepo: vmRepo, vendorRepo: vendorRepo, currencyRepo: currencyRepo}
}

func (s *service) Create(mc *entities.MaterialCost) (*entities.MaterialCost, error) {
	if mc == nil {
		return nil, errors.New("material cost is required")
	}

	if mc.MaterialID == uuid.Nil {
		return nil, errors.New("material_id is required")
	}
	if mc.VendorID == uuid.Nil {
		return nil, errors.New("vendor_id is required")
	}
	if mc.CurrencyID == uuid.Nil {
		return nil, errors.New("currency_id is required")
	}
	if mc.Cost <= 0 {
		return nil, errors.New("cost must be greater than zero")
	}
	if mc.CostDate.IsZero() || mc.CostDate.Equal(time.Time{}) {
		return nil, errors.New("cost_date is required")
	}

	return s.repo.Create(mc)
}

func (s *service) FindAll() ([]entities.MaterialCost, error) {
	return s.repo.FindAll()
}

func (s *service) FindByID(id uuid.UUID) (*entities.MaterialCost, error) {
	return s.repo.FindByID(id)
}

func (s *service) FindByMaterial(materialID uuid.UUID) ([]entities.MaterialCost, error) {
	return s.repo.FindByMaterial(materialID)
}

func (s *service) FindLatestByMaterialAndVendor(materialID uuid.UUID, vendorID uuid.UUID) (*entities.MaterialCost, error) {
	return s.repo.FindLatestByMaterialAndVendor(materialID, vendorID)
}

// GetCurrentCost returns the material cost valid at 'at' (or now if nil).
func (s *service) GetCurrentCost(materialID uuid.UUID, vendorID uuid.UUID, at *time.Time) (*entities.MaterialCost, error) {
	var when time.Time
	if at == nil {
		when = time.Now()
	} else {
		when = *at
	}

	mc, err := s.repo.FindCurrentCost(materialID, vendorID, when)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no current cost found for material/vendor at given date")
		}
		return nil, err
	}
	return mc, nil
}

// CompareVendorsByMaterial returns active vendors for a material with their current cost (omits vendors without cost).
func (s *service) CompareVendorsByMaterial(materialID uuid.UUID, at *time.Time) ([]VendorCostComparison, error) {
	var when time.Time
	if at == nil {
		when = time.Now()
	} else {
		when = *at
	}

	vms, err := s.vmRepo.FindActiveByMaterial(materialID)
	if err != nil {
		return nil, err
	}

	out := make([]VendorCostComparison, 0, len(vms))
	for _, vm := range vms {
		mc, err := s.GetCurrentCost(vm.MaterialID, vm.VendorID, &when)
		if err != nil {
			// skip vendors without current cost
			continue
		}

		v, err := s.vendorRepo.FindByID(vm.VendorID)
		if err != nil {
			continue
		}

		cur, err := s.currencyRepo.FindByID(mc.CurrencyID)
		if err != nil {
			continue
		}

		out = append(out, VendorCostComparison{
			VendorID:   vm.VendorID,
			VendorName: v.Name,
			Currency:   cur.Code,
			Cost:       mc.Cost,
			CostDate:   mc.CostDate,
		})
	}

	return out, nil
}
