package serve

import (
	"net/http"
	"path/filepath"
)

type commonHandlerBuilder struct{}

func (chb *commonHandlerBuilder) Build(module Module) {

	module.Handlers["/"] = func(ctx Context, w http.ResponseWriter, r *http.Request) {
		//fmt.Println(ctx.GetConfig("modules"))
		server := ctx.Server
		appFolder := ""
		if ctx.Module.AppEnabled {
			appFolder = "app"
		}
		path2file := filepath.Join(ctx.Module.Path, appFolder, ctx.Path)
		//fmt.Println(path2file)
		server.ServeFile(w, r, path2file)
	}

}
