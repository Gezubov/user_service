package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Gezubov/file_storage/config"
)

var db *sql.DB

func InitDB() {
	cfg := config.GetConfig().Database

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database is ready")
}

func GetDB() *sql.DB {
	return db
}
