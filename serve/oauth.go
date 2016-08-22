package serve

import (
	"encoding/json"
	"net/http"
)

// OAuth2 object
type OAuth2 struct{}

// Authorize authorize
func (oauth2 *OAuth2) Authorize(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	grantType := values["grant_type"]
	if len(grantType) == 1 && grantType[0] == "pa	ssword" {
		userName := values["username"]
		password := values["password"]
		clientID := values["client_id"]
		clientSecret := values["client_secret"]
		if len(userName) == 1 && len(password) == 1 && len(clientID) == 1 && len(clientSecret) == 1 {
			w.Write([]byte("authorize"))
		}
	}
}

// Token token
func (oauth2 *OAuth2) Token(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	auth := new(Authentication)
	grantType := values["grant_type"]
	if len(grantType) == 1 && grantType[0] == "password" {
		clientID := values["client_id"]
		clientSecret := values["client_secret"]
		username := values["username"]
		password := values["password"]
		if len(username) == 1 && len(password) == 1 && len(clientID) == 1 && len(clientSecret) == 1 {
			if auth.Authenticate(username[0], password[0], clientID[0], clientSecret[0]) {
				res := map[string]interface{}{
					"access_token": "access_token",
					"id":           "id",
					"issued_at":    "issued_at",
					"signature":    "signature",
				}
				res2, _ := json.Marshal(res)
				w.WriteHeader(200)
				w.Write(res2)
			} else {
				w.WriteHeader(301)
				w.Write([]byte("invalid request"))
			}
		} else {
			w.WriteHeader(301)
			w.Write([]byte("invalid request"))
		}
	}
}

// Register register oauth handlers
func (oauth2 *OAuth2) Register(mux *http.ServeMux) {
	mux.HandleFunc("/oauth2/authorize", func(w http.ResponseWriter, r *http.Request) {
		oauth2.Authorize(w, r)
	})
	mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		oauth2.Token(w, r)
	})
}
