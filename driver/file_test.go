package driver_test

import (
	"net/http"
	"testing"

	"github.com/kirubasankars/serve/driver"
	"github.com/kirubasankars/serve/serve"
)

// func TestBuildNamespace(t *testing.T) {
//
// 	fs := driver.NewFileSystem(func(path string) bool {
// 		if path == "/serve/kiruba" || path == "/serve/kiruba/apps/app" || path == "/serve/kiruba/modules/module" {
// 			return true
// 		}
// 		return false
// 	})
//
// 	ctx := new(serve.Context)
// 	ctx.Server = serve.NewServer("", "/serve", fs)
// 	ctx.URL = "/kiruba/app///module/path/to/file"
//
// 	fs.Build(ctx, ctx.URL)
//
// 	if !(ctx.Namespace.Name == "kiruba" && ctx.Application.Name == "app" && ctx.Module.Name == "module") {
// 		t.Fail()
// 	}
//
// }
//
// func TestBuildNamespaceDefault(t *testing.T) {
// 	fs := driver.NewFileSystem(func(path string) bool {
// 		if path == "/serve" || path == "/serve/apps/app" || path == "/serve/modules/home" {
// 			return true
// 		}
// 		return false
// 	})
//
// 	ctx := new(serve.Context)
// 	ctx.Server = serve.NewServer("", "/serve", fs)
// 	ctx.URL = "/kiruba/app/module/path/to/file"
//
// 	fs.Build(ctx, ctx.URL)
//
// 	if !(ctx.Namespace.Name == "." && ctx.Application == nil && ctx.Module.Name == "home" && ctx.Path == "/kiruba/app/module/path/to/file") {
// 		t.Fail()
// 	}
// }
//
// func TestBuildApp(t *testing.T) {
// 	fs := driver.NewFileSystem(func(path string) bool {
// 		//fmt.Println("stat", path)
// 		if path == "/serve" || path == "/serve/apps/app" || path == "/serve/modules/module" {
// 			return true
// 		}
// 		return false
// 	})
//
// 	ctx := new(serve.Context)
// 	ctx.Server = serve.NewServer("", "/serve", fs)
// 	ctx.URL = "/app/module/path/to/file"
//
// 	fs.Build(ctx, ctx.URL)
//
// 	fmt.Println(ctx.Path)
//
// 	if !(ctx.Namespace.Name == "." && ctx.Application.Name == "app" && ctx.Module.Name == "module" && ctx.Path == "/path/to/file") {
// 		t.Fail()
// 	}
// }
//
// func TestBuildDefaultApp(t *testing.T) {
// 	fs := driver.NewFileSystem(func(path string) bool {
// 		if path == "/serve" || path == "/serve/modules/home" {
// 			return true
// 		}
// 		return false
// 	})
//
// 	ctx := new(serve.Context)
// 	ctx.Server = serve.NewServer("", "/serve", fs)
// 	ctx.URL = "//path/to/file"
//
// 	fs.Build(ctx, ctx.URL)
//
// 	if !(ctx.Namespace.Name == "." && ctx.Application == nil && ctx.Module.Name == "home" && ctx.Path == "/path/to/file") {
// 		t.Fail()
// 	}
// }

func TestBuildRedirectApp(t *testing.T) {
	getConfig := func(path string) *[]byte {
		if path == "/serve/apps/app" {
			ba := []byte("{ \"modules\" : [\"home\"] }")
			return &ba
		}
		return nil
	}

	stat := func(path string) bool {
		if path == "/serve/apps/app" || path == "/ctx *serve.Context, w http.ResponseWriter, r *http.Requestrve/modules/home" {
			return true
		}
		return false
	}

	serveFile := func(ctx *serve.Context, w http.ResponseWriter, r *http.Request) {

	}

	fs := driver.NewFileSystem(stat, getConfig, serveFile)
	ctx := new(serve.Context)
	ctx.Server = serve.NewServer("", "/serve", fs)
	ctx.URL = "/app"

	fs.Build(ctx, ctx.URL)

	//fmt.Println(ctx.Namespace, ctx.Application, ctx.Module, ctx.Path)

	if !(ctx.Namespace.Name == "." && ctx.Application.Name == "app" && ctx.Module.Name == "home" && ctx.Path == "") {
		t.Fail()
	}
}
