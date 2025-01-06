package presenter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Company es el objeto de presentación que será pasado en las respuestas del controlador
type Company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	CUIT uint64 `json:"cuit"`
}

// CompanySuccessResponse devuelve una respuesta de éxito singular
func CompanySuccessResponse(data *entities.Company) *fiber.Map {
	company := Company{
		ID:   data.ID,
		Name: data.Name,
		CUIT: data.CUIT,
	}
	return &fiber.Map{
		"status": true,
		"data":   company,
		"error":  nil,
	}
}

// CompaniesSuccessResponse devuelve una respuesta de éxito con una lista de compañías
func CompaniesSuccessResponse(data *[]entities.Company) *fiber.Map {
	var companies []Company
	for _, c := range *data {
		companies = append(companies, Company{
			ID:   c.ID,
			Name: c.Name,
			CUIT: c.CUIT,
		})
	}

	return &fiber.Map{
		"status": true,
		"data":   companies,
		"error":  nil,
	}
}

// CompanyErrorResponse devuelve una respuesta de error
func CompanyErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   nil,
		"error":  err.Error(),
	}
}
