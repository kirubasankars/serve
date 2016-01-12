package serve

import "github.com/kirubasankars/metal"

func (site *Site) Model(path string, method string, input *metal.Metal) *metal.Metal {
	server := site.server
	return server.IO.API(site, path, method, input)
}
