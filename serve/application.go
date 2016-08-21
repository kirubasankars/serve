package serve

import "encoding/json"

// Application struct
type Application struct {
	Name   string
	Path   string
	Config *ApplicationConfigration

	roles map[string][]string

	server *Server
}

// GetConfig get config from application
// func (app *Application) GetConfig(key string) interface{} {
// 	if app.config == nil {
// 		return nil
// 	}
// 	return app.config.Get(key)
// }

// Build build application
func (app *Application) Build() {
	// if roles, _ := app.GetConfig("roles").(*metal.Metal); roles != nil {
	// 	props := roles.Properties()
	// 	if app.roles == nil {
	// 		app.roles = make(map[string][]string, len(props))
	// 	}
	// 	appRoles := app.roles
	// 	for name := range props {
	// 		if role, done := roles.Get(name).(*metal.Metal); done == true {
	// 			for _, v := range role.Properties() {
	// 				if auth, done := v.(string); done == true {
	// 					if _, p := appRoles[name]; p == true {
	// 						appRoles[name] = make([]string, 0)
	// 					}
	// 					appRoles[name] = append(appRoles[name], auth)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
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
