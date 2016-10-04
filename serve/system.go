package serve

import "net/http"

//System inferface for driver
type System interface {
	GetConfig(path string) *[]byte
	Build(ctx *Context, uri string)
	GetOAuthProvider(server *Server) OAuth2Provider
	ServeFile(ctx *Context, w http.ResponseWriter, r *http.Request)
}
