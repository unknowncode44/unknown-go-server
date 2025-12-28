package material

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository define las operaciones de persistencia para Material.
// Se usa por la capa de servicio para abstraer detalles de GORM.
type Repository interface {
	Create(material *entities.Material) (*entities.Material, error)
	FindAll() ([]entities.Material, error)
	FindByID(id uuid.UUID) (*entities.Material, error)
	Update(material *entities.Material) (*entities.Material, error)
	Delete(id uuid.UUID) error
}

// repository es la implementación basada en GORM de Repository.
type repository struct {
	db *gorm.DB
}

// NewRepo crea una instancia del repositorio de Material con la DB proporcionada.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserta un nuevo registro de material en la base de datos.
func (r *repository) Create(material *entities.Material) (*entities.Material, error) {
	if err := r.db.
		Create(material).Error; err != nil {
		return nil, err
	}
	return material, nil
}

// FindAll recupera todos los registros de materiales.
func (r *repository) FindAll() ([]entities.Material, error) {
	var materials []entities.Material
	if err := r.db.
		Find(&materials).Error; err != nil {
		return nil, err
	}
	return materials, nil
}

// FindByID busca un material por su UUID.
func (r *repository) FindByID(id uuid.UUID) (*entities.Material, error) {
	var material entities.Material
	if err := r.db.
		First(&material, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &material, nil
}

// Update aplica los cambios de la entidad en la base y devuelve la entidad actualizada.
func (r *repository) Update(material *entities.Material) (*entities.Material, error) {
	if err := r.db.
		Model(&entities.Material{}).
		Where("id = ?", material.ID).
		Updates(material).Error; err != nil {
		return nil, err
	}
	return material, nil
}

// Delete elimina el registro identificado por el UUID.
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Model(&entities.Material{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}
