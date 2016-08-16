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

// GetConfig get config data from module/app/namespace/server
func (ctx *Context) GetConfig(key string) interface{} {

	if module := ctx.Module; module != nil {
		if value := module.GetConfig(key); value != nil {
			return value
		}
	}

	if app := ctx.Application; app != nil {
		if value := app.GetConfig(key); value != nil {
			return value
		}
	}

	if ns := ctx.Namespace; ns != nil {
		if value := ns.GetConfig(key); value != nil {
			return value
		}
	}

	return nil
}

// NewContext for create new context
func newContext(server *Server, r *http.Request) *Context {
	ctx := new(Context)
	url := r.URL.Path
	ctx.URL = url
	ctx.Server = server
	ctx.User = newUser("")

	system := server.System
	system.Build(ctx, url)

	auth := new(authenticator)
	auth.Validate(ctx, r)

	return ctx
}
