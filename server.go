package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type Server struct {
	path string
	host string

	sites map[string]*Site
	http  *http.ServeMux
}

func (server *Server) SetConfig(path string, port string) {
	server.path = filepath.ToSlash(path)
	server.host = "localhost:" + port
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {
		w.Write([]byte("Root"))
		return
	}

	var (
		mux  = server.http
		re   = regexp.MustCompile("^/\\w+")
		path = ""
	)
	path = re.FindString(r.URL.Path)
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

		handler, _ := mux.Handler(r)
		handler.ServeHTTP(w, r)

		return
	} else {
		http.NotFound(w, r)

		return
	}
}

func (server *Server) Start() {
	mux := http.NewServeMux()
	server.http = mux
	server.sites = make(map[string]*Site)
	mux.Handle("/", server)

	if err := http.ListenAndServe(server.host, server.http); err != nil {
		fmt.Println("Unable to start server")
		fmt.Println(err)
	}
}
