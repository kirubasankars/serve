package serve

import (
	"log"
	"net/http"
	"time"
)

type Site struct {
	name string
	path string
	uri  string

	server  *Server
	builder SiteBuilder
}

func (site *Site) Build() {
	site.builder.Build(site)
}

func (site *Site) Path() string {
	return site.path
}

func (site *Site) Name() string {
	return site.path
}

func (site *Site) SetupHandler(pattern string, handler HttpHandler) {
	http := site.server.http
	http.Handle(pattern, &SiteHandler{site: site, handler: handler})
}

func (site *Site) IsAuthEnabled() bool {
	return site.server.IO.IsExists(site.path + "/_auth")
}

type HttpHandler interface {
	ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request)
}

type SiteHandler struct {
	site    *Site
	handler HttpHandler
}

func (sitehandler *SiteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	site := sitehandler.site
	t1 := time.Now()

	if site.IsAuthEnabled() {
		if cookie, err := r.Cookie("_auth"); err == nil {
			value := make(map[string]string)
			if err = site.server.jar.Decode("_auth", cookie.Value, &value); err == nil {
				sitehandler.handler.ServeHTTP(sitehandler.site, w, r)
			} else {
				http.Redirect(w, r, "/_auth?redirectUrl="+r.URL.Path, http.StatusFound)
			}
		} else {
			http.Redirect(w, r, "/_auth?redirectUrl="+r.URL.Path, http.StatusFound)
		}
	} else {
		sitehandler.handler.ServeHTTP(sitehandler.site, w, r)
	}

	t2 := time.Now()
	log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
}
