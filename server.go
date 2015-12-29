package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type Server struct {
	path string
	host string
}

func (server *Server) SetConfig(path string, port string) {
	server.path = filepath.ToSlash(path)
	server.host = "localhost:" + port
}

func (server *Server) Start() {
	appsPath := server.path + "/sites"
	sitesDir, _ := ioutil.ReadDir(appsPath)

	for _, f := range sitesDir {
		if f.IsDir() {
			site := new(Site)
			site.name = f.Name()
			site.path = server.path + "/sites" + "/" + f.Name()
			site.uri = "/" + f.Name() + "/"
			site.server = server
			site.Build()
		}
	}

	if err := http.ListenAndServe(server.host, nil); err != nil {
		fmt.Println("unable to start server")
		fmt.Println(err)
	}
}
