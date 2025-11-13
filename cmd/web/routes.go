package main

import (
	"net/http"

	"github.com/Chocolate529/nevarol/internal/config"
	"github.com/Chocolate529/nevarol/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/time/rate"
)

func routes(app *config.AppConfig) http.Handler {
	// mux := pat.New()

	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	mux := chi.NewRouter()

	// Create rate limiter: 100 requests per minute with burst of 200
	rateLimiter := NewRateLimiter(rate.Limit(100.0/60.0), 200)
	go rateLimiter.CleanupVisitors()

	mux.Use(middleware.Recoverer)
	mux.Use(SecurityHeaders)
	mux.Use(RateLimit(rateLimiter))
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// Page routes
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/store", handlers.Repo.Store)
	mux.Get("/shipping", handlers.Repo.Shipping)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/checkout", handlers.Repo.Checkout)
	mux.Get("/account", handlers.Repo.Account)
	mux.Get("/login", handlers.Repo.Login)

	// API routes
	mux.Route("/api", func(r chi.Router) {
		// Auth routes
		r.Post("/register", handlers.Repo.Register)
		r.Post("/login", handlers.Repo.LoginAPI)
		r.Post("/logout", handlers.Repo.LogoutAPI)
		r.Get("/user", handlers.Repo.GetCurrentUser)

		// Product routes
		r.Get("/products", handlers.Repo.GetProducts)

		// Cart routes
		r.Get("/cart", handlers.Repo.GetCart)
		r.Post("/cart", handlers.Repo.AddToCart)
		r.Put("/cart/{id}", handlers.Repo.UpdateCartItem)
		r.Delete("/cart/{id}", handlers.Repo.RemoveFromCart)
		r.Delete("/cart", handlers.Repo.ClearCart)

		// Order routes
		r.Post("/orders", handlers.Repo.CreateOrder)
		r.Get("/orders", handlers.Repo.GetOrders)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
