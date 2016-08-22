package serve

// Authentication forms authentication
type Authentication struct{}

// Authenticate authenticate with username and password
func (auth *Authentication) Authenticate(username string, password string, clientID string, clientSecret string) bool {
	if username == "admin" && password == "admin" && clientID == "client_id" && clientSecret == "client_secret" {
		return true
	}
	return false
}
