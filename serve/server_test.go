package serve_test

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kirubasankars/serve/driver"
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
	req, err := http.NewRequest("GET", "http://localhost:3000/path/to/file", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/config.json") {
			b := []byte("{ \"roles\" : { \"admin\" : [\"home:permission\"] } }")
			return &b
		}
		if path == filepath.FromSlash("/serve/modules/home/config.json") {
			b := []byte("{ \"permissions\" : { \"permission\" : [\"admin\",\"url(GET /path/to/file)\"] } }")
			return &b
		}
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

	if !(w.Code == 200 && strings.TrimSpace(w.Body.String()) == ". . home") {
		t.Error("return code is not 200")
	}
}

func TestServeHttpModuleRootRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:3000/module", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/config.json") {
			b := []byte("{ \"modules\" : [\"module\"], \"roles\" : { \"admin\" : [\"module:permission\"] } }")
			return &b
		}
		if path == filepath.FromSlash("/serve/modules/module/config.json") {
			b := []byte("{ \"permissions\" : { \"permission\" : [\"url(GET /?)\"] } }")
			return &b
		}
		return nil
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

	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/apps/app/config.json") {
			ba := []byte("{ \"modules\" : [\"module\"], \"roles\" : { \"admin\" : [\"module:permission\"] } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/modules/module/config.json") {
			ba := []byte("{ \"permissions\" : { \"permission\" : [\"url(GET /?)\"] } }")
			return &ba
		}
		return nil
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
	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/apps/app/config.json") {
			ba := []byte("{ \"modules\" : [\"module\"], \"roles\" : { \"admin\" : [\"module:permission\"] } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/modules/module/config.json") {
			ba := []byte("{ \"permissions\" : { \"permission\" : [\"url(GET /path/2/file)\"] } }")
			return &ba
		}
		return nil
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
	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/config.json") {
			ba := []byte("{ \"modules\" : [\"module\"], \"roles\" : { \"admin\" : [\"module:permission\"] } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/modules/module/config.json") {
			ba := []byte("{ \"permissions\" : { \"permission\" : [\"url(GET /)\"] } }")
			return &ba
		}
		return nil
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
	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/apps/app/config.json") {
			ba := []byte("{ \"roles\" : { \"admin\" : [\"home:permission\"] } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/modules/home/config.json") {
			ba := []byte("{ \"permissions\" : { \"permission\" : [\"url(GET /)\"] } }")
			return &ba
		}
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
		if path == filepath.FromSlash("/serve/namespace") || path == filepath.FromSlash("/serve/namespace/apps/app") || path == filepath.FromSlash("/serve/namespace/modules/module") {
			return true
		}
		return false
	}
	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/namespace/apps/app/config.json") {
			ba := []byte("{ \"modules\" : [\"module\"], \"roles\" : { \"admin\" : [\"module:permission\"] } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/namespace/modules/module/config.json") {
			ba := []byte("{ \"permissions\" : { \"permission\" : [\"url(GET /)\"] } }")
			return &ba
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

	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/namespace/config.json") {
			ba := []byte("{ \"modules\" : [\"module\"], \"roles\" : { \"admin\" : [\"module:admin\"] } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/namespace/modules/module/config.json") {
			ba := []byte("{ \"permissions\" : { \"admin\" : [\"url(GET /?)\"] } }")
			return &ba
		}
		return nil
	}

	stat := func(path string) bool {
		if path == filepath.FromSlash("/serve/namespace") || path == filepath.FromSlash("/serve/namespace/modules/module") {
			return true
		}
		return false
	}

	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))
	server.ServeHTTP(w, req)

	//fmt.Printf("%d - %s", w.Code, w.Body.String())

	if w.Code != 200 || w.Body.String() != "namespace . module" {
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
		if path == filepath.FromSlash("/serve/namespace") || path == filepath.FromSlash("/serve/namespace/apps/app") || path == filepath.FromSlash("/serve/namespace/modules/module") {
			return true
		}
		return false
	}

	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/namespace/apps/app/config.json") {
			ba := []byte("{ \"modules\" : [\"module\"], \"roles\" : { \"admin\" : [\"module:admin\"] } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/namespace/modules/module/config.json") {
			ba := []byte("{ \"permissions\" : { \"admin\" : [\"url(GET /?)\"] } }")
			return &ba
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

func TestServeHttpOAuthUserPassword(t *testing.T) {
	query := "?grant_type=password&client_id=client_id&client_secret=client_secret&username=admin&password=admin"
	req, err := http.NewRequest("GET", "http://localhost:3000/oauth2/token"+query, nil)
	req.Method = "POST"
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/users.json") {
			ba := []byte("{ \"admin\" : { \"id\" : \"admin\", \"password\" : \"admin\", \"namespaces\" : { \"namespace\" : { \"apps\" : { \"app\" : { \"roles\" : [ \"\"] } } } } } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/clients.json") {
			ba := []byte("{ \"client_id\" : { \"id\" : \"client_id\", \"secret\" : \"client_secret\" } }")
			return &ba
		}
		return nil
	}
	stat := func(path string) bool {
		return false
	}

	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))
	server.ServeHTTP(w, req)

	if !(w.Code == 200 && strings.TrimSpace(w.Body.String()) == "{\"access_token\":\"access_token\",\"issued_at\":\"issued_at\",\"signature\":\"signature\"}") {
		t.Error("return code is not 200")
	}
}

