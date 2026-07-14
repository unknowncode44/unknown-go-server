package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"golang.org/x/crypto/bcrypt"
)

// Package user provides the domain service and repository contracts
// and default implementations for system user accounts (auth + roles).

// ErrLastActiveAdmin se devuelve al intentar desactivar al último ADMIN activo.
var ErrLastActiveAdmin = errors.New("no se puede desactivar al último administrador activo")

// Service defines the business API for user operations. Handlers use
// this interface to perform domain operations without depending on
// persistence details.
type Service interface {
	Create(user *entities.User, plainPassword string) (*entities.User, error)
	FindAll() ([]entities.User, error)
	FindByID(id uuid.UUID) (*entities.User, error)
	Update(user *entities.User) (*entities.User, error)
	SetPassword(id uuid.UUID, plainPassword string) error
	Deactivate(id uuid.UUID) error
	Authenticate(email, plainPassword string) (*entities.User, error)
}

// service is the default Service implementation.
type service struct {
	repo Repository
}

// NewService returns a new Service using the provided Repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// isValidRole verifica que el rol sea uno de los 3 valores permitidos.
func isValidRole(role entities.UserRole) bool {
	switch role {
	case entities.UserRoleAdmin, entities.UserRoleUser, entities.UserRoleOperativeUser:
		return true
	}
	return false
}

// Create valida las reglas de negocio, hashea el password con bcrypt y
// delega la persistencia al repository.
func (s *service) Create(user *entities.User, plainPassword string) (*entities.User, error) {
	if user == nil {
		return nil, errors.New("user is required")
	}
	if strings.TrimSpace(user.Name) == "" {
		return nil, errors.New("user name is required")
	}
	if strings.TrimSpace(user.Email) == "" {
		return nil, errors.New("user email is required")
	}
	if !isValidRole(user.Role) {
		return nil, errors.New("role inválido: debe ser ADMIN, USER u OPERATIVE_USER")
	}
	if len(plainPassword) < 8 {
		return nil, errors.New("el password debe tener al menos 8 caracteres")
	}

	// el email debe ser único
	existing, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("ya existe un usuario con ese email")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = string(hash)
	user.IsActive = true
	return s.repo.Create(user)
}

// FindAll returns all users.
func (s *service) FindAll() ([]entities.User, error) {
	return s.repo.FindAll()
}

// FindByID returns a user by UUID.
func (s *service) FindByID(id uuid.UUID) (*entities.User, error) {
	return s.repo.FindByID(id)
}

// Update valida la entidad y delega el update al repository.
// No toca password ni email (el email no se edita en este alcance:
// si hace falta, se da de baja el usuario y se crea uno nuevo).
func (s *service) Update(user *entities.User) (*entities.User, error) {
	if user == nil || user.ID == uuid.Nil {
		return nil, errors.New("user id is required")
	}
	if strings.TrimSpace(user.Name) == "" {
		return nil, errors.New("user name is required")
	}
	if !isValidRole(user.Role) {
		return nil, errors.New("role inválido: debe ser ADMIN, USER u OPERATIVE_USER")
	}

	return s.repo.Update(user)
}

// SetPassword valida la longitud mínima, hashea y actualiza solo el password.
func (s *service) SetPassword(id uuid.UUID, plainPassword string) error {
	if len(plainPassword) < 8 {
		return errors.New("el password debe tener al menos 8 caracteres")
	}
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(id, string(hash))
}

// Deactivate performs a logical delete by setting IsActive to false.
// The operation is idempotent. Si el usuario es ADMIN y es el último admin
// activo, se rechaza con ErrLastActiveAdmin sin aplicar el cambio.
func (s *service) Deactivate(id uuid.UUID) error {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !u.IsActive {
		return nil
	}

	if u.Role == entities.UserRoleAdmin {
		count, err := s.repo.CountActiveByRole(entities.UserRoleAdmin)
		if err != nil {
			return err
		}
		if count <= 1 {
			return ErrLastActiveAdmin
		}
	}

	u.IsActive = false
	_, err = s.repo.Update(u)
	return err
}

// Authenticate busca el usuario por email y compara el password con bcrypt.
// Devuelve siempre el mismo error genérico ante email inexistente, usuario
// inactivo o password incorrecta, para evitar enumeración de usuarios.
func (s *service) Authenticate(email, plainPassword string) (*entities.User, error) {
	invalidCredentials := errors.New("credenciales inválidas")

	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if u == nil || !u.IsActive {
		return nil, invalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword)); err != nil {
		return nil, invalidCredentials
	}
	return u, nil
}

// SeedInitialAdmin crea el primer usuario ADMIN si no existe ninguno activo.
// Es idempotente: si ya hay un admin (o no hay credenciales configuradas) no
// hace nada, así que puede llamarse en cada arranque sin riesgo.
func SeedInitialAdmin(repo Repository, email, plainPassword string) error {
	if email == "" || plainPassword == "" {
		return nil // no hay nada configurado, no-op
	}
	count, err := repo.CountActiveByRole(entities.UserRoleAdmin)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // ya existe al menos un admin, no hacer nada
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = repo.Create(&entities.User{
		Name:         "Admin",
		Email:        email,
		PasswordHash: string(hash),
		Role:         entities.UserRoleAdmin,
		IsActive:     true,
	})
	if err == nil {
		fmt.Printf("Admin inicial creado: %s\n", email)
	}
	return err
}
