package material

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Service define la API del módulo material. Permite a los handlers usar
// la lógica de negocio sin depender de detalles de persistencia.
type Service interface {
	Create(material *entities.Material) (*entities.Material, error)
	FindAll() ([]entities.Material, error)
	FindByID(id uuid.UUID) (*entities.Material, error)
	Update(material *entities.Material) (*entities.Material, error)
	Deactivate(id uuid.UUID) error
}

// service es la implementación por defecto de Service.
type service struct {
	repo Repository
}

// NewService crea una nueva instancia de Service con el repositorio inyectado.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Create aplica las validaciones de negocio y delega la persistencia al repositorio.
// Valida que los campos obligatorios no estén vacíos (trimmed).
func (s *service) Create(material *entities.Material) (*entities.Material, error) {
	if material == nil {
		return nil, errors.New("material es requerido")
	}

	// Trim y validaciones básicas
	if strings.TrimSpace(material.Name) == "" {
		return nil, errors.New("Nombre del material es requerido")
	}
	if strings.TrimSpace(material.Sector) == "" {
		return nil, errors.New("Rubro del material es requerido")
	}
	if strings.TrimSpace(material.UnitOfMeasure) == "" {
		return nil, errors.New("Unidad de medida es requerida")
	}

	// Definir estado inicial
	material.IsActive = true

	return s.repo.Create(material)
}

// FindAll devuelve todos los materiales.
func (s *service) FindAll() ([]entities.Material, error) {
	return s.repo.FindAll()
}

// FindByID devuelve un material por su UUID.
func (s *service) FindByID(id uuid.UUID) (*entities.Material, error) {
	return s.repo.FindByID(id)
}

// Update valida la entidad y delega la actualización al repositorio.
func (s *service) Update(material *entities.Material) (*entities.Material, error) {
	if material == nil || material.ID == uuid.Nil {
		return nil, errors.New("La id del material es requerida")
	}
	return s.repo.Update(material)
}

// Deactivate realiza un borrado lógico (cambia IsActive a false).
func (s *service) Deactivate(id uuid.UUID) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !m.IsActive {
		// Ya estaba desactivado, consideramos la operación idempotente
		return nil
	}
	m.IsActive = false
	_, err = s.repo.Update(m)
	return err
}
