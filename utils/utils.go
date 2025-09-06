package utils

import (
	"net/http"
	"time"
)

const (
	// Cookie constants
	SessionCookieName = "goths_session"
	CookieMaxAge      = 24 * time.Hour
)

// SetSessionCookie creates a secure session cookie for the authenticated user
func SetSessionCookie(w http.ResponseWriter, userEmail string) {
	http.SetCookie(w, &http.Cookie{
		Name: SessionCookieName,

		// Do not do this in production code. This is just for learning.
		Value: userEmail,

		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(CookieMaxAge.Seconds()),
	})
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
}

// GetAuthenticatedUser retrieves the authenticated user's email from the session cookie
func GetAuthenticatedUser(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}
