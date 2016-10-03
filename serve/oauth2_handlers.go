package serve

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OAuth2 object
type OAuth2 struct {
	server *Server
}

// Authorize authorize
func (oauth2 *OAuth2) Authorize(w http.ResponseWriter, r *http.Request) {
	server := oauth2.server
	authProvider := server.System.GetOAuthProvider(server)
	values := r.URL.Query()

	r.ParseForm()

	responseCode := values["response_code"]

	clientID := values["client_id"]
	redirectURI := values["redirect_uri"]

	username := r.Form["username"]
	password := r.Form["password"]
	state := values["state"]

	if len(responseCode) == 1 && responseCode[0] == "code" {
		if len(clientID) == 1 && len(redirectURI) == 1 && len(username) == 1 && len(password) == 1 {
			ws := authProvider.GetWebserver()
			code := ws.GetAuthorizationCode(clientID[0], redirectURI[0], username[0], password[0])
			query := "?code=" + code
			if len(state) > 0 {
				query += "&state=" + state[0]
			}
			http.Redirect(w, r, "/oauth/code_callback"+query+"&redirect_uri="+redirectURI[0], 302)
		}
	}

	if len(responseCode) == 1 && responseCode[0] == "token" {
		if len(clientID) == 1 && len(redirectURI) == 1 && len(username) == 1 && len(password) == 1 {
			up := authProvider.GetUserAgent()
			t := up.GetAccessToken(clientID[0], redirectURI[0], username[0], password[0])
			query := "?access_token=" + t.Token + "&issued_at=" + "issuedAt"
			if len(state) > 0 {
				query += "&state=" + state[0]
			}
			http.Redirect(w, r, "/"+redirectURI[0]+query, 302)
		}
	}
}

// Token token
func (oauth2 *OAuth2) Token(w http.ResponseWriter, r *http.Request) {
	server := oauth2.server
	authProvider := server.System.GetOAuthProvider(server)

	values := r.URL.Query()
	if r.Method == "POST" {
		grantType := values["grant_type"]
		clientID := values["client_id"]
		clientSecret := values["client_secret"]

		if len(grantType) == 1 && grantType[0] == "authorization_code" {
			redirectURI := values["redirect_uri"]
			code := values["code"]
			if len(clientID) == 1 && len(clientSecret) == 1 && len(redirectURI) == 1 && len(code) == 1 {
				ws := authProvider.GetWebserver()
				t := ws.GetAccessToken(code[0], clientID[0], clientSecret[0], redirectURI[0])
				res := map[string]string{
					"access_token": t.Token,
					"issued_at":    "issued_at",
					"signature":    "signature",
				}
				res2, _ := json.Marshal(res)
				w.WriteHeader(200)
				w.Write(res2)
			}
		}

		fmt.Println("fdasfsddasdasd")
		if len(grantType) == 1 && grantType[0] == "password" {
			username := values["username"]
			password := values["password"]
			if len(username) == 1 && len(password) == 1 && len(clientID) == 1 && len(clientSecret) == 1 {
				up := authProvider.GetUserPassword()
				t := up.GetAccessToken(clientID[0], clientSecret[0], username[0], password[0])
				res := map[string]string{
					"access_token": t.Token,
					"issued_at":    "issued_at",
					"signature":    "signature",
				}
				res2, _ := json.Marshal(res)
				w.WriteHeader(200)
				w.Write(res2)
			}
		}

		if len(grantType) == 1 && grantType[0] == "refresh_token" {
			refreshToken := values["refresh_token"]
			if len(clientID) == 1 && len(clientSecret) == 1 {
				rt := authProvider.GetRefreshToken()
				res := rt.GetAccessToken(refreshToken[0], clientID[0], clientSecret[0])
				res2, _ := json.Marshal(res)
				w.WriteHeader(200)
				w.Write(res2)
			}
		}
	}
}
