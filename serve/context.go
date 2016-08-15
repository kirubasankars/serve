package serve

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

// UserInRole check weather user part this role
func (ctx *Context) UserInRole(role string) bool {
	roles := *ctx.User.Roles
	for idx := range roles {
		if role == roles[idx] {
			return true
		}
	}
	return false
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
func newContext(server *Server, url string) *Context {

	ctx := new(Context)
	ctx.URL = url
	ctx.Server = server

	system := server.System
	system.Build(ctx, url)

	return ctx
}
