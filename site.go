package main

import (
	"fmt"
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

	http.HandleFunc(site.uri, func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, site.uri)
		filePath := ""
		appDir, _ := exists(site.path + "/app")
		if appDir {
			filePath = site.path + "/app/" + parts[1]
		} else {
			filePath = site.path + "/" + parts[1]
		}
		http.ServeFile(w, r, filePath)
	})

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

	fmt.Println(site.server.host + site.uri)
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