func TestServeHttpOAuthUserAgent(t *testing.T) {
	query := "?response_code=token&client_id=client_id&redirect_uri=namespace/app&username=admin&password=admin"
	req, err := http.NewRequest("GET", "http://localhost:3000/oauth2/authorize"+query, nil)
	req.Method = "POST"
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/users.json") {
			ba := []byte("{ \"admin\" : { \"id\" : \"admin\", \"password\" : \"admin\", \"roles\" : { \"namespace:app\" : [ \"\"]  } } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/clients.json") {
			ba := []byte("{ \"client_id\" : { \"id\" : \"client_id\", \"secret\" : \"client_secret\" } }")
			return &ba
		}
		return nil
	}
	stat := func(path string) bool {
		if path == filepath.FromSlash("/serve/namespace") {
			return true
		}
		if path == filepath.FromSlash("/serve/namespace/apps/app") {
			return true
		}
		return false
	}

	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))
	server.ServeHTTP(w, req)

	if !(w.Code == 302 && w.Header().Get("Location") == "/namespace/app?access_token=access_token&issued_at=issuedAt") {
		t.Error("return code is not 302")
	}
}

func TestServeHttpOAuthWebServer(t *testing.T) {
	query := "?response_code=code&client_id=client_id&redirect_uri=namespace/app&username=admin&password=admin"
	req, err := http.NewRequest("GET", "http://localhost:3000/oauth2/authorize"+query, nil)
	req.Method = "POST"
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	getConfig := func(path string) *[]byte {
		if path == filepath.FromSlash("/serve/users.json") {
			ba := []byte("{ \"admin\" : { \"id\" : \"admin\", \"password\" : \"admin\", \"roles\" : { \"namespace:app\" : [ \"\"]  } } }")
			return &ba
		}
		if path == filepath.FromSlash("/serve/clients.json") {
			ba := []byte("{ \"client_id\" : { \"id\" : \"client_id\", \"secret\" : \"client_secret\" } }")
			return &ba
		}
		return nil
	}
	stat := func(path string) bool {
		if path == filepath.FromSlash("/serve/namespace") {
			return true
		}
		if path == filepath.FromSlash("/serve/namespace/apps/app") {
			return true
		}
		return false
	}

	d := driver.NewFileSystem(stat, getConfig)
	server := serve.NewServer("3000", "/serve", d)
	server.RegisterProvider(".", new(CommonSiteHandler))
	server.ServeHTTP(w, req)

	if !(w.Code == 302 && w.Header().Get("Location") == "/oauth/code_callback?code=12345678&redirect_uri=namespace/app") {
		t.Error("return code is not 302")
	}
}
