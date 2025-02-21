package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/Gezubov/user_service/config"
)

var db *sql.DB

func InitDB() {
	dbHost := config.GetConfig().Database.Host
	dbPort := config.GetConfig().Database.Port
	dbUser := config.GetConfig().Database.Username
	dbPassword := config.GetConfig().Database.Password
	dbName := config.GetConfig().Database.Database

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
	}

	slog.Info("Connecting to database...")
	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		slog.Error("Failed to ping database, retrying in 5 seconds...", "error", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		slog.Error("Failed to ping database after 5 retries", "error", err)
	}
	slog.Info("Database connection established")
}

func GetDB() *sql.DB {
	slog.Info("Returning database connection")
	return db
}

func CloseDB() {
	if db != nil {
		slog.Info("Closing database connection...")
		if err := db.Close(); err != nil {
			slog.Error("Error closing database", "error", err)
		} else {
			slog.Info("Database connection closed")
		}
	}
}
