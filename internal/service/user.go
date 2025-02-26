package service

import (
	"context"
	"errors"
	"time"

	"github.com/Gezubov/user_service/config"
	"github.com/Gezubov/user_service/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type UserStorage interface {
	Create(ctx context.Context, user *models.User) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, uuid uuid.UUID) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type UserService struct {
	userRepo UserStorage
	ctx      context.Context
}

func NewUserService(ctx context.Context, userRepo UserStorage) *UserService {
	return &UserService{
		ctx:      ctx,
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User, password string) error {
	existingUserByUsername, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUserByUsername != nil {
		return errors.New("имя пользователя уже занято")
	}

	existingUserByEmail, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUserByEmail != nil {
		return errors.New("почта уже занята")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	user.Role = "user"
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, uuid uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByUUID(ctx, uuid)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, uuid uuid.UUID) error {
	return s.userRepo.Delete(ctx, uuid)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *UserService) Authenticate(ctx context.Context, identifier, password string) (string, error) {
	var user *models.User
	var err error

	user, err = s.userRepo.GetByEmail(ctx, identifier)
	if err != nil {
		user, err = s.userRepo.GetByUsername(ctx, identifier)
	}

	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
	}
	return generateJWT(user)
}

func generateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(config.GetConfig().JWT.Expiration) * time.Second)
	claims := jwt.MapClaims{
		"user_id": user.UUID,
		"exp":     expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWT.Secret))
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
