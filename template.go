package main

import (
	"net/http"
	"path"
	"strings"

	html "html/template"
	text "text/template"

	"github.com/serve/lib/metal"
)

type HtmlTemplateHandler struct {
}

func (handler *HtmlTemplateHandler) ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/html")
	model := site.Model(parts[1] + "/" + strings.ToLower(r.Method))
	handler.Handle(site, model, w, r)
}

func (handler *HtmlTemplateHandler) Handle(site *Site, model *metal.Metal, w http.ResponseWriter, r *http.Request) {

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

type TextTemplateHandler struct {
}

func (handler *TextTemplateHandler) ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/text")
	model := site.Model(parts[1] + "/" + strings.ToLower(r.Method))
	handler.Handle(site, model, w, r)
}

func (handler *TextTemplateHandler) Handle(site *Site, model *metal.Metal, w http.ResponseWriter, r *http.Request) {
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
