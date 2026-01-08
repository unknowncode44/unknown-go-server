package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
)

func VendorMaterialRoutes(router fiber.Router, handler *handlers.VendorMaterialHandler) {
	vms := router.Group("/vendor_materials")

	vms.Post("/", handler.Create)
	vms.Get("/", handler.GetAll)
	vms.Get(":id", handler.GetByID)
	vms.Put(":id", handler.Update)
	vms.Delete(":id", handler.Deactivate)
}
