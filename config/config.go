package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_PATH"))
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

// MigrateDB выполняет автоматические миграции.
func MigrateDB(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("[ERROR] не удалось создать драйвер базы данных: %w", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"song_go",
		driver,
	)
	if err != nil {
		return fmt.Errorf("[ERROR] ошибка миграции: %w", err)
	}

	// Применяем миграции
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("[ERROR] ошибка миграции: %w", err)
	}

	log.Println("[INFO] миграция прошла успешно")
	return nil
}