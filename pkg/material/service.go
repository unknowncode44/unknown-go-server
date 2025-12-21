package material

import (
	"errors"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Service es una interfaz que permite a nuestro módulo API acceder al repositorio de Material
type Service interface {
	Create(material *entities.Material) (*entities.Material, error)
	FindAll() ([]entities.Material, error)
	// FindByID(id uuid.UUID) (*entities.Material, error)
	Update(material *entities.Material) (*entities.Material, error)
	// Deactivate(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Logica para crear nuevo material
func (s *service) Create(material *entities.Material) (*entities.Material, error) {

	// las validaciones se realizan en el servicio y no el handler
	if material.Name == "" {
		return nil, errors.New("Nombre del Material es requerido")
	}
	if material.Sector == "" {
		return nil, errors.New("Rubro del Material es requerido")
	}
	if material.UnitOfMeasure == "" {
		return nil, errors.New("Unidad de medida es requerido")
	}

	// Definimos al material como activo
	material.IsActive = true

	return s.repo.Create(material)
}

// Logica para obtener todos los materiales
func (s *service) FindAll() ([]entities.Material, error) {
	return s.repo.FindAll()
}

// Logica para actualizar un material
func (s *service) Update(material *entities.Material) (*entities.Material, error) {

	// validaciones se aplica en el servicio y no el handler
	if material.ID == uuid.Nil {
		return nil, errors.New("La id del material es requerida")
	}
	return s.repo.Update(material)
}
