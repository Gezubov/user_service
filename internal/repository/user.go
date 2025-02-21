package repository

import (
	"database/sql"

	"log/slog"
	"time"

	"github.com/Gezubov/user_service/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	slog.Info("Creating user", "username", user.Username)
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
		slog.Error("Error creating user", "error", err)
		return err
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	slog.Info("Getting user with ID", "id", id)
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
		slog.Warn("User not found", "id", id)
		return nil, ErrUserNotFound
	}
	if err != nil {
		slog.Error("Error fetching user", "id", id, "error", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	slog.Info("Updating user", "id", user.ID)
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
		slog.Error("Error updating user", "id", user.ID, "error", err)

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Error checking update result", "id", user.ID, "error", err)

		return err
	}
	if rowsAffected == 0 {
		slog.Warn("No rows updated", "id", user.ID)
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(id int64) error {
	slog.Info("Deleting user", "id", id)
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		slog.Error("Error deleting user", "id", id, "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Error checking delete result", "id", id, "error", err)
		return err
	}
	if rowsAffected == 0 {
		slog.Warn("User not found during deletion", "id", id)
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	slog.Info("Fetching all users from database")

	query := `SELECT id, username, email, created_at, updated_at FROM users`
	rows, err := r.db.Query(query)

	if err != nil {
		slog.Error("Error executing query to fetch users", "error", err)
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			slog.Error("Error scanning user row", "error", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating over user rows", "error", err)
		return nil, err
	}

	slog.Info("Successfully fetched users", "count", len(users))
	return users, nil
}
