package serve

import (
	"encoding/json"
	"sync"
)

//Namespace sturct
type Namespace struct {
	Name   string
	Path   string
	Config *NamespaceConfigration

	Modules      map[string]*Module
	Applications map[string]*Application

	server     *Server
	sync.Mutex // <-- this mutex protects

	roles map[string][]string
}

// GetConfig get config from namespace
// func (ns *Namespace) GetConfig(key string) interface{} {
// 	if ns.config == nil {
// 		return nil
// 	}
// 	return ns.config.Get(key)
// }

// Build builing namespace
func (ns *Namespace) Build() {

}

// NewNamespace create namespace
func NewNamespace(name string, path string, config *[]byte, server *Server) *Namespace {
	namespace := new(Namespace)
	namespace.Name = name
	namespace.Path = path

	if config != nil {
		nc := new(NamespaceConfigration)
		if err := json.Unmarshal(*config, &nc); err != nil {
			panic(err)
		}
		namespace.Config = nc
	}

	namespace.Applications = make(map[string]*Application)
	namespace.Modules = make(map[string]*Module)
	namespace.server = server

	return namespace
}
