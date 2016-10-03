package serve

import "encoding/json"

// Application struct
type Application struct {
	Name   string
	Path   string
	Config *ApplicationConfigration

	server *Server
}

// Build build application
func (app *Application) Build() {

}

// NewApplication create application
func NewApplication(name string, path string, config *[]byte, server *Server) *Application {
	app := new(Application)
	app.Name = name
	app.Path = path

	if config != nil {
		ac := new(ApplicationConfigration)
		if err := json.Unmarshal(*config, &ac); err != nil {
			panic(err)
		}
		app.Config = ac
	}

	app.server = server

	return app
}
