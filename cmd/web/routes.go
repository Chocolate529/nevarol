package main

import (
	"net/http"

	"github.com/Chocolate529/nevarol/internal/config"
	"github.com/Chocolate529/nevarol/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	// mux := pat.New()

	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/store", handlers.Repo.Store)
	mux.Get("/shipping", handlers.Repo.Shipping)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/checkout", handlers.Repo.Checkout)
	mux.Get("/account", handlers.Repo.Account)
	mux.Get("/login", handlers.Repo.Login)


	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
