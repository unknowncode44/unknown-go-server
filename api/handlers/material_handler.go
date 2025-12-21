package handlers

// Las responsabilidades del handler seran
// # Validar el formato de la peticion / respuesta
// # Convertir las DTO en Entity
// # Llamar a la instancia de service
// # Convertir Entity en Response

import (
	"github.com/unknowncode44/unknown-go-server/pkg/material"
)

type MaterialHandler struct {
	service material.Service
}

func NewMaterialHandler(service material.Service) *MaterialHandler {
	return &MaterialHandler{service: service}
}
