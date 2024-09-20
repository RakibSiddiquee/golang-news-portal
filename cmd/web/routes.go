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
	mux.Post("/login", a.loginPostHandler)
	mux.Get("/signup", a.signupHandler)
	mux.Post("/signup", a.signupPostHandler)
	mux.Get("/logout", a.authRequired(a.logoutHandler))

	mux.Get("/vote", a.authRequired(a.voteHandler))
	mux.Get("/submit", a.authRequired(a.submitHandler))

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return mux
}
