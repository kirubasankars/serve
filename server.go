package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	path string
	host string

	sites map[string]*Site
}

func (server *Server) SetConfig(path string, port string) {
	server.path = filepath.ToSlash(path)
	server.host = "localhost:" + port
}

func (server *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	var (
		idx    = strings.Index(r.URL.Path[1:], "/")
		length = len(r.URL.Path)
		path   = ""
	)

	switch {
	case idx == -1 && length > 1:
		path = r.URL.Path[:length]
	case idx == -1 && length == 1:
		path = "/"
	default:
		path = r.URL.Path[:idx+1]
	}

	appPath := server.path + "/sites" + path
	f, err := os.Stat(appPath)
	if err != nil || f == nil {
		http.NotFound(w, r)
		return
	}

	if f.IsDir() {
		site := new(Site)
		site.name = f.Name()
		site.path = server.path + "/sites" + "/" + f.Name()
		site.uri = "/" + f.Name() + "/"
		site.server = server
		site.Build()

		server.sites[f.Name()] = site

		fmt.Println(server.host + site.uri)

		http.Redirect(w, r, r.URL.Path, http.StatusFound)
		return
	} else {
		http.NotFound(w, r)
		return
	}
}

func (server *Server) Start() {
	server.sites = make(map[string]*Site)

	http.HandleFunc("/", server.HandleRequest)

	if err := http.ListenAndServe(server.host, nil); err != nil {
		fmt.Println("unable to start server")
		fmt.Println(err)
	}
}
