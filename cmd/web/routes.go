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

	if a.debug {
		mux.Use(middleware.Logger)
	}

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(a.appName))
	})
	mux.Get("/comments", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("comments"))
	})

	return mux
}
