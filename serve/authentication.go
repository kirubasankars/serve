package serve

// Authentication forms authentication
type Authentication struct {
	server *Server
}

// Authenticate authenticate with username and password
func (auth *Authentication) Authenticate(param map[string]string) int {
	//var count int

	// server := auth.server
	// system := server.System
	// client := system.GetClient(server, param.clientID)

	// if client == nil {
	// 	return 0
	// }

	// if param.redirectURI != "" {
	// 	r, _ := http.NewRequest("GET", param.redirectURI, nil)
	// 	redirectURICtx := newContext(server, r)
	//
	// 	if redirectURICtx != nil && redirectURICtx.Namespace != nil {
	// 		for _, ns := range client.Namespaces {
	// 			if ns == redirectURICtx.Namespace.Name {
	// 				count++
	// 				break
	// 			}
	// 		}
	// 	}
	// }

	// if param.refreshToken == "refresh_token" {
	// 	count++
	// }
	// if param.authorizationCode == "code" {
	// 	count++
	// }
	// if param.username == "admin" && param.password == "admin" {
	// 	count++
	// }
	//
	// if param.clientID == client.ID && param.clientSecret == client.Secret {
	// 	count++
	// }

	return 1
}
