package serve

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	//	"strings"
	"sync"
	"time"
)

// Server - manage server operations
type Server struct {
	port string
	path string

	mux    *http.ServeMux
	System System

	Namespaces     map[string]*Namespace
	moduleProvider map[string]ModuleHandlerProvider

	serveHandler
	sync.Mutex // <-- this mutex protects
}

// Path get path
func (server *Server) Path() string {
	return server.path
}

// Start new server
func (server *Server) Start() {
	if err := http.ListenAndServe("localhost:"+server.port, server.mux); err != nil {
		fmt.Println(err)
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t1 := time.Now()

	url := r.URL.Path

	if url == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	ctx := newContext(server, r)
	if ctx == nil || ctx.Module == nil {
		fmt.Println(r.URL.Path)
		http.NotFound(w, r)
		return
	}

	server.serveHandler.ServeHTTP(*ctx, w, r)

	t2 := time.Now()

	log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
}

// NewServer for create new server
func NewServer(port string, rootPath string, driver System) *Server {
	server := new(Server)
	server.port = port
	server.path = rootPath
	server.System = driver

	server.mux = http.NewServeMux()
	server.mux.Handle("/", server)

	server.moduleProvider = make(map[string]ModuleHandlerProvider)
	server.Namespaces = make(map[string]*Namespace)

	server.RegisterProvider(".", new(commonHandlerBuilder))

	return server
}

//RegisterProvider to register handler provider
func (server *Server) RegisterProvider(name string, provider ModuleHandlerProvider) {
	server.moduleProvider[name] = provider
}

// ServeFile serve file
func (server *Server) ServeFile(w http.ResponseWriter, r *http.Request, file string) {
	path := filepath.Join(server.Path(), file)
	//fmt.Println(path)
	http.ServeFile(w, r, path)
}
