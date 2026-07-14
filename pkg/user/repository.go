package user

import (
	"errors"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/gorm"
)

// Repository defines persistence operations for User entities.
// The service layer depends on this interface to remain persistence-agnostic.
type Repository interface {
	Create(user *entities.User) (*entities.User, error)
	FindAll() ([]entities.User, error)
	FindByID(id uuid.UUID) (*entities.User, error)
	// FindByEmail devuelve (nil, nil) si no existe un usuario con ese email,
	// para que el service pueda distinguir "no encontrado" de un error real
	// sin depender de GORM.
	FindByEmail(email string) (*entities.User, error)
	Update(user *entities.User) (*entities.User, error)
	UpdatePassword(id uuid.UUID, passwordHash string) error
	Delete(id uuid.UUID) error
	CountActiveByRole(role entities.UserRole) (int64, error)
}

// repository is the GORM-backed Repository implementation.
type repository struct {
	db *gorm.DB
}

// NewRepo returns a new Repository using the provided GORM DB.
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new user record into the database.
func (r *repository) Create(user *entities.User) (*entities.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindAll retrieves all user records.
func (r *repository) FindAll() ([]entities.User, error) {
	var users []entities.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FindByID returns a user by its UUID.
func (r *repository) FindByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail returns a user by email, or (nil, nil) if none exists.
func (r *repository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update applies changes to the user record and returns the updated entity.
// Uses an updates map so falsey values (e.g. IsActive=false) persist correctly.
// El email y el password NO se actualizan por esta vía.
func (r *repository) Update(user *entities.User) (*entities.User, error) {
	if err := r.db.
		Model(&entities.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"name":      user.Name,
			"role":      user.Role,
			"is_active": user.IsActive,
		}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// UpdatePassword actualiza únicamente el hash de password del usuario.
func (r *repository) UpdatePassword(id uuid.UUID, passwordHash string) error {
	return r.db.Model(&entities.User{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).Error
}

// Delete performs a logical delete by setting `is_active` to false.
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Model(&entities.User{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// CountActiveByRole cuenta los usuarios activos con el rol dado.
// El service lo usa para el guard de "no desactivar al último Admin".
func (r *repository) CountActiveByRole(role entities.UserRole) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.User{}).
		Where("role = ? AND is_active = ?", role, true).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
