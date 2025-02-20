package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Gezubov/user_service/config"
)

var db *sql.DB

func InitDB() {
	dbHost := config.GetConfig().Database.Host
	dbPort := config.GetConfig().Database.Port
	fmt.Println("port and host > ", dbHost, dbPort)
	dbUser := config.GetConfig().Database.Username
	dbPassword := config.GetConfig().Database.Password
	dbName := config.GetConfig().Database.Database

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connecting to database...")

	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		log.Printf("Failed to ping database, retrying in 5 seconds... Error: %v", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to ping database after 5 retries: %v", err)
	}
	log.Println("Database is ready")
}

func GetDB() *sql.DB {
	return db
}
