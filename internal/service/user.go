package service

import (
	"errors"
	"time"

	"github.com/Gezubov/user_service/config"
	"github.com/Gezubov/user_service/internal/models"
	"github.com/Gezubov/user_service/pkg/hash"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(user *models.User, password string) error {
	existingUserByUsername, err := s.userRepo.GetByUsername(user.Username)
	if err == nil && existingUserByUsername != nil {
		return errors.New("имя пользователя уже занято")
	}

	existingUserByEmail, err := s.userRepo.GetByEmail(user.Email)
	if err == nil && existingUserByEmail != nil {
		return errors.New("почта уже занята")
	}

	hash, err := hash.HashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	user.Role = "user"
	return s.userRepo.Create(user)
}

func (s *UserService) GetUserByID(uuid uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByUUID(uuid)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(uuid uuid.UUID) error {
	return s.userRepo.Delete(uuid)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *UserService) Authenticate(identifier, password string) (string, error) {
	var user *models.User
	var err error

	user, err = s.userRepo.GetByEmail(identifier)
	if err != nil {
		user, err = s.userRepo.GetByUsername(identifier)
	}

	if err != nil {
		return "", ErrInvalidCredentials
	}

	if !hash.CheckPasswordHash(password, user.PasswordHash) {
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
