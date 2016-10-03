package serve

import (
	"net/http"
	"path/filepath"
)

//OAuth2Builder oauth2builder
type OAuth2Builder struct{}

//Build oauth2 builder
func (oauth2 *OAuth2Builder) Build(module Module) {
	hanlder := new(OAuth2)
	hanlder.server = module.server

	module.Handlers["/authorize"] = func(ctx Context, w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			hanlder.Authorize(w, r)
			return
		}
	}

	module.Handlers["/authorize/"] = func(ctx Context, w http.ResponseWriter, r *http.Request) {

		server := ctx.Server
		appFolder := ""
		if ctx.Module.AppEnabled {
			appFolder = "app"
		}
		path2file := filepath.Join(ctx.Module.Path, appFolder, ctx.Path)
		server.ServeFile(w, r, path2file)
	}

	module.Handlers["/token"] = func(ctx Context, w http.ResponseWriter, r *http.Request) {

		hanlder.Token(w, r)
	}

}
