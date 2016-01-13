package serve

import (
	"net/http"
	"strings"
)

type FileHandler struct{}

func (handler *FileHandler) ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, site.uri)
	server := site.server
	appDir := server.IO.IsExists(site.path + "/app")

	filePath := ""
	if appDir {
		filePath = site.path + "/app/" + parts[1]
	} else {
		filePath = site.path + "/" + parts[1]
	}

	server.IO.ServeFile(w, r, filePath)
}
