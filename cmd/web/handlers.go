package main

import (
	"github.com/RakibSiddiquee/golang-news-portal/models"
	"log"
	"net/http"
)

func (a *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	filter := models.Filter{
		Query:    "",
		Page:     1,
		PageSize: 1,
		OrderBy:  "",
	}
}
