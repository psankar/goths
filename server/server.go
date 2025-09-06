package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"psankar-goths-demo/sqlc/db"
	"psankar-goths-demo/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

// requireAuth is a middleware that checks if the user is authenticated
func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, authenticated := utils.GetAuthenticatedUser(r)
		if !authenticated {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

type server struct {
	mux     *http.ServeMux
	queries *db.Queries
}

func RunServer() {
	pgURL := os.Getenv("POSTGRES_URL")
	migrationsDir := os.Getenv("MIGRATIONS_DIR")

	slog.Debug("env", "POSTGRES_URL:", pgURL, "MIGRATIONS_DIR:", migrationsDir)

	conn, err := pgx.Connect(context.Background(), pgURL)
	if err != nil {
		log.Fatal("Database connection failure", err)
		return
	}

	queries := db.New(conn)

	m, err := migrate.New(migrationsDir, pgURL)
	if err != nil {
		log.Fatal("Error creating migrate instance", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migrations", err)
		return
	}

	srv := server{
		mux:     http.DefaultServeMux,
		queries: queries,
	}

	// Redirects to the HomePage if the user is authenticated; else to the login page
	srv.mux.HandleFunc("/", RootHandler)

	// Show the login form
	srv.mux.HandleFunc("GET /login", LoginGetHandler)

	// Authenticate the user
	srv.mux.HandleFunc("POST /login", srv.LoginPostHandler)

	// Logs out the user and clear the cookie
	srv.mux.HandleFunc("POST /logout", LogoutHandler)

	// Loads the home page if the user is authenticated
	srv.mux.HandleFunc("/home", requireAuth(HomeHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
