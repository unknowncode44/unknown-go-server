package presenter

import (
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// CreateUserRequest representa el cuerpo esperado para crear un usuario.
// Las validaciones adicionales (rol válido, longitud de password, email único)
// se aplican en la capa de servicio.
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// UpdateUserRequest define el cuerpo esperado para actualizar un usuario.
// El ID viene por path, no en el body. No admite password ni email:
// el password se cambia por /users/:id/password y el email no se edita.
type UpdateUserRequest struct {
	Name string `json:"name,omitempty"`
	Role string `json:"role,omitempty"`
}

// SetPasswordRequest define el cuerpo esperado para resetear el password
// de un usuario (operación de Admin).
type SetPasswordRequest struct {
	Password string `json:"password"`
}

// UserResponse representa un usuario en las respuestas HTTP.
// Nunca incluye el hash de password.
type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// UserSuccessResponse representa una respuesta exitosa con uno o más usuarios.
type UserSuccessResponse struct {
	Success bool        `json:"ok"`
	Data    interface{} `json:"data"`
}

// ToUserResponse transforma una entidad User en UserResponse.
func ToUserResponse(u *entities.User) *UserResponse {
	if u == nil {
		return nil
	}

	return &UserResponse{
		ID:       u.ID.String(),
		Name:     u.Name,
		Email:    u.Email,
		Role:     string(u.Role),
		IsActive: u.IsActive,
	}
}

// ToUserListResponse transforma una lista de entidades User.
func ToUserListResponse(users []entities.User) []UserResponse {
	response := make([]UserResponse, 0, len(users))

	for _, u := range users {
		response = append(response, UserResponse{
			ID:       u.ID.String(),
			Name:     u.Name,
			Email:    u.Email,
			Role:     string(u.Role),
			IsActive: u.IsActive,
		})
	}

	return response
}
