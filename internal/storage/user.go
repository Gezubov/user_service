package storage

import (
	"context"
	"database/sql"

	"log/slog"
	"time"

	"github.com/Gezubov/user_service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type UserStorage struct {
	db  *pgx.Conn
	ctx context.Context
}

func NewUserStorage(ctx context.Context, db *pgx.Conn) *UserStorage {
	return &UserStorage{ctx: ctx, db: db}
}

func (r *UserStorage) Create(ctx context.Context, user *models.User) error {
	slog.Info("Creating user", "username", user.Username)
	query := `
		INSERT INTO users (uuid, username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid`

	now := time.Now()
	user.UUID = uuid.New()
	err := r.db.QueryRow(ctx,
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

func (r *UserStorage) GetByUUID(ctx context.Context, uuid uuid.UUID) (*models.User, error) {
	slog.Info("Getting user with UUID", "uuid", uuid)
	user := &models.User{}

	query := `
		SELECT uuid, username, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE uuid = $1`

	err := r.db.QueryRow(ctx, query, uuid).Scan(
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

func (r *UserStorage) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	slog.Info("Getting user with email", "email", email)
	user := &models.User{}

	query := `SELECT uuid, username, email, password_hash, role, created_at, updated_at 
	FROM users 
	WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
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

func (r *UserStorage) Update(ctx context.Context, user *models.User) error {
	slog.Info("Updating user", "uuid", user.UUID)
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, role = $4, updated_at = $5
		WHERE uuid = $6`

	result, err := r.db.Exec(
		ctx,
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

	rowsAffected := result.RowsAffected()
	// if err != nil {
	// 	slog.Error("Error checking update result", "uuid", user.UUID, "error", err)

	// 	return err
	// }
	if rowsAffected == 0 {
		slog.Warn("No rows updated", "uuid", user.UUID)
		return ErrUserNotFound
	}

	return nil
}

func (r *UserStorage) Delete(ctx context.Context, uuid uuid.UUID) error {
	slog.Info("Deleting user", "uuid", uuid)
	query := `DELETE FROM users WHERE uuid = $1`

	result, err := r.db.Exec(ctx, query, uuid)
	if err != nil {
		slog.Error("Error deleting user", "uuid", uuid, "error", err)
		return err
	}

	rowsAffected := result.RowsAffected()
	// if err != nil {
	// 	slog.Error("Error checking delete result", "uuid", uuid, "error", err)
	// 	return err
	// }
	if rowsAffected == 0 {
		slog.Warn("User not found during deletion", "uuid", uuid)
		return ErrUserNotFound
	}

	return nil
}

func (r *UserStorage) GetAllUsers(ctx context.Context) ([]models.User, error) {
	slog.Info("Fetching all users from database")

	query := `SELECT uuid, username, email, role, created_at, updated_at FROM users`
	rows, err := r.db.Query(ctx, query)

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

func (r *UserStorage) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	slog.Info("Getting user with username", "username", username)
	user := &models.User{}

	query := `SELECT uuid, username, email, password_hash, role, created_at, updated_at 
	FROM users 
	WHERE username = $1`

	err := r.db.QueryRow(ctx, query, username).Scan(
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
