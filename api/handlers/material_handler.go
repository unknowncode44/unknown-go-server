package handlers

// Las responsabilidades del handler seran
// # Validar el formato de la peticion / respuesta
// # Convertir las DTO en Entity
// # Llamar a la instancia de service
// # Convertir Entity en Response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
)

type MaterialHandler struct {
	service material.Service
}

func NewMaterialHandler(service material.Service) *MaterialHandler {
	return &MaterialHandler{service: service}
}

// Crear nuevo material (POST)

// los handler usan el contexto de la consulta que nos llega, lo pasamos usando fiber.Ctx
func (h *MaterialHandler) Create(c *fiber.Ctx) error {

	// nuestra variable req (de Request) deberia cumplir con el struct CreateMaterialRequest de nuestro modulo presenter
	var req presenter.CreateMaterialRequest

	// trataremos de parsear el body de la solicitud y en caso de que no sea posible
	// devolveremos un BadRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cuerpo de la consulta es invalido")
	}

	// en caso de que no exista error parseamos el DTO a una Entity
	materialEntity := &entities.Material{
		Name:          req.Name,
		Sector:        req.Sector,
		UnitOfMeasure: req.UnitOfMeasure,
	}

	// llamamos a nuestra instancia de servicio y a su metodo Create para crear un nuevo material
	created, err := h.service.Create(materialEntity)

	// si hay algun error por ahora devolvemos un BadRequest
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// retornamos un status Created y devolvemos la respuesta en JSON usando nuestra funcion auxiliar
	return c.Status(fiber.StatusCreated).JSON(toMaterialResponse(created))
}

// Listar todos los materiales (GET)
func (h *MaterialHandler) GetAll(c *fiber.Ctx) error {
	materials, err := h.service.FindAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var response []presenter.MaterialResponse
	for _, m := range materials {
		response = append(response, toMaterialResponse(&m))
	}

	return c.JSON(response)
}

// Obtener material por id (GET /:id)
func (h *MaterialHandler) GetById(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID de material es invalido")
	}

	materialEntity, err := h.service.FindByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "No se encontró el material")
	}

	return c.JSON(toMaterialResponse(materialEntity))
}

// Actualizar material (PUT)
func (h *MaterialHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID de material invalido")
	}

	var req presenter.UpdateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cuerpo de la peticion invalido")
	}

	materialEntity, err := h.service.FindByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "No se encontró el material")
	}

	// Aplicamos cambios
	materialEntity.Name = req.Name
	materialEntity.Sector = req.Sector
	materialEntity.UnitOfMeasure = req.UnitOfMeasure
	materialEntity.IsActive = req.IsActive

	updated, err := h.service.Update(materialEntity)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(toMaterialResponse(updated))
}

// Desactivar material (DELETE logico)
func (h *MaterialHandler) Deactivate(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID de material invalido")
	}

	if err := h.service.Deactivate(id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "No se encontró el material")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// funcion auxiliar que nos ayudara a transformar nuestro entidad en una response que
// cumpla con la structura del DTO MaterialResponse
func toMaterialResponse(m *entities.Material) presenter.MaterialResponse {
	return presenter.MaterialResponse{
		ID:            m.ID,
		Name:          m.Name,
		Sector:        m.Sector,
		UnitOfMeasure: m.UnitOfMeasure,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}
