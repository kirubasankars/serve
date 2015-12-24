package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	path string
	port string
}

var rootPath = ""

func (s *Server) SetConfig(path string, port string) {
	s.path = filepath.ToSlash(path)
	s.port = port
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

func (s *Server) Start() {
	rootPath = s.path + "/sites"
	sites, _ := ioutil.ReadDir(rootPath)
	for _, f := range sites {
		if f.IsDir() {
			var siteDir http.Handler

			siteRoot := "/" + f.Name()
			site := new(Site)
			site.name = f.Name()
			site.path = s.path + "/sites" + siteRoot
			site.uri = siteRoot
			r, _ := exists(site.path + "/app")

			if r {
				siteDir = http.FileServer(http.Dir(site.path + "/app"))
			} else {
				siteDir = http.FileServer(http.Dir(site.path))
			}

			http.Handle(site.uri, http.StripPrefix("/"+site.uri, siteDir))

			http.HandleFunc(site.uri+"/api/", func(w http.ResponseWriter, r *http.Request) {
				HandleAPI(site, w, r)
			})
			http.HandleFunc(site.uri+"/tpl/", func(w http.ResponseWriter, r *http.Request) {
				HandleTemplate(site, w, r)
			})

			fmt.Println("http://localhost:" + s.port + "/" + site.uri)
		}
	}

	if err := http.ListenAndServe(":"+s.port, nil); err != nil {
		fmt.Println("unable to start server")
		fmt.Println(err)
	}
}
