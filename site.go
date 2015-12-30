package main

import (
	"net/http"
	"os"
	"strings"
)

type Site struct {
	name string
	path string
	uri  string

	server *Server
}

func (site *Site) Build() {

	http.Handle(site.uri, &FileServe{site})

	http.HandleFunc(site.uri+"api/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/api")
		model := site.GetModel(parts[1] + "/" + strings.ToLower(r.Method))
		site.HandleAPI(model, w, r)
	})

	http.HandleFunc(site.uri+"html/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/html")
		model := site.GetModel(parts[1] + "/" + strings.ToLower(r.Method))
		site.HandleHTMLTemplate(model, w, r)
	})

	http.HandleFunc(site.uri+"text/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/text")
		model := site.GetModel(parts[1] + "/" + strings.ToLower(r.Method))
		site.HandleTextTemplate(model, w, r)
	})

}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
