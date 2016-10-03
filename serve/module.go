package serve

import (
	"encoding/json"
	"net/http"
	"regexp"
)

// HTTPHandler its handler for serve
type HTTPHandler func(ctx Context, w http.ResponseWriter, r *http.Request)

// Module for serve
type Module struct {
	Name   string
	Path   string
	Config *ModuleConfigration

	AppEnabled  bool
	AuthEnabled bool
	Handlers    map[string]HTTPHandler

	permissions     map[string][]string
	permittedRoutes map[string]*regexp.Regexp

	mux    *http.ServeMux
	server *Server
}

// Build build a module
func (module *Module) Build() {
	server := module.server

	// if module.Config != nil {
	// 	permissions := module.Config.Permissions
	// 	if permissions != nil {
	// 		if module.permissions == nil {
	// 			module.permissions = make(map[string][]string)
	// 		}
	//
	// 		for permission, values := range permissions {
	// 			exp := ""
	// 			for idx := range values {
	// 				auth := values[idx]
	// 				le := len(auth)
	// 				if le > 6 && auth[0:4] == "url(" && auth[le-1:] == ")" {
	// 					exp += "^(" + auth[4:le-1] + ")$|"
	// 				} else {
	// 					if _, p := module.permissions[permission]; p == false {
	// 						module.permissions[permission] = make([]string, 0)
	// 					}
	// 					module.permissions[permission] = append(module.permissions[permission], auth)
	// 				}
	// 			}
	// 			exp = strings.TrimSuffix(exp, "|")
	// 			if len(exp) > 0 {
	// 				if module.permittedRoutes == nil {
	// 					module.permittedRoutes = make(map[string]*regexp.Regexp)
	// 				}
	// 				module.permittedRoutes[permission] = regexp.MustCompile(exp)
	// 			}
	// 		}
	// 	}
	// }

	if provider, p := server.moduleProvider["."]; p {
		provider.Build(*module)
	}

	if provider, p := server.moduleProvider[module.Name]; p {
		provider.Build(*module)
	}

	mux := module.mux

	for pattern, handler := range module.Handlers {
		if pattern != "" {
			var mh = new(moduleHandler)
			mh.handler = handler
			mh.module = module
			uri := "/" + module.Name + pattern

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
	server := mh.module.server
	//TODO: need to refactor.
	mh.handler(*server.contexts[r], w, r)
}

// NewModule create module
func NewModule(name string, path string, config *[]byte, server *Server) *Module {
	module := new(Module)
	module.Name = name
	module.Path = path

	if config != nil {
		nc := new(ModuleConfigration)
		if err := json.Unmarshal(*config, &nc); err != nil {
			panic(err)
		}
		module.Config = nc
	}

	module.server = server

	module.Handlers = make(map[string]HTTPHandler)
	module.mux = http.NewServeMux()

	return module
}
