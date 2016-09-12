package driver

import (
	"path/filepath"
	"strings"

	"github.com/kirubasankars/serve/serve"
)

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
		} else {
			fs.GetApplication(ctx, ".")
		}
	}

	if ctx.Module == nil {

		var modules []string
		appConfig := ctx.Application.Config
		if appConfig != nil {
			modules = appConfig.Modules
		}

		if currentIdx <= urlLen {
			name := parts[currentIdx]
			for idx := range modules {
				if modules[idx] == name {
					fs.GetModule(ctx, name)
					break
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
			if len(modules) >= 1 {
				fs.GetModule(ctx, modules[0])
			} else {
				fs.GetModule(ctx, "home")
			}
		}
	}

	if l == 1 {
		l = 0
	}

	ctx.Path = uri[l:]
	//TODO ::
	//ctx.User.Roles = append(ctx.User.Roles, "admin")
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
		ns := serve.NewNamespace(name, name, nil, server)
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

	var path string
	if name == "." {
		path = filepath.Join(ctx.Namespace.Path, name)
	} else {
		path = filepath.Join(ctx.Namespace.Path, APPS, name)
	}

	loc := filepath.Join(server.Path(), path)
	if fs.stat(loc) {
		locC := filepath.Join(loc, "config.json")
		app := serve.NewApplication(name, path, fs.getConfig(locC), server)
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
		locC := filepath.Join(loc, "config.json")
		modules[name] = serve.NewModule(name, path, fs.getConfig(locC), server)
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
