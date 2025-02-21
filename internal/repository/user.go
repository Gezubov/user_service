package repository

import (
	"database/sql"

	"log/slog"
	"time"

	"github.com/Gezubov/user_service/internal/models"
	"github.com/google/uuid"
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
		INSERT INTO users (uuid, username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid`

	now := time.Now()
	user.UUID = uuid.New()
	err := r.db.QueryRow(
		query,
		user.UUID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		now,
		now,
	).Scan(&user.UUID)

	if err != nil {
		slog.Error("Error creating user", "error", err)
		return err
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *UserRepository) GetByUUID(uuid uuid.UUID) (*models.User, error) {
	slog.Info("Getting user with UUID", "uuid", uuid)
	user := &models.User{}

	query := `
		SELECT uuid, username, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE uuid = $1`

	err := r.db.QueryRow(query, uuid).Scan(
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		slog.Warn("User not found", "uuid", uuid)
		return nil, ErrUserNotFound
	}
	if err != nil {
		slog.Error("Error fetching user", "uuid", uuid, "error", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	slog.Info("Getting user with email", "email", email)
	user := &models.User{}

	query := `SELECT uuid, username, email, password_hash, role, created_at, updated_at 
	FROM users 
	WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return user, err
}

func (r *UserRepository) Update(user *models.User) error {
	slog.Info("Updating user", "uuid", user.UUID)
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, role = $4, updated_at = $5
		WHERE uuid = $6`

	result, err := r.db.Exec(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		time.Now(),
		user.UUID,
	)
	if err != nil {
		slog.Error("Error updating user", "uuid", user.UUID, "error", err)

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Error checking update result", "uuid", user.UUID, "error", err)

		return err
	}
	if rowsAffected == 0 {
		slog.Warn("No rows updated", "uuid", user.UUID)
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(uuid uuid.UUID) error {
	slog.Info("Deleting user", "uuid", uuid)
	query := `DELETE FROM users WHERE uuid = $1`

	result, err := r.db.Exec(query, uuid)
	if err != nil {
		slog.Error("Error deleting user", "uuid", uuid, "error", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("Error checking delete result", "uuid", uuid, "error", err)
		return err
	}
	if rowsAffected == 0 {
		slog.Warn("User not found during deletion", "uuid", uuid)
		return ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	slog.Info("Fetching all users from database")

	query := `SELECT uuid, username, email, role, created_at, updated_at FROM users`
	rows, err := r.db.Query(query)

	if err != nil {
		slog.Error("Error executing query to fetch users", "error", err)
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UUID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
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

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	slog.Info("Getting user with username", "username", username)
	user := &models.User{}

	query := `SELECT uuid, username, email, password_hash, role, created_at, updated_at 
	FROM users 
	WHERE username = $1`

	err := r.db.QueryRow(query, username).Scan(
		&user.UUID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return user, err
}
