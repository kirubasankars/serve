package main

import (
	html "html/template"
	"net/http"
	"path"
	"strings"
	text "text/template"

	"github.com/metal"
)

func (site *Site) HandleHTMLTemplate(model *metal.Metal, w http.ResponseWriter, r *http.Request) {

	if model == nil {
		http.Error(w, "No api found.", http.StatusInternalServerError)
		return
	}

	paths := strings.Split(r.URL.Path, "html/")

	w.Header().Set("Content-Type", "text/html")
	templatePath := path.Join(site.path, "tpl/html/", paths[1], "get.html")

	tmpl, err := html.ParseFiles(templatePath)

	if err != nil {
		http.Error(w, "No Template found.", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, model.Raw()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (site *Site) HandleTextTemplate(model *metal.Metal, w http.ResponseWriter, r *http.Request) {
	if model == nil {
		http.Error(w, "No api found.", http.StatusInternalServerError)
		return
	}

	paths := strings.Split(r.URL.Path, "text/")

	w.Header().Set("Content-Type", "text/plain")
	templatePath := path.Join(site.path, "tpl/text/", paths[1], "get.txt")

	tmpl, err := text.ParseFiles(templatePath)

	if err != nil {
		http.Error(w, "No Template found.", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, model.Raw()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
