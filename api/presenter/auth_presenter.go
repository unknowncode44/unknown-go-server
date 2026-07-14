package presenter

// LoginRequest representa el cuerpo esperado para iniciar sesión.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse representa la respuesta de un login exitoso: el JWT firmado
// y los datos públicos del usuario autenticado.
type LoginResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}
