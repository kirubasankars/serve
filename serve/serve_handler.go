package serve

import "net/http"

type serveHandler struct{}

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
	server := ctx.Server
	server.contexts[r] = &ctx
	fakeR, _ := http.NewRequest(r.Method, "/"+ctx.Module.Name+ctx.Path, nil)
	handler, _ := ctx.Module.mux.Handler(fakeR)
	//fmt.Println("/" + ctx.Module.Name + ctx.Path)
	fakeR.URL.Path = r.URL.Path
	handler.ServeHTTP(w, r)
	delete(server.contexts, r)
}
