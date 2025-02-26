package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gezubov/user_service/config"
	"github.com/Gezubov/user_service/internal/controller"
	"github.com/Gezubov/user_service/internal/infrastructure/db"
	"github.com/Gezubov/user_service/internal/middlewares"
	"github.com/Gezubov/user_service/internal/service"
	"github.com/Gezubov/user_service/internal/storage"
	"github.com/go-chi/chi"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Loading configuration...")
	config.Load()

	ctx := context.Background()

	slog.Info("Initializing database...")
	db.InitDB(ctx, &config.GetConfig().Database)
	database := db.GetDB()

	userRepo := storage.NewUserStorage(ctx, database)
	userService := service.NewUserService(ctx, userRepo)
	userController := controller.NewUserController(ctx, userService)

	r := SetupRoutes(userController)

	port := config.GetConfig().Server.Port
	serverAddr := ":" + port
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Server started", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server crashed", "error", err)
		}
	}()

	<-stop
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	} else {
		slog.Info("Server exited properly")
	}

	db.CloseDB(ctx)
}

func SetupRoutes(userController *controller.UserController) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.CorsMiddleware())

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", userController.Register)
		r.Post("/login", userController.Login)
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/{id}", userController.GetUser)
		r.With(middlewares.AuthMiddleware).Patch("/{id}", userController.UpdateUser)
		r.With(middlewares.AuthMiddleware).Delete("/{id}", userController.DeleteUser)
	})
	r.Get("/users", userController.GetUsers)

	return r
}
