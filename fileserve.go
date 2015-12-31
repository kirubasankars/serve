package main

import (
	"net/http"
	"os"
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
