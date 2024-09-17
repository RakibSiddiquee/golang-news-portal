package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

// routes function returns all the routes
func (a *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(a.LoadSession)

	if a.debug {
		mux.Use(middleware.Logger)
	}

	// Register routes
	mux.Get("/", a.homeHandler)
	mux.Get("/comments/{postId}", a.commentHandler)

	mux.Get("/login", a.loginHandler)

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return mux
}
