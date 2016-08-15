package serve

import "github.com/kirubasankars/serve/metal"

// Application struct
type Application struct {
	Name   string
	Path   string
	config *metal.Metal

	server *Server
}

// GetConfig get config from application
func (app *Application) GetConfig(key string) interface{} {
	if app.config == nil {
		return nil
	}
	return app.config.Get(key)
}

// Build build application
func (app *Application) Build() {

}

// NewApplication create application
func NewApplication(name string, path string, config *metal.Metal, server *Server) *Application {
	app := new(Application)
	app.Name = name
	app.Path = path
	app.config = config

	app.server = server

	return app
}
