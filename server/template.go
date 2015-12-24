package server

import (
	htmlTeamplate "html/template"
	"net/http"
	"path"
	"strings"
	textTeamplate "text/template"

	"github.com/serve/metal"
)

func HandleTemplate(site *Site, w http.ResponseWriter, r *http.Request) {
	m := metal.NewMetal()
	m.Set("Profile.Name", "Kiruba")
	m.Set("Profile.Hobbies.@0", "Programming")
	m.Set("Profile.Hobbies.@1", "Music")

	paths := strings.Split(r.URL.Path, "tpl/")

	if strings.Contains(r.URL.Path, "html") {
		w.Header().Set("Content-Type", "text/html")
		templatePath := path.Join(site.path, "tpl", paths[1], "get.html")
		handleHTML(templatePath, m.Raw(), w)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		templatePath := path.Join(site.path, "tpl", paths[1], "get.txt")
		handleText(templatePath, m.Raw(), w)
	}
}

func handleHTML(path string, model interface{}, w http.ResponseWriter) {
	tmpl, err := htmlTeamplate.ParseFiles(path)

	if err != nil {
		http.Error(w, "No Template found.", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleText(path string, model interface{}, w http.ResponseWriter) {
	tmpl, err := textTeamplate.ParseFiles(path)

	if err != nil {
		http.Error(w, "No Template found.", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, model); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
