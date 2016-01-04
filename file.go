package main

import (
	"net/http"
	"strings"
)

type FileHandler struct {
}

func (handler *FileHandler) ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, site.uri)
	filePath := ""
	appDir := IsExists(site.path + "/app")
	if appDir {
		filePath = site.path + "/app/" + parts[1]
	} else {
		filePath = site.path + "/" + parts[1]
	}
	http.ServeFile(w, r, filePath) //TODO:Should be using ServeContent
}
