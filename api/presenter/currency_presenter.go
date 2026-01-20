package presenter

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// CreateCurrencyRequest representa el cuerpo esperado para crear una moneda.
type CreateCurrencyRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// UpdateCurrencyRequest representa los campos permitidos para actualizar una moneda.
type UpdateCurrencyRequest struct {
	Name string `json:"name"`
}

// CurrencyResponse define la respuesta que devuelve la API para una moneda.
type CurrencyResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CurrencySuccessResponse representa una respuesta exitosa.
type CurrencySuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

// CurrencyErrorResponse representa una respuesta de error estándar.
type CurrencyErrorResponse struct {
	Success bool   `json:"ok"`
	Error   string `json:"error"`
}

// ToCurrencyResponse transforma una entidad Currency en CurrencyResponse.
func ToCurrencyResponse(c *entities.Currency) CurrencyResponse {
	if c == nil {
		return CurrencyResponse{}
	}
	return CurrencyResponse{
		ID:        c.ID.String(),
		Code:      c.Code,
		Name:      c.Name,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToCurrencyListResponse transforma una lista de entidades Currency.
func ToCurrencyListResponse(items []entities.Currency) []CurrencyResponse {
	resp := make([]CurrencyResponse, 0, len(items))
	for _, c := range items {
		resp = append(resp, CurrencyResponse{
			ID:        c.ID.String(),
			Code:      c.Code,
			Name:      c.Name,
			IsActive:  c.IsActive,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	return resp
}
