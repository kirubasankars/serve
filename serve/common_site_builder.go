package serve

import "net/http"

type commonHandlerBuilder struct{}

func (chb *commonHandlerBuilder) Build(module Module) {

	module.Handlers["/"] = func(ctx Context, w http.ResponseWriter, r *http.Request) {
		server := ctx.Server
		server.System.ServeFile(&ctx, w, r)
	}
}
