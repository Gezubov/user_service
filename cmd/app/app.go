package main

import (
	"log"
	"net/http"

	"github.com/Gezubov/user_service/config"
	"github.com/Gezubov/user_service/internal/controller"
	"github.com/Gezubov/user_service/internal/infrastructure/db"
	"github.com/Gezubov/user_service/internal/repository"
	"github.com/Gezubov/user_service/internal/service"
	_ "github.com/lib/pq"
)

func main() {
	config.Load()

	db.InitDB()
	database := db.GetDB()

	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	http.HandleFunc("/user/get", userController.GetUser)
	http.HandleFunc("/user/create", userController.CreateUser)
	http.HandleFunc("/user/update", userController.UpdateUser)
	http.HandleFunc("/user/delete", userController.DeleteUser)
	http.HandleFunc("/user/list", userController.GetUsers)

	serverAddr := ":" + config.GetConfig().Server.Port
	log.Printf("Server started on port %s", config.GetConfig().Server.Port)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
