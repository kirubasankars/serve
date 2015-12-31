package main

import (
	"net/http"
	"strings"
)

type ApiServe struct {
	site *Site
}

func (apiServe *ApiServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var site = apiServe.site
	parts := strings.Split(r.URL.Path, "/api")
	model := site.GetModel(parts[1] + "/" + strings.ToLower(r.Method))
	site.HandleAPI(model, w, r)
}
