package main

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/justinas/nosurf"
	"net/http"
)

// TemplateData is used to hold some metadata
type TemplateData struct {
	URL             string
	IsAuthenticated bool
	AuthUser        string
	Flash           string
	Error           string
	CSRFToken       string
}

// defaultData is used to for some default data
func (a *application) defaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.URL = a.server.url

	if a.session != nil {
		if a.session.Exists(r.Context(), sessionKeyUserId) {
			td.IsAuthenticated = true
			td.AuthUser = a.session.GetString(r.Context(), sessionKeyUserName)
		}

		td.Flash = a.session.PopString(r.Context(), "flash")
	}

	td.CSRFToken = nosurf.Token(r)
	return td
}

// render is used to render the template
func (a *application) render(w http.ResponseWriter, r *http.Request, view string, vars jet.VarMap) error {
	td := &TemplateData{}

	td = a.defaultData(td, r)

	tp, err := a.view.GetTemplate(fmt.Sprintf("%s.html", view))
	if err != nil {
		return err
	}

	if err = tp.Execute(w, vars, td); err != nil {
		return err
	}

	return nil
}
