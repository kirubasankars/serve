package serve

import (
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	html "html/template"
	text "text/template"

	"github.com/kirubasankars/metal"
)

type HtmlTemplateHandler struct{}

func (handler *HtmlTemplateHandler) ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/html/")
	var input *metal.Metal
	json, _ := ioutil.ReadAll(r.Body)
	if len(json) > 0 {
		input = metal.NewMetal()
		input.Parse(json)
	}
	model := site.Model(r.Method, parts[1], input)
	handler.handle(site, model, w, r)
}

func (handler *HtmlTemplateHandler) handle(site *Site, model *metal.Metal, w http.ResponseWriter, r *http.Request) {
	server := site.server

	if model == nil {
		http.Error(w, "No api found.", http.StatusInternalServerError)
		return
	}

	paths := strings.Split(r.URL.Path, "html/")

	w.Header().Set("Content-Type", "text/html")
	templatePath := path.Join(site.path, "tpl", "html", paths[1], "get.html")

	tpl := server.IO.Template(templatePath)
	if tpl != nil {
		tmpl, err := html.New(templatePath).Parse(string(*tpl))
		if err != nil {
			http.Error(w, "No Template found.", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, model.Raw()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "No Template found.", http.StatusInternalServerError)
		return
	}

}

type TextTemplateHandler struct{}

func (handler *TextTemplateHandler) ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/text/")
	var input *metal.Metal
	json, _ := ioutil.ReadAll(r.Body)
	if len(json) > 0 {
		input = metal.NewMetal()
		input.Parse(json)
	}
	model := site.Model(r.Method, parts[1], input)
	handler.handle(site, model, w, r)
}

func (handler *TextTemplateHandler) handle(site *Site, model *metal.Metal, w http.ResponseWriter, r *http.Request) {
	server := site.server

	if model == nil {
		http.Error(w, "No api found.", http.StatusInternalServerError)
		return
	}

	paths := strings.Split(r.URL.Path, "text/")

	w.Header().Set("Content-Type", "text/plain")
	templatePath := path.Join(site.path, "tpl", "text", paths[1], "get.txt")

	tpl := server.IO.Template(templatePath)
	if tpl != nil {
		tmpl, err := text.New(templatePath).Parse(string(*tpl))
		if err != nil {
			http.Error(w, "No Template found.", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, model.Raw()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "No Template found.", http.StatusInternalServerError)
		return
	}

}
