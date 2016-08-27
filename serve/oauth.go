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

	if r.Method == "POST" {
		r.ParseForm()

		auth := new(Authentication)
		responseCode := values["response_code"]

		clientID := values["client_id"]
		redirectURI := values["redirect_uri"]

		username := r.Form["username"]
		password := r.Form["password"]
		state := values["state"]

		if len(responseCode) == 1 && responseCode[0] == "code" {
			if len(clientID) == 1 && len(redirectURI) == 1 && len(username) == 1 && len(password) == 1 {
				if auth.Authenticate(AuthParameter{username: username[0], password: password[0], clientID: clientID[0], redirectURI: redirectURI[0]}) >= 1 {
					query := "?code=code"
					if len(state) > 0 {
						query += "&state=" + state[0]
					}
					http.Redirect(w, r, redirectURI[0]+query, 302)
				}
			}
		}

		if len(responseCode) == 1 && responseCode[0] == "token" {
			if len(clientID) == 1 && len(redirectURI) == 1 && len(username) == 1 && len(password) == 1 {
				if auth.Authenticate(AuthParameter{username: username[0], password: password[0], clientID: clientID[0], redirectURI: redirectURI[0]}) >= 1 {
					query := "?access_token=access_token&expires_in=expires_in&refresh_token=refresh_token&issued_at=issued_at&signature=signature"
					if len(state) > 0 {
						query += "&state=" + state[0]
					}
					http.Redirect(w, r, redirectURI[0]+query, 302)
				}
			}
		}
	}

	if r.Method == "GET" {
		w.Write([]byte("login screen"))
	}
}

// Token token
func (oauth2 *OAuth2) Token(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	if r.Method == "POST" {
		auth := new(Authentication)

		grantType := values["grant_type"]
		clientID := values["client_id"]
		clientSecret := values["client_secret"]

		if len(grantType) == 1 && grantType[0] == "authorization_code" {
			redirectURI := values["redirect_uri"]
			code := values["code"]
			if len(clientID) == 1 && len(clientSecret) == 1 && len(redirectURI) == 1 && len(code) == 1 {
				if auth.Authenticate(AuthParameter{clientID: clientID[0], clientSecret: clientSecret[0], redirectURI: redirectURI[0], authorizationCode: code[0]}) >= 1 {
					res := map[string]interface{}{
						"access_token":  "access_token",
						"refresh_token": "refresh_token",
						"issued_at":     "issued_at",
						"signature":     "signature",
					}
					res2, _ := json.Marshal(res)
					w.WriteHeader(200)
					w.Write(res2)
				}
			}
		}

		if len(grantType) == 1 && grantType[0] == "password" {
			username := values["username"]
			password := values["password"]
			if len(username) == 1 && len(password) == 1 && len(clientID) == 1 && len(clientSecret) == 1 {
				if auth.Authenticate(AuthParameter{username: username[0], password: password[0], clientID: clientID[0], clientSecret: clientSecret[0]}) >= 1 {
					res := map[string]interface{}{
						"access_token": "access_token",
						"issued_at":    "issued_at",
						"signature":    "signature",
					}
					res2, _ := json.Marshal(res)
					w.WriteHeader(200)
					w.Write(res2)
				}
			}
		}

		if len(grantType) == 1 && grantType[0] == "refresh_token" {
			refreshToken := values["refresh_token"]
			if len(clientID) == 1 && len(clientSecret) == 1 {
				if auth.Authenticate(AuthParameter{clientID: clientID[0], clientSecret: clientSecret[0], refreshToken: refreshToken[0]}) >= 1 {
					res := map[string]interface{}{
						"access_token": "access_token",
						"issued_at":    "issued_at",
						"signature":    "signature",
					}
					res2, _ := json.Marshal(res)
					w.WriteHeader(200)
					w.Write(res2)
				}
			}
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
