package material

import (
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Service es una interfaz que permite a nuestro módulo API acceder al repositorio de Material
type Service interface {
	Create(material *entities.Material) (*entities.Material, error)
	FindAll() ([]entities.Material, error)
	FindByID(id uuid.UUID) (*entities.Material, error)
	Update(material *entities.Material) (*entities.Material, error)
	Deactivate(id uuid.UUID) error
}

type service struct {
	repo Repository
}
