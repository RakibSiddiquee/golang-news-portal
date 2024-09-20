package main

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/RakibSiddiquee/golang-news-portal/forms"
	"github.com/RakibSiddiquee/golang-news-portal/models"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func (a *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	/*	u := models.User{
			Email:    "ad@example.com",
			Password: "password",
			Name:     "Rakib",
		}
		err := a.Models.Users.Insert(&u)
		if err != nil {
			log.Fatal(err)
		}

		a.Models.Posts.Insert("Test Title", "https://example.com", u.ID)
		a.Models.Posts.Insert("Test Title 2", "https://example.com", u.ID)
		a.Models.Posts.Insert("Test Title 3", "https://example.com", u.ID)
		a.Models.Posts.Insert("Test Title 4", "https://example.com", u.ID)
		a.Models.Posts.Insert("Test Title 5", "https://example.com", u.ID)*/

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

	fmt.Println(posts)

	vars := make(jet.VarMap)
	vars.Set("posts", posts)
	vars.Set("meta", meta)
	vars.Set("nextUrl", nextUrl)
	vars.Set("prevUrl", prevUrl)
	vars.Set("form", forms.New(r.Form))

	err = a.render(w, r, "index", vars)

	if err != nil {
		log.Fatal(err)
	}
}

func (a *application) commentHandler(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)

	postId, err := strconv.Atoi(chi.URLParam(r, "postId"))
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	post, err := a.Models.Posts.Get(postId)
	if err != nil {
		a.serverError(w, err)
		return
	}

	comments, err := a.Models.Comments.GetForPost(post.ID)
	if err != nil {
		a.serverError(w, err)
		return
	}

	vars.Set("post", post)
	vars.Set("comments", comments)

	err = a.render(w, r, "comments", vars)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	err := a.render(w, r, "login", nil)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)
	vars.Set("form", forms.New(r.PostForm))

	err := a.render(w, r, "signup", vars)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *application) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*2)
	err := r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	form.Email("email")
	form.MinLength("password", 3)

	if !form.Valid() {
		vars := make(jet.VarMap)
		vars.Set("errors", form.Errors)
		err := a.render(w, r, "login", vars)
		if err != nil {
			a.serverError(w, err)
			return
		}
	}

	user, err := a.Models.Users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		a.session.Put(r.Context(), "flash", "Login error: "+err.Error())
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	a.session.RenewToken(r.Context())
	a.session.Put(r.Context(), sessionKeyUserId, user.ID)
	a.session.Put(r.Context(), sessionKeyUserName, user.Name)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	a.session.Remove(r.Context(), sessionKeyUserId)
	a.session.Remove(r.Context(), sessionKeyUserName)
	a.session.Destroy(r.Context())
	a.session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (a *application) signupPostHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*2)

	err := r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	vars := make(jet.VarMap)

	form := forms.New(r.PostForm)
	vars.Set("form", form)

	form.Required("name", "email", "password").
		Email("email")

	if !form.Valid() {
		vars.Set("errors", form.Errors)
		err := a.render(w, r, "signup", vars)
		if err != nil {
			a.serverError(w, err)
		}
		return
	}

	user := models.User{
		Name:      form.Get("name"),
		Email:     form.Get("email"),
		Password:  form.Get("password"),
		Activated: true, // move to database/config
	}

	err = a.Models.Users.Insert(&user)
	if err != nil {
		form.Fail("signup", "failed to create account "+err.Error())
		vars.Set("errors", form.Errors)
		err := a.render(w, r, "signup", vars)
		if err != nil {
			a.serverError(w, err)
		}
		return
	}

	a.session.Put(r.Context(), "flash", "account created!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (a *application) voteHandler(w http.ResponseWriter, r *http.Request) {
	id := a.readIntDefault(r, "id", 0)

	post, err := a.Models.Posts.Get(id)
	if err != nil {
		a.session.Put(r.Context(), "flash", "Error: "+err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userId := a.session.GetInt(r.Context(), sessionKeyUserId)
	err = a.Models.Posts.Vote(post.ID, userId)
	if err != nil {
		a.session.Put(r.Context(), "flash", "Error: "+err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	a.session.Put(r.Context(), "flash", "Voted successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *application) submitHandler(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)
	vars.Set("form", forms.New(r.PostForm))
	err := a.render(w, r, "submit", vars)
	if err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *application) submitPostHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*2)

	err := r.ParseForm()
	if err != nil {
		a.serverError(w, err)
		return
	}

	userId := a.session.GetInt(r.Context(), sessionKeyUserId)
	vars := make(jet.VarMap)
	form := forms.New(r.PostForm)

	form.Required("title", "url").Url("url").
		MaxLength("title", 255).MaxLength("url", 255)

	vars.Set("form", form)

	if !form.Valid() {
		vars.Set("errors", form.Errors)
		err := a.render(w, r, "submit", vars)
		if err != nil {
			a.serverError(w, err)
		}
		return
	}

	_, err = a.Models.Posts.Insert(form.Get("title"), form.Get("url"), userId)
	if err != nil {
		form.Fail("form", "failed due to "+err.Error())
		vars.Set("errors", form.Errors)
		err := a.render(w, r, "submit", vars)
		if err != nil {
			a.serverError(w, err)
		}
		return
	}

	a.session.Put(r.Context(), "flash", "post submitted successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
