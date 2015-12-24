package server

import (
	"encoding/json"
	"net/http"

	"github.com/serve/metal"
)

func HandleAPI(site *Site, w http.ResponseWriter, r *http.Request) {
	m := metal.NewMetal()
	m.Set("Profile.Name", "Kiruba")
	m.Set("Profile.Hobbies.@0", "Programming")
	m.Set("Profile.Hobbies.@1", "Music")
	//paths := strings.Split(r.URL.Path, "/api")
	rs, err := json.Marshal(m.Raw())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(rs)
}
