package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"sync"
)

type Server struct {
	path string
	host string

	sites map[string]*Site
	http  *http.ServeMux
	mutex *sync.Mutex
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

	mux := server.http
	re := regexp.MustCompile("^/\\w+")
	path := re.FindString(r.URL.Path)
	appPath := server.path + "/sites" + path

	f := Stat(appPath)

	server.mutex.Lock()
	defer server.mutex.Unlock()

	if _, p := server.sites[f.Name]; !p {
		if f != nil {
			site := new(Site)
			site.name = f.Name
			site.path = server.path + "/sites" + "/" + site.name
			site.uri = "/" + site.name + "/"
			site.server = server
			site.Build()
			server.sites[site.name] = site
		} else {
			http.NotFound(w, r)
			return
		}
	}

	handler, _ := mux.Handler(r)
	handler.ServeHTTP(w, r)
}

func (server *Server) Start() {
	mux := http.NewServeMux()
	server.http = mux
	server.sites = make(map[string]*Site)
	server.mutex = &sync.Mutex{}
	mux.Handle("/", server)

	if err := http.ListenAndServe(server.host, server.http); err != nil {
		fmt.Println("Unable to start server")
		fmt.Println(err)
	}
}
