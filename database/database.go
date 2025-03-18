package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://browz_user:123123123@localhost:5432/browz"
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	DB = conn
	fmt.Println("✅ Подключение к БД успешно!")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close(context.Background())
	}
}
