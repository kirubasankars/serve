package serve

import (
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
