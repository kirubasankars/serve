package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	html "html/template"
	text "text/template"

	"github.com/serve/lib/metal"
)

type Site struct {
	name string
	path string
	uri  string

	server *Server
}

func (site *Site) Build() {
	var mux = site.server.http

	mux.Handle(site.uri, Handler(&FileServe{site}))
	mux.Handle(site.uri+"api/", Handler(&ApiServe{site}))
	mux.Handle(site.uri+"html/", Handler(&HtmlServe{site}))
	mux.Handle(site.uri+"text/", Handler(&TextServe{site}))
}

func (site *Site) HandleAPI(model *metal.Metal, w http.ResponseWriter, r *http.Request) {
	if model == nil {
		http.Error(w, "No api found.", http.StatusInternalServerError)
		return
	}
	rs, err := json.Marshal(model.Raw())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(rs)
}

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

func Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
	h := http.HandlerFunc(fn)
	return loggingHandler(h)
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}
