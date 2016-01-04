package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/serve/lib/metal"
)

type ApiHandler struct {
	site *Site
}

func (apiHandler *ApiHandler) ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/api")
	model := site.Model(parts[1] + "/" + strings.ToLower(r.Method))
	apiHandler.Handle(model, w, r)
}

func (apiHandler *ApiHandler) Handle(model *metal.Metal, w http.ResponseWriter, r *http.Request) {
	if model == nil {
		http.Error(w, "No api found.", http.StatusInternalServerError)
		return
	}
	rs, err := json.Marshal(model.Raw())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(rs)
}
