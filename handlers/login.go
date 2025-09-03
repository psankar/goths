package handlers

import (
	"net/http"

	"psankar-goths-demo/templ"
)

func Login(w http.ResponseWriter, r *http.Request) {
	loginForm := templ.LoginForm()
	loginForm.Render(r.Context(), w)
}
