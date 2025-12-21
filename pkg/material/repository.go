package material

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository es la interfaz para realizar operaciones CRUD sobre la entidad Material
type Repository interface {
	Create(material *entities.Material) (*entities.Material, error)
	FindAll() ([]entities.Material, error)
	FindByID(id uuid.UUID) (*entities.Material, error)
	Update(material *entities.Material) (*entities.Material, error)
	Delete(id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

// NewRepo crea una instancia única del repositorio de Material
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Implementacion de metodos

// Crear nuevo material
func (r *repository) Create(material *entities.Material) (*entities.Material, error) {
	if err := r.db.Create(material).Error; err != nil {
		return nil, err
	}
	return material, nil
}

// Leer todos los materiales
func (r *repository) FindAll() ([]entities.Material, error) {
	var materials []entities.Material
	if err := r.db.Find(&materials).Error; err != nil {
		return nil, err
	}
	return materials, nil
}

// Obtener material por id
func (r *repository) FindByID(id uuid.UUID) (*entities.Material, error) {
	var material entities.Material
	if err := r.db.First(&material, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &material, nil
}

// Actualizar material
func (r *repository) Update(material *entities.Material) (*entities.Material, error) {
	if err := r.db.Model(material).Updates(material).Error; err != nil {
		return nil, err
	}
	return material, nil
}

// Eliminar material
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.Material{}, "id = ?", id).Error
}
