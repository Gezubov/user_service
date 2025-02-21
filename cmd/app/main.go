package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Gezubov/user_service/config"
	"github.com/Gezubov/user_service/internal/controller"
	"github.com/Gezubov/user_service/internal/infrastructure/db"
	"github.com/Gezubov/user_service/internal/middlewares"
	"github.com/Gezubov/user_service/internal/repository"
	"github.com/Gezubov/user_service/internal/service"
	"github.com/go-chi/chi"

	_ "github.com/lib/pq"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Loading configuration...")
	config.Load()

	slog.Info("Initializing database...")
	db.InitDB()

	database := db.GetDB()
	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r := chi.NewRouter()
	r.Use(middlewares.CorsMiddleware())
	r.Post("/user", userController.CreateUser)
	r.Get("/user/{id}", userController.GetUser)
	r.Patch("/user/{id}", userController.UpdateUser)
	r.Delete("/user/{id}", userController.DeleteUser)
	r.Get("/users", userController.GetUsers)

	port := config.GetConfig().Server.Port
	serverAddr := ":" + port

	slog.Info("Server started", "port", port)
	slog.Error("Server crashed", "error", http.ListenAndServe(serverAddr, r))
}
