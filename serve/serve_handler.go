package serve

import "net/http"

type serveHandler struct{}

var contexts map[*http.Request]*Context

func init() {
	contexts = make(map[*http.Request]*Context)
}

func (serveHandler *serveHandler) ServeHTTP(ctx Context, w http.ResponseWriter, r *http.Request) {

	if ctx.Module == nil {
		http.NotFound(w, r)
		return
	}

	// if ctx.Module.Name == "_refresh" {
	// 	delete(ctx.Server.Namespaces, ctx.Namespace.Name)
	// 	return
	// }

	serveHandler.Serve(ctx, w, r)
}

func (serveHandler *serveHandler) Serve(ctx Context, w http.ResponseWriter, r *http.Request) {
	contexts[r] = &ctx
	fakeR, _ := http.NewRequest(r.Method, "/"+ctx.Module.Name+ctx.Path, nil)
	handler, _ := ctx.Module.mux.Handler(fakeR)
	fakeR.URL.Path = r.URL.Path
	handler.ServeHTTP(w, r)
	delete(contexts, r)
}
