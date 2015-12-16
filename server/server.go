package server

import (
	"net/http"
	"io/ioutil"
	"fmt"	
	"os"
	"html"
	"github.com/gravity"
	"strings"
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

func HandleApi(w http.ResponseWriter, r *http.Request) {	
	paths := strings.Split(r.URL.Path, "/api")
	v := gravity.Get(strings.ToLower(r.Method) ,html.EscapeString(paths[1]))
	fmt.Fprintf(w, "Hello, %q", v)
}

func (s *Server) Start() {
	sites, _ := ioutil.ReadDir(s.path + "/sites")
    for _, f := range sites {
		if f.IsDir() {
			var siteDir http.Handler
			
			siteRoot := "/" + f.Name() + "/"
			site := new(Site)
			site.name = f.Name()
			site.path = s.path + "/sites" + siteRoot
			site.uri = siteRoot
			r, _ := exists(site.path + "/app")
						
			if(r) {
				siteDir = http.FileServer(http.Dir(site.path + "/app"))
			} else {
				siteDir = http.FileServer(http.Dir(site.path))
			}
			
			http.HandleFunc(site.uri + "api/", HandleApi)
			http.Handle(site.uri, http.StripPrefix(site.uri, siteDir))						
			
			fmt.Println("http://localhost:" + s.port + site.uri)
		}
    }
	
	if err := http.ListenAndServe(":" + s.port, nil); err != nil {
		fmt.Println("unable to start server")
		fmt.Println(err)
	}
}