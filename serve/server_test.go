package serve_test

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"path/filepath"
	
	"github.com/kirubasankars/serve/driver"
	"github.com/kirubasankars/serve/metal"
	"github.com/kirubasankars/serve/serve"
)

type CommonSiteHandler struct{}

func (csh *CommonSiteHandler) Build(module serve.Module) {
	module.Handlers["/"] = func(ctx serve.Context, w http.ResponseWriter, r *http.Request) {

		n := ctx.Namespace
		a := ctx.Application
		m := ctx.Module

		var (
			namespace string
			app       string
			module    string
		)

		if n != nil {
			namespace = n.Name
		}
		if a != nil {
			app = a.Name
		}
		if m != nil {
			module = m.Name
		}

		//fmt.Println(ctx.Namespace, ctx.Application, app, ctx.Module, ctx.Path)

		fmt.Fprintf(w, "%s %s %s", namespace, app, module)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func TestServeHttp(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/app/module/path/to/file", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	getConfig := func(path string) *metal.Metal {
		m := metal.NewMetal()
		m.Set("modules.@0", "module")
		return m
	}
	stat := func(path string) bool {		
		if path == filepath.FromSlash("/serve") || path == filepath.FromSlash("/serve/apps/app") || path == filepath.FromSlash("/serve/modules/module") {
			return true
		}
		return false
	}

	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))
	server.ServeHTTP(w, req)

	//fmt.Printf("%d - %s", w.Code, w.Body.String())

	if !(w.Code == 200 || strings.TrimSpace(w.Body.String()) != ". app module") {
		t.Error("return code is not 200")
	}
}

func TestServeHttpModuleRootRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/module", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	getConfig := func(path string) *metal.Metal {
		m := metal.NewMetal()
		m.Set("modules.@0", "module")
		return m
	}
	stat := func(path string) bool {		
		if path == filepath.FromSlash("/serve") || path == filepath.FromSlash("/serve/modules/module") {
			return true
		}
		return false
	}

	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.ServeHTTP(w, req)

	if !(w.Code == 301 && strings.TrimSpace(w.Body.String()) == "<a href=\"/module/\">Moved Permanently</a>.") {
		fmt.Printf("%d - %s", w.Code, w.Body.String())
		t.Error("return code is not 301")
	}
}

func TestServeHttpAppModuleRootRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/app/module", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	getConfig := func(path string) *metal.Metal {
		m := metal.NewMetal()
		m.Set("modules.@0", "module")
		return m
	}
	stat := func(path string) bool {
		if path == filepath.FromSlash("/serve") || path == filepath.FromSlash("/serve/apps/app") || path == filepath.FromSlash("/serve/modules/module") {
			return true
		}
		return false
	}
	d := driver.NewFileSystem(stat, getConfig)

	server := serve.NewServer("3000", "/serve", d)
	server.ServeHTTP(w, req)

	//fmt.Printf("%d - %s", w.Code, w.Body.String())

	if !(w.Code == 301 || strings.TrimSpace(w.Body.String()) == "<a href=\"/app/module/\">Moved Permanently</a>.") {
		t.Error("return code is not 301")
	}
}

func TestServeHttpAppModuleRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/app/module/path/2/file", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	stat := func(path string) bool {
		if path == filepath.FromSlash("/serve") || path == filepath.FromSlash("/serve/apps/app") || path == filepath.FromSlash("/serve/modules/module") {
			return true
		}
		return false
	}
	getConfig := func(path string) *metal.Metal {
		m := metal.NewMetal()
		m.Set("modules.@0", "module")
		return m
	}

	d := driver.NewFileSystem(stat, getConfig)

	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))

	server.ServeHTTP(w, req)

	//fmt.Printf("%d - %s", w.Code, w.Body)

	if !(w.Code == 200 && w.Body.String() == ". app module") {
		log.Printf("%d - %s", w.Code, w.Body)
		t.Error("return code is not 200")
	}
}

func TestServeHttpModuleRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/module/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	getConfig := func(path string) *metal.Metal {
		m := metal.NewMetal()
		m.Set("modules.@0", "module")
		return m
	}
	stat := func(path string) bool {
		if path == filepath.FromSlash("/serve") || path == filepath.FromSlash("/serve/modules/module") {
			return true
		}
		return false
	}
	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))
	server.ServeHTTP(w, req)

	//fmt.Printf("%d - %s", w.Code, w.Body.String())

	if !(w.Code == 200 || w.Body.String() == ".  module") {
		t.Logf("%d - %s", w.Code, w.Body.String())
		t.Error("return code is not 301")
	}
}

func TestServeHttpApp(t *testing.T) {

	req, err := http.NewRequest("GET", "http://localhost:3000/app/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	getConfig := func(path string) *metal.Metal {
		return nil
	}
	stat := func(path string) bool {
		if path == filepath.FromSlash("/serve") || path == filepath.FromSlash("/serve/apps/app") || path == filepath.FromSlash("/serve/modules/home") {
			return true
		}
		return false
	}
	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))
	server.ServeHTTP(w, req)

	//fmt.Printf("%d - %s", w.Code, w.Body.String())

	if !(w.Code == 200 && w.Body.String() == ". app home") {
		t.Logf("%d - %s", w.Code, w.Body.String())
		t.Error("return code is not 301")
	}
}

func TestServeHttpNamespaceAppNamespaceModuleRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/namespace/app/module/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	stat := func(path string) bool {
		//fmt.Println(path)
		if path == filepath.FromSlash("/serve/namespace") || path == filepath.FromSlash("/serve/namespace/apps/app") || path == filepath.FromSlash("/serve/namespace/modules/module") {
			return true
		}
		return false
	}
	getConfig := func(path string) *metal.Metal {
		if path == filepath.FromSlash("/serve/namespace") {
			m := metal.NewMetal()
			m.Set("modules.@0", "module")
			return m
		}
		return nil
	}
	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))

	server.ServeHTTP(w, req)

	//fmt.Printf("%d - %s", w.Code, w.Body.String())

	if w.Code != 200 || w.Body.String() != "namespace app module" {
		t.Logf("%d - %s", w.Code, w.Body.String())
		t.Error("return code is not 301")
	}
}

func TestServeHttpNamespaceModuleRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/namespace/module/", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	getConfig := func(path string) *metal.Metal {
		if path == filepath.FromSlash("/serve/namespace") {
			m := metal.NewMetal()
			m.Set("modules.@0", "module")
			m.Set("roles.adminstrator.@0", "check:admin")
			return m
		}
		if path == filepath.FromSlash("/serve/modules/module") {
			m := metal.NewMetal()
			m.Set("modules.@0", "module")
			m.Set("permissions.admin.@0", "url(GET:path.to.file)")
			return m
		}
		return nil
	}

	stat := func(path string) bool {
		if path == filepath.FromSlash("/serve/namespace") || path == filepath.FromSlash("/serve/modules/module") {
			return true
		}
		return false
	}

	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))
	server.ServeHTTP(w, req)

	//fmt.Printf("%d - %s", w.Code, w.Body.String())

	if w.Code != 200 || w.Body.String() != "namespace  module" {
		t.Logf("%d - %s", w.Code, w.Body.String())
		t.Error("return code is not 301")
	}
}

func TestServeHttpNamespcaeAppModuleRootRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/namespace/app/module", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	stat := func(path string) bool {
		if path == filepath.FromSlash("/serve/namespace") || path == filepath.FromSlash("/serve/namespace/apps/app") || path == filepath.FromSlash("/serve/modules/module") {
			return true
		}
		return false
	}

	getConfig := func(path string) *metal.Metal {
		if path == filepath.FromSlash("/serve/namespace") {
			m := metal.NewMetal()
			m.Set("modules.@0", "module")
			return m
		}
		return nil
	}

	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.ServeHTTP(w, req)

	//fmt.Printf("%d - %s", w.Code, w.Body.String())

	if !(w.Code == 301 && strings.TrimSpace(w.Body.String()) == "<a href=\"/namespace/app/module/\">Moved Permanently</a>.") {
		log.Printf("%d - %s", w.Code, w.Body.String())
		t.Error("return code is not 301")
	}
}
