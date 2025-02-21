package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Gezubov/user_service/internal/middlewares"
	"github.com/Gezubov/user_service/internal/models"
	"github.com/Gezubov/user_service/internal/service"
	"github.com/Gezubov/user_service/pkg/hash"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.userService.GetUserByID(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := c.userService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	currentUser, err := c.userService.GetUserByID(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var user models.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username != "" {
		currentUser.Username = user.Username
	}
	if user.Email != "" {
		existingUser, err := c.userService.GetByEmail(user.Email)
		if err == nil && existingUser.UUID != currentUser.UUID {
			http.Error(w, "Email already in use", http.StatusConflict)
			return
		}
		currentUser.Email = user.Email
	}
	if user.Password != "" {
		hash, err := hash.HashPassword(user.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		currentUser.PasswordHash = hash
	}

	if err := c.userService.UpdateUser(currentUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentUser)
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if userID != idStr {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	uuid, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	if err := c.userService.DeleteUser(uuid); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var input models.UserRegister

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if input.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}
	user := models.User{Username: input.Username, Email: input.Email}

	if err := c.userService.CreateUser(&user, input.Password); err != nil {
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := c.userService.Authenticate(input.Identifier, input.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, err = c.userService.GetByEmail(input.Identifier)
	if err != nil {
		user, err = c.userService.GetByUsername(input.Identifier)
	}

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
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
