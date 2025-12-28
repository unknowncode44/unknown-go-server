package vendor

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository define las operaciones de persistencia para Vendor.
// Se usa por la capa de servicio para abstraer detalles de GORM.
type Repository interface {
	Create(vendor *entities.Vendor) (*entities.Vendor, error)
	FindAll() ([]entities.Vendor, error)
	FindByID(id uuid.UUID) (*entities.Vendor, error)
	Update(vendor *entities.Vendor) (*entities.Vendor, error)
	Delete(id uuid.UUID) error
}

// repository es la implementación basada en GORM de Repository.
type repository struct {
	db *gorm.DB
}

// NewRepo crea una instancia del repositorio de Vendor con la DB proporcionada.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserta un nuevo registro de vendor en la base de datos.
func (r *repository) Create(vendor *entities.Vendor) (*entities.Vendor, error) {
	if err := r.db.Create(vendor).Error; err != nil {
		return nil, err
	}
	return vendor, nil
}

// FindAll recupera todos los registros de vendor.
func (r *repository) FindAll() ([]entities.Vendor, error) {
	var vendors []entities.Vendor
	if err := r.db.Find(&vendors).Error; err != nil {
		return nil, err
	}

	return vendors, nil
}

// FindByID busca un material por su UUID.
func (r *repository) FindByID(id uuid.UUID) (*entities.Vendor, error) {
	var vendor entities.Vendor
	if err := r.db.First(&vendor, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &vendor, nil
}

// Update aplica los cambios de la entidad en la base y devuelve la entidad actualizada.
func (r *repository) Update(vendor *entities.Vendor) (*entities.Vendor, error) {
	if err := r.db.Model(vendor).Updates(vendor).Error; err != nil {
		return nil, err
	}

	return vendor, nil
}

// Delete elimina el registro identificado por el UUID.
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Model(&entities.Vendor{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
