package main

import (
	"log"
	"net/http"
)

func (a *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
}
