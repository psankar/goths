package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"psankar-goths-demo/sqlc/db"
	"psankar-goths-demo/templ"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

var queries *db.Queries

func main() {
	pgURL := os.Getenv("POSTGRES_URL")
	migrationsDir := os.Getenv("MIGRATIONS_DIR")

	slog.Debug("env", "POSTGRES_URL:", pgURL, "MIGRATIONS_DIR:", migrationsDir)

	conn, err := pgx.Connect(context.Background(), pgURL)
	if err != nil {
		log.Fatal("Database connection failure", err)
		return
	}

	queries = db.New(conn)

	m, err := migrate.New(migrationsDir, pgURL)
	if err != nil {
		log.Fatal("Error creating migrate instance", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migrations", err)
		return
	}

	slog.Info("Migrations applied successfully. Launching server...")
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/home", HomeHandler)

	http.ListenAndServe(":8080", nil)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	homePage := templ.HomePage()
	homePage.Render(r.Context(), w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginForm := templ.LoginForm()
	loginForm.Render(r.Context(), w)
}
