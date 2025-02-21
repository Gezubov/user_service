package service

import "github.com/Gezubov/user_service/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int64) (*models.User, error)
	Update(user *models.User) error
	Delete(id int64) error
	GetAllUsers() ([]models.User, error)
}
