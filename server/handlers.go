package server

import (
	"log/slog"
	"net/http"
	"psankar-goths-demo/libgoths"
	"psankar-goths-demo/sqlc/db"
	"psankar-goths-demo/templ"
	"psankar-goths-demo/utils"

	"github.com/jackc/pgx/v5"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	email, _ := utils.GetAuthenticatedUser(r)
	homePage := templ.HomePage(email)
	homePage.Render(r.Context(), w)
}

func LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	loginForm := templ.LoginForm("")
	loginForm.Render(r.Context(), w)
}

func (s *server) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := s.queries.Login(r.Context(), db.LoginParams{
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
	utils.SetSessionCookie(w, user.Email)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	utils.ClearSessionCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
