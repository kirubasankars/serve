package main

import (
	"net/http"
	"strings"
)

type FileServe struct {
	site *Site
}

func (fileServe *FileServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var site = fileServe.site
	parts := strings.Split(r.URL.Path, site.uri)
	filePath := ""
	appDir, _ := exists(site.path + "/app")
	if appDir {
		filePath = site.path + "/app/" + parts[1]
	} else {
		filePath = site.path + "/" + parts[1]
	}
	http.ServeFile(w, r, filePath)
}
