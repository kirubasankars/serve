package serve

import (
	"regexp"
	"sync"

	"github.com/kirubasankars/serve/metal"
)

//Namespace sturct
type Namespace struct {
	Name   string
	Path   string
	config *metal.Metal

	Modules      map[string]*Module
	Applications map[string]*Application

	server     *Server
	sync.Mutex // <-- this mutex protects

	roles map[string][]string
}

// GetConfig get config from namespace
func (ns *Namespace) GetConfig(key string) interface{} {
	if ns.config == nil {
		return nil
	}
	return ns.config.Get(key)
}

// Build builing namespace
func (ns *Namespace) Build() {
	if roles, _ := ns.GetConfig("roles").(*metal.Metal); roles != nil {
		props := roles.Properties()
		if ns.roles == nil {
			ns.roles = make(map[string][]string, len(props))
		}
		for name := range props {
			if role, done := roles.Get(name).(*metal.Metal); done == true {
				for _, v := range role.Properties() {
					if auth, done := v.(string); done == true {
						if _, p := ns.roles[name]; p == true {
							ns.roles[name] = make([]string, 0)
						}
						ns.roles[name] = append(ns.roles[name], auth)
					}
				}
			}
		}
	}
}

// le := len(auth)
// if le > 4 && auth[0:4] == "url(" && auth[le-1:] == ")" {
// 	if _, p := ns.authURL[name]; p == false {
// 		ns.authURL[name] = make([]string, 0)
// 	}
// 	rr, _ := ns.authURL[name]
// 	ns.authURL[name] = append(rr, auth[4:le-1])
// }

// ValidateURL validate URL
func (ns *Namespace) ValidateURL(ctx *Context, authString string) bool {

	if roles, _ := ns.GetConfig("roles").(*metal.Metal); roles != nil {
		for r := range roles.Properties() {
			if ctx.UserInRole(r) {
				if role, done := roles.Get(r).(*metal.Metal); done == true {
					for _, v := range role.Properties() {
						if auth, done := v.(string); done == true {
							le := len(auth)
							if le > 4 && auth[0:4] == "url(" && auth[le-1:] == ")" {
								r1 := auth[4 : le-1]
								if p, _ := regexp.MatchString(r1, authString); p == true {
									return true
								}
								return false
							}
						}
					}
				}
			}
		}
	}
	return true
}

// NewNamespace create namespace
func NewNamespace(name string, path string, config *metal.Metal, server *Server) *Namespace {
	namespace := new(Namespace)
	namespace.Name = name
	namespace.Path = path
	namespace.config = config

	namespace.Applications = make(map[string]*Application)
	namespace.Modules = make(map[string]*Module)
	namespace.server = server

	return namespace
}
