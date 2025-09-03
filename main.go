package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"psankar-goths-demo/libgoths"
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
	http.HandleFunc("GET /login", LoginGetHandler)
	http.HandleFunc("POST /login", LoginPostHandler)
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

func LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	loginForm := templ.LoginForm("")
	loginForm.Render(r.Context(), w)
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := queries.Login(r.Context(), db.LoginParams{
		Email:    email,
		Password: password,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			templ.LoginForm(libgoths.LoginFailed).Render(r.Context(), w)
			return
		}

		slog.Error("Error during login", "error", err)
		templ.LoginForm(libgoths.InternalError).Render(r.Context(), w)
		return
	}

	slog.Info("User logged in", "email", user.Email)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
