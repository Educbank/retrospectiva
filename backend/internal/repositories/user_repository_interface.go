package repositories

import (
	"educ-retro/internal/models"

	"github.com/google/uuid"
)

// UserRepositoryInterface define a interface para o UserRepository
type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	GetUsersByIDs(ids []uuid.UUID) ([]models.User, error)
}
