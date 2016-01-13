package serve

import (
	"net/http"

	"github.com/kirubasankars/metal"
)

type ServeHandler interface {
	IsSitePath(path string) bool
	IsExists(path string) bool
	ServeFile(w http.ResponseWriter, r *http.Request, path string)

	Template(path string) *[]byte
	API(site *Site, method string, path string, input *metal.Metal) *metal.Metal
}
