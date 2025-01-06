package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
	"github.com/unknowncode44/unknown-go-server/pkg/company"
)

func CompanyRouter(app fiber.Router, service company.Service) {
	app.Get("/companies", handlers.GetCompanies(service))
	app.Post("/companies", handlers.AddCompany(service))
	app.Put("/companies", handlers.UpdateCompany(service))
	app.Delete("/companies", handlers.RemoveCompany(service))
}
