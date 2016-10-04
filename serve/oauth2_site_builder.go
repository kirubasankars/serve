package serve

import "net/http"

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
		server.System.ServeFile(&ctx, w, r)
	}

	module.Handlers["/token"] = func(ctx Context, w http.ResponseWriter, r *http.Request) {
		hanlder.Token(w, r)
	}

}
