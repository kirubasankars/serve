package serve

import "net/http"

// Authentication forms authentication
type Authentication struct{}

// Authenticate authenticate with username and password
func (auth *Authentication) Authenticate(username string, password string, clientID string, clientSecret string) bool {
	if username == "admin" && password == "admin" && clientID == "client_id" && clientSecret == "client_secret" {
		return true
	}
	return false
}

// SetAuthCookie setAuthCookie
func (auth *Authentication) SetAuthCookie(username string, ctx Context, w http.ResponseWriter) {

	// cookieValue := make(map[string]string)
	// cookieValue["username"] = username
	//
	// jar, path := ctx.GetJar()
	//
	// if encoded, err := jar.Encode("_auth", cookieValue); err == nil {
	// 	cookie := &http.Cookie{
	// 		Name:  "_auth",
	// 		Value: encoded,
	// 		Path:  path,
	// 	}
	// 	http.SetCookie(w, cookie)
	// } else {
	// 	fmt.Println(err)
	// }
}

// Signout signout
func (auth *Authentication) Signout(ctx *Context, w http.ResponseWriter, r *http.Request) {
	// _, p := ctx.GetJar()
	// cookie := &http.Cookie{
	// 	Name:   "_auth",
	// 	Value:  "",
	// 	Path:   p,
	// 	MaxAge: -1,
	// }
	// http.SetCookie(w, cookie)
}

// RedirectToLogin Redirect To Login
func (auth *Authentication) RedirectToLogin(ctx Context, w http.ResponseWriter, r *http.Request) {
	// _, path := ctx.GetJar()
	// http.Redirect(w, r, path+"/_auth/?redirectUrl="+r.URL.Path, http.StatusFound)
}

// Validate Validate
func (auth *Authentication) Validate(ctx *Context, w http.ResponseWriter, r *http.Request) bool {
	// jar, _ := ctx.GetJar()
	// if cookie, err := r.Cookie("_auth"); err == nil {
	// 	value := make(map[string]string)
	// 	if err = jar.Decode("_auth", cookie.Value, &value); err == nil {
	// 		return true
	// 	}
	// }
	// return false
	return true
}
