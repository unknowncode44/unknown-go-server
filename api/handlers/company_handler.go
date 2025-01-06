package handlers

import (
	"net/http"

	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/company"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

// AddCompany maneja la creación de una nueva compañía
func AddCompany(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.Company
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.CompanyErrorResponse(err))
		}

		// Validar los campos requeridos
		if requestBody.Name == "" || requestBody.CUIT == 0 {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.CompanyErrorResponse(errors.New(
				"Por favor especifica nombre y CUIT")))
		}

		result, err := service.InsertCompany(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.CompanyErrorResponse(err))
		}

		return c.JSON(presenter.CompanySuccessResponse(result))
	}
}

// UpdateCompany maneja la actualización de una compañía existente
func UpdateCompany(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.Company
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.CompanyErrorResponse(err))
		}

		result, err := service.UpdateCompany(&requestBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.CompanyErrorResponse(err))
		}

		return c.JSON(presenter.CompanySuccessResponse(result))
	}
}

// RemoveCompany maneja la eliminación de una compañía
func RemoveCompany(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type DeleteRequest struct {
			ID string `json:"id"`
		}

		var requestBody DeleteRequest
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.CompanyErrorResponse(err))
		}

		err = service.RemoveCompany(requestBody.ID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.CompanyErrorResponse(err))
		}

		return c.JSON(&fiber.Map{
			"status": true,
			"data":   "Eliminado exitosamente",
			"err":    nil,
		})
	}
}

// GetCompanies maneja la obtención de todas las compañías
func GetCompanies(service company.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fetched, err := service.FetchCompanies()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.CompanyErrorResponse(err))
		}

		return c.JSON(presenter.CompaniesSuccessResponse(fetched))
	}
}
