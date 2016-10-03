package serve

import "net/http"

// Context for serve
type Context struct {
	Server      *Server
	Namespace   *Namespace
	Application *Application
	Module      *Module

	User *User

	URL  string
	Path string
}

// NewContext for create new context
func NewContext(server *Server, r *http.Request) *Context {
	ctx := new(Context)
	url := r.URL.Path
	ctx.URL = url
	ctx.Server = server
	ctx.User = NewUser("")

	system := server.System
	system.Build(ctx, url)

	//auth := new(routeAuthenticator)
	//auth.Validate(ctx, r)

	return ctx
}
