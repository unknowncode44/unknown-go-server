package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserRole define los roles posibles de un usuario del sistema.
type UserRole string

const (
	UserRoleAdmin         UserRole = "ADMIN"          // gestiona depósitos/sectores/estanterías y usuarios
	UserRoleUser          UserRole = "USER"           // gestiona materiales / Bienes de Cambio
	UserRoleOperativeUser UserRole = "OPERATIVE_USER" // gestiona Activos / Bienes de Uso (BDU)
)

// User representa una cuenta con acceso al sistema (rol Usuario Público excluido:
// ese caso se resuelve como request sin token contra /public).
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Role         UserRole  `gorm:"type:varchar(20);not null;default:'USER'"`
	IsActive     bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
