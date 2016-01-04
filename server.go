package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/serve/lib/securecookie"
)

type Server struct {
	path string
	host string

	sites map[string]*Site
	http  *http.ServeMux
	mutex *sync.Mutex

	jar *securecookie.SecureCookie

	siteBuilders map[string]SiteBuilder
}

func (server *Server) SetConfig(path string, port string) {
	server.path = filepath.ToSlash(path)
	server.host = "localhost:" + port
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" || r.URL.Path == "/favicon" {
		http.NotFound(w, r)
		return
	}

	mux := server.http
	re := regexp.MustCompile("^/\\w+")
	path := re.FindString(r.URL.Path)
	appPath := server.path + "/sites" + path

	f := Stat(appPath)

	if f != nil {

		server.mutex.Lock()
		defer server.mutex.Unlock()

		if _, p := server.sites[f.Name]; p == false {

			site := new(Site)
			site.name = f.Name
			site.path = server.path + "/sites/" + site.name
			site.uri = "/" + site.name
			site.server = server

			builder, ok := server.siteBuilders[site.name]
			if ok {
				site.builder = builder
			} else {
				site.builder = server.siteBuilders["."]
			}
			site.Build()

			server.sites[site.name] = site
		}

		handler, _ := mux.Handler(r)
		handler.ServeHTTP(w, r)
		return
	}

	http.NotFound(w, r)
	return
}

func (server *Server) Start() {
	mux := http.NewServeMux()
	server.http = mux
	server.sites = make(map[string]*Site)
	server.mutex = &sync.Mutex{}

	server.SetupSiteBuilders()

	var hashKey = []byte(securecookie.GenerateRandomKey(16))
	var blockKey = []byte(securecookie.GenerateRandomKey(16))
	server.jar = securecookie.New(hashKey, blockKey)

	mux.Handle("/", server)

	if err := http.ListenAndServe(server.host, server.http); err != nil {
		fmt.Println("Unable to start server")
		fmt.Println(err)
	}
}

func (server *Server) SetupSiteBuilders() {
	server.siteBuilders = make(map[string]SiteBuilder)
	server.siteBuilders["."] = new(CommonSiteBuilder)
	server.siteBuilders["_auth"] = new(AuthSiteBuilder)
}
