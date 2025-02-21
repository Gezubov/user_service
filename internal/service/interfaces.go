package service

import (
	"github.com/Gezubov/user_service/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByUUID(uuid uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(uuid uuid.UUID) error
	GetAllUsers() ([]models.User, error)
}
