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

	//mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
	//	a.session.Put(r.Context(), "test", "Hello World!")
	//	err := a.render(w, r, "index", nil)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//})
	//
	//mux.Get("/comments", func(w http.ResponseWriter, r *http.Request) {
	//	vars := make(jet.VarMap)
	//	tt := a.session.GetString(r.Context(), "test")
	//	fmt.Println("tt", tt)
	//	vars.Set("test", a.session.GetString(r.Context(), "test"))
	//	err := a.render(w, r, "index", vars)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//})

	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return mux
}
