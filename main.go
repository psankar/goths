package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"psankar-goths-demo/handlers"
	"psankar-goths-demo/sqlc/db"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

func main() {
	pgURL := os.Getenv("POSTGRES_URL")
	migrationsDir := os.Getenv("MIGRATIONS_DIR")

	slog.Debug("env", "POSTGRES_URL:", pgURL, "MIGRATIONS_DIR:", migrationsDir)

	conn, err := pgx.Connect(context.Background(), pgURL)
	if err != nil {
		log.Fatal("Database connection failure", err)
		return
	}

	queries := db.New(conn)
	_ = queries

	m, err := migrate.New(migrationsDir, pgURL)
	if err != nil {
		log.Fatal("Error creating migrate instance", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migrations", err)
		return
	}

	http.HandleFunc("/", handlers.Login)

	slog.Info("Launched goths on :8080")
	http.ListenAndServe(":8080", nil)
}
