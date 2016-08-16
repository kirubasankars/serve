package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
_	"fmt"
	"github.com/kirubasankars/serve/metal"
	"github.com/kirubasankars/serve/serve"
)

// APPS name of the folder
const APPS string = "apps"

// MODULES name of the folder
const MODULES string = "modules"

var re = regexp.MustCompile("[^A-Za-z0-9/.]+")

type statFunction func(path string) bool
type getConfigFunction func(path string) *metal.Metal

// FileSystem dads
type FileSystem struct {
	stat      statFunction
	getConfig getConfigFunction
}

// Build dads
func (fs *FileSystem) Build(ctx *serve.Context, uri string) {

	for {
		if strings.Contains(uri, "//") == false && strings.Contains(uri, "..") == false {
			break
		}
		uri = strings.Replace(uri, "//", "/", -1)
		uri = strings.Replace(uri, "..", ".", -1)
	}

	if re.ReplaceAllString(uri, "") != uri {
		return
	}

	parts := strings.Split(uri, "/")
	urlLen := len(parts) - 1
	l := 1
	currentIdx := 1
	
	if ctx.Namespace == nil {
		if currentIdx <= urlLen {
			fs.GetNamespace(ctx, parts[currentIdx])
		}
		
		if ctx.Namespace != nil {			
			currentIdx++
			l += len(ctx.Namespace.Name)
		} else {		
			fs.GetNamespace(ctx, ".")			
		}		
	}

	if ctx.Application == nil {
		if currentIdx <= urlLen {
			fs.GetApplication(ctx, parts[currentIdx])
		}
		if ctx.Application != nil {
			if currentIdx == 2 {
				l++
			}

			currentIdx++
			l += len(ctx.Application.Name)
		}
	}

	if ctx.Module == nil {

		if currentIdx <= urlLen {
			name := parts[currentIdx]

			if ctx.GetConfig("modules") != nil {
				l, _ := ctx.GetConfig("modules.$length").(int)
				for i := 0; i < l; i++ {
					mname, _ := ctx.GetConfig("modules.@" + strconv.Itoa(i)).(string)
					if mname == name {
						fs.GetModule(ctx, name)
					}
				}
			}

		}

		if ctx.Module != nil {
			if currentIdx == 2 || currentIdx == 3 {
				l++
			}
			currentIdx++
			l += len(ctx.Module.Name)
		} else {
			if ctx.GetConfig("modules.@0") != nil {
				name, _ := ctx.GetConfig("modules.@0").(string)
				fs.GetModule(ctx, name)
			} else {
				fs.GetModule(ctx, "home")
			}
		}
	}

	if l == 1 {
		l = 0
	}	

	ctx.Path = uri[l:]

	//fmt.Println(ctx.Namespace, ctx.Application, ctx.Module, ctx.Path)
}

// GetNamespace dad
func (fs *FileSystem) GetNamespace(ctx *serve.Context, name string) {
	if name == "" || name == APPS || name == MODULES || strings.Contains(name, "..") {
		return
	}

	server := ctx.Server

	if ns := ctx.Server.Namespaces[name]; ns != nil {
		ctx.Namespace = ns
		return
	}

	server.Lock()
	defer server.Unlock()

	if ns := ctx.Server.Namespaces[name]; ns != nil {
		ctx.Namespace = ns
		return
	}

	loc := filepath.Join(ctx.Server.Path(), name)		
	if fs.stat(loc) {		
		ns := serve.NewNamespace(name, name, fs.getConfig(loc), server)
		ns.Build()
		server.Namespaces[name] = ns
		ctx.Namespace = ns
	}
}

// GetApplication dad
func (fs *FileSystem) GetApplication(ctx *serve.Context, name string) {
	if name == "" || name == APPS || name == MODULES || strings.Contains(name, "..") {
		return
	}

	server := ctx.Server
	namespace := ctx.Namespace
	apps := namespace.Applications
	if app := apps[name]; app != nil {
		ctx.Application = app
		return
	}

	namespace.Lock()
	defer namespace.Unlock()

	if app := apps[name]; app != nil {
		ctx.Application = app
		return
	}

	path := filepath.Join(ctx.Namespace.Path, APPS, name)
	loc := filepath.Join(server.Path(), path)
	if fs.stat(loc) {
		app := serve.NewApplication(name, path, fs.getConfig(loc), server)
		app.Build()
		apps[name] = app
		ctx.Application = app
	}
}

// GetModule dad
func (fs *FileSystem) GetModule(ctx *serve.Context, name string) {
	if name == "" || name == APPS || name == MODULES || strings.Contains(name, "..") {
		return
	}

	server := ctx.Server
	namespace := ctx.Namespace
	modules := ctx.Namespace.Modules

	var module *serve.Module

	if module = modules[name]; module != nil {
		ctx.Module = module
		return
	}

	namespace.Lock()
	defer namespace.Unlock()

	if module = modules[name]; module != nil {
		ctx.Module = module
		return
	}

	path := filepath.Join(ctx.Namespace.Path, MODULES, name)
	loc := filepath.Join(server.Path(), path)
	if fs.stat(loc) {
		modules[name] = serve.NewModule(name, path, fs.getConfig(loc), server)
	} else {
		path = filepath.Join(MODULES, name)
		loc = filepath.Join(server.Path(), path)
		if fs.stat(loc) {
			modules[name] = serve.NewModule(name, path, fs.getConfig(loc), server)
		}
	}

	module = modules[name]
	if module != nil {
		appFolder := filepath.Join(ctx.Server.Path(), module.Path, "app")
		if fs.stat(appFolder) {
			module.AppEnabled = true
		}
	}

	if module != nil {
		module.Build()
		ctx.Module = module
	}
}

//LoadConfig used for load config
func LoadConfig(path string) *metal.Metal {
	if s, err := ioutil.ReadFile(filepath.Join(path, "config.json")); err == nil {
		m := metal.NewMetal()
		m.Parse(s)
		return m
	}
	return nil
}

// NewFileSystem dadsasf
func NewFileSystem(stat statFunction, getConfig getConfigFunction) *FileSystem {
	fs := new(FileSystem)
	fs.stat = stat
	fs.getConfig = getConfig
	return fs
}

// Stat aas
func Stat(path string) bool {
	if f, _ := os.Stat(path); f != nil {
		return true
	}
	return false
}
