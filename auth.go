package main

import (
	"log"
	"net/http"
)

type AuthSiteHandler struct {
	file FileHandler
}

func (handler *AuthSiteHandler) ServeHTTP(site *Site, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/_auth" && r.Method == "POST" {
		server := site.server

		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}

		formData := make(map[string]string)
		for key, values := range r.Form { // range over map
			for _, value := range values { // range over []string
				formData[key] = value
			}
		}

		cookieValue := make(map[string]string)
		cookieValue["username"] = formData["username"]

		if encoded, err := server.jar.Encode("_auth", cookieValue); err == nil {
			cookie := &http.Cookie{
				Name:  "_auth",
				Value: encoded,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, formData["redirectUrl"], http.StatusFound)
		}
	} else {
		handler.file.ServeHTTP(site, w, r)
	}
}
