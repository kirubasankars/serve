package serve

import (
	"fmt"
	"net/http"

	"github.com/kirubasankars/serve/metal"
)

// HTTPHandler its handler for serve
type HTTPHandler func(ctx Context, w http.ResponseWriter, r *http.Request)

// Module for serve
type Module struct {
	Name   string
	Path   string
	config *metal.Metal

	AppEnabled  bool
	AuthEnabled bool
	Handlers    map[string]HTTPHandler

	permissions map[string][]string

	mux    *http.ServeMux
	server *Server
}

// GetConfig get config from module
func (module *Module) GetConfig(key string) interface{} {
	if module.config == nil {
		return nil
	}
	return module.config.Get(key)
}

// Build build a module
func (module *Module) Build() {
	server := module.server

	if permissions, _ := module.GetConfig("permissions").(*metal.Metal); permissions != nil {
		props := permissions.Properties()
		if module.permissions == nil {
			module.permissions = make(map[string][]string, len(props))
		}

		for name := range props {
			if permission, done := permissions.Get(name).(*metal.Metal); done == true {
				for _, v := range permission.Properties() {
					if auth, done := v.(string); done == true {
						if _, p := module.permissions[name]; p == true {
							module.permissions[name] = make([]string, 0)
						}
						module.permissions[name] = append(module.permissions[name], auth)
					}
				}
			}
		}

		fmt.Println(module.permissions)
	}

	if provider, p := server.moduleProvider["."]; p {
		provider.Build(*module)
	}

	if provider, p := server.moduleProvider[module.Name]; p {
		provider.Build(*module)
	}

	mux := module.mux

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, r.URL.Path+"/"+module.Name+"/", 301)
	// })

	for pattern, handler := range module.Handlers {

		if pattern != "" {
			var mh = new(moduleHandler)
			mh.handler = handler
			mh.module = module
			uri := "/" + module.Name + pattern

			//fmt.Println("/"+module.Name+pattern, "build")

			mux.Handle(uri, mh)
		}
	}

	mux.HandleFunc("/"+module.Name, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.URL.Path+"/", 301)
	})
}

type moduleHandler struct {
	handler HTTPHandler
	module  *Module
}

func (mh *moduleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//ctx := newContext(mh.module.server, r.URL.Path)
	//fmt.Println(uri)
	mh.handler(*contexts[r], w, r)
}

// NewModule create module
func NewModule(name string, path string, config *metal.Metal, server *Server) *Module {
	module := new(Module)
	module.Name = name
	module.Path = path
	module.config = config
	module.server = server

	module.Handlers = make(map[string]HTTPHandler)
	module.mux = http.NewServeMux()

	return module
}
