package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)


///add CSRF protection to the application
// NoSurf adds CSRF protection to all POST requests.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

/// SessionLoad loads and saves the session for each request. 
func SessionLoad(next http.Handler) http.Handler {

	return session.LoadAndSave(next)
}