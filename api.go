package main

import (
	"encoding/json"
	"net/http"

	"github.com/serve/libs/metal"
)

func (site *Site) HandleAPI(model *metal.Metal, w http.ResponseWriter, r *http.Request) {

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
