package serve

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kirubasankars/metal"
)

type ApiHandler struct {
	site *Site
}

func (handler *ApiHandler) ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/api/")

	var input *metal.Metal
	json, _ := ioutil.ReadAll(r.Body)
	if len(json) > 0 {
		input = metal.NewMetal()
		input.Parse(json)
	}
	model := site.Model(r.Method, parts[1], input)
	handler.handle(model, w, r)
}

func (handler *ApiHandler) handle(model *metal.Metal, w http.ResponseWriter, r *http.Request) {
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
