package vendor

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Service define la API del módulo vendor. Permite a los handlers usar
// la lógica de negocio sin depender de detalles de persistencia.
type Service interface {
	Create(vendor *entities.Vendor) (*entities.Vendor, error)
	FindAll() ([]entities.Vendor, error)
	FindByID(id uuid.UUID) (*entities.Vendor, error)
	Update(vendor *entities.Vendor) (*entities.Vendor, error)
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
func (s *service) Create(vendor *entities.Vendor) (*entities.Vendor, error) {
	if vendor == nil {
		return nil, errors.New("Proveedor es requerido")
	}

	// Trim y validaciones basicas
	if strings.TrimSpace(vendor.Name) == "" {
		return nil, errors.New("Nombre del proveedor es requerido")
	}

	// Lo definimos como un vendor activo enseguida
	vendor.IsActive = true

	return s.repo.Create(vendor)
}

// FindAll devuelve todos los proveedores
func (s *service) FindAll() ([]entities.Vendor, error) {
	return s.repo.FindAll()
}

// FindByID devuelve un proveedor por su UUID.
func (s *service) FindByID(id uuid.UUID) (*entities.Vendor, error) {
	return s.repo.FindByID(id)
}

// Update valida la entidad y delega la actualización al repositorio.
func (s *service) Update(vendor *entities.Vendor) (*entities.Vendor, error) {

	// Validaciones basicas y Trim
	if vendor == nil || vendor.ID == uuid.Nil {
		return nil, errors.New("La id del proveedor es requerida")
	}
	if strings.TrimSpace(vendor.Name) == "" {
		return nil, errors.New("Nombre del proveedor es requerido")
	}

	return s.repo.Update(vendor)
}

// Deactivate realiza un borrado lógico (cambia IsActive a false).
func (s *service) Deactivate(id uuid.UUID) error {
	v, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !v.IsActive {
		// Ya estaba desactivado, consideramos la operación idempotente
		return nil
	}

	v.IsActive = false
	_, err = s.repo.Update(v)

	return err
}
