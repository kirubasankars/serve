package serve

import (
	"fmt"
	"log"
	"net/http"
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
	AccessTokens   map[string]*AccessToken

	serveHandler
	sync.Mutex // <-- this mutex protects

	contexts map[*http.Request]*Context
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

	handler, _ := server.mux.Handler(r)
	handler.ServeHTTP(w, r)

	t2 := time.Now()

	log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
}

// Serve serve
func (server *Server) serve(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	if url == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	ctx := NewContext(server, r)
	if ctx == nil || ctx.Module == nil {
		http.NotFound(w, r)
		return
	}

	server.serveHandler.ServeHTTP(*ctx, w, r)
}

// NewServer for create new server
func NewServer(port string, rootPath string, driver System) *Server {
	server := new(Server)

	server.path = rootPath
	server.port = port
	server.System = driver
	server.AccessTokens = make(map[string]*AccessToken)
	server.contexts = make(map[*http.Request]*Context)
	server.moduleProvider = make(map[string]ModuleHandlerProvider)
	server.Namespaces = make(map[string]*Namespace)

	server.mux = http.NewServeMux()
	server.mux.HandleFunc("/", server.serve)

	server.RegisterProvider(".", new(commonHandlerBuilder))
	server.RegisterProvider("_oauth2", new(OAuth2Builder))

	return server
}

//RegisterProvider to register handler provider
func (server *Server) RegisterProvider(name string, provider ModuleHandlerProvider) {
	server.moduleProvider[name] = provider
}
