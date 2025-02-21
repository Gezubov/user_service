package main

import (
	"log"
	"net/http"

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
	config.Load()
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

	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}
