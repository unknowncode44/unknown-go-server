package company

import (
	"errors"
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository es la interfaz para realizar operaciones CRUD sobre la entidad Company
type Repository interface {
	CreateCompany(company *entities.Company) (*entities.Company, error)
	ReadCompanies() (*[]entities.Company, error)
	UpdateCompany(company *entities.Company) (*entities.Company, error)
	DeleteCompany(ID string) error
}

type repository struct {
	db *gorm.DB
}

// NewRepo crea una instancia única del repositorio de Company
func NewRepo(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// CreateCompany agrega una nueva compañía a la base de datos
func (r *repository) CreateCompany(company *entities.Company) (*entities.Company, error) {
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()
	if err := r.db.Create(&company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

// ReadCompanies obtiene todas las compañías de la base de datos
func (r *repository) ReadCompanies() (*[]entities.Company, error) {
	var companies []entities.Company
	if err := r.db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return &companies, nil
}

// UpdateCompany actualiza una compañía existente en la base de datos
func (r *repository) UpdateCompany(company *entities.Company) (*entities.Company, error) {
	company.UpdatedAt = time.Now()
	if err := r.db.Save(&company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

// DeleteCompany elimina una compañía de la base de datos
func (r *repository) DeleteCompany(ID string) error {
	result := r.db.Delete(&entities.Company{}, "id = ?", ID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("company not found")
	}
	return nil
}
