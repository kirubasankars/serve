package driver

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/kirubasankars/serve/serve"
)

// APPS name of the folder
const APPS string = "apps"

// MODULES name of the folder
const MODULES string = "modules"

var re = regexp.MustCompile("[^A-Za-z0-9/._]+")

type statFunction func(path string) bool
type getConfigFunction func(path string) *[]byte
type serveFileFunction func(ctx *serve.Context, w http.ResponseWriter, r *http.Request)

// FileSystem dads
type FileSystem struct {
	stat      statFunction
	getConfig getConfigFunction
	serveFile serveFileFunction
}

//GetConfig get config
func (fs *FileSystem) GetConfig(path string) *[]byte {
	return fs.getConfig(path)
}

//ServeFile - serveFile
func (fs *FileSystem) ServeFile(ctx *serve.Context, w http.ResponseWriter, r *http.Request) {
	fs.serveFile(ctx, w, r)
}

//ServeFile - serveFile
func ServeFile(ctx *serve.Context, w http.ResponseWriter, r *http.Request) {
	appFolder := ""
	if ctx.Module.AppEnabled {
		appFolder = "app"
	}
	path := filepath.Join(ctx.Server.Path(), ctx.Module.Path, appFolder, ctx.Path)
	http.ServeFile(w, r, path)
}

//LoadConfig used for load config
func LoadConfig(path string) *[]byte {
	if s, err := ioutil.ReadFile(path); err == nil {
		return &s
	}
	return nil
}

// NewFileSystem dadsasf
func NewFileSystem(stat statFunction, getConfig getConfigFunction, serveFile serveFileFunction) *FileSystem {
	fs := new(FileSystem)
	fs.stat = stat
	fs.getConfig = getConfig
	fs.serveFile = serveFile
	return fs
}

// Stat aas
func Stat(path string) bool {
	if f, _ := os.Stat(path); f != nil {
		return true
	}
	return false
}
