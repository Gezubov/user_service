package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Gezubov/user_service/internal/middlewares"
	"github.com/Gezubov/user_service/internal/models"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

var ErrInvalidUserID = errors.New("invalid user ID")
var ErrInvalidRequestBody = errors.New("invalid request body")
var ErrUnauthorized = errors.New("unauthorized")
var ErrForbidden = errors.New("forbidden")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrEmailAlreadyInUse = errors.New("email already in use")
var ErrMethodNotAllowed = errors.New("method not allowed")
var ErrUserIDRequired = errors.New("user ID is required")
var ErrPasswordRequired = errors.New("password is required")

type UserService interface {
	CreateUser(ctx context.Context, user *models.User, password string) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	Authenticate(ctx context.Context, identifier, password string) (string, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type UserController struct {
	userService UserService
	ctx         context.Context
}

func NewUserController(ctx context.Context, userService UserService) *UserController {
	return &UserController{
		ctx:         ctx,
		userService: userService,
	}
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, ErrInvalidUserID.Error(), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, ErrInvalidUserID.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.userService.GetUserByID(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	users, err := c.userService.GetAllUsers(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, ErrUserIDRequired.Error(), http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, ErrInvalidUserID.Error(), http.StatusBadRequest)
		return
	}

	currentUser, err := c.userService.GetUserByID(context.Background(), uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var user models.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	if user.Username != "" {
		currentUser.Username = user.Username
	}
	if user.Email != "" {
		existingUser, err := c.userService.GetByEmail(context.Background(), user.Email)
		if err == nil && existingUser.UUID != currentUser.UUID {
			http.Error(w, ErrEmailAlreadyInUse.Error(), http.StatusConflict)
			return
		}
		currentUser.Email = user.Email
	}

	if err := c.userService.UpdateUser(context.Background(), currentUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentUser)
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, ErrUserIDRequired.Error(), http.StatusBadRequest)
		return
	}

	if userID != idStr {
		http.Error(w, ErrForbidden.Error(), http.StatusForbidden)
		return
	}

	uuid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, ErrInvalidUserID.Error(), http.StatusBadRequest)
		return
	}

	if err := c.userService.DeleteUser(context.Background(), uuid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var input models.UserRegister

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}
	if input.Password == "" {
		http.Error(w, ErrPasswordRequired.Error(), http.StatusBadRequest)
		return
	}
	user := models.User{Username: input.Username, Email: input.Email}

	if err := c.userService.CreateUser(context.Background(), &user, input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"uuid":    user.UUID,
		"message": "User registered successfully",
	})
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var input models.UserLogin

	var err error

	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	token, err := c.userService.Authenticate(context.Background(), input.Identifier, input.Password)
	if err != nil {
		http.Error(w, ErrInvalidCredentials.Error(), http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, err = c.userService.GetByEmail(context.Background(), input.Identifier)
	if err != nil {
		user, err = c.userService.GetByUsername(context.Background(), input.Identifier)
	}

	if err != nil {
		http.Error(w, ErrInvalidCredentials.Error(), http.StatusUnauthorized)
		return
	}

	if user == nil {
		http.Error(w, ErrUserNotFound.Error(), http.StatusNotFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"uuid":    user.UUID,
	})
}
