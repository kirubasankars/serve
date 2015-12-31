package main

import (
	"log"
	"net/http"
	"time"
)

type Site struct {
	name string
	path string
	uri  string

	server *Server
}

func (site *Site) Build() {
	var mux = site.server.http

	mux.Handle(site.uri, Handler(&FileServe{site}))
	mux.Handle(site.uri+"api/", Handler(&ApiServe{site}))
	mux.Handle(site.uri+"html/", Handler(&HtmlServe{site}))
	mux.Handle(site.uri+"text/", Handler(&TextServe{site}))
}

func Handler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
	h := http.HandlerFunc(fn)
	return loggingHandler(h)
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}
