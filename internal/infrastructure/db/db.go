package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Gezubov/user_service/config"
	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

func InitDB(ctx context.Context, dbConfig *config.DatabaseConfig) {
	var err error

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Database)

	conn, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		return
	}
	slog.Info("Connected to database")
}

func GetDB() *pgx.Conn {
	return conn
}

func CloseDB(ctx context.Context) {
	if conn != nil {
		conn.Close(ctx)
	}
}
