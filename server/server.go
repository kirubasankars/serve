package server

import (
	"net/http"
	"io/ioutil"
	"fmt"	
	"os"
)

type Server struct {
	path string
	port string
}

type Site struct {
	name string
	path string
	uri string
}

func (s *Server) SetConfig(path string, port string) {
	s.path = path
	s.port = port
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func (s *Server) Start() {
	sites, _ := ioutil.ReadDir(s.path + "/sites")
    for _, f := range sites {
		if f.IsDir() {
			siteRoot := "/" + f.Name() + "/"
			site := new(Site)
			site.name = f.Name()
			site.path = s.path + "/sites" + siteRoot
			site.uri = siteRoot
			var r, _ = exists(site.path + "/app")
			var siteDir http.Handler
			fmt.Println("http://localhost:" + s.port + site.uri)
			if(r) {
				siteDir = http.FileServer(http.Dir(site.path + "/app"))
			} else {
				siteDir = http.FileServer(http.Dir(site.path))
			}
			http.Handle(site.uri, http.StripPrefix(site.uri, siteDir))
			
		}
    }
	if err := http.ListenAndServe(":" + s.port, nil); err != nil {
		fmt.Println("unable to start server")
		fmt.Println(err)
	}
}