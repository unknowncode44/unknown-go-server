// Package auth provee la generación/validación de JWT y los middlewares de
// Fiber para proteger rutas por autenticación y por rol. No tiene entidad ni
// repository propios: se apoya en la entidad User y en la config global.
package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/config"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

// Claims son los datos que viajan firmados dentro del JWT.
// El rol viaja en el token: el middleware no vuelve a consultar la DB en cada
// request (trade-off aceptado: un cambio de rol no aplica hasta el próximo login
// o hasta que expire el token vigente).
type Claims struct {
	UserID uuid.UUID         `json:"user_id"`
	Email  string            `json:"email"`
	Role   entities.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken firma un JWT HS256 para el usuario dado.
func GenerateToken(user *entities.User) (string, error) {
	cfg := config.GetConfig().Auth
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.JWTExpiryHours) * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ParseToken valida la firma y expiración y devuelve los claims.
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig().Auth
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token inválido o expirado")
	}
	return claims, nil
}

// RequireAuth exige un Bearer token válido y deja los claims en fiber.Locals
// bajo las keys "user_id", "user_email", "user_role".
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": "falta el header Authorization"})
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := ParseToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)
		return c.Next()
	}
}

// RequireRole exige que el rol ya inyectado por RequireAuth esté en la lista dada.
// Debe montarse siempre DESPUÉS de RequireAuth en la cadena de middlewares.
func RequireRole(roles ...entities.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("user_role").(entities.UserRole)
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "error": "no tenés permisos para esta operación"})
	}
}
