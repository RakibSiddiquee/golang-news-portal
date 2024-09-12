package main

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/RakibSiddiquee/golang-news-portal/models"
	"log"
	"net/http"
)

func (a *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	//u := models.User{
	//	Email:    "ad@example.com",
	//	Password: "password",
	//	Name:     "Rakib",
	//}
	//err := a.Models.Users.Insert(&u)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//a.Models.Posts.Insert("Test Title", "https://example.com", u.ID)
	//a.Models.Posts.Insert("Test Title 2", "https://example.com", u.ID)
	//a.Models.Posts.Insert("Test Title 3", "https://example.com", u.ID)
	//a.Models.Posts.Insert("Test Title 4", "https://example.com", u.ID)
	//a.Models.Posts.Insert("Test Title 5", "https://example.com", u.ID)

	err := r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	filter := models.Filter{
		Query:    r.URL.Query().Get("q"),
		Page:     a.readIntDefault(r, "page", 1),
		PageSize: a.readIntDefault(r, "page_size", 5),
		OrderBy:  r.URL.Query().Get("order_by"),
	}

	posts, meta, err := a.Models.Posts.GetAll(filter)
	if err != nil {
		a.serverError(w, err)
		return
	}

	queryUrl := fmt.Sprintf("page_size=%d&order_by=%s&q=%s", meta.PageSize, filter.OrderBy, filter.Query)
	nextUrl := fmt.Sprintf("%s&page=%d", queryUrl, meta.NextPage)
	prevUrl := fmt.Sprintf("%s&page=%d", queryUrl, meta.PrevPage)

	vars := make(jet.VarMap)
	vars.Set("posts", posts)
	vars.Set("meta", meta)
	vars.Set("nextUrl", nextUrl)
	vars.Set("prevUrl", prevUrl)

	err = a.render(w, r, "index", vars)

	if err != nil {
		log.Fatal(err)
	}
}
