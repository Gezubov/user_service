package repository

import (
	"database/sql"
	"time"

	"log"

	"github.com/Gezubov/user_service/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	log.Printf("Creating user: %s", user.Username)
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.Password,
		now,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	log.Printf("Getting user with ID: %d", id)
	user := &models.User{}
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	log.Printf("Updating user with ID: %d", user.ID)
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5`

	result, err := r.db.Exec(
		query,
		user.Username,
		user.Email,
		user.Password,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(id int64) error {
	log.Printf("Deleting user with ID: %d", id)
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	log.Println("Getting all users")
	query := `SELECT id, username, email, created_at, updated_at FROM users`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
