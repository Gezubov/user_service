package controller

import "github.com/Gezubov/user_service/internal/models"

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id int64) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
	GetAllUsers() ([]models.User, error)
}
