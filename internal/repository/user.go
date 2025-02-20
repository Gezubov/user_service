package repository

import (
	"database/sql"

	"github.com/Gezubov/file_storage/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {

	return nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}

	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {

	return nil
}

func (r *UserRepository) Delete(id int64) error {

	return nil
}
