package company

import "github.com/unknowncode44/unknown-go-server/pkg/entities"

// Service es una interfaz que permite a nuestro módulo API acceder al repositorio de Company
type Service interface {
	InsertCompany(company *entities.Company) (*entities.Company, error)
	FetchCompanies() (*[]entities.Company, error)
	UpdateCompany(company *entities.Company) (*entities.Company, error)
	RemoveCompany(ID string) error
}

type service struct {
	repository Repository
}

// NewService se usa para crear una instancia única del servicio
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

// InsertCompany es una capa de servicio que ayuda a insertar una nueva compañía
func (s *service) InsertCompany(company *entities.Company) (*entities.Company, error) {
	return s.repository.CreateCompany(company)
}

// FetchCompanies es una capa de servicio que ayuda a obtener todas las compañías
func (s *service) FetchCompanies() (*[]entities.Company, error) {
	return s.repository.ReadCompanies()
}

// UpdateCompany es una capa de servicio que ayuda a actualizar una compañía existente
func (s *service) UpdateCompany(company *entities.Company) (*entities.Company, error) {
	return s.repository.UpdateCompany(company)
}

// RemoveCompany es una capa de servicio que ayuda a eliminar una compañía
func (s *service) RemoveCompany(ID string) error {
	return s.repository.DeleteCompany(ID)
}
